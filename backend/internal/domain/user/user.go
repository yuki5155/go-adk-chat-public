package user

import (
	"time"

	"github.com/yuki5155/go-google-auth/internal/domain/shared"
)

// User represents the User aggregate root
type User struct {
	id        UserID
	email     Email
	profile   Profile
	createdAt time.Time
	updatedAt time.Time
	events    []shared.DomainEvent
}

// NewUser creates a new User with validation
func NewUser(id UserID, email Email, profile Profile) (*User, error) {
	if id.IsEmpty() {
		return nil, shared.ErrEmptyUserID
	}

	if !email.IsVerified() {
		return nil, shared.ErrUnverifiedEmail
	}

	now := time.Now()
	user := &User{
		id:        id,
		email:     email,
		profile:   profile,
		createdAt: now,
		updatedAt: now,
		events:    make([]shared.DomainEvent, 0),
	}

	// Record domain event
	user.addEvent(NewUserRegisteredEvent(id.Value(), email.Value(), profile.Name()))

	return user, nil
}

// ReconstructUser reconstructs a User from persistence (without domain events)
func ReconstructUser(id UserID, email Email, profile Profile, createdAt, updatedAt time.Time) *User {
	return &User{
		id:        id,
		email:     email,
		profile:   profile,
		createdAt: createdAt,
		updatedAt: updatedAt,
		events:    make([]shared.DomainEvent, 0),
	}
}

// ID returns the user's ID
func (u *User) ID() UserID {
	return u.id
}

// Email returns the user's email
func (u *User) Email() Email {
	return u.email
}

// Profile returns the user's profile
func (u *User) Profile() Profile {
	return u.profile
}

// CreatedAt returns when the user was created
func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

// UpdatedAt returns when the user was last updated
func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

// UpdateProfile updates the user's profile
func (u *User) UpdateProfile(profile Profile) {
	u.profile = profile
	u.updatedAt = time.Now()
}

// UpdateEmail updates the user's email (must be verified)
func (u *User) UpdateEmail(email Email) error {
	if !email.IsVerified() {
		return shared.ErrUnverifiedEmail
	}

	u.email = email
	u.updatedAt = time.Now()
	return nil
}

// RecordLogin records a login event
func (u *User) RecordLogin() {
	u.addEvent(NewUserLoggedInEvent(u.id.Value(), u.email.Value()))
	u.updatedAt = time.Now()
}

// DomainEvents returns all domain events
func (u *User) DomainEvents() []shared.DomainEvent {
	return u.events
}

// ClearDomainEvents clears all domain events (after they've been published)
func (u *User) ClearDomainEvents() {
	u.events = make([]shared.DomainEvent, 0)
}

// addEvent adds a domain event
func (u *User) addEvent(event shared.DomainEvent) {
	u.events = append(u.events, event)
}
