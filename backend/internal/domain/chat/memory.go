package chat

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ImportanceLevel represents the importance of a memory
type ImportanceLevel int

const (
	ImportanceLow    ImportanceLevel = 1
	ImportanceMedium ImportanceLevel = 2
	ImportanceHigh   ImportanceLevel = 3
)

// Memory represents extracted knowledge from conversations
type Memory struct {
	memoryID        string
	threadID        string
	userID          string
	content         string
	keywords        []string
	sourceSessionID string
	importance      ImportanceLevel
	timestamp       time.Time
}

// NewMemory creates a new Memory
func NewMemory(threadID, userID, content string, keywords []string, sourceSessionID string, importance ImportanceLevel) (*Memory, error) {
	if threadID == "" {
		return nil, ErrInvalidThreadID
	}
	if userID == "" {
		return nil, ErrInvalidUserID
	}
	if content == "" {
		return nil, ErrEmptyContent
	}

	now := time.Now()
	// Generate sortable memory ID (timestamp prefix + uuid)
	memoryID := fmt.Sprintf("%d_%s", now.UnixNano(), uuid.New().String())

	return &Memory{
		memoryID:        memoryID,
		threadID:        threadID,
		userID:          userID,
		content:         content,
		keywords:        keywords,
		sourceSessionID: sourceSessionID,
		importance:      importance,
		timestamp:       now,
	}, nil
}

// ReconstructMemory reconstructs a Memory from persistence
func ReconstructMemory(
	memoryID, threadID, userID, content string,
	keywords []string,
	sourceSessionID string,
	importance ImportanceLevel,
	timestamp time.Time,
) *Memory {
	return &Memory{
		memoryID:        memoryID,
		threadID:        threadID,
		userID:          userID,
		content:         content,
		keywords:        keywords,
		sourceSessionID: sourceSessionID,
		importance:      importance,
		timestamp:       timestamp,
	}
}

// Getters for Memory
func (m *Memory) MemoryID() string        { return m.memoryID }
func (m *Memory) ThreadID() string        { return m.threadID }
func (m *Memory) UserID() string          { return m.userID }
func (m *Memory) Content() string         { return m.content }
func (m *Memory) Keywords() []string      { return m.keywords }
func (m *Memory) SourceSessionID() string { return m.sourceSessionID }
func (m *Memory) Importance() ImportanceLevel { return m.importance }
func (m *Memory) Timestamp() time.Time    { return m.timestamp }

// BelongsTo checks if the memory belongs to the given user
func (m *Memory) BelongsTo(userID string) bool {
	return m.userID == userID
}

// HasKeyword checks if the memory contains a specific keyword
func (m *Memory) HasKeyword(keyword string) bool {
	for _, k := range m.keywords {
		if k == keyword {
			return true
		}
	}
	return false
}
