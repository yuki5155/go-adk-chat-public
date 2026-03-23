package chat

import (
	"context"
	"testing"

	domainChat "github.com/yuki5155/go-google-auth/internal/domain/chat"
)

func TestNewListModelsUseCase(t *testing.T) {
	uc := NewListModelsUseCase(map[string]bool{"gemini": true})

	if uc == nil {
		t.Fatal("expected use case to be created")
	}
}

func TestListModelsUseCase_Execute(t *testing.T) {
	uc := NewListModelsUseCase(map[string]bool{"gemini": true, "openai": true, "anthropic": true})
	ctx := context.Background()

	result := uc.Execute(ctx)

	if result == nil {
		t.Fatal("expected result to be returned")
	}

	// Should return at least one model
	if len(result.Models) == 0 {
		t.Error("expected at least one model")
	}

	// Default should be set
	if result.Default == "" {
		t.Error("expected default model to be set")
	}

	// Default should be one of the returned models
	found := false
	for _, m := range result.Models {
		if m.ID == result.Default {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("default model %q not found in returned models", result.Default)
	}

	// All returned models should have a known provider
	for _, m := range result.Models {
		if domainChat.ProviderForModel(m.ID) == "" {
			t.Errorf("model %q has unknown provider", m.ID)
		}
	}

	// Check that all models have required fields
	for i, model := range result.Models {
		if model.ID == "" {
			t.Errorf("model %d has empty ID", i)
		}
		if model.Name == "" {
			t.Errorf("model %d has empty Name", i)
		}
		if model.Description == "" {
			t.Errorf("model %d has empty Description", i)
		}
	}
}
