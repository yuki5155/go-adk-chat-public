// Package multi provides a multi-provider AI runner dispatcher.
package multi

import (
	"context"
	"fmt"

	"github.com/yuki5155/go-google-auth/internal/application/ports"
	"github.com/yuki5155/go-google-auth/internal/domain/chat"
)

// Dispatcher implements ports.AIRunner by routing each call to the
// correct provider runner based on the model ID.
type Dispatcher struct {
	runners map[string]ports.AIRunner // keyed by provider name
}

// NewDispatcher creates a Dispatcher from a map of provider -> runner.
// Providers with a nil runner are ignored.
func NewDispatcher(runners map[string]ports.AIRunner) *Dispatcher {
	active := make(map[string]ports.AIRunner, len(runners))
	for provider, r := range runners {
		if r != nil {
			active[provider] = r
		}
	}
	return &Dispatcher{runners: active}
}

// ActiveProviders returns the set of provider names that have a configured runner.
func (d *Dispatcher) ActiveProviders() map[string]bool {
	set := make(map[string]bool, len(d.runners))
	for p := range d.runners {
		set[p] = true
	}
	return set
}

func (d *Dispatcher) runnerFor(model string) (ports.AIRunner, error) {
	provider := chat.ProviderForModel(model)
	if provider == "" {
		return nil, fmt.Errorf("unknown model %q", model)
	}
	r, ok := d.runners[provider]
	if !ok {
		return nil, fmt.Errorf("provider %q is not configured (model: %s)", provider, model)
	}
	return r, nil
}

// SendMessage implements ports.AIRunner.
func (d *Dispatcher) SendMessage(ctx context.Context, history []ports.ChatMessage, userMessage string, model string) (string, error) {
	r, err := d.runnerFor(model)
	if err != nil {
		return "", err
	}
	return r.SendMessage(ctx, history, userMessage, model)
}

// StreamMessage implements ports.AIRunner.
func (d *Dispatcher) StreamMessage(ctx context.Context, history []ports.ChatMessage, userMessage string, model string, callback func(chunk string) error) error {
	r, err := d.runnerFor(model)
	if err != nil {
		return err
	}
	return r.StreamMessage(ctx, history, userMessage, model, callback)
}

// StreamMessageWithTools implements ports.AIRunner.
func (d *Dispatcher) StreamMessageWithTools(ctx context.Context, history []ports.ChatMessage, userMessage string, model string, callback ports.StreamEventCallback) error {
	r, err := d.runnerFor(model)
	if err != nil {
		return err
	}
	return r.StreamMessageWithTools(ctx, history, userMessage, model, callback)
}

// RegisterTool registers the tool on all configured runners.
func (d *Dispatcher) RegisterTool(def ports.ToolDefinition, handler ports.ToolHandler) {
	for _, r := range d.runners {
		r.RegisterTool(def, handler)
	}
}

// Config implements ports.AIRunner.
func (d *Dispatcher) Config() *ports.AIRunnerConfig {
	return &ports.AIRunnerConfig{AppName: "go-adk-chat"}
}
