package chat

import (
	"context"
	"fmt"

	"github.com/yuki5155/go-google-auth/internal/application/ports"
	"github.com/yuki5155/go-google-auth/internal/domain/chat"
)

// SendMessageUseCase handles sending messages and getting AI responses
type SendMessageUseCase struct {
	threadRepo  chat.ThreadRepository
	sessionRepo chat.SessionRepository
	eventRepo   chat.EventRepository
	aiRunner    ports.AIRunner
}

// NewSendMessageUseCase creates a new SendMessageUseCase
func NewSendMessageUseCase(
	threadRepo chat.ThreadRepository,
	sessionRepo chat.SessionRepository,
	eventRepo chat.EventRepository,
	aiRunner ports.AIRunner,
) *SendMessageUseCase {
	return &SendMessageUseCase{
		threadRepo:  threadRepo,
		sessionRepo: sessionRepo,
		eventRepo:   eventRepo,
		aiRunner:    aiRunner,
	}
}

// messageContext holds the context needed for message processing
type messageContext struct {
	thread    *chat.Thread
	session   *chat.Session
	userEvent *chat.Event
	history   []ports.ChatMessage
}

// prepareMessage handles common setup for both Execute and ExecuteStream
func (uc *SendMessageUseCase) prepareMessage(ctx context.Context, cmd SendMessageCommand) (*messageContext, error) {
	// Check if AI runner is configured
	if uc.aiRunner == nil {
		return nil, fmt.Errorf("AI runner is not configured - check GEMINI_API_KEY")
	}

	// Fetch thread
	thread, err := uc.threadRepo.FindByID(ctx, cmd.UserID, cmd.ThreadID)
	if err != nil {
		return nil, fmt.Errorf("failed to get thread: %w", err)
	}

	// Verify ownership
	if !thread.BelongsTo(cmd.UserID) {
		return nil, chat.ErrThreadUnauthorized
	}

	// Get or create active session
	session, err := uc.getOrCreateSession(ctx, cmd.ThreadID, cmd.UserID)
	if err != nil {
		return nil, err
	}

	// Create and save user message event
	userEvent, err := chat.NewEvent(session.SessionID(), cmd.ThreadID, chat.EventRoleUser, cmd.Content, cmd.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to create user event: %w", err)
	}

	if err := uc.eventRepo.Save(ctx, userEvent); err != nil {
		return nil, fmt.Errorf("failed to save user event: %w", err)
	}

	// Build conversation history
	history, err := uc.buildHistory(ctx, session.SessionID(), userEvent.EventID())
	if err != nil {
		return nil, err
	}

	return &messageContext{
		thread:    thread,
		session:   session,
		userEvent: userEvent,
		history:   history,
	}, nil
}

// getOrCreateSession gets an existing active session or creates a new one
func (uc *SendMessageUseCase) getOrCreateSession(ctx context.Context, threadID, userID string) (*chat.Session, error) {
	session, err := uc.sessionRepo.FindActiveByThread(ctx, threadID)
	if err != nil {
		return nil, fmt.Errorf("failed to find session: %w", err)
	}

	if session == nil {
		session, err = chat.NewSession(threadID, userID, uc.aiRunner.Config().AppName)
		if err != nil {
			return nil, fmt.Errorf("failed to create session: %w", err)
		}
		if err := uc.sessionRepo.Save(ctx, session); err != nil {
			return nil, fmt.Errorf("failed to save session: %w", err)
		}
	}

	return session, nil
}

// buildHistory builds the conversation history for AI context
func (uc *SendMessageUseCase) buildHistory(ctx context.Context, sessionID, excludeEventID string) ([]ports.ChatMessage, error) {
	events, err := uc.eventRepo.ListBySession(ctx, sessionID, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}

	history := make([]ports.ChatMessage, 0, len(events))
	for _, event := range events {
		if event.EventID() == excludeEventID {
			continue
		}
		history = append(history, ports.ChatMessage{
			Role:    string(event.Role()),
			Content: event.Content(),
		})
	}

	return history, nil
}

// finalizeMessage saves the assistant response and updates thread/session
func (uc *SendMessageUseCase) finalizeMessage(ctx context.Context, mc *messageContext, response string, userContent string) (*chat.Event, error) {
	// Create assistant message event
	assistantEvent, err := chat.NewEvent(mc.session.SessionID(), mc.thread.ThreadID(), chat.EventRoleAssistant, response, "assistant")
	if err != nil {
		return nil, fmt.Errorf("failed to create assistant event: %w", err)
	}

	// Save assistant event
	if err := uc.eventRepo.Save(ctx, assistantEvent); err != nil {
		return nil, fmt.Errorf("failed to save assistant event: %w", err)
	}

	// Update thread
	mc.thread.UpdateLastMessage(response)
	if mc.thread.MessageCount() == 1 {
		mc.thread.GenerateTitle(userContent)
	}
	if err := uc.threadRepo.Save(ctx, mc.thread); err != nil {
		return nil, fmt.Errorf("failed to update thread: %w", err)
	}

	// Update session event count (for both user and assistant)
	mc.session.IncrementEventCount()
	mc.session.IncrementEventCount()
	if err := uc.sessionRepo.Save(ctx, mc.session); err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	return assistantEvent, nil
}

// Execute sends a message and returns the AI response
func (uc *SendMessageUseCase) Execute(ctx context.Context, cmd SendMessageCommand) (*SendMessageResponseDTO, error) {
	mc, err := uc.prepareMessage(ctx, cmd)
	if err != nil {
		return nil, err
	}

	// Send to AI using the thread's model
	response, err := uc.aiRunner.SendMessage(ctx, mc.history, cmd.Content, mc.thread.Model())
	if err != nil {
		return nil, fmt.Errorf("failed to get AI response: %w", err)
	}

	// Finalize and save response
	assistantEvent, err := uc.finalizeMessage(ctx, mc, response, cmd.Content)
	if err != nil {
		return nil, err
	}

	return &SendMessageResponseDTO{
		Message:  ToMessageDTO(mc.userEvent),
		Response: ToMessageDTO(assistantEvent),
	}, nil
}

// StreamCallback is the callback function for streaming responses
type StreamCallback func(event ports.StreamEvent) error

// ExecuteStream sends a message and streams the AI response
func (uc *SendMessageUseCase) ExecuteStream(ctx context.Context, cmd SendMessageCommand, callback StreamCallback) (*SendMessageResponseDTO, error) {
	mc, err := uc.prepareMessage(ctx, cmd)
	if err != nil {
		return nil, err
	}

	// Stream from AI using the thread's model with tool support
	var fullResponse string
	err = uc.aiRunner.StreamMessageWithTools(ctx, mc.history, cmd.Content, mc.thread.Model(), func(event ports.StreamEvent) error {
		if event.Type == ports.StreamEventChunk {
			fullResponse += event.Content
		}
		return callback(event)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to stream AI response: %w", err)
	}

	// Finalize and save response
	assistantEvent, err := uc.finalizeMessage(ctx, mc, fullResponse, cmd.Content)
	if err != nil {
		return nil, err
	}

	return &SendMessageResponseDTO{
		Message:  ToMessageDTO(mc.userEvent),
		Response: ToMessageDTO(assistantEvent),
	}, nil
}
