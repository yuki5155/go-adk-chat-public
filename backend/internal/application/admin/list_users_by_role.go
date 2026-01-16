package admin

import (
	"context"

	"github.com/yuki5155/go-google-auth/internal/domain/role"
)

// ListUsersByRoleUseCase handles listing users by their role
type ListUsersByRoleUseCase struct {
	roleRepo role.Repository
}

// NewListUsersByRoleUseCase creates a new ListUsersByRoleUseCase
func NewListUsersByRoleUseCase(roleRepo role.Repository) *ListUsersByRoleUseCase {
	return &ListUsersByRoleUseCase{
		roleRepo: roleRepo,
	}
}

// Execute retrieves all users with a specific role
func (uc *ListUsersByRoleUseCase) Execute(ctx context.Context, query ListUsersByRoleQuery) (*UserRoleListDTO, error) {
	// Validate query
	if err := query.Validate(); err != nil {
		return nil, err
	}

	userRoles, err := uc.roleRepo.ListUsersByRole(ctx, query.Role)
	if err != nil {
		return nil, err
	}

	// Convert to DTOs
	dtos := make([]*UserRoleDTO, 0, len(userRoles))
	for _, userRole := range userRoles {
		dtos = append(dtos, NewUserRoleDTO(userRole))
	}

	return &UserRoleListDTO{
		Users: dtos,
		Count: len(dtos),
	}, nil
}
