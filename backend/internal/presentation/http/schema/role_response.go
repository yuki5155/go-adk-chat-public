package schema

// RequestRoleResponse represents the HTTP response for role request
type RequestRoleResponse struct {
	Success bool                   `json:"success"`
	Data    *RequestRoleData       `json:"data,omitempty"`
	Error   *ErrorDetail           `json:"error,omitempty"`
}

// RequestRoleData contains the role request data
type RequestRoleData struct {
	RequestID     string `json:"request_id"`
	UserID        string `json:"user_id"`
	UserEmail     string `json:"user_email"`
	RequestedRole string `json:"requested_role"`
	Status        string `json:"status"`
	RequestedAt   int64  `json:"requested_at"`
}

// GetMyRoleResponse represents the HTTP response for getting user's role
type GetMyRoleResponse struct {
	Success bool              `json:"success"`
	Data    *UserRoleData     `json:"data,omitempty"`
	Error   *ErrorDetail      `json:"error,omitempty"`
}

// UserRoleData contains user role information
type UserRoleData struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Status string `json:"status"`
}

// ListPendingRequestsResponse represents the HTTP response for listing pending requests
type ListPendingRequestsResponse struct {
	Success bool                       `json:"success"`
	Data    *PendingRequestsData       `json:"data,omitempty"`
	Error   *ErrorDetail               `json:"error,omitempty"`
}

// PendingRequestsData contains the list of pending requests
type PendingRequestsData struct {
	Requests []*RoleRequestItem `json:"requests"`
	Count    int                `json:"count"`
}

// RoleRequestItem represents a single role request in the list
type RoleRequestItem struct {
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

// ApproveRejectResponse represents the HTTP response for approve/reject operations
type ApproveRejectResponse struct {
	Success bool         `json:"success"`
	Data    *MessageData `json:"data,omitempty"`
	Error   *ErrorDetail `json:"error,omitempty"`
}

// MessageData contains a simple message
type MessageData struct {
	Message string `json:"message"`
}

// ListUsersResponse represents the HTTP response for listing users
type ListUsersResponse struct {
	Success bool         `json:"success"`
	Data    *UsersData   `json:"data,omitempty"`
	Error   *ErrorDetail `json:"error,omitempty"`
}

// UsersData contains the list of users
type UsersData struct {
	Users []*UserRoleItem `json:"users"`
	Count int             `json:"count"`
}

// UserRoleItem represents a single user role in the list
type UserRoleItem struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Status    string `json:"status"`
	GrantedAt int64  `json:"granted_at"`
	GrantedBy string `json:"granted_by"`
}

// ErrorDetail represents error information
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
