package schema

// RequestRoleErrorResponse represents the HTTP error response for role request
type RequestRoleErrorResponse struct {
	Success bool         `json:"success"`
	Error   *ErrorDetail `json:"error"`
}

// GetMyRoleErrorResponse represents the HTTP error response for get my role
type GetMyRoleErrorResponse struct {
	Success bool         `json:"success"`
	Error   *ErrorDetail `json:"error"`
}

// ListPendingRequestsErrorResponse represents the HTTP error response for list pending requests
type ListPendingRequestsErrorResponse struct {
	Success bool         `json:"success"`
	Error   *ErrorDetail `json:"error"`
}

// ApproveRejectErrorResponse represents the HTTP error response for approve/reject operations
type ApproveRejectErrorResponse struct {
	Success bool         `json:"success"`
	Error   *ErrorDetail `json:"error"`
}

// ListUsersErrorResponse represents the HTTP error response for list users
type ListUsersErrorResponse struct {
	Success bool         `json:"success"`
	Error   *ErrorDetail `json:"error"`
}

// NewRequestRoleError creates an error response for role request
func NewRequestRoleError(code string, message string) *RequestRoleErrorResponse {
	return &RequestRoleErrorResponse{
		Success: false,
		Error: &ErrorDetail{
			Code:    code,
			Message: message,
		},
	}
}

// NewGetMyRoleError creates an error response for get my role
func NewGetMyRoleError(code string, message string) *GetMyRoleErrorResponse {
	return &GetMyRoleErrorResponse{
		Success: false,
		Error: &ErrorDetail{
			Code:    code,
			Message: message,
		},
	}
}

// NewListPendingRequestsError creates an error response for list pending requests
func NewListPendingRequestsError(code string, message string) *ListPendingRequestsErrorResponse {
	return &ListPendingRequestsErrorResponse{
		Success: false,
		Error: &ErrorDetail{
			Code:    code,
			Message: message,
		},
	}
}

// NewApproveRejectError creates an error response for approve/reject operations
func NewApproveRejectError(code string, message string) *ApproveRejectErrorResponse {
	return &ApproveRejectErrorResponse{
		Success: false,
		Error: &ErrorDetail{
			Code:    code,
			Message: message,
		},
	}
}

// NewListUsersError creates an error response for list users
func NewListUsersError(code string, message string) *ListUsersErrorResponse {
	return &ListUsersErrorResponse{
		Success: false,
		Error: &ErrorDetail{
			Code:    code,
			Message: message,
		},
	}
}
