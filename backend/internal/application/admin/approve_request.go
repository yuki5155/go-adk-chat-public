package admin

import (
	"context"

	"github.com/yuki5155/go-google-auth/internal/domain/role"
)

// ApproveRequestUseCase handles approval of role requests
type ApproveRequestUseCase struct {
	roleRepo role.Repository
}

// NewApproveRequestUseCase creates a new ApproveRequestUseCase
func NewApproveRequestUseCase(roleRepo role.Repository) *ApproveRequestUseCase {
	return &ApproveRequestUseCase{
		roleRepo: roleRepo,
	}
}

// Execute approves a role request and grants the role to the user
func (uc *ApproveRequestUseCase) Execute(ctx context.Context, cmd ApproveRequestCommand) error {
	// Validate command
	if err := cmd.Validate(); err != nil {
		return err
	}

	// Get the role request
	request, err := uc.roleRepo.GetRoleRequest(ctx, cmd.RequestID)
	if err != nil {
		return err
	}

	// Approve the request
	if err := request.Approve(cmd.ApprovedBy, cmd.Notes); err != nil {
		return err
	}

	// Update the request in the repository
	if err := uc.roleRepo.UpdateRoleRequest(ctx, request); err != nil {
		return err
	}

	// Create or update user role
	userRole, err := role.NewUserRole(
		request.UserID(),
		request.UserEmail(),
		request.RequestedRole(),
		cmd.ApprovedBy,
	)
	if err != nil {
		return err
	}

	// Save user role
	if err := uc.roleRepo.UpsertUserRole(ctx, userRole); err != nil {
		return err
	}

	return nil
}
