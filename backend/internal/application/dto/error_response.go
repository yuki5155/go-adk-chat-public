package dto

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// SuccessResponse represents a standard success response
type SuccessResponse struct {
	Message string `json:"message"`
}

// NewErrorResponse creates a new error response
func NewErrorResponse(errorCode string, message string) *ErrorResponse {
	return &ErrorResponse{
		Error:   errorCode,
		Message: message,
	}
}

// NewSuccessResponse creates a new success response
func NewSuccessResponse(message string) *SuccessResponse {
	return &SuccessResponse{
		Message: message,
	}
}

// Common error codes
const (
	ErrCodeInvalidRequest       = "invalid_request"
	ErrCodeInvalidRole          = "invalid_role"
	ErrCodeDuplicateRequest     = "duplicate_request"
	ErrCodeInternalError        = "internal_error"
	ErrCodeNotFound             = "not_found"
	ErrCodeAlreadyProcessed     = "already_processed"
	ErrCodeUnauthorized         = "unauthorized"
	ErrCodeForbidden            = "forbidden"
)
