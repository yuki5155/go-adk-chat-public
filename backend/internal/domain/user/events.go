package user

import "github.com/yuki5155/go-google-auth/internal/domain/shared"

// Event type constants
const (
	EventTypeUserRegistered = "user.registered"
	EventTypeUserLoggedIn   = "user.logged_in"
)

// UserRegisteredEvent is emitted when a new user is registered
type UserRegisteredEvent struct {
	shared.BaseDomainEvent
	UserID  string
	Email   string
	Name    string
}

// NewUserRegisteredEvent creates a new UserRegisteredEvent
func NewUserRegisteredEvent(userID, email, name string) UserRegisteredEvent {
	return UserRegisteredEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent(EventTypeUserRegistered, userID),
		UserID:          userID,
		Email:           email,
		Name:            name,
	}
}

// UserLoggedInEvent is emitted when a user logs in
type UserLoggedInEvent struct {
	shared.BaseDomainEvent
	UserID string
	Email  string
}

// NewUserLoggedInEvent creates a new UserLoggedInEvent
func NewUserLoggedInEvent(userID, email string) UserLoggedInEvent {
	return UserLoggedInEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent(EventTypeUserLoggedIn, userID),
		UserID:          userID,
		Email:           email,
	}
}
