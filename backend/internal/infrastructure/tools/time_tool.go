package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/yuki5155/go-google-auth/internal/application/ports"
)

// GetCurrentTimeTool returns the tool definition and handler for get_current_time
func GetCurrentTimeTool() (ports.ToolDefinition, ports.ToolHandler) {
	def := ports.ToolDefinition{
		Name:        "get_current_time",
		Description: "Get the current date and time. Optionally specify a timezone in IANA format (e.g. America/New_York, Asia/Tokyo).",
		Parameters: []ports.ToolParameter{
			{
				Name:        "timezone",
				Type:        "string",
				Description: "IANA timezone name (e.g. America/New_York, Asia/Tokyo, Europe/London). Defaults to UTC if not specified.",
				Required:    false,
			},
		},
	}

	handler := func(ctx context.Context, args map[string]any) (*ports.ToolResult, error) {
		tz := "UTC"
		if v, ok := args["timezone"]; ok {
			if s, ok := v.(string); ok && s != "" {
				tz = s
			}
		}

		loc, err := time.LoadLocation(tz)
		if err != nil {
			return nil, fmt.Errorf("invalid timezone %q: %w", tz, err)
		}

		now := time.Now().In(loc)
		return &ports.ToolResult{
			Name: "get_current_time",
			Content: map[string]any{
				"datetime": now.Format(time.RFC3339),
				"timezone": tz,
			},
		}, nil
	}

	return def, handler
}
