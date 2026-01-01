package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/yuki5155/go-google-auth/internal/application/ports"
	"github.com/yuki5155/go-google-auth/internal/mocks"
)

func TestRefreshTokenUseCase_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)

	mockTokenGen.EXPECT().
		RefreshAccessToken("valid-refresh-token").
		Return("new-access-token", nil)

	useCase := NewRefreshTokenUseCase(mockTokenGen)

	result, err := useCase.Execute(ctx, "valid-refresh-token")

	require.NoError(t, err)
	assert.Equal(t, "new-access-token", result.AccessToken)
	assert.Equal(t, "Token refreshed successfully", result.Message)
}

func TestRefreshTokenUseCase_EmptyToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)

	useCase := NewRefreshTokenUseCase(mockTokenGen)

	result, err := useCase.Execute(ctx, "")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "refresh token is required")
}

func TestRefreshTokenUseCase_InvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)

	mockTokenGen.EXPECT().
		RefreshAccessToken("invalid-token").
		Return("", errors.New("invalid refresh token"))

	useCase := NewRefreshTokenUseCase(mockTokenGen)

	result, err := useCase.Execute(ctx, "invalid-token")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid refresh token")
}

func TestRefreshTokenUseCase_ExpiredToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)

	mockTokenGen.EXPECT().
		RefreshAccessToken("expired-token").
		Return("", ports.ErrExpiredToken)

	useCase := NewRefreshTokenUseCase(mockTokenGen)

	result, err := useCase.Execute(ctx, "expired-token")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ports.ErrExpiredToken, err)
}
