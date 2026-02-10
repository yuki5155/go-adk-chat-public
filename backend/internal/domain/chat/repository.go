package chat

import "context"

// ThreadRepository defines the interface for thread persistence
type ThreadRepository interface {
	// Save creates or updates a thread
	Save(ctx context.Context, thread *Thread) error

	// FindByID finds a thread by user ID and thread ID
	FindByID(ctx context.Context, userID, threadID string) (*Thread, error)

	// ListByUser lists threads for a user, ordered by updated_at desc
	ListByUser(ctx context.Context, userID string, limit int, lastKey string) ([]*Thread, string, error)

	// Delete deletes a thread
	Delete(ctx context.Context, userID, threadID string) error
}

// SessionRepository defines the interface for session persistence
type SessionRepository interface {
	// Save creates or updates a session
	Save(ctx context.Context, session *Session) error

	// FindByID finds a session by thread ID and session ID
	FindByID(ctx context.Context, threadID, sessionID string) (*Session, error)

	// FindActiveByThread finds the active session for a thread
	FindActiveByThread(ctx context.Context, threadID string) (*Session, error)

	// ListByThread lists sessions for a thread
	ListByThread(ctx context.Context, threadID string) ([]*Session, error)
}

// EventRepository defines the interface for event persistence
type EventRepository interface {
	// Save saves an event
	Save(ctx context.Context, event *Event) error

	// ListBySession lists events for a session, ordered by event_id
	ListBySession(ctx context.Context, sessionID string, limit int) ([]*Event, error)

	// DeleteBySession deletes all events for a session
	DeleteBySession(ctx context.Context, sessionID string) error
}

// MemoryRepository defines the interface for memory persistence
type MemoryRepository interface {
	// Save saves a memory
	Save(ctx context.Context, memory *Memory) error

	// FindByID finds a memory by thread ID and memory ID
	FindByID(ctx context.Context, threadID, memoryID string) (*Memory, error)

	// ListByThread lists memories for a thread
	ListByThread(ctx context.Context, threadID string, limit int) ([]*Memory, error)

	// ListByUser lists memories for a user
	ListByUser(ctx context.Context, userID string, limit int) ([]*Memory, error)

	// DeleteByThread deletes all memories for a thread
	DeleteByThread(ctx context.Context, threadID string) error
}
