package adk

import (
	"os"
	"testing"
)

func TestNewConfigFromEnv(t *testing.T) {
	// Save original env vars
	origAPIKey := os.Getenv("GOOGLE_AI_API_KEY")
	origModel := os.Getenv("GEMINI_MODEL")
	origTemp := os.Getenv("GEMINI_TEMPERATURE")
	origMaxTokens := os.Getenv("GEMINI_MAX_TOKENS")
	origTopP := os.Getenv("GEMINI_TOP_P")
	origAppName := os.Getenv("ADK_APP_NAME")
	origAppVersion := os.Getenv("ADK_APP_VERSION")
	origDescription := os.Getenv("ADK_APP_DESCRIPTION")

	// Restore env vars after test
	defer func() {
		os.Setenv("GOOGLE_AI_API_KEY", origAPIKey)
		os.Setenv("GEMINI_MODEL", origModel)
		os.Setenv("GEMINI_TEMPERATURE", origTemp)
		os.Setenv("GEMINI_MAX_TOKENS", origMaxTokens)
		os.Setenv("GEMINI_TOP_P", origTopP)
		os.Setenv("ADK_APP_NAME", origAppName)
		os.Setenv("ADK_APP_VERSION", origAppVersion)
		os.Setenv("ADK_APP_DESCRIPTION", origDescription)
	}()

	t.Run("with environment variables set", func(t *testing.T) {
		os.Setenv("GOOGLE_AI_API_KEY", "test-api-key")
		os.Setenv("GEMINI_MODEL", "gemini-pro")
		os.Setenv("GEMINI_TEMPERATURE", "0.5")
		os.Setenv("GEMINI_MAX_TOKENS", "4096")
		os.Setenv("GEMINI_TOP_P", "0.9")
		os.Setenv("ADK_APP_NAME", "test-app")
		os.Setenv("ADK_APP_VERSION", "2.0.0")
		os.Setenv("ADK_APP_DESCRIPTION", "Test Description")

		config := NewConfigFromEnv()

		if config.GetAPIKey() != "test-api-key" {
			t.Errorf("expected APIKey 'test-api-key', got %s", config.GetAPIKey())
		}
		if config.Model != "gemini-pro" {
			t.Errorf("expected Model 'gemini-pro', got %s", config.Model)
		}
		if config.Temperature != 0.5 {
			t.Errorf("expected Temperature 0.5, got %f", config.Temperature)
		}
		if config.MaxTokens != 4096 {
			t.Errorf("expected MaxTokens 4096, got %d", config.MaxTokens)
		}
		if config.TopP != 0.9 {
			t.Errorf("expected TopP 0.9, got %f", config.TopP)
		}
		if config.AppName != "test-app" {
			t.Errorf("expected AppName 'test-app', got %s", config.AppName)
		}
		if config.AppVersion != "2.0.0" {
			t.Errorf("expected AppVersion '2.0.0', got %s", config.AppVersion)
		}
		if config.Description != "Test Description" {
			t.Errorf("expected Description 'Test Description', got %s", config.Description)
		}
	})

	t.Run("with default values", func(t *testing.T) {
		os.Setenv("GOOGLE_AI_API_KEY", "")
		os.Setenv("GEMINI_MODEL", "")
		os.Setenv("GEMINI_TEMPERATURE", "")
		os.Setenv("GEMINI_MAX_TOKENS", "")
		os.Setenv("GEMINI_TOP_P", "")
		os.Setenv("ADK_APP_NAME", "")
		os.Setenv("ADK_APP_VERSION", "")
		os.Setenv("ADK_APP_DESCRIPTION", "")

		config := NewConfigFromEnv()

		if config.GetAPIKey() != "" {
			t.Errorf("expected empty APIKey, got %s", config.GetAPIKey())
		}
		if config.Model != "gemini-2.0-flash" {
			t.Errorf("expected default Model 'gemini-2.0-flash', got %s", config.Model)
		}
		if config.Temperature != 0.7 {
			t.Errorf("expected default Temperature 0.7, got %f", config.Temperature)
		}
		if config.MaxTokens != 8192 {
			t.Errorf("expected default MaxTokens 8192, got %d", config.MaxTokens)
		}
		if config.TopP != 0.95 {
			t.Errorf("expected default TopP 0.95, got %f", config.TopP)
		}
		if config.AppName != "go-adk-chat" {
			t.Errorf("expected default AppName 'go-adk-chat', got %s", config.AppName)
		}
	})

	t.Run("with invalid numeric values", func(t *testing.T) {
		os.Setenv("GOOGLE_AI_API_KEY", "key")
		os.Setenv("GEMINI_TEMPERATURE", "invalid")
		os.Setenv("GEMINI_MAX_TOKENS", "invalid")
		os.Setenv("GEMINI_TOP_P", "invalid")

		config := NewConfigFromEnv()

		// Should fall back to defaults on parse error
		if config.Temperature != 0.7 {
			t.Errorf("expected default Temperature on parse error, got %f", config.Temperature)
		}
		if config.MaxTokens != 8192 {
			t.Errorf("expected default MaxTokens on parse error, got %d", config.MaxTokens)
		}
		if config.TopP != 0.95 {
			t.Errorf("expected default TopP on parse error, got %f", config.TopP)
		}
	})
}

func TestConfig_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		apiKey   string
		model    string
		expected bool
	}{
		{
			name:     "valid config",
			apiKey:   "test-key",
			model:    "gemini-pro",
			expected: true,
		},
		{
			name:     "empty API key",
			apiKey:   "",
			model:    "gemini-pro",
			expected: false,
		},
		{
			name:     "empty model",
			apiKey:   "test-key",
			model:    "",
			expected: false,
		},
		{
			name:     "both empty",
			apiKey:   "",
			model:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				apiKey: tt.apiKey,
				Model:  tt.model,
			}

			if config.IsValid() != tt.expected {
				t.Errorf("expected IsValid() = %v, got %v", tt.expected, config.IsValid())
			}
		})
	}
}
