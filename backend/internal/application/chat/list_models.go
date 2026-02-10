package chat

import (
	"context"

	"github.com/yuki5155/go-google-auth/internal/domain/chat"
)

// ListModelsUseCase handles listing available LLM models
type ListModelsUseCase struct{}

// NewListModelsUseCase creates a new ListModelsUseCase
func NewListModelsUseCase() *ListModelsUseCase {
	return &ListModelsUseCase{}
}

// Execute returns the list of available models
func (uc *ListModelsUseCase) Execute(_ context.Context) *ModelListDTO {
	models := make([]ModelDTO, len(chat.AvailableModels))
	for i, m := range chat.AvailableModels {
		models[i] = ModelDTO{
			ID:          m.ID,
			Name:        m.Name,
			Description: m.Description,
		}
	}

	return &ModelListDTO{
		Models:  models,
		Default: chat.DefaultModel(),
	}
}
