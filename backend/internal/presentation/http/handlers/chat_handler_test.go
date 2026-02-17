package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/yuki5155/go-google-auth/internal/application/ports"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestNewChatHandler(t *testing.T) {
	h := NewChatHandler(nil, nil, nil, nil, nil, nil)
	assert.NotNil(t, h)
}

func TestWriteSSEEvent(t *testing.T) {
	tests := []struct {
		name     string
		event    string
		data     []byte
		expected string
	}{
		{
			name:     "done event",
			event:    "done",
			data:     []byte(`{"id":"123"}`),
			expected: "event: done\ndata: {\"id\":\"123\"}\n\n",
		},
		{
			name:     "error event",
			event:    "error",
			data:     []byte(`"something went wrong"`),
			expected: "event: error\ndata: \"something went wrong\"\n\n",
		},
		{
			name:     "tool_start event",
			event:    "tool_start",
			data:     []byte(`{"tool":"get_current_time"}`),
			expected: "event: tool_start\ndata: {\"tool\":\"get_current_time\"}\n\n",
		},
		{
			name:     "tool_end event",
			event:    "tool_end",
			data:     []byte(`{"tool":"get_current_time"}`),
			expected: "event: tool_end\ndata: {\"tool\":\"get_current_time\"}\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			writeSSEEvent(&buf, tt.event, tt.data)
			assert.Equal(t, tt.expected, buf.String())
		})
	}
}

func TestGetClaims_NoClaims(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	h := NewChatHandler(nil, nil, nil, nil, nil, nil)
	claims := h.getClaims(c)

	assert.Nil(t, claims)
	assert.Len(t, c.Errors, 1)
}

func TestGetClaims_InvalidType(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Set("claims", "not-a-claims-struct")

	h := NewChatHandler(nil, nil, nil, nil, nil, nil)
	claims := h.getClaims(c)

	assert.Nil(t, claims)
	assert.Len(t, c.Errors, 1)
}

func TestGetClaims_ValidClaims(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Set("claims", &ports.TokenClaims{
		UserID: "user-123",
		Email:  "test@example.com",
		Role:   "subscriber",
	})

	h := NewChatHandler(nil, nil, nil, nil, nil, nil)
	claims := h.getClaims(c)

	assert.NotNil(t, claims)
	assert.Equal(t, "user-123", claims.UserID)
}

func TestCreateThread_NoClaims(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	h := NewChatHandler(nil, nil, nil, nil, nil, nil)
	h.CreateThread(c)

	assert.Len(t, c.Errors, 1)
}

func TestListThreads_NoClaims(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	h := NewChatHandler(nil, nil, nil, nil, nil, nil)
	h.ListThreads(c)

	assert.Len(t, c.Errors, 1)
}

func TestGetThread_NoClaims(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	h := NewChatHandler(nil, nil, nil, nil, nil, nil)
	h.GetThread(c)

	assert.Len(t, c.Errors, 1)
}

func TestDeleteThread_NoClaims(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/", nil)

	h := NewChatHandler(nil, nil, nil, nil, nil, nil)
	h.DeleteThread(c)

	assert.Len(t, c.Errors, 1)
}

func TestSendMessage_NoClaims(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	h := NewChatHandler(nil, nil, nil, nil, nil, nil)
	h.SendMessage(c)

	assert.Len(t, c.Errors, 1)
}

func TestStreamMessage_NoClaims(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	h := NewChatHandler(nil, nil, nil, nil, nil, nil)
	h.StreamMessage(c)

	assert.Len(t, c.Errors, 1)
}
