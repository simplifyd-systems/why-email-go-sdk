package smail

import "encoding/json"

// ── Request types ─────────────────────────────────────────────────────────────

type RegisterRequest struct {
	Email    string  `json:"email"`
	Password string  `json:"password"`
	FullName string  `json:"full_name"`
	Company  *string `json:"company,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type VerifyEmailRequest struct {
	Token string `json:"token"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

type UpdateProfileRequest struct {
	FullName *string `json:"full_name,omitempty"`
	Company  *string `json:"company,omitempty"`
}

type CreateAPIKeyRequest struct {
	Name   string   `json:"name"`
	Scopes []string `json:"scopes,omitempty"`
}

type AddDomainRequest struct {
	Domain string `json:"domain"`
}

type WebhookEvent string

const (
	WebhookEventEmailQueued     WebhookEvent = "email.queued"
	WebhookEventEmailDelivered  WebhookEvent = "email.delivered"
	WebhookEventEmailFailed     WebhookEvent = "email.failed"
	WebhookEventEmailSuppressed WebhookEvent = "email.suppressed"
)

type CreateWebhookRequest struct {
	URL    string         `json:"url"`
	Events []WebhookEvent `json:"events"`
}

type CreateSuppressionRequest struct {
	Email  string `json:"email"`
	Reason string `json:"reason,omitempty"`
	Note   string `json:"note,omitempty"`
}

type TemplateVariable struct {
	Key           string      `json:"key"`
	Type          string      `json:"type,omitempty"`
	FallbackValue interface{} `json:"fallback_value,omitempty"`
}

type CreateTemplateRequest struct {
	Name      string             `json:"name"`
	Alias     string             `json:"alias,omitempty"`
	From      string             `json:"from,omitempty"`
	Subject   string             `json:"subject,omitempty"`
	ReplyTo   string             `json:"reply_to,omitempty"`
	HTML      string             `json:"html"`
	Text      string             `json:"text,omitempty"`
	Variables []TemplateVariable `json:"variables,omitempty"`
}

type UpdateTemplateRequest struct {
	Name      string             `json:"name,omitempty"`
	Alias     string             `json:"alias,omitempty"`
	From      string             `json:"from,omitempty"`
	Subject   string             `json:"subject,omitempty"`
	ReplyTo   string             `json:"reply_to,omitempty"`
	HTML      string             `json:"html,omitempty"`
	Text      string             `json:"text,omitempty"`
	Variables []TemplateVariable `json:"variables,omitempty"`
}

type SendEmailRequest struct {
	From        string            `json:"from"`
	To          []string          `json:"to"`
	Cc          []string          `json:"cc,omitempty"`
	Bcc         []string          `json:"bcc,omitempty"`
	ReplyTo     string            `json:"reply_to,omitempty"`
	Subject     string            `json:"subject"`
	HTML        string            `json:"html,omitempty"`
	Text        string            `json:"text,omitempty"`
	Template    *SendTemplate     `json:"template,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
	Attachments []EmailAttachment `json:"attachments,omitempty"`
	// ScheduledAt is an RFC3339 timestamp. When empty, the email sends immediately.
	ScheduledAt string `json:"scheduled_at,omitempty"`
	// Type is "transactional" (default) or "marketing".
	Type string `json:"type,omitempty"`
}

type SendTemplate struct {
	ID        string                 `json:"id,omitempty"`
	Alias     string                 `json:"alias,omitempty"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type EmailAttachment struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type,omitempty"`
	// Content is base64-encoded file content.
	Content string `json:"content"`
}

type BatchSendRequest struct {
	Emails []SendEmailRequest `json:"emails"`
}

type ListEmailsOptions struct {
	Limit  int
	Offset int
}

type SendOptions struct {
	// IdempotencyKey deduplicates retries for 24 hours. Maximum 256 characters.
	IdempotencyKey string
}

type ListWebhookDeliveriesOptions struct {
	Limit  int
	Offset int
}

type ListSuppressionsOptions struct {
	Limit  int
	Offset int
}

type ListTemplatesOptions struct {
	Limit  int
	Offset int
}

// ── Response types ────────────────────────────────────────────────────────────

type Developer struct {
	ID        string  `json:"id"`
	Email     string  `json:"email"`
	FullName  string  `json:"full_name"`
	Company   *string `json:"company"`
	Plan      string  `json:"plan"`
	Status    string  `json:"status"`
	CreatedAt string  `json:"created_at"`
}

type AuthResponse struct {
	Token     string    `json:"token"`
	Developer Developer `json:"developer"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type APIKey struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	KeyPrefix string   `json:"key_prefix"`
	Scopes    []string `json:"scopes"`
	CreatedAt string   `json:"created_at"`
	RevokedAt *string  `json:"revoked_at"`
}

type CreateAPIKeyResponse struct {
	// Key is the full API key — shown only once.
	Key    string `json:"key"`
	APIKey APIKey `json:"api_key"`
	// Warning reminds the caller to store the key securely.
	Warning string `json:"warning"`
}

type DNSRecord struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type DNSRecords struct {
	Verification DNSRecord `json:"verification"`
	DKIM         DNSRecord `json:"dkim"`
	SPF          DNSRecord `json:"spf"`
	DMARC        DNSRecord `json:"dmarc"`
}

type SendingDomain struct {
	ID           string  `json:"id"`
	Domain       string  `json:"domain"`
	Status       string  `json:"status"`
	DkimSelector string  `json:"dkim_selector"`
	VerifiedAt   *string `json:"verified_at"`
	CreatedAt    string  `json:"created_at"`
}

type DomainResponse struct {
	Domain     SendingDomain `json:"domain"`
	DNSRecords DNSRecords    `json:"dns_records"`
}

type EmailSend struct {
	ID          string  `json:"id"`
	Status      string  `json:"status"`
	From        string  `json:"from"`
	To          string  `json:"to"`
	Subject     string  `json:"subject"`
	CreatedAt   string  `json:"created_at"`
	ScheduledAt *string `json:"scheduled_at,omitempty"`
}

type BatchSendResponse struct {
	Emails []EmailSend `json:"emails"`
	Count  int         `json:"count"`
}

type ListEmailsResponse struct {
	Emails []EmailSend `json:"emails"`
	Count  int         `json:"count"`
}

type Suppression struct {
	ID        string  `json:"id"`
	Email     string  `json:"email"`
	Reason    string  `json:"reason"`
	Source    string  `json:"source"`
	Note      *string `json:"note,omitempty"`
	CreatedAt string  `json:"created_at"`
}

type Template struct {
	ID          string             `json:"id"`
	Alias       *string            `json:"alias,omitempty"`
	Name        string             `json:"name"`
	From        *string            `json:"from,omitempty"`
	Subject     *string            `json:"subject,omitempty"`
	ReplyTo     *string            `json:"reply_to,omitempty"`
	HTML        string             `json:"html"`
	Text        *string            `json:"text,omitempty"`
	Variables   []TemplateVariable `json:"variables"`
	Status      string             `json:"status"`
	PublishedAt *string            `json:"published_at,omitempty"`
	CreatedAt   string             `json:"created_at"`
	UpdatedAt   string             `json:"updated_at"`
}

type WebhookEndpoint struct {
	ID        string         `json:"id"`
	URL       string         `json:"url"`
	Events    []WebhookEvent `json:"events"`
	Status    string         `json:"status"`
	Secret    string         `json:"secret,omitempty"`
	CreatedAt string         `json:"created_at"`
}

type WebhookDelivery struct {
	ID             string          `json:"id"`
	EndpointID     string          `json:"endpoint_id"`
	EmailSendID    *string         `json:"email_send_id,omitempty"`
	EventType      WebhookEvent    `json:"event_type"`
	Payload        json.RawMessage `json:"payload"`
	Status         string          `json:"status"`
	Attempts       int32           `json:"attempts"`
	MaxAttempts    int32           `json:"max_attempts"`
	NextAttemptAt  string          `json:"next_attempt_at"`
	LastError      *string         `json:"last_error,omitempty"`
	ResponseStatus *int32          `json:"response_status,omitempty"`
	DeliveredAt    *string         `json:"delivered_at,omitempty"`
	CreatedAt      string          `json:"created_at"`
}

// ── Internal envelope types ───────────────────────────────────────────────────

type developerEnvelope struct {
	Developer Developer `json:"developer"`
}

type tokenRefreshEnvelope struct {
	Token string `json:"token"`
}

type apiKeyListEnvelope struct {
	APIKeys []APIKey `json:"api_keys"`
}

type domainListEnvelope struct {
	Domains []SendingDomain `json:"domains"`
}

type emailSendEnvelope struct {
	Email EmailSend `json:"email"`
}

type webhookEnvelope struct {
	Webhook WebhookEndpoint `json:"webhook"`
}

type webhookListEnvelope struct {
	Webhooks []WebhookEndpoint `json:"webhooks"`
}

type webhookDeliveryListEnvelope struct {
	Deliveries []WebhookDelivery `json:"deliveries"`
	Count      int               `json:"count"`
}

type suppressionEnvelope struct {
	Suppression Suppression `json:"suppression"`
}

type suppressionListEnvelope struct {
	Suppressions []Suppression `json:"suppressions"`
	Count        int           `json:"count"`
}

type templateEnvelope struct {
	Template Template `json:"template"`
}

type templateListEnvelope struct {
	Templates []Template `json:"templates"`
	Count     int        `json:"count"`
}

type apiErrorEnvelope struct {
	Error string `json:"error"`
}
