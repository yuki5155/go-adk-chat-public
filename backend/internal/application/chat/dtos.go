package chat

import (
	"time"

	"github.com/yuki5155/go-google-auth/internal/domain/chat"
)

// ThreadDTO represents a thread for API responses
type ThreadDTO struct {
	ThreadID     string    `json:"thread_id"`
	Title        string    `json:"title"`
	Model        string    `json:"model"`
	Status       string    `json:"status"`
	MessageCount int       `json:"message_count"`
	LastMessage  string    `json:"last_message"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ThreadListDTO represents a paginated list of threads
type ThreadListDTO struct {
	Threads []ThreadDTO `json:"threads"`
	NextKey string      `json:"next_key,omitempty"`
}

// MessageDTO represents a message in a conversation
type MessageDTO struct {
	MessageID string    `json:"message_id"`
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// ThreadDetailDTO represents a thread with its messages
type ThreadDetailDTO struct {
	ThreadDTO
	Messages []MessageDTO `json:"messages"`
}

// SessionDTO represents a session for API responses
type SessionDTO struct {
	SessionID  string    `json:"session_id"`
	ThreadID   string    `json:"thread_id"`
	State      string    `json:"state"`
	EventCount int       `json:"event_count"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// SendMessageResponseDTO represents the response from sending a message
type SendMessageResponseDTO struct {
	Message  MessageDTO `json:"message"`
	Response MessageDTO `json:"response"`
}

// ModelDTO represents an available LLM model for API responses
type ModelDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ModelListDTO represents the list of available models
type ModelListDTO struct {
	Models  []ModelDTO `json:"models"`
	Default string     `json:"default"`
}

// ToThreadDTO converts a domain Thread to ThreadDTO
func ToThreadDTO(thread *chat.Thread) ThreadDTO {
	return ThreadDTO{
		ThreadID:     thread.ThreadID(),
		Title:        thread.Title(),
		Model:        thread.Model(),
		Status:       string(thread.Status()),
		MessageCount: thread.MessageCount(),
		LastMessage:  thread.LastMessage(),
		CreatedAt:    thread.CreatedAt(),
		UpdatedAt:    thread.UpdatedAt(),
	}
}

// ToThreadDTOs converts a slice of domain Threads to ThreadDTOs
func ToThreadDTOs(threads []*chat.Thread) []ThreadDTO {
	dtos := make([]ThreadDTO, 0, len(threads))
	for _, thread := range threads {
		dtos = append(dtos, ToThreadDTO(thread))
	}
	return dtos
}

// ToMessageDTO converts a domain Event to MessageDTO
func ToMessageDTO(event *chat.Event) MessageDTO {
	return MessageDTO{
		MessageID: event.EventID(),
		Role:      string(event.Role()),
		Content:   event.Content(),
		Timestamp: event.Timestamp(),
	}
}

// ToMessageDTOs converts a slice of domain Events to MessageDTOs
func ToMessageDTOs(events []*chat.Event) []MessageDTO {
	dtos := make([]MessageDTO, 0, len(events))
	for _, event := range events {
		dtos = append(dtos, ToMessageDTO(event))
	}
	return dtos
}
