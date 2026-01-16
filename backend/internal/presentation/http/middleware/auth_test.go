package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/yuki5155/go-google-auth/internal/application/ports"
	"github.com/yuki5155/go-google-auth/internal/mocks"
)

func TestAuth(t *testing.T) {
	tests := []struct {
		name           string
		setupCookie    func(*http.Request)
		setupMock      func(*mocks.MockTokenGenerator)
		wantStatus     int
		wantAborted    bool
		wantClaimsSet  bool
	}{
		{
			name: "Valid token",
			setupCookie: func(req *http.Request) {
				req.AddCookie(&http.Cookie{
					Name:  "access_token",
					Value: "valid-token",
				})
			},
			setupMock: func(m *mocks.MockTokenGenerator) {
				m.EXPECT().
					ValidateAccessToken("valid-token").
					Return(&ports.TokenClaims{
						UserID:  "user-123",
						Email:   "test@example.com",
						Name:    "Test User",
						Picture: "https://example.com/photo.jpg",
					}, nil)
			},
			wantStatus:    http.StatusOK,
			wantAborted:   false,
			wantClaimsSet: true,
		},
		{
			name: "Missing token",
			setupCookie: func(req *http.Request) {
				// No cookie added
			},
			setupMock: func(m *mocks.MockTokenGenerator) {
				// No expectations - should not be called
			},
			wantStatus:    http.StatusUnauthorized,
			wantAborted:   true,
			wantClaimsSet: false,
		},
		{
			name: "Expired token",
			setupCookie: func(req *http.Request) {
				req.AddCookie(&http.Cookie{
					Name:  "access_token",
					Value: "expired-token",
				})
			},
			setupMock: func(m *mocks.MockTokenGenerator) {
				m.EXPECT().
					ValidateAccessToken("expired-token").
					Return(nil, ports.ErrExpiredToken)
			},
			wantStatus:    http.StatusUnauthorized,
			wantAborted:   true,
			wantClaimsSet: false,
		},
		{
			name: "Invalid token",
			setupCookie: func(req *http.Request) {
				req.AddCookie(&http.Cookie{
					Name:  "access_token",
					Value: "invalid-token",
				})
			},
			setupMock: func(m *mocks.MockTokenGenerator) {
				m.EXPECT().
					ValidateAccessToken("invalid-token").
					Return(nil, errors.New("invalid token"))
			},
			wantStatus:    http.StatusUnauthorized,
			wantAborted:   true,
			wantClaimsSet: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockTokenGen := mocks.NewMockTokenGenerator(ctrl)
			tt.setupMock(mockTokenGen)

			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.Use(Auth(mockTokenGen))
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test", nil)
			tt.setupCookie(req)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantAborted {
				// Check that the handler was not executed (response doesn't contain "success")
				assert.NotContains(t, w.Body.String(), "success")
			}
		})
	}
}

func TestOptionalAuth(t *testing.T) {
	tests := []struct {
		name           string
		setupCookie    func(*http.Request)
		setupMock      func(*mocks.MockTokenGenerator)
		wantClaimsSet  bool
		wantAuthenticated bool
	}{
		{
			name: "Valid token",
			setupCookie: func(req *http.Request) {
				req.AddCookie(&http.Cookie{
					Name:  "access_token",
					Value: "valid-token",
				})
			},
			setupMock: func(m *mocks.MockTokenGenerator) {
				m.EXPECT().
					ValidateAccessToken("valid-token").
					Return(&ports.TokenClaims{
						UserID:  "user-123",
						Email:   "test@example.com",
						Name:    "Test User",
						Picture: "https://example.com/photo.jpg",
					}, nil)
			},
			wantClaimsSet:     true,
			wantAuthenticated: true,
		},
		{
			name: "Missing token - continues without auth",
			setupCookie: func(req *http.Request) {
				// No cookie added
			},
			setupMock: func(m *mocks.MockTokenGenerator) {
				// No expectations - should not be called
			},
			wantClaimsSet:     false,
			wantAuthenticated: false,
		},
		{
			name: "Invalid token - continues without auth",
			setupCookie: func(req *http.Request) {
				req.AddCookie(&http.Cookie{
					Name:  "access_token",
					Value: "invalid-token",
				})
			},
			setupMock: func(m *mocks.MockTokenGenerator) {
				m.EXPECT().
					ValidateAccessToken("invalid-token").
					Return(nil, errors.New("invalid token"))
			},
			wantClaimsSet:     false,
			wantAuthenticated: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockTokenGen := mocks.NewMockTokenGenerator(ctrl)
			tt.setupMock(mockTokenGen)

			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.Use(OptionalAuth(mockTokenGen))
			router.GET("/test", func(c *gin.Context) {
				claims, claimsExist := c.Get("claims")
				authenticated, authExists := c.Get("authenticated")

				c.JSON(http.StatusOK, gin.H{
					"claims_set":    claimsExist,
					"authenticated": authExists && authenticated.(bool),
					"claims":        claims,
				})
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test", nil)
			tt.setupCookie(req)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			if tt.wantClaimsSet {
				assert.Contains(t, w.Body.String(), `"claims_set":true`)
			} else {
				assert.Contains(t, w.Body.String(), `"claims_set":false`)
			}

			if tt.wantAuthenticated {
				assert.Contains(t, w.Body.String(), `"authenticated":true`)
			} else {
				assert.Contains(t, w.Body.String(), `"authenticated":false`)
			}
		})
	}
}
