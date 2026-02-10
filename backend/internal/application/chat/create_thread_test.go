package chat

import (
	"context"
	"testing"

	"github.com/yuki5155/go-google-auth/internal/domain/chat"
)

func TestNewCreateThreadUseCase(t *testing.T) {
	mockRepo := &MockThreadRepository{}
	uc := NewCreateThreadUseCase(mockRepo)

	if uc == nil {
		t.Fatal("expected use case to be created")
	}
	if uc.threadRepo == nil {
		t.Error("expected threadRepo to be set")
	}
}

func TestCreateThreadUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := &MockThreadRepository{}
		uc := NewCreateThreadUseCase(mockRepo)

		cmd := CreateThreadCommand{
			UserID: "user-123",
			Title:  "Test Thread",
		}

		dto, err := uc.Execute(ctx, cmd)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if dto == nil {
			t.Fatal("expected dto to be returned")
		}
		if dto.Title != "Test Thread" {
			t.Errorf("expected title 'Test Thread', got %s", dto.Title)
		}
		if dto.Status != "active" {
			t.Errorf("expected status 'active', got %s", dto.Status)
		}
		if !mockRepo.SaveCalled {
			t.Error("expected Save to be called")
		}
	})

	t.Run("success with empty title", func(t *testing.T) {
		mockRepo := &MockThreadRepository{}
		uc := NewCreateThreadUseCase(mockRepo)

		cmd := CreateThreadCommand{
			UserID: "user-123",
			Title:  "",
		}

		dto, err := uc.Execute(ctx, cmd)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if dto.Title != "New Conversation" {
			t.Errorf("expected default title 'New Conversation', got %s", dto.Title)
		}
	})

	t.Run("error on invalid user ID", func(t *testing.T) {
		mockRepo := &MockThreadRepository{}
		uc := NewCreateThreadUseCase(mockRepo)

		cmd := CreateThreadCommand{
			UserID: "",
			Title:  "Test",
		}

		_, err := uc.Execute(ctx, cmd)

		if err == nil {
			t.Error("expected error for empty user ID")
		}
	})

	t.Run("error on save failure", func(t *testing.T) {
		mockRepo := &MockThreadRepository{
			SaveFunc: func(ctx context.Context, thread *chat.Thread) error {
				return errMockSave
			},
		}
		uc := NewCreateThreadUseCase(mockRepo)

		cmd := CreateThreadCommand{
			UserID: "user-123",
			Title:  "Test",
		}

		_, err := uc.Execute(ctx, cmd)

		if err == nil {
			t.Error("expected error on save failure")
		}
	})
}
