package admin

import (
	"context"

	"github.com/yuki5155/go-google-auth/internal/domain/role"
)

// CheckUserRoleUseCase handles checking a user's current role
type CheckUserRoleUseCase struct {
	roleRepo role.Repository
}

// NewCheckUserRoleUseCase creates a new CheckUserRoleUseCase
func NewCheckUserRoleUseCase(roleRepo role.Repository) *CheckUserRoleUseCase {
	return &CheckUserRoleUseCase{
		roleRepo: roleRepo,
	}
}

// Execute retrieves a user's current role
func (uc *CheckUserRoleUseCase) Execute(ctx context.Context, userID string) (*role.UserRole, error) {
	userRole, err := uc.roleRepo.GetUserRole(ctx, userID)
	if err != nil {
		// If user role not found, user has default role (user)
		if err == role.ErrUserRoleNotFound {
			return nil, nil
		}
		return nil, err
	}

	return userRole, nil
}
