package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	openaisdk "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/yuki5155/go-google-auth/internal/application/ports"
)

const systemPrompt = "You are a helpful AI assistant. Be concise, accurate, and helpful in your responses."

// Runner wraps the OpenAI client for chat operations
type Runner struct {
	config   *Config
	client   *openaisdk.Client
	registry *toolRegistry
}

// NewRunner creates a new OpenAI Runner
func NewRunner(config *Config) (*Runner, error) {
	if !config.IsValid() {
		return nil, fmt.Errorf("invalid OpenAI configuration: API key or model not set")
	}
	client := openaisdk.NewClient(option.WithAPIKey(config.GetAPIKey()))
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
	messages := r.buildMessages(history, userMessage)

	for {
		params := openaisdk.ChatCompletionNewParams{
			Model:    openaisdk.ChatModel(model),
			Messages: messages,
		}
		if tools := r.registry.openAITools(); len(tools) > 0 {
			params.Tools = tools
		}

		resp, err := r.client.Chat.Completions.New(ctx, params)
		if err != nil {
			return "", fmt.Errorf("openai: failed to send message: %w", err)
		}
		if len(resp.Choices) == 0 {
			return "", fmt.Errorf("openai: empty response")
		}

		choice := resp.Choices[0]

		// No tool calls → return text
		if len(choice.Message.ToolCalls) == 0 {
			return choice.Message.Content, nil
		}

		// Append assistant message with tool calls
		messages = append(messages, choice.Message.ToParam())

		// Execute tool calls and append results
		for _, tc := range choice.Message.ToolCalls {
			var args map[string]any
			_ = json.Unmarshal([]byte(tc.Function.Arguments), &args)

			log.Printf("[OpenAI] Executing tool: %s", tc.Function.Name)
			result, execErr := r.registry.executeTool(ctx, tc.Function.Name, args)

			var content string
			if execErr != nil {
				content = fmt.Sprintf(`{"error": %q}`, execErr.Error())
			} else {
				b, _ := json.Marshal(result.Content)
				content = string(b)
			}
			messages = append(messages, openaisdk.ToolMessage(content, tc.ID))
		}
	}
}

// StreamMessage sends a message and streams text chunks via callback
func (r *Runner) StreamMessage(ctx context.Context, history []ports.ChatMessage, userMessage string, modelName string, callback func(chunk string) error) error {
	model := r.resolveModel(modelName)
	messages := r.buildMessages(history, userMessage)

	params := openaisdk.ChatCompletionNewParams{
		Model:    openaisdk.ChatModel(model),
		Messages: messages,
	}

	stream := r.client.Chat.Completions.NewStreaming(ctx, params)
	acc := openaisdk.ChatCompletionAccumulator{}
	chunkCount := 0

	for stream.Next() {
		chunk := stream.Current()
		acc.AddChunk(chunk)
		if len(chunk.Choices) == 0 {
			continue
		}
		delta := chunk.Choices[0].Delta.Content
		if delta == "" {
			continue
		}
		chunkCount++
		if err := callback(delta); err != nil {
			return err
		}
	}
	if err := stream.Err(); err != nil {
		return fmt.Errorf("openai: stream error: %w", err)
	}

	log.Printf("[OpenAI] Stream completed with %d chunks", chunkCount)
	return nil
}

// StreamMessageWithTools sends a message and streams the response, executing tools as needed
func (r *Runner) StreamMessageWithTools(ctx context.Context, history []ports.ChatMessage, userMessage string, modelName string, callback ports.StreamEventCallback) error {
	model := r.resolveModel(modelName)
	messages := r.buildMessages(history, userMessage)

	for {
		params := openaisdk.ChatCompletionNewParams{
			Model:    openaisdk.ChatModel(model),
			Messages: messages,
		}
		if tools := r.registry.openAITools(); len(tools) > 0 {
			params.Tools = tools
		}

		stream := r.client.Chat.Completions.NewStreaming(ctx, params)
		acc := openaisdk.ChatCompletionAccumulator{}

		for stream.Next() {
			chunk := stream.Current()
			acc.AddChunk(chunk)
			if len(chunk.Choices) == 0 {
				continue
			}
			delta := chunk.Choices[0].Delta.Content
			if delta == "" {
				continue
			}
			if err := callback(ports.StreamEvent{
				Type:    ports.StreamEventChunk,
				Content: delta,
			}); err != nil {
				return err
			}
		}
		if err := stream.Err(); err != nil {
			return fmt.Errorf("openai: stream error: %w", err)
		}

		// No tool calls → done
		if len(acc.Choices) == 0 || len(acc.Choices[0].Message.ToolCalls) == 0 {
			break
		}

		// Append assistant message with tool calls
		messages = append(messages, acc.Choices[0].Message.ToParam())

		// Execute tool calls
		for _, tc := range acc.Choices[0].Message.ToolCalls {
			var args map[string]any
			_ = json.Unmarshal([]byte(tc.Function.Arguments), &args)

			toolCall := &ports.ToolCall{Name: tc.Function.Name, Args: args}
			if err := callback(ports.StreamEvent{Type: ports.StreamEventToolStart, ToolCall: toolCall}); err != nil {
				return err
			}

			log.Printf("[OpenAI] Executing tool: %s", tc.Function.Name)
			result, execErr := r.registry.executeTool(ctx, tc.Function.Name, args)

			var content string
			if execErr != nil {
				content = fmt.Sprintf(`{"error": %q}`, execErr.Error())
			} else {
				b, _ := json.Marshal(result.Content)
				content = string(b)
			}
			messages = append(messages, openaisdk.ToolMessage(content, tc.ID))

			if err := callback(ports.StreamEvent{Type: ports.StreamEventToolEnd, ToolCall: toolCall}); err != nil {
				return err
			}
		}
	}

	return nil
}

// resolveModel returns the model name to use, falling back to config default
func (r *Runner) resolveModel(modelName string) string {
	if modelName != "" {
		return modelName
	}
	return r.config.Model
}

// buildMessages converts chat history and the new user message into OpenAI message params
func (r *Runner) buildMessages(history []ports.ChatMessage, userMessage string) []openaisdk.ChatCompletionMessageParamUnion {
	messages := []openaisdk.ChatCompletionMessageParamUnion{
		openaisdk.SystemMessage(systemPrompt),
	}
	for _, msg := range history {
		switch msg.Role {
		case "assistant", "model":
			messages = append(messages, openaisdk.AssistantMessage(msg.Content))
		default:
			messages = append(messages, openaisdk.UserMessage(msg.Content))
		}
	}
	messages = append(messages, openaisdk.UserMessage(userMessage))
	return messages
}
