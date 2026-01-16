package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/yuki5155/go-google-auth/internal/application/auth"
	"github.com/yuki5155/go-google-auth/internal/application/ports"
	"github.com/yuki5155/go-google-auth/internal/infrastructure/config"
	"github.com/yuki5155/go-google-auth/internal/mocks"
)

func TestAuthHandler_GoogleLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOAuth := mocks.NewMockOAuthValidator(ctrl)
	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)
	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)

	oauthInfo := &ports.OAuthUserInfo{
		UserID:        "google-user-123",
		Email:         "test@example.com",
		EmailVerified: true,
		Name:          "Test User",
		Picture:       "https://example.com/photo.jpg",
	}

	mockOAuth.EXPECT().
		ValidateToken(gomock.Any(), "valid-google-token", "test-client-id").
		Return(oauthInfo, nil)

	// Mock role repository - no role assigned
	mockRoleRepo.EXPECT().
		GetUserRole(gomock.Any(), "google-user-123").
		Return(nil, nil)

	mockTokenGen.EXPECT().
		GenerateTokenPair(gomock.Any()).
		Return("mock-access-token", "mock-refresh-token", nil)

	// Mock cookie expiry methods
	mockTokenGen.EXPECT().
		GetAccessTokenExpiry().
		Return(3600)

	mockTokenGen.EXPECT().
		GetRefreshTokenExpiry().
		Return(604800)

	googleLoginUC := auth.NewGoogleLoginUseCase(mockOAuth, mockTokenGen, mockRoleRepo, "test-client-id", "root@example.com")

	cfg := &config.Config{
		Environment: "development",
	}
	handler := NewAuthHandler(googleLoginUC, nil, nil, mockTokenGen, cfg)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/auth/google", handler.GoogleLogin)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/auth/google",
		strings.NewReader(`{"credential": "valid-google-token"}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Check cookies are set
	cookies := w.Result().Cookies()
	var foundAccessToken, foundRefreshToken bool
	for _, cookie := range cookies {
		if cookie.Name == "access_token" {
			foundAccessToken = true
		}
		if cookie.Name == "refresh_token" {
			foundRefreshToken = true
		}
	}
	assert.True(t, foundAccessToken, "access_token cookie should be set")
	assert.True(t, foundRefreshToken, "refresh_token cookie should be set")
}

func TestAuthHandler_GoogleLogin_InvalidJSON(t *testing.T) {
	handler := NewAuthHandler(nil, nil, nil, nil, nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/auth/google", handler.GoogleLogin)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/auth/google",
		strings.NewReader(`{invalid json}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_GoogleLogin_UnverifiedEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOAuth := mocks.NewMockOAuthValidator(ctrl)
	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)
	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)

	oauthInfo := &ports.OAuthUserInfo{
		UserID:        "google-user-123",
		Email:         "test@example.com",
		EmailVerified: false,
		Name:          "Test User",
	}

	mockOAuth.EXPECT().
		ValidateToken(gomock.Any(), "valid-token", "test-client-id").
		Return(oauthInfo, nil)

	googleLoginUC := auth.NewGoogleLoginUseCase(mockOAuth, mockTokenGen, mockRoleRepo, "test-client-id", "root@example.com")
	handler := NewAuthHandler(googleLoginUC, nil, nil, nil, nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/auth/google", handler.GoogleLogin)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/auth/google",
		strings.NewReader(`{"credential": "valid-token"}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthHandler_GoogleLogin_AuthenticationFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOAuth := mocks.NewMockOAuthValidator(ctrl)
	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)
	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)

	mockOAuth.EXPECT().
		ValidateToken(gomock.Any(), "invalid-token", "test-client-id").
		Return(nil, errors.New("invalid token"))

	googleLoginUC := auth.NewGoogleLoginUseCase(mockOAuth, mockTokenGen, mockRoleRepo, "test-client-id", "root@example.com")
	handler := NewAuthHandler(googleLoginUC, nil, nil, nil, nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/auth/google", handler.GoogleLogin)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/auth/google",
		strings.NewReader(`{"credential": "invalid-token"}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthHandler_RefreshToken_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)

	mockTokenGen.EXPECT().
		RefreshAccessToken("valid-refresh-token").
		Return("new-access-token", nil)

	// Mock cookie expiry methods
	mockTokenGen.EXPECT().
		GetAccessTokenExpiry().
		Return(3600)

	refreshTokenUC := auth.NewRefreshTokenUseCase(mockTokenGen)
	cfg := &config.Config{
		Environment: "development",
	}
	handler := NewAuthHandler(nil, refreshTokenUC, nil, mockTokenGen, cfg)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/auth/refresh", handler.RefreshToken)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/auth/refresh", nil)
	req.AddCookie(&http.Cookie{
		Name:  "refresh_token",
		Value: "valid-refresh-token",
	})
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthHandler_RefreshToken_MissingCookie(t *testing.T) {
	handler := NewAuthHandler(nil, nil, nil, nil, nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/auth/refresh", handler.RefreshToken)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/auth/refresh", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthHandler_RefreshToken_InvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)

	mockTokenGen.EXPECT().
		RefreshAccessToken("invalid-token").
		Return("", errors.New("invalid refresh token"))

	refreshTokenUC := auth.NewRefreshTokenUseCase(mockTokenGen)
	handler := NewAuthHandler(nil, refreshTokenUC, nil, nil, nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/auth/refresh", handler.RefreshToken)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/auth/refresh", nil)
	req.AddCookie(&http.Cookie{
		Name:  "refresh_token",
		Value: "invalid-token",
	})
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthHandler_Logout_Success(t *testing.T) {
	logoutUC := auth.NewLogoutUseCase()
	cfg := &config.Config{
		Environment: "development",
	}
	handler := NewAuthHandler(nil, nil, logoutUC, nil, cfg)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/auth/logout", handler.Logout)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/auth/logout", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Check that cookies are cleared
	cookies := w.Result().Cookies()
	var accessTokenCleared, refreshTokenCleared bool
	for _, cookie := range cookies {
		if cookie.Name == "access_token" && cookie.MaxAge == -1 {
			accessTokenCleared = true
		}
		if cookie.Name == "refresh_token" && cookie.MaxAge == -1 {
			refreshTokenCleared = true
		}
	}
	assert.True(t, accessTokenCleared, "access_token should be cleared")
	assert.True(t, refreshTokenCleared, "refresh_token should be cleared")
}

func TestAuthHandler_GetCurrentUser_Success(t *testing.T) {
	handler := NewAuthHandler(nil, nil, nil, nil, nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	// Simulate auth middleware setting claims
	router.Use(func(c *gin.Context) {
		c.Set("claims", &ports.TokenClaims{
			UserID: "user-123",
			Email:  "test@example.com",
			Name:   "Test User",
			Role:   "user",
		})
		c.Next()
	})
	router.GET("/api/me", handler.GetCurrentUser)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/me", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthHandler_GetCurrentUser_NoClaimsInContext(t *testing.T) {
	handler := NewAuthHandler(nil, nil, nil, nil, nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/me", handler.GetCurrentUser)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/me", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
