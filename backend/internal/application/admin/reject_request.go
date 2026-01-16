package admin

import (
	"context"

	"github.com/yuki5155/go-google-auth/internal/domain/role"
)

// RejectRequestUseCase handles rejection of role requests
type RejectRequestUseCase struct {
	roleRepo role.Repository
}

// NewRejectRequestUseCase creates a new RejectRequestUseCase
func NewRejectRequestUseCase(roleRepo role.Repository) *RejectRequestUseCase {
	return &RejectRequestUseCase{
		roleRepo: roleRepo,
	}
}

// Execute rejects a role request
func (uc *RejectRequestUseCase) Execute(ctx context.Context, cmd RejectRequestCommand) error {
	// Validate command
	if err := cmd.Validate(); err != nil {
		return err
	}

	// Get the role request
	request, err := uc.roleRepo.GetRoleRequest(ctx, cmd.RequestID)
	if err != nil {
		return err
	}

	// Reject the request
	if err := request.Reject(cmd.RejectedBy, cmd.Notes); err != nil {
		return err
	}

	// Update the request in the repository
	if err := uc.roleRepo.UpdateRoleRequest(ctx, request); err != nil {
		return err
	}

	return nil
}
