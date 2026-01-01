package memory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuki5155/go-google-auth/internal/domain/shared"
	"github.com/yuki5155/go-google-auth/internal/domain/user"
)

func TestNewUserRepository(t *testing.T) {
	repo := NewUserRepository()

	assert.NotNil(t, repo)
	assert.NotNil(t, repo.users)
	assert.NotNil(t, repo.emails)
}

func TestUserRepository_Save_Success(t *testing.T) {
	ctx := context.Background()
	repo := NewUserRepository()

	userID, _ := user.NewUserID("test-user-123")
	email, _ := user.NewEmail("test@example.com", true)
	profile := user.NewProfile("Test User", "https://example.com/photo.jpg")
	u, _ := user.NewUser(userID, email, profile)

	err := repo.Save(ctx, u)

	require.NoError(t, err)
	assert.Len(t, repo.users, 1)
	assert.Len(t, repo.emails, 1)
}

func TestUserRepository_Save_DuplicateEmail(t *testing.T) {
	ctx := context.Background()
	repo := NewUserRepository()

	// Create first user
	userID1, _ := user.NewUserID("user-1")
	email, _ := user.NewEmail("test@example.com", true)
	profile := user.NewProfile("User 1", "")
	u1, _ := user.NewUser(userID1, email, profile)

	err := repo.Save(ctx, u1)
	require.NoError(t, err)

	// Try to create second user with same email
	userID2, _ := user.NewUserID("user-2")
	u2, _ := user.NewUser(userID2, email, profile)

	err = repo.Save(ctx, u2)

	assert.Error(t, err)
	assert.Equal(t, shared.ErrUserAlreadyExists, err)
}

func TestUserRepository_Save_UpdateExisting(t *testing.T) {
	ctx := context.Background()
	repo := NewUserRepository()

	userID, _ := user.NewUserID("test-user-123")
	email, _ := user.NewEmail("test@example.com", true)
	profile := user.NewProfile("Old Name", "")
	u, _ := user.NewUser(userID, email, profile)

	// Save initial user
	err := repo.Save(ctx, u)
	require.NoError(t, err)

	// Update profile
	newProfile := user.NewProfile("New Name", "https://example.com/new.jpg")
	u.UpdateProfile(newProfile)

	// Save updated user
	err = repo.Save(ctx, u)

	require.NoError(t, err)
	assert.Len(t, repo.users, 1) // Still only one user

	// Verify update
	savedUser, _ := repo.FindByID(ctx, userID)
	assert.Equal(t, "New Name", savedUser.Profile().Name())
}

func TestUserRepository_FindByID_Success(t *testing.T) {
	ctx := context.Background()
	repo := NewUserRepository()

	userID, _ := user.NewUserID("test-user-123")
	email, _ := user.NewEmail("test@example.com", true)
	profile := user.NewProfile("Test User", "https://example.com/photo.jpg")
	u, _ := user.NewUser(userID, email, profile)

	_ = repo.Save(ctx, u)

	foundUser, err := repo.FindByID(ctx, userID)

	require.NoError(t, err)
	assert.Equal(t, userID.Value(), foundUser.ID().Value())
	assert.Equal(t, email.Value(), foundUser.Email().Value())
	assert.Equal(t, "Test User", foundUser.Profile().Name())
}

func TestUserRepository_FindByID_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := NewUserRepository()

	userID, _ := user.NewUserID("nonexistent-user")

	foundUser, err := repo.FindByID(ctx, userID)

	assert.Error(t, err)
	assert.Equal(t, shared.ErrUserNotFound, err)
	assert.Nil(t, foundUser)
}

func TestUserRepository_FindByEmail_Success(t *testing.T) {
	ctx := context.Background()
	repo := NewUserRepository()

	userID, _ := user.NewUserID("test-user-123")
	email, _ := user.NewEmail("test@example.com", true)
	profile := user.NewProfile("Test User", "")
	u, _ := user.NewUser(userID, email, profile)

	_ = repo.Save(ctx, u)

	foundUser, err := repo.FindByEmail(ctx, email)

	require.NoError(t, err)
	assert.Equal(t, userID.Value(), foundUser.ID().Value())
	assert.Equal(t, email.Value(), foundUser.Email().Value())
}

func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := NewUserRepository()

	email, _ := user.NewEmail("nonexistent@example.com", true)

	foundUser, err := repo.FindByEmail(ctx, email)

	assert.Error(t, err)
	assert.Equal(t, shared.ErrUserNotFound, err)
	assert.Nil(t, foundUser)
}

func TestUserRepository_Delete_Success(t *testing.T) {
	ctx := context.Background()
	repo := NewUserRepository()

	userID, _ := user.NewUserID("test-user-123")
	email, _ := user.NewEmail("test@example.com", true)
	profile := user.NewProfile("Test User", "")
	u, _ := user.NewUser(userID, email, profile)

	_ = repo.Save(ctx, u)

	err := repo.Delete(ctx, userID)

	require.NoError(t, err)
	assert.Len(t, repo.users, 0)
	assert.Len(t, repo.emails, 0)

	// Verify user is gone
	_, err = repo.FindByID(ctx, userID)
	assert.Equal(t, shared.ErrUserNotFound, err)
}

func TestUserRepository_Delete_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := NewUserRepository()

	userID, _ := user.NewUserID("nonexistent-user")

	err := repo.Delete(ctx, userID)

	assert.Error(t, err)
	assert.Equal(t, shared.ErrUserNotFound, err)
}

func TestUserRepository_Exists_True(t *testing.T) {
	ctx := context.Background()
	repo := NewUserRepository()

	userID, _ := user.NewUserID("test-user-123")
	email, _ := user.NewEmail("test@example.com", true)
	profile := user.NewProfile("Test User", "")
	u, _ := user.NewUser(userID, email, profile)

	_ = repo.Save(ctx, u)

	exists, err := repo.Exists(ctx, userID)

	require.NoError(t, err)
	assert.True(t, exists)
}

func TestUserRepository_Exists_False(t *testing.T) {
	ctx := context.Background()
	repo := NewUserRepository()

	userID, _ := user.NewUserID("nonexistent-user")

	exists, err := repo.Exists(ctx, userID)

	require.NoError(t, err)
	assert.False(t, exists)
}

func TestUserRepository_ExistsByEmail_True(t *testing.T) {
	ctx := context.Background()
	repo := NewUserRepository()

	userID, _ := user.NewUserID("test-user-123")
	email, _ := user.NewEmail("test@example.com", true)
	profile := user.NewProfile("Test User", "")
	u, _ := user.NewUser(userID, email, profile)

	_ = repo.Save(ctx, u)

	exists, err := repo.ExistsByEmail(ctx, email)

	require.NoError(t, err)
	assert.True(t, exists)
}

func TestUserRepository_ExistsByEmail_False(t *testing.T) {
	ctx := context.Background()
	repo := NewUserRepository()

	email, _ := user.NewEmail("nonexistent@example.com", true)

	exists, err := repo.ExistsByEmail(ctx, email)

	require.NoError(t, err)
	assert.False(t, exists)
}

func TestUserRepository_ConcurrentAccess(t *testing.T) {
	ctx := context.Background()
	repo := NewUserRepository()

	// Create multiple users concurrently
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(index int) {
			userID, _ := user.NewUserID(string(rune('a' + index)))
			email, _ := user.NewEmail(string(rune('a'+index))+"@example.com", true)
			profile := user.NewProfile("User "+string(rune('a'+index)), "")
			u, _ := user.NewUser(userID, email, profile)

			_ = repo.Save(ctx, u)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all users were saved
	assert.Len(t, repo.users, 10)
	assert.Len(t, repo.emails, 10)
}
