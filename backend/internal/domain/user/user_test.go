package user

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuki5155/go-google-auth/internal/domain/shared"
)

func TestNewUser_Success(t *testing.T) {
	userID, _ := NewUserID("google-user-123")
	email, _ := NewEmail("test@example.com", true)
	profile := NewProfile("Test User", "https://example.com/photo.jpg")

	user, err := NewUser(userID, email, profile)

	require.NoError(t, err)
	assert.Equal(t, userID, user.ID())
	assert.Equal(t, email, user.Email())
	assert.Equal(t, profile, user.Profile())
	assert.False(t, user.CreatedAt().IsZero())
	assert.False(t, user.UpdatedAt().IsZero())
	assert.Equal(t, user.CreatedAt(), user.UpdatedAt())

	// Check domain event was recorded
	events := user.DomainEvents()
	require.Len(t, events, 1)
	assert.Equal(t, EventTypeUserRegistered, events[0].EventType())
}

func TestNewUser_EmptyUserID(t *testing.T) {
	emptyID := UserID{}
	email, _ := NewEmail("test@example.com", true)
	profile := NewProfile("Test User", "https://example.com/photo.jpg")

	user, err := NewUser(emptyID, email, profile)

	assert.Error(t, err)
	assert.Equal(t, shared.ErrEmptyUserID, err)
	assert.Nil(t, user)
}

func TestNewUser_UnverifiedEmail(t *testing.T) {
	userID, _ := NewUserID("google-user-123")
	unverifiedEmail, _ := NewEmail("test@example.com", false)
	profile := NewProfile("Test User", "https://example.com/photo.jpg")

	user, err := NewUser(userID, unverifiedEmail, profile)

	assert.Error(t, err)
	assert.Equal(t, shared.ErrUnverifiedEmail, err)
	assert.Nil(t, user)
}

func TestReconstructUser(t *testing.T) {
	userID, _ := NewUserID("google-user-123")
	email, _ := NewEmail("test@example.com", true)
	profile := NewProfile("Test User", "https://example.com/photo.jpg")
	createdAt := time.Now().Add(-24 * time.Hour)
	updatedAt := time.Now()

	user := ReconstructUser(userID, email, profile, createdAt, updatedAt)

	assert.Equal(t, userID, user.ID())
	assert.Equal(t, email, user.Email())
	assert.Equal(t, profile, user.Profile())
	assert.Equal(t, createdAt, user.CreatedAt())
	assert.Equal(t, updatedAt, user.UpdatedAt())

	// Reconstructed users should not have events
	assert.Empty(t, user.DomainEvents())
}

func TestUser_UpdateProfile(t *testing.T) {
	userID, _ := NewUserID("google-user-123")
	email, _ := NewEmail("test@example.com", true)
	originalProfile := NewProfile("Old Name", "https://example.com/old.jpg")

	user, _ := NewUser(userID, email, originalProfile)
	originalUpdatedAt := user.UpdatedAt()

	// Wait a moment to ensure timestamp changes
	time.Sleep(time.Millisecond)

	newProfile := NewProfile("New Name", "https://example.com/new.jpg")
	user.UpdateProfile(newProfile)

	assert.Equal(t, newProfile, user.Profile())
	assert.True(t, user.UpdatedAt().After(originalUpdatedAt))
}

func TestUser_UpdateEmail_Success(t *testing.T) {
	userID, _ := NewUserID("google-user-123")
	originalEmail, _ := NewEmail("old@example.com", true)
	profile := NewProfile("Test User", "https://example.com/photo.jpg")

	user, _ := NewUser(userID, originalEmail, profile)
	originalUpdatedAt := user.UpdatedAt()

	// Wait a moment to ensure timestamp changes
	time.Sleep(time.Millisecond)

	newEmail, _ := NewEmail("new@example.com", true)
	err := user.UpdateEmail(newEmail)

	assert.NoError(t, err)
	assert.Equal(t, newEmail, user.Email())
	assert.True(t, user.UpdatedAt().After(originalUpdatedAt))
}

func TestUser_UpdateEmail_Unverified(t *testing.T) {
	userID, _ := NewUserID("google-user-123")
	email, _ := NewEmail("test@example.com", true)
	profile := NewProfile("Test User", "https://example.com/photo.jpg")

	user, _ := NewUser(userID, email, profile)

	unverifiedEmail, _ := NewEmail("unverified@example.com", false)
	err := user.UpdateEmail(unverifiedEmail)

	assert.Error(t, err)
	assert.Equal(t, shared.ErrUnverifiedEmail, err)
	// Email should remain unchanged
	assert.Equal(t, email, user.Email())
}

func TestUser_RecordLogin(t *testing.T) {
	userID, _ := NewUserID("google-user-123")
	email, _ := NewEmail("test@example.com", true)
	profile := NewProfile("Test User", "https://example.com/photo.jpg")

	user, _ := NewUser(userID, email, profile)
	user.ClearDomainEvents() // Clear registration event

	originalUpdatedAt := user.UpdatedAt()
	time.Sleep(time.Millisecond)

	user.RecordLogin()

	assert.True(t, user.UpdatedAt().After(originalUpdatedAt))

	// Check login event was recorded
	events := user.DomainEvents()
	require.Len(t, events, 1)
	assert.Equal(t, EventTypeUserLoggedIn, events[0].EventType())

	loggedInEvent, ok := events[0].(UserLoggedInEvent)
	require.True(t, ok)
	assert.Equal(t, userID.Value(), loggedInEvent.UserID)
	assert.Equal(t, email.Value(), loggedInEvent.Email)
}

func TestUser_DomainEvents(t *testing.T) {
	userID, _ := NewUserID("google-user-123")
	email, _ := NewEmail("test@example.com", true)
	profile := NewProfile("Test User", "https://example.com/photo.jpg")

	user, _ := NewUser(userID, email, profile)

	// Should have registration event
	events := user.DomainEvents()
	require.Len(t, events, 1)
	assert.Equal(t, EventTypeUserRegistered, events[0].EventType())

	// Add login event
	user.RecordLogin()
	events = user.DomainEvents()
	require.Len(t, events, 2)
	assert.Equal(t, EventTypeUserRegistered, events[0].EventType())
	assert.Equal(t, EventTypeUserLoggedIn, events[1].EventType())

	// Clear events
	user.ClearDomainEvents()
	events = user.DomainEvents()
	assert.Empty(t, events)
}

func TestUserRegisteredEvent(t *testing.T) {
	event := NewUserRegisteredEvent("user-123", "test@example.com", "Test User")

	assert.Equal(t, EventTypeUserRegistered, event.EventType())
	assert.Equal(t, "user-123", event.AggregateID())
	assert.Equal(t, "user-123", event.UserID)
	assert.Equal(t, "test@example.com", event.Email)
	assert.Equal(t, "Test User", event.Name)
	assert.False(t, event.OccurredAt().IsZero())
}

func TestUserLoggedInEvent(t *testing.T) {
	event := NewUserLoggedInEvent("user-123", "test@example.com")

	assert.Equal(t, EventTypeUserLoggedIn, event.EventType())
	assert.Equal(t, "user-123", event.AggregateID())
	assert.Equal(t, "user-123", event.UserID)
	assert.Equal(t, "test@example.com", event.Email)
	assert.False(t, event.OccurredAt().IsZero())
}
