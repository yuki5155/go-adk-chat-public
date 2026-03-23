package chat

import (
	"context"

	"github.com/yuki5155/go-google-auth/internal/domain/chat"
)

// ListModelsUseCase handles listing available LLM models
type ListModelsUseCase struct {
	activeProviders map[string]bool
}

// NewListModelsUseCase creates a new ListModelsUseCase.
// activeProviders is the set of provider names that have a configured runner.
func NewListModelsUseCase(activeProviders map[string]bool) *ListModelsUseCase {
	return &ListModelsUseCase{activeProviders: activeProviders}
}

// Execute returns the list of models for all active providers
func (uc *ListModelsUseCase) Execute(_ context.Context) *ModelListDTO {
	var available []chat.Model
	for _, m := range chat.AvailableModels {
		if uc.activeProviders[m.Provider] {
			available = append(available, m)
		}
	}
	if len(available) == 0 {
		available = chat.AvailableModels
	}

	models := make([]ModelDTO, len(available))
	for i, m := range available {
		models[i] = ModelDTO{
			ID:          m.ID,
			Name:        m.Name,
			Description: m.Description,
		}
	}

	return &ModelListDTO{
		Models:  models,
		Default: available[0].ID,
	}
}
