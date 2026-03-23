package openai

import (
	"context"

	"github.com/yuki5155/go-google-auth/internal/application/ports"
)

// RunnerAdapter adapts Runner to implement ports.AIRunner
type RunnerAdapter struct {
	runner *Runner
}

// NewRunnerAdapter creates a new adapter for the Runner
func NewRunnerAdapter(runner *Runner) *RunnerAdapter {
	return &RunnerAdapter{runner: runner}
}

// SendMessage implements ports.AIRunner
func (a *RunnerAdapter) SendMessage(ctx context.Context, history []ports.ChatMessage, userMessage string, model string) (string, error) {
	return a.runner.SendMessage(ctx, history, userMessage, model)
}

// StreamMessage implements ports.AIRunner
func (a *RunnerAdapter) StreamMessage(ctx context.Context, history []ports.ChatMessage, userMessage string, model string, callback func(chunk string) error) error {
	return a.runner.StreamMessage(ctx, history, userMessage, model, callback)
}

// RegisterTool implements ports.AIRunner
func (a *RunnerAdapter) RegisterTool(def ports.ToolDefinition, handler ports.ToolHandler) {
	a.runner.RegisterTool(def, handler)
}

// StreamMessageWithTools implements ports.AIRunner
func (a *RunnerAdapter) StreamMessageWithTools(ctx context.Context, history []ports.ChatMessage, userMessage string, model string, callback ports.StreamEventCallback) error {
	return a.runner.StreamMessageWithTools(ctx, history, userMessage, model, callback)
}

// Config implements ports.AIRunner
func (a *RunnerAdapter) Config() *ports.AIRunnerConfig {
	return &ports.AIRunnerConfig{AppName: a.runner.Config().AppName}
}
