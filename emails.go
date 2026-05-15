package smail

import (
	"context"
	"fmt"
	"net/url"
)

// EmailService handles transactional and marketing email sending.
type EmailService struct{ c *Client }

// Send delivers a single email. The from address must belong to a verified
// domain on the account. Returns the queued email record.
func (s *EmailService) Send(ctx context.Context, req SendEmailRequest) (EmailSend, error) {
	return s.SendWithOptions(ctx, req, SendOptions{})
}

// SendWithOptions delivers a single email with optional request controls.
func (s *EmailService) SendWithOptions(ctx context.Context, req SendEmailRequest, opts SendOptions) (EmailSend, error) {
	var env emailSendEnvelope
	if err := s.c.doWithHeaders(ctx, "POST", "/api/v1/emails", req, &env, idempotencyHeaders(opts)); err != nil {
		return EmailSend{}, err
	}
	return env.Email, nil
}

// Batch sends up to 1000 emails in a single API call. Each email is processed
// independently; failures for individual emails do not abort the batch.
func (s *EmailService) Batch(ctx context.Context, emails []SendEmailRequest) (BatchSendResponse, error) {
	return s.BatchWithOptions(ctx, emails, SendOptions{})
}

// BatchWithOptions sends up to 1000 emails with optional request controls.
func (s *EmailService) BatchWithOptions(ctx context.Context, emails []SendEmailRequest, opts SendOptions) (BatchSendResponse, error) {
	var resp BatchSendResponse
	if err := s.c.doWithHeaders(ctx, "POST", "/api/v1/emails/batch", BatchSendRequest{Emails: emails}, &resp, idempotencyHeaders(opts)); err != nil {
		return BatchSendResponse{}, err
	}
	return resp, nil
}

// List returns a paginated list of sent emails for the authenticated developer.
func (s *EmailService) List(ctx context.Context, opts ListEmailsOptions) (ListEmailsResponse, error) {
	limit := opts.Limit
	if limit <= 0 {
		limit = 50
	}
	offset := opts.Offset

	path := fmt.Sprintf("/api/v1/emails?%s", url.Values{
		"limit":  []string{fmt.Sprintf("%d", limit)},
		"offset": []string{fmt.Sprintf("%d", offset)},
	}.Encode())

	var resp ListEmailsResponse
	if err := s.c.do(ctx, "GET", path, nil, &resp); err != nil {
		return ListEmailsResponse{}, err
	}
	return resp, nil
}

// Get returns the status and details of a single email send by ID.
func (s *EmailService) Get(ctx context.Context, emailID string) (EmailSend, error) {
	var env emailSendEnvelope
	if err := s.c.do(ctx, "GET", fmt.Sprintf("/api/v1/emails/%s", emailID), nil, &env); err != nil {
		return EmailSend{}, err
	}
	return env.Email, nil
}

func idempotencyHeaders(opts SendOptions) map[string]string {
	if opts.IdempotencyKey == "" {
		return nil
	}
	return map[string]string{"Idempotency-Key": opts.IdempotencyKey}
}
