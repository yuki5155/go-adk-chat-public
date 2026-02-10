package chat

import (
	"context"
	"testing"
	"time"

	"github.com/yuki5155/go-google-auth/internal/domain/chat"
)

func TestNewListThreadsUseCase(t *testing.T) {
	mockRepo := &MockThreadRepository{}
	uc := NewListThreadsUseCase(mockRepo)

	if uc == nil {
		t.Fatal("expected use case to be created")
	}
}

func TestListThreadsUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("success with threads", func(t *testing.T) {
		now := time.Now()
		threads := []*chat.Thread{
			chat.ReconstructThread("thread-1", "user-123", "Thread 1", "", chat.ThreadStatusActive, 5, "Last message 1", now, now),
			chat.ReconstructThread("thread-2", "user-123", "Thread 2", "", chat.ThreadStatusActive, 3, "Last message 2", now, now),
		}

		mockRepo := &MockThreadRepository{
			ListByUserFunc: func(ctx context.Context, userID string, limit int, lastKey string) ([]*chat.Thread, string, error) {
				return threads, "next-key", nil
			},
		}
		uc := NewListThreadsUseCase(mockRepo)

		query := ListThreadsQuery{
			UserID: "user-123",
			Limit:  20,
		}

		result, err := uc.Execute(ctx, query)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Threads) != 2 {
			t.Errorf("expected 2 threads, got %d", len(result.Threads))
		}
		if result.NextKey != "next-key" {
			t.Errorf("expected next key 'next-key', got %s", result.NextKey)
		}
	})

	t.Run("success with empty list", func(t *testing.T) {
		mockRepo := &MockThreadRepository{
			ListByUserFunc: func(ctx context.Context, userID string, limit int, lastKey string) ([]*chat.Thread, string, error) {
				return []*chat.Thread{}, "", nil
			},
		}
		uc := NewListThreadsUseCase(mockRepo)

		query := ListThreadsQuery{
			UserID: "user-123",
		}

		result, err := uc.Execute(ctx, query)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Threads) != 0 {
			t.Errorf("expected 0 threads, got %d", len(result.Threads))
		}
	})

	t.Run("default limit applied", func(t *testing.T) {
		var capturedLimit int
		mockRepo := &MockThreadRepository{
			ListByUserFunc: func(ctx context.Context, userID string, limit int, lastKey string) ([]*chat.Thread, string, error) {
				capturedLimit = limit
				return []*chat.Thread{}, "", nil
			},
		}
		uc := NewListThreadsUseCase(mockRepo)

		query := ListThreadsQuery{
			UserID: "user-123",
			Limit:  0, // Should default to 20
		}

		_, _ = uc.Execute(ctx, query)

		if capturedLimit != 20 {
			t.Errorf("expected default limit 20, got %d", capturedLimit)
		}
	})

	t.Run("max limit enforced", func(t *testing.T) {
		var capturedLimit int
		mockRepo := &MockThreadRepository{
			ListByUserFunc: func(ctx context.Context, userID string, limit int, lastKey string) ([]*chat.Thread, string, error) {
				capturedLimit = limit
				return []*chat.Thread{}, "", nil
			},
		}
		uc := NewListThreadsUseCase(mockRepo)

		query := ListThreadsQuery{
			UserID: "user-123",
			Limit:  500, // Should be capped to 100
		}

		_, _ = uc.Execute(ctx, query)

		if capturedLimit != 100 {
			t.Errorf("expected max limit 100, got %d", capturedLimit)
		}
	})

	t.Run("error on repository failure", func(t *testing.T) {
		mockRepo := &MockThreadRepository{
			ListByUserFunc: func(ctx context.Context, userID string, limit int, lastKey string) ([]*chat.Thread, string, error) {
				return nil, "", errMockFind
			},
		}
		uc := NewListThreadsUseCase(mockRepo)

		query := ListThreadsQuery{
			UserID: "user-123",
		}

		_, err := uc.Execute(ctx, query)

		if err == nil {
			t.Error("expected error on repository failure")
		}
	})
}
