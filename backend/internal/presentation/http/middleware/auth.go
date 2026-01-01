package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yuki5155/go-google-auth/internal/application/ports"
)

// Auth creates a middleware for JWT authentication using ports.TokenGenerator
func Auth(tokenGen ports.TokenGenerator) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Cookie("access_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Access token not found",
			})
			c.Abort()
			return
		}

		claims, err := tokenGen.ValidateAccessToken(accessToken)
		if err != nil {
			if ports.IsTokenExpired(err) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error":   "token_expired",
					"message": "Access token has expired",
				})
				c.Abort()
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "invalid_token",
				"message": "Invalid access token",
			})
			c.Abort()
			return
		}

		// Set user information in context for handlers to use
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("name", claims.Name)
		c.Set("picture", claims.Picture)
		c.Set("claims", claims)

		c.Next()
	}
}

// OptionalAuth creates a middleware that validates JWT if present but doesn't require it
func OptionalAuth(tokenGen ports.TokenGenerator) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Cookie("access_token")
		if err != nil {
			// No token, but that's okay - continue without authentication
			c.Next()
			return
		}

		claims, err := tokenGen.ValidateAccessToken(accessToken)
		if err != nil {
			// Invalid token, but still continue - user is just not authenticated
			c.Next()
			return
		}

		// Set user information in context
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("name", claims.Name)
		c.Set("picture", claims.Picture)
		c.Set("claims", claims)
		c.Set("authenticated", true)

		c.Next()
	}
}
