package container

import (
	"context"
	"log"

	"github.com/yuki5155/go-google-auth/internal/application/admin"
	"github.com/yuki5155/go-google-auth/internal/application/auth"
	"github.com/yuki5155/go-google-auth/internal/application/ports"
	"github.com/yuki5155/go-google-auth/internal/domain/role"
	"github.com/yuki5155/go-google-auth/internal/infrastructure/auth/google"
	"github.com/yuki5155/go-google-auth/internal/infrastructure/auth/jwt"
	"github.com/yuki5155/go-google-auth/internal/infrastructure/config"
	"github.com/yuki5155/go-google-auth/internal/infrastructure/dynamodb"
	"github.com/yuki5155/go-google-auth/internal/infrastructure/persistence"
)

// Container holds all application dependencies
type Container struct {
	// Config
	Config *config.Config

	// Infrastructure
	TokenGenerator ports.TokenGenerator
	OAuthValidator ports.OAuthValidator
	RoleRepository role.Repository

	// Auth Use Cases
	GoogleLoginUseCase  *auth.GoogleLoginUseCase
	RefreshTokenUseCase *auth.RefreshTokenUseCase
	LogoutUseCase       *auth.LogoutUseCase

	// Admin Use Cases
	RequestRoleUseCase       *admin.RequestRoleUseCase
	ListPendingRequestsUseCase *admin.ListPendingRequestsUseCase
	ApproveRequestUseCase    *admin.ApproveRequestUseCase
	RejectRequestUseCase     *admin.RejectRequestUseCase
	CheckUserRoleUseCase     *admin.CheckUserRoleUseCase
	ListUsersByRoleUseCase   *admin.ListUsersByRoleUseCase
}

// NewContainer creates and wires all dependencies
func NewContainer(cfg *config.Config) *Container {
	// Infrastructure layer
	tokenGen := jwt.NewService(cfg.JWTSecret)
	oauthValidator := google.NewValidator()

	// DynamoDB client
	ctx := context.Background()
	dynamoClient, err := dynamodb.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create DynamoDB client: %v", err)
	}
	roleRepo := persistence.NewRoleRepository(dynamoClient)

	// Application layer - Auth use cases
	googleLoginUC := auth.NewGoogleLoginUseCase(
		oauthValidator,
		tokenGen,
		roleRepo,
		cfg.GoogleClientID,
		cfg.RootUserEmail,
	)
	refreshTokenUC := auth.NewRefreshTokenUseCase(tokenGen)
	logoutUC := auth.NewLogoutUseCase()

	// Application layer - Admin use cases
	requestRoleUC := admin.NewRequestRoleUseCase(roleRepo)
	listPendingUC := admin.NewListPendingRequestsUseCase(roleRepo)
	approveRequestUC := admin.NewApproveRequestUseCase(roleRepo)
	rejectRequestUC := admin.NewRejectRequestUseCase(roleRepo)
	checkUserRoleUC := admin.NewCheckUserRoleUseCase(roleRepo)
	listUsersByRoleUC := admin.NewListUsersByRoleUseCase(roleRepo)

	return &Container{
		Config:              cfg,
		TokenGenerator:      tokenGen,
		OAuthValidator:      oauthValidator,
		RoleRepository:      roleRepo,
		GoogleLoginUseCase:  googleLoginUC,
		RefreshTokenUseCase: refreshTokenUC,
		LogoutUseCase:       logoutUC,
		RequestRoleUseCase:       requestRoleUC,
		ListPendingRequestsUseCase: listPendingUC,
		ApproveRequestUseCase:    approveRequestUC,
		RejectRequestUseCase:     rejectRequestUC,
		CheckUserRoleUseCase:     checkUserRoleUC,
		ListUsersByRoleUseCase:   listUsersByRoleUC,
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
