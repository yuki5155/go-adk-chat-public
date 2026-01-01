package user

import (
	"regexp"
	"strings"

	"github.com/yuki5155/go-google-auth/internal/domain/shared"
)

// Email represents a validated email address
type Email struct {
	value    string
	verified bool
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// NewEmail creates a new Email with validation
func NewEmail(address string, verified bool) (Email, error) {
	if address == "" {
		return Email{}, shared.ErrEmptyEmail
	}

	// Trim whitespace and convert to lowercase
	address = strings.ToLower(strings.TrimSpace(address))
	if address == "" {
		return Email{}, shared.ErrEmptyEmail
	}

	// Validate email format
	if !emailRegex.MatchString(address) {
		return Email{}, shared.ErrInvalidEmail
	}

	return Email{
		value:    address,
		verified: verified,
	}, nil
}

// Value returns the string value of the email
func (e Email) Value() string {
	return e.value
}

// String implements the Stringer interface
func (e Email) String() string {
	return e.value
}

// IsVerified returns whether the email has been verified
func (e Email) IsVerified() bool {
	return e.verified
}

// Equals compares two Emails for equality
func (e Email) Equals(other Email) bool {
	return e.value == other.value
}

// Verify marks the email as verified
func (e Email) Verify() Email {
	return Email{
		value:    e.value,
		verified: true,
	}
}
