package admin

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/yuki5155/go-google-auth/internal/domain/role"
	"github.com/yuki5155/go-google-auth/internal/domain/user"
	"github.com/yuki5155/go-google-auth/internal/mocks"
)

// RequestRoleUseCase tests
func TestRequestRoleUseCase_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockRoleRepository(ctrl)

	cmd := RequestRoleCommand{
		UserID:        "user-123",
		UserEmail:     "test@example.com",
		RequestedRole: user.RoleSubscriber,
	}

	// Expect check for existing pending request (none found - return nil without error)
	mockRepo.EXPECT().
		GetPendingRequestByUserID(ctx, "user-123").
		Return(nil, nil)

	// Expect the repository to create a role request
	mockRepo.EXPECT().
		CreateRoleRequest(ctx, gomock.Any()).
		Return(nil)

	useCase := NewRequestRoleUseCase(mockRepo)
	result, err := useCase.Execute(ctx, cmd)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "user-123", result.UserID)
	assert.Equal(t, "test@example.com", result.UserEmail)
	assert.Equal(t, "subscriber", result.RequestedRole)
	assert.Equal(t, "pending", result.Status)
}

func TestRequestRoleUseCase_Execute_DuplicateRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockRoleRepository(ctrl)

	cmd := RequestRoleCommand{
		UserID:        "user-123",
		UserEmail:     "test@example.com",
		RequestedRole: user.RoleSubscriber,
	}

	// Mock existing pending request
	existingRequest, _ := role.NewRoleRequest("existing-req", "user-123", "test@example.com", user.RoleSubscriber)

	mockRepo.EXPECT().
		GetPendingRequestByUserID(ctx, "user-123").
		Return(existingRequest, nil)

	useCase := NewRequestRoleUseCase(mockRepo)
	result, err := useCase.Execute(ctx, cmd)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, role.ErrDuplicateRoleRequest, err)
}

// ListPendingRequestsUseCase tests
func TestListPendingRequestsUseCase_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockRoleRepository(ctrl)

	now := time.Now()
	mockRequests := []*role.RoleRequest{
		role.ReconstructRoleRequest(
			"req-1",
			"user-1",
			"user1@example.com",
			user.RoleSubscriber,
			role.RequestStatusPending,
			now,
			nil,
			"",
			"",
			now,
			now,
		),
	}

	mockRepo.EXPECT().
		ListRoleRequestsByStatus(ctx, role.RequestStatusPending).
		Return(mockRequests, nil)

	useCase := NewListPendingRequestsUseCase(mockRepo)
	result, err := useCase.Execute(ctx)

	require.NoError(t, err)
	assert.Len(t, result.Requests, 1)
	assert.Equal(t, "req-1", result.Requests[0].RequestID)
	assert.Equal(t, "user-1", result.Requests[0].UserID)
}

// ApproveRequestUseCase tests
func TestApproveRequestUseCase_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockRoleRepository(ctrl)

	cmd := ApproveRequestCommand{
		RequestID:  "req-123",
		ApprovedBy: "admin@example.com",
		Notes:      "Approved",
	}

	// Mock get request
	mockRequest, _ := role.NewRoleRequest("req-123", "user-123", "test@example.com", user.RoleSubscriber)
	mockRepo.EXPECT().
		GetRoleRequest(ctx, "req-123").
		Return(mockRequest, nil)

	// Expect update request
	mockRepo.EXPECT().
		UpdateRoleRequest(ctx, gomock.Any()).
		Return(nil)

	// Expect upsert user role
	mockRepo.EXPECT().
		UpsertUserRole(ctx, gomock.Any()).
		Return(nil)

	useCase := NewApproveRequestUseCase(mockRepo)
	err := useCase.Execute(ctx, cmd)

	require.NoError(t, err)
}

// RejectRequestUseCase tests
func TestRejectRequestUseCase_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockRoleRepository(ctrl)

	cmd := RejectRequestCommand{
		RequestID:  "req-123",
		RejectedBy: "admin@example.com",
		Notes:      "Insufficient justification",
	}

	// Mock get request
	mockRequest, _ := role.NewRoleRequest("req-123", "user-123", "test@example.com", user.RoleSubscriber)
	mockRepo.EXPECT().
		GetRoleRequest(ctx, "req-123").
		Return(mockRequest, nil)

	// Expect update request
	mockRepo.EXPECT().
		UpdateRoleRequest(ctx, gomock.Any()).
		Return(nil)

	useCase := NewRejectRequestUseCase(mockRepo)
	err := useCase.Execute(ctx, cmd)

	require.NoError(t, err)
}

// CheckUserRoleUseCase tests
func TestCheckUserRoleUseCase_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockRoleRepository(ctrl)

	mockUserRole, _ := role.NewUserRole("user-123", "test@example.com", user.RoleSubscriber, "admin@example.com")
	mockRepo.EXPECT().
		GetUserRole(ctx, "user-123").
		Return(mockUserRole, nil)

	useCase := NewCheckUserRoleUseCase(mockRepo)
	result, err := useCase.Execute(ctx, "user-123")

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "user-123", result.UserID())
	assert.Equal(t, user.RoleSubscriber, result.Role())
}

// ListUsersByRoleUseCase tests
func TestListUsersByRoleUseCase_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockRoleRepository(ctrl)

	query := ListUsersByRoleQuery{
		Role: user.RoleSubscriber,
	}

	mockUserRoles := []*role.UserRole{
		func() *role.UserRole {
			ur, _ := role.NewUserRole("user-1", "user1@example.com", user.RoleSubscriber, "admin@example.com")
			return ur
		}(),
	}

	mockRepo.EXPECT().
		ListUsersByRole(ctx, user.RoleSubscriber).
		Return(mockUserRoles, nil)

	useCase := NewListUsersByRoleUseCase(mockRepo)
	result, err := useCase.Execute(ctx, query)

	require.NoError(t, err)
	assert.Len(t, result.Users, 1)
	assert.Equal(t, "user-1", result.Users[0].UserID)
}

// Error cases
func TestRequestRoleUseCase_Execute_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockRoleRepository(ctrl)

	cmd := RequestRoleCommand{
		UserID:        "user-123",
		UserEmail:     "test@example.com",
		RequestedRole: user.RoleSubscriber,
	}

	mockRepo.EXPECT().
		GetPendingRequestByUserID(ctx, "user-123").
		Return(nil, errors.New("database error"))

	useCase := NewRequestRoleUseCase(mockRepo)
	result, err := useCase.Execute(ctx, cmd)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestCheckUserRoleUseCase_Execute_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockRoleRepository(ctrl)

	mockRepo.EXPECT().
		GetUserRole(ctx, "user-123").
		Return(nil, role.ErrUserRoleNotFound)

	useCase := NewCheckUserRoleUseCase(mockRepo)
	result, err := useCase.Execute(ctx, "user-123")

	// When user role not found, Execute returns (nil, nil) - user has default role
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestCheckUserRoleUseCase_Execute_DatabaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockRoleRepository(ctrl)

	mockRepo.EXPECT().
		GetUserRole(ctx, "user-123").
		Return(nil, errors.New("database error"))

	useCase := NewCheckUserRoleUseCase(mockRepo)
	result, err := useCase.Execute(ctx, "user-123")

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestApproveRequestUseCase_Execute_GetRequestError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockRoleRepository(ctrl)

	cmd := ApproveRequestCommand{
		RequestID:  "req-123",
		ApprovedBy: "admin@example.com",
		Notes:      "Approved",
	}

	mockRepo.EXPECT().
		GetRoleRequest(ctx, "req-123").
		Return(nil, errors.New("database error"))

	useCase := NewApproveRequestUseCase(mockRepo)
	err := useCase.Execute(ctx, cmd)

	require.Error(t, err)
}

func TestApproveRequestUseCase_Execute_UpdateRequestError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockRoleRepository(ctrl)

	cmd := ApproveRequestCommand{
		RequestID:  "req-123",
		ApprovedBy: "admin@example.com",
		Notes:      "Approved",
	}

	mockRequest, _ := role.NewRoleRequest("req-123", "user-123", "test@example.com", user.RoleSubscriber)
	mockRepo.EXPECT().
		GetRoleRequest(ctx, "req-123").
		Return(mockRequest, nil)

	mockRepo.EXPECT().
		UpdateRoleRequest(ctx, gomock.Any()).
		Return(errors.New("update error"))

	useCase := NewApproveRequestUseCase(mockRepo)
	err := useCase.Execute(ctx, cmd)

	require.Error(t, err)
}

func TestApproveRequestUseCase_Execute_UpsertUserRoleError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockRoleRepository(ctrl)

	cmd := ApproveRequestCommand{
		RequestID:  "req-123",
		ApprovedBy: "admin@example.com",
		Notes:      "Approved",
	}

	mockRequest, _ := role.NewRoleRequest("req-123", "user-123", "test@example.com", user.RoleSubscriber)
	mockRepo.EXPECT().
		GetRoleRequest(ctx, "req-123").
		Return(mockRequest, nil)

	mockRepo.EXPECT().
		UpdateRoleRequest(ctx, gomock.Any()).
		Return(nil)

	mockRepo.EXPECT().
		UpsertUserRole(ctx, gomock.Any()).
		Return(errors.New("upsert error"))

	useCase := NewApproveRequestUseCase(mockRepo)
	err := useCase.Execute(ctx, cmd)

	require.Error(t, err)
}

func TestRejectRequestUseCase_Execute_GetRequestError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockRoleRepository(ctrl)

	cmd := RejectRequestCommand{
		RequestID:  "req-123",
		RejectedBy: "admin@example.com",
		Notes:      "Rejected",
	}

	mockRepo.EXPECT().
		GetRoleRequest(ctx, "req-123").
		Return(nil, errors.New("database error"))

	useCase := NewRejectRequestUseCase(mockRepo)
	err := useCase.Execute(ctx, cmd)

	require.Error(t, err)
}

func TestRejectRequestUseCase_Execute_UpdateRequestError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockRoleRepository(ctrl)

	cmd := RejectRequestCommand{
		RequestID:  "req-123",
		RejectedBy: "admin@example.com",
		Notes:      "Rejected",
	}

	mockRequest, _ := role.NewRoleRequest("req-123", "user-123", "test@example.com", user.RoleSubscriber)
	mockRepo.EXPECT().
		GetRoleRequest(ctx, "req-123").
		Return(mockRequest, nil)

	mockRepo.EXPECT().
		UpdateRoleRequest(ctx, gomock.Any()).
		Return(errors.New("update error"))

	useCase := NewRejectRequestUseCase(mockRepo)
	err := useCase.Execute(ctx, cmd)

	require.Error(t, err)
}

func TestListPendingRequestsUseCase_Execute_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockRoleRepository(ctrl)

	mockRepo.EXPECT().
		ListRoleRequestsByStatus(ctx, role.RequestStatusPending).
		Return(nil, errors.New("database error"))

	useCase := NewListPendingRequestsUseCase(mockRepo)
	result, err := useCase.Execute(ctx)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestListUsersByRoleUseCase_Execute_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockRoleRepository(ctrl)

	query := ListUsersByRoleQuery{
		Role: user.RoleSubscriber,
	}

	mockRepo.EXPECT().
		ListUsersByRole(ctx, user.RoleSubscriber).
		Return(nil, errors.New("database error"))

	useCase := NewListUsersByRoleUseCase(mockRepo)
	result, err := useCase.Execute(ctx, query)

	require.Error(t, err)
	assert.Nil(t, result)
}
