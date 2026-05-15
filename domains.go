package smail

import (
	"context"
	"fmt"
)

// DomainService manages custom sending domains.
type DomainService struct{ c *Client }

// Add registers a new custom domain. The response includes the DNS records
// you must configure before calling Verify.
func (s *DomainService) Add(ctx context.Context, req AddDomainRequest) (DomainResponse, error) {
	var resp DomainResponse
	if err := s.c.do(ctx, "POST", "/api/v1/domains", req, &resp); err != nil {
		return DomainResponse{}, err
	}
	return resp, nil
}

// List returns all sending domains for the authenticated developer.
func (s *DomainService) List(ctx context.Context) ([]SendingDomain, error) {
	var env domainListEnvelope
	if err := s.c.do(ctx, "GET", "/api/v1/domains", nil, &env); err != nil {
		return nil, err
	}
	return env.Domains, nil
}

// Get returns a single domain by ID, including its required DNS records.
func (s *DomainService) Get(ctx context.Context, domainID string) (DomainResponse, error) {
	var resp DomainResponse
	if err := s.c.do(ctx, "GET", fmt.Sprintf("/api/v1/domains/%s", domainID), nil, &resp); err != nil {
		return DomainResponse{}, err
	}
	return resp, nil
}

// Verify triggers a live DNS check to confirm the domain's TXT records are in
// place. Returns the updated domain on success, or an error with detail if the
// DNS records are not yet propagated.
func (s *DomainService) Verify(ctx context.Context, domainID string) (DomainResponse, error) {
	var resp DomainResponse
	if err := s.c.do(ctx, "POST", fmt.Sprintf("/api/v1/domains/%s/verify", domainID), nil, &resp); err != nil {
		return DomainResponse{}, err
	}
	return resp, nil
}

// Delete removes the domain from the account. Emails can no longer be sent
// from addresses on this domain after deletion.
func (s *DomainService) Delete(ctx context.Context, domainID string) error {
	return s.c.do(ctx, "DELETE", fmt.Sprintf("/api/v1/domains/%s", domainID), nil, nil)
}
