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
	"github.com/yuki5155/go-google-auth/internal/mocks"
)

func TestGoogleLoginUseCase_RegularUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockOAuth := mocks.NewMockOAuthValidator(ctrl)
	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)

	oauthInfo := &ports.OAuthUserInfo{
		UserID:        "google-user-123",
		Email:         "user@example.com",
		EmailVerified: true,
		Name:          "Regular User",
		Picture:       "https://example.com/photo.jpg",
	}

	// Set expectations
	mockOAuth.EXPECT().
		ValidateToken(ctx, "valid-google-token", "test-client-id").
		Return(oauthInfo, nil)

	mockTokenGen.EXPECT().
		GenerateTokenPair(ports.UserInfo{
			UserID:  "google-user-123",
			Email:   "user@example.com",
			Name:    "Regular User",
			Picture: "https://example.com/photo.jpg",
			Role:    "user",
		}).
		Return("mock-access-token", "mock-refresh-token", nil)

	useCase := NewGoogleLoginUseCase(mockOAuth, mockTokenGen, "test-client-id", "root@example.com")

	result, err := useCase.Execute(ctx, "valid-google-token")

	require.NoError(t, err)
	assert.Equal(t, "Login successful", result.Message)
	assert.Equal(t, "google-user-123", result.User.ID)
	assert.Equal(t, "user@example.com", result.User.Email)
	assert.Equal(t, "Regular User", result.User.Name)
	assert.Equal(t, "https://example.com/photo.jpg", result.User.Picture)
	assert.Equal(t, "user", result.User.Role)
	assert.Equal(t, "mock-access-token", result.AccessToken)
	assert.Equal(t, "mock-refresh-token", result.RefreshToken)
}

func TestGoogleLoginUseCase_RootUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockOAuth := mocks.NewMockOAuthValidator(ctrl)
	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)

	oauthInfo := &ports.OAuthUserInfo{
		UserID:        "google-user-456",
		Email:         "root@example.com",
		EmailVerified: true,
		Name:          "Root User",
		Picture:       "https://example.com/root.jpg",
	}

	// Set expectations
	mockOAuth.EXPECT().
		ValidateToken(ctx, "valid-google-token", "test-client-id").
		Return(oauthInfo, nil)

	mockTokenGen.EXPECT().
		GenerateTokenPair(ports.UserInfo{
			UserID:  "google-user-456",
			Email:   "root@example.com",
			Name:    "Root User",
			Picture: "https://example.com/root.jpg",
			Role:    "root",
		}).
		Return("mock-access-token", "mock-refresh-token", nil)

	useCase := NewGoogleLoginUseCase(mockOAuth, mockTokenGen, "test-client-id", "root@example.com")

	result, err := useCase.Execute(ctx, "valid-google-token")

	require.NoError(t, err)
	assert.Equal(t, "Login successful", result.Message)
	assert.Equal(t, "google-user-456", result.User.ID)
	assert.Equal(t, "root@example.com", result.User.Email)
	assert.Equal(t, "Root User", result.User.Name)
	assert.Equal(t, "root", result.User.Role)
	assert.Equal(t, "mock-access-token", result.AccessToken)
	assert.Equal(t, "mock-refresh-token", result.RefreshToken)
}

func TestGoogleLoginUseCase_InvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockOAuth := mocks.NewMockOAuthValidator(ctrl)
	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)

	mockOAuth.EXPECT().
		ValidateToken(ctx, "invalid-token", "test-client-id").
		Return(nil, errors.New("invalid token"))

	useCase := NewGoogleLoginUseCase(mockOAuth, mockTokenGen, "test-client-id", "root@example.com")

	result, err := useCase.Execute(ctx, "invalid-token")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to verify Google ID token")
}

func TestGoogleLoginUseCase_UnverifiedEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockOAuth := mocks.NewMockOAuthValidator(ctrl)
	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)

	oauthInfo := &ports.OAuthUserInfo{
		UserID:        "google-user-123",
		Email:         "unverified@example.com",
		EmailVerified: false,
		Name:          "Unverified User",
		Picture:       "https://example.com/photo.jpg",
	}

	mockOAuth.EXPECT().
		ValidateToken(ctx, "valid-token", "test-client-id").
		Return(oauthInfo, nil)

	useCase := NewGoogleLoginUseCase(mockOAuth, mockTokenGen, "test-client-id", "root@example.com")

	result, err := useCase.Execute(ctx, "valid-token")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, shared.ErrUnverifiedEmail, err)
}

func TestGoogleLoginUseCase_TokenGenerationFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockOAuth := mocks.NewMockOAuthValidator(ctrl)
	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)

	oauthInfo := &ports.OAuthUserInfo{
		UserID:        "google-user-123",
		Email:         "user@example.com",
		EmailVerified: true,
		Name:          "User",
		Picture:       "https://example.com/photo.jpg",
	}

	mockOAuth.EXPECT().
		ValidateToken(ctx, "valid-token", "test-client-id").
		Return(oauthInfo, nil)

	mockTokenGen.EXPECT().
		GenerateTokenPair(gomock.Any()).
		Return("", "", errors.New("token generation failed"))

	useCase := NewGoogleLoginUseCase(mockOAuth, mockTokenGen, "test-client-id", "root@example.com")

	result, err := useCase.Execute(ctx, "valid-token")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to generate authentication tokens")
}
