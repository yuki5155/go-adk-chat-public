package chat

import (
	"context"
	"fmt"

	"github.com/yuki5155/go-google-auth/internal/domain/chat"
)

// DeleteThreadUseCase handles deleting chat threads
type DeleteThreadUseCase struct {
	threadRepo  chat.ThreadRepository
	sessionRepo chat.SessionRepository
	eventRepo   chat.EventRepository
	memoryRepo  chat.MemoryRepository
}

// NewDeleteThreadUseCase creates a new DeleteThreadUseCase
func NewDeleteThreadUseCase(
	threadRepo chat.ThreadRepository,
	sessionRepo chat.SessionRepository,
	eventRepo chat.EventRepository,
	memoryRepo chat.MemoryRepository,
) *DeleteThreadUseCase {
	return &DeleteThreadUseCase{
		threadRepo:  threadRepo,
		sessionRepo: sessionRepo,
		eventRepo:   eventRepo,
		memoryRepo:  memoryRepo,
	}
}

// Execute deletes a thread and all associated data
func (uc *DeleteThreadUseCase) Execute(ctx context.Context, cmd DeleteThreadCommand) error {
	// Fetch thread
	thread, err := uc.threadRepo.FindByID(ctx, cmd.UserID, cmd.ThreadID)
	if err != nil {
		return fmt.Errorf("failed to get thread: %w", err)
	}

	// Verify ownership
	if !thread.BelongsTo(cmd.UserID) {
		return chat.ErrThreadUnauthorized
	}

	// Get all sessions for the thread
	sessions, err := uc.sessionRepo.ListByThread(ctx, cmd.ThreadID)
	if err != nil {
		return fmt.Errorf("failed to list sessions: %w", err)
	}

	// Delete all events for each session
	for _, session := range sessions {
		if err := uc.eventRepo.DeleteBySession(ctx, session.SessionID()); err != nil {
			return fmt.Errorf("failed to delete events for session %s: %w", session.SessionID(), err)
		}
	}

	// Delete all memories for the thread
	if err := uc.memoryRepo.DeleteByThread(ctx, cmd.ThreadID); err != nil {
		return fmt.Errorf("failed to delete memories: %w", err)
	}

	// Delete the thread
	if err := uc.threadRepo.Delete(ctx, cmd.UserID, cmd.ThreadID); err != nil {
		return fmt.Errorf("failed to delete thread: %w", err)
	}

	return nil
}
