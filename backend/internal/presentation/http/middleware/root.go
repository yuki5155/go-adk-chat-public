package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yuki5155/go-google-auth/internal/application/ports"
	"github.com/yuki5155/go-google-auth/internal/domain/user"
)

// RequireRoot middleware checks if the authenticated user has root privileges
// Must be used after the Auth middleware
func RequireRoot() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get claims from context (set by Auth middleware)
		claimsInterface, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "User not authenticated",
			})
			c.Abort()
			return
		}

		claims, ok := claimsInterface.(*ports.TokenClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Invalid authentication",
			})
			c.Abort()
			return
		}

		// Check if user has root role from JWT claims
		if claims.Role != user.RoleRoot.String() {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "forbidden",
				"message": "Root privileges required",
			})
			c.Abort()
			return
		}

		// User is root, continue
		c.Next()
	}
}
