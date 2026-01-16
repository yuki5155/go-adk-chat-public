package role

import (
	"time"

	"github.com/yuki5155/go-google-auth/internal/domain/user"
)

// RoleRequest represents a user's request for a role subscription
type RoleRequest struct {
	requestID     string
	userID        string
	userEmail     string
	requestedRole user.Role
	status        RequestStatus
	requestedAt   time.Time
	processedAt   *time.Time
	processedBy   string
	notes         string
	createdAt     time.Time
	updatedAt     time.Time
}

// RequestStatus represents the status of a role request
type RequestStatus string

const (
	// RequestStatusPending represents a pending request
	RequestStatusPending RequestStatus = "pending"
	// RequestStatusApproved represents an approved request
	RequestStatusApproved RequestStatus = "approved"
	// RequestStatusRejected represents a rejected request
	RequestStatusRejected RequestStatus = "rejected"
)

// NewRoleRequest creates a new RoleRequest entity
func NewRoleRequest(
	requestID string,
	userID string,
	userEmail string,
	requestedRole user.Role,
) (*RoleRequest, error) {
	if requestID == "" {
		return nil, ErrInvalidRequestID
	}

	if userID == "" {
		return nil, ErrInvalidUserID
	}

	if userEmail == "" {
		return nil, ErrInvalidEmail
	}

	if !requestedRole.IsValid() {
		return nil, ErrInvalidRole
	}

	now := time.Now()

	return &RoleRequest{
		requestID:     requestID,
		userID:        userID,
		userEmail:     userEmail,
		requestedRole: requestedRole,
		status:        RequestStatusPending,
		requestedAt:   now,
		createdAt:     now,
		updatedAt:     now,
	}, nil
}

// ReconstructRoleRequest reconstructs a RoleRequest from persistence
func ReconstructRoleRequest(
	requestID string,
	userID string,
	userEmail string,
	requestedRole user.Role,
	status RequestStatus,
	requestedAt time.Time,
	processedAt *time.Time,
	processedBy string,
	notes string,
	createdAt time.Time,
	updatedAt time.Time,
) *RoleRequest {
	return &RoleRequest{
		requestID:     requestID,
		userID:        userID,
		userEmail:     userEmail,
		requestedRole: requestedRole,
		status:        status,
		requestedAt:   requestedAt,
		processedAt:   processedAt,
		processedBy:   processedBy,
		notes:         notes,
		createdAt:     createdAt,
		updatedAt:     updatedAt,
	}
}

// RequestID returns the request ID
func (rr *RoleRequest) RequestID() string {
	return rr.requestID
}

// UserID returns the user ID
func (rr *RoleRequest) UserID() string {
	return rr.userID
}

// UserEmail returns the user's email
func (rr *RoleRequest) UserEmail() string {
	return rr.userEmail
}

// RequestedRole returns the requested role
func (rr *RoleRequest) RequestedRole() user.Role {
	return rr.requestedRole
}

// Status returns the request status
func (rr *RoleRequest) Status() RequestStatus {
	return rr.status
}

// RequestedAt returns when the request was made
func (rr *RoleRequest) RequestedAt() time.Time {
	return rr.requestedAt
}

// ProcessedAt returns when the request was processed
func (rr *RoleRequest) ProcessedAt() *time.Time {
	return rr.processedAt
}

// ProcessedBy returns who processed the request
func (rr *RoleRequest) ProcessedBy() string {
	return rr.processedBy
}

// Notes returns the request notes
func (rr *RoleRequest) Notes() string {
	return rr.notes
}

// CreatedAt returns when the request was created
func (rr *RoleRequest) CreatedAt() time.Time {
	return rr.createdAt
}

// UpdatedAt returns when the request was last updated
func (rr *RoleRequest) UpdatedAt() time.Time {
	return rr.updatedAt
}

// IsPending returns whether the request is pending
func (rr *RoleRequest) IsPending() bool {
	return rr.status == RequestStatusPending
}

// IsApproved returns whether the request is approved
func (rr *RoleRequest) IsApproved() bool {
	return rr.status == RequestStatusApproved
}

// IsRejected returns whether the request is rejected
func (rr *RoleRequest) IsRejected() bool {
	return rr.status == RequestStatusRejected
}

// Approve approves the role request
func (rr *RoleRequest) Approve(approvedBy string, notes string) error {
	if !rr.IsPending() {
		return ErrRequestAlreadyProcessed
	}

	now := time.Now()
	rr.status = RequestStatusApproved
	rr.processedAt = &now
	rr.processedBy = approvedBy
	rr.notes = notes
	rr.updatedAt = now

	return nil
}

// Reject rejects the role request
func (rr *RoleRequest) Reject(rejectedBy string, notes string) error {
	if !rr.IsPending() {
		return ErrRequestAlreadyProcessed
	}

	now := time.Now()
	rr.status = RequestStatusRejected
	rr.processedAt = &now
	rr.processedBy = rejectedBy
	rr.notes = notes
	rr.updatedAt = now

	return nil
}
