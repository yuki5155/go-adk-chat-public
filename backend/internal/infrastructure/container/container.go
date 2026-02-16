// Package container provides dependency injection for the application
package container

import (
	"context"
	"log"

	"github.com/yuki5155/go-google-auth/internal/application/admin"
	"github.com/yuki5155/go-google-auth/internal/application/auth"
	chatApp "github.com/yuki5155/go-google-auth/internal/application/chat"
	"github.com/yuki5155/go-google-auth/internal/application/ports"
	"github.com/yuki5155/go-google-auth/internal/domain/chat"
	"github.com/yuki5155/go-google-auth/internal/domain/role"
	"github.com/yuki5155/go-google-auth/internal/infrastructure/adk"
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

	// Chat Repositories
	ThreadRepository  chat.ThreadRepository
	SessionRepository chat.SessionRepository
	EventRepository   chat.EventRepository
	MemoryRepository  chat.MemoryRepository

	// ADK Runner
	ADKRunner *adk.Runner

	// Auth Use Cases
	GoogleLoginUseCase  *auth.GoogleLoginUseCase
	RefreshTokenUseCase *auth.RefreshTokenUseCase
	LogoutUseCase       *auth.LogoutUseCase

	// Admin Use Cases
	RequestRoleUseCase         *admin.RequestRoleUseCase
	ListPendingRequestsUseCase *admin.ListPendingRequestsUseCase
	ApproveRequestUseCase      *admin.ApproveRequestUseCase
	RejectRequestUseCase       *admin.RejectRequestUseCase
	CheckUserRoleUseCase       *admin.CheckUserRoleUseCase
	ListUsersByRoleUseCase     *admin.ListUsersByRoleUseCase

	// Chat Use Cases
	CreateThreadUseCase *chatApp.CreateThreadUseCase
	ListThreadsUseCase  *chatApp.ListThreadsUseCase
	GetThreadUseCase    *chatApp.GetThreadUseCase
	SendMessageUseCase  *chatApp.SendMessageUseCase
	DeleteThreadUseCase *chatApp.DeleteThreadUseCase
	ListModelsUseCase   *chatApp.ListModelsUseCase
}

// NewContainer creates and wires all dependencies
func NewContainer(cfg *config.Config) *Container {
	// Infrastructure layer
	tokenGen := jwt.NewService(cfg.GetJWTSecret())
	oauthValidator := google.NewValidator()

	// DynamoDB client
	ctx := context.Background()
	dynamoClient, err := dynamodb.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create DynamoDB client: %v", err)
	}
	roleRepo := persistence.NewRoleRepository(dynamoClient)

	// Chat repositories
	threadRepo := persistence.NewChatThreadRepository(dynamoClient)
	sessionRepo := persistence.NewChatSessionRepository(dynamoClient)
	eventRepo := persistence.NewChatEventRepository(dynamoClient)
	memoryRepo := persistence.NewChatMemoryRepository(dynamoClient)

	// ADK Runner
	adkConfig := adk.NewConfigFromEnv()
	adkRunner, err := adk.NewRunner(adkConfig)
	if err != nil {
		log.Printf("Warning: Failed to create ADK runner: %v (chat will not work)", err)
	}

	// Create ADK adapter for AIRunner interface
	var aiRunner *adk.RunnerAdapter
	if adkRunner != nil {
		aiRunner = adk.NewRunnerAdapter(adkRunner)
	}

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

	// Application layer - Chat use cases
	createThreadUC := chatApp.NewCreateThreadUseCase(threadRepo)
	listThreadsUC := chatApp.NewListThreadsUseCase(threadRepo)
	getThreadUC := chatApp.NewGetThreadUseCase(threadRepo, sessionRepo, eventRepo)
	sendMessageUC := chatApp.NewSendMessageUseCase(threadRepo, sessionRepo, eventRepo, aiRunner)
	deleteThreadUC := chatApp.NewDeleteThreadUseCase(threadRepo, sessionRepo, eventRepo, memoryRepo)
	listModelsUC := chatApp.NewListModelsUseCase()

	return &Container{
		Config:         cfg,
		TokenGenerator: tokenGen,
		OAuthValidator: oauthValidator,
		RoleRepository: roleRepo,

		// Chat Repositories
		ThreadRepository:  threadRepo,
		SessionRepository: sessionRepo,
		EventRepository:   eventRepo,
		MemoryRepository:  memoryRepo,

		// ADK Runner
		ADKRunner: adkRunner,

		// Auth Use Cases
		GoogleLoginUseCase:  googleLoginUC,
		RefreshTokenUseCase: refreshTokenUC,
		LogoutUseCase:       logoutUC,

		// Admin Use Cases
		RequestRoleUseCase:         requestRoleUC,
		ListPendingRequestsUseCase: listPendingUC,
		ApproveRequestUseCase:      approveRequestUC,
		RejectRequestUseCase:       rejectRequestUC,
		CheckUserRoleUseCase:       checkUserRoleUC,
		ListUsersByRoleUseCase:     listUsersByRoleUC,

		// Chat Use Cases
		CreateThreadUseCase: createThreadUC,
		ListThreadsUseCase:  listThreadsUC,
		GetThreadUseCase:    getThreadUC,
		SendMessageUseCase:  sendMessageUC,
		DeleteThreadUseCase: deleteThreadUC,
		ListModelsUseCase:   listModelsUC,
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
