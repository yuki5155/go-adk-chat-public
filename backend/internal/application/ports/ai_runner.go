// Package ports defines interfaces for external dependencies
package ports

import "context"

// ChatMessage represents a message in the chat history
type ChatMessage struct {
	Role    string
	Content string
}

// AIRunnerConfig holds the configuration for the AI runner
type AIRunnerConfig struct {
	AppName string
}

// AIRunner defines the interface for AI chat operations
type AIRunner interface {
	// SendMessage sends a message and returns the complete response
	// model parameter allows specifying which AI model to use (empty string uses default)
	SendMessage(ctx context.Context, history []ChatMessage, userMessage string, model string) (string, error)

	// StreamMessage sends a message and streams the response via callback
	// model parameter allows specifying which AI model to use (empty string uses default)
	StreamMessage(ctx context.Context, history []ChatMessage, userMessage string, model string, callback func(chunk string) error) error

	// RegisterTool registers a tool definition with an optional handler
	RegisterTool(def ToolDefinition, handler ToolHandler)

	// StreamMessageWithTools sends a message and streams the response, executing tools as needed
	StreamMessageWithTools(ctx context.Context, history []ChatMessage, userMessage string, model string, callback StreamEventCallback) error

	// Config returns the runner configuration
	Config() *AIRunnerConfig
}
