package schema

// RequestRoleRequest represents the HTTP request body for requesting a role
type RequestRoleRequest struct {
	RequestedRole string `json:"requested_role" binding:"required"`
}

// ApproveRequestRequest represents the HTTP request body for approving a request
type ApproveRequestRequest struct {
	RequestID string `json:"request_id" binding:"required"`
	Notes     string `json:"notes"`
}

// RejectRequestRequest represents the HTTP request body for rejecting a request
type RejectRequestRequest struct {
	RequestID string `json:"request_id" binding:"required"`
	Notes     string `json:"notes"`
}
