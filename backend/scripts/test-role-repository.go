//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
	"github.com/yuki5155/go-google-auth/internal/domain/role"
	"github.com/yuki5155/go-google-auth/internal/domain/user"
	"github.com/yuki5155/go-google-auth/internal/infrastructure/persistence"
)

const (
	endpoint = "http://localhost:8000"
	region   = "ap-northeast-1"
)

func main() {
	ctx := context.Background()

	// Set environment variables for table names
	os.Setenv("PROJECT_NAME", "go-adk-chat")
	os.Setenv("GO_ENV", "dev")

	// Configure AWS SDK for DynamoDB Local
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
	repo := persistence.NewRoleRepository(client)

	fmt.Println("=== Testing Role Repository Implementation ===")
	fmt.Printf("Endpoint: %s\n\n", endpoint)

	// Test UserRole operations
	testUserRoleOperations(ctx, repo)

	// Test RoleRequest operations
	testRoleRequestOperations(ctx, repo)

	fmt.Println("\n=== All Repository Tests Completed Successfully! ===")
}

func testUserRoleOperations(ctx context.Context, repo *persistence.RoleRepository) {
	fmt.Println("--- UserRole Repository Operations ---")

	testUserID := fmt.Sprintf("test-user-%s", uuid.New().String()[:8])
	testEmail := fmt.Sprintf("%s@example.com", testUserID)

	// Test 1: Create a new UserRole
	fmt.Println("\n1. Testing UpsertUserRole (Create)...")
	userRole, err := role.NewUserRole(testUserID, testEmail, user.RoleSubscriber, "admin@example.com")
	if err != nil {
		log.Fatalf("❌ Failed to create UserRole entity: %v", err)
	}

	if err := repo.UpsertUserRole(ctx, userRole); err != nil {
		log.Fatalf("❌ Failed to upsert user role: %v", err)
	}
	fmt.Printf("✅ UserRole created: user_id=%s, email=%s, role=%s\n",
		userRole.UserID(), userRole.Email(), userRole.Role())

	// Test 2: GetUserRole by ID
	fmt.Println("\n2. Testing GetUserRole...")
	retrieved, err := repo.GetUserRole(ctx, testUserID)
	if err != nil {
		log.Fatalf("❌ Failed to get user role: %v", err)
	}
	fmt.Printf("✅ UserRole retrieved: user_id=%s, role=%s, status=%s\n",
		retrieved.UserID(), retrieved.Role(), retrieved.Status())

	// Test 3: GetUserRoleByEmail
	fmt.Println("\n3. Testing GetUserRoleByEmail...")
	byEmail, err := repo.GetUserRoleByEmail(ctx, testEmail)
	if err != nil {
		log.Fatalf("❌ Failed to get user role by email: %v", err)
	}
	fmt.Printf("✅ UserRole retrieved by email: user_id=%s, email=%s\n",
		byEmail.UserID(), byEmail.Email())

	// Test 4: Update UserRole (change role)
	fmt.Println("\n4. Testing UpsertUserRole (Update - change role)...")
	if err := userRole.ChangeRole(user.RoleAdmin, "root@example.com"); err != nil {
		log.Fatalf("❌ Failed to change role: %v", err)
	}
	if err := repo.UpsertUserRole(ctx, userRole); err != nil {
		log.Fatalf("❌ Failed to update user role: %v", err)
	}
	fmt.Println("✅ UserRole updated: role changed to admin")

	// Verify update
	updated, err := repo.GetUserRole(ctx, testUserID)
	if err != nil {
		log.Fatalf("❌ Failed to verify update: %v", err)
	}
	fmt.Printf("   → Verified: role=%s, granted_by=%s\n",
		updated.Role(), updated.GrantedBy())

	// Test 5: ListUsersByRole
	fmt.Println("\n5. Testing ListUsersByRole...")
	admins, err := repo.ListUsersByRole(ctx, user.RoleAdmin)
	if err != nil {
		log.Fatalf("❌ Failed to list users by role: %v", err)
	}
	fmt.Printf("✅ Found %d admin user(s)\n", len(admins))
	for _, admin := range admins {
		fmt.Printf("   → user_id=%s, email=%s, role=%s\n",
			admin.UserID(), admin.Email(), admin.Role())
	}

	// Test 6: Test not found scenario
	fmt.Println("\n6. Testing GetUserRole (not found)...")
	_, err = repo.GetUserRole(ctx, "non-existent-user")
	if err == role.ErrUserRoleNotFound {
		fmt.Println("✅ Correctly returned ErrUserRoleNotFound")
	} else {
		log.Fatalf("❌ Expected ErrUserRoleNotFound, got: %v", err)
	}
}

func testRoleRequestOperations(ctx context.Context, repo *persistence.RoleRepository) {
	fmt.Println("\n--- RoleRequest Repository Operations ---")

	testRequestID := uuid.New().String()
	testUserID := fmt.Sprintf("test-user-%s", uuid.New().String()[:8])
	testEmail := fmt.Sprintf("%s@example.com", testUserID)

	// Test 1: Create a new RoleRequest
	fmt.Println("\n1. Testing CreateRoleRequest...")
	roleRequest, err := role.NewRoleRequest(testRequestID, testUserID, testEmail, user.RoleSubscriber)
	if err != nil {
		log.Fatalf("❌ Failed to create RoleRequest entity: %v", err)
	}

	if err := repo.CreateRoleRequest(ctx, roleRequest); err != nil {
		log.Fatalf("❌ Failed to create role request: %v", err)
	}
	fmt.Printf("✅ RoleRequest created: request_id=%s, user_id=%s, status=%s\n",
		roleRequest.RequestID(), roleRequest.UserID(), roleRequest.Status())

	// Test 2: GetRoleRequest by ID
	fmt.Println("\n2. Testing GetRoleRequest...")
	retrieved, err := repo.GetRoleRequest(ctx, testRequestID)
	if err != nil {
		log.Fatalf("❌ Failed to get role request: %v", err)
	}
	fmt.Printf("✅ RoleRequest retrieved: request_id=%s, status=%s\n",
		retrieved.RequestID(), retrieved.Status())

	// Test 3: GetPendingRequestByUserID
	fmt.Println("\n3. Testing GetPendingRequestByUserID...")
	pending, err := repo.GetPendingRequestByUserID(ctx, testUserID)
	if err != nil {
		log.Fatalf("❌ Failed to get pending request: %v", err)
	}
	if pending != nil {
		fmt.Printf("✅ Found pending request: request_id=%s, user_id=%s\n",
			pending.RequestID(), pending.UserID())
	} else {
		log.Fatalf("❌ Expected to find pending request, got nil")
	}

	// Test 4: ListRoleRequestsByStatus (pending)
	fmt.Println("\n4. Testing ListRoleRequestsByStatus (pending)...")
	pendingRequests, err := repo.ListRoleRequestsByStatus(ctx, role.RequestStatusPending)
	if err != nil {
		log.Fatalf("❌ Failed to list pending requests: %v", err)
	}
	fmt.Printf("✅ Found %d pending request(s)\n", len(pendingRequests))
	for _, req := range pendingRequests {
		fmt.Printf("   → request_id=%s, user_email=%s, requested_role=%s\n",
			req.RequestID(), req.UserEmail(), req.RequestedRole())
	}

	// Test 5: ListRoleRequestsByUserID
	fmt.Println("\n5. Testing ListRoleRequestsByUserID...")
	userRequests, err := repo.ListRoleRequestsByUserID(ctx, testUserID)
	if err != nil {
		log.Fatalf("❌ Failed to list requests by user: %v", err)
	}
	fmt.Printf("✅ Found %d request(s) for user_id=%s\n", len(userRequests), testUserID)

	// Test 6: Update RoleRequest (approve)
	fmt.Println("\n6. Testing UpdateRoleRequest (approve)...")
	if err := roleRequest.Approve("admin@example.com", "Approved for testing"); err != nil {
		log.Fatalf("❌ Failed to approve request: %v", err)
	}
	if err := repo.UpdateRoleRequest(ctx, roleRequest); err != nil {
		log.Fatalf("❌ Failed to update role request: %v", err)
	}
	fmt.Println("✅ RoleRequest approved and updated")

	// Verify update
	updated, err := repo.GetRoleRequest(ctx, testRequestID)
	if err != nil {
		log.Fatalf("❌ Failed to verify update: %v", err)
	}
	fmt.Printf("   → Verified: status=%s, processed_by=%s, notes=%s\n",
		updated.Status(), updated.ProcessedBy(), updated.Notes())

	// Test 7: Test duplicate approval prevention
	fmt.Println("\n7. Testing duplicate approval prevention...")
	err = roleRequest.Approve("another-admin@example.com", "Trying to approve again")
	if err == role.ErrRequestAlreadyProcessed {
		fmt.Println("✅ Correctly prevented duplicate approval")
	} else {
		log.Fatalf("❌ Expected ErrRequestAlreadyProcessed, got: %v", err)
	}

	// Test 8: Create and reject a request
	fmt.Println("\n8. Testing RoleRequest rejection flow...")
	rejectRequestID := uuid.New().String()
	rejectUserID := fmt.Sprintf("reject-user-%s", uuid.New().String()[:8])
	rejectEmail := fmt.Sprintf("%s@example.com", rejectUserID)

	rejectRequest, err := role.NewRoleRequest(rejectRequestID, rejectUserID, rejectEmail, user.RoleSubscriber)
	if err != nil {
		log.Fatalf("❌ Failed to create request for rejection test: %v", err)
	}
	if err := repo.CreateRoleRequest(ctx, rejectRequest); err != nil {
		log.Fatalf("❌ Failed to create request: %v", err)
	}

	if err := rejectRequest.Reject("admin@example.com", "Not eligible at this time"); err != nil {
		log.Fatalf("❌ Failed to reject request: %v", err)
	}
	if err := repo.UpdateRoleRequest(ctx, rejectRequest); err != nil {
		log.Fatalf("❌ Failed to update rejected request: %v", err)
	}
	fmt.Printf("✅ RoleRequest rejected: status=%s\n", rejectRequest.Status())

	// Test 9: Test not found scenario
	fmt.Println("\n9. Testing GetRoleRequest (not found)...")
	_, err = repo.GetRoleRequest(ctx, "non-existent-request-id")
	if err == role.ErrRoleRequestNotFound {
		fmt.Println("✅ Correctly returned ErrRoleRequestNotFound")
	} else {
		log.Fatalf("❌ Expected ErrRoleRequestNotFound, got: %v", err)
	}

	// Test 10: Test GetPendingRequestByUserID with no pending requests
	fmt.Println("\n10. Testing GetPendingRequestByUserID (no pending)...")
	noPending, err := repo.GetPendingRequestByUserID(ctx, testUserID)
	if err != nil {
		log.Fatalf("❌ Failed to query: %v", err)
	}
	if noPending == nil {
		fmt.Println("✅ Correctly returned nil for user with no pending requests")
	} else {
		log.Fatalf("❌ Expected nil, got request: %s", noPending.RequestID())
	}
}
