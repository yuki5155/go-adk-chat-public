package ports

import "time"

// UserInfo represents user information for token generation
type UserInfo struct {
	UserID  string
	Email   string
	Name    string
	Picture string
}

// TokenPair represents an access token and refresh token pair
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

// TokenClaims represents the claims extracted from a token
type TokenClaims struct {
	UserID  string
	Email   string
	Name    string
	Picture string
}

// TokenGenerator defines the interface for JWT token operations
type TokenGenerator interface {
	// GenerateTokenPair generates both access and refresh tokens
	GenerateTokenPair(userInfo UserInfo) (accessToken, refreshToken string, err error)

	// RefreshAccessToken generates a new access token from a valid refresh token
	RefreshAccessToken(refreshToken string) (string, error)

	// ValidateAccessToken validates an access token and returns the claims
	ValidateAccessToken(accessToken string) (*TokenClaims, error)

	// GetAccessTokenExpiry returns the access token expiry duration in seconds
	GetAccessTokenExpiry() int

	// GetRefreshTokenExpiry returns the refresh token expiry duration in seconds
	GetRefreshTokenExpiry() int
}

// Common errors for token operations
var (
	ErrInvalidToken = &TokenError{Message: "invalid token"}
	ErrExpiredToken = &TokenError{Message: "token has expired"}
)

// TokenError represents a token-related error
type TokenError struct {
	Message string
}

func (e *TokenError) Error() string {
	return e.Message
}

// IsTokenExpired checks if an error is an expired token error
func IsTokenExpired(err error) bool {
	if err == nil {
		return false
	}
	if tokenErr, ok := err.(*TokenError); ok {
		return tokenErr == ErrExpiredToken
	}
	return err.Error() == "token has expired"
}

// TokenMetadata contains metadata about token lifetimes
type TokenMetadata struct {
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}
