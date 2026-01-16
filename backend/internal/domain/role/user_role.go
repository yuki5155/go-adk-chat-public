package role

import (
	"time"

	"github.com/yuki5155/go-google-auth/internal/domain/user"
)

// UserRole represents a user's role assignment in the system
// This entity tracks role assignments stored in DynamoDB (excluding root user)
type UserRole struct {
	userID    string
	email     string
	role      user.Role
	status    Status
	grantedAt time.Time
	grantedBy string
	createdAt time.Time
	updatedAt time.Time
}

// Status represents the status of a user role
type Status string

const (
	// StatusActive represents an active user role
	StatusActive Status = "active"
	// StatusSuspended represents a suspended user role
	StatusSuspended Status = "suspended"
)

// NewUserRole creates a new UserRole entity
func NewUserRole(
	userID string,
	email string,
	role user.Role,
	grantedBy string,
) (*UserRole, error) {
	if userID == "" {
		return nil, ErrInvalidUserID
	}

	if email == "" {
		return nil, ErrInvalidEmail
	}

	if !role.IsValid() {
		return nil, ErrInvalidRole
	}

	now := time.Now()

	return &UserRole{
		userID:    userID,
		email:     email,
		role:      role,
		status:    StatusActive,
		grantedAt: now,
		grantedBy: grantedBy,
		createdAt: now,
		updatedAt: now,
	}, nil
}

// ReconstructUserRole reconstructs a UserRole from persistence
func ReconstructUserRole(
	userID string,
	email string,
	role user.Role,
	status Status,
	grantedAt time.Time,
	grantedBy string,
	createdAt time.Time,
	updatedAt time.Time,
) *UserRole {
	return &UserRole{
		userID:    userID,
		email:     email,
		role:      role,
		status:    status,
		grantedAt: grantedAt,
		grantedBy: grantedBy,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

// UserID returns the user ID
func (ur *UserRole) UserID() string {
	return ur.userID
}

// Email returns the user's email
func (ur *UserRole) Email() string {
	return ur.email
}

// Role returns the user's role
func (ur *UserRole) Role() user.Role {
	return ur.role
}

// Status returns the user role status
func (ur *UserRole) Status() Status {
	return ur.status
}

// GrantedAt returns when the role was granted
func (ur *UserRole) GrantedAt() time.Time {
	return ur.grantedAt
}

// GrantedBy returns who granted the role
func (ur *UserRole) GrantedBy() string {
	return ur.grantedBy
}

// CreatedAt returns when the role was created
func (ur *UserRole) CreatedAt() time.Time {
	return ur.createdAt
}

// UpdatedAt returns when the role was last updated
func (ur *UserRole) UpdatedAt() time.Time {
	return ur.updatedAt
}

// IsActive returns whether the role is active
func (ur *UserRole) IsActive() bool {
	return ur.status == StatusActive
}

// Suspend suspends the user role
func (ur *UserRole) Suspend() {
	ur.status = StatusSuspended
	ur.updatedAt = time.Now()
}

// Activate activates the user role
func (ur *UserRole) Activate() {
	ur.status = StatusActive
	ur.updatedAt = time.Now()
}

// ChangeRole changes the user's role
func (ur *UserRole) ChangeRole(newRole user.Role, changedBy string) error {
	if !newRole.IsValid() {
		return ErrInvalidRole
	}

	ur.role = newRole
	ur.grantedBy = changedBy
	ur.grantedAt = time.Now()
	ur.updatedAt = time.Now()

	return nil
}
