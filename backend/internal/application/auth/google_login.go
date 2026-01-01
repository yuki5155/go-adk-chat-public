package auth

import (
	"context"
	"fmt"
	"log"

	"github.com/yuki5155/go-google-auth/internal/application/dto"
	"github.com/yuki5155/go-google-auth/internal/application/ports"
	"github.com/yuki5155/go-google-auth/internal/domain/shared"
	"github.com/yuki5155/go-google-auth/internal/domain/user"
)

// GoogleLoginUseCase handles Google OAuth login flow
type GoogleLoginUseCase struct {
	userRepo       user.Repository
	oauthValidator ports.OAuthValidator
	tokenGenerator ports.TokenGenerator
	clientID       string
}

// NewGoogleLoginUseCase creates a new GoogleLoginUseCase
func NewGoogleLoginUseCase(
	userRepo user.Repository,
	oauthValidator ports.OAuthValidator,
	tokenGenerator ports.TokenGenerator,
	clientID string,
) *GoogleLoginUseCase {
	return &GoogleLoginUseCase{
		userRepo:       userRepo,
		oauthValidator: oauthValidator,
		tokenGenerator: tokenGenerator,
		clientID:       clientID,
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

	// Create domain value objects
	userID, err := user.NewUserID(oauthUser.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID from OAuth: %w", err)
	}

	email, err := user.NewEmail(oauthUser.Email, oauthUser.EmailVerified)
	if err != nil {
		return nil, fmt.Errorf("invalid email from OAuth: %w", err)
	}

	profile := user.NewProfile(oauthUser.Name, oauthUser.Picture)

	// Check if user already exists
	existingUser, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil && err != shared.ErrUserNotFound {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}

	var domainUser *user.User
	if existingUser != nil {
		// User exists - update and record login
		domainUser = existingUser
		domainUser.UpdateProfile(profile)
		domainUser.RecordLogin()

		if err := uc.userRepo.Save(ctx, domainUser); err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}

		log.Printf("Existing user logged in: %s (%s)", email.Value(), userID.Value())
	} else {
		// New user - create and save
		domainUser, err = user.NewUser(userID, email, profile)
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}

		if err := uc.userRepo.Save(ctx, domainUser); err != nil {
			return nil, fmt.Errorf("failed to save new user: %w", err)
		}

		log.Printf("New user registered: %s (%s)", email.Value(), userID.Value())
	}

	// Generate JWT tokens
	userInfo := ports.UserInfo{
		UserID:  domainUser.ID().Value(),
		Email:   domainUser.Email().Value(),
		Name:    domainUser.Profile().Name(),
		Picture: domainUser.Profile().Picture(),
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
		User:         dto.FromDomain(domainUser),
		Message:      "Login successful",
	}, nil
}
