// Package ports defines interfaces for external dependencies
package ports

import "context"

// ToolDefinition describes a tool that the AI model can invoke
type ToolDefinition struct {
	Name        string
	Description string
	Parameters  []ToolParameter
}

// ToolParameter describes a single parameter for a tool
type ToolParameter struct {
	Name        string
	Type        string // "string", "integer", "number", "boolean"
	Description string
	Required    bool
	Enum        []string
}

// ToolCall represents the model requesting execution of a tool
type ToolCall struct {
	Name string
	Args map[string]any
}

// ToolResult represents the result of executing a tool
type ToolResult struct {
	Name    string
	Content map[string]any
}

// ToolHandler is a function that executes a tool
type ToolHandler func(ctx context.Context, args map[string]any) (*ToolResult, error)

// StreamEventType identifies the kind of stream event
type StreamEventType int

const (
	// StreamEventChunk is a text chunk from the model
	StreamEventChunk StreamEventType = iota
	// StreamEventToolStart indicates a tool is about to be executed
	StreamEventToolStart
	// StreamEventToolEnd indicates a tool has finished executing
	StreamEventToolEnd
)

// StreamEvent is a typed event emitted during streaming
type StreamEvent struct {
	Type     StreamEventType
	Content  string
	ToolCall *ToolCall
}

// StreamEventCallback is a callback for stream events
type StreamEventCallback func(event StreamEvent) error
