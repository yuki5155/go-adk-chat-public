package middleware

import (
	"github.com/gin-gonic/gin"
	adminapp "github.com/yuki5155/go-google-auth/internal/application/admin"
	"github.com/yuki5155/go-google-auth/internal/application/ports"
	"github.com/yuki5155/go-google-auth/internal/domain/shared"
	"github.com/yuki5155/go-google-auth/internal/domain/user"
)

// RequireSubscriber creates middleware that allows subscriber, admin, and root users
func RequireSubscriber(checkRoleUC *adminapp.CheckUserRoleUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		role := user.Role(claims.Role)
		if role.IsSubscriber() {
			c.Next()
			return
		}

		// Check database for subscriber role (in case JWT is outdated)
		userRole, err := checkRoleUC.Execute(c.Request.Context(), claims.UserID)
		if err == nil && userRole != nil && userRole.Role().IsSubscriber() {
			c.Next()
			return
		}

		_ = c.Error(shared.NewForbiddenError("FORBIDDEN", "Subscriber access required", nil))
		c.Abort()
	}
}
