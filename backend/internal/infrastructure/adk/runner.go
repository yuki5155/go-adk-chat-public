package adk

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/generative-ai-go/genai"
	"github.com/yuki5155/go-google-auth/internal/application/ports"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// Runner wraps the Gemini generative AI client for chat operations
type Runner struct {
	config   *Config
	client   *genai.Client
	model    *genai.GenerativeModel
	registry *toolRegistry
}

// NewRunner creates a new ADK Runner
func NewRunner(config *Config) (*Runner, error) {
	if !config.IsValid() {
		return nil, fmt.Errorf("invalid ADK configuration: API key or model not set")
	}

	return &Runner{
		config:   config,
		registry: newToolRegistry(),
	}, nil
}

// RegisterTool registers a tool definition with an optional handler
func (r *Runner) RegisterTool(def ports.ToolDefinition, handler ports.ToolHandler) {
	r.registry.register(def, handler)
}

// Initialize initializes the Gemini client
func (r *Runner) Initialize(ctx context.Context) error {
	client, err := genai.NewClient(ctx, option.WithAPIKey(r.config.GetAPIKey()))
	if err != nil {
		return fmt.Errorf("failed to create Gemini client: %w", err)
	}
	r.client = client

	// Configure the model
	r.model = client.GenerativeModel(r.config.Model)
	r.model.SetTemperature(r.config.Temperature)
	r.model.SetMaxOutputTokens(r.config.MaxTokens)
	r.model.SetTopP(r.config.TopP)

	// Set system instruction
	r.model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{
			genai.Text("You are a helpful AI assistant. Be concise, accurate, and helpful in your responses."),
		},
	}

	return nil
}

// Close closes the Gemini client
func (r *Runner) Close() error {
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}

// ChatMessage represents a message in the chat
type ChatMessage struct {
	Role    string
	Content string
}

// getModel returns the model to use, either the specified one or the default
func (r *Runner) getModel(modelName string) *genai.GenerativeModel {
	if modelName == "" {
		// Attach tools to the default model
		if tools := r.registry.genaiTools(); tools != nil {
			r.model.Tools = tools
		}
		return r.model
	}

	// Create a new model with the specified name
	model := r.client.GenerativeModel(modelName)
	model.SetTemperature(r.config.Temperature)
	model.SetMaxOutputTokens(int32(r.config.MaxTokens))
	model.SetTopP(r.config.TopP)

	// Set system instruction
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{
			genai.Text("You are a helpful AI assistant. Be concise, accurate, and helpful in your responses."),
		},
	}

	// Attach tools
	if tools := r.registry.genaiTools(); tools != nil {
		model.Tools = tools
	}

	return model
}

// SendMessage sends a message and returns the response
// model parameter allows specifying which AI model to use (empty string uses default)
func (r *Runner) SendMessage(ctx context.Context, history []ChatMessage, userMessage string, model string) (string, error) {
	if r.client == nil {
		if err := r.Initialize(ctx); err != nil {
			return "", err
		}
	}

	// Get the model to use
	genModel := r.getModel(model)

	// Start a chat session
	cs := genModel.StartChat()

	// Build history
	for _, msg := range history {
		role := msg.Role
		if role == "assistant" {
			role = "model"
		}
		cs.History = append(cs.History, &genai.Content{
			Parts: []genai.Part{genai.Text(msg.Content)},
			Role:  role,
		})
	}

	// Send the user message and handle function call loop
	resp, err := cs.SendMessage(ctx, genai.Text(userMessage))
	if err != nil {
		return "", fmt.Errorf("failed to send message: %w", err)
	}

	// Function call loop: keep processing until we get a text response
	for {
		if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
			return "", fmt.Errorf("no response received from model")
		}

		// Check for function calls
		funcCalls := extractFunctionCalls(resp)
		if len(funcCalls) == 0 {
			break // No function calls, extract text
		}

		// Execute function calls and collect responses
		var funcResponses []genai.Part
		for _, fc := range funcCalls {
			log.Printf("[ADK] Executing tool: %s", fc.Name)
			result, execErr := r.registry.executeTool(ctx, fc.Name, fc.Args)
			if execErr != nil {
				funcResponses = append(funcResponses, genai.FunctionResponse{
					Name:     fc.Name,
					Response: map[string]any{"error": execErr.Error()},
				})
			} else {
				funcResponses = append(funcResponses, genai.FunctionResponse{
					Name:     fc.Name,
					Response: result.Content,
				})
			}
		}

		// Send function responses back to model
		resp, err = cs.SendMessage(ctx, funcResponses...)
		if err != nil {
			return "", fmt.Errorf("failed to send function response: %w", err)
		}
	}

	// Extract text response
	var responseText string
	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			responseText += string(text)
		}
	}

	return responseText, nil
}

// StreamMessage sends a message and streams the response
// model parameter allows specifying which AI model to use (empty string uses default)
func (r *Runner) StreamMessage(ctx context.Context, history []ChatMessage, userMessage string, model string, callback func(chunk string) error) error {
	if r.client == nil {
		if err := r.Initialize(ctx); err != nil {
			return err
		}
	}

	// Get the model to use
	genModel := r.getModel(model)

	// Start a chat session
	cs := genModel.StartChat()

	// Build history
	for _, msg := range history {
		role := msg.Role
		if role == "assistant" {
			role = "model"
		}
		cs.History = append(cs.History, &genai.Content{
			Parts: []genai.Part{genai.Text(msg.Content)},
			Role:  role,
		})
	}

	// Send the user message with streaming
	iter := cs.SendMessageStream(ctx, genai.Text(userMessage))

	chunkCount := 0
	for {
		resp, err := iter.Next()
		if err != nil {
			// Check if we've reached the end of the stream
			if errors.Is(err, iterator.Done) {
				break
			}
			return fmt.Errorf("stream error: %w", err)
		}

		// Extract text from response
		for _, candidate := range resp.Candidates {
			for _, part := range candidate.Content.Parts {
				if text, ok := part.(genai.Text); ok {
					chunkCount++
					if err := callback(string(text)); err != nil {
						return err
					}
				}
			}
		}
	}

	log.Printf("[ADK] Stream completed with %d chunks", chunkCount)
	return nil
}

// StreamMessageWithTools sends a message and streams the response, executing tools as needed
func (r *Runner) StreamMessageWithTools(ctx context.Context, history []ChatMessage, userMessage string, model string, callback ports.StreamEventCallback) error {
	if r.client == nil {
		if err := r.Initialize(ctx); err != nil {
			return err
		}
	}

	genModel := r.getModel(model)
	cs := genModel.StartChat()

	// Build history
	for _, msg := range history {
		role := msg.Role
		if role == "assistant" {
			role = "model"
		}
		cs.History = append(cs.History, &genai.Content{
			Parts: []genai.Part{genai.Text(msg.Content)},
			Role:  role,
		})
	}

	// Initial message parts
	var sendParts []genai.Part
	sendParts = append(sendParts, genai.Text(userMessage))

	for {
		// Stream the response
		iter := cs.SendMessageStream(ctx, sendParts...)

		var accumulatedParts []genai.Part
		for {
			resp, err := iter.Next()
			if err != nil {
				if errors.Is(err, iterator.Done) {
					break
				}
				return fmt.Errorf("stream error: %w", err)
			}

			for _, candidate := range resp.Candidates {
				for _, part := range candidate.Content.Parts {
					accumulatedParts = append(accumulatedParts, part)
					if text, ok := part.(genai.Text); ok {
						if err := callback(ports.StreamEvent{
							Type:    ports.StreamEventChunk,
							Content: string(text),
						}); err != nil {
							return err
						}
					}
				}
			}
		}

		// Check for function calls in accumulated parts
		var funcCalls []genai.FunctionCall
		for _, part := range accumulatedParts {
			if fc, ok := part.(genai.FunctionCall); ok {
				funcCalls = append(funcCalls, fc)
			}
		}

		if len(funcCalls) == 0 {
			break // No function calls, done
		}

		// Execute function calls
		var funcResponses []genai.Part
		for _, fc := range funcCalls {
			toolCall := &ports.ToolCall{Name: fc.Name, Args: fc.Args}

			// Emit tool start event
			if err := callback(ports.StreamEvent{
				Type:     ports.StreamEventToolStart,
				ToolCall: toolCall,
			}); err != nil {
				return err
			}

			log.Printf("[ADK] Executing tool: %s", fc.Name)
			result, execErr := r.registry.executeTool(ctx, fc.Name, fc.Args)
			if execErr != nil {
				funcResponses = append(funcResponses, genai.FunctionResponse{
					Name:     fc.Name,
					Response: map[string]any{"error": execErr.Error()},
				})
			} else {
				funcResponses = append(funcResponses, genai.FunctionResponse{
					Name:     fc.Name,
					Response: result.Content,
				})
			}

			// Emit tool end event
			if err := callback(ports.StreamEvent{
				Type:     ports.StreamEventToolEnd,
				ToolCall: toolCall,
			}); err != nil {
				return err
			}
		}

		// Send function responses back and loop
		sendParts = funcResponses
	}

	return nil
}

// Config returns the runner configuration
func (r *Runner) Config() *Config {
	return r.config
}

// extractFunctionCalls extracts FunctionCall parts from a response
func extractFunctionCalls(resp *genai.GenerateContentResponse) []genai.FunctionCall {
	var calls []genai.FunctionCall
	if len(resp.Candidates) == 0 {
		return calls
	}
	for _, part := range resp.Candidates[0].Content.Parts {
		if fc, ok := part.(genai.FunctionCall); ok {
			calls = append(calls, fc)
		}
	}
	return calls
}
