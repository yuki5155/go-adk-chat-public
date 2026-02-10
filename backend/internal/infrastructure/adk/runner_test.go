package adk

import (
	"testing"
)

func TestNewRunner(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		config := &Config{
			APIKey: "test-api-key",
			Model:  "gemini-pro",
		}

		runner, err := NewRunner(config)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if runner == nil {
			t.Fatal("expected runner to be created")
		}
		if runner.config != config {
			t.Error("expected config to be set")
		}
	})

	t.Run("invalid config - empty API key", func(t *testing.T) {
		config := &Config{
			APIKey: "",
			Model:  "gemini-pro",
		}

		_, err := NewRunner(config)

		if err == nil {
			t.Error("expected error for invalid config")
		}
	})

	t.Run("invalid config - empty model", func(t *testing.T) {
		config := &Config{
			APIKey: "test-api-key",
			Model:  "",
		}

		_, err := NewRunner(config)

		if err == nil {
			t.Error("expected error for invalid config")
		}
	})
}

func TestRunner_Config(t *testing.T) {
	config := &Config{
		APIKey:  "test-api-key",
		Model:   "gemini-pro",
		AppName: "test-app",
	}

	runner, _ := NewRunner(config)

	if runner.Config() != config {
		t.Error("expected Config() to return the same config")
	}
	if runner.Config().AppName != "test-app" {
		t.Errorf("expected AppName 'test-app', got %s", runner.Config().AppName)
	}
}

func TestRunner_Close(t *testing.T) {
	t.Run("close uninitialized runner", func(t *testing.T) {
		config := &Config{
			APIKey: "test-api-key",
			Model:  "gemini-pro",
		}

		runner, _ := NewRunner(config)

		// Should not panic or error when client is nil
		err := runner.Close()

		if err != nil {
			t.Errorf("unexpected error closing uninitialized runner: %v", err)
		}
	})
}

func TestChatMessage(t *testing.T) {
	msg := ChatMessage{
		Role:    "user",
		Content: "Hello, world!",
	}

	if msg.Role != "user" {
		t.Errorf("expected Role 'user', got %s", msg.Role)
	}
	if msg.Content != "Hello, world!" {
		t.Errorf("expected Content 'Hello, world!', got %s", msg.Content)
	}
}
