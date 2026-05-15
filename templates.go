package smail

import (
	"context"
	"fmt"
	"net/url"
)

// TemplateService manages stored email templates.
type TemplateService struct{ c *Client }

func (s *TemplateService) Create(ctx context.Context, req CreateTemplateRequest) (Template, error) {
	var env templateEnvelope
	if err := s.c.do(ctx, "POST", "/api/v1/templates", req, &env); err != nil {
		return Template{}, err
	}
	return env.Template, nil
}

func (s *TemplateService) List(ctx context.Context, opts ListTemplatesOptions) ([]Template, error) {
	limit := opts.Limit
	if limit <= 0 {
		limit = 50
	}
	offset := opts.Offset

	path := fmt.Sprintf("/api/v1/templates?%s", url.Values{
		"limit":  []string{fmt.Sprintf("%d", limit)},
		"offset": []string{fmt.Sprintf("%d", offset)},
	}.Encode())

	var env templateListEnvelope
	if err := s.c.do(ctx, "GET", path, nil, &env); err != nil {
		return nil, err
	}
	return env.Templates, nil
}

func (s *TemplateService) Get(ctx context.Context, templateID string) (Template, error) {
	var env templateEnvelope
	if err := s.c.do(ctx, "GET", fmt.Sprintf("/api/v1/templates/%s", templateID), nil, &env); err != nil {
		return Template{}, err
	}
	return env.Template, nil
}

func (s *TemplateService) Update(ctx context.Context, templateID string, req UpdateTemplateRequest) (Template, error) {
	var env templateEnvelope
	if err := s.c.do(ctx, "PUT", fmt.Sprintf("/api/v1/templates/%s", templateID), req, &env); err != nil {
		return Template{}, err
	}
	return env.Template, nil
}

func (s *TemplateService) Publish(ctx context.Context, templateID string) (Template, error) {
	var env templateEnvelope
	if err := s.c.do(ctx, "POST", fmt.Sprintf("/api/v1/templates/%s/publish", templateID), nil, &env); err != nil {
		return Template{}, err
	}
	return env.Template, nil
}

func (s *TemplateService) Delete(ctx context.Context, templateID string) error {
	return s.c.do(ctx, "DELETE", fmt.Sprintf("/api/v1/templates/%s", templateID), nil, nil)
}
