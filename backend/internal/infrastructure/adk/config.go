package adk

import (
	"os"
	"strconv"
)

// Config holds ADK configuration
type Config struct {
	// Gemini API configuration
	APIKey string
	Model  string

	// Model parameters
	Temperature float32
	MaxTokens   int32
	TopP        float32

	// Application settings
	AppName     string
	AppVersion  string
	Description string
}

// NewConfigFromEnv creates a Config from environment variables
func NewConfigFromEnv() *Config {
	return &Config{
		APIKey:      getEnv("GOOGLE_AI_API_KEY", ""),
		Model:       getEnv("GEMINI_MODEL", "gemini-2.0-flash"),
		Temperature: getEnvFloat32("GEMINI_TEMPERATURE", 0.7),
		MaxTokens:   getEnvInt32("GEMINI_MAX_TOKENS", 8192),
		TopP:        getEnvFloat32("GEMINI_TOP_P", 0.95),
		AppName:     getEnv("ADK_APP_NAME", "go-adk-chat"),
		AppVersion:  getEnv("ADK_APP_VERSION", "1.0.0"),
		Description: getEnv("ADK_APP_DESCRIPTION", "AI Chat Assistant with Memory"),
	}
}

// IsValid checks if the configuration is valid
func (c *Config) IsValid() bool {
	return c.APIKey != "" && c.Model != ""
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt32(key string, defaultValue int32) int32 {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.ParseInt(value, 10, 32); err == nil {
			return int32(intVal)
		}
	}
	return defaultValue
}

func getEnvFloat32(key string, defaultValue float32) float32 {
	if value := os.Getenv(key); value != "" {
		if floatVal, err := strconv.ParseFloat(value, 32); err == nil {
			return float32(floatVal)
		}
	}
	return defaultValue
}
