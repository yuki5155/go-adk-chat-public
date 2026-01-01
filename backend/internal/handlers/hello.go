package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HelloHandler struct {
	Path string
}

func NewHelloHandler() *HelloHandler {
	return &HelloHandler{
		Path: "/hello",
	}
}

func (h *HelloHandler) Handle(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello World",
	})
}
