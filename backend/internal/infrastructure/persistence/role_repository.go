package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/yuki5155/go-google-auth/internal/domain/role"
	"github.com/yuki5155/go-google-auth/internal/domain/user"
	dynamodbInfra "github.com/yuki5155/go-google-auth/internal/infrastructure/dynamodb"
)

// RoleRepository implements the role.Repository interface using DynamoDB
type RoleRepository struct {
	client           *dynamodb.Client
	userRolesTable   string
	roleRequestsTable string
}

// NewRoleRepository creates a new DynamoDB-based role repository
func NewRoleRepository(client *dynamodb.Client) *RoleRepository {
	tables := dynamodbInfra.GetTableNames()
	return &RoleRepository{
		client:            client,
		userRolesTable:    tables.UserRoles,
		roleRequestsTable: tables.RoleRequests,
	}
}

// UserRole DynamoDB schema
type userRoleModel struct {
	UserID    string `dynamodbav:"user_id"`
	Email     string `dynamodbav:"email"`
	Role      string `dynamodbav:"role"`
	Status    string `dynamodbav:"status"`
	GrantedAt int64  `dynamodbav:"granted_at"`
	GrantedBy string `dynamodbav:"granted_by"`
	CreatedAt int64  `dynamodbav:"created_at"`
	UpdatedAt int64  `dynamodbav:"updated_at"`
}

// RoleRequest DynamoDB schema
type roleRequestModel struct {
	RequestID     string  `dynamodbav:"request_id"`
	UserID        string  `dynamodbav:"user_id"`
	UserEmail     string  `dynamodbav:"user_email"`
	RequestedRole string  `dynamodbav:"requested_role"`
	Status        string  `dynamodbav:"status"`
	RequestedAt   int64   `dynamodbav:"requested_at"`
	ProcessedAt   *int64  `dynamodbav:"processed_at,omitempty"`
	ProcessedBy   string  `dynamodbav:"processed_by,omitempty"`
	Notes         string  `dynamodbav:"notes,omitempty"`
	CreatedAt     int64   `dynamodbav:"created_at"`
	UpdatedAt     int64   `dynamodbav:"updated_at"`
}

// ============================================================================
// UserRole operations
// ============================================================================

// GetUserRole retrieves a user role by user ID
func (r *RoleRepository) GetUserRole(ctx context.Context, userID string) (*role.UserRole, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.userRolesTable),
		Key: map[string]types.AttributeValue{
			"user_id": &types.AttributeValueMemberS{Value: userID},
		},
	}

	result, err := r.client.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get user role: %w", err)
	}

	if result.Item == nil {
		return nil, role.ErrUserRoleNotFound
	}

	var model userRoleModel
	if err := attributevalue.UnmarshalMap(result.Item, &model); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user role: %w", err)
	}

	return r.userRoleModelToDomain(&model), nil
}

// GetUserRoleByEmail retrieves a user role by email using GSI
func (r *RoleRepository) GetUserRoleByEmail(ctx context.Context, email string) (*role.UserRole, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.userRolesTable),
		IndexName:              aws.String("email-index"),
		KeyConditionExpression: aws.String("email = :email"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":email": &types.AttributeValueMemberS{Value: email},
		},
		Limit: aws.Int32(1),
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query user role by email: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, role.ErrUserRoleNotFound
	}

	var model userRoleModel
	if err := attributevalue.UnmarshalMap(result.Items[0], &model); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user role: %w", err)
	}

	return r.userRoleModelToDomain(&model), nil
}

// UpsertUserRole creates or updates a user role
func (r *RoleRepository) UpsertUserRole(ctx context.Context, userRole *role.UserRole) error {
	model := r.userRoleDomainToModel(userRole)

	item, err := attributevalue.MarshalMap(model)
	if err != nil {
		return fmt.Errorf("failed to marshal user role: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.userRolesTable),
		Item:      item,
	}

	_, err = r.client.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to upsert user role: %w", err)
	}

	return nil
}

// ListUsersByRole retrieves all users with a specific role using GSI
func (r *RoleRepository) ListUsersByRole(ctx context.Context, roleType user.Role) ([]*role.UserRole, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.userRolesTable),
		IndexName:              aws.String("role-index"),
		KeyConditionExpression: aws.String("#role = :role"),
		ExpressionAttributeNames: map[string]string{
			"#role": "role",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":role": &types.AttributeValueMemberS{Value: roleType.String()},
		},
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list users by role: %w", err)
	}

	userRoles := make([]*role.UserRole, 0, len(result.Items))
	for _, item := range result.Items {
		var model userRoleModel
		if err := attributevalue.UnmarshalMap(item, &model); err != nil {
			return nil, fmt.Errorf("failed to unmarshal user role: %w", err)
		}
		userRoles = append(userRoles, r.userRoleModelToDomain(&model))
	}

	return userRoles, nil
}

// ============================================================================
// RoleRequest operations
// ============================================================================

// CreateRoleRequest creates a new role request
func (r *RoleRepository) CreateRoleRequest(ctx context.Context, request *role.RoleRequest) error {
	model := r.roleRequestDomainToModel(request)

	item, err := attributevalue.MarshalMap(model)
	if err != nil {
		return fmt.Errorf("failed to marshal role request: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.roleRequestsTable),
		Item:      item,
	}

	_, err = r.client.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create role request: %w", err)
	}

	return nil
}

// GetRoleRequest retrieves a role request by request ID
func (r *RoleRepository) GetRoleRequest(ctx context.Context, requestID string) (*role.RoleRequest, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.roleRequestsTable),
		Key: map[string]types.AttributeValue{
			"request_id": &types.AttributeValueMemberS{Value: requestID},
		},
	}

	result, err := r.client.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get role request: %w", err)
	}

	if result.Item == nil {
		return nil, role.ErrRoleRequestNotFound
	}

	var model roleRequestModel
	if err := attributevalue.UnmarshalMap(result.Item, &model); err != nil {
		return nil, fmt.Errorf("failed to unmarshal role request: %w", err)
	}

	return r.roleRequestModelToDomain(&model), nil
}

// GetPendingRequestByUserID retrieves a pending role request for a user
func (r *RoleRepository) GetPendingRequestByUserID(ctx context.Context, userID string) (*role.RoleRequest, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.roleRequestsTable),
		IndexName:              aws.String("user_id-index"),
		KeyConditionExpression: aws.String("user_id = :user_id"),
		FilterExpression:       aws.String("#status = :status"),
		ExpressionAttributeNames: map[string]string{
			"#status": "status",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":user_id": &types.AttributeValueMemberS{Value: userID},
			":status":  &types.AttributeValueMemberS{Value: string(role.RequestStatusPending)},
		},
		Limit: aws.Int32(1),
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query pending request: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, nil // No pending request found
	}

	var model roleRequestModel
	if err := attributevalue.UnmarshalMap(result.Items[0], &model); err != nil {
		return nil, fmt.Errorf("failed to unmarshal role request: %w", err)
	}

	return r.roleRequestModelToDomain(&model), nil
}

// UpdateRoleRequest updates an existing role request
func (r *RoleRepository) UpdateRoleRequest(ctx context.Context, request *role.RoleRequest) error {
	model := r.roleRequestDomainToModel(request)

	item, err := attributevalue.MarshalMap(model)
	if err != nil {
		return fmt.Errorf("failed to marshal role request: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.roleRequestsTable),
		Item:      item,
	}

	_, err = r.client.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update role request: %w", err)
	}

	return nil
}

// ListRoleRequestsByStatus retrieves all role requests with a specific status
func (r *RoleRepository) ListRoleRequestsByStatus(ctx context.Context, status role.RequestStatus) ([]*role.RoleRequest, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.roleRequestsTable),
		IndexName:              aws.String("status-requested_at-index"),
		KeyConditionExpression: aws.String("#status = :status"),
		ExpressionAttributeNames: map[string]string{
			"#status": "status",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":status": &types.AttributeValueMemberS{Value: string(status)},
		},
		ScanIndexForward: aws.Bool(false), // Sort descending (newest first)
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list role requests by status: %w", err)
	}

	requests := make([]*role.RoleRequest, 0, len(result.Items))
	for _, item := range result.Items {
		var model roleRequestModel
		if err := attributevalue.UnmarshalMap(item, &model); err != nil {
			return nil, fmt.Errorf("failed to unmarshal role request: %w", err)
		}
		requests = append(requests, r.roleRequestModelToDomain(&model))
	}

	return requests, nil
}

// ListRoleRequestsByUserID retrieves all role requests for a specific user
func (r *RoleRepository) ListRoleRequestsByUserID(ctx context.Context, userID string) ([]*role.RoleRequest, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.roleRequestsTable),
		IndexName:              aws.String("user_id-index"),
		KeyConditionExpression: aws.String("user_id = :user_id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":user_id": &types.AttributeValueMemberS{Value: userID},
		},
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list role requests by user: %w", err)
	}

	requests := make([]*role.RoleRequest, 0, len(result.Items))
	for _, item := range result.Items {
		var model roleRequestModel
		if err := attributevalue.UnmarshalMap(item, &model); err != nil {
			return nil, fmt.Errorf("failed to unmarshal role request: %w", err)
		}
		requests = append(requests, r.roleRequestModelToDomain(&model))
	}

	return requests, nil
}

// ============================================================================
// Model conversion helpers
// ============================================================================

func (r *RoleRepository) userRoleModelToDomain(model *userRoleModel) *role.UserRole {
	return role.ReconstructUserRole(
		model.UserID,
		model.Email,
		user.Role(model.Role),
		role.Status(model.Status),
		time.Unix(model.GrantedAt, 0),
		model.GrantedBy,
		time.Unix(model.CreatedAt, 0),
		time.Unix(model.UpdatedAt, 0),
	)
}

func (r *RoleRepository) userRoleDomainToModel(userRole *role.UserRole) *userRoleModel {
	return &userRoleModel{
		UserID:    userRole.UserID(),
		Email:     userRole.Email(),
		Role:      userRole.Role().String(),
		Status:    string(userRole.Status()),
		GrantedAt: userRole.GrantedAt().Unix(),
		GrantedBy: userRole.GrantedBy(),
		CreatedAt: userRole.CreatedAt().Unix(),
		UpdatedAt: userRole.UpdatedAt().Unix(),
	}
}

func (r *RoleRepository) roleRequestModelToDomain(model *roleRequestModel) *role.RoleRequest {
	var processedAt *time.Time
	if model.ProcessedAt != nil {
		t := time.Unix(*model.ProcessedAt, 0)
		processedAt = &t
	}

	return role.ReconstructRoleRequest(
		model.RequestID,
		model.UserID,
		model.UserEmail,
		user.Role(model.RequestedRole),
		role.RequestStatus(model.Status),
		time.Unix(model.RequestedAt, 0),
		processedAt,
		model.ProcessedBy,
		model.Notes,
		time.Unix(model.CreatedAt, 0),
		time.Unix(model.UpdatedAt, 0),
	)
}

func (r *RoleRepository) roleRequestDomainToModel(request *role.RoleRequest) *roleRequestModel {
	var processedAt *int64
	if request.ProcessedAt() != nil {
		t := request.ProcessedAt().Unix()
		processedAt = &t
	}

	return &roleRequestModel{
		RequestID:     request.RequestID(),
		UserID:        request.UserID(),
		UserEmail:     request.UserEmail(),
		RequestedRole: request.RequestedRole().String(),
		Status:        string(request.Status()),
		RequestedAt:   request.RequestedAt().Unix(),
		ProcessedAt:   processedAt,
		ProcessedBy:   request.ProcessedBy(),
		Notes:         request.Notes(),
		CreatedAt:     request.CreatedAt().Unix(),
		UpdatedAt:     request.UpdatedAt().Unix(),
	}
}
