package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuki5155/go-google-auth/internal/infrastructure/config"
)

func TestHealthHandler_Handle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHealthHandler()

	router := gin.New()
	router.GET(handler.Path, handler.Handle)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/health", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "ok", response["status"])
}

func TestNewHealthHandler(t *testing.T) {
	handler := NewHealthHandler()

	assert.NotNil(t, handler)
	assert.Equal(t, "/health", handler.Path)
}

func TestHelloHandler_Handle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHelloHandler()

	router := gin.New()
	router.GET(handler.Path, handler.Handle)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/hello", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "Hello World", response["message"])
}

func TestNewHelloHandler(t *testing.T) {
	handler := NewHelloHandler()

	assert.NotNil(t, handler)
	assert.Equal(t, "/hello", handler.Path)
}

func TestCookieHandler_Handle_Development(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create development config
	cfg := &config.Config{
		Environment: "development",
	}

	handler := NewCookieHandler(cfg)

	router := gin.New()
	router.GET(handler.Path, handler.Handle)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/set-cookie", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Check JSON response
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "Cookie has been set", response["message"])
	assert.Equal(t, "test_session", response["cookieName"])
	assert.Equal(t, "session_value_12345", response["cookieValue"])
	assert.Equal(t, "Lax", response["sameSite"])
	assert.Equal(t, false, response["secure"]) // Development = not secure
	assert.NotEmpty(t, response["timestamp"])

	// Check Set-Cookie header
	cookies := w.Result().Cookies()
	assert.Len(t, cookies, 1)

	cookie := cookies[0]
	assert.Equal(t, "test_session", cookie.Name)
	assert.Equal(t, "session_value_12345", cookie.Value)
	assert.Equal(t, "/", cookie.Path)
	assert.Equal(t, 3600, cookie.MaxAge)
	assert.True(t, cookie.HttpOnly)
	assert.False(t, cookie.Secure) // Development
	assert.Equal(t, http.SameSiteLaxMode, cookie.SameSite)
}

func TestCookieHandler_Handle_Production(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create production config
	cfg := &config.Config{
		Environment: "production",
	}

	handler := NewCookieHandler(cfg)

	router := gin.New()
	router.GET(handler.Path, handler.Handle)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/set-cookie", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Check JSON response
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, true, response["secure"]) // Production = secure

	// Check Set-Cookie header
	cookies := w.Result().Cookies()
	assert.Len(t, cookies, 1)

	cookie := cookies[0]
	assert.True(t, cookie.Secure) // Production
}

func TestNewCookieHandler(t *testing.T) {
	cfg := &config.Config{}
	handler := NewCookieHandler(cfg)

	assert.NotNil(t, handler)
	assert.Equal(t, "/api/set-cookie", handler.Path)
	assert.Equal(t, cfg, handler.config)
}

func TestCheckCookieHandler_Handle_WithCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewCheckCookieHandler()

	router := gin.New()
	router.GET(handler.Path, handler.Handle)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/check-cookie", nil)

	// Add cookie to request
	req.AddCookie(&http.Cookie{
		Name:  "test_session",
		Value: "session_value_12345",
	})

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Check JSON response
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, true, response["cookieReceived"])
	assert.NotEmpty(t, response["cookies"])
	assert.Equal(t, "session_value_12345", response["testSession"])
	assert.NotEmpty(t, response["timestamp"])
}

func TestCheckCookieHandler_Handle_WithoutCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewCheckCookieHandler()

	router := gin.New()
	router.GET(handler.Path, handler.Handle)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/check-cookie", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Check JSON response
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, false, response["cookieReceived"])
	assert.Empty(t, response["cookies"])
	assert.Empty(t, response["testSession"])
	assert.NotEmpty(t, response["timestamp"])
}

func TestCheckCookieHandler_Handle_WithMultipleCookies(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewCheckCookieHandler()

	router := gin.New()
	router.GET(handler.Path, handler.Handle)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/check-cookie", nil)

	// Add multiple cookies
	req.AddCookie(&http.Cookie{
		Name:  "test_session",
		Value: "session_value_12345",
	})
	req.AddCookie(&http.Cookie{
		Name:  "other_cookie",
		Value: "other_value",
	})

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Check JSON response
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, true, response["cookieReceived"])
	assert.Equal(t, "session_value_12345", response["testSession"])
	assert.Contains(t, response["cookies"], "test_session")
	assert.Contains(t, response["cookies"], "other_cookie")
}

func TestNewCheckCookieHandler(t *testing.T) {
	handler := NewCheckCookieHandler()

	assert.NotNil(t, handler)
	assert.Equal(t, "/api/check-cookie", handler.Path)
}
