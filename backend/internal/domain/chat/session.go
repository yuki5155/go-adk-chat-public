package chat

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// SessionState represents the state of a session
type SessionState string

const (
	SessionStateActive   SessionState = "active"
	SessionStateComplete SessionState = "complete"
	SessionStateError    SessionState = "error"
)

// EventRole represents the role of an event author
type EventRole string

const (
	EventRoleUser      EventRole = "user"
	EventRoleAssistant EventRole = "assistant"
	EventRoleSystem    EventRole = "system"
)

// IsValid checks if the role is valid
func (r EventRole) IsValid() bool {
	return r == EventRoleUser || r == EventRoleAssistant || r == EventRoleSystem
}

// Event represents a single message/event in a session
type Event struct {
	eventID      string
	sessionID    string
	threadID     string
	role         EventRole
	content      string
	author       string
	invocationID string
	timestamp    time.Time
}

// NewEvent creates a new Event
func NewEvent(sessionID, threadID string, role EventRole, content, author string) (*Event, error) {
	if sessionID == "" {
		return nil, ErrInvalidSessionID
	}
	if !role.IsValid() {
		return nil, ErrInvalidRole
	}
	if content == "" {
		return nil, ErrEmptyContent
	}

	now := time.Now()
	// Generate sortable event ID (timestamp prefix + uuid)
	eventID := fmt.Sprintf("%d_%s", now.UnixNano(), uuid.New().String())

	return &Event{
		eventID:      eventID,
		sessionID:    sessionID,
		threadID:     threadID,
		role:         role,
		content:      content,
		author:       author,
		invocationID: uuid.New().String(),
		timestamp:    now,
	}, nil
}

// ReconstructEvent reconstructs an Event from persistence
func ReconstructEvent(
	eventID, sessionID, threadID string,
	role EventRole,
	content, author, invocationID string,
	timestamp time.Time,
) *Event {
	return &Event{
		eventID:      eventID,
		sessionID:    sessionID,
		threadID:     threadID,
		role:         role,
		content:      content,
		author:       author,
		invocationID: invocationID,
		timestamp:    timestamp,
	}
}

// Getters for Event
func (e *Event) EventID() string      { return e.eventID }
func (e *Event) SessionID() string    { return e.sessionID }
func (e *Event) ThreadID() string     { return e.threadID }
func (e *Event) Role() EventRole      { return e.role }
func (e *Event) Content() string      { return e.content }
func (e *Event) Author() string       { return e.author }
func (e *Event) InvocationID() string { return e.invocationID }
func (e *Event) Timestamp() time.Time { return e.timestamp }

// Session represents a chat session entity
type Session struct {
	sessionID  string
	threadID   string
	userID     string
	appName    string
	state      SessionState
	eventCount int
	createdAt  time.Time
	updatedAt  time.Time
}

// NewSession creates a new Session
func NewSession(threadID, userID, appName string) (*Session, error) {
	if threadID == "" {
		return nil, ErrInvalidThreadID
	}
	if userID == "" {
		return nil, ErrInvalidUserID
	}

	now := time.Now()
	return &Session{
		sessionID:  uuid.New().String(),
		threadID:   threadID,
		userID:     userID,
		appName:    appName,
		state:      SessionStateActive,
		eventCount: 0,
		createdAt:  now,
		updatedAt:  now,
	}, nil
}

// ReconstructSession reconstructs a Session from persistence
func ReconstructSession(
	sessionID, threadID, userID, appName string,
	state SessionState,
	eventCount int,
	createdAt, updatedAt time.Time,
) *Session {
	return &Session{
		sessionID:  sessionID,
		threadID:   threadID,
		userID:     userID,
		appName:    appName,
		state:      state,
		eventCount: eventCount,
		createdAt:  createdAt,
		updatedAt:  updatedAt,
	}
}

// Getters for Session
func (s *Session) SessionID() string    { return s.sessionID }
func (s *Session) ThreadID() string     { return s.threadID }
func (s *Session) UserID() string       { return s.userID }
func (s *Session) AppName() string      { return s.appName }
func (s *Session) State() SessionState  { return s.state }
func (s *Session) EventCount() int      { return s.eventCount }
func (s *Session) CreatedAt() time.Time { return s.createdAt }
func (s *Session) UpdatedAt() time.Time { return s.updatedAt }

// IsActive returns whether the session is active
func (s *Session) IsActive() bool {
	return s.state == SessionStateActive
}

// IncrementEventCount increments the event count
func (s *Session) IncrementEventCount() {
	s.eventCount++
	s.updatedAt = time.Now()
}

// Complete marks the session as complete
func (s *Session) Complete() {
	s.state = SessionStateComplete
	s.updatedAt = time.Now()
}

// MarkError marks the session as having an error
func (s *Session) MarkError() {
	s.state = SessionStateError
	s.updatedAt = time.Now()
}

// BelongsTo checks if the session belongs to the given user
func (s *Session) BelongsTo(userID string) bool {
	return s.userID == userID
}
