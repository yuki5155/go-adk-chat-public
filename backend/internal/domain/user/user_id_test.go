package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yuki5155/go-google-auth/internal/domain/shared"
)

func TestNewUserID_Success(t *testing.T) {
	id, err := NewUserID("user-123")
	assert.NoError(t, err)
	assert.Equal(t, "user-123", id.Value())
	assert.Equal(t, "user-123", id.String())
	assert.False(t, id.IsEmpty())
}

func TestNewUserID_EmptyString(t *testing.T) {
	id, err := NewUserID("")
	assert.Error(t, err)
	assert.Equal(t, shared.ErrEmptyUserID, err)
	assert.True(t, id.IsEmpty())
}

func TestNewUserID_WhitespaceOnly(t *testing.T) {
	id, err := NewUserID("   ")
	assert.Error(t, err)
	assert.Equal(t, shared.ErrEmptyUserID, err)
	assert.True(t, id.IsEmpty())
}

func TestNewUserID_TrimsWhitespace(t *testing.T) {
	id, err := NewUserID("  user-123  ")
	assert.NoError(t, err)
	assert.Equal(t, "user-123", id.Value())
}

func TestUserID_Equals(t *testing.T) {
	id1, _ := NewUserID("user-123")
	id2, _ := NewUserID("user-123")
	id3, _ := NewUserID("user-456")

	assert.True(t, id1.Equals(id2))
	assert.False(t, id1.Equals(id3))
}

func TestUserID_String(t *testing.T) {
	id, _ := NewUserID("google-user-123")
	assert.Equal(t, "google-user-123", id.String())
}
