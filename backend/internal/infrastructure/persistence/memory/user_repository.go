package memory

import (
	"context"
	"sync"

	"github.com/yuki5155/go-google-auth/internal/domain/shared"
	"github.com/yuki5155/go-google-auth/internal/domain/user"
)

// UserRepository is an in-memory implementation of user.Repository
type UserRepository struct {
	mu    sync.RWMutex
	users map[string]*user.User // key: user ID
	emails map[string]string    // key: email, value: user ID
}

// NewUserRepository creates a new in-memory user repository
func NewUserRepository() *UserRepository {
	return &UserRepository{
		users:  make(map[string]*user.User),
		emails: make(map[string]string),
	}
}

// Save persists a user to the in-memory store
func (r *UserRepository) Save(ctx context.Context, u *user.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	userID := u.ID().Value()
	email := u.Email().Value()

	// Check if email is already taken by another user
	if existingUserID, exists := r.emails[email]; exists && existingUserID != userID {
		return shared.ErrUserAlreadyExists
	}

	r.users[userID] = u
	r.emails[email] = userID

	return nil
}

// FindByID retrieves a user by their ID
func (r *UserRepository) FindByID(ctx context.Context, id user.UserID) (*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, exists := r.users[id.Value()]
	if !exists {
		return nil, shared.ErrUserNotFound
	}

	return u, nil
}

// FindByEmail retrieves a user by their email
func (r *UserRepository) FindByEmail(ctx context.Context, email user.Email) (*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	userID, exists := r.emails[email.Value()]
	if !exists {
		return nil, shared.ErrUserNotFound
	}

	u, exists := r.users[userID]
	if !exists {
		return nil, shared.ErrUserNotFound
	}

	return u, nil
}

// Delete removes a user from the repository
func (r *UserRepository) Delete(ctx context.Context, id user.UserID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	userID := id.Value()
	u, exists := r.users[userID]
	if !exists {
		return shared.ErrUserNotFound
	}

	// Remove from email index
	delete(r.emails, u.Email().Value())
	delete(r.users, userID)

	return nil
}

// Exists checks if a user exists by ID
func (r *UserRepository) Exists(ctx context.Context, id user.UserID) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.users[id.Value()]
	return exists, nil
}

// ExistsByEmail checks if a user exists by email
func (r *UserRepository) ExistsByEmail(ctx context.Context, email user.Email) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.emails[email.Value()]
	return exists, nil
}
