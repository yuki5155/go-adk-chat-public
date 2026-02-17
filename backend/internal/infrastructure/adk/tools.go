package adk

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/generative-ai-go/genai"
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
	r.entries[def.Name] = toolEntry{
		definition: def,
		handler:    handler,
	}
}

// genaiTools converts registered tools to genai.Tool slice
func (r *toolRegistry) genaiTools() []*genai.Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.entries) == 0 {
		return nil
	}

	var declarations []*genai.FunctionDeclaration
	for _, entry := range r.entries {
		declarations = append(declarations, toFunctionDeclaration(entry.definition))
	}

	return []*genai.Tool{
		{FunctionDeclarations: declarations},
	}
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

// toFunctionDeclaration converts a ToolDefinition to a genai.FunctionDeclaration
func toFunctionDeclaration(def ports.ToolDefinition) *genai.FunctionDeclaration {
	fd := &genai.FunctionDeclaration{
		Name:        def.Name,
		Description: def.Description,
	}

	if len(def.Parameters) > 0 {
		properties := make(map[string]*genai.Schema)
		var required []string

		for _, param := range def.Parameters {
			schema := &genai.Schema{
				Type:        toGenaiType(param.Type),
				Description: param.Description,
			}
			if len(param.Enum) > 0 {
				schema.Enum = param.Enum
			}
			properties[param.Name] = schema
			if param.Required {
				required = append(required, param.Name)
			}
		}

		fd.Parameters = &genai.Schema{
			Type:       genai.TypeObject,
			Properties: properties,
			Required:   required,
		}
	}

	return fd
}

// toGenaiType converts a string type to genai.Type
func toGenaiType(t string) genai.Type {
	switch t {
	case "string":
		return genai.TypeString
	case "integer":
		return genai.TypeInteger
	case "number":
		return genai.TypeNumber
	case "boolean":
		return genai.TypeBoolean
	default:
		return genai.TypeString
	}
}
