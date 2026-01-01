package user

import "context"

// Repository defines the interface for user persistence
type Repository interface {
	// Save persists a user
	Save(ctx context.Context, user *User) error

	// FindByID retrieves a user by their ID
	FindByID(ctx context.Context, id UserID) (*User, error)

	// FindByEmail retrieves a user by their email
	FindByEmail(ctx context.Context, email Email) (*User, error)

	// Delete removes a user
	Delete(ctx context.Context, id UserID) error

	// Exists checks if a user exists by ID
	Exists(ctx context.Context, id UserID) (bool, error)

	// ExistsByEmail checks if a user exists by email
	ExistsByEmail(ctx context.Context, email Email) (bool, error)
}
