package chat

import (
	"context"
	"errors"

	"github.com/yuki5155/go-google-auth/internal/application/ports"
	"github.com/yuki5155/go-google-auth/internal/domain/chat"
)

// MockThreadRepository is a mock implementation of chat.ThreadRepository
type MockThreadRepository struct {
	SaveFunc       func(ctx context.Context, thread *chat.Thread) error
	FindByIDFunc   func(ctx context.Context, userID, threadID string) (*chat.Thread, error)
	ListByUserFunc func(ctx context.Context, userID string, limit int, lastKey string) ([]*chat.Thread, string, error)
	DeleteFunc     func(ctx context.Context, userID, threadID string) error

	SaveCalled     bool
	FindByIDCalled bool
	DeleteCalled   bool
}

func (m *MockThreadRepository) Save(ctx context.Context, thread *chat.Thread) error {
	m.SaveCalled = true
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, thread)
	}
	return nil
}

func (m *MockThreadRepository) FindByID(ctx context.Context, userID, threadID string) (*chat.Thread, error) {
	m.FindByIDCalled = true
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, userID, threadID)
	}
	return nil, nil
}

func (m *MockThreadRepository) ListByUser(ctx context.Context, userID string, limit int, lastKey string) ([]*chat.Thread, string, error) {
	if m.ListByUserFunc != nil {
		return m.ListByUserFunc(ctx, userID, limit, lastKey)
	}
	return nil, "", nil
}

func (m *MockThreadRepository) Delete(ctx context.Context, userID, threadID string) error {
	m.DeleteCalled = true
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, userID, threadID)
	}
	return nil
}

// MockSessionRepository is a mock implementation of chat.SessionRepository
type MockSessionRepository struct {
	SaveFunc               func(ctx context.Context, session *chat.Session) error
	FindByIDFunc           func(ctx context.Context, threadID, sessionID string) (*chat.Session, error)
	FindActiveByThreadFunc func(ctx context.Context, threadID string) (*chat.Session, error)
	ListByThreadFunc       func(ctx context.Context, threadID string) ([]*chat.Session, error)

	SaveCalled bool
}

func (m *MockSessionRepository) Save(ctx context.Context, session *chat.Session) error {
	m.SaveCalled = true
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, session)
	}
	return nil
}

func (m *MockSessionRepository) FindByID(ctx context.Context, threadID, sessionID string) (*chat.Session, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, threadID, sessionID)
	}
	return nil, nil
}

func (m *MockSessionRepository) FindActiveByThread(ctx context.Context, threadID string) (*chat.Session, error) {
	if m.FindActiveByThreadFunc != nil {
		return m.FindActiveByThreadFunc(ctx, threadID)
	}
	return nil, nil
}

func (m *MockSessionRepository) ListByThread(ctx context.Context, threadID string) ([]*chat.Session, error) {
	if m.ListByThreadFunc != nil {
		return m.ListByThreadFunc(ctx, threadID)
	}
	return nil, nil
}

// MockEventRepository is a mock implementation of chat.EventRepository
type MockEventRepository struct {
	SaveFunc            func(ctx context.Context, event *chat.Event) error
	ListBySessionFunc   func(ctx context.Context, sessionID string, limit int) ([]*chat.Event, error)
	DeleteBySessionFunc func(ctx context.Context, sessionID string) error

	SaveCalled bool
}

func (m *MockEventRepository) Save(ctx context.Context, event *chat.Event) error {
	m.SaveCalled = true
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, event)
	}
	return nil
}

func (m *MockEventRepository) ListBySession(ctx context.Context, sessionID string, limit int) ([]*chat.Event, error) {
	if m.ListBySessionFunc != nil {
		return m.ListBySessionFunc(ctx, sessionID, limit)
	}
	return nil, nil
}

func (m *MockEventRepository) DeleteBySession(ctx context.Context, sessionID string) error {
	if m.DeleteBySessionFunc != nil {
		return m.DeleteBySessionFunc(ctx, sessionID)
	}
	return nil
}

// MockMemoryRepository is a mock implementation of chat.MemoryRepository
type MockMemoryRepository struct {
	SaveFunc           func(ctx context.Context, memory *chat.Memory) error
	FindByIDFunc       func(ctx context.Context, threadID, memoryID string) (*chat.Memory, error)
	ListByThreadFunc   func(ctx context.Context, threadID string, limit int) ([]*chat.Memory, error)
	ListByUserFunc     func(ctx context.Context, userID string, limit int) ([]*chat.Memory, error)
	DeleteByThreadFunc func(ctx context.Context, threadID string) error
}

func (m *MockMemoryRepository) Save(ctx context.Context, memory *chat.Memory) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, memory)
	}
	return nil
}

func (m *MockMemoryRepository) FindByID(ctx context.Context, threadID, memoryID string) (*chat.Memory, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, threadID, memoryID)
	}
	return nil, nil
}

func (m *MockMemoryRepository) ListByThread(ctx context.Context, threadID string, limit int) ([]*chat.Memory, error) {
	if m.ListByThreadFunc != nil {
		return m.ListByThreadFunc(ctx, threadID, limit)
	}
	return nil, nil
}

func (m *MockMemoryRepository) ListByUser(ctx context.Context, userID string, limit int) ([]*chat.Memory, error) {
	if m.ListByUserFunc != nil {
		return m.ListByUserFunc(ctx, userID, limit)
	}
	return nil, nil
}

func (m *MockMemoryRepository) DeleteByThread(ctx context.Context, threadID string) error {
	if m.DeleteByThreadFunc != nil {
		return m.DeleteByThreadFunc(ctx, threadID)
	}
	return nil
}

// MockAIRunner is a mock implementation of ports.AIRunner interface
type MockAIRunner struct {
	SendMessageFunc   func(ctx context.Context, history []ports.ChatMessage, userMessage string, model string) (string, error)
	StreamMessageFunc func(ctx context.Context, history []ports.ChatMessage, userMessage string, model string, callback func(chunk string) error) error
	ConfigFunc        func() *ports.AIRunnerConfig
}

func (m *MockAIRunner) SendMessage(ctx context.Context, history []ports.ChatMessage, userMessage string, model string) (string, error) {
	if m.SendMessageFunc != nil {
		return m.SendMessageFunc(ctx, history, userMessage, model)
	}
	return "Mock AI response", nil
}

func (m *MockAIRunner) StreamMessage(ctx context.Context, history []ports.ChatMessage, userMessage string, model string, callback func(chunk string) error) error {
	if m.StreamMessageFunc != nil {
		return m.StreamMessageFunc(ctx, history, userMessage, model, callback)
	}
	// Default: send response in chunks
	chunks := []string{"Hello", " from", " AI!"}
	for _, chunk := range chunks {
		if err := callback(chunk); err != nil {
			return err
		}
	}
	return nil
}

func (m *MockAIRunner) Config() *ports.AIRunnerConfig {
	if m.ConfigFunc != nil {
		return m.ConfigFunc()
	}
	return &ports.AIRunnerConfig{
		AppName: "test-app",
	}
}

// Common test errors
var (
	errMockSave   = errors.New("mock save error")
	errMockFind   = errors.New("mock find error")
	errMockDelete = errors.New("mock delete error")
	errMockAI     = errors.New("mock AI error")
)
