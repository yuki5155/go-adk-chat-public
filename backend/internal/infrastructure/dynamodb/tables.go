package dynamodb

import (
	"fmt"
	"os"
)

// TableNames holds the names of all DynamoDB tables
type TableNames struct {
	UserRoles    string
	RoleRequests string
}

// GetTableNames returns the table names based on environment configuration
func GetTableNames() TableNames {
	// Try to use explicit environment variables first (for Lambda)
	userRolesTable := os.Getenv("DYNAMODB_USER_ROLES_TABLE")
	roleRequestsTable := os.Getenv("DYNAMODB_ROLE_REQUESTS_TABLE")

	// Fallback to constructed names for local development
	if userRolesTable == "" {
		projectName := getEnv("PROJECT_NAME", "go-adk-chat")
		environment := getEnv("GO_ENV", "dev")
		userRolesTable = fmt.Sprintf("%s-%s-user-roles", projectName, environment)
	}

	if roleRequestsTable == "" {
		projectName := getEnv("PROJECT_NAME", "go-adk-chat")
		environment := getEnv("GO_ENV", "dev")
		roleRequestsTable = fmt.Sprintf("%s-%s-role-requests", projectName, environment)
	}

	return TableNames{
		UserRoles:    userRolesTable,
		RoleRequests: roleRequestsTable,
	}
}

// GetUserRolesTableName returns the user roles table name
func GetUserRolesTableName() string {
	return GetTableNames().UserRoles
}

// GetRoleRequestsTableName returns the role requests table name
func GetRoleRequestsTableName() string {
	return GetTableNames().RoleRequests
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
