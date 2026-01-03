package admin

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/yuki5155/go-google-auth/internal/domain/shared"
	"github.com/yuki5155/go-google-auth/internal/domain/user"
)

// InitializeRootUserUseCase initializes the root user at application startup
type InitializeRootUserUseCase struct {
	userRepo  user.Repository
	rootEmail string
}

// NewInitializeRootUserUseCase creates a new use case for initializing the root user
func NewInitializeRootUserUseCase(userRepo user.Repository, rootEmail string) *InitializeRootUserUseCase {
	return &InitializeRootUserUseCase{
		userRepo:  userRepo,
		rootEmail: rootEmail,
	}
}

// Execute initializes the root user if it doesn't exist
func (uc *InitializeRootUserUseCase) Execute(ctx context.Context) error {
	if uc.rootEmail == "" {
		log.Println("ROOT_USER_EMAIL not configured, skipping root user initialization")
		return nil
	}

	// Create email value object (not verified yet, will be verified in NewRootUser)
	email, err := user.NewEmail(uc.rootEmail, false)
	if err != nil {
		return fmt.Errorf("invalid root email: %w", err)
	}

	// Check if user with root email already exists
	existingUser, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil && err != shared.ErrUserNotFound {
		return fmt.Errorf("failed to check if root user exists: %w", err)
	}

	if existingUser != nil {
		// User exists - check if they're already root
		if existingUser.IsRoot() {
			log.Printf("Root user with email %s already exists with root role", uc.rootEmail)
			return nil
		}
		// User exists but is not root - they need to re-login to get root role
		log.Printf("User with email %s exists but will need to re-login after backend restart to become root", uc.rootEmail)
		return nil
	}

	// Generate user ID from email hash
	hash := sha256.Sum256([]byte(uc.rootEmail))
	userIDStr := "root_" + hex.EncodeToString(hash[:])[:16]
	userID, err := user.NewUserID(userIDStr)
	if err != nil {
		return fmt.Errorf("failed to create user ID: %w", err)
	}

	// Create profile
	profile := user.NewProfile("Root Administrator", "")

	// Create root user
	rootUser, err := user.NewRootUser(userID, email, profile)
	if err != nil {
		return fmt.Errorf("failed to create root user: %w", err)
	}

	// Save root user
	if err := uc.userRepo.Save(ctx, rootUser); err != nil {
		return fmt.Errorf("failed to save root user: %w", err)
	}

	log.Printf("Root user initialized successfully with email: %s", uc.rootEmail)
	return nil
}
