package admin

import (
	"github.com/yuki5155/go-google-auth/internal/domain/shared"
	"github.com/yuki5155/go-google-auth/internal/domain/user"
)

// RequestRoleCommand represents a command to request a role
type RequestRoleCommand struct {
	UserID        string
	UserEmail     string
	RequestedRole string
}

// Validate validates the command
func (c *RequestRoleCommand) Validate() error {
	if c.UserID == "" {
		return shared.NewBadRequestError("INVALID_REQUEST", "User ID is required", nil)
	}
	if c.UserEmail == "" {
		return shared.NewBadRequestError("INVALID_REQUEST", "User email is required", nil)
	}
	role := user.Role(c.RequestedRole)
	if !role.IsValid() {
		return shared.NewBadRequestError("INVALID_REQUEST", "Invalid role type", nil)
	}
	if role != user.RoleSubscriber && role != user.RoleUser {
		return shared.NewBadRequestError("INVALID_REQUEST", "Only 'subscriber' or 'user' roles can be requested", nil)
	}
	return nil
}

// ApproveRequestCommand represents a command to approve a role request
type ApproveRequestCommand struct {
	RequestID   string
	ApprovedBy  string
	Notes       string
}

// Validate validates the command
func (c *ApproveRequestCommand) Validate() error {
	if c.RequestID == "" {
		return shared.NewBadRequestError("INVALID_REQUEST", "Request ID is required", nil)
	}
	if c.ApprovedBy == "" {
		return shared.NewBadRequestError("INVALID_REQUEST", "Approver email is required", nil)
	}
	return nil
}

// RejectRequestCommand represents a command to reject a role request
type RejectRequestCommand struct {
	RequestID  string
	RejectedBy string
	Notes      string
}

// Validate validates the command
func (c *RejectRequestCommand) Validate() error {
	if c.RequestID == "" {
		return shared.NewBadRequestError("INVALID_REQUEST", "Request ID is required", nil)
	}
	if c.RejectedBy == "" {
		return shared.NewBadRequestError("INVALID_REQUEST", "Rejecter email is required", nil)
	}
	if c.Notes == "" {
		return shared.NewBadRequestError("INVALID_REQUEST", "Rejection notes are required", nil)
	}
	return nil
}

// ListUsersByRoleQuery represents a query to list users by role
type ListUsersByRoleQuery struct {
	Role string
}

// Validate validates the query
func (q *ListUsersByRoleQuery) Validate() error {
	if !user.Role(q.Role).IsValid() {
		return shared.NewBadRequestError("INVALID_REQUEST", "Invalid role parameter", nil)
	}
	return nil
}
