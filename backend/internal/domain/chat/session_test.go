package chat

import (
	"testing"
	"time"
)

func TestNewSession(t *testing.T) {
	threadID := "thread-123"
	userID := "user-123"
	appName := "test-app"

	session, err := NewSession(threadID, userID, appName)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if session.ThreadID() != threadID {
		t.Errorf("expected ThreadID %s, got %s", threadID, session.ThreadID())
	}
	if session.UserID() != userID {
		t.Errorf("expected UserID %s, got %s", userID, session.UserID())
	}
	if session.AppName() != appName {
		t.Errorf("expected AppName %s, got %s", appName, session.AppName())
	}
	if session.SessionID() == "" {
		t.Error("expected SessionID to be generated")
	}
	if session.State() != SessionStateActive {
		t.Errorf("expected State %s, got %s", SessionStateActive, session.State())
	}
	if session.EventCount() != 0 {
		t.Errorf("expected EventCount 0, got %d", session.EventCount())
	}
}

func TestNewSession_Validation(t *testing.T) {
	tests := []struct {
		name        string
		threadID    string
		userID      string
		expectedErr error
	}{
		{
			name:        "empty thread ID",
			threadID:    "",
			userID:      "user-123",
			expectedErr: ErrInvalidThreadID,
		},
		{
			name:        "empty user ID",
			threadID:    "thread-123",
			userID:      "",
			expectedErr: ErrInvalidUserID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewSession(tt.threadID, tt.userID, "test-app")
			if err != tt.expectedErr {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestReconstructSession(t *testing.T) {
	sessionID := "session-123"
	threadID := "thread-123"
	userID := "user-123"
	appName := "test-app"
	state := SessionStateComplete
	eventCount := 10
	createdAt := time.Now().Add(-1 * time.Hour)
	updatedAt := time.Now()

	session := ReconstructSession(sessionID, threadID, userID, appName, state, eventCount, createdAt, updatedAt)

	if session.SessionID() != sessionID {
		t.Errorf("expected SessionID %s, got %s", sessionID, session.SessionID())
	}
	if session.ThreadID() != threadID {
		t.Errorf("expected ThreadID %s, got %s", threadID, session.ThreadID())
	}
	if session.State() != state {
		t.Errorf("expected State %s, got %s", state, session.State())
	}
	if session.EventCount() != eventCount {
		t.Errorf("expected EventCount %d, got %d", eventCount, session.EventCount())
	}
}

func TestSession_IncrementEventCount(t *testing.T) {
	session, _ := NewSession("thread-123", "user-123", "test-app")
	originalEventCount := session.EventCount()

	session.IncrementEventCount()

	if session.EventCount() != originalEventCount+1 {
		t.Errorf("expected EventCount %d, got %d", originalEventCount+1, session.EventCount())
	}
}

func TestSession_Complete(t *testing.T) {
	session, _ := NewSession("thread-123", "user-123", "test-app")

	session.Complete()

	if session.State() != SessionStateComplete {
		t.Errorf("expected State %s, got %s", SessionStateComplete, session.State())
	}
}

func TestSession_MarkError(t *testing.T) {
	session, _ := NewSession("thread-123", "user-123", "test-app")

	session.MarkError()

	if session.State() != SessionStateError {
		t.Errorf("expected State %s, got %s", SessionStateError, session.State())
	}
}

func TestSession_IsActive(t *testing.T) {
	session, _ := NewSession("thread-123", "user-123", "test-app")

	if !session.IsActive() {
		t.Error("expected session to be active")
	}

	session.Complete()

	if session.IsActive() {
		t.Error("expected session to not be active after completion")
	}
}

func TestSession_BelongsTo(t *testing.T) {
	session, _ := NewSession("thread-123", "user-123", "test-app")

	if !session.BelongsTo("user-123") {
		t.Error("expected session to belong to user-123")
	}
	if session.BelongsTo("other-user") {
		t.Error("expected session to not belong to other-user")
	}
}

func TestNewEvent(t *testing.T) {
	sessionID := "session-123"
	threadID := "thread-123"
	role := EventRoleAssistant
	content := "Hello, how can I help?"
	author := "ai-agent"

	event, err := NewEvent(sessionID, threadID, role, content, author)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if event.SessionID() != sessionID {
		t.Errorf("expected SessionID %s, got %s", sessionID, event.SessionID())
	}
	if event.ThreadID() != threadID {
		t.Errorf("expected ThreadID %s, got %s", threadID, event.ThreadID())
	}
	if event.Role() != role {
		t.Errorf("expected Role %s, got %s", role, event.Role())
	}
	if event.Content() != content {
		t.Errorf("expected Content %s, got %s", content, event.Content())
	}
	if event.Author() != author {
		t.Errorf("expected Author %s, got %s", author, event.Author())
	}
	if event.EventID() == "" {
		t.Error("expected EventID to be generated")
	}
	if event.InvocationID() == "" {
		t.Error("expected InvocationID to be generated")
	}
	if event.Timestamp().IsZero() {
		t.Error("expected Timestamp to be set")
	}
}

func TestNewEvent_Validation(t *testing.T) {
	tests := []struct {
		name        string
		sessionID   string
		role        EventRole
		content     string
		expectedErr error
	}{
		{
			name:        "empty session ID",
			sessionID:   "",
			role:        EventRoleUser,
			content:     "content",
			expectedErr: ErrInvalidSessionID,
		},
		{
			name:        "invalid role",
			sessionID:   "session-123",
			role:        EventRole("invalid"),
			content:     "content",
			expectedErr: ErrInvalidRole,
		},
		{
			name:        "empty content",
			sessionID:   "session-123",
			role:        EventRoleUser,
			content:     "",
			expectedErr: ErrEmptyContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewEvent(tt.sessionID, "thread-123", tt.role, tt.content, "author")
			if err != tt.expectedErr {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestReconstructEvent(t *testing.T) {
	eventID := "event-123"
	sessionID := "session-123"
	threadID := "thread-123"
	role := EventRoleSystem
	content := "System message"
	author := "system"
	invocationID := "inv-123"
	timestamp := time.Now()

	event := ReconstructEvent(eventID, sessionID, threadID, role, content, author, invocationID, timestamp)

	if event.EventID() != eventID {
		t.Errorf("expected EventID %s, got %s", eventID, event.EventID())
	}
	if event.SessionID() != sessionID {
		t.Errorf("expected SessionID %s, got %s", sessionID, event.SessionID())
	}
	if event.ThreadID() != threadID {
		t.Errorf("expected ThreadID %s, got %s", threadID, event.ThreadID())
	}
	if event.InvocationID() != invocationID {
		t.Errorf("expected InvocationID %s, got %s", invocationID, event.InvocationID())
	}
}

func TestEventRole_IsValid(t *testing.T) {
	tests := []struct {
		role     EventRole
		expected bool
	}{
		{EventRoleUser, true},
		{EventRoleAssistant, true},
		{EventRoleSystem, true},
		{EventRole("invalid"), false},
		{EventRole(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.role), func(t *testing.T) {
			if tt.role.IsValid() != tt.expected {
				t.Errorf("expected IsValid() = %v for role %s", tt.expected, tt.role)
			}
		})
	}
}

func TestEventRole_Constants(t *testing.T) {
	if EventRoleUser != "user" {
		t.Errorf("expected EventRoleUser 'user', got %s", EventRoleUser)
	}
	if EventRoleAssistant != "assistant" {
		t.Errorf("expected EventRoleAssistant 'assistant', got %s", EventRoleAssistant)
	}
	if EventRoleSystem != "system" {
		t.Errorf("expected EventRoleSystem 'system', got %s", EventRoleSystem)
	}
}

func TestSessionState_Constants(t *testing.T) {
	if SessionStateActive != "active" {
		t.Errorf("expected SessionStateActive 'active', got %s", SessionStateActive)
	}
	if SessionStateComplete != "complete" {
		t.Errorf("expected SessionStateComplete 'complete', got %s", SessionStateComplete)
	}
	if SessionStateError != "error" {
		t.Errorf("expected SessionStateError 'error', got %s", SessionStateError)
	}
}
