package adk

import (
	"context"

	"github.com/yuki5155/go-google-auth/internal/application/ports"
)

// RunnerAdapter adapts the Runner to implement ports.AIRunner interface
type RunnerAdapter struct {
	runner *Runner
}

// NewRunnerAdapter creates a new adapter for the Runner
func NewRunnerAdapter(runner *Runner) *RunnerAdapter {
	return &RunnerAdapter{runner: runner}
}

// SendMessage implements ports.AIRunner
func (a *RunnerAdapter) SendMessage(ctx context.Context, history []ports.ChatMessage, userMessage string, model string) (string, error) {
	// Convert ports.ChatMessage to adk.ChatMessage
	adkHistory := make([]ChatMessage, len(history))
	for i, msg := range history {
		adkHistory[i] = ChatMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}
	return a.runner.SendMessage(ctx, adkHistory, userMessage, model)
}

// StreamMessage implements ports.AIRunner
func (a *RunnerAdapter) StreamMessage(ctx context.Context, history []ports.ChatMessage, userMessage string, model string, callback func(chunk string) error) error {
	// Convert ports.ChatMessage to adk.ChatMessage
	adkHistory := make([]ChatMessage, len(history))
	for i, msg := range history {
		adkHistory[i] = ChatMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}
	return a.runner.StreamMessage(ctx, adkHistory, userMessage, model, callback)
}

// Config implements ports.AIRunner
func (a *RunnerAdapter) Config() *ports.AIRunnerConfig {
	cfg := a.runner.Config()
	return &ports.AIRunnerConfig{
		AppName: cfg.AppName,
	}
}
