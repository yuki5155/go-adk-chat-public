package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/yuki5155/go-google-auth/internal/infrastructure/config"
)

func TestCORS(t *testing.T) {
	tests := []struct {
		name                   string
		allowedOrigins         []string
		requestOrigin          string
		requestMethod          string
		wantOrigin             string
		wantStatus             int
		wantAllowCredentials   string
		wantAllowMethods       string
	}{
		{
			name:                 "Allowed origin",
			allowedOrigins:       []string{"https://example.com", "https://app.example.com"},
			requestOrigin:        "https://example.com",
			requestMethod:        "GET",
			wantOrigin:           "https://example.com",
			wantStatus:           http.StatusOK,
			wantAllowCredentials: "true",
			wantAllowMethods:     "POST, OPTIONS, GET, PUT, DELETE",
		},
		{
			name:                 "Wildcard origin",
			allowedOrigins:       []string{"*"},
			requestOrigin:        "https://any-origin.com",
			requestMethod:        "POST",
			wantOrigin:           "https://any-origin.com",
			wantStatus:           http.StatusOK,
			wantAllowCredentials: "true",
			wantAllowMethods:     "POST, OPTIONS, GET, PUT, DELETE",
		},
		{
			name:                 "Unallowed origin - uses default",
			allowedOrigins:       []string{"https://example.com"},
			requestOrigin:        "https://malicious.com",
			requestMethod:        "GET",
			wantOrigin:           "https://example.com",
			wantStatus:           http.StatusOK,
			wantAllowCredentials: "true",
			wantAllowMethods:     "POST, OPTIONS, GET, PUT, DELETE",
		},
		{
			name:                 "OPTIONS preflight request",
			allowedOrigins:       []string{"https://example.com"},
			requestOrigin:        "https://example.com",
			requestMethod:        "OPTIONS",
			wantOrigin:           "https://example.com",
			wantStatus:           http.StatusNoContent,
			wantAllowCredentials: "true",
			wantAllowMethods:     "POST, OPTIONS, GET, PUT, DELETE",
		},
		{
			name:                 "Empty allowed origins",
			allowedOrigins:       []string{},
			requestOrigin:        "https://example.com",
			requestMethod:        "GET",
			wantOrigin:           "",
			wantStatus:           http.StatusOK,
			wantAllowCredentials: "true",
			wantAllowMethods:     "POST, OPTIONS, GET, PUT, DELETE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				AllowedOrigins: tt.allowedOrigins,
			}

			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.Use(CORS(cfg))
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})
			router.POST("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})
			router.OPTIONS("/test", func(c *gin.Context) {
				// OPTIONS handled by middleware
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.requestMethod, "/test", nil)
			req.Header.Set("Origin", tt.requestOrigin)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantOrigin != "" {
				assert.Equal(t, tt.wantOrigin, w.Header().Get("Access-Control-Allow-Origin"))
			} else {
				// If no origin expected, the header might not be set
				origin := w.Header().Get("Access-Control-Allow-Origin")
				assert.True(t, origin == "" || origin == tt.requestOrigin, "Unexpected origin header: %s", origin)
			}

			assert.Equal(t, tt.wantAllowCredentials, w.Header().Get("Access-Control-Allow-Credentials"))
			assert.Equal(t, tt.wantAllowMethods, w.Header().Get("Access-Control-Allow-Methods"))
			assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Headers"))
		})
	}
}
