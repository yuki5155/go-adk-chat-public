package adk

import (
	"testing"

	"github.com/yuki5155/go-google-auth/internal/application/ports"
)

func TestNewRunnerAdapter(t *testing.T) {
	cfg := &Config{apiKey: "test-key", Model: "gemini-2.0-flash", AppName: "test-app"}
	runner, _ := NewRunner(cfg)
	adapter := NewRunnerAdapter(runner)

	if adapter == nil {
		t.Fatal("expected non-nil adapter")
	}
	if adapter.runner != runner {
		t.Error("expected adapter to reference the runner")
	}
}

func TestRunnerAdapter_Config(t *testing.T) {
	cfg := &Config{apiKey: "test-key", Model: "gemini-2.0-flash", AppName: "my-app"}
	runner, _ := NewRunner(cfg)
	adapter := NewRunnerAdapter(runner)

	result := adapter.Config()
	if result == nil {
		t.Fatal("expected non-nil config")
	}
	if result.AppName != "my-app" {
		t.Errorf("expected AppName 'my-app', got %s", result.AppName)
	}
}

func TestRunnerAdapter_RegisterTool(t *testing.T) {
	cfg := &Config{apiKey: "test-key", Model: "gemini-2.0-flash"}
	runner, _ := NewRunner(cfg)
	adapter := NewRunnerAdapter(runner)

	def := ports.ToolDefinition{Name: "test_tool", Description: "Test"}
	adapter.RegisterTool(def, nil)

	if len(runner.registry.entries) != 1 {
		t.Errorf("expected 1 tool registered, got %d", len(runner.registry.entries))
	}
}
