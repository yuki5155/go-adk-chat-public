package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yuki5155/go-google-auth/internal/application/admin"
	"github.com/yuki5155/go-google-auth/internal/application/dto"
	"github.com/yuki5155/go-google-auth/internal/application/ports"
	"github.com/yuki5155/go-google-auth/internal/domain/shared"
	"github.com/yuki5155/go-google-auth/internal/domain/user"
)

// RoleHandler is a thin handler that delegates to use cases
type RoleHandler struct {
	requestRoleUC     *admin.RequestRoleUseCase
	listPendingUC     *admin.ListPendingRequestsUseCase
	approveRequestUC  *admin.ApproveRequestUseCase
	rejectRequestUC   *admin.RejectRequestUseCase
	listUsersByRoleUC *admin.ListUsersByRoleUseCase
}

// NewRoleHandler creates a new thin role handler
func NewRoleHandler(
	requestRoleUC *admin.RequestRoleUseCase,
	listPendingUC *admin.ListPendingRequestsUseCase,
	approveRequestUC *admin.ApproveRequestUseCase,
	rejectRequestUC *admin.RejectRequestUseCase,
	listUsersByRoleUC *admin.ListUsersByRoleUseCase,
) *RoleHandler {
	return &RoleHandler{
		requestRoleUC:     requestRoleUC,
		listPendingUC:     listPendingUC,
		approveRequestUC:  approveRequestUC,
		rejectRequestUC:   rejectRequestUC,
		listUsersByRoleUC: listUsersByRoleUC,
	}
}

// RequestRole handles role request creation
func (h *RoleHandler) RequestRole(c *gin.Context) {
	// Get claims from context
	claims := getClaims(c)
	if claims == nil {
		return // Error already handled by getClaims
	}

	// Parse request body
	var req dto.RequestRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(shared.NewBadRequestError("INVALID_REQUEST", "Missing or invalid requested_role field", err))
		return
	}

	// Build command
	cmd := admin.RequestRoleCommand{
		UserID:        claims.UserID,
		UserEmail:     claims.Email,
		RequestedRole: user.Role(req.RequestedRole),
	}

	// Execute use case
	dto, err := h.requestRoleUC.Execute(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    dto,
	})
}

// ListPendingRequests lists all pending role requests (admin only)
func (h *RoleHandler) ListPendingRequests(c *gin.Context) {
	// Execute use case
	dto, err := h.listPendingUC.Execute(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    dto,
	})
}

// ApproveRequest approves a role request (admin only)
func (h *RoleHandler) ApproveRequest(c *gin.Context) {
	// Get claims from context
	claims := getClaims(c)
	if claims == nil {
		return
	}

	// Parse request body
	var req dto.ApproveRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(shared.NewBadRequestError("INVALID_REQUEST", "Missing request_id field", err))
		return
	}

	// Build command
	cmd := admin.ApproveRequestCommand{
		RequestID:  req.RequestID,
		ApprovedBy: claims.Email,
		Notes:      req.Notes,
	}

	// Execute use case
	if err := h.approveRequestUC.Execute(c.Request.Context(), cmd); err != nil {
		_ = c.Error(err)
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"message": "Request approved successfully",
		},
	})
}

// RejectRequest rejects a role request (admin only)
func (h *RoleHandler) RejectRequest(c *gin.Context) {
	// Get claims from context
	claims := getClaims(c)
	if claims == nil {
		return
	}

	// Parse request body
	var req dto.RejectRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(shared.NewBadRequestError("INVALID_REQUEST", "Missing required fields (request_id, notes)", err))
		return
	}

	// Build command
	cmd := admin.RejectRequestCommand{
		RequestID:  req.RequestID,
		RejectedBy: claims.Email,
		Notes:      req.Notes,
	}

	// Execute use case
	if err := h.rejectRequestUC.Execute(c.Request.Context(), cmd); err != nil {
		_ = c.Error(err)
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"message": "Request rejected successfully",
		},
	})
}

// ListUsers lists users filtered by role (admin only)
func (h *RoleHandler) ListUsers(c *gin.Context) {
	// Get role query parameter
	roleParam := c.Query("role")
	if roleParam == "" {
		_ = c.Error(shared.NewBadRequestError("INVALID_REQUEST", "Role parameter is required", nil))
		return
	}

	// Build query
	query := admin.ListUsersByRoleQuery{
		Role: user.Role(roleParam),
	}

	// Execute use case
	dto, err := h.listUsersByRoleUC.Execute(c.Request.Context(), query)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    dto,
	})
}

// getClaims is a helper function to extract claims from context
func getClaims(c *gin.Context) *ports.TokenClaims {
	claimsInterface, exists := c.Get("claims")
	if !exists {
		_ = c.Error(shared.NewUnauthorizedError("UNAUTHORIZED", "User not authenticated", nil))
		return nil
	}

	claims, ok := claimsInterface.(*ports.TokenClaims)
	if !ok {
		_ = c.Error(shared.NewUnauthorizedError("UNAUTHORIZED", "Invalid authentication", nil))
		return nil
	}

	return claims
}
