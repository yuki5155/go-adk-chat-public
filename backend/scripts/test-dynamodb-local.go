//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

const (
	endpoint    = "http://localhost:8000"
	region      = "ap-northeast-1"
	projectName = "go-adk-chat"
	environment = "dev"
)

type UserRole struct {
	UserID    string `dynamodbav:"user_id"`
	Email     string `dynamodbav:"email"`
	Role      string `dynamodbav:"role"`
	Status    string `dynamodbav:"status"`
	GrantedAt int64  `dynamodbav:"granted_at"`
	GrantedBy string `dynamodbav:"granted_by"`
	CreatedAt int64  `dynamodbav:"created_at"`
	UpdatedAt int64  `dynamodbav:"updated_at"`
}

type RoleRequest struct {
	RequestID     string `dynamodbav:"request_id"`
	UserID        string `dynamodbav:"user_id"`
	UserEmail     string `dynamodbav:"user_email"`
	RequestedRole string `dynamodbav:"requested_role"`
	Status        string `dynamodbav:"status"`
	RequestedAt   int64  `dynamodbav:"requested_at"`
	ProcessedAt   *int64 `dynamodbav:"processed_at,omitempty"`
	ProcessedBy   string `dynamodbav:"processed_by,omitempty"`
	Notes         string `dynamodbav:"notes,omitempty"`
	CreatedAt     int64  `dynamodbav:"created_at"`
	UpdatedAt     int64  `dynamodbav:"updated_at"`
}

func main() {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           endpoint,
					SigningRegion: region,
				}, nil
			}),
		),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummy", "dummy", "")),
	)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	client := dynamodb.NewFromConfig(cfg)

	fmt.Println("=== Testing DynamoDB Local ===")
	fmt.Printf("Endpoint: %s\n\n", endpoint)

	// Test 1: Insert a user role
	fmt.Println("Test 1: Inserting a user role...")
	if err := insertUserRole(ctx, client); err != nil {
		log.Printf("❌ Failed: %v\n", err)
	} else {
		fmt.Println("✅ Success: User role inserted")
	}

	// Test 2: Query user role by email
	fmt.Println("\nTest 2: Querying user role by email...")
	if err := queryUserRoleByEmail(ctx, client, "test.user@example.com"); err != nil {
		log.Printf("❌ Failed: %v\n", err)
	} else {
		fmt.Println("✅ Success: User role queried by email")
	}

	// Test 3: Insert a role request
	fmt.Println("\nTest 3: Inserting a role request...")
	requestID := uuid.New().String()
	if err := insertRoleRequest(ctx, client, requestID); err != nil {
		log.Printf("❌ Failed: %v\n", err)
	} else {
		fmt.Println("✅ Success: Role request inserted")
	}

	// Test 4: Query pending role requests
	fmt.Println("\nTest 4: Querying pending role requests...")
	if err := queryPendingRoleRequests(ctx, client); err != nil {
		log.Printf("❌ Failed: %v\n", err)
	} else {
		fmt.Println("✅ Success: Pending role requests queried")
	}

	fmt.Println("\n=== All Tests Completed ===")
}

func insertUserRole(ctx context.Context, client *dynamodb.Client) error {
	tableName := fmt.Sprintf("%s-%s-user-roles", projectName, environment)
	now := time.Now().Unix()

	userRole := UserRole{
		UserID:    "user-123",
		Email:     "test.user@example.com",
		Role:      "subscriber",
		Status:    "active",
		GrantedAt: now,
		GrantedBy: "admin@example.com",
		CreatedAt: now,
		UpdatedAt: now,
	}

	item, err := attributevalue.MarshalMap(userRole)
	if err != nil {
		return fmt.Errorf("failed to marshal user role: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	}

	_, err = client.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put item: %w", err)
	}

	fmt.Printf("  → Inserted user_id: %s, email: %s, role: %s\n",
		userRole.UserID, userRole.Email, userRole.Role)
	return nil
}

func queryUserRoleByEmail(ctx context.Context, client *dynamodb.Client, email string) error {
	tableName := fmt.Sprintf("%s-%s-user-roles", projectName, environment)

	input := &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		IndexName:              aws.String("email-index"),
		KeyConditionExpression: aws.String("email = :email"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":email": &types.AttributeValueMemberS{Value: email},
		},
	}

	result, err := client.Query(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to query: %w", err)
	}

	fmt.Printf("  → Found %d user(s) with email: %s\n", len(result.Items), email)

	for _, item := range result.Items {
		var userRole UserRole
		if err := attributevalue.UnmarshalMap(item, &userRole); err != nil {
			return fmt.Errorf("failed to unmarshal: %w", err)
			}
		fmt.Printf("     - user_id: %s, role: %s, status: %s\n",
			userRole.UserID, userRole.Role, userRole.Status)
	}

	return nil
}

func insertRoleRequest(ctx context.Context, client *dynamodb.Client, requestID string) error {
	tableName := fmt.Sprintf("%s-%s-role-requests", projectName, environment)
	now := time.Now().Unix()

	roleRequest := RoleRequest{
		RequestID:     requestID,
		UserID:        "user-456",
		UserEmail:     "new.user@example.com",
		RequestedRole: "subscriber",
		Status:        "pending",
		RequestedAt:   now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	item, err := attributevalue.MarshalMap(roleRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal role request: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	}

	_, err = client.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put item: %w", err)
	}

	fmt.Printf("  → Inserted request_id: %s, user_id: %s, status: %s\n",
		roleRequest.RequestID, roleRequest.UserID, roleRequest.Status)
	return nil
}

func queryPendingRoleRequests(ctx context.Context, client *dynamodb.Client) error {
	tableName := fmt.Sprintf("%s-%s-role-requests", projectName, environment)

	input := &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		IndexName:              aws.String("status-requested_at-index"),
		KeyConditionExpression: aws.String("#status = :status"),
		ExpressionAttributeNames: map[string]string{
			"#status": "status",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":status": &types.AttributeValueMemberS{Value: "pending"},
		},
		ScanIndexForward: aws.Bool(false), // Sort descending (newest first)
	}

	result, err := client.Query(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to query: %w", err)
	}

	fmt.Printf("  → Found %d pending request(s)\n", len(result.Items))

	for _, item := range result.Items {
		var roleRequest RoleRequest
		if err := attributevalue.UnmarshalMap(item, &roleRequest); err != nil {
			return fmt.Errorf("failed to unmarshal: %w", err)
		}
		fmt.Printf("     - request_id: %s, user_email: %s, requested_role: %s\n",
			roleRequest.RequestID, roleRequest.UserEmail, roleRequest.RequestedRole)
	}

	return nil
}
