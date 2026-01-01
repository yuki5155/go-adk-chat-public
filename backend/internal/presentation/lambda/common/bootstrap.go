package common

import (
	"github.com/gin-gonic/gin"
	"github.com/yuki5155/go-google-auth/internal/infrastructure/config"
	"github.com/yuki5155/go-google-auth/internal/infrastructure/container"
	"github.com/yuki5155/go-google-auth/internal/presentation/http/middleware"
)

// Bootstrap initializes the Gin router and dependency container for Lambda functions
// Returns a configured router and the dependency container
func Bootstrap() (*gin.Engine, *container.Container) {
	// Load configuration
	cfg := config.Load()

	// Create dependency injection container
	c := container.NewContainer(cfg)

	// Set Gin mode based on environment
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router with default middleware
	r := gin.Default()

	// Apply shared CORS middleware
	r.Use(middleware.CORS(cfg))

	return r, c
}
