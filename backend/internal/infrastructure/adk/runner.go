package adk

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// Runner wraps the Gemini generative AI client for chat operations
type Runner struct {
	config *Config
	client *genai.Client
	model  *genai.GenerativeModel
}

// NewRunner creates a new ADK Runner
func NewRunner(config *Config) (*Runner, error) {
	if !config.IsValid() {
		return nil, fmt.Errorf("invalid ADK configuration: API key or model not set")
	}

	return &Runner{
		config: config,
	}, nil
}

// Initialize initializes the Gemini client
func (r *Runner) Initialize(ctx context.Context) error {
	client, err := genai.NewClient(ctx, option.WithAPIKey(r.config.APIKey))
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

	// Send the user message
	resp, err := cs.SendMessage(ctx, genai.Text(userMessage))
	if err != nil {
		return "", fmt.Errorf("failed to send message: %w", err)
	}

	// Extract text response
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response received from model")
	}

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

// Config returns the runner configuration
func (r *Runner) Config() *Config {
	return r.config
}
