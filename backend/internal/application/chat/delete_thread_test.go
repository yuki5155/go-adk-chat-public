package chat

import (
	"context"
	"testing"
	"time"

	"github.com/yuki5155/go-google-auth/internal/domain/chat"
)

func TestNewDeleteThreadUseCase(t *testing.T) {
	threadRepo := &MockThreadRepository{}
	sessionRepo := &MockSessionRepository{}
	eventRepo := &MockEventRepository{}
	memoryRepo := &MockMemoryRepository{}

	uc := NewDeleteThreadUseCase(threadRepo, sessionRepo, eventRepo, memoryRepo)

	if uc == nil {
		t.Fatal("expected use case to be created")
	}
}

func TestDeleteThreadUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		thread := chat.ReconstructThread("thread-123", "user-123", "Test Thread", "", chat.ThreadStatusActive, 2, "Last message", now, now)
		session, _ := chat.NewSession("thread-123", "user-123", "test-app")

		threadRepo := &MockThreadRepository{
			FindByIDFunc: func(ctx context.Context, userID, threadID string) (*chat.Thread, error) {
				return thread, nil
			},
		}
		sessionRepo := &MockSessionRepository{
			ListByThreadFunc: func(ctx context.Context, threadID string) ([]*chat.Session, error) {
				return []*chat.Session{session}, nil
			},
		}
		eventRepo := &MockEventRepository{}
		memoryRepo := &MockMemoryRepository{}

		uc := NewDeleteThreadUseCase(threadRepo, sessionRepo, eventRepo, memoryRepo)

		cmd := DeleteThreadCommand{
			UserID:   "user-123",
			ThreadID: "thread-123",
		}

		err := uc.Execute(ctx, cmd)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !threadRepo.DeleteCalled {
			t.Error("expected Delete to be called on thread repository")
		}
	})

	t.Run("success with no sessions", func(t *testing.T) {
		thread := chat.ReconstructThread("thread-123", "user-123", "Test Thread", "", chat.ThreadStatusActive, 0, "", now, now)

		threadRepo := &MockThreadRepository{
			FindByIDFunc: func(ctx context.Context, userID, threadID string) (*chat.Thread, error) {
				return thread, nil
			},
		}
		sessionRepo := &MockSessionRepository{
			ListByThreadFunc: func(ctx context.Context, threadID string) ([]*chat.Session, error) {
				return []*chat.Session{}, nil
			},
		}
		eventRepo := &MockEventRepository{}
		memoryRepo := &MockMemoryRepository{}

		uc := NewDeleteThreadUseCase(threadRepo, sessionRepo, eventRepo, memoryRepo)

		cmd := DeleteThreadCommand{
			UserID:   "user-123",
			ThreadID: "thread-123",
		}

		err := uc.Execute(ctx, cmd)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("error on thread not found", func(t *testing.T) {
		threadRepo := &MockThreadRepository{
			FindByIDFunc: func(ctx context.Context, userID, threadID string) (*chat.Thread, error) {
				return nil, errMockFind
			},
		}
		sessionRepo := &MockSessionRepository{}
		eventRepo := &MockEventRepository{}
		memoryRepo := &MockMemoryRepository{}

		uc := NewDeleteThreadUseCase(threadRepo, sessionRepo, eventRepo, memoryRepo)

		cmd := DeleteThreadCommand{
			UserID:   "user-123",
			ThreadID: "thread-123",
		}

		err := uc.Execute(ctx, cmd)

		if err == nil {
			t.Error("expected error on thread not found")
		}
	})

	t.Run("error on unauthorized access", func(t *testing.T) {
		thread := chat.ReconstructThread("thread-123", "other-user", "Test Thread", "", chat.ThreadStatusActive, 0, "", now, now)

		threadRepo := &MockThreadRepository{
			FindByIDFunc: func(ctx context.Context, userID, threadID string) (*chat.Thread, error) {
				return thread, nil
			},
		}
		sessionRepo := &MockSessionRepository{}
		eventRepo := &MockEventRepository{}
		memoryRepo := &MockMemoryRepository{}

		uc := NewDeleteThreadUseCase(threadRepo, sessionRepo, eventRepo, memoryRepo)

		cmd := DeleteThreadCommand{
			UserID:   "user-123",
			ThreadID: "thread-123",
		}

		err := uc.Execute(ctx, cmd)

		if err != chat.ErrThreadUnauthorized {
			t.Errorf("expected ErrThreadUnauthorized, got %v", err)
		}
	})

	t.Run("error on event delete failure", func(t *testing.T) {
		thread := chat.ReconstructThread("thread-123", "user-123", "Test Thread", "", chat.ThreadStatusActive, 2, "Last message", now, now)
		session, _ := chat.NewSession("thread-123", "user-123", "test-app")

		threadRepo := &MockThreadRepository{
			FindByIDFunc: func(ctx context.Context, userID, threadID string) (*chat.Thread, error) {
				return thread, nil
			},
		}
		sessionRepo := &MockSessionRepository{
			ListByThreadFunc: func(ctx context.Context, threadID string) ([]*chat.Session, error) {
				return []*chat.Session{session}, nil
			},
		}
		eventRepo := &MockEventRepository{
			DeleteBySessionFunc: func(ctx context.Context, sessionID string) error {
				return errMockDelete
			},
		}
		memoryRepo := &MockMemoryRepository{}

		uc := NewDeleteThreadUseCase(threadRepo, sessionRepo, eventRepo, memoryRepo)

		cmd := DeleteThreadCommand{
			UserID:   "user-123",
			ThreadID: "thread-123",
		}

		err := uc.Execute(ctx, cmd)

		if err == nil {
			t.Error("expected error on event delete failure")
		}
	})

	t.Run("error on memory delete failure", func(t *testing.T) {
		thread := chat.ReconstructThread("thread-123", "user-123", "Test Thread", "", chat.ThreadStatusActive, 0, "", now, now)

		threadRepo := &MockThreadRepository{
			FindByIDFunc: func(ctx context.Context, userID, threadID string) (*chat.Thread, error) {
				return thread, nil
			},
		}
		sessionRepo := &MockSessionRepository{
			ListByThreadFunc: func(ctx context.Context, threadID string) ([]*chat.Session, error) {
				return []*chat.Session{}, nil
			},
		}
		eventRepo := &MockEventRepository{}
		memoryRepo := &MockMemoryRepository{
			DeleteByThreadFunc: func(ctx context.Context, threadID string) error {
				return errMockDelete
			},
		}

		uc := NewDeleteThreadUseCase(threadRepo, sessionRepo, eventRepo, memoryRepo)

		cmd := DeleteThreadCommand{
			UserID:   "user-123",
			ThreadID: "thread-123",
		}

		err := uc.Execute(ctx, cmd)

		if err == nil {
			t.Error("expected error on memory delete failure")
		}
	})

	t.Run("error on thread delete failure", func(t *testing.T) {
		thread := chat.ReconstructThread("thread-123", "user-123", "Test Thread", "", chat.ThreadStatusActive, 0, "", now, now)

		threadRepo := &MockThreadRepository{
			FindByIDFunc: func(ctx context.Context, userID, threadID string) (*chat.Thread, error) {
				return thread, nil
			},
			DeleteFunc: func(ctx context.Context, userID, threadID string) error {
				return errMockDelete
			},
		}
		sessionRepo := &MockSessionRepository{
			ListByThreadFunc: func(ctx context.Context, threadID string) ([]*chat.Session, error) {
				return []*chat.Session{}, nil
			},
		}
		eventRepo := &MockEventRepository{}
		memoryRepo := &MockMemoryRepository{}

		uc := NewDeleteThreadUseCase(threadRepo, sessionRepo, eventRepo, memoryRepo)

		cmd := DeleteThreadCommand{
			UserID:   "user-123",
			ThreadID: "thread-123",
		}

		err := uc.Execute(ctx, cmd)

		if err == nil {
			t.Error("expected error on thread delete failure")
		}
	})
}
