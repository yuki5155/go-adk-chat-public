package router

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yuki5155/go-google-auth/internal/handlers"
	"github.com/yuki5155/go-google-auth/internal/infrastructure/container"
	presentationHandlers "github.com/yuki5155/go-google-auth/internal/presentation/http/handlers"
	"github.com/yuki5155/go-google-auth/internal/presentation/http/middleware"
)

// Setup initializes and configures the Gin router with all routes and middleware
func Setup(c *container.Container) *gin.Engine {
	cfg := c.GetConfig()

	// Set Gin mode
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	r := gin.Default()

	// Apply CORS middleware
	r.Use(middleware.CORS(cfg))

	// Initialize new presentation layer handler
	authHandler := presentationHandlers.NewAuthHandler(
		c.GoogleLoginUseCase,
		c.RefreshTokenUseCase,
		c.GetCurrentUserUseCase,
		c.LogoutUseCase,
		c.TokenGenerator,
		cfg,
	)

	// Initialize old handlers (to be migrated)
	helloHandler := handlers.NewHelloHandler()
	healthHandler := handlers.NewHealthHandler()
	setCookieHandler := handlers.NewCookieHandler(cfg)
	checkCookieHandler := handlers.NewCheckCookieHandler()

	// Public routes
	r.GET(helloHandler.Path, helloHandler.Handle)
	r.GET(healthHandler.Path, healthHandler.Handle)
	r.GET("/health/ready", healthHandler.Handle)

	// Cookie test endpoints
	r.GET(setCookieHandler.Path, setCookieHandler.Handle)
	r.GET(checkCookieHandler.Path, checkCookieHandler.Handle)

	// Auth endpoints (public) - using new presentation layer handlers
	r.POST("/auth/google", authHandler.GoogleLogin)
	r.POST("/auth/refresh", authHandler.RefreshToken)
	r.POST("/auth/logout", authHandler.Logout)

	// Protected routes (require authentication)
	protected := r.Group("/api")
	protected.Use(middleware.Auth(c.TokenGenerator))
	{
		protected.GET("/me", authHandler.GetCurrentUser)
	}

	log.Printf("Router configured (environment: %s)", cfg.Environment)
	log.Printf("Allowed CORS origins: %s", strings.Join(cfg.AllowedOrigins, ", "))
	log.Printf("Google Client ID configured: %v", cfg.GoogleClientID != "")

	return r
}
