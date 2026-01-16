package shared

import (
	"errors"
	"net/http"
)

var (
	// User validation errors
	ErrInvalidUserID    = errors.New("invalid user ID")
	ErrEmptyUserID      = errors.New("user ID cannot be empty")
	ErrInvalidEmail     = errors.New("invalid email format")
	ErrEmptyEmail       = errors.New("email cannot be empty")
	ErrUnverifiedEmail  = errors.New("email address is not verified")

	// User state errors
	ErrUserNotFound     = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")

	// Authentication errors
	ErrInvalidToken     = errors.New("invalid token")
	ErrExpiredToken     = errors.New("token has expired")
	ErrMissingToken     = errors.New("token not found")
	ErrUnauthorized     = errors.New("unauthorized access")

	// Profile errors
	ErrInvalidProfile   = errors.New("invalid profile data")
)

// AppError represents an application error with HTTP status code and error code
type AppError struct {
	StatusCode int
	Code       string
	Message    string
	Err        error
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new application error
func NewAppError(statusCode int, code, message string, err error) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
		Err:        err,
	}
}

// Common error constructors
func NewBadRequestError(code, message string, err error) *AppError {
	return NewAppError(http.StatusBadRequest, code, message, err)
}

func NewUnauthorizedError(code, message string, err error) *AppError {
	return NewAppError(http.StatusUnauthorized, code, message, err)
}

func NewForbiddenError(code, message string, err error) *AppError {
	return NewAppError(http.StatusForbidden, code, message, err)
}

func NewNotFoundError(code, message string, err error) *AppError {
	return NewAppError(http.StatusNotFound, code, message, err)
}

func NewConflictError(code, message string, err error) *AppError {
	return NewAppError(http.StatusConflict, code, message, err)
}

func NewInternalError(code, message string, err error) *AppError {
	return NewAppError(http.StatusInternalServerError, code, message, err)
}
