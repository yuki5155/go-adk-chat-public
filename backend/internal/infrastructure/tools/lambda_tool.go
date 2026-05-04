package tools

import (
	"context"

	"github.com/yuki5155/go-google-auth/internal/application/ports"
	"github.com/yuki5155/go-google-auth/internal/utils"
)

// LambdaTool builds a ToolDefinition + ToolHandler from a loaded LambdaToolEntry.
func LambdaTool(entry utils.LambdaToolEntry, apiURL, apiKey string) (ports.ToolDefinition, ports.ToolHandler) {
	params := make([]ports.ToolParameter, len(entry.Schema.Parameters))
	for i, p := range entry.Schema.Parameters {
		params[i] = ports.ToolParameter{
			Name:        p.Name,
			Type:        p.Type,
			Description: p.Description,
			Required:    p.Required,
		}
	}

	def := ports.ToolDefinition{
		Name:        entry.Schema.Name,
		Description: entry.Schema.Description,
		Parameters:  params,
	}

	client := utils.NewLambdaToolClient(apiURL, apiKey)

	handler := func(ctx context.Context, args map[string]any) (*ports.ToolResult, error) {
		result, err := client.Call(ctx, args)
		if err != nil {
			return nil, err
		}
		return &ports.ToolResult{Name: entry.Schema.Name, Content: result}, nil
	}

	return def, handler
}
