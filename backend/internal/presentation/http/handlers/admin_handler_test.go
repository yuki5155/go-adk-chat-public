package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAdminHandler_GetDashboard(t *testing.T) {
	handler := NewAdminHandler()

	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Simulate auth middleware setting user info
	router.Use(func(c *gin.Context) {
		c.Set("userID", "admin-123")
		c.Set("email", "admin@example.com")
		c.Set("name", "Admin User")
		c.Next()
	})

	router.GET("/api/admin/dashboard", handler.GetDashboard)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/admin/dashboard", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Welcome to the admin dashboard")
	assert.Contains(t, w.Body.String(), "admin-123")
	assert.Contains(t, w.Body.String(), "admin@example.com")
	assert.Contains(t, w.Body.String(), "Admin User")
	assert.Contains(t, w.Body.String(), `"isRoot":true`)
}
