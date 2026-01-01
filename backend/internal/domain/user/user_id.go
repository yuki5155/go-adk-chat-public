package user

import (
	"strings"

	"github.com/yuki5155/go-google-auth/internal/domain/shared"
)

// UserID represents a unique identifier for a user
type UserID struct {
	value string
}

// NewUserID creates a new UserID with validation
func NewUserID(id string) (UserID, error) {
	if id == "" {
		return UserID{}, shared.ErrEmptyUserID
	}

	// Trim whitespace
	id = strings.TrimSpace(id)
	if id == "" {
		return UserID{}, shared.ErrEmptyUserID
	}

	return UserID{value: id}, nil
}

// Value returns the string value of the UserID
func (u UserID) Value() string {
	return u.value
}

// String implements the Stringer interface
func (u UserID) String() string {
	return u.value
}

// Equals compares two UserIDs for equality
func (u UserID) Equals(other UserID) bool {
	return u.value == other.value
}

// IsEmpty returns true if the UserID is empty
func (u UserID) IsEmpty() bool {
	return u.value == ""
}
