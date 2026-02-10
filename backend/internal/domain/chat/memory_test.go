package chat

import (
	"testing"
	"time"
)

func TestNewMemory(t *testing.T) {
	threadID := "thread-123"
	userID := "user-123"
	content := "Important information to remember"
	keywords := []string{"important", "remember"}
	sourceSessionID := "session-123"
	importance := ImportanceHigh

	memory, err := NewMemory(threadID, userID, content, keywords, sourceSessionID, importance)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if memory.ThreadID() != threadID {
		t.Errorf("expected ThreadID %s, got %s", threadID, memory.ThreadID())
	}
	if memory.UserID() != userID {
		t.Errorf("expected UserID %s, got %s", userID, memory.UserID())
	}
	if memory.Content() != content {
		t.Errorf("expected Content %s, got %s", content, memory.Content())
	}
	if len(memory.Keywords()) != len(keywords) {
		t.Errorf("expected %d keywords, got %d", len(keywords), len(memory.Keywords()))
	}
	if memory.SourceSessionID() != sourceSessionID {
		t.Errorf("expected SourceSessionID %s, got %s", sourceSessionID, memory.SourceSessionID())
	}
	if memory.Importance() != importance {
		t.Errorf("expected Importance %d, got %d", importance, memory.Importance())
	}
	if memory.MemoryID() == "" {
		t.Error("expected MemoryID to be generated")
	}
	if memory.Timestamp().IsZero() {
		t.Error("expected Timestamp to be set")
	}
}

func TestNewMemory_Validation(t *testing.T) {
	tests := []struct {
		name        string
		threadID    string
		userID      string
		content     string
		expectedErr error
	}{
		{
			name:        "empty thread ID",
			threadID:    "",
			userID:      "user-123",
			content:     "content",
			expectedErr: ErrInvalidThreadID,
		},
		{
			name:        "empty user ID",
			threadID:    "thread-123",
			userID:      "",
			content:     "content",
			expectedErr: ErrInvalidUserID,
		},
		{
			name:        "empty content",
			threadID:    "thread-123",
			userID:      "user-123",
			content:     "",
			expectedErr: ErrEmptyContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewMemory(tt.threadID, tt.userID, tt.content, nil, "", ImportanceMedium)
			if err != tt.expectedErr {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestReconstructMemory(t *testing.T) {
	memoryID := "memory-123"
	threadID := "thread-123"
	userID := "user-123"
	content := "Reconstructed memory"
	keywords := []string{"test", "memory"}
	sourceSessionID := "session-123"
	importance := ImportanceHigh
	timestamp := time.Now()

	memory := ReconstructMemory(memoryID, threadID, userID, content, keywords, sourceSessionID, importance, timestamp)

	if memory.MemoryID() != memoryID {
		t.Errorf("expected MemoryID %s, got %s", memoryID, memory.MemoryID())
	}
	if memory.ThreadID() != threadID {
		t.Errorf("expected ThreadID %s, got %s", threadID, memory.ThreadID())
	}
	if memory.UserID() != userID {
		t.Errorf("expected UserID %s, got %s", userID, memory.UserID())
	}
	if memory.Content() != content {
		t.Errorf("expected Content %s, got %s", content, memory.Content())
	}
	if memory.Importance() != importance {
		t.Errorf("expected Importance %d, got %d", importance, memory.Importance())
	}
	if !memory.Timestamp().Equal(timestamp) {
		t.Errorf("expected Timestamp %v, got %v", timestamp, memory.Timestamp())
	}
}

func TestMemory_HasKeyword(t *testing.T) {
	memory, _ := NewMemory("thread-123", "user-123", "Content", []string{"test", "keyword", "example"}, "session-123", ImportanceMedium)

	tests := []struct {
		keyword  string
		expected bool
	}{
		{"test", true},
		{"keyword", true},
		{"example", true},
		{"notfound", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.keyword, func(t *testing.T) {
			result := memory.HasKeyword(tt.keyword)
			if result != tt.expected {
				t.Errorf("HasKeyword(%s) = %v, expected %v", tt.keyword, result, tt.expected)
			}
		})
	}
}

func TestMemory_BelongsTo(t *testing.T) {
	memory, _ := NewMemory("thread-123", "user-123", "Content", nil, "session-123", ImportanceMedium)

	if !memory.BelongsTo("user-123") {
		t.Error("expected memory to belong to user-123")
	}
	if memory.BelongsTo("other-user") {
		t.Error("expected memory to not belong to other-user")
	}
}

func TestImportanceLevel_Constants(t *testing.T) {
	if ImportanceLow != 1 {
		t.Errorf("expected ImportanceLow 1, got %d", ImportanceLow)
	}
	if ImportanceMedium != 2 {
		t.Errorf("expected ImportanceMedium 2, got %d", ImportanceMedium)
	}
	if ImportanceHigh != 3 {
		t.Errorf("expected ImportanceHigh 3, got %d", ImportanceHigh)
	}
}
