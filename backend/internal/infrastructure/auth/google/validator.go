package google

import (
	"context"

	"github.com/yuki5155/go-google-auth/internal/application/ports"
	"google.golang.org/api/idtoken"
)

// Validator validates Google OAuth ID tokens and implements ports.OAuthValidator
type Validator struct{}

// NewValidator creates a new Google OAuth Validator
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateToken validates a Google ID token and returns user information
func (v *Validator) ValidateToken(ctx context.Context, idToken string, audience string) (*ports.OAuthUserInfo, error) {
	payload, err := idtoken.Validate(ctx, idToken, audience)
	if err != nil {
		return nil, err
	}

	// Extract user information from payload claims
	userID, _ := payload.Claims["sub"].(string)
	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)
	picture, _ := payload.Claims["picture"].(string)
	emailVerified, _ := payload.Claims["email_verified"].(bool)

	return &ports.OAuthUserInfo{
		UserID:        userID,
		Email:         email,
		EmailVerified: emailVerified,
		Name:          name,
		Picture:       picture,
	}, nil
}
