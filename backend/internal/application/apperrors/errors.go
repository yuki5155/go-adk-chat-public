package apperrors

import "github.com/yuki5155/go-google-auth/internal/domain/shared"

func NewBadRequestError(code, message string, err error) error {
	return shared.NewBadRequestError(code, message, err)
}

func NewUnauthorizedError(code, message string, err error) error {
	return shared.NewUnauthorizedError(code, message, err)
}

func NewForbiddenError(code, message string, err error) error {
	return shared.NewForbiddenError(code, message, err)
}

func NewInternalError(code, message string, err error) error {
	return shared.NewInternalError(code, message, err)
}
