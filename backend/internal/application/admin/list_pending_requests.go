package admin

import (
	"context"

	"github.com/yuki5155/go-google-auth/internal/domain/role"
)

// ListPendingRequestsUseCase handles listing of pending role requests
type ListPendingRequestsUseCase struct {
	roleRepo role.Repository
}

// NewListPendingRequestsUseCase creates a new ListPendingRequestsUseCase
func NewListPendingRequestsUseCase(roleRepo role.Repository) *ListPendingRequestsUseCase {
	return &ListPendingRequestsUseCase{
		roleRepo: roleRepo,
	}
}

// Execute retrieves all pending role requests
func (uc *ListPendingRequestsUseCase) Execute(ctx context.Context) (*RoleRequestListDTO, error) {
	requests, err := uc.roleRepo.ListRoleRequestsByStatus(ctx, role.RequestStatusPending)
	if err != nil {
		return nil, err
	}

	// Convert to DTOs
	dtos := make([]*RoleRequestDTO, 0, len(requests))
	for _, req := range requests {
		dtos = append(dtos, NewRoleRequestDTO(req))
	}

	return &RoleRequestListDTO{
		Requests: dtos,
		Count:    len(dtos),
	}, nil
}
