package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yuki5155/go-google-auth/internal/infrastructure/config"
)

// CookieHandler はCookieテスト用のハンドラー
type CookieHandler struct {
	Path   string
	config *config.Config
}

// NewCookieHandler はCookieHandlerの新しいインスタンスを返す
func NewCookieHandler(cfg *config.Config) *CookieHandler {
	return &CookieHandler{
		Path:   "/api/set-cookie",
		config: cfg,
	}
}

// Handle はSet-Cookieのテストを行う
func (h *CookieHandler) Handle(c *gin.Context) {
	// 本番環境ではSecureフラグを有効にする
	secure := h.config.IsProduction()

	// CookieをSameSite=Lax属性で設定
	cookie := &http.Cookie{
		Name:     "test_session",
		Value:    "session_value_12345",
		Path:     "/",
		MaxAge:   3600, // 1 hour
		HttpOnly: true,
		Secure:   secure, // 本番環境（HTTPS）ではtrue
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(c.Writer, cookie)

	c.JSON(http.StatusOK, gin.H{
		"message":     "Cookie has been set",
		"cookieName":  "test_session",
		"cookieValue": "session_value_12345",
		"sameSite":    "Lax",
		"secure":      secure,
		"timestamp":   time.Now().Format(time.RFC3339),
	})
}

// CheckCookieHandler はCookieが送信されているか確認するハンドラー
type CheckCookieHandler struct {
	Path string
}

// NewCheckCookieHandler はCheckCookieHandlerの新しいインスタンスを返す
func NewCheckCookieHandler() *CheckCookieHandler {
	return &CheckCookieHandler{
		Path: "/api/check-cookie",
	}
}

// Handle はCookieが正しく送信されているか確認する
func (h *CheckCookieHandler) Handle(c *gin.Context) {
	// リクエストから全てのCookieを取得
	cookies := c.Request.Header.Get("Cookie")

	// 特定のCookieを取得
	testSession, err := c.Cookie("test_session")

	cookieReceived := err == nil && testSession != ""

	c.JSON(http.StatusOK, gin.H{
		"cookieReceived": cookieReceived,
		"cookies":        cookies,
		"testSession":    testSession,
		"timestamp":      time.Now().Format(time.RFC3339),
	})
}
