package chat

import (
	"context"
	"testing"
	"time"

	"github.com/yuki5155/go-google-auth/internal/domain/chat"
)

func TestNewGetThreadUseCase(t *testing.T) {
	threadRepo := &MockThreadRepository{}
	sessionRepo := &MockSessionRepository{}
	eventRepo := &MockEventRepository{}

	uc := NewGetThreadUseCase(threadRepo, sessionRepo, eventRepo)

	if uc == nil {
		t.Fatal("expected use case to be created")
	}
}

func TestGetThreadUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("success with messages", func(t *testing.T) {
		thread := chat.ReconstructThread("thread-123", "user-123", "Test Thread", "", chat.ThreadStatusActive, 2, "Last message", now, now)
		session, _ := chat.NewSession("thread-123", "user-123", "test-app")

		event1 := chat.ReconstructEvent("event-1", session.SessionID(), "thread-123", chat.EventRoleUser, "Hello", "user", "inv-1", now)
		event2 := chat.ReconstructEvent("event-2", session.SessionID(), "thread-123", chat.EventRoleAssistant, "Hi there!", "assistant", "inv-2", now)

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
			ListBySessionFunc: func(ctx context.Context, sessionID string, limit int) ([]*chat.Event, error) {
				return []*chat.Event{event1, event2}, nil
			},
		}

		uc := NewGetThreadUseCase(threadRepo, sessionRepo, eventRepo)

		query := GetThreadQuery{
			UserID:   "user-123",
			ThreadID: "thread-123",
		}

		result, err := uc.Execute(ctx, query)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.ThreadID != "thread-123" {
			t.Errorf("expected thread ID 'thread-123', got %s", result.ThreadID)
		}
		if len(result.Messages) != 2 {
			t.Errorf("expected 2 messages, got %d", len(result.Messages))
		}
		if result.Messages[0].Role != "user" {
			t.Errorf("expected first message role 'user', got %s", result.Messages[0].Role)
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

		uc := NewGetThreadUseCase(threadRepo, sessionRepo, eventRepo)

		query := GetThreadQuery{
			UserID:   "user-123",
			ThreadID: "thread-123",
		}

		result, err := uc.Execute(ctx, query)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Messages) != 0 {
			t.Errorf("expected 0 messages, got %d", len(result.Messages))
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

		uc := NewGetThreadUseCase(threadRepo, sessionRepo, eventRepo)

		query := GetThreadQuery{
			UserID:   "user-123",
			ThreadID: "thread-123",
		}

		_, err := uc.Execute(ctx, query)

		if err == nil {
			t.Error("expected error on thread not found")
		}
	})

	t.Run("error on unauthorized access", func(t *testing.T) {
		// Thread belongs to different user
		thread := chat.ReconstructThread("thread-123", "other-user", "Test Thread", "", chat.ThreadStatusActive, 0, "", now, now)

		threadRepo := &MockThreadRepository{
			FindByIDFunc: func(ctx context.Context, userID, threadID string) (*chat.Thread, error) {
				return thread, nil
			},
		}
		sessionRepo := &MockSessionRepository{}
		eventRepo := &MockEventRepository{}

		uc := NewGetThreadUseCase(threadRepo, sessionRepo, eventRepo)

		query := GetThreadQuery{
			UserID:   "user-123",
			ThreadID: "thread-123",
		}

		_, err := uc.Execute(ctx, query)

		if err != chat.ErrThreadUnauthorized {
			t.Errorf("expected ErrThreadUnauthorized, got %v", err)
		}
	})

	t.Run("error on session list failure", func(t *testing.T) {
		thread := chat.ReconstructThread("thread-123", "user-123", "Test Thread", "", chat.ThreadStatusActive, 0, "", now, now)

		threadRepo := &MockThreadRepository{
			FindByIDFunc: func(ctx context.Context, userID, threadID string) (*chat.Thread, error) {
				return thread, nil
			},
		}
		sessionRepo := &MockSessionRepository{
			ListByThreadFunc: func(ctx context.Context, threadID string) ([]*chat.Session, error) {
				return nil, errMockFind
			},
		}
		eventRepo := &MockEventRepository{}

		uc := NewGetThreadUseCase(threadRepo, sessionRepo, eventRepo)

		query := GetThreadQuery{
			UserID:   "user-123",
			ThreadID: "thread-123",
		}

		_, err := uc.Execute(ctx, query)

		if err == nil {
			t.Error("expected error on session list failure")
		}
	})
}
