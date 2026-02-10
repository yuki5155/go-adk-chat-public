package chat

import (
	"time"

	"github.com/google/uuid"
)

// ThreadStatus represents the status of a thread
type ThreadStatus string

const (
	ThreadStatusActive   ThreadStatus = "active"
	ThreadStatusArchived ThreadStatus = "archived"
)

// DefaultModel returns the default model ID for chat threads
func DefaultModel() string {
	return GetDefaultModel().ID
}

// Thread represents a conversation thread aggregate root
type Thread struct {
	threadID     string
	userID       string
	title        string
	model        string
	status       ThreadStatus
	messageCount int
	lastMessage  string
	createdAt    time.Time
	updatedAt    time.Time
}

// NewThread creates a new Thread with a generated ID
func NewThread(userID, title, model string) (*Thread, error) {
	if userID == "" {
		return nil, ErrInvalidUserID
	}

	if title == "" {
		title = "New Conversation"
	}

	if model == "" {
		model = DefaultModel()
	}

	now := time.Now()
	return &Thread{
		threadID:     uuid.New().String(),
		userID:       userID,
		title:        title,
		model:        model,
		status:       ThreadStatusActive,
		messageCount: 0,
		lastMessage:  "",
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

// ReconstructThread reconstructs a Thread from persistence
func ReconstructThread(
	threadID, userID, title, model string,
	status ThreadStatus,
	messageCount int,
	lastMessage string,
	createdAt, updatedAt time.Time,
) *Thread {
	if model == "" {
		model = DefaultModel()
	}
	return &Thread{
		threadID:     threadID,
		userID:       userID,
		title:        title,
		model:        model,
		status:       status,
		messageCount: messageCount,
		lastMessage:  lastMessage,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

// ThreadID returns the thread ID
func (t *Thread) ThreadID() string {
	return t.threadID
}

// UserID returns the user ID
func (t *Thread) UserID() string {
	return t.userID
}

// Title returns the thread title
func (t *Thread) Title() string {
	return t.title
}

// Model returns the LLM model for this thread
func (t *Thread) Model() string {
	return t.model
}

// Status returns the thread status
func (t *Thread) Status() ThreadStatus {
	return t.status
}

// MessageCount returns the message count
func (t *Thread) MessageCount() int {
	return t.messageCount
}

// LastMessage returns the last message preview
func (t *Thread) LastMessage() string {
	return t.lastMessage
}

// CreatedAt returns when the thread was created
func (t *Thread) CreatedAt() time.Time {
	return t.createdAt
}

// UpdatedAt returns when the thread was last updated
func (t *Thread) UpdatedAt() time.Time {
	return t.updatedAt
}

// IsActive returns whether the thread is active
func (t *Thread) IsActive() bool {
	return t.status == ThreadStatusActive
}

// UpdateTitle updates the thread title
func (t *Thread) UpdateTitle(title string) error {
	if title == "" {
		return ErrEmptyTitle
	}
	t.title = title
	t.updatedAt = time.Now()
	return nil
}

// UpdateLastMessage updates the last message and increments count
func (t *Thread) UpdateLastMessage(content string) {
	// Truncate for preview (first 100 chars)
	if len(content) > 100 {
		t.lastMessage = content[:100] + "..."
	} else {
		t.lastMessage = content
	}
	t.messageCount++
	t.updatedAt = time.Now()
}

// Archive marks the thread as archived
func (t *Thread) Archive() {
	t.status = ThreadStatusArchived
	t.updatedAt = time.Now()
}

// BelongsTo checks if the thread belongs to the given user
func (t *Thread) BelongsTo(userID string) bool {
	return t.userID == userID
}

// GenerateTitle generates a title from the first message
func (t *Thread) GenerateTitle(firstMessage string) {
	// Use first 50 chars of first message as title
	if len(firstMessage) > 50 {
		t.title = firstMessage[:50] + "..."
	} else {
		t.title = firstMessage
	}
	t.updatedAt = time.Now()
}
