package role

import (
	"context"

	"github.com/yuki5155/go-google-auth/internal/domain/user"
)

// Repository defines the interface for role persistence
type Repository interface {
	// UserRole operations
	GetUserRole(ctx context.Context, userID string) (*UserRole, error)
	GetUserRoleByEmail(ctx context.Context, email string) (*UserRole, error)
	UpsertUserRole(ctx context.Context, userRole *UserRole) error
	ListUsersByRole(ctx context.Context, role user.Role) ([]*UserRole, error)

	// RoleRequest operations
	CreateRoleRequest(ctx context.Context, request *RoleRequest) error
	GetRoleRequest(ctx context.Context, requestID string) (*RoleRequest, error)
	GetPendingRequestByUserID(ctx context.Context, userID string) (*RoleRequest, error)
	UpdateRoleRequest(ctx context.Context, request *RoleRequest) error
	ListRoleRequestsByStatus(ctx context.Context, status RequestStatus) ([]*RoleRequest, error)
	ListRoleRequestsByUserID(ctx context.Context, userID string) ([]*RoleRequest, error)
}
