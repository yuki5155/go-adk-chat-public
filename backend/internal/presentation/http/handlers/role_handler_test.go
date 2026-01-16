package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/yuki5155/go-google-auth/internal/application/admin"
	"github.com/yuki5155/go-google-auth/internal/application/ports"
	"github.com/yuki5155/go-google-auth/internal/domain/role"
	"github.com/yuki5155/go-google-auth/internal/domain/user"
	"github.com/yuki5155/go-google-auth/internal/mocks"
	"github.com/yuki5155/go-google-auth/internal/presentation/http/middleware"
	"go.uber.org/mock/gomock"
)

func TestRoleHandler_RequestRole(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		setupAuth      func(*gin.Context)
		setupMock      func(*mocks.MockRoleRepository)
		wantStatusCode int
		wantErrorCode  string
	}{
		{
			name:        "Successful role request",
			requestBody: `{"requested_role": "subscriber"}`,
			setupAuth: func(c *gin.Context) {
				claims := &ports.TokenClaims{
					UserID: "user-123",
					Email:  "user@example.com",
				}
				c.Set("claims", claims)
			},
			setupMock: func(repo *mocks.MockRoleRepository) {
				// No pending request
				repo.EXPECT().
					GetPendingRequestByUserID(gomock.Any(), "user-123").
					Return(nil, nil)
				// Create succeeds
				repo.EXPECT().
					CreateRoleRequest(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantStatusCode: 201,
		},
		{
			name:        "Missing authentication",
			requestBody: `{"requested_role": "subscriber"}`,
			setupAuth: func(c *gin.Context) {
				// No claims set
			},
			setupMock:      func(repo *mocks.MockRoleRepository) {},
			wantStatusCode: 401,
			wantErrorCode:  "UNAUTHORIZED",
		},
		{
			name:        "Invalid JSON body",
			requestBody: `invalid json`,
			setupAuth: func(c *gin.Context) {
				claims := &ports.TokenClaims{
					UserID: "user-123",
					Email:  "user@example.com",
				}
				c.Set("claims", claims)
			},
			setupMock:      func(repo *mocks.MockRoleRepository) {},
			wantStatusCode: 400,
			wantErrorCode:  "INVALID_REQUEST",
		},
		{
			name:        "Missing requested_role field",
			requestBody: `{}`,
			setupAuth: func(c *gin.Context) {
				claims := &ports.TokenClaims{
					UserID: "user-123",
					Email:  "user@example.com",
				}
				c.Set("claims", claims)
			},
			setupMock:      func(repo *mocks.MockRoleRepository) {},
			wantStatusCode: 400,
			wantErrorCode:  "INVALID_REQUEST",
		},
		{
			name:        "Invalid role type - admin",
			requestBody: `{"requested_role": "admin"}`,
			setupAuth: func(c *gin.Context) {
				claims := &ports.TokenClaims{
					UserID: "user-123",
					Email:  "user@example.com",
				}
				c.Set("claims", claims)
			},
			setupMock:      func(repo *mocks.MockRoleRepository) {},
			wantStatusCode: 400,
			wantErrorCode:  "INVALID_REQUEST",
		},
		{
			name:        "Invalid role type - custom",
			requestBody: `{"requested_role": "superuser"}`,
			setupAuth: func(c *gin.Context) {
				claims := &ports.TokenClaims{
					UserID: "user-123",
					Email:  "user@example.com",
				}
				c.Set("claims", claims)
			},
			setupMock:      func(repo *mocks.MockRoleRepository) {},
			wantStatusCode: 400,
			wantErrorCode:  "INVALID_REQUEST",
		},
		{
			name:        "Duplicate pending request",
			requestBody: `{"requested_role": "subscriber"}`,
			setupAuth: func(c *gin.Context) {
				claims := &ports.TokenClaims{
					UserID: "user-123",
					Email:  "user@example.com",
				}
				c.Set("claims", claims)
			},
			setupMock: func(repo *mocks.MockRoleRepository) {
				// Pending request exists
				existingRequest, _ := role.NewRoleRequest("existing-req", "user-123", "user@example.com", user.RoleSubscriber)
				repo.EXPECT().
					GetPendingRequestByUserID(gomock.Any(), "user-123").
					Return(existingRequest, nil)
			},
			wantStatusCode: 409,
			wantErrorCode:  "DUPLICATE_REQUEST",
		},
		{
			name:        "Repository error",
			requestBody: `{"requested_role": "subscriber"}`,
			setupAuth: func(c *gin.Context) {
				claims := &ports.TokenClaims{
					UserID: "user-123",
					Email:  "user@example.com",
				}
				c.Set("claims", claims)
			},
			setupMock: func(repo *mocks.MockRoleRepository) {
				repo.EXPECT().
					GetPendingRequestByUserID(gomock.Any(), "user-123").
					Return(nil, errors.New("database error"))
			},
			wantStatusCode: 500,
			wantErrorCode:  "INTERNAL_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := mocks.NewMockRoleRepository(ctrl)
			tt.setupMock(mockRepo)

			// Create use case
			requestRoleUC := admin.NewRequestRoleUseCase(mockRepo)

			// Create handler
			handler := NewRoleHandler(
				requestRoleUC,
				nil, nil, nil, nil, // Other use cases not needed for this test
			)

			// Setup router
			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.Use(middleware.ErrorHandler())
			router.Use(func(c *gin.Context) {
				tt.setupAuth(c)
				c.Next()
			})
			router.POST("/api/role/request", handler.RequestRole)

			// Make request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/role/request", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			// Assertions
			if w.Code != tt.wantStatusCode {
				t.Errorf("Status code = %v, want %v", w.Code, tt.wantStatusCode)
				t.Logf("Response body: %s", w.Body.String())
			}

			// Verify response structure
			var response map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to parse response: %v", err)
			}

			// Check success field
			if tt.wantStatusCode < 300 {
				if success, ok := response["success"].(bool); !ok || !success {
					t.Error("Expected success=true for successful response")
				}
				if response["data"] == nil {
					t.Error("Expected data field in successful response")
				}
			} else {
				if success, ok := response["success"].(bool); !ok || success {
					t.Error("Expected success=false for error response")
				}
				if tt.wantErrorCode != "" {
					errorMap, ok := response["error"].(map[string]interface{})
					if !ok {
						t.Errorf("Expected error object in response")
					} else if code, ok := errorMap["code"].(string); !ok || code != tt.wantErrorCode {
						t.Errorf("Error code = %v, want %v", code, tt.wantErrorCode)
					}
				}
			}
		})
	}
}

func TestRoleHandler_RequestRole_ResponseFormat(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockRoleRepository(ctrl)

	mockRepo.EXPECT().GetPendingRequestByUserID(gomock.Any(), "user-123").Return(nil, nil)
	mockRepo.EXPECT().CreateRoleRequest(gomock.Any(), gomock.Any()).Return(nil)

	requestRoleUC := admin.NewRequestRoleUseCase(mockRepo)
	handler := NewRoleHandler(requestRoleUC, nil, nil, nil, nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.ErrorHandler())
	router.Use(func(c *gin.Context) {
		claims := &ports.TokenClaims{
			UserID: "user-123",
			Email:  "user@example.com",
		}
		c.Set("claims", claims)
		c.Next()
	})
	router.POST("/api/role/request", handler.RequestRole)

	// Make request
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/role/request", strings.NewReader(`{"requested_role": "subscriber"}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Parse response
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Verify response structure
	if success, ok := response["success"].(bool); !ok || !success {
		t.Error("Expected success=true")
	}

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected data object in response")
	}

	requiredFields := []string{"request_id", "user_id", "user_email", "requested_role", "status", "requested_at"}
	for _, field := range requiredFields {
		if data[field] == nil {
			t.Errorf("Missing required field in response: %s", field)
		}
	}

	// Verify specific values
	if userID, ok := data["user_id"].(string); !ok || userID != "user-123" {
		t.Errorf("user_id = %v, want user-123", data["user_id"])
	}
	if userEmail, ok := data["user_email"].(string); !ok || userEmail != "user@example.com" {
		t.Errorf("user_email = %v, want user@example.com", data["user_email"])
	}
	if requestedRole, ok := data["requested_role"].(string); !ok || requestedRole != "subscriber" {
		t.Errorf("requested_role = %v, want subscriber", data["requested_role"])
	}
	if status, ok := data["status"].(string); !ok || status != "pending" {
		t.Errorf("status = %v, want pending", data["status"])
	}
}

func TestRoleHandler_ListPendingRequests(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRoleRepository(ctrl)
	
	// Setup mock to return pending requests
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
		ListRoleRequestsByStatus(gomock.Any(), role.RequestStatusPending).
		Return(mockRequests, nil)

	listPendingUC := admin.NewListPendingRequestsUseCase(mockRepo)
	handler := NewRoleHandler(nil, listPendingUC, nil, nil, nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.ErrorHandler())
	router.GET("/api/admin/role-requests/pending", handler.ListPendingRequests)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/admin/role-requests/pending", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoleHandler_ApproveRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRoleRepository(ctrl)
	
	// Mock get request
	mockRequest, _ := role.NewRoleRequest("req-123", "user-123", "test@example.com", user.RoleSubscriber)
	mockRepo.EXPECT().
		GetRoleRequest(gomock.Any(), "req-123").
		Return(mockRequest, nil)
	
	// Mock update request
	mockRepo.EXPECT().
		UpdateRoleRequest(gomock.Any(), gomock.Any()).
		Return(nil)
	
	// Mock upsert user role
	mockRepo.EXPECT().
		UpsertUserRole(gomock.Any(), gomock.Any()).
		Return(nil)

	approveUC := admin.NewApproveRequestUseCase(mockRepo)
	handler := NewRoleHandler(nil, nil, approveUC, nil, nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.ErrorHandler())
	router.Use(func(c *gin.Context) {
		c.Set("claims", &ports.TokenClaims{
			UserID: "admin-1",
			Email:  "admin@example.com",
			Role:   "admin",
		})
		c.Next()
	})
	router.POST("/api/admin/role-requests/approve", handler.ApproveRequest)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/admin/role-requests/approve", 
		strings.NewReader(`{"request_id": "req-123", "notes": "Approved"}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoleHandler_RejectRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRoleRepository(ctrl)
	
	// Mock get request
	mockRequest, _ := role.NewRoleRequest("req-123", "user-123", "test@example.com", user.RoleSubscriber)
	mockRepo.EXPECT().
		GetRoleRequest(gomock.Any(), "req-123").
		Return(mockRequest, nil)
	
	// Mock update request
	mockRepo.EXPECT().
		UpdateRoleRequest(gomock.Any(), gomock.Any()).
		Return(nil)

	rejectUC := admin.NewRejectRequestUseCase(mockRepo)
	handler := NewRoleHandler(nil, nil, nil, rejectUC, nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.ErrorHandler())
	router.Use(func(c *gin.Context) {
		c.Set("claims", &ports.TokenClaims{
			Email: "admin@example.com",
		})
		c.Next()
	})
	router.POST("/api/admin/role-requests/reject", handler.RejectRequest)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/admin/role-requests/reject", 
		strings.NewReader(`{"request_id": "req-123", "notes": "Insufficient justification"}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoleHandler_ListUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRoleRepository(ctrl)
	
	// Mock list users by role
	mockUserRoles := []*role.UserRole{
		func() *role.UserRole {
			ur, _ := role.NewUserRole("user-1", "user1@example.com", user.RoleSubscriber, "admin@example.com")
			return ur
		}(),
	}
	
	mockRepo.EXPECT().
		ListUsersByRole(gomock.Any(), user.RoleSubscriber).
		Return(mockUserRoles, nil)

	listUsersUC := admin.NewListUsersByRoleUseCase(mockRepo)
	handler := NewRoleHandler(nil, nil, nil, nil, listUsersUC)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.ErrorHandler())
	router.GET("/api/admin/users", handler.ListUsers)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/admin/users?role=subscriber", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoleHandler_ListPendingRequests_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRoleRepository(ctrl)

	mockRepo.EXPECT().
		ListRoleRequestsByStatus(gomock.Any(), role.RequestStatusPending).
		Return(nil, errors.New("database error"))

	listPendingUC := admin.NewListPendingRequestsUseCase(mockRepo)
	handler := NewRoleHandler(nil, listPendingUC, nil, nil, nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.ErrorHandler())
	router.GET("/api/admin/role-requests/pending", handler.ListPendingRequests)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/admin/role-requests/pending", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestRoleHandler_ApproveRequest_MissingRequestID(t *testing.T) {
	handler := NewRoleHandler(nil, nil, nil, nil, nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.ErrorHandler())
	router.Use(func(c *gin.Context) {
		c.Set("claims", &ports.TokenClaims{
			UserID: "admin-1",
			Email:  "admin@example.com",
			Role:   "admin",
		})
		c.Next()
	})
	router.POST("/api/admin/role-requests/approve", handler.ApproveRequest)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/admin/role-requests/approve",
		strings.NewReader(`{"notes": "Approved"}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRoleHandler_ApproveRequest_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRoleRepository(ctrl)

	mockRepo.EXPECT().
		GetRoleRequest(gomock.Any(), "req-123").
		Return(nil, errors.New("database error"))

	approveUC := admin.NewApproveRequestUseCase(mockRepo)
	handler := NewRoleHandler(nil, nil, approveUC, nil, nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.ErrorHandler())
	router.Use(func(c *gin.Context) {
		c.Set("claims", &ports.TokenClaims{
			UserID: "admin-1",
			Email:  "admin@example.com",
			Role:   "admin",
		})
		c.Next()
	})
	router.POST("/api/admin/role-requests/approve", handler.ApproveRequest)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/admin/role-requests/approve",
		strings.NewReader(`{"request_id": "req-123", "notes": "Approved"}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestRoleHandler_RejectRequest_MissingRequestID(t *testing.T) {
	handler := NewRoleHandler(nil, nil, nil, nil, nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.ErrorHandler())
	router.Use(func(c *gin.Context) {
		c.Set("claims", &ports.TokenClaims{
			Email: "admin@example.com",
		})
		c.Next()
	})
	router.POST("/api/admin/role-requests/reject", handler.RejectRequest)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/admin/role-requests/reject",
		strings.NewReader(`{"notes": "Rejected"}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRoleHandler_RejectRequest_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRoleRepository(ctrl)

	mockRepo.EXPECT().
		GetRoleRequest(gomock.Any(), "req-123").
		Return(nil, errors.New("database error"))

	rejectUC := admin.NewRejectRequestUseCase(mockRepo)
	handler := NewRoleHandler(nil, nil, nil, rejectUC, nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.ErrorHandler())
	router.Use(func(c *gin.Context) {
		c.Set("claims", &ports.TokenClaims{
			Email: "admin@example.com",
		})
		c.Next()
	})
	router.POST("/api/admin/role-requests/reject", handler.RejectRequest)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/admin/role-requests/reject",
		strings.NewReader(`{"request_id": "req-123", "notes": "Rejected"}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestRoleHandler_ListUsers_InvalidRole(t *testing.T) {
	handler := NewRoleHandler(nil, nil, nil, nil, nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.ErrorHandler())
	router.GET("/api/admin/users", handler.ListUsers)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/admin/users?role=invalid_role", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRoleHandler_ListUsers_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRoleRepository(ctrl)

	mockRepo.EXPECT().
		ListUsersByRole(gomock.Any(), user.RoleSubscriber).
		Return(nil, errors.New("database error"))

	listUsersUC := admin.NewListUsersByRoleUseCase(mockRepo)
	handler := NewRoleHandler(nil, nil, nil, nil, listUsersUC)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.ErrorHandler())
	router.GET("/api/admin/users", handler.ListUsers)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/admin/users?role=subscriber", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
