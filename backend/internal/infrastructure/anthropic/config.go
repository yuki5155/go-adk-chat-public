package anthropic

import (
	"os"
	"strconv"
)

// Config holds Anthropic configuration
type Config struct {
	apiKey      string
	Model       string
	Temperature float32
	MaxTokens   int64
	TopP        float32
	AppName     string
}

// NewConfigFromEnv creates a Config from environment variables
func NewConfigFromEnv() *Config {
	return &Config{
		apiKey:      getEnv("ANTHROPIC_API_KEY", ""),
		Model:       getEnv("ANTHROPIC_MODEL", "claude-sonnet-4-6"),
		Temperature: getEnvFloat32("ANTHROPIC_TEMPERATURE", 0.7),
		MaxTokens:   getEnvInt64("ANTHROPIC_MAX_TOKENS", 8192),
		TopP:        getEnvFloat32("ANTHROPIC_TOP_P", 0.95),
		AppName:     getEnv("ADK_APP_NAME", "go-adk-chat"),
	}
}

// GetAPIKey returns the API key
func (c *Config) GetAPIKey() string {
	return c.apiKey
}

// IsValid checks if the configuration is valid
func (c *Config) IsValid() bool {
	return c.apiKey != "" && c.Model != ""
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intVal
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
