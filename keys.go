package smail

import (
	"context"
	"fmt"
)

// KeyService manages developer API keys.
type KeyService struct{ c *Client }

// Create generates a new API key. The full key value is returned only in this
// response — store it immediately, as it cannot be retrieved again.
func (s *KeyService) Create(ctx context.Context, req CreateAPIKeyRequest) (CreateAPIKeyResponse, error) {
	var resp CreateAPIKeyResponse
	if err := s.c.do(ctx, "POST", "/api/v1/keys", req, &resp); err != nil {
		return CreateAPIKeyResponse{}, err
	}
	return resp, nil
}

// List returns all API keys for the authenticated developer. The full key value
// is never included; only the prefix is shown.
func (s *KeyService) List(ctx context.Context) ([]APIKey, error) {
	var env apiKeyListEnvelope
	if err := s.c.do(ctx, "GET", "/api/v1/keys", nil, &env); err != nil {
		return nil, err
	}
	return env.APIKeys, nil
}

// Revoke permanently revokes the API key with the given ID. Revoked keys cannot
// be re-activated.
func (s *KeyService) Revoke(ctx context.Context, keyID string) error {
	return s.c.do(ctx, "DELETE", fmt.Sprintf("/api/v1/keys/%s", keyID), nil, nil)
}
