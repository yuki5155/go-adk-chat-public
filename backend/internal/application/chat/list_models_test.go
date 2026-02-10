package chat

import (
	"context"
	"testing"

	domainChat "github.com/yuki5155/go-google-auth/internal/domain/chat"
)

func TestNewListModelsUseCase(t *testing.T) {
	uc := NewListModelsUseCase()

	if uc == nil {
		t.Fatal("expected use case to be created")
	}
}

func TestListModelsUseCase_Execute(t *testing.T) {
	uc := NewListModelsUseCase()
	ctx := context.Background()

	result := uc.Execute(ctx)

	if result == nil {
		t.Fatal("expected result to be returned")
	}

	// Should return at least one model
	if len(result.Models) == 0 {
		t.Error("expected at least one model")
	}

	// Should match domain models count
	if len(result.Models) != len(domainChat.AvailableModels) {
		t.Errorf("expected %d models, got %d", len(domainChat.AvailableModels), len(result.Models))
	}

	// Default should be set
	if result.Default == "" {
		t.Error("expected default model to be set")
	}

	// Default should be the first model
	if result.Default != domainChat.DefaultModel() {
		t.Errorf("expected default %q, got %q", domainChat.DefaultModel(), result.Default)
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
