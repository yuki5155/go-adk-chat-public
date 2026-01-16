package admin

import (
	"time"

	"github.com/yuki5155/go-google-auth/internal/domain/role"
)

// RoleRequestDTO represents a role request data transfer object
type RoleRequestDTO struct {
	RequestID     string    `json:"request_id"`
	UserID        string    `json:"user_id"`
	UserEmail     string    `json:"user_email"`
	RequestedRole string    `json:"requested_role"`
	Status        string    `json:"status"`
	RequestedAt   time.Time `json:"requested_at"`
	ProcessedAt   *time.Time `json:"processed_at,omitempty"`
	ProcessedBy   string    `json:"processed_by,omitempty"`
	Notes         string    `json:"notes,omitempty"`
}

// NewRoleRequestDTO creates a DTO from a domain RoleRequest
func NewRoleRequestDTO(req *role.RoleRequest) *RoleRequestDTO {
	return &RoleRequestDTO{
		RequestID:     req.RequestID(),
		UserID:        req.UserID(),
		UserEmail:     req.UserEmail(),
		RequestedRole: req.RequestedRole().String(),
		Status:        string(req.Status()),
		RequestedAt:   req.RequestedAt(),
		ProcessedAt:   req.ProcessedAt(),
		ProcessedBy:   req.ProcessedBy(),
		Notes:         req.Notes(),
	}
}

// UserRoleDTO represents a user role data transfer object
type UserRoleDTO struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewUserRoleDTO creates a DTO from a domain UserRole
func NewUserRoleDTO(userRole *role.UserRole) *UserRoleDTO {
	return &UserRoleDTO{
		UserID:    userRole.UserID(),
		Email:     userRole.Email(),
		Role:      userRole.Role().String(),
		Status:    string(userRole.Status()),
		CreatedAt: userRole.CreatedAt(),
		UpdatedAt: userRole.UpdatedAt(),
	}
}

// RoleRequestListDTO represents a list of role requests
type RoleRequestListDTO struct {
	Requests []*RoleRequestDTO `json:"requests"`
	Count    int               `json:"count"`
}

// UserRoleListDTO represents a list of user roles
type UserRoleListDTO struct {
	Users []*UserRoleDTO `json:"users"`
	Count int            `json:"count"`
}
