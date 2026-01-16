package dto

// RequestRoleRequest represents a request to request a role
type RequestRoleRequest struct {
	RequestedRole string `json:"requested_role" binding:"required"`
}

// RequestRoleResponse represents the response after requesting a role
type RequestRoleResponse struct {
	RequestID     string `json:"request_id"`
	UserID        string `json:"user_id"`
	UserEmail     string `json:"user_email"`
	RequestedRole string `json:"requested_role"`
	Status        string `json:"status"`
	RequestedAt   int64  `json:"requested_at"`
}

// ApproveRequestRequest represents a request to approve a role request
type ApproveRequestRequest struct {
	RequestID string `json:"request_id" binding:"required"`
	Notes     string `json:"notes"`
}

// RejectRequestRequest represents a request to reject a role request
type RejectRequestRequest struct{
	RequestID string `json:"request_id" binding:"required"`
	Notes     string `json:"notes"`
}

// RoleRequestDTO represents a role request in API responses
type RoleRequestDTO struct {
	RequestID     string  `json:"request_id"`
	UserID        string  `json:"user_id"`
	UserEmail     string  `json:"user_email"`
	RequestedRole string  `json:"requested_role"`
	Status        string  `json:"status"`
	RequestedAt   int64   `json:"requested_at"`
	ProcessedAt   *int64  `json:"processed_at,omitempty"`
	ProcessedBy   string  `json:"processed_by,omitempty"`
	Notes         string  `json:"notes,omitempty"`
}

// UserRoleDTO represents a user role in API responses
type UserRoleDTO struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Status    string `json:"status"`
	GrantedAt int64  `json:"granted_at"`
	GrantedBy string `json:"granted_by"`
}

// ListPendingRequestsResponse represents the response for listing pending requests
type ListPendingRequestsResponse struct {
	Requests []*RoleRequestDTO `json:"requests"`
	Count    int               `json:"count"`
}

// ListUsersResponse represents the response for listing users by role
type ListUsersResponse struct {
	Users []*UserRoleDTO `json:"users"`
	Count int            `json:"count"`
}

// CheckUserRoleResponse represents the response for checking user's role
type CheckUserRoleResponse struct {
	UserID string  `json:"user_id"`
	Email  string  `json:"email"`
	Role   string  `json:"role"`
	Status string  `json:"status"`
}
