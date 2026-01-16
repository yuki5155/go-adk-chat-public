package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/yuki5155/go-google-auth/internal/domain/role"
	"github.com/yuki5155/go-google-auth/internal/domain/shared"
)

func TestErrorHandler_DomainErrors(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		wantStatusCode int
		wantErrorCode  string
	}{
		{
			name:           "User role not found",
			err:            role.ErrUserRoleNotFound,
			wantStatusCode: http.StatusNotFound,
			wantErrorCode:  "NOT_FOUND",
		},
		{
			name:           "Role request not found",
			err:            role.ErrRoleRequestNotFound,
			wantStatusCode: http.StatusNotFound,
			wantErrorCode:  "NOT_FOUND",
		},
		{
			name:           "Duplicate role request",
			err:            role.ErrDuplicateRoleRequest,
			wantStatusCode: http.StatusConflict,
			wantErrorCode:  "DUPLICATE_REQUEST",
		},
		{
			name:           "Request already processed",
			err:            role.ErrRequestAlreadyProcessed,
			wantStatusCode: http.StatusConflict,
			wantErrorCode:  "ALREADY_PROCESSED",
		},
		{
			name:           "Invalid role",
			err:            role.ErrInvalidRole,
			wantStatusCode: http.StatusBadRequest,
			wantErrorCode:  "INVALID_REQUEST",
		},
		{
			name:           "Unauthorized",
			err:            shared.ErrUnauthorized,
			wantStatusCode: http.StatusUnauthorized,
			wantErrorCode:  "UNAUTHORIZED",
		},
		{
			name:           "Invalid email",
			err:            shared.ErrInvalidEmail,
			wantStatusCode: http.StatusBadRequest,
			wantErrorCode:  "INVALID_REQUEST",
		},
		{
			name:           "User not found",
			err:            shared.ErrUserNotFound,
			wantStatusCode: http.StatusNotFound,
			wantErrorCode:  "NOT_FOUND",
		},
		{
			name:           "Unmapped error",
			err:            errors.New("some random error"),
			wantStatusCode: http.StatusInternalServerError,
			wantErrorCode:  "INTERNAL_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.Use(ErrorHandler())
			router.GET("/test", func(c *gin.Context) {
				_ = c.Error(tt.err)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatusCode, w.Code)
			assert.Contains(t, w.Body.String(), tt.wantErrorCode)
			assert.Contains(t, w.Body.String(), `"success":false`)
		})
	}
}

func TestErrorHandler_AppError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(ErrorHandler())
	router.GET("/test", func(c *gin.Context) {
		appErr := shared.NewBadRequestError("CUSTOM_CODE", "Custom message", errors.New("underlying error"))
		_ = c.Error(appErr)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "CUSTOM_CODE")
	assert.Contains(t, w.Body.String(), "Custom message")
	assert.Contains(t, w.Body.String(), `"success":false`)
}

func TestErrorHandler_Panic(t *testing.T) {
	tests := []struct {
		name          string
		panicValue    interface{}
		wantErrorCode string
	}{
		{
			name:          "Panic with error type",
			panicValue:    errors.New("panic error"),
			wantErrorCode: "INTERNAL_ERROR",
		},
		{
			name:          "Panic with string type",
			panicValue:    "panic string",
			wantErrorCode: "INTERNAL_ERROR",
		},
		{
			name:          "Panic with unknown type",
			panicValue:    123,
			wantErrorCode: "INTERNAL_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.Use(ErrorHandler())
			router.GET("/test", func(c *gin.Context) {
				panic(tt.panicValue)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusInternalServerError, w.Code)
			assert.Contains(t, w.Body.String(), tt.wantErrorCode)
			assert.Contains(t, w.Body.String(), `"success":false`)
		})
	}
}

func TestErrorHandler_NoError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(ErrorHandler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}
