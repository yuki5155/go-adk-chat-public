package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/yuki5155/go-google-auth/internal/application/ports"
)

func TestRequireRoot(t *testing.T) {
	tests := []struct {
		name        string
		setupAuth   func(*gin.Context)
		wantStatus  int
		wantAborted bool
	}{
		{
			name: "Root user allowed",
			setupAuth: func(c *gin.Context) {
				c.Set("claims", &ports.TokenClaims{
					UserID: "root-user-1",
					Email:  "root@example.com",
					Name:   "Root User",
					Role:   "root",
				})
			},
			wantStatus:  http.StatusOK,
			wantAborted: false,
		},
		{
			name: "Admin user forbidden",
			setupAuth: func(c *gin.Context) {
				c.Set("claims", &ports.TokenClaims{
					UserID: "admin-1",
					Email:  "admin@example.com",
					Name:   "Admin User",
					Role:   "admin",
				})
			},
			wantStatus:  http.StatusForbidden,
			wantAborted: true,
		},
		{
			name: "Regular user forbidden",
			setupAuth: func(c *gin.Context) {
				c.Set("claims", &ports.TokenClaims{
					UserID: "user-1",
					Email:  "user@example.com",
					Name:   "Regular User",
					Role:   "user",
				})
			},
			wantStatus:  http.StatusForbidden,
			wantAborted: true,
		},
		{
			name: "No claims in context",
			setupAuth: func(c *gin.Context) {
				// No claims set
			},
			wantStatus:  http.StatusUnauthorized,
			wantAborted: true,
		},
		{
			name: "Invalid claims type",
			setupAuth: func(c *gin.Context) {
				c.Set("claims", "invalid-claims")
			},
			wantStatus:  http.StatusUnauthorized,
			wantAborted: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.Use(func(c *gin.Context) {
				tt.setupAuth(c)
				c.Next()
			})
			router.Use(RequireRoot())
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantAborted {
				assert.NotContains(t, w.Body.String(), "success")
			} else {
				assert.Contains(t, w.Body.String(), "success")
			}
		})
	}
}
