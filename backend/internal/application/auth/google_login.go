package auth

import (
	"context"
	"fmt"
	"log"

	"github.com/yuki5155/go-google-auth/internal/application/dto"
	"github.com/yuki5155/go-google-auth/internal/application/ports"
	"github.com/yuki5155/go-google-auth/internal/domain/shared"
)

// GoogleLoginUseCase handles Google OAuth login flow
type GoogleLoginUseCase struct {
	oauthValidator ports.OAuthValidator
	tokenGenerator ports.TokenGenerator
	clientID       string
	rootUserEmail  string
}

// NewGoogleLoginUseCase creates a new GoogleLoginUseCase
func NewGoogleLoginUseCase(
	oauthValidator ports.OAuthValidator,
	tokenGenerator ports.TokenGenerator,
	clientID string,
	rootUserEmail string,
) *GoogleLoginUseCase {
	return &GoogleLoginUseCase{
		oauthValidator: oauthValidator,
		tokenGenerator: tokenGenerator,
		clientID:       clientID,
		rootUserEmail:  rootUserEmail,
	}
}

// Execute performs the Google login flow
func (uc *GoogleLoginUseCase) Execute(ctx context.Context, credential string) (*dto.LoginResponse, error) {
	// Validate the Google ID token
	oauthUser, err := uc.oauthValidator.ValidateToken(ctx, credential, uc.clientID)
	if err != nil {
		log.Printf("Failed to verify Google ID token: %v", err)
		return nil, fmt.Errorf("failed to verify Google ID token: %w", err)
	}

	// Check if email is verified
	if !oauthUser.EmailVerified {
		return nil, shared.ErrUnverifiedEmail
	}

	// Determine role based on email
	log.Printf("DEBUG: Checking root email - configured: '%s', oauth email: '%s', match: %v",
		uc.rootUserEmail, oauthUser.Email, uc.rootUserEmail == oauthUser.Email)

	role := "user"
	if uc.rootUserEmail != "" && oauthUser.Email == uc.rootUserEmail {
		role = "root"
		log.Printf("ROOT user logged in: %s (%s)", oauthUser.Email, oauthUser.UserID)
	} else {
		log.Printf("User logged in: %s (%s)", oauthUser.Email, oauthUser.UserID)
	}

	// Generate JWT tokens directly from OAuth data
	userInfo := ports.UserInfo{
		UserID:  oauthUser.UserID,
		Email:   oauthUser.Email,
		Name:    oauthUser.Name,
		Picture: oauthUser.Picture,
		Role:    role,
	}

	accessToken, refreshToken, err := uc.tokenGenerator.GenerateTokenPair(userInfo)
	if err != nil {
		log.Printf("Failed to generate JWT tokens: %v", err)
		return nil, fmt.Errorf("failed to generate authentication tokens: %w", err)
	}

	// Return response
	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserResponse{
			ID:      oauthUser.UserID,
			Email:   oauthUser.Email,
			Name:    oauthUser.Name,
			Picture: oauthUser.Picture,
			Role:    role,
		},
		Message: "Login successful",
	}, nil
}
