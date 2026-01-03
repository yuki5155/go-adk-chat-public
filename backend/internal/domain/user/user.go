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
	role      Role
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
		role:      RoleUser,
		createdAt: now,
		updatedAt: now,
		events:    make([]shared.DomainEvent, 0),
	}

	// Record domain event
	user.addEvent(NewUserRegisteredEvent(id.Value(), email.Value(), profile.Name()))

	return user, nil
}

// NewRootUser creates a root/admin user with verified email
func NewRootUser(id UserID, email Email, profile Profile) (*User, error) {
	if id.IsEmpty() {
		return nil, shared.ErrEmptyUserID
	}

	// Create verified email for root user
	verifiedEmail := email.Verify()

	now := time.Now()
	user := &User{
		id:        id,
		email:     verifiedEmail,
		profile:   profile,
		role:      RoleRoot,
		createdAt: now,
		updatedAt: now,
		events:    make([]shared.DomainEvent, 0),
	}

	// Record domain event
	user.addEvent(NewUserRegisteredEvent(id.Value(), verifiedEmail.Value(), profile.Name()))

	return user, nil
}

// ReconstructUser reconstructs a User from persistence (without domain events)
func ReconstructUser(id UserID, email Email, profile Profile, role Role, createdAt, updatedAt time.Time) *User {
	return &User{
		id:        id,
		email:     email,
		profile:   profile,
		role:      role,
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

// Role returns the user's role
func (u *User) Role() Role {
	return u.role
}

// IsRoot returns whether the user is a root/admin user
func (u *User) IsRoot() bool {
	return u.role.IsRoot()
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
