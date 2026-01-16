package role

import "errors"

var (
	// ErrInvalidUserID is returned when user ID is invalid
	ErrInvalidUserID = errors.New("invalid user ID")

	// ErrInvalidEmail is returned when email is invalid
	ErrInvalidEmail = errors.New("invalid email")

	// ErrInvalidRole is returned when role is invalid
	ErrInvalidRole = errors.New("invalid role")

	// ErrInvalidRequestID is returned when request ID is invalid
	ErrInvalidRequestID = errors.New("invalid request ID")

	// ErrRequestAlreadyProcessed is returned when trying to process an already processed request
	ErrRequestAlreadyProcessed = errors.New("request already processed")

	// ErrUserRoleNotFound is returned when user role is not found
	ErrUserRoleNotFound = errors.New("user role not found")

	// ErrRoleRequestNotFound is returned when role request is not found
	ErrRoleRequestNotFound = errors.New("role request not found")

	// ErrDuplicateRoleRequest is returned when a pending request already exists
	ErrDuplicateRoleRequest = errors.New("duplicate role request: a pending request already exists for this user")
)
