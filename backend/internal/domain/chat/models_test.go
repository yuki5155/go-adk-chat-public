package chat

import "testing"

func TestAvailableModels(t *testing.T) {
	if len(AvailableModels) == 0 {
		t.Error("expected at least one available model")
	}

	// Check that all models have required fields
	for i, model := range AvailableModels {
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

func TestIsValidModel(t *testing.T) {
	tests := []struct {
		modelID string
		want    bool
	}{
		{"gemini-2.0-flash", true},
		{"gemini-2.5-flash", true},
		{"gemini-2.5-pro", true},
		{"invalid-model", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.modelID, func(t *testing.T) {
			got := IsValidModel(tt.modelID)
			if got != tt.want {
				t.Errorf("IsValidModel(%q) = %v, want %v", tt.modelID, got, tt.want)
			}
		})
	}
}

func TestGetDefaultModel(t *testing.T) {
	model := GetDefaultModel()

	if model.ID == "" {
		t.Error("expected default model to have an ID")
	}
	if model.Name == "" {
		t.Error("expected default model to have a Name")
	}

	// Default model should be the first in the list
	if model.ID != AvailableModels[0].ID {
		t.Errorf("expected default model ID %q, got %q", AvailableModels[0].ID, model.ID)
	}
}

func TestDefaultModelFunction(t *testing.T) {
	defaultModelID := DefaultModel()

	if defaultModelID == "" {
		t.Error("expected DefaultModel() to return a non-empty string")
	}

	// Should match the first available model
	if defaultModelID != AvailableModels[0].ID {
		t.Errorf("expected DefaultModel() = %q, got %q", AvailableModels[0].ID, defaultModelID)
	}
}
