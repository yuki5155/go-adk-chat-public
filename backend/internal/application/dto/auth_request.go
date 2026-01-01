package dto

// GoogleLoginRequest represents a Google OAuth login request
type GoogleLoginRequest struct {
	Credential string `json:"credential" binding:"required"`
}

// RefreshTokenRequest represents a token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
