package smail

import "context"

// AuthService handles developer account and authentication endpoints.
type AuthService struct{ c *Client }

// Register creates a new developer account. The account starts in "pending"
// status; the developer must verify their email before logging in.
func (s *AuthService) Register(ctx context.Context, req RegisterRequest) (Developer, error) {
	var env struct {
		Developer Developer `json:"developer"`
		Message   string    `json:"message"`
	}
	if err := s.c.do(ctx, "POST", "/api/v1/auth/register", req, &env); err != nil {
		return Developer{}, err
	}
	return env.Developer, nil
}

// Login authenticates with email and password and returns a JWT token and
// developer profile. Store the token and pass it to WithToken or SetToken.
func (s *AuthService) Login(ctx context.Context, req LoginRequest) (AuthResponse, error) {
	var resp AuthResponse
	if err := s.c.do(ctx, "POST", "/api/v1/auth/login", req, &resp); err != nil {
		return AuthResponse{}, err
	}
	return resp, nil
}

// VerifyEmail confirms email ownership using the token sent during registration.
// On success it returns a JWT token so the developer can log in immediately.
func (s *AuthService) VerifyEmail(ctx context.Context, req VerifyEmailRequest) (AuthResponse, error) {
	var resp AuthResponse
	if err := s.c.do(ctx, "POST", "/api/v1/auth/verify-email", req, &resp); err != nil {
		return AuthResponse{}, err
	}
	return resp, nil
}

// ForgotPassword requests a password reset email. Always returns success to
// prevent email enumeration.
func (s *AuthService) ForgotPassword(ctx context.Context, req ForgotPasswordRequest) error {
	return s.c.do(ctx, "POST", "/api/v1/auth/forgot-password", req, nil)
}

// ResetPassword sets a new password using a reset token from email.
func (s *AuthService) ResetPassword(ctx context.Context, req ResetPasswordRequest) error {
	return s.c.do(ctx, "POST", "/api/v1/auth/reset-password", req, nil)
}

// Refresh exchanges a valid JWT for a fresh one. The client must be configured
// with a Bearer token via WithToken or SetToken.
func (s *AuthService) Refresh(ctx context.Context) (string, error) {
	var env tokenRefreshEnvelope
	if err := s.c.do(ctx, "POST", "/api/v1/auth/refresh", nil, &env); err != nil {
		return "", err
	}
	return env.Token, nil
}

// Me returns the profile of the currently authenticated developer.
func (s *AuthService) Me(ctx context.Context) (Developer, error) {
	var env developerEnvelope
	if err := s.c.do(ctx, "GET", "/api/v1/auth/me", nil, &env); err != nil {
		return Developer{}, err
	}
	return env.Developer, nil
}

// UpdateMe updates the profile of the currently authenticated developer.
func (s *AuthService) UpdateMe(ctx context.Context, req UpdateProfileRequest) (Developer, error) {
	var env developerEnvelope
	if err := s.c.do(ctx, "PUT", "/api/v1/auth/me", req, &env); err != nil {
		return Developer{}, err
	}
	return env.Developer, nil
}
