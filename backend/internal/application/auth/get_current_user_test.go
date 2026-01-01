package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/yuki5155/go-google-auth/internal/application/ports"
	"github.com/yuki5155/go-google-auth/internal/domain/shared"
	"github.com/yuki5155/go-google-auth/internal/domain/user"
	"github.com/yuki5155/go-google-auth/internal/mocks"
)

func TestGetCurrentUserUseCase_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockRepository(ctrl)
	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)

	// Pre-create user
	userID, _ := user.NewUserID("test-user-123")
	email, _ := user.NewEmail("test@example.com", true)
	profile := user.NewProfile("Test User", "https://example.com/photo.jpg")
	domainUser, _ := user.NewUser(userID, email, profile)

	tokenClaims := &ports.TokenClaims{
		UserID:  "test-user-123",
		Email:   "test@example.com",
		Name:    "Test User",
		Picture: "https://example.com/photo.jpg",
	}

	// Set expectations
	mockTokenGen.EXPECT().
		ValidateAccessToken("valid-access-token").
		Return(tokenClaims, nil)

	mockRepo.EXPECT().
		FindByID(ctx, userID).
		Return(domainUser, nil)

	useCase := NewGetCurrentUserUseCase(mockRepo, mockTokenGen)

	result, err := useCase.Execute(ctx, "valid-access-token")

	require.NoError(t, err)
	assert.Equal(t, "test-user-123", result.ID)
	assert.Equal(t, "test@example.com", result.Email)
	assert.Equal(t, "Test User", result.Name)
	assert.Equal(t, "https://example.com/photo.jpg", result.Picture)
}

func TestGetCurrentUserUseCase_EmptyToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockRepository(ctrl)
	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)

	useCase := NewGetCurrentUserUseCase(mockRepo, mockTokenGen)

	result, err := useCase.Execute(ctx, "")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, shared.ErrMissingToken, err)
}

func TestGetCurrentUserUseCase_InvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockRepository(ctrl)
	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)

	mockTokenGen.EXPECT().
		ValidateAccessToken("invalid-token").
		Return(nil, errors.New("invalid token"))

	useCase := NewGetCurrentUserUseCase(mockRepo, mockTokenGen)

	result, err := useCase.Execute(ctx, "invalid-token")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid access token")
}

func TestGetCurrentUserUseCase_ExpiredToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockRepository(ctrl)
	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)

	mockTokenGen.EXPECT().
		ValidateAccessToken("expired-token").
		Return(nil, ports.ErrExpiredToken)

	useCase := NewGetCurrentUserUseCase(mockRepo, mockTokenGen)

	result, err := useCase.Execute(ctx, "expired-token")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, shared.ErrExpiredToken, err)
}

func TestGetCurrentUserUseCase_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockRepository(ctrl)
	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)

	tokenClaims := &ports.TokenClaims{
		UserID:  "nonexistent-user",
		Email:   "test@example.com",
		Name:    "Test User",
		Picture: "https://example.com/photo.jpg",
	}

	userID, _ := user.NewUserID("nonexistent-user")

	mockTokenGen.EXPECT().
		ValidateAccessToken("valid-token").
		Return(tokenClaims, nil)

	mockRepo.EXPECT().
		FindByID(ctx, userID).
		Return(nil, shared.ErrUserNotFound)

	useCase := NewGetCurrentUserUseCase(mockRepo, mockTokenGen)

	result, err := useCase.Execute(ctx, "valid-token")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, shared.ErrUnauthorized, err)
}

func TestGetCurrentUserUseCase_ExecuteFromClaims_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockRepository(ctrl)
	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)

	// Pre-create user
	userID, _ := user.NewUserID("test-user-123")
	email, _ := user.NewEmail("test@example.com", true)
	profile := user.NewProfile("Test User", "https://example.com/photo.jpg")
	domainUser, _ := user.NewUser(userID, email, profile)

	mockRepo.EXPECT().
		FindByID(ctx, userID).
		Return(domainUser, nil)

	useCase := NewGetCurrentUserUseCase(mockRepo, mockTokenGen)

	claims := &ports.TokenClaims{
		UserID:  "test-user-123",
		Email:   "test@example.com",
		Name:    "Test User",
		Picture: "https://example.com/photo.jpg",
	}

	result, err := useCase.ExecuteFromClaims(ctx, claims)

	require.NoError(t, err)
	assert.Equal(t, "test-user-123", result.ID)
	assert.Equal(t, "test@example.com", result.Email)
}

func TestGetCurrentUserUseCase_ExecuteFromClaims_NilClaims(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockRepository(ctrl)
	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)

	useCase := NewGetCurrentUserUseCase(mockRepo, mockTokenGen)

	result, err := useCase.ExecuteFromClaims(ctx, nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, shared.ErrMissingToken, err)
}

func TestGetCurrentUserUseCase_ExecuteFromClaims_InvalidUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockRepository(ctrl)
	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)

	useCase := NewGetCurrentUserUseCase(mockRepo, mockTokenGen)

	claims := &ports.TokenClaims{
		UserID:  "",  // Empty user ID
		Email:   "test@example.com",
		Name:    "Test User",
		Picture: "https://example.com/photo.jpg",
	}

	result, err := useCase.ExecuteFromClaims(ctx, claims)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid user ID in token")
}

func TestGetCurrentUserUseCase_ExecuteFromClaims_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockRepository(ctrl)
	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)

	userID, _ := user.NewUserID("test-user-123")

	mockRepo.EXPECT().
		FindByID(ctx, userID).
		Return(nil, errors.New("database connection error"))

	useCase := NewGetCurrentUserUseCase(mockRepo, mockTokenGen)

	claims := &ports.TokenClaims{
		UserID:  "test-user-123",
		Email:   "test@example.com",
		Name:    "Test User",
		Picture: "https://example.com/photo.jpg",
	}

	result, err := useCase.ExecuteFromClaims(ctx, claims)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to retrieve user")
}

func TestGetCurrentUserUseCase_ExecuteFromClaims_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockRepository(ctrl)
	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)

	userID, _ := user.NewUserID("nonexistent-user")

	mockRepo.EXPECT().
		FindByID(ctx, userID).
		Return(nil, shared.ErrUserNotFound)

	useCase := NewGetCurrentUserUseCase(mockRepo, mockTokenGen)

	claims := &ports.TokenClaims{
		UserID:  "nonexistent-user",
		Email:   "test@example.com",
		Name:    "Test User",
		Picture: "https://example.com/photo.jpg",
	}

	result, err := useCase.ExecuteFromClaims(ctx, claims)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, shared.ErrUnauthorized, err)
}
