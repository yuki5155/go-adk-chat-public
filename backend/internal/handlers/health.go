package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	Path string
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{
		Path: "/health",
	}
}

func (h *HealthHandler) Handle(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
