package container

import (
	"github.com/yuki5155/go-google-auth/internal/application/auth"
	"github.com/yuki5155/go-google-auth/internal/application/ports"
	"github.com/yuki5155/go-google-auth/internal/domain/user"
	"github.com/yuki5155/go-google-auth/internal/infrastructure/auth/google"
	"github.com/yuki5155/go-google-auth/internal/infrastructure/auth/jwt"
	"github.com/yuki5155/go-google-auth/internal/infrastructure/config"
	"github.com/yuki5155/go-google-auth/internal/infrastructure/persistence/memory"
)

// Container holds all application dependencies
type Container struct {
	// Config
	Config *config.Config

	// Infrastructure
	UserRepository user.Repository
	TokenGenerator ports.TokenGenerator
	OAuthValidator ports.OAuthValidator

	// Use Cases
	GoogleLoginUseCase    *auth.GoogleLoginUseCase
	RefreshTokenUseCase   *auth.RefreshTokenUseCase
	GetCurrentUserUseCase *auth.GetCurrentUserUseCase
	LogoutUseCase         *auth.LogoutUseCase
}

// NewContainer creates and wires all dependencies
func NewContainer(cfg *config.Config) *Container {
	// Infrastructure layer
	userRepo := memory.NewUserRepository()
	tokenGen := jwt.NewService(cfg.JWTSecret)
	oauthValidator := google.NewValidator()

	// Application layer - Use cases
	googleLoginUC := auth.NewGoogleLoginUseCase(
		userRepo,
		oauthValidator,
		tokenGen,
		cfg.GoogleClientID,
	)
	refreshTokenUC := auth.NewRefreshTokenUseCase(tokenGen)
	getCurrentUserUC := auth.NewGetCurrentUserUseCase(userRepo, tokenGen)
	logoutUC := auth.NewLogoutUseCase()

	return &Container{
		Config:                cfg,
		UserRepository:        userRepo,
		TokenGenerator:        tokenGen,
		OAuthValidator:        oauthValidator,
		GoogleLoginUseCase:    googleLoginUC,
		RefreshTokenUseCase:   refreshTokenUC,
		GetCurrentUserUseCase: getCurrentUserUC,
		LogoutUseCase:         logoutUC,
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
