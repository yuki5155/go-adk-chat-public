package adk

import (
	"context"
	"testing"

	"github.com/google/generative-ai-go/genai"
	"github.com/yuki5155/go-google-auth/internal/application/ports"
)

func TestToolRegistry_Register(t *testing.T) {
	reg := newToolRegistry()

	def := ports.ToolDefinition{
		Name:        "test_tool",
		Description: "A test tool",
		Parameters: []ports.ToolParameter{
			{Name: "param1", Type: "string", Description: "A parameter", Required: true},
		},
	}
	handler := func(ctx context.Context, args map[string]any) (*ports.ToolResult, error) {
		return &ports.ToolResult{Name: "test_tool", Content: map[string]any{"result": "ok"}}, nil
	}

	reg.register(def, handler)

	if len(reg.entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(reg.entries))
	}
	if _, ok := reg.entries["test_tool"]; !ok {
		t.Error("expected test_tool to be registered")
	}
}

func TestToolRegistry_GenaiTools(t *testing.T) {
	reg := newToolRegistry()

	// Empty registry should return nil
	tools := reg.genaiTools()
	if tools != nil {
		t.Error("expected nil tools for empty registry")
	}

	// Register a tool
	reg.register(ports.ToolDefinition{
		Name:        "my_tool",
		Description: "My tool",
		Parameters: []ports.ToolParameter{
			{Name: "query", Type: "string", Description: "Search query", Required: true},
			{Name: "limit", Type: "integer", Description: "Max results", Required: false},
		},
	}, nil)

	tools = reg.genaiTools()
	if tools == nil {
		t.Fatal("expected non-nil tools")
	}
	if len(tools) != 1 {
		t.Fatalf("expected 1 tool group, got %d", len(tools))
	}
	if len(tools[0].FunctionDeclarations) != 1 {
		t.Fatalf("expected 1 function declaration, got %d", len(tools[0].FunctionDeclarations))
	}

	fd := tools[0].FunctionDeclarations[0]
	if fd.Name != "my_tool" {
		t.Errorf("expected name 'my_tool', got %s", fd.Name)
	}
	if fd.Parameters == nil {
		t.Fatal("expected parameters")
	}
	if len(fd.Parameters.Properties) != 2 {
		t.Errorf("expected 2 properties, got %d", len(fd.Parameters.Properties))
	}
	if len(fd.Parameters.Required) != 1 {
		t.Errorf("expected 1 required param, got %d", len(fd.Parameters.Required))
	}
}

func TestToolRegistry_ExecuteTool(t *testing.T) {
	reg := newToolRegistry()
	ctx := context.Background()

	handler := func(ctx context.Context, args map[string]any) (*ports.ToolResult, error) {
		return &ports.ToolResult{
			Name:    "test_tool",
			Content: map[string]any{"greeting": "hello " + args["name"].(string)},
		}, nil
	}

	reg.register(ports.ToolDefinition{
		Name:        "test_tool",
		Description: "Test tool",
	}, handler)

	// Successful execution
	result, err := reg.executeTool(ctx, "test_tool", map[string]any{"name": "world"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Content["greeting"] != "hello world" {
		t.Errorf("expected 'hello world', got %v", result.Content["greeting"])
	}

	// Unknown tool
	_, err = reg.executeTool(ctx, "unknown_tool", nil)
	if err == nil {
		t.Error("expected error for unknown tool")
	}

	// Tool with nil handler
	reg.register(ports.ToolDefinition{
		Name:        "no_handler",
		Description: "No handler tool",
	}, nil)

	_, err = reg.executeTool(ctx, "no_handler", nil)
	if err == nil {
		t.Error("expected error for tool with nil handler")
	}
}

func TestToGenaiType(t *testing.T) {
	tests := []struct {
		input    string
		expected genai.Type
	}{
		{"string", genai.TypeString},
		{"integer", genai.TypeInteger},
		{"number", genai.TypeNumber},
		{"boolean", genai.TypeBoolean},
		{"unknown", genai.TypeString},
		{"", genai.TypeString},
	}

	for _, tt := range tests {
		result := toGenaiType(tt.input)
		if result != tt.expected {
			t.Errorf("toGenaiType(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestToFunctionDeclaration(t *testing.T) {
	t.Run("with parameters and enum", func(t *testing.T) {
		def := ports.ToolDefinition{
			Name:        "search",
			Description: "Search the web",
			Parameters: []ports.ToolParameter{
				{Name: "query", Type: "string", Description: "Search query", Required: true},
				{Name: "format", Type: "string", Description: "Output format", Required: false, Enum: []string{"json", "text"}},
			},
		}

		fd := toFunctionDeclaration(def)
		if fd.Name != "search" {
			t.Errorf("expected name 'search', got %s", fd.Name)
		}
		if fd.Description != "Search the web" {
			t.Errorf("expected description 'Search the web', got %s", fd.Description)
		}
		if fd.Parameters == nil {
			t.Fatal("expected parameters")
		}
		if fd.Parameters.Type != genai.TypeObject {
			t.Errorf("expected TypeObject, got %v", fd.Parameters.Type)
		}
		formatSchema := fd.Parameters.Properties["format"]
		if formatSchema == nil {
			t.Fatal("expected format property")
		}
		if len(formatSchema.Enum) != 2 {
			t.Errorf("expected 2 enum values, got %d", len(formatSchema.Enum))
		}
	})

	t.Run("without parameters", func(t *testing.T) {
		def := ports.ToolDefinition{
			Name:        "ping",
			Description: "Ping the server",
		}

		fd := toFunctionDeclaration(def)
		if fd.Name != "ping" {
			t.Errorf("expected name 'ping', got %s", fd.Name)
		}
		if fd.Parameters != nil {
			t.Error("expected nil parameters")
		}
	})
}

func TestToolRegistry_MultipleTools(t *testing.T) {
	reg := newToolRegistry()

	reg.register(ports.ToolDefinition{Name: "tool_a", Description: "Tool A"}, nil)
	reg.register(ports.ToolDefinition{Name: "tool_b", Description: "Tool B"}, nil)
	reg.register(ports.ToolDefinition{Name: "tool_c", Description: "Tool C"}, nil)

	if len(reg.entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(reg.entries))
	}

	tools := reg.genaiTools()
	if tools == nil {
		t.Fatal("expected non-nil tools")
	}
	if len(tools[0].FunctionDeclarations) != 3 {
		t.Errorf("expected 3 declarations, got %d", len(tools[0].FunctionDeclarations))
	}
}

func TestToolRegistry_OverwriteExisting(t *testing.T) {
	reg := newToolRegistry()

	reg.register(ports.ToolDefinition{Name: "tool", Description: "Version 1"}, nil)
	reg.register(ports.ToolDefinition{Name: "tool", Description: "Version 2"}, nil)

	if len(reg.entries) != 1 {
		t.Errorf("expected 1 entry after overwrite, got %d", len(reg.entries))
	}
	if reg.entries["tool"].definition.Description != "Version 2" {
		t.Errorf("expected overwritten description 'Version 2', got %s", reg.entries["tool"].definition.Description)
	}
}

func TestExtractFunctionCalls(t *testing.T) {
	t.Run("no candidates", func(t *testing.T) {
		resp := &genai.GenerateContentResponse{}
		calls := extractFunctionCalls(resp)
		if len(calls) != 0 {
			t.Errorf("expected 0 calls, got %d", len(calls))
		}
	})

	t.Run("text only response", func(t *testing.T) {
		resp := &genai.GenerateContentResponse{
			Candidates: []*genai.Candidate{
				{Content: &genai.Content{
					Parts: []genai.Part{genai.Text("Hello")},
				}},
			},
		}
		calls := extractFunctionCalls(resp)
		if len(calls) != 0 {
			t.Errorf("expected 0 calls, got %d", len(calls))
		}
	})

	t.Run("with function calls", func(t *testing.T) {
		resp := &genai.GenerateContentResponse{
			Candidates: []*genai.Candidate{
				{Content: &genai.Content{
					Parts: []genai.Part{
						genai.FunctionCall{Name: "get_time", Args: map[string]any{"tz": "UTC"}},
						genai.FunctionCall{Name: "search", Args: map[string]any{"q": "hello"}},
					},
				}},
			},
		}
		calls := extractFunctionCalls(resp)
		if len(calls) != 2 {
			t.Fatalf("expected 2 calls, got %d", len(calls))
		}
		if calls[0].Name != "get_time" {
			t.Errorf("expected 'get_time', got %s", calls[0].Name)
		}
		if calls[1].Name != "search" {
			t.Errorf("expected 'search', got %s", calls[1].Name)
		}
	})
}

func TestRunner_RegisterTool(t *testing.T) {
	cfg := &Config{apiKey: "test-key", Model: "gemini-2.0-flash"}
	runner, _ := NewRunner(cfg)

	def := ports.ToolDefinition{Name: "my_tool", Description: "My tool"}
	handler := func(ctx context.Context, args map[string]any) (*ports.ToolResult, error) {
		return &ports.ToolResult{Name: "my_tool", Content: map[string]any{}}, nil
	}

	runner.RegisterTool(def, handler)

	if len(runner.registry.entries) != 1 {
		t.Errorf("expected 1 registered tool, got %d", len(runner.registry.entries))
	}
}

func TestRunner_HasRegistry(t *testing.T) {
	cfg := &Config{apiKey: "test-key", Model: "gemini-2.0-flash"}
	runner, err := NewRunner(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if runner.registry == nil {
		t.Error("expected non-nil registry on new runner")
	}
}
