package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuki5155/go-google-auth/internal/application/ports"
)

const testSecretKey = "test-secret-key-for-jwt-testing"

func TestNewService(t *testing.T) {
	service := NewService(testSecretKey)

	assert.NotNil(t, service)
	assert.Equal(t, []byte(testSecretKey), service.secretKey)
	assert.Equal(t, 15*time.Minute, service.accessTokenExpiry)
	assert.Equal(t, 7*24*time.Hour, service.refreshTokenExpiry)
}

func TestGenerateTokenPair(t *testing.T) {
	service := NewService(testSecretKey)
	user := ports.UserInfo{
		UserID:  "user123",
		Email:   "test@example.com",
		Name:    "Test User",
		Picture: "https://example.com/picture.jpg",
	}

	accessToken, refreshToken, err := service.GenerateTokenPair(user)

	require.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)
	assert.NotEqual(t, accessToken, refreshToken)

	// Validate access token
	accessClaims, err := service.ValidateAccessToken(accessToken)
	require.NoError(t, err)
	assert.Equal(t, user.UserID, accessClaims.UserID)
	assert.Equal(t, user.Email, accessClaims.Email)
	assert.Equal(t, user.Name, accessClaims.Name)
	assert.Equal(t, user.Picture, accessClaims.Picture)
}

func TestValidateAccessToken_ValidToken(t *testing.T) {
	service := NewService(testSecretKey)
	user := ports.UserInfo{
		UserID:  "user123",
		Email:   "test@example.com",
		Name:    "Test User",
		Picture: "https://example.com/picture.jpg",
	}

	accessToken, _, err := service.GenerateTokenPair(user)
	require.NoError(t, err)

	claims, err := service.ValidateAccessToken(accessToken)

	require.NoError(t, err)
	assert.Equal(t, user.UserID, claims.UserID)
	assert.Equal(t, user.Email, claims.Email)
	assert.Equal(t, user.Name, claims.Name)
	assert.Equal(t, user.Picture, claims.Picture)
}

func TestValidateAccessToken_InvalidToken(t *testing.T) {
	service := NewService(testSecretKey)

	tests := []struct {
		name  string
		token string
	}{
		{name: "empty token", token: ""},
		{name: "malformed token", token: "not.a.valid.jwt"},
		{name: "random string", token: "totally-invalid-token"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := service.ValidateAccessToken(tt.token)

			assert.Error(t, err)
			assert.Nil(t, claims)
		})
	}
}

func TestValidateAccessToken_RefreshTokenRejected(t *testing.T) {
	service := NewService(testSecretKey)
	user := ports.UserInfo{
		UserID: "user123",
		Email:  "test@example.com",
	}

	_, refreshToken, err := service.GenerateTokenPair(user)
	require.NoError(t, err)

	// Try to validate refresh token as access token
	claims, err := service.ValidateAccessToken(refreshToken)

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestValidateAccessToken_ExpiredToken(t *testing.T) {
	service := NewService(testSecretKey)
	service.accessTokenExpiry = -1 * time.Hour // Expired

	user := ports.UserInfo{
		UserID: "user123",
		Email:  "test@example.com",
	}

	token, err := service.generateToken(user, "access", service.accessTokenExpiry)
	require.NoError(t, err)

	service.accessTokenExpiry = 15 * time.Minute // Reset

	claims, err := service.ValidateAccessToken(token)

	assert.Error(t, err)
	assert.Equal(t, ports.ErrExpiredToken, err)
	assert.Nil(t, claims)
}

func TestRefreshAccessToken_ValidRefreshToken(t *testing.T) {
	service := NewService(testSecretKey)
	user := ports.UserInfo{
		UserID:  "user123",
		Email:   "test@example.com",
		Name:    "Test User",
		Picture: "https://example.com/picture.jpg",
	}

	_, refreshToken, err := service.GenerateTokenPair(user)
	require.NoError(t, err)

	newAccessToken, err := service.RefreshAccessToken(refreshToken)

	require.NoError(t, err)
	assert.NotEmpty(t, newAccessToken)

	// Validate the new access token
	claims, err := service.ValidateAccessToken(newAccessToken)
	require.NoError(t, err)
	assert.Equal(t, user.UserID, claims.UserID)
	assert.Equal(t, user.Email, claims.Email)
	assert.Equal(t, user.Name, claims.Name)
	assert.Equal(t, user.Picture, claims.Picture)
}

func TestRefreshAccessToken_InvalidRefreshToken(t *testing.T) {
	service := NewService(testSecretKey)

	tests := []struct {
		name  string
		token string
	}{
		{name: "empty token", token: ""},
		{name: "invalid token", token: "invalid-token"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newAccessToken, err := service.RefreshAccessToken(tt.token)

			assert.Error(t, err)
			assert.Empty(t, newAccessToken)
		})
	}
}

func TestRefreshAccessToken_AccessTokenRejected(t *testing.T) {
	service := NewService(testSecretKey)
	user := ports.UserInfo{
		UserID: "user123",
		Email:  "test@example.com",
	}

	accessToken, _, err := service.GenerateTokenPair(user)
	require.NoError(t, err)

	newAccessToken, err := service.RefreshAccessToken(accessToken)

	assert.Error(t, err)
	assert.Empty(t, newAccessToken)
}

func TestRefreshAccessToken_ExpiredRefreshToken(t *testing.T) {
	service := NewService(testSecretKey)
	user := ports.UserInfo{
		UserID: "user123",
		Email:  "test@example.com",
	}

	service.refreshTokenExpiry = -1 * time.Hour
	expiredRefreshToken, err := service.generateToken(user, "refresh", service.refreshTokenExpiry)
	require.NoError(t, err)

	service.refreshTokenExpiry = 7 * 24 * time.Hour

	newAccessToken, err := service.RefreshAccessToken(expiredRefreshToken)

	assert.Error(t, err)
	assert.Equal(t, ports.ErrExpiredToken, err)
	assert.Empty(t, newAccessToken)
}

func TestGetAccessTokenExpiry(t *testing.T) {
	service := NewService(testSecretKey)

	expiry := service.GetAccessTokenExpiry()

	assert.Equal(t, 900, expiry) // 15 minutes = 900 seconds
}

func TestGetRefreshTokenExpiry(t *testing.T) {
	service := NewService(testSecretKey)

	expiry := service.GetRefreshTokenExpiry()

	assert.Equal(t, 604800, expiry) // 7 days = 604800 seconds
}

func TestTokenLifecycle(t *testing.T) {
	service := NewService(testSecretKey)
	user := ports.UserInfo{
		UserID:  "user123",
		Email:   "test@example.com",
		Name:    "Test User",
		Picture: "https://example.com/picture.jpg",
	}

	// 1. Generate token pair
	accessToken, refreshToken, err := service.GenerateTokenPair(user)
	require.NoError(t, err)

	// 2. Validate access token works
	accessClaims, err := service.ValidateAccessToken(accessToken)
	require.NoError(t, err)
	assert.Equal(t, user.UserID, accessClaims.UserID)

	// 3. Refresh access token
	newAccessToken, err := service.RefreshAccessToken(refreshToken)
	require.NoError(t, err)

	// 4. Validate new access token
	newAccessClaims, err := service.ValidateAccessToken(newAccessToken)
	require.NoError(t, err)
	assert.Equal(t, user.UserID, newAccessClaims.UserID)
	assert.Equal(t, user.Email, newAccessClaims.Email)

	// 5. Old access token should still be valid
	oldAccessClaims, err := service.ValidateAccessToken(accessToken)
	require.NoError(t, err)
	assert.Equal(t, user.UserID, oldAccessClaims.UserID)
}

func TestUserInfoWithEmptyFields(t *testing.T) {
	service := NewService(testSecretKey)
	user := ports.UserInfo{
		UserID: "user123",
		Email:  "test@example.com",
		// Name and Picture empty
	}

	accessToken, refreshToken, err := service.GenerateTokenPair(user)

	require.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)

	claims, err := service.ValidateAccessToken(accessToken)
	require.NoError(t, err)
	assert.Equal(t, user.UserID, claims.UserID)
	assert.Equal(t, user.Email, claims.Email)
	assert.Empty(t, claims.Name)
	assert.Empty(t, claims.Picture)
}

func TestWrongSigningKey(t *testing.T) {
	service := NewService(testSecretKey)
	differentService := NewService("different-secret-key")

	user := ports.UserInfo{
		UserID: "user123",
		Email:  "test@example.com",
	}

	// Generate token with different secret
	token, _, err := differentService.GenerateTokenPair(user)
	require.NoError(t, err)

	// Try to validate with original service
	claims, err := service.ValidateAccessToken(token)

	assert.Error(t, err)
	assert.Nil(t, claims)
}
