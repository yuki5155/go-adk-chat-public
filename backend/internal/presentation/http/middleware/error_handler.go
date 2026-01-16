package middleware

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yuki5155/go-google-auth/internal/domain/role"
	"github.com/yuki5155/go-google-auth/internal/domain/shared"
)

// ErrorHandler middleware handles errors consistently across all handlers
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				var err error
				switch v := r.(type) {
				case error:
					err = v
				case string:
					err = errors.New(v)
				default:
					err = errors.New("unknown panic")
				}

				log.Printf("Panic recovered: %v", err)
				handleError(c, shared.NewInternalError("INTERNAL_ERROR", "An unexpected error occurred", err))
				c.Abort()
			}
		}()

		c.Next()

		// Handle errors set by handlers via c.Error()
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			handleError(c, err)
		}
	}
}

// handleError converts errors to appropriate HTTP responses
func handleError(c *gin.Context, err error) {
	// Check if it's already an AppError
	var appErr *shared.AppError
	if errors.As(err, &appErr) {
		c.JSON(appErr.StatusCode, gin.H{
			"success": false,
			"error": gin.H{
				"code":    appErr.Code,
				"message": appErr.Message,
			},
		})
		return
	}

	// Map domain errors to HTTP responses
	statusCode, code, message := mapDomainError(err)
	c.JSON(statusCode, gin.H{
		"success": false,
		"error": gin.H{
			"code":    code,
			"message": message,
		},
	})
}

// mapDomainError maps domain errors to HTTP status codes and error codes
func mapDomainError(err error) (statusCode int, code string, message string) {
	switch {
	// Role domain errors
	case errors.Is(err, role.ErrUserRoleNotFound):
		return http.StatusNotFound, "NOT_FOUND", "User role not found"
	case errors.Is(err, role.ErrRoleRequestNotFound):
		return http.StatusNotFound, "NOT_FOUND", "Role request not found"
	case errors.Is(err, role.ErrDuplicateRoleRequest):
		return http.StatusConflict, "DUPLICATE_REQUEST", "You already have a pending role request"
	case errors.Is(err, role.ErrRequestAlreadyProcessed):
		return http.StatusConflict, "ALREADY_PROCESSED", "Request has already been processed"
	case errors.Is(err, role.ErrInvalidRole):
		return http.StatusBadRequest, "INVALID_REQUEST", "Invalid role type"

	// User/shared errors
	case errors.Is(err, shared.ErrUnauthorized):
		return http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated"
	case errors.Is(err, shared.ErrInvalidEmail):
		return http.StatusBadRequest, "INVALID_REQUEST", "Invalid email format"
	case errors.Is(err, shared.ErrUserNotFound):
		return http.StatusNotFound, "NOT_FOUND", "User not found"

	// Default error
	default:
		log.Printf("Unmapped error: %v", err)
		return http.StatusInternalServerError, "INTERNAL_ERROR", "An internal error occurred"
	}
}
