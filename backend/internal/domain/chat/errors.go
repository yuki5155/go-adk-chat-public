package chat

import "errors"

// Domain errors for chat
var (
	// Thread errors
	ErrThreadNotFound    = errors.New("thread not found")
	ErrThreadUnauthorized = errors.New("unauthorized access to thread")
	ErrInvalidThreadID   = errors.New("invalid thread ID")
	ErrInvalidUserID     = errors.New("invalid user ID")
	ErrEmptyTitle        = errors.New("thread title cannot be empty")

	// Session errors
	ErrSessionNotFound    = errors.New("session not found")
	ErrInvalidSessionID   = errors.New("invalid session ID")
	ErrSessionExpired     = errors.New("session has expired")

	// Event errors
	ErrInvalidEventID    = errors.New("invalid event ID")
	ErrEmptyContent      = errors.New("message content cannot be empty")
	ErrInvalidRole       = errors.New("invalid message role")

	// Memory errors
	ErrMemoryNotFound    = errors.New("memory not found")
	ErrInvalidMemoryID   = errors.New("invalid memory ID")
)
