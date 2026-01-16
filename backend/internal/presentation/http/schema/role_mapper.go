package schema

import (
	"github.com/yuki5155/go-google-auth/internal/application/dto"
)

// ToRequestRoleResponse converts DTO to HTTP response schema
func ToRequestRoleResponse(roleReqDTO *dto.RequestRoleResponse) *RequestRoleResponse {
	return &RequestRoleResponse{
		Success: true,
		Data: &RequestRoleData{
			RequestID:     roleReqDTO.RequestID,
			UserID:        roleReqDTO.UserID,
			UserEmail:     roleReqDTO.UserEmail,
			RequestedRole: roleReqDTO.RequestedRole,
			Status:        roleReqDTO.Status,
			RequestedAt:   roleReqDTO.RequestedAt,
		},
	}
}

// ToGetMyRoleResponse converts DTO to HTTP response schema
func ToGetMyRoleResponse(roleDTO *dto.CheckUserRoleResponse) *GetMyRoleResponse {
	return &GetMyRoleResponse{
		Success: true,
		Data: &UserRoleData{
			UserID: roleDTO.UserID,
			Email:  roleDTO.Email,
			Role:   roleDTO.Role,
			Status: roleDTO.Status,
		},
	}
}

// ToListPendingRequestsResponse converts DTO to HTTP response schema
func ToListPendingRequestsResponse(listDTO *dto.ListPendingRequestsResponse) *ListPendingRequestsResponse {
	requests := make([]*RoleRequestItem, len(listDTO.Requests))
	for i, req := range listDTO.Requests {
		requests[i] = &RoleRequestItem{
			RequestID:     req.RequestID,
			UserID:        req.UserID,
			UserEmail:     req.UserEmail,
			RequestedRole: req.RequestedRole,
			Status:        req.Status,
			RequestedAt:   req.RequestedAt,
			ProcessedAt:   req.ProcessedAt,
			ProcessedBy:   req.ProcessedBy,
			Notes:         req.Notes,
		}
	}

	return &ListPendingRequestsResponse{
		Success: true,
		Data: &PendingRequestsData{
			Requests: requests,
			Count:    listDTO.Count,
		},
	}
}

// ToApproveRejectResponse creates a success response for approve/reject operations
func ToApproveRejectResponse(message string) *ApproveRejectResponse {
	return &ApproveRejectResponse{
		Success: true,
		Data: &MessageData{
			Message: message,
		},
	}
}

// ToListUsersResponse converts DTO to HTTP response schema
func ToListUsersResponse(listDTO *dto.ListUsersResponse) *ListUsersResponse {
	users := make([]*UserRoleItem, len(listDTO.Users))
	for i, u := range listDTO.Users {
		users[i] = &UserRoleItem{
			UserID:    u.UserID,
			Email:     u.Email,
			Role:      u.Role,
			Status:    u.Status,
			GrantedAt: u.GrantedAt,
			GrantedBy: u.GrantedBy,
		}
	}

	return &ListUsersResponse{
		Success: true,
		Data: &UsersData{
			Users: users,
			Count: listDTO.Count,
		},
	}
}

// ToErrorResponse creates an error response
func ToErrorResponse(code string, message string) interface{} {
	return map[string]interface{}{
		"success": false,
		"error": &ErrorDetail{
			Code:    code,
			Message: message,
		},
	}
}
