package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yuki5155/go-google-auth/internal/application/auth"
	"github.com/yuki5155/go-google-auth/internal/application/dto"
	"github.com/yuki5155/go-google-auth/internal/application/ports"
	"github.com/yuki5155/go-google-auth/internal/domain/shared"
	"github.com/yuki5155/go-google-auth/internal/infrastructure/config"
)

// AuthHandler handles HTTP authentication requests (thin controller)
type AuthHandler struct {
	googleLoginUC    *auth.GoogleLoginUseCase
	refreshTokenUC   *auth.RefreshTokenUseCase
	getCurrentUserUC *auth.GetCurrentUserUseCase
	logoutUC         *auth.LogoutUseCase
	tokenGenerator   ports.TokenGenerator
	config           *config.Config
}

// NewAuthHandler creates a new thin AuthHandler
func NewAuthHandler(
	googleLoginUC *auth.GoogleLoginUseCase,
	refreshTokenUC *auth.RefreshTokenUseCase,
	getCurrentUserUC *auth.GetCurrentUserUseCase,
	logoutUC *auth.LogoutUseCase,
	tokenGenerator ports.TokenGenerator,
	config *config.Config,
) *AuthHandler {
	return &AuthHandler{
		googleLoginUC:    googleLoginUC,
		refreshTokenUC:   refreshTokenUC,
		getCurrentUserUC: getCurrentUserUC,
		logoutUC:         logoutUC,
		tokenGenerator:   tokenGenerator,
		config:           config,
	}
}

// GoogleLogin handles Google OAuth login requests
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	var req dto.GoogleLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Missing or invalid credential",
		})
		return
	}

	result, err := h.googleLoginUC.Execute(c.Request.Context(), req.Credential)
	if err != nil {
		if err == shared.ErrUnverifiedEmail {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unverified_email",
				"message": "Email address is not verified",
			})
			return
		}
		log.Printf("Google login failed: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "authentication_failed",
			"message": "Failed to authenticate with Google",
		})
		return
	}

	h.setAuthCookies(c, result.AccessToken, result.RefreshToken)

	c.JSON(http.StatusOK, gin.H{
		"message": result.Message,
		"user":    result.User,
	})
}

// RefreshToken handles token refresh requests
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "missing_refresh_token",
			"message": "Refresh token not found",
		})
		return
	}

	result, err := h.refreshTokenUC.Execute(c.Request.Context(), refreshToken)
	if err != nil {
		if err == ports.ErrExpiredToken {
			h.clearAuthCookies(c)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "refresh_token_expired",
				"message": "Refresh token has expired, please login again",
			})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "invalid_refresh_token",
			"message": "Invalid refresh token",
		})
		return
	}

	h.setAccessTokenCookie(c, result.AccessToken)

	c.JSON(http.StatusOK, gin.H{
		"message": result.Message,
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	result := h.logoutUC.Execute(c.Request.Context())

	h.clearAuthCookies(c)

	c.JSON(http.StatusOK, gin.H{
		"message": result.Message,
	})
}

// GetCurrentUser returns the current authenticated user's information
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	// Get claims from context (set by auth middleware)
	claimsInterface, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "User not authenticated",
		})
		return
	}

	claims, ok := claimsInterface.(*ports.TokenClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "Invalid authentication",
		})
		return
	}

	// Return user info directly from JWT claims
	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":      claims.UserID,
			"email":   claims.Email,
			"name":    claims.Name,
			"picture": claims.Picture,
		},
	})
}

// setAuthCookies sets both access and refresh token cookies
func (h *AuthHandler) setAuthCookies(c *gin.Context, accessToken, refreshToken string) {
	secure := h.config.IsProduction()

	c.SetCookie(
		"access_token",
		accessToken,
		h.tokenGenerator.GetAccessTokenExpiry(),
		"/",
		"",
		secure,
		true, // HttpOnly
	)

	c.SetCookie(
		"refresh_token",
		refreshToken,
		h.tokenGenerator.GetRefreshTokenExpiry(),
		"/",
		"",
		secure,
		true, // HttpOnly
	)
}

// setAccessTokenCookie sets only the access token cookie
func (h *AuthHandler) setAccessTokenCookie(c *gin.Context, accessToken string) {
	secure := h.config.IsProduction()

	c.SetCookie(
		"access_token",
		accessToken,
		h.tokenGenerator.GetAccessTokenExpiry(),
		"/",
		"",
		secure,
		true, // HttpOnly
	)
}

// clearAuthCookies removes authentication cookies
func (h *AuthHandler) clearAuthCookies(c *gin.Context) {
	secure := h.config.IsProduction()

	c.SetCookie("access_token", "", -1, "/", "", secure, true)
	c.SetCookie("refresh_token", "", -1, "/", "", secure, true)
}
