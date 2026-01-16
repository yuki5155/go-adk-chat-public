package admin

import (
	"context"

	"github.com/google/uuid"
	"github.com/yuki5155/go-google-auth/internal/domain/role"
	"github.com/yuki5155/go-google-auth/internal/domain/user"
)

// RequestRoleUseCase handles user role requests
type RequestRoleUseCase struct {
	roleRepo role.Repository
}

// NewRequestRoleUseCase creates a new RequestRoleUseCase
func NewRequestRoleUseCase(roleRepo role.Repository) *RequestRoleUseCase {
	return &RequestRoleUseCase{
		roleRepo: roleRepo,
	}
}

// Execute creates a new role request for a user
func (uc *RequestRoleUseCase) Execute(ctx context.Context, cmd RequestRoleCommand) (*RoleRequestDTO, error) {
	// Validate command
	if err := cmd.Validate(); err != nil {
		return nil, err
	}

	// Check if user already has a pending request
	existingRequest, err := uc.roleRepo.GetPendingRequestByUserID(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}

	if existingRequest != nil {
		return nil, role.ErrDuplicateRoleRequest
	}

	// Create new role request
	requestID := uuid.New().String()
	roleRequest, err := role.NewRoleRequest(requestID, cmd.UserID, cmd.UserEmail, cmd.RequestedRole)
	if err != nil {
		return nil, err
	}

	// Save to repository
	if err := uc.roleRepo.CreateRoleRequest(ctx, roleRequest); err != nil {
		return nil, err
	}

	return NewRoleRequestDTO(roleRequest), nil
}

// ExecuteLegacy is the old signature for backward compatibility
func (uc *RequestRoleUseCase) ExecuteLegacy(ctx context.Context, userID string, userEmail string, requestedRole user.Role) (*role.RoleRequest, error) {
	cmd := RequestRoleCommand{
		UserID:        userID,
		UserEmail:     userEmail,
		RequestedRole: requestedRole,
	}
	dto, err := uc.Execute(ctx, cmd)
	if err != nil {
		return nil, err
	}
	// Return domain object for backward compatibility
	return uc.roleRepo.GetRoleRequest(ctx, dto.RequestID)
}
