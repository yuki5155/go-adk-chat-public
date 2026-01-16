package middleware

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/yuki5155/go-google-auth/internal/application/ports"
	"github.com/yuki5155/go-google-auth/internal/domain/role"
	"github.com/yuki5155/go-google-auth/internal/domain/user"
	"github.com/yuki5155/go-google-auth/internal/mocks"
	"go.uber.org/mock/gomock"
)

func TestRequireAdmin(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*mocks.MockRoleRepository)
		email          string
		userID         string
		rootEmail      string
		wantStatusCode int
		wantErrorCode  string
	}{
		{
			name:      "Root user allowed",
			email:     "root@example.com",
			userID:    "root-123",
			rootEmail: "root@example.com",
			setupMock: func(repo *mocks.MockRoleRepository) {
				// No DB call needed for root user
			},
			wantStatusCode: 200, // Passes to next handler
		},
		{
			name:      "Admin user allowed",
			email:     "admin@example.com",
			userID:    "admin-123",
			rootEmail: "root@example.com",
			setupMock: func(repo *mocks.MockRoleRepository) {
				adminRole, _ := role.NewUserRole("admin-123", "admin@example.com", user.RoleAdmin, "root@example.com")
				repo.EXPECT().
					GetUserRoleByEmail(gomock.Any(), "admin@example.com").
					Return(adminRole, nil)
			},
			wantStatusCode: 200,
		},
		{
			name:      "Regular user forbidden",
			email:     "user@example.com",
			userID:    "user-123",
			rootEmail: "root@example.com",
			setupMock: func(repo *mocks.MockRoleRepository) {
				userRole, _ := role.NewUserRole("user-123", "user@example.com", user.RoleUser, "admin@example.com")
				repo.EXPECT().
					GetUserRoleByEmail(gomock.Any(), "user@example.com").
					Return(userRole, nil)
			},
			wantStatusCode: 403,
			wantErrorCode:  "forbidden",
		},
		{
			name:      "Subscriber user forbidden",
			email:     "subscriber@example.com",
			userID:    "sub-123",
			rootEmail: "root@example.com",
			setupMock: func(repo *mocks.MockRoleRepository) {
				subRole, _ := role.NewUserRole("sub-123", "subscriber@example.com", user.RoleSubscriber, "admin@example.com")
				repo.EXPECT().
					GetUserRoleByEmail(gomock.Any(), "subscriber@example.com").
					Return(subRole, nil)
			},
			wantStatusCode: 403,
			wantErrorCode:  "forbidden",
		},
		{
			name:      "User not found in DB forbidden",
			email:     "newuser@example.com",
			userID:    "new-123",
			rootEmail: "root@example.com",
			setupMock: func(repo *mocks.MockRoleRepository) {
				repo.EXPECT().
					GetUserRoleByEmail(gomock.Any(), "newuser@example.com").
					Return(nil, role.ErrUserRoleNotFound)
			},
			wantStatusCode: 403,
			wantErrorCode:  "forbidden",
		},
		{
			name:      "No claims in context",
			email:     "",
			userID:    "",
			rootEmail: "root@example.com",
			setupMock: func(repo *mocks.MockRoleRepository) {
				// No DB call
			},
			wantStatusCode: 401,
			wantErrorCode:  "unauthorized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := mocks.NewMockRoleRepository(ctrl)
			tt.setupMock(mockRepo)

			// Set environment variable
			if err := os.Setenv("ROOT_USER_EMAIL", tt.rootEmail); err != nil {
				t.Fatalf("Failed to set ROOT_USER_EMAIL: %v", err)
			}
			defer func() {
				if err := os.Unsetenv("ROOT_USER_EMAIL"); err != nil {
					t.Logf("Failed to unset ROOT_USER_EMAIL: %v", err)
				}
			}()

			// Create test router
			gin.SetMode(gin.TestMode)
			router := gin.New()

			// Add middleware to set claims (simulating Auth middleware)
			router.Use(func(c *gin.Context) {
				if tt.email != "" {
					claims := &ports.TokenClaims{
						UserID: tt.userID,
						Email:  tt.email,
					}
					c.Set("claims", claims)
				}
				c.Next()
			})
			router.Use(RequireAdmin(mockRepo))

			// Add test endpoint
			router.GET("/test", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "success"})
			})

			// Make request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test", nil)
			router.ServeHTTP(w, req)

			// Assertions
			if w.Code != tt.wantStatusCode {
				t.Errorf("Status code = %v, want %v", w.Code, tt.wantStatusCode)
				t.Logf("Response body: %s", w.Body.String())
			}

			if tt.wantErrorCode != "" {
				var response map[string]interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Fatalf("Failed to parse response: %v", err)
				}
				if errorCode, ok := response["error"].(string); !ok || errorCode != tt.wantErrorCode {
					t.Errorf("Error code = %v, want %v", errorCode, tt.wantErrorCode)
				}
			}
		})
	}
}
