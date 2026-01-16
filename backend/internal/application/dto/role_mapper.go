package dto

import (
	"github.com/yuki5155/go-google-auth/internal/domain/role"
)

// ToRoleRequestDTO converts a domain RoleRequest to DTO
func ToRoleRequestDTO(request *role.RoleRequest) *RoleRequestDTO {
	var processedAt *int64
	if request.ProcessedAt() != nil {
		t := request.ProcessedAt().Unix()
		processedAt = &t
	}

	return &RoleRequestDTO{
		RequestID:     request.RequestID(),
		UserID:        request.UserID(),
		UserEmail:     request.UserEmail(),
		RequestedRole: request.RequestedRole().String(),
		Status:        string(request.Status()),
		RequestedAt:   request.RequestedAt().Unix(),
		ProcessedAt:   processedAt,
		ProcessedBy:   request.ProcessedBy(),
		Notes:         request.Notes(),
	}
}

// ToRequestRoleResponse converts a domain RoleRequest to RequestRoleResponse
func ToRequestRoleResponse(request *role.RoleRequest) *RequestRoleResponse {
	return &RequestRoleResponse{
		RequestID:     request.RequestID(),
		UserID:        request.UserID(),
		UserEmail:     request.UserEmail(),
		RequestedRole: request.RequestedRole().String(),
		Status:        string(request.Status()),
		RequestedAt:   request.RequestedAt().Unix(),
	}
}

// ToUserRoleDTO converts a domain UserRole to DTO
func ToUserRoleDTO(userRole *role.UserRole) *UserRoleDTO {
	return &UserRoleDTO{
		UserID:    userRole.UserID(),
		Email:     userRole.Email(),
		Role:      userRole.Role().String(),
		Status:    string(userRole.Status()),
		GrantedAt: userRole.GrantedAt().Unix(),
		GrantedBy: userRole.GrantedBy(),
	}
}

// ToCheckUserRoleResponse converts a domain UserRole to CheckUserRoleResponse
func ToCheckUserRoleResponse(userRole *role.UserRole) *CheckUserRoleResponse {
	return &CheckUserRoleResponse{
		UserID: userRole.UserID(),
		Email:  userRole.Email(),
		Role:   userRole.Role().String(),
		Status: string(userRole.Status()),
	}
}

// ToListPendingRequestsResponse converts a slice of RoleRequests to ListPendingRequestsResponse
func ToListPendingRequestsResponse(requests []*role.RoleRequest) *ListPendingRequestsResponse {
	requestDTOs := make([]*RoleRequestDTO, len(requests))
	for i, req := range requests {
		requestDTOs[i] = ToRoleRequestDTO(req)
	}

	return &ListPendingRequestsResponse{
		Requests: requestDTOs,
		Count:    len(requestDTOs),
	}
}

// ToListUsersResponse converts a slice of UserRoles to ListUsersResponse
func ToListUsersResponse(userRoles []*role.UserRole) *ListUsersResponse {
	userDTOs := make([]*UserRoleDTO, len(userRoles))
	for i, ur := range userRoles {
		userDTOs[i] = ToUserRoleDTO(ur)
	}

	return &ListUsersResponse{
		Users: userDTOs,
		Count: len(userDTOs),
	}
}
