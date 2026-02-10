package chat

import (
	"context"
	"fmt"

	"github.com/yuki5155/go-google-auth/internal/domain/chat"
)

// GetThreadUseCase handles getting a thread with its messages
type GetThreadUseCase struct {
	threadRepo  chat.ThreadRepository
	sessionRepo chat.SessionRepository
	eventRepo   chat.EventRepository
}

// NewGetThreadUseCase creates a new GetThreadUseCase
func NewGetThreadUseCase(
	threadRepo chat.ThreadRepository,
	sessionRepo chat.SessionRepository,
	eventRepo chat.EventRepository,
) *GetThreadUseCase {
	return &GetThreadUseCase{
		threadRepo:  threadRepo,
		sessionRepo: sessionRepo,
		eventRepo:   eventRepo,
	}
}

// Execute gets a thread with its message history
func (uc *GetThreadUseCase) Execute(ctx context.Context, query GetThreadQuery) (*ThreadDetailDTO, error) {
	// Fetch thread
	thread, err := uc.threadRepo.FindByID(ctx, query.UserID, query.ThreadID)
	if err != nil {
		return nil, fmt.Errorf("failed to get thread: %w", err)
	}

	// Verify ownership
	if !thread.BelongsTo(query.UserID) {
		return nil, chat.ErrThreadUnauthorized
	}

	// Get all sessions for the thread
	sessions, err := uc.sessionRepo.ListByThread(ctx, query.ThreadID)
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}

	// Collect all messages from all sessions
	var messages []MessageDTO
	for _, session := range sessions {
		events, err := uc.eventRepo.ListBySession(ctx, session.SessionID(), query.Limit)
		if err != nil {
			return nil, fmt.Errorf("failed to list events: %w", err)
		}
		messages = append(messages, ToMessageDTOs(events)...)
	}

	return &ThreadDetailDTO{
		ThreadDTO: ToThreadDTO(thread),
		Messages:  messages,
	}, nil
}
