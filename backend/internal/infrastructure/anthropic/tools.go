package anthropic

import (
	"context"
	"fmt"
	"sync"

	anthropicsdk "github.com/anthropics/anthropic-sdk-go"
	"github.com/yuki5155/go-google-auth/internal/application/ports"
)

// toolEntry holds a tool definition and its handler
type toolEntry struct {
	definition ports.ToolDefinition
	handler    ports.ToolHandler
}

// toolRegistry manages registered tools
type toolRegistry struct {
	mu      sync.RWMutex
	entries map[string]toolEntry
}

// newToolRegistry creates a new tool registry
func newToolRegistry() *toolRegistry {
	return &toolRegistry{
		entries: make(map[string]toolEntry),
	}
}

// register adds a tool to the registry
func (r *toolRegistry) register(def ports.ToolDefinition, handler ports.ToolHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[def.Name] = toolEntry{definition: def, handler: handler}
}

// anthropicTools converts registered tools to Anthropic tool params
func (r *toolRegistry) anthropicTools() []anthropicsdk.ToolUnionParam {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.entries) == 0 {
		return nil
	}

	var tools []anthropicsdk.ToolUnionParam
	for _, entry := range r.entries {
		tools = append(tools, toAnthropicTool(entry.definition))
	}
	return tools
}

// executeTool dispatches a function call to the correct handler
func (r *toolRegistry) executeTool(ctx context.Context, name string, args map[string]any) (*ports.ToolResult, error) {
	r.mu.RLock()
	entry, ok := r.entries[name]
	r.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
	if entry.handler == nil {
		return nil, fmt.Errorf("tool %s has no handler", name)
	}
	return entry.handler(ctx, args)
}

// toAnthropicTool converts a ToolDefinition to an Anthropic tool param
func toAnthropicTool(def ports.ToolDefinition) anthropicsdk.ToolUnionParam {
	properties := make(map[string]any)
	var required []string

	for _, param := range def.Parameters {
		prop := map[string]any{
			"type":        param.Type,
			"description": param.Description,
		}
		if len(param.Enum) > 0 {
			prop["enum"] = param.Enum
		}
		properties[param.Name] = prop
		if param.Required {
			required = append(required, param.Name)
		}
	}

	toolParam := anthropicsdk.ToolParam{
		Name: def.Name,
		InputSchema: anthropicsdk.ToolInputSchemaParam{
			Properties: properties,
		},
	}
	if def.Description != "" {
		toolParam.Description = anthropicsdk.String(def.Description)
	}
	if len(required) > 0 {
		toolParam.InputSchema.Required = required
	}

	return anthropicsdk.ToolUnionParam{OfTool: &toolParam}
}
