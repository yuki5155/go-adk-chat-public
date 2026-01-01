package config

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"os"
	"strings"
)

type Config struct {
	Port              string
	Environment       string
	AllowedOrigins    []string
	FrontendURL       string
	GoogleClientID    string
	GoogleSecret      string
	GoogleRedirectURL string
	JWTSecret         string
}

func Load() *Config {
	// CORS Allowed Origins - comma separated
	allowedOriginsStr := getEnv("ALLOWED_ORIGINS", "http://localhost:5173")
	allowedOrigins := strings.Split(allowedOriginsStr, ",")
	for i := range allowedOrigins {
		allowedOrigins[i] = strings.TrimSpace(allowedOrigins[i])
	}

	// JWT Secret - generate a random one if not provided (for development only)
	jwtSecret := getEnv("JWT_SECRET", "")
	if jwtSecret == "" {
		jwtSecret = generateRandomSecret()
		log.Println("WARNING: JWT_SECRET not set, using auto-generated secret. Set JWT_SECRET environment variable in production.")
	}

	return &Config{
		Port:              getEnv("PORT", "8080"),
		Environment:       getEnv("GO_ENV", "development"),
		AllowedOrigins:    allowedOrigins,
		FrontendURL:       getEnv("FRONTEND_URL", "http://localhost:5173"),
		GoogleClientID:    getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleSecret:      getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL: getEnv("GOOGLE_REDIRECT_URL", ""),
		JWTSecret:         jwtSecret,
	}
}

// generateRandomSecret generates a random 32-byte secret for development
func generateRandomSecret() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "development-fallback-secret-change-in-prod"
	}
	return base64.StdEncoding.EncodeToString(bytes)
}

// IsProduction returns true if running in production environment
func (c *Config) IsProduction() bool {
	return c.Environment == "production" || c.Environment == "prod"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
