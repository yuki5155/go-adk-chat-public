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

func TestGoogleLoginUseCase_NewUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockRepository(ctrl)
	mockOAuth := mocks.NewMockOAuthValidator(ctrl)
	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)

	oauthInfo := &ports.OAuthUserInfo{
		UserID:        "google-user-123",
		Email:         "newuser@example.com",
		EmailVerified: true,
		Name:          "New User",
		Picture:       "https://example.com/photo.jpg",
	}

	// Set expectations
	mockOAuth.EXPECT().
		ValidateToken(ctx, "valid-google-token", "test-client-id").
		Return(oauthInfo, nil)

	userID, _ := user.NewUserID("google-user-123")
	mockRepo.EXPECT().
		FindByID(ctx, userID).
		Return(nil, shared.ErrUserNotFound)

	mockRepo.EXPECT().
		Save(ctx, gomock.Any()).
		Return(nil)

	mockTokenGen.EXPECT().
		GenerateTokenPair(gomock.Any()).
		Return("mock-access-token", "mock-refresh-token", nil)

	useCase := NewGoogleLoginUseCase(mockRepo, mockOAuth, mockTokenGen, "test-client-id")

	result, err := useCase.Execute(ctx, "valid-google-token")

	require.NoError(t, err)
	assert.Equal(t, "Login successful", result.Message)
	assert.Equal(t, "google-user-123", result.User.ID)
	assert.Equal(t, "newuser@example.com", result.User.Email)
	assert.Equal(t, "New User", result.User.Name)
	assert.Equal(t, "https://example.com/photo.jpg", result.User.Picture)
	assert.Equal(t, "mock-access-token", result.AccessToken)
	assert.Equal(t, "mock-refresh-token", result.RefreshToken)
}

func TestGoogleLoginUseCase_ExistingUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockRepository(ctrl)
	mockOAuth := mocks.NewMockOAuthValidator(ctrl)
	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)

	// Pre-create existing user
	userID, _ := user.NewUserID("google-user-123")
	email, _ := user.NewEmail("existing@example.com", true)
	profile := user.NewProfile("Old Name", "https://example.com/old.jpg")
	existingUser, _ := user.NewUser(userID, email, profile)

	oauthInfo := &ports.OAuthUserInfo{
		UserID:        "google-user-123",
		Email:         "existing@example.com",
		EmailVerified: true,
		Name:          "Updated Name",
		Picture:       "https://example.com/new.jpg",
	}

	// Set expectations
	mockOAuth.EXPECT().
		ValidateToken(ctx, "valid-google-token", "test-client-id").
		Return(oauthInfo, nil)

	mockRepo.EXPECT().
		FindByID(ctx, userID).
		Return(existingUser, nil)

	mockRepo.EXPECT().
		Save(ctx, gomock.Any()).
		Return(nil)

	mockTokenGen.EXPECT().
		GenerateTokenPair(gomock.Any()).
		Return("mock-access-token", "mock-refresh-token", nil)

	useCase := NewGoogleLoginUseCase(mockRepo, mockOAuth, mockTokenGen, "test-client-id")

	result, err := useCase.Execute(ctx, "valid-google-token")

	require.NoError(t, err)
	assert.Equal(t, "Login successful", result.Message)
	assert.Equal(t, "google-user-123", result.User.ID)
	assert.Equal(t, "existing@example.com", result.User.Email)
	assert.Equal(t, "Updated Name", result.User.Name)
	assert.Equal(t, "https://example.com/new.jpg", result.User.Picture)
}

func TestGoogleLoginUseCase_InvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockRepository(ctrl)
	mockOAuth := mocks.NewMockOAuthValidator(ctrl)
	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)

	mockOAuth.EXPECT().
		ValidateToken(ctx, "invalid-token", "test-client-id").
		Return(nil, errors.New("invalid token"))

	useCase := NewGoogleLoginUseCase(mockRepo, mockOAuth, mockTokenGen, "test-client-id")

	result, err := useCase.Execute(ctx, "invalid-token")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to verify Google ID token")
}

func TestGoogleLoginUseCase_UnverifiedEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockRepository(ctrl)
	mockOAuth := mocks.NewMockOAuthValidator(ctrl)
	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)

	oauthInfo := &ports.OAuthUserInfo{
		UserID:        "google-user-123",
		Email:         "unverified@example.com",
		EmailVerified: false, // Not verified
		Name:          "Test User",
		Picture:       "https://example.com/photo.jpg",
	}

	mockOAuth.EXPECT().
		ValidateToken(ctx, "valid-token", "test-client-id").
		Return(oauthInfo, nil)

	useCase := NewGoogleLoginUseCase(mockRepo, mockOAuth, mockTokenGen, "test-client-id")

	result, err := useCase.Execute(ctx, "valid-token")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, shared.ErrUnverifiedEmail, err)
}

func TestGoogleLoginUseCase_TokenGenerationFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockRepository(ctrl)
	mockOAuth := mocks.NewMockOAuthValidator(ctrl)
	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)

	oauthInfo := &ports.OAuthUserInfo{
		UserID:        "google-user-123",
		Email:         "test@example.com",
		EmailVerified: true,
		Name:          "Test User",
		Picture:       "https://example.com/photo.jpg",
	}

	userID, _ := user.NewUserID("google-user-123")

	mockOAuth.EXPECT().
		ValidateToken(ctx, "valid-token", "test-client-id").
		Return(oauthInfo, nil)

	mockRepo.EXPECT().
		FindByID(ctx, userID).
		Return(nil, shared.ErrUserNotFound)

	mockRepo.EXPECT().
		Save(ctx, gomock.Any()).
		Return(nil)

	mockTokenGen.EXPECT().
		GenerateTokenPair(gomock.Any()).
		Return("", "", errors.New("token generation failed"))

	useCase := NewGoogleLoginUseCase(mockRepo, mockOAuth, mockTokenGen, "test-client-id")

	result, err := useCase.Execute(ctx, "valid-token")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to generate authentication tokens")
}

func TestGoogleLoginUseCase_RepositorySaveFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockRepository(ctrl)
	mockOAuth := mocks.NewMockOAuthValidator(ctrl)
	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)

	oauthInfo := &ports.OAuthUserInfo{
		UserID:        "google-user-123",
		Email:         "test@example.com",
		EmailVerified: true,
		Name:          "Test User",
		Picture:       "https://example.com/photo.jpg",
	}

	userID, _ := user.NewUserID("google-user-123")

	mockOAuth.EXPECT().
		ValidateToken(ctx, "valid-token", "test-client-id").
		Return(oauthInfo, nil)

	mockRepo.EXPECT().
		FindByID(ctx, userID).
		Return(nil, shared.ErrUserNotFound)

	mockRepo.EXPECT().
		Save(ctx, gomock.Any()).
		Return(errors.New("database error"))

	useCase := NewGoogleLoginUseCase(mockRepo, mockOAuth, mockTokenGen, "test-client-id")

	result, err := useCase.Execute(ctx, "valid-token")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to save new user")
}

func TestGoogleLoginUseCase_RepositoryUpdateFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockRepository(ctrl)
	mockOAuth := mocks.NewMockOAuthValidator(ctrl)
	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)

	// Pre-create existing user
	userID, _ := user.NewUserID("google-user-123")
	email, _ := user.NewEmail("existing@example.com", true)
	profile := user.NewProfile("Old Name", "https://example.com/old.jpg")
	existingUser, _ := user.NewUser(userID, email, profile)

	oauthInfo := &ports.OAuthUserInfo{
		UserID:        "google-user-123",
		Email:         "existing@example.com",
		EmailVerified: true,
		Name:          "Updated Name",
		Picture:       "https://example.com/new.jpg",
	}

	mockOAuth.EXPECT().
		ValidateToken(ctx, "valid-token", "test-client-id").
		Return(oauthInfo, nil)

	mockRepo.EXPECT().
		FindByID(ctx, userID).
		Return(existingUser, nil)

	mockRepo.EXPECT().
		Save(ctx, gomock.Any()).
		Return(errors.New("database update error"))

	useCase := NewGoogleLoginUseCase(mockRepo, mockOAuth, mockTokenGen, "test-client-id")

	result, err := useCase.Execute(ctx, "valid-token")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to update user")
}

func TestGoogleLoginUseCase_RepositoryFindError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockRepository(ctrl)
	mockOAuth := mocks.NewMockOAuthValidator(ctrl)
	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)

	oauthInfo := &ports.OAuthUserInfo{
		UserID:        "google-user-123",
		Email:         "test@example.com",
		EmailVerified: true,
		Name:          "Test User",
		Picture:       "https://example.com/photo.jpg",
	}

	userID, _ := user.NewUserID("google-user-123")

	mockOAuth.EXPECT().
		ValidateToken(ctx, "valid-token", "test-client-id").
		Return(oauthInfo, nil)

	mockRepo.EXPECT().
		FindByID(ctx, userID).
		Return(nil, errors.New("database connection error"))

	useCase := NewGoogleLoginUseCase(mockRepo, mockOAuth, mockTokenGen, "test-client-id")

	result, err := useCase.Execute(ctx, "valid-token")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to check user existence")
}
