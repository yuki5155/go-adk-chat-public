package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yuki5155/go-google-auth/internal/domain/shared"
)

func TestNewEmail_Success(t *testing.T) {
	tests := []struct {
		name     string
		address  string
		verified bool
		expected string
	}{
		{
			name:     "valid email verified",
			address:  "test@example.com",
			verified: true,
			expected: "test@example.com",
		},
		{
			name:     "valid email unverified",
			address:  "user@domain.org",
			verified: false,
			expected: "user@domain.org",
		},
		{
			name:     "email with plus sign",
			address:  "user+tag@example.com",
			verified: true,
			expected: "user+tag@example.com",
		},
		{
			name:     "email with numbers",
			address:  "user123@example.com",
			verified: true,
			expected: "user123@example.com",
		},
		{
			name:     "uppercase converted to lowercase",
			address:  "USER@EXAMPLE.COM",
			verified: true,
			expected: "user@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, err := NewEmail(tt.address, tt.verified)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, email.Value())
			assert.Equal(t, tt.expected, email.String())
			assert.Equal(t, tt.verified, email.IsVerified())
		})
	}
}

func TestNewEmail_Invalid(t *testing.T) {
	tests := []struct {
		name    string
		address string
		errType error
	}{
		{
			name:    "empty string",
			address: "",
			errType: shared.ErrEmptyEmail,
		},
		{
			name:    "whitespace only",
			address: "   ",
			errType: shared.ErrEmptyEmail,
		},
		{
			name:    "missing @",
			address: "userexample.com",
			errType: shared.ErrInvalidEmail,
		},
		{
			name:    "missing domain",
			address: "user@",
			errType: shared.ErrInvalidEmail,
		},
		{
			name:    "missing local part",
			address: "@example.com",
			errType: shared.ErrInvalidEmail,
		},
		{
			name:    "no TLD",
			address: "user@example",
			errType: shared.ErrInvalidEmail,
		},
		{
			name:    "invalid characters",
			address: "user name@example.com",
			errType: shared.ErrInvalidEmail,
		},
		{
			name:    "multiple @",
			address: "user@@example.com",
			errType: shared.ErrInvalidEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, err := NewEmail(tt.address, false)
			assert.Error(t, err)
			assert.Equal(t, tt.errType, err)
			assert.Equal(t, "", email.Value())
		})
	}
}

func TestEmail_Verify(t *testing.T) {
	email, _ := NewEmail("test@example.com", false)
	assert.False(t, email.IsVerified())

	verified := email.Verify()
	assert.True(t, verified.IsVerified())
	assert.Equal(t, email.Value(), verified.Value())

	// Original should be unchanged (immutability)
	assert.False(t, email.IsVerified())
}

func TestEmail_Equals(t *testing.T) {
	email1, _ := NewEmail("test@example.com", true)
	email2, _ := NewEmail("test@example.com", false)
	email3, _ := NewEmail("other@example.com", true)

	// Same email address, different verification status
	assert.True(t, email1.Equals(email2))

	// Different email addresses
	assert.False(t, email1.Equals(email3))
}

func TestEmail_TrimsAndLowercase(t *testing.T) {
	email, err := NewEmail("  TEST@EXAMPLE.COM  ", true)
	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", email.Value())
}
