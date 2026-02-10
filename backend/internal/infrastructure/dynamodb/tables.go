package dynamodb

import (
	"fmt"
	"os"
)

// TableNames holds the names of all DynamoDB tables
type TableNames struct {
	UserRoles    string
	RoleRequests string
	ChatThreads  string
	ChatSessions string
	ChatEvents   string
	ChatMemories string
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

	// Chat tables
	chatThreadsTable := os.Getenv("DYNAMODB_CHAT_THREADS_TABLE")
	chatSessionsTable := os.Getenv("DYNAMODB_CHAT_SESSIONS_TABLE")
	chatEventsTable := os.Getenv("DYNAMODB_CHAT_EVENTS_TABLE")
	chatMemoriesTable := os.Getenv("DYNAMODB_CHAT_MEMORIES_TABLE")

	projectName := getEnv("PROJECT_NAME", "go-adk-chat")
	environment := getEnv("GO_ENV", "dev")

	if chatThreadsTable == "" {
		chatThreadsTable = fmt.Sprintf("%s-%s-chat-threads", projectName, environment)
	}
	if chatSessionsTable == "" {
		chatSessionsTable = fmt.Sprintf("%s-%s-chat-sessions", projectName, environment)
	}
	if chatEventsTable == "" {
		chatEventsTable = fmt.Sprintf("%s-%s-chat-events", projectName, environment)
	}
	if chatMemoriesTable == "" {
		chatMemoriesTable = fmt.Sprintf("%s-%s-chat-memories", projectName, environment)
	}

	return TableNames{
		UserRoles:    userRolesTable,
		RoleRequests: roleRequestsTable,
		ChatThreads:  chatThreadsTable,
		ChatSessions: chatSessionsTable,
		ChatEvents:   chatEventsTable,
		ChatMemories: chatMemoriesTable,
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

// GetChatThreadsTableName returns the chat threads table name
func GetChatThreadsTableName() string {
	return GetTableNames().ChatThreads
}

// GetChatSessionsTableName returns the chat sessions table name
func GetChatSessionsTableName() string {
	return GetTableNames().ChatSessions
}

// GetChatEventsTableName returns the chat events table name
func GetChatEventsTableName() string {
	return GetTableNames().ChatEvents
}

// GetChatMemoriesTableName returns the chat memories table name
func GetChatMemoriesTableName() string {
	return GetTableNames().ChatMemories
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
