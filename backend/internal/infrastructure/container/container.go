package container

import (
	"github.com/yuki5155/go-google-auth/internal/application/auth"
	"github.com/yuki5155/go-google-auth/internal/application/ports"
	"github.com/yuki5155/go-google-auth/internal/infrastructure/auth/google"
	"github.com/yuki5155/go-google-auth/internal/infrastructure/auth/jwt"
	"github.com/yuki5155/go-google-auth/internal/infrastructure/config"
)

// Container holds all application dependencies
type Container struct {
	// Config
	Config *config.Config

	// Infrastructure
	TokenGenerator ports.TokenGenerator
	OAuthValidator ports.OAuthValidator

	// Use Cases
	GoogleLoginUseCase  *auth.GoogleLoginUseCase
	RefreshTokenUseCase *auth.RefreshTokenUseCase
	LogoutUseCase       *auth.LogoutUseCase
}

// NewContainer creates and wires all dependencies
func NewContainer(cfg *config.Config) *Container {
	// Infrastructure layer
	tokenGen := jwt.NewService(cfg.JWTSecret)
	oauthValidator := google.NewValidator()

	// Application layer - Use cases
	googleLoginUC := auth.NewGoogleLoginUseCase(
		oauthValidator,
		tokenGen,
		cfg.GoogleClientID,
		cfg.RootUserEmail,
	)
	refreshTokenUC := auth.NewRefreshTokenUseCase(tokenGen)
	logoutUC := auth.NewLogoutUseCase()

	return &Container{
		Config:              cfg,
		TokenGenerator:      tokenGen,
		OAuthValidator:      oauthValidator,
		GoogleLoginUseCase:  googleLoginUC,
		RefreshTokenUseCase: refreshTokenUC,
		LogoutUseCase:       logoutUC,
	}
}

// GetTokenGenerator returns the token generator (for middleware)
func (c *Container) GetTokenGenerator() ports.TokenGenerator {
	return c.TokenGenerator
}

// GetConfig returns the config
func (c *Container) GetConfig() *config.Config {
	return c.Config
}
