package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	adminapp "github.com/yuki5155/go-google-auth/internal/application/admin"
	"github.com/yuki5155/go-google-auth/internal/application/ports"
	"github.com/yuki5155/go-google-auth/internal/domain/user"
)

// RequireAdmin middleware checks if the authenticated user has admin or root privileges
// Must be used after the Auth middleware
func RequireAdmin(checkRoleUC *adminapp.CheckUserRoleUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
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
			c.Next()
			return
		}

		userRole, err := checkRoleUC.Execute(c.Request.Context(), claims.UserID)
		if err != nil || userRole == nil {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "forbidden",
				"message": "Admin privileges required",
			})
			c.Abort()
			return
		}

		if userRole.Role() != user.RoleAdmin {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "forbidden",
				"message": "Admin privileges required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
