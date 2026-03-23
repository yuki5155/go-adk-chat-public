package anthropic

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	anthropicsdk "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/yuki5155/go-google-auth/internal/application/ports"
)

const systemPrompt = "You are a helpful AI assistant. Be concise, accurate, and helpful in your responses."

// Runner wraps the Anthropic client for chat operations
type Runner struct {
	config   *Config
	client   *anthropicsdk.Client
	registry *toolRegistry
}

// NewRunner creates a new Anthropic Runner
func NewRunner(config *Config) (*Runner, error) {
	if !config.IsValid() {
		return nil, fmt.Errorf("invalid Anthropic configuration: API key or model not set")
	}
	client := anthropicsdk.NewClient(option.WithAPIKey(config.GetAPIKey()))
	return &Runner{
		config:   config,
		client:   &client,
		registry: newToolRegistry(),
	}, nil
}

// RegisterTool registers a tool definition with an optional handler
func (r *Runner) RegisterTool(def ports.ToolDefinition, handler ports.ToolHandler) {
	r.registry.register(def, handler)
}

// Config returns the runner config
func (r *Runner) Config() *Config {
	return r.config
}

// SendMessage sends a message and returns the complete response
func (r *Runner) SendMessage(ctx context.Context, history []ports.ChatMessage, userMessage string, modelName string) (string, error) {
	model := r.resolveModel(modelName)
	messages := buildMessages(history, userMessage)

	for {
		params := anthropicsdk.MessageNewParams{
			Model:     anthropicsdk.Model(model),
			Messages:  messages,
			MaxTokens: r.config.MaxTokens,
			System: []anthropicsdk.TextBlockParam{
				{Text: systemPrompt},
			},
		}
		if tools := r.registry.anthropicTools(); len(tools) > 0 {
			params.Tools = tools
		}

		resp, err := r.client.Messages.New(ctx, params)
		if err != nil {
			return "", fmt.Errorf("anthropic: failed to send message: %w", err)
		}

		// No tool use → extract text and return
		if resp.StopReason != anthropicsdk.StopReasonToolUse {
			return extractText(resp.Content), nil
		}

		// Append assistant message and process tool calls
		messages = append(messages, anthropicsdk.NewAssistantMessage(responseToBlocks(resp.Content)...))

		toolResults := processToolCalls(ctx, resp.Content, r.registry)
		messages = append(messages, anthropicsdk.NewUserMessage(toolResults...))
	}
}

// StreamMessage sends a message and streams text chunks via callback
func (r *Runner) StreamMessage(ctx context.Context, history []ports.ChatMessage, userMessage string, modelName string, callback func(chunk string) error) error {
	model := r.resolveModel(modelName)
	messages := buildMessages(history, userMessage)

	params := anthropicsdk.MessageNewParams{
		Model:     anthropicsdk.Model(model),
		Messages:  messages,
		MaxTokens: r.config.MaxTokens,
		System: []anthropicsdk.TextBlockParam{
			{Text: systemPrompt},
		},
	}

	stream := r.client.Messages.NewStreaming(ctx, params)
	acc := anthropicsdk.Message{}
	chunkCount := 0

	for stream.Next() {
		event := stream.Current()
		_ = acc.Accumulate(event)
		if event.Type == "content_block_delta" && event.Delta.Type == "text_delta" && event.Delta.Text != "" {
			chunkCount++
			if err := callback(event.Delta.Text); err != nil {
				return err
			}
		}
	}
	if err := stream.Err(); err != nil {
		return fmt.Errorf("anthropic: stream error: %w", err)
	}

	log.Printf("[Anthropic] Stream completed with %d chunks", chunkCount)
	return nil
}

// StreamMessageWithTools sends a message and streams the response, executing tools as needed
func (r *Runner) StreamMessageWithTools(ctx context.Context, history []ports.ChatMessage, userMessage string, modelName string, callback ports.StreamEventCallback) error {
	model := r.resolveModel(modelName)
	messages := buildMessages(history, userMessage)

	for {
		params := anthropicsdk.MessageNewParams{
			Model:     anthropicsdk.Model(model),
			Messages:  messages,
			MaxTokens: r.config.MaxTokens,
			System: []anthropicsdk.TextBlockParam{
				{Text: systemPrompt},
			},
		}
		if tools := r.registry.anthropicTools(); len(tools) > 0 {
			params.Tools = tools
		}

		stream := r.client.Messages.NewStreaming(ctx, params)
		acc := anthropicsdk.Message{}

		for stream.Next() {
			event := stream.Current()
			_ = acc.Accumulate(event)
			if event.Type == "content_block_delta" && event.Delta.Type == "text_delta" && event.Delta.Text != "" {
				if err := callback(ports.StreamEvent{
					Type:    ports.StreamEventChunk,
					Content: event.Delta.Text,
				}); err != nil {
					return err
				}
			}
		}
		if err := stream.Err(); err != nil {
			return fmt.Errorf("anthropic: stream error: %w", err)
		}

		// No tool use → done
		if acc.StopReason != anthropicsdk.StopReasonToolUse {
			break
		}

		// Append assistant turn and execute tools
		messages = append(messages, anthropicsdk.NewAssistantMessage(responseToBlocks(acc.Content)...))

		var toolResultBlocks []anthropicsdk.ContentBlockParamUnion
		for _, block := range acc.Content {
			toolUse, ok := block.AsAny().(anthropicsdk.ToolUseBlock)
			if !ok {
				continue
			}

			var args map[string]any
			_ = json.Unmarshal(toolUse.Input, &args)

			toolCall := &ports.ToolCall{Name: toolUse.Name, Args: args}
			if err := callback(ports.StreamEvent{Type: ports.StreamEventToolStart, ToolCall: toolCall}); err != nil {
				return err
			}

			log.Printf("[Anthropic] Executing tool: %s", toolUse.Name)
			result, execErr := r.registry.executeTool(ctx, toolUse.Name, args)

			var resultText string
			if execErr != nil {
				resultText = fmt.Sprintf(`{"error": %q}`, execErr.Error())
			} else {
				b, _ := json.Marshal(result.Content)
				resultText = string(b)
			}

			toolResultBlocks = append(toolResultBlocks, anthropicsdk.ContentBlockParamUnion{
				OfToolResult: &anthropicsdk.ToolResultBlockParam{
					ToolUseID: toolUse.ID,
					Content: []anthropicsdk.ToolResultBlockParamContentUnion{
						{OfText: &anthropicsdk.TextBlockParam{Text: resultText}},
					},
				},
			})

			if err := callback(ports.StreamEvent{Type: ports.StreamEventToolEnd, ToolCall: toolCall}); err != nil {
				return err
			}
		}

		messages = append(messages, anthropicsdk.NewUserMessage(toolResultBlocks...))
	}

	return nil
}

// resolveModel returns modelName if set, otherwise falls back to config default
func (r *Runner) resolveModel(modelName string) string {
	if modelName != "" {
		return modelName
	}
	return r.config.Model
}

// buildMessages converts history and the new user message into Anthropic message params
func buildMessages(history []ports.ChatMessage, userMessage string) []anthropicsdk.MessageParam {
	var messages []anthropicsdk.MessageParam
	for _, msg := range history {
		text := anthropicsdk.ContentBlockParamUnion{
			OfText: &anthropicsdk.TextBlockParam{Text: msg.Content},
		}
		switch msg.Role {
		case "assistant", "model":
			messages = append(messages, anthropicsdk.NewAssistantMessage(text))
		default:
			messages = append(messages, anthropicsdk.NewUserMessage(text))
		}
	}
	messages = append(messages, anthropicsdk.NewUserMessage(
		anthropicsdk.ContentBlockParamUnion{
			OfText: &anthropicsdk.TextBlockParam{Text: userMessage},
		},
	))
	return messages
}

// extractText pulls all text out of a response content block list
func extractText(content []anthropicsdk.ContentBlockUnion) string {
	var text string
	for _, block := range content {
		if tb, ok := block.AsAny().(anthropicsdk.TextBlock); ok {
			text += tb.Text
		}
	}
	return text
}

// responseToBlocks converts response content to ContentBlockParamUnion for the next turn
func responseToBlocks(content []anthropicsdk.ContentBlockUnion) []anthropicsdk.ContentBlockParamUnion {
	var blocks []anthropicsdk.ContentBlockParamUnion
	for _, block := range content {
		switch b := block.AsAny().(type) {
		case anthropicsdk.TextBlock:
			blocks = append(blocks, anthropicsdk.ContentBlockParamUnion{
				OfText: &anthropicsdk.TextBlockParam{Text: b.Text},
			})
		case anthropicsdk.ToolUseBlock:
			var input map[string]any
			_ = json.Unmarshal(b.Input, &input)
			blocks = append(blocks, anthropicsdk.ContentBlockParamUnion{
				OfToolUse: &anthropicsdk.ToolUseBlockParam{
					ID:    b.ID,
					Name:  b.Name,
					Input: input,
				},
			})
		}
	}
	return blocks
}

// processToolCalls executes all tool-use blocks and returns tool result blocks
func processToolCalls(ctx context.Context, content []anthropicsdk.ContentBlockUnion, registry *toolRegistry) []anthropicsdk.ContentBlockParamUnion {
	var results []anthropicsdk.ContentBlockParamUnion
	for _, block := range content {
		toolUse, ok := block.AsAny().(anthropicsdk.ToolUseBlock)
		if !ok {
			continue
		}

		var args map[string]any
		_ = json.Unmarshal(toolUse.Input, &args)

		log.Printf("[Anthropic] Executing tool: %s", toolUse.Name)
		result, execErr := registry.executeTool(ctx, toolUse.Name, args)

		var resultText string
		if execErr != nil {
			resultText = fmt.Sprintf(`{"error": %q}`, execErr.Error())
		} else {
			b, _ := json.Marshal(result.Content)
			resultText = string(b)
		}

		results = append(results, anthropicsdk.ContentBlockParamUnion{
			OfToolResult: &anthropicsdk.ToolResultBlockParam{
				ToolUseID: toolUse.ID,
				Content: []anthropicsdk.ToolResultBlockParamContentUnion{
					{OfText: &anthropicsdk.TextBlockParam{Text: resultText}},
				},
			},
		})
	}
	return results
}
