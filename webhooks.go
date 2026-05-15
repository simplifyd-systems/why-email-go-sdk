package smail

import (
	"context"
	"fmt"
	"net/url"
)

// WebhookService manages webhook endpoints and delivery attempts.
type WebhookService struct{ c *Client }

// Create registers a webhook endpoint. The returned Secret is shown only once.
func (s *WebhookService) Create(ctx context.Context, req CreateWebhookRequest) (WebhookEndpoint, error) {
	var env webhookEnvelope
	if err := s.c.do(ctx, "POST", "/api/v1/webhooks", req, &env); err != nil {
		return WebhookEndpoint{}, err
	}
	return env.Webhook, nil
}

// List returns configured webhook endpoints for the authenticated developer.
func (s *WebhookService) List(ctx context.Context) ([]WebhookEndpoint, error) {
	var env webhookListEnvelope
	if err := s.c.do(ctx, "GET", "/api/v1/webhooks", nil, &env); err != nil {
		return nil, err
	}
	return env.Webhooks, nil
}

// Delete removes a webhook endpoint and stops future deliveries.
func (s *WebhookService) Delete(ctx context.Context, webhookID string) error {
	return s.c.do(ctx, "DELETE", fmt.Sprintf("/api/v1/webhooks/%s", webhookID), nil, nil)
}

// Deliveries returns recent webhook delivery attempts.
func (s *WebhookService) Deliveries(ctx context.Context, opts ListWebhookDeliveriesOptions) ([]WebhookDelivery, error) {
	limit := opts.Limit
	if limit <= 0 {
		limit = 50
	}
	offset := opts.Offset

	path := fmt.Sprintf("/api/v1/webhook-deliveries?%s", url.Values{
		"limit":  []string{fmt.Sprintf("%d", limit)},
		"offset": []string{fmt.Sprintf("%d", offset)},
	}.Encode())

	var env webhookDeliveryListEnvelope
	if err := s.c.do(ctx, "GET", path, nil, &env); err != nil {
		return nil, err
	}
	return env.Deliveries, nil
}
