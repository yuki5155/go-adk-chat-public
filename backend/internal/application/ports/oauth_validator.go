package ports

import "context"

// OAuthUserInfo represents user information from OAuth provider
type OAuthUserInfo struct {
	UserID        string
	Email         string
	EmailVerified bool
	Name          string
	Picture       string
}

// OAuthValidator defines the interface for OAuth token validation
type OAuthValidator interface {
	// ValidateToken validates an OAuth ID token and returns user information
	// The audience parameter is the OAuth client ID that the token should be issued to
	ValidateToken(ctx context.Context, idToken string, audience string) (*OAuthUserInfo, error)
}

// Common OAuth errors
var (
	ErrInvalidOAuthToken     = &OAuthError{Message: "invalid OAuth token"}
	ErrUnverifiedEmail       = &OAuthError{Message: "email address is not verified"}
	ErrInvalidAudience       = &OAuthError{Message: "invalid audience"}
	ErrExpiredOAuthToken     = &OAuthError{Message: "OAuth token has expired"}
)

// OAuthError represents an OAuth-related error
type OAuthError struct {
	Message string
}

func (e *OAuthError) Error() string {
	return e.Message
}
