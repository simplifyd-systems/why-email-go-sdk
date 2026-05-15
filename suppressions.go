package smail

import (
	"context"
	"fmt"
	"net/url"
)

// SuppressionService manages account-level suppressed recipients.
type SuppressionService struct{ c *Client }

func (s *SuppressionService) Create(ctx context.Context, req CreateSuppressionRequest) (Suppression, error) {
	var env suppressionEnvelope
	if err := s.c.do(ctx, "POST", "/api/v1/suppressions", req, &env); err != nil {
		return Suppression{}, err
	}
	return env.Suppression, nil
}

func (s *SuppressionService) List(ctx context.Context, opts ListSuppressionsOptions) ([]Suppression, error) {
	limit := opts.Limit
	if limit <= 0 {
		limit = 50
	}
	offset := opts.Offset

	path := fmt.Sprintf("/api/v1/suppressions?%s", url.Values{
		"limit":  []string{fmt.Sprintf("%d", limit)},
		"offset": []string{fmt.Sprintf("%d", offset)},
	}.Encode())

	var env suppressionListEnvelope
	if err := s.c.do(ctx, "GET", path, nil, &env); err != nil {
		return nil, err
	}
	return env.Suppressions, nil
}

func (s *SuppressionService) Delete(ctx context.Context, email string) error {
	return s.c.do(ctx, "DELETE", "/api/v1/suppressions/"+url.PathEscape(email), nil, nil)
}
