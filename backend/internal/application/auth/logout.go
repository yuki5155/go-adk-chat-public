package auth

import (
	"context"

	"github.com/yuki5155/go-google-auth/internal/application/dto"
)

// LogoutUseCase handles user logout operations
type LogoutUseCase struct {
	// No dependencies needed - logout is handled by clearing cookies in presentation layer
}

// NewLogoutUseCase creates a new LogoutUseCase
func NewLogoutUseCase() *LogoutUseCase {
	return &LogoutUseCase{}
}

// Execute performs logout operation
// Note: Actual token invalidation happens in the presentation layer by clearing cookies
// This use case can be extended to support token blacklisting or session invalidation
func (uc *LogoutUseCase) Execute(ctx context.Context) *dto.LogoutResponse {
	return &dto.LogoutResponse{
		Message: "Logged out successfully",
	}
}
