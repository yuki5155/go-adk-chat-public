package chat

import (
	"testing"
	"time"
)

func TestNewThread(t *testing.T) {
	userID := "user-123"
	title := "Test Thread"

	thread, err := NewThread(userID, title, "")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if thread.UserID() != userID {
		t.Errorf("expected UserID %s, got %s", userID, thread.UserID())
	}
	if thread.Title() != title {
		t.Errorf("expected Title %s, got %s", title, thread.Title())
	}
	if thread.Status() != ThreadStatusActive {
		t.Errorf("expected Status %s, got %s", ThreadStatusActive, thread.Status())
	}
	if thread.MessageCount() != 0 {
		t.Errorf("expected MessageCount 0, got %d", thread.MessageCount())
	}
	if thread.ThreadID() == "" {
		t.Error("expected ThreadID to be generated")
	}
	if thread.CreatedAt().IsZero() {
		t.Error("expected CreatedAt to be set")
	}
	if thread.UpdatedAt().IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
}

func TestNewThread_EmptyTitle(t *testing.T) {
	userID := "user-123"

	thread, err := NewThread(userID, "", "")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if thread.Title() != "New Conversation" {
		t.Errorf("expected default Title 'New Conversation', got %s", thread.Title())
	}
}

func TestNewThread_EmptyUserID(t *testing.T) {
	_, err := NewThread("", "Test", "")

	if err != ErrInvalidUserID {
		t.Errorf("expected error ErrInvalidUserID, got %v", err)
	}
}

func TestReconstructThread(t *testing.T) {
	threadID := "thread-123"
	userID := "user-123"
	title := "Reconstructed Thread"
	status := ThreadStatusArchived
	messageCount := 5
	lastMessage := "Last message content"
	createdAt := time.Now().Add(-24 * time.Hour)
	updatedAt := time.Now()

	thread := ReconstructThread(threadID, userID, title, "", status, messageCount, lastMessage, createdAt, updatedAt)

	if thread.ThreadID() != threadID {
		t.Errorf("expected ThreadID %s, got %s", threadID, thread.ThreadID())
	}
	if thread.UserID() != userID {
		t.Errorf("expected UserID %s, got %s", userID, thread.UserID())
	}
	if thread.Title() != title {
		t.Errorf("expected Title %s, got %s", title, thread.Title())
	}
	if thread.Status() != status {
		t.Errorf("expected Status %s, got %s", status, thread.Status())
	}
	if thread.MessageCount() != messageCount {
		t.Errorf("expected MessageCount %d, got %d", messageCount, thread.MessageCount())
	}
	if thread.LastMessage() != lastMessage {
		t.Errorf("expected LastMessage %s, got %s", lastMessage, thread.LastMessage())
	}
}

func TestThread_UpdateLastMessage(t *testing.T) {
	thread, _ := NewThread("user-123", "Test", "")
	originalUpdatedAt := thread.UpdatedAt()

	time.Sleep(1 * time.Millisecond)
	thread.UpdateLastMessage("New last message")

	if thread.LastMessage() != "New last message" {
		t.Errorf("expected LastMessage 'New last message', got %s", thread.LastMessage())
	}
	if thread.MessageCount() != 1 {
		t.Errorf("expected MessageCount 1, got %d", thread.MessageCount())
	}
	if !thread.UpdatedAt().After(originalUpdatedAt) {
		t.Error("expected UpdatedAt to be updated")
	}
}

func TestThread_UpdateLastMessage_Truncation(t *testing.T) {
	thread, _ := NewThread("user-123", "Test", "")
	longMessage := "This is a very long message that exceeds one hundred characters and should be truncated with ellipsis at the end of the preview"

	thread.UpdateLastMessage(longMessage)

	if len(thread.LastMessage()) > 103 { // 100 chars + "..."
		t.Errorf("expected LastMessage to be truncated, got length %d", len(thread.LastMessage()))
	}
	if thread.LastMessage()[len(thread.LastMessage())-3:] != "..." {
		t.Error("expected LastMessage to end with '...'")
	}
}

func TestThread_Archive(t *testing.T) {
	thread, _ := NewThread("user-123", "Test", "")

	thread.Archive()

	if thread.Status() != ThreadStatusArchived {
		t.Errorf("expected Status %s, got %s", ThreadStatusArchived, thread.Status())
	}
}

func TestThread_IsActive(t *testing.T) {
	thread, _ := NewThread("user-123", "Test", "")

	if !thread.IsActive() {
		t.Error("expected thread to be active")
	}

	thread.Archive()

	if thread.IsActive() {
		t.Error("expected thread to not be active after archiving")
	}
}

func TestThread_BelongsTo(t *testing.T) {
	thread, _ := NewThread("user-123", "Test", "")

	if !thread.BelongsTo("user-123") {
		t.Error("expected thread to belong to user-123")
	}
	if thread.BelongsTo("other-user") {
		t.Error("expected thread to not belong to other-user")
	}
}

func TestThread_UpdateTitle(t *testing.T) {
	thread, _ := NewThread("user-123", "Original Title", "")

	err := thread.UpdateTitle("New Title")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if thread.Title() != "New Title" {
		t.Errorf("expected Title 'New Title', got %s", thread.Title())
	}
}

func TestThread_UpdateTitle_Empty(t *testing.T) {
	thread, _ := NewThread("user-123", "Original Title", "")

	err := thread.UpdateTitle("")

	if err != ErrEmptyTitle {
		t.Errorf("expected ErrEmptyTitle, got %v", err)
	}
}

func TestThread_GenerateTitle(t *testing.T) {
	tests := []struct {
		name          string
		firstMessage  string
		expectedTitle string
	}{
		{
			name:          "short message",
			firstMessage:  "Hello",
			expectedTitle: "Hello",
		},
		{
			name:          "message under 50 chars",
			firstMessage:  "This is a test message",
			expectedTitle: "This is a test message",
		},
		{
			name:          "message over 50 chars",
			firstMessage:  "This is a very long message that should be truncated because it exceeds fifty characters",
			expectedTitle: "This is a very long message that should be truncat...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			thread, _ := NewThread("user-123", "", "")
			thread.GenerateTitle(tt.firstMessage)

			if thread.Title() != tt.expectedTitle {
				t.Errorf("expected Title %q, got %q", tt.expectedTitle, thread.Title())
			}
		})
	}
}

func TestThreadStatus_Constants(t *testing.T) {
	if ThreadStatusActive != "active" {
		t.Errorf("expected ThreadStatusActive 'active', got %s", ThreadStatusActive)
	}
	if ThreadStatusArchived != "archived" {
		t.Errorf("expected ThreadStatusArchived 'archived', got %s", ThreadStatusArchived)
	}
}
