package dto

// GoogleLoginRequest represents a Google OAuth login request
type GoogleLoginRequest struct {
	Credential string `json:"credential" binding:"required"`
}

