package chat

import (
	"context"
	"fmt"

	"github.com/yuki5155/go-google-auth/internal/domain/chat"
)

const (
	// DefaultThreadListLimit is the default number of threads to return
	DefaultThreadListLimit = 20
	// MaxThreadListLimit is the maximum number of threads to return
	MaxThreadListLimit = 100
)

// ListThreadsUseCase handles listing chat threads for a user
type ListThreadsUseCase struct {
	threadRepo chat.ThreadRepository
}

// NewListThreadsUseCase creates a new ListThreadsUseCase
func NewListThreadsUseCase(threadRepo chat.ThreadRepository) *ListThreadsUseCase {
	return &ListThreadsUseCase{
		threadRepo: threadRepo,
	}
}

// Execute lists threads for a user with pagination
func (uc *ListThreadsUseCase) Execute(ctx context.Context, query ListThreadsQuery) (*ThreadListDTO, error) {
	// Set default limit
	limit := query.Limit
	if limit <= 0 {
		limit = DefaultThreadListLimit
	}
	if limit > MaxThreadListLimit {
		limit = MaxThreadListLimit
	}

	// Fetch threads
	threads, nextKey, err := uc.threadRepo.ListByUser(ctx, query.UserID, limit, query.LastKey)
	if err != nil {
		return nil, fmt.Errorf("failed to list threads: %w", err)
	}

	return &ThreadListDTO{
		Threads: ToThreadDTOs(threads),
		NextKey: nextKey,
	}, nil
}
