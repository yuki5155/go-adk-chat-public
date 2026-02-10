package chat

import (
	"context"
	"fmt"

	"github.com/yuki5155/go-google-auth/internal/domain/chat"
)

// CreateThreadUseCase handles creating new chat threads
type CreateThreadUseCase struct {
	threadRepo chat.ThreadRepository
}

// NewCreateThreadUseCase creates a new CreateThreadUseCase
func NewCreateThreadUseCase(threadRepo chat.ThreadRepository) *CreateThreadUseCase {
	return &CreateThreadUseCase{
		threadRepo: threadRepo,
	}
}

// Execute creates a new thread
func (uc *CreateThreadUseCase) Execute(ctx context.Context, cmd CreateThreadCommand) (*ThreadDTO, error) {
	// Create domain entity
	thread, err := chat.NewThread(cmd.UserID, cmd.Title, cmd.Model)
	if err != nil {
		return nil, fmt.Errorf("failed to create thread: %w", err)
	}

	// Persist
	if err := uc.threadRepo.Save(ctx, thread); err != nil {
		return nil, fmt.Errorf("failed to save thread: %w", err)
	}

	// Return DTO
	dto := ToThreadDTO(thread)
	return &dto, nil
}
