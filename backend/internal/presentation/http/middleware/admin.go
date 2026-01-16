package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/yuki5155/go-google-auth/internal/application/ports"
	"github.com/yuki5155/go-google-auth/internal/domain/role"
	"github.com/yuki5155/go-google-auth/internal/domain/user"
)

// RequireAdmin middleware checks if the authenticated user has admin or root privileges
// Must be used after the Auth middleware
func RequireAdmin(roleRepo role.Repository) gin.HandlerFunc {
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

		// Check if user is root (from environment variable)
		rootEmail := os.Getenv("ROOT_USER_EMAIL")
		if rootEmail != "" && claims.Email == rootEmail {
			// User is root, continue
			c.Next()
			return
		}

		// Check if user has admin role in database
		userRole, err := roleRepo.GetUserRoleByEmail(c.Request.Context(), claims.Email)
		if err != nil {
			// User not found or database error - forbid access
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "forbidden",
				"message": "Admin privileges required",
			})
			c.Abort()
			return
		}

		// Check if the role is Admin
		if userRole.Role() != user.RoleAdmin {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "forbidden",
				"message": "Admin privileges required",
			})
			c.Abort()
			return
		}

		// User is admin, continue
		c.Next()
	}
}
