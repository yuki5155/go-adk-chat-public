package chat

import (
	"context"
	"testing"
	"time"

	"github.com/yuki5155/go-google-auth/internal/application/ports"
	"github.com/yuki5155/go-google-auth/internal/domain/chat"
)

// Tests for SendMessageUseCase

func TestNewSendMessageUseCase(t *testing.T) {
	threadRepo := &MockThreadRepository{}
	sessionRepo := &MockSessionRepository{}
	eventRepo := &MockEventRepository{}
	aiRunner := &MockAIRunner{}

	uc := NewSendMessageUseCase(threadRepo, sessionRepo, eventRepo, aiRunner)

	if uc == nil {
		t.Fatal("expected use case to be created")
	}
}

func TestSendMessageUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("success with new session", func(t *testing.T) {
		thread := chat.ReconstructThread("thread-123", "user-123", "Test Thread", "", chat.ThreadStatusActive, 0, "", now, now)

		threadRepo := &MockThreadRepository{
			FindByIDFunc: func(ctx context.Context, userID, threadID string) (*chat.Thread, error) {
				return thread, nil
			},
		}
		sessionRepo := &MockSessionRepository{
			FindActiveByThreadFunc: func(ctx context.Context, threadID string) (*chat.Session, error) {
				return nil, nil // No active session
			},
		}
		eventRepo := &MockEventRepository{
			ListBySessionFunc: func(ctx context.Context, sessionID string, limit int) ([]*chat.Event, error) {
				return []*chat.Event{}, nil
			},
		}
		aiRunner := &MockAIRunner{
			SendMessageFunc: func(ctx context.Context, history []ports.ChatMessage, userMessage string, model string) (string, error) {
				return "Hello from AI!", nil
			},
		}

		uc := NewSendMessageUseCase(threadRepo, sessionRepo, eventRepo, aiRunner)

		cmd := SendMessageCommand{
			UserID:   "user-123",
			ThreadID: "thread-123",
			Content:  "Hello",
		}

		result, err := uc.Execute(ctx, cmd)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil {
			t.Fatal("expected result")
		}
		if result.Message.Content != "Hello" {
			t.Errorf("expected message content 'Hello', got %s", result.Message.Content)
		}
		if result.Response.Content != "Hello from AI!" {
			t.Errorf("expected response content 'Hello from AI!', got %s", result.Response.Content)
		}
	})

	t.Run("success with existing session", func(t *testing.T) {
		thread := chat.ReconstructThread("thread-123", "user-123", "Test Thread", "", chat.ThreadStatusActive, 2, "Last msg", now, now)
		session, _ := chat.NewSession("thread-123", "user-123", "test-app")

		threadRepo := &MockThreadRepository{
			FindByIDFunc: func(ctx context.Context, userID, threadID string) (*chat.Thread, error) {
				return thread, nil
			},
		}
		sessionRepo := &MockSessionRepository{
			FindActiveByThreadFunc: func(ctx context.Context, threadID string) (*chat.Session, error) {
				return session, nil
			},
		}
		eventRepo := &MockEventRepository{
			ListBySessionFunc: func(ctx context.Context, sessionID string, limit int) ([]*chat.Event, error) {
				return []*chat.Event{}, nil
			},
		}
		aiRunner := &MockAIRunner{}

		uc := NewSendMessageUseCase(threadRepo, sessionRepo, eventRepo, aiRunner)

		cmd := SendMessageCommand{
			UserID:   "user-123",
			ThreadID: "thread-123",
			Content:  "Follow up message",
		}

		result, err := uc.Execute(ctx, cmd)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil {
			t.Fatal("expected result")
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
		aiRunner := &MockAIRunner{}

		uc := NewSendMessageUseCase(threadRepo, sessionRepo, eventRepo, aiRunner)

		cmd := SendMessageCommand{
			UserID:   "user-123",
			ThreadID: "thread-123",
			Content:  "Hello",
		}

		_, err := uc.Execute(ctx, cmd)

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
		aiRunner := &MockAIRunner{}

		uc := NewSendMessageUseCase(threadRepo, sessionRepo, eventRepo, aiRunner)

		cmd := SendMessageCommand{
			UserID:   "user-123",
			ThreadID: "thread-123",
			Content:  "Hello",
		}

		_, err := uc.Execute(ctx, cmd)

		if err != chat.ErrThreadUnauthorized {
			t.Errorf("expected ErrThreadUnauthorized, got %v", err)
		}
	})

	t.Run("error on AI failure", func(t *testing.T) {
		thread := chat.ReconstructThread("thread-123", "user-123", "Test Thread", "", chat.ThreadStatusActive, 0, "", now, now)

		threadRepo := &MockThreadRepository{
			FindByIDFunc: func(ctx context.Context, userID, threadID string) (*chat.Thread, error) {
				return thread, nil
			},
		}
		sessionRepo := &MockSessionRepository{
			FindActiveByThreadFunc: func(ctx context.Context, threadID string) (*chat.Session, error) {
				return nil, nil
			},
		}
		eventRepo := &MockEventRepository{
			ListBySessionFunc: func(ctx context.Context, sessionID string, limit int) ([]*chat.Event, error) {
				return []*chat.Event{}, nil
			},
		}
		aiRunner := &MockAIRunner{
			SendMessageFunc: func(ctx context.Context, history []ports.ChatMessage, userMessage string, model string) (string, error) {
				return "", errMockAI
			},
		}

		uc := NewSendMessageUseCase(threadRepo, sessionRepo, eventRepo, aiRunner)

		cmd := SendMessageCommand{
			UserID:   "user-123",
			ThreadID: "thread-123",
			Content:  "Hello",
		}

		_, err := uc.Execute(ctx, cmd)

		if err == nil {
			t.Error("expected error on AI failure")
		}
	})
}

func TestSendMessageUseCase_ExecuteStream(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("success streaming", func(t *testing.T) {
		thread := chat.ReconstructThread("thread-123", "user-123", "Test Thread", "", chat.ThreadStatusActive, 0, "", now, now)

		threadRepo := &MockThreadRepository{
			FindByIDFunc: func(ctx context.Context, userID, threadID string) (*chat.Thread, error) {
				return thread, nil
			},
		}
		sessionRepo := &MockSessionRepository{
			FindActiveByThreadFunc: func(ctx context.Context, threadID string) (*chat.Session, error) {
				return nil, nil
			},
		}
		eventRepo := &MockEventRepository{
			ListBySessionFunc: func(ctx context.Context, sessionID string, limit int) ([]*chat.Event, error) {
				return []*chat.Event{}, nil
			},
		}
		aiRunner := &MockAIRunner{
			StreamMessageFunc: func(ctx context.Context, history []ports.ChatMessage, userMessage string, model string, callback func(chunk string) error) error {
				chunks := []string{"Hello", " from", " streaming!"}
				for _, chunk := range chunks {
					if err := callback(chunk); err != nil {
						return err
					}
				}
				return nil
			},
		}

		uc := NewSendMessageUseCase(threadRepo, sessionRepo, eventRepo, aiRunner)

		cmd := SendMessageCommand{
			UserID:   "user-123",
			ThreadID: "thread-123",
			Content:  "Hello",
		}

		var receivedChunks []string
		result, err := uc.ExecuteStream(ctx, cmd, func(chunk string) error {
			receivedChunks = append(receivedChunks, chunk)
			return nil
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil {
			t.Fatal("expected result")
		}
		if len(receivedChunks) != 3 {
			t.Errorf("expected 3 chunks, got %d", len(receivedChunks))
		}
		if result.Response.Content != "Hello from streaming!" {
			t.Errorf("expected response content 'Hello from streaming!', got %s", result.Response.Content)
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
		aiRunner := &MockAIRunner{}

		uc := NewSendMessageUseCase(threadRepo, sessionRepo, eventRepo, aiRunner)

		cmd := SendMessageCommand{
			UserID:   "user-123",
			ThreadID: "thread-123",
			Content:  "Hello",
		}

		_, err := uc.ExecuteStream(ctx, cmd, func(chunk string) error { return nil })

		if err != chat.ErrThreadUnauthorized {
			t.Errorf("expected ErrThreadUnauthorized, got %v", err)
		}
	})

	t.Run("error on stream failure", func(t *testing.T) {
		thread := chat.ReconstructThread("thread-123", "user-123", "Test Thread", "", chat.ThreadStatusActive, 0, "", now, now)

		threadRepo := &MockThreadRepository{
			FindByIDFunc: func(ctx context.Context, userID, threadID string) (*chat.Thread, error) {
				return thread, nil
			},
		}
		sessionRepo := &MockSessionRepository{
			FindActiveByThreadFunc: func(ctx context.Context, threadID string) (*chat.Session, error) {
				return nil, nil
			},
		}
		eventRepo := &MockEventRepository{
			ListBySessionFunc: func(ctx context.Context, sessionID string, limit int) ([]*chat.Event, error) {
				return []*chat.Event{}, nil
			},
		}
		aiRunner := &MockAIRunner{
			StreamMessageFunc: func(ctx context.Context, history []ports.ChatMessage, userMessage string, model string, callback func(chunk string) error) error {
				return errMockAI
			},
		}

		uc := NewSendMessageUseCase(threadRepo, sessionRepo, eventRepo, aiRunner)

		cmd := SendMessageCommand{
			UserID:   "user-123",
			ThreadID: "thread-123",
			Content:  "Hello",
		}

		_, err := uc.ExecuteStream(ctx, cmd, func(chunk string) error { return nil })

		if err == nil {
			t.Error("expected error on stream failure")
		}
	})
}
