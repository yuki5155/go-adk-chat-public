package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProfile(t *testing.T) {
	profile := NewProfile("John Doe", "https://example.com/photo.jpg")
	assert.Equal(t, "John Doe", profile.Name())
	assert.Equal(t, "https://example.com/photo.jpg", profile.Picture())
	assert.False(t, profile.IsEmpty())
}

func TestProfile_EmptyProfile(t *testing.T) {
	profile := NewProfile("", "")
	assert.True(t, profile.IsEmpty())
	assert.Equal(t, "", profile.Name())
	assert.Equal(t, "", profile.Picture())
}

func TestProfile_WithName(t *testing.T) {
	original := NewProfile("John Doe", "https://example.com/photo.jpg")
	updated := original.WithName("Jane Smith")

	// Updated profile has new name
	assert.Equal(t, "Jane Smith", updated.Name())
	assert.Equal(t, "https://example.com/photo.jpg", updated.Picture())

	// Original is unchanged (immutability)
	assert.Equal(t, "John Doe", original.Name())
	assert.Equal(t, "https://example.com/photo.jpg", original.Picture())
}

func TestProfile_WithPicture(t *testing.T) {
	original := NewProfile("John Doe", "https://example.com/photo.jpg")
	updated := original.WithPicture("https://example.com/new-photo.jpg")

	// Updated profile has new picture
	assert.Equal(t, "John Doe", updated.Name())
	assert.Equal(t, "https://example.com/new-photo.jpg", updated.Picture())

	// Original is unchanged (immutability)
	assert.Equal(t, "John Doe", original.Name())
	assert.Equal(t, "https://example.com/photo.jpg", original.Picture())
}

func TestProfile_PartiallyEmpty(t *testing.T) {
	profileWithNameOnly := NewProfile("John Doe", "")
	assert.False(t, profileWithNameOnly.IsEmpty())

	profileWithPictureOnly := NewProfile("", "https://example.com/photo.jpg")
	assert.False(t, profileWithPictureOnly.IsEmpty())
}
