package auth

import (
	"context"
	"fmt"

	"github.com/yuki5155/go-google-auth/internal/application/dto"
	"github.com/yuki5155/go-google-auth/internal/application/ports"
	"github.com/yuki5155/go-google-auth/internal/domain/shared"
	"github.com/yuki5155/go-google-auth/internal/domain/user"
)

// GetCurrentUserUseCase handles retrieving current user information
type GetCurrentUserUseCase struct {
	userRepo       user.Repository
	tokenGenerator ports.TokenGenerator
}

// NewGetCurrentUserUseCase creates a new GetCurrentUserUseCase
func NewGetCurrentUserUseCase(userRepo user.Repository, tokenGenerator ports.TokenGenerator) *GetCurrentUserUseCase {
	return &GetCurrentUserUseCase{
		userRepo:       userRepo,
		tokenGenerator: tokenGenerator,
	}
}

// Execute retrieves the current user's information from an access token
func (uc *GetCurrentUserUseCase) Execute(ctx context.Context, accessToken string) (*dto.UserResponse, error) {
	if accessToken == "" {
		return nil, shared.ErrMissingToken
	}

	// Validate access token and extract claims
	claims, err := uc.tokenGenerator.ValidateAccessToken(accessToken)
	if err != nil {
		if ports.IsTokenExpired(err) {
			return nil, shared.ErrExpiredToken
		}
		return nil, fmt.Errorf("invalid access token: %w", err)
	}

	// Get user from repository
	userID, err := user.NewUserID(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID in token: %w", err)
	}

	domainUser, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		if err == shared.ErrUserNotFound {
			return nil, shared.ErrUnauthorized
		}
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}

	response := dto.FromDomain(domainUser)
	return &response, nil
}

// ExecuteFromClaims retrieves user information using token claims directly (for middleware use)
func (uc *GetCurrentUserUseCase) ExecuteFromClaims(ctx context.Context, claims *ports.TokenClaims) (*dto.UserResponse, error) {
	if claims == nil {
		return nil, shared.ErrMissingToken
	}

	// Get user from repository
	userID, err := user.NewUserID(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID in token: %w", err)
	}

	domainUser, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		if err == shared.ErrUserNotFound {
			return nil, shared.ErrUnauthorized
		}
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}

	response := dto.FromDomain(domainUser)
	return &response, nil
}
