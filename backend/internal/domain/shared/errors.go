package shared

import "errors"

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
