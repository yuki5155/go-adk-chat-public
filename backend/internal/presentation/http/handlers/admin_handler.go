package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminHandler handles admin-only endpoints
type AdminHandler struct {
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

// GetDashboard returns admin dashboard information
func (h *AdminHandler) GetDashboard(c *gin.Context) {
	// Get user info from context (set by Auth middleware)
	userID, _ := c.Get("userID")
	email, _ := c.Get("email")
	name, _ := c.Get("name")

	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to the admin dashboard",
		"user": gin.H{
			"id":     userID,
			"email":  email,
			"name":   name,
			"isRoot": true,
		},
	})
}
