// Package smail provides a Go client for the Smail developer email API.
//
// Usage with an API key (sending emails):
//
//	client := smail.New("https://mail.example.com", smail.WithAPIKey("sk_live_..."))
//	resp, err := client.Emails.Send(ctx, smail.SendEmailRequest{
//	    From:    "hello@yourdomain.com",
//	    To:      []string{"user@example.com"},
//	    Subject: "Hello",
//	    Text:    "Hello from Smail!",
//	})
//
// Usage with a developer JWT (account management):
//
//	client := smail.New("https://mail.example.com", smail.WithToken("eyJ..."))
//	me, err := client.Auth.Me(ctx)
package smail

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Client is the root Smail API client. Construct it with New().
type Client struct {
	baseURL    string
	token      string // Bearer JWT
	apiKey     string // X-API-Key
	httpClient *http.Client

	Auth         *AuthService
	Emails       *EmailService
	Keys         *KeyService
	Domains      *DomainService
	Webhooks     *WebhookService
	Suppressions *SuppressionService
	Templates    *TemplateService
}

// Option configures a Client.
type Option func(*Client)

// WithToken sets a developer JWT for authenticated management endpoints.
func WithToken(token string) Option {
	return func(c *Client) { c.token = token }
}

// WithAPIKey sets an API key used via the X-API-Key header for send endpoints.
func WithAPIKey(key string) Option {
	return func(c *Client) { c.apiKey = key }
}

// WithHTTPClient replaces the default HTTP client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.httpClient = hc }
}

// New creates a Smail API client targeting baseURL.
func New(baseURL string, opts ...Option) *Client {
	c := &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{},
	}
	for _, o := range opts {
		o(c)
	}
	c.Auth = &AuthService{c}
	c.Emails = &EmailService{c}
	c.Keys = &KeyService{c}
	c.Domains = &DomainService{c}
	c.Webhooks = &WebhookService{c}
	c.Suppressions = &SuppressionService{c}
	c.Templates = &TemplateService{c}
	return c
}

// SetToken updates the Bearer JWT on an existing client (e.g. after login).
func (c *Client) SetToken(token string) { c.token = token }

// SetAPIKey updates the API key on an existing client.
func (c *Client) SetAPIKey(key string) { c.apiKey = key }

// do executes an HTTP request and decodes the JSON response into result.
// body may be nil for requests without a payload.
// result may be nil if the response body should be discarded.
func (c *Client) do(ctx context.Context, method, path string, body, result interface{}) error {
	return c.doWithHeaders(ctx, method, path, body, result, nil)
}

func (c *Client) doWithHeaders(ctx context.Context, method, path string, body, result interface{}, headers map[string]string) error {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("smail: marshal request: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return fmt.Errorf("smail: create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// Auth: API key takes precedence over JWT.
	if c.apiKey != "" {
		req.Header.Set("X-API-Key", c.apiKey)
	} else if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("smail: http: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("smail: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var envelope apiErrorEnvelope
		_ = json.Unmarshal(respBytes, &envelope)
		msg := envelope.Error
		if msg == "" {
			msg = http.StatusText(resp.StatusCode)
		}
		return &APIError{StatusCode: resp.StatusCode, Message: msg}
	}

	if result != nil && len(respBytes) > 0 {
		if err := json.Unmarshal(respBytes, result); err != nil {
			return fmt.Errorf("smail: decode response: %w", err)
		}
	}

	return nil
}
