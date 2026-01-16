package role

import (
	"testing"
	"time"

	"github.com/yuki5155/go-google-auth/internal/domain/user"
)

func TestNewRoleRequest(t *testing.T) {
	tests := []struct {
		name          string
		requestID     string
		userID        string
		userEmail     string
		requestedRole user.Role
		wantErr       error
	}{
		{
			name:          "Valid role request creation",
			requestID:     "req-123",
			userID:        "user-123",
			userEmail:     "test@example.com",
			requestedRole: user.RoleSubscriber,
			wantErr:       nil,
		},
		{
			name:          "Empty request ID",
			requestID:     "",
			userID:        "user-123",
			userEmail:     "test@example.com",
			requestedRole: user.RoleSubscriber,
			wantErr:       ErrInvalidRequestID,
		},
		{
			name:          "Empty user ID",
			requestID:     "req-123",
			userID:        "",
			userEmail:     "test@example.com",
			requestedRole: user.RoleSubscriber,
			wantErr:       ErrInvalidUserID,
		},
		{
			name:          "Empty email",
			requestID:     "req-123",
			userID:        "user-123",
			userEmail:     "",
			requestedRole: user.RoleSubscriber,
			wantErr:       ErrInvalidEmail,
		},
		{
			name:          "Invalid role",
			requestID:     "req-123",
			userID:        "user-123",
			userEmail:     "test@example.com",
			requestedRole: user.Role("invalid"),
			wantErr:       ErrInvalidRole,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRoleRequest(tt.requestID, tt.userID, tt.userEmail, tt.requestedRole)

			if err != tt.wantErr {
				t.Errorf("NewRoleRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				if got.RequestID() != tt.requestID {
					t.Errorf("RequestID() = %v, want %v", got.RequestID(), tt.requestID)
				}
				if got.UserID() != tt.userID {
					t.Errorf("UserID() = %v, want %v", got.UserID(), tt.userID)
				}
				if got.UserEmail() != tt.userEmail {
					t.Errorf("UserEmail() = %v, want %v", got.UserEmail(), tt.userEmail)
				}
				if got.RequestedRole() != tt.requestedRole {
					t.Errorf("RequestedRole() = %v, want %v", got.RequestedRole(), tt.requestedRole)
				}
				if got.Status() != RequestStatusPending {
					t.Errorf("Status() = %v, want %v", got.Status(), RequestStatusPending)
				}
				if !got.IsPending() {
					t.Error("Expected request to be pending")
				}
			}
		})
	}
}

func TestRoleRequest_Approve(t *testing.T) {
	request, _ := NewRoleRequest("req-123", "user-123", "test@example.com", user.RoleSubscriber)

	err := request.Approve("admin@example.com", "Approved for testing")
	if err != nil {
		t.Errorf("Approve() error = %v", err)
	}

	if request.Status() != RequestStatusApproved {
		t.Errorf("Status() = %v, want %v", request.Status(), RequestStatusApproved)
	}

	if !request.IsApproved() {
		t.Error("Expected request to be approved")
	}

	if request.IsPending() {
		t.Error("Expected request to not be pending")
	}

	if request.ProcessedBy() != "admin@example.com" {
		t.Errorf("ProcessedBy() = %v, want %v", request.ProcessedBy(), "admin@example.com")
	}

	if request.Notes() != "Approved for testing" {
		t.Errorf("Notes() = %v, want %v", request.Notes(), "Approved for testing")
	}

	if request.ProcessedAt() == nil {
		t.Error("Expected ProcessedAt to be set")
	}
}

func TestRoleRequest_Reject(t *testing.T) {
	request, _ := NewRoleRequest("req-123", "user-123", "test@example.com", user.RoleSubscriber)

	err := request.Reject("admin@example.com", "Rejected for testing")
	if err != nil {
		t.Errorf("Reject() error = %v", err)
	}

	if request.Status() != RequestStatusRejected {
		t.Errorf("Status() = %v, want %v", request.Status(), RequestStatusRejected)
	}

	if !request.IsRejected() {
		t.Error("Expected request to be rejected")
	}

	if request.IsPending() {
		t.Error("Expected request to not be pending")
	}

	if request.ProcessedBy() != "admin@example.com" {
		t.Errorf("ProcessedBy() = %v, want %v", request.ProcessedBy(), "admin@example.com")
	}

	if request.Notes() != "Rejected for testing" {
		t.Errorf("Notes() = %v, want %v", request.Notes(), "Rejected for testing")
	}

	if request.ProcessedAt() == nil {
		t.Error("Expected ProcessedAt to be set")
	}
}

func TestRoleRequest_ApproveAlreadyProcessed(t *testing.T) {
	request, _ := NewRoleRequest("req-123", "user-123", "test@example.com", user.RoleSubscriber)
	err := request.Approve("admin@example.com", "First approval")
	if err != nil {
		t.Fatalf("First Approve() should not error, got: %v", err)
	}

	err = request.Approve("admin2@example.com", "Second approval")
	if err != ErrRequestAlreadyProcessed {
		t.Errorf("Approve() error = %v, want %v", err, ErrRequestAlreadyProcessed)
	}
}

func TestRoleRequest_RejectAlreadyProcessed(t *testing.T) {
	request, _ := NewRoleRequest("req-123", "user-123", "test@example.com", user.RoleSubscriber)
	err := request.Reject("admin@example.com", "First rejection")
	if err != nil {
		t.Fatalf("First Reject() should not error, got: %v", err)
	}

	err = request.Reject("admin2@example.com", "Second rejection")
	if err != ErrRequestAlreadyProcessed {
		t.Errorf("Reject() error = %v, want %v", err, ErrRequestAlreadyProcessed)
	}
}

func TestReconstructRoleRequest(t *testing.T) {
	now := time.Now()
	requestedAt := now.Add(-1 * time.Hour)
	processedAt := now

	request := ReconstructRoleRequest(
		"req-123",
		"user-123",
		"test@example.com",
		user.RoleSubscriber,
		RequestStatusApproved,
		requestedAt,
		&processedAt,
		"admin@example.com",
		"Test notes",
		now,
		now,
	)

	if request.RequestID() != "req-123" {
		t.Errorf("RequestID() = %v, want %v", request.RequestID(), "req-123")
	}
	if request.Status() != RequestStatusApproved {
		t.Errorf("Status() = %v, want %v", request.Status(), RequestStatusApproved)
	}
	if request.ProcessedBy() != "admin@example.com" {
		t.Errorf("ProcessedBy() = %v, want %v", request.ProcessedBy(), "admin@example.com")
	}
	if request.Notes() != "Test notes" {
		t.Errorf("Notes() = %v, want %v", request.Notes(), "Test notes")
	}
}
