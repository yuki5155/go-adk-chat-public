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

	fmt.Println("=== Comprehensive DynamoDB Operations Verification ===")
	fmt.Printf("Endpoint: %s\n\n", endpoint)

	// Test user role operations
	testUserRoleOperations(ctx, client)

	// Test role request operations
	testRoleRequestOperations(ctx, client)

	fmt.Println("\n=== All Verification Tests Completed Successfully! ===")
}

func testUserRoleOperations(ctx context.Context, client *dynamodb.Client) {
	fmt.Println("--- User Role Operations ---")

	testUserID := fmt.Sprintf("verify-user-%s", uuid.New().String()[:8])
	testEmail := fmt.Sprintf("%s@example.com", testUserID)

	// Test 1: Create (PutItem)
	fmt.Println("\n1. Testing CREATE operation...")
	if err := createUserRole(ctx, client, testUserID, testEmail); err != nil {
		log.Fatalf("❌ Create failed: %v", err)
	}
	fmt.Println("✅ Create successful")

	// Test 2: Read (GetItem)
	fmt.Println("\n2. Testing READ operation...")
	userRole, err := getUserRole(ctx, client, testUserID)
	if err != nil {
		log.Fatalf("❌ Read failed: %v", err)
	}
	fmt.Printf("✅ Read successful: user_id=%s, role=%s, status=%s\n",
		userRole.UserID, userRole.Role, userRole.Status)

	// Test 3: Update (UpdateItem)
	fmt.Println("\n3. Testing UPDATE operation...")
	if err := updateUserRole(ctx, client, testUserID); err != nil {
		log.Fatalf("❌ Update failed: %v", err)
	}
	fmt.Println("✅ Update successful")

	// Verify update
	userRole, err = getUserRole(ctx, client, testUserID)
	if err != nil {
		log.Fatalf("❌ Read after update failed: %v", err)
	}
	fmt.Printf("   → Updated role: %s, status: %s\n", userRole.Role, userRole.Status)

	// Test 4: Query by email (GSI)
	fmt.Println("\n4. Testing QUERY by email (using GSI)...")
	if err := queryUserRoleByEmail(ctx, client, testEmail); err != nil {
		log.Fatalf("❌ Query by email failed: %v", err)
	}
	fmt.Println("✅ Query by email successful")

	// Test 5: Query by role (GSI)
	fmt.Println("\n5. Testing QUERY by role (using GSI)...")
	if err := queryUserRolesByRole(ctx, client, "admin"); err != nil {
		log.Fatalf("❌ Query by role failed: %v", err)
	}
	fmt.Println("✅ Query by role successful")

	// Test 6: Delete (DeleteItem)
	fmt.Println("\n6. Testing DELETE operation...")
	if err := deleteUserRole(ctx, client, testUserID); err != nil {
		log.Fatalf("❌ Delete failed: %v", err)
	}
	fmt.Println("✅ Delete successful")

	// Verify deletion
	_, err = getUserRole(ctx, client, testUserID)
	if err == nil {
		log.Fatalf("❌ User role still exists after delete")
	}
	fmt.Println("   → Verified: User role deleted")
}

func testRoleRequestOperations(ctx context.Context, client *dynamodb.Client) {
	fmt.Println("\n--- Role Request Operations ---")

	testRequestID := uuid.New().String()
	testUserID := fmt.Sprintf("verify-user-%s", uuid.New().String()[:8])
	testEmail := fmt.Sprintf("%s@example.com", testUserID)

	// Test 1: Create role request
	fmt.Println("\n1. Testing CREATE role request...")
	if err := createRoleRequest(ctx, client, testRequestID, testUserID, testEmail); err != nil {
		log.Fatalf("❌ Create failed: %v", err)
	}
	fmt.Println("✅ Create successful")

	// Test 2: Read role request
	fmt.Println("\n2. Testing READ role request...")
	roleRequest, err := getRoleRequest(ctx, client, testRequestID)
	if err != nil {
		log.Fatalf("❌ Read failed: %v", err)
	}
	fmt.Printf("✅ Read successful: request_id=%s, status=%s\n",
		roleRequest.RequestID, roleRequest.Status)

	// Test 3: Query pending requests (GSI)
	fmt.Println("\n3. Testing QUERY pending requests (using GSI)...")
	if err := queryPendingRequests(ctx, client); err != nil {
		log.Fatalf("❌ Query pending failed: %v", err)
	}
	fmt.Println("✅ Query pending successful")

	// Test 4: Update role request (approve)
	fmt.Println("\n4. Testing UPDATE role request (approve)...")
	if err := updateRoleRequest(ctx, client, testRequestID); err != nil {
		log.Fatalf("❌ Update failed: %v", err)
	}
	fmt.Println("✅ Update successful")

	// Verify update
	roleRequest, err = getRoleRequest(ctx, client, testRequestID)
	if err != nil {
		log.Fatalf("❌ Read after update failed: %v", err)
	}
	fmt.Printf("   → Updated status: %s, processed_by: %s\n",
		roleRequest.Status, roleRequest.ProcessedBy)

	// Test 5: Query by user_id (GSI)
	fmt.Println("\n5. Testing QUERY by user_id (using GSI)...")
	if err := queryRoleRequestsByUserID(ctx, client, testUserID); err != nil {
		log.Fatalf("❌ Query by user_id failed: %v", err)
	}
	fmt.Println("✅ Query by user_id successful")

	// Test 6: Delete role request
	fmt.Println("\n6. Testing DELETE role request...")
	if err := deleteRoleRequest(ctx, client, testRequestID); err != nil {
		log.Fatalf("❌ Delete failed: %v", err)
	}
	fmt.Println("✅ Delete successful")
}

// User Role operations
func createUserRole(ctx context.Context, client *dynamodb.Client, userID, email string) error {
	tableName := fmt.Sprintf("%s-%s-user-roles", projectName, environment)
	now := time.Now().Unix()

	userRole := UserRole{
		UserID:    userID,
		Email:     email,
		Role:      "subscriber",
		Status:    "active",
		GrantedAt: now,
		GrantedBy: "admin@example.com",
		CreatedAt: now,
		UpdatedAt: now,
	}

	item, err := attributevalue.MarshalMap(userRole)
	if err != nil {
		return err
	}

	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	return err
}

func getUserRole(ctx context.Context, client *dynamodb.Client, userID string) (*UserRole, error) {
	tableName := fmt.Sprintf("%s-%s-user-roles", projectName, environment)

	result, err := client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"user_id": &types.AttributeValueMemberS{Value: userID},
		},
	})
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, fmt.Errorf("user role not found")
	}

	var userRole UserRole
	if err := attributevalue.UnmarshalMap(result.Item, &userRole); err != nil {
		return nil, err
	}

	return &userRole, nil
}

func updateUserRole(ctx context.Context, client *dynamodb.Client, userID string) error {
	tableName := fmt.Sprintf("%s-%s-user-roles", projectName, environment)

	_, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"user_id": &types.AttributeValueMemberS{Value: userID},
		},
		UpdateExpression: aws.String("SET #role = :role, #status = :status, #updated_at = :updated_at"),
		ExpressionAttributeNames: map[string]string{
			"#role":       "role",
			"#status":     "status",
			"#updated_at": "updated_at",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":role":       &types.AttributeValueMemberS{Value: "admin"},
			":status":     &types.AttributeValueMemberS{Value: "active"},
			":updated_at": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", time.Now().Unix())},
		},
	})
	return err
}

func queryUserRoleByEmail(ctx context.Context, client *dynamodb.Client, email string) error {
	tableName := fmt.Sprintf("%s-%s-user-roles", projectName, environment)

	result, err := client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		IndexName:              aws.String("email-index"),
		KeyConditionExpression: aws.String("email = :email"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":email": &types.AttributeValueMemberS{Value: email},
		},
	})
	if err != nil {
		return err
	}

	fmt.Printf("   → Found %d user(s) with email: %s\n", len(result.Items), email)
	return nil
}

func queryUserRolesByRole(ctx context.Context, client *dynamodb.Client, role string) error {
	tableName := fmt.Sprintf("%s-%s-user-roles", projectName, environment)

	result, err := client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		IndexName:              aws.String("role-index"),
		KeyConditionExpression: aws.String("#role = :role"),
		ExpressionAttributeNames: map[string]string{
			"#role": "role",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":role": &types.AttributeValueMemberS{Value: role},
		},
	})
	if err != nil {
		return err
	}

	fmt.Printf("   → Found %d user(s) with role: %s\n", len(result.Items), role)
	return nil
}

func deleteUserRole(ctx context.Context, client *dynamodb.Client, userID string) error {
	tableName := fmt.Sprintf("%s-%s-user-roles", projectName, environment)

	_, err := client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"user_id": &types.AttributeValueMemberS{Value: userID},
		},
	})
	return err
}

// Role Request operations
func createRoleRequest(ctx context.Context, client *dynamodb.Client, requestID, userID, email string) error {
	tableName := fmt.Sprintf("%s-%s-role-requests", projectName, environment)
	now := time.Now().Unix()

	roleRequest := RoleRequest{
		RequestID:     requestID,
		UserID:        userID,
		UserEmail:     email,
		RequestedRole: "subscriber",
		Status:        "pending",
		RequestedAt:   now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	item, err := attributevalue.MarshalMap(roleRequest)
	if err != nil {
		return err
	}

	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	return err
}

func getRoleRequest(ctx context.Context, client *dynamodb.Client, requestID string) (*RoleRequest, error) {
	tableName := fmt.Sprintf("%s-%s-role-requests", projectName, environment)

	result, err := client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"request_id": &types.AttributeValueMemberS{Value: requestID},
		},
	})
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, fmt.Errorf("role request not found")
	}

	var roleRequest RoleRequest
	if err := attributevalue.UnmarshalMap(result.Item, &roleRequest); err != nil {
		return nil, err
	}

	return &roleRequest, nil
}

func updateRoleRequest(ctx context.Context, client *dynamodb.Client, requestID string) error {
	tableName := fmt.Sprintf("%s-%s-role-requests", projectName, environment)
	processedAt := time.Now().Unix()

	_, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"request_id": &types.AttributeValueMemberS{Value: requestID},
		},
		UpdateExpression: aws.String("SET #status = :status, #processed_at = :processed_at, #processed_by = :processed_by, #updated_at = :updated_at"),
		ExpressionAttributeNames: map[string]string{
			"#status":       "status",
			"#processed_at": "processed_at",
			"#processed_by": "processed_by",
			"#updated_at":   "updated_at",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":status":       &types.AttributeValueMemberS{Value: "approved"},
			":processed_at": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", processedAt)},
			":processed_by": &types.AttributeValueMemberS{Value: "admin@example.com"},
			":updated_at":   &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", time.Now().Unix())},
		},
	})
	return err
}

func queryPendingRequests(ctx context.Context, client *dynamodb.Client) error {
	tableName := fmt.Sprintf("%s-%s-role-requests", projectName, environment)

	result, err := client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		IndexName:              aws.String("status-requested_at-index"),
		KeyConditionExpression: aws.String("#status = :status"),
		ExpressionAttributeNames: map[string]string{
			"#status": "status",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":status": &types.AttributeValueMemberS{Value: "pending"},
		},
	})
	if err != nil {
		return err
	}

	fmt.Printf("   → Found %d pending request(s)\n", len(result.Items))
	return nil
}

func queryRoleRequestsByUserID(ctx context.Context, client *dynamodb.Client, userID string) error {
	tableName := fmt.Sprintf("%s-%s-role-requests", projectName, environment)

	result, err := client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		IndexName:              aws.String("user_id-index"),
		KeyConditionExpression: aws.String("user_id = :user_id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":user_id": &types.AttributeValueMemberS{Value: userID},
		},
	})
	if err != nil {
		return err
	}

	fmt.Printf("   → Found %d request(s) for user_id: %s\n", len(result.Items), userID)
	return nil
}

func deleteRoleRequest(ctx context.Context, client *dynamodb.Client, requestID string) error {
	tableName := fmt.Sprintf("%s-%s-role-requests", projectName, environment)

	_, err := client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"request_id": &types.AttributeValueMemberS{Value: requestID},
		},
	})
	return err
}
