package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/yuki5155/go-google-auth/internal/application/ports"
	"github.com/yuki5155/go-google-auth/internal/domain/role"
	"github.com/yuki5155/go-google-auth/internal/domain/user"
	"github.com/yuki5155/go-google-auth/internal/mocks"
	"go.uber.org/mock/gomock"
)

func TestRequireSubscriber(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*mocks.MockRoleRepository)
		claims         *ports.TokenClaims
		wantStatusCode int
	}{
		{
			name: "Root user allowed",
			claims: &ports.TokenClaims{
				UserID: "root-123",
				Email:  "root@example.com",
				Role:   "root",
			},
			setupMock: func(repo *mocks.MockRoleRepository) {
				// No DB call needed for root user
			},
			wantStatusCode: 200,
		},
		{
			name: "Admin user allowed",
			claims: &ports.TokenClaims{
				UserID: "admin-123",
				Email:  "admin@example.com",
				Role:   "admin",
			},
			setupMock: func(repo *mocks.MockRoleRepository) {
				// No DB call needed for admin user with role in token
			},
			wantStatusCode: 200,
		},
		{
			name: "Subscriber user allowed",
			claims: &ports.TokenClaims{
				UserID: "sub-123",
				Email:  "subscriber@example.com",
				Role:   "subscriber",
			},
			setupMock: func(repo *mocks.MockRoleRepository) {
				// No DB call needed for subscriber user with role in token
			},
			wantStatusCode: 200,
		},
		{
			name: "Regular user with subscriber role in DB allowed",
			claims: &ports.TokenClaims{
				UserID: "user-123",
				Email:  "user@example.com",
				Role:   "user",
			},
			setupMock: func(repo *mocks.MockRoleRepository) {
				subRole, _ := role.NewUserRole("user-123", "user@example.com", user.RoleSubscriber, "admin@example.com")
				repo.EXPECT().
					GetUserRole(gomock.Any(), "user-123").
					Return(subRole, nil)
			},
			wantStatusCode: 200,
		},
		{
			name: "Regular user with admin role in DB allowed",
			claims: &ports.TokenClaims{
				UserID: "user-123",
				Email:  "user@example.com",
				Role:   "user",
			},
			setupMock: func(repo *mocks.MockRoleRepository) {
				adminRole, _ := role.NewUserRole("user-123", "user@example.com", user.RoleAdmin, "root@example.com")
				repo.EXPECT().
					GetUserRole(gomock.Any(), "user-123").
					Return(adminRole, nil)
			},
			wantStatusCode: 200,
		},
		{
			name: "Regular user forbidden",
			claims: &ports.TokenClaims{
				UserID: "user-123",
				Email:  "user@example.com",
				Role:   "user",
			},
			setupMock: func(repo *mocks.MockRoleRepository) {
				userRole, _ := role.NewUserRole("user-123", "user@example.com", user.RoleUser, "admin@example.com")
				repo.EXPECT().
					GetUserRole(gomock.Any(), "user-123").
					Return(userRole, nil)
			},
			wantStatusCode: 403,
		},
		{
			name: "User not found in DB forbidden",
			claims: &ports.TokenClaims{
				UserID: "new-123",
				Email:  "newuser@example.com",
				Role:   "user",
			},
			setupMock: func(repo *mocks.MockRoleRepository) {
				repo.EXPECT().
					GetUserRole(gomock.Any(), "new-123").
					Return(nil, role.ErrUserRoleNotFound)
			},
			wantStatusCode: 403,
		},
		{
			name:   "No claims in context",
			claims: nil,
			setupMock: func(repo *mocks.MockRoleRepository) {
				// No DB call
			},
			wantStatusCode: 401,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := mocks.NewMockRoleRepository(ctrl)
			tt.setupMock(mockRepo)

			gin.SetMode(gin.TestMode)
			router := gin.New()

			// Add error handler middleware first to properly handle errors
			router.Use(ErrorHandler())

			router.Use(func(c *gin.Context) {
				if tt.claims != nil {
					c.Set("claims", tt.claims)
				}
				c.Next()
			})
			router.Use(RequireSubscriber(mockRepo))

			router.GET("/test", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "success"})
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test", nil)
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("Status code = %v, want %v", w.Code, tt.wantStatusCode)
				t.Logf("Response body: %s", w.Body.String())
			}
		})
	}
}

func TestRequireSubscriber_InvalidClaimsType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockRoleRepository(ctrl)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(ErrorHandler())

	// Set wrong type of claims
	router.Use(func(c *gin.Context) {
		c.Set("claims", "wrong type")
		c.Next()
	})
	router.Use(RequireSubscriber(mockRepo))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("Status code = %v, want %v", w.Code, 401)
	}
}
