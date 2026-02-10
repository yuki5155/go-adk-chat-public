package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/yuki5155/go-google-auth/internal/application/ports"
	"github.com/yuki5155/go-google-auth/internal/domain/role"
	"github.com/yuki5155/go-google-auth/internal/domain/shared"
)

// RequireSubscriber creates middleware that allows subscriber, admin, and root users
func RequireSubscriber(roleRepo role.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get claims from context (set by Auth middleware)
		claimsInterface, exists := c.Get("claims")
		if !exists {
			_ = c.Error(shared.NewUnauthorizedError("UNAUTHORIZED", "User not authenticated", nil))
			c.Abort()
			return
		}

		claims, ok := claimsInterface.(*ports.TokenClaims)
		if !ok {
			_ = c.Error(shared.NewUnauthorizedError("UNAUTHORIZED", "Invalid authentication", nil))
			c.Abort()
			return
		}

		// Root users always have access
		if claims.Role == "root" {
			c.Next()
			return
		}

		// Admin users have access
		if claims.Role == "admin" {
			c.Next()
			return
		}

		// Check if user is a subscriber
		if claims.Role == "subscriber" {
			c.Next()
			return
		}

		// Check database for subscriber role (in case JWT is outdated)
		userRole, err := roleRepo.GetUserRole(c.Request.Context(), claims.UserID)
		if err == nil && userRole != nil {
			roleStr := userRole.Role().String()
			if roleStr == "subscriber" || roleStr == "admin" || roleStr == "root" {
				c.Next()
				return
			}
		}

		_ = c.Error(shared.NewForbiddenError("FORBIDDEN", "Subscriber access required", nil))
		c.Abort()
	}
}
