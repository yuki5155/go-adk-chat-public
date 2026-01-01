package auth

import (
	"context"
	"fmt"

	"github.com/yuki5155/go-google-auth/internal/application/dto"
	"github.com/yuki5155/go-google-auth/internal/application/ports"
)

// RefreshTokenUseCase handles token refresh operations
type RefreshTokenUseCase struct {
	tokenGenerator ports.TokenGenerator
}

// NewRefreshTokenUseCase creates a new RefreshTokenUseCase
func NewRefreshTokenUseCase(tokenGenerator ports.TokenGenerator) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{
		tokenGenerator: tokenGenerator,
	}
}

// Execute generates a new access token from a valid refresh token
func (uc *RefreshTokenUseCase) Execute(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error) {
	if refreshToken == "" {
		return nil, fmt.Errorf("refresh token is required")
	}

	// Generate new access token from refresh token
	newAccessToken, err := uc.tokenGenerator.RefreshAccessToken(refreshToken)
	if err != nil {
		if ports.IsTokenExpired(err) {
			return nil, ports.ErrExpiredToken
		}
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	return &dto.RefreshResponse{
		AccessToken: newAccessToken,
		Message:     "Token refreshed successfully",
	}, nil
}
