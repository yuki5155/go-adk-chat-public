package role

import (
	"testing"
	"time"

	"github.com/yuki5155/go-google-auth/internal/domain/user"
)

func TestNewUserRole(t *testing.T) {
	tests := []struct {
		name      string
		userID    string
		email     string
		role      user.Role
		grantedBy string
		wantErr   error
	}{
		{
			name:      "Valid user role creation",
			userID:    "user-123",
			email:     "test@example.com",
			role:      user.RoleSubscriber,
			grantedBy: "admin@example.com",
			wantErr:   nil,
		},
		{
			name:      "Empty user ID",
			userID:    "",
			email:     "test@example.com",
			role:      user.RoleSubscriber,
			grantedBy: "admin@example.com",
			wantErr:   ErrInvalidUserID,
		},
		{
			name:      "Empty email",
			userID:    "user-123",
			email:     "",
			role:      user.RoleSubscriber,
			grantedBy: "admin@example.com",
			wantErr:   ErrInvalidEmail,
		},
		{
			name:      "Invalid role",
			userID:    "user-123",
			email:     "test@example.com",
			role:      user.Role("invalid"),
			grantedBy: "admin@example.com",
			wantErr:   ErrInvalidRole,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUserRole(tt.userID, tt.email, tt.role, tt.grantedBy)

			if err != tt.wantErr {
				t.Errorf("NewUserRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				if got.UserID() != tt.userID {
					t.Errorf("UserID() = %v, want %v", got.UserID(), tt.userID)
				}
				if got.Email() != tt.email {
					t.Errorf("Email() = %v, want %v", got.Email(), tt.email)
				}
				if got.Role() != tt.role {
					t.Errorf("Role() = %v, want %v", got.Role(), tt.role)
				}
				if got.Status() != StatusActive {
					t.Errorf("Status() = %v, want %v", got.Status(), StatusActive)
				}
				if got.GrantedBy() != tt.grantedBy {
					t.Errorf("GrantedBy() = %v, want %v", got.GrantedBy(), tt.grantedBy)
				}
			}
		})
	}
}

func TestUserRole_Suspend(t *testing.T) {
	userRole, _ := NewUserRole("user-123", "test@example.com", user.RoleSubscriber, "admin@example.com")

	if !userRole.IsActive() {
		t.Error("Expected user role to be active initially")
	}

	userRole.Suspend()

	if userRole.Status() != StatusSuspended {
		t.Errorf("Status() = %v, want %v", userRole.Status(), StatusSuspended)
	}

	if userRole.IsActive() {
		t.Error("Expected user role to not be active after suspension")
	}
}

func TestUserRole_Activate(t *testing.T) {
	userRole, _ := NewUserRole("user-123", "test@example.com", user.RoleSubscriber, "admin@example.com")
	userRole.Suspend()

	userRole.Activate()

	if userRole.Status() != StatusActive {
		t.Errorf("Status() = %v, want %v", userRole.Status(), StatusActive)
	}

	if !userRole.IsActive() {
		t.Error("Expected user role to be active after activation")
	}
}

func TestUserRole_ChangeRole(t *testing.T) {
	userRole, _ := NewUserRole("user-123", "test@example.com", user.RoleSubscriber, "admin@example.com")

	err := userRole.ChangeRole(user.RoleAdmin, "superadmin@example.com")
	if err != nil {
		t.Errorf("ChangeRole() error = %v", err)
	}

	if userRole.Role() != user.RoleAdmin {
		t.Errorf("Role() = %v, want %v", userRole.Role(), user.RoleAdmin)
	}

	if userRole.GrantedBy() != "superadmin@example.com" {
		t.Errorf("GrantedBy() = %v, want %v", userRole.GrantedBy(), "superadmin@example.com")
	}
}

func TestUserRole_ChangeRole_Invalid(t *testing.T) {
	userRole, _ := NewUserRole("user-123", "test@example.com", user.RoleSubscriber, "admin@example.com")

	err := userRole.ChangeRole(user.Role("invalid"), "admin@example.com")
	if err != ErrInvalidRole {
		t.Errorf("ChangeRole() error = %v, want %v", err, ErrInvalidRole)
	}
}

func TestReconstructUserRole(t *testing.T) {
	now := time.Now()
	grantedAt := now.Add(-24 * time.Hour)

	userRole := ReconstructUserRole(
		"user-123",
		"test@example.com",
		user.RoleSubscriber,
		StatusActive,
		grantedAt,
		"admin@example.com",
		now,
		now,
	)

	if userRole.UserID() != "user-123" {
		t.Errorf("UserID() = %v, want %v", userRole.UserID(), "user-123")
	}
	if userRole.Email() != "test@example.com" {
		t.Errorf("Email() = %v, want %v", userRole.Email(), "test@example.com")
	}
	if userRole.Role() != user.RoleSubscriber {
		t.Errorf("Role() = %v, want %v", userRole.Role(), user.RoleSubscriber)
	}
	if userRole.Status() != StatusActive {
		t.Errorf("Status() = %v, want %v", userRole.Status(), StatusActive)
	}
}
