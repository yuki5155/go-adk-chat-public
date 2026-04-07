package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	chatApp "github.com/yuki5155/go-google-auth/internal/application/chat"
	"github.com/yuki5155/go-google-auth/internal/application/ports"
	"github.com/yuki5155/go-google-auth/internal/domain/shared"
)

// ChatHandler handles chat-related HTTP requests
type ChatHandler struct {
	createThreadUC  *chatApp.CreateThreadUseCase
	listThreadsUC   *chatApp.ListThreadsUseCase
	getThreadUC     *chatApp.GetThreadUseCase
	sendMessageUC   *chatApp.SendMessageUseCase
	deleteThreadUC  *chatApp.DeleteThreadUseCase
	listModelsUC    *chatApp.ListModelsUseCase
}

// NewChatHandler creates a new ChatHandler
func NewChatHandler(
	createThreadUC *chatApp.CreateThreadUseCase,
	listThreadsUC *chatApp.ListThreadsUseCase,
	getThreadUC *chatApp.GetThreadUseCase,
	sendMessageUC *chatApp.SendMessageUseCase,
	deleteThreadUC *chatApp.DeleteThreadUseCase,
	listModelsUC *chatApp.ListModelsUseCase,
) *ChatHandler {
	return &ChatHandler{
		createThreadUC:  createThreadUC,
		listThreadsUC:   listThreadsUC,
		getThreadUC:     getThreadUC,
		sendMessageUC:   sendMessageUC,
		deleteThreadUC:  deleteThreadUC,
		listModelsUC:    listModelsUC,
	}
}

// ListModels returns the list of available LLM models
func (h *ChatHandler) ListModels(c *gin.Context) {
	dto := h.listModelsUC.Execute(c.Request.Context())

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    dto,
	})
}

// CreateThread creates a new chat thread
func (h *ChatHandler) CreateThread(c *gin.Context) {
	claims := h.getClaims(c)
	if claims == nil {
		return
	}

	var req struct {
		Title string `json:"title"`
		Model string `json:"model"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Title and Model are optional, so we don't error out
		req.Title = ""
		req.Model = ""
	}

	cmd := chatApp.CreateThreadCommand{
		UserID: claims.UserID,
		Title:  req.Title,
		Model:  req.Model,
	}

	dto, err := h.createThreadUC.Execute(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    dto,
	})
}

// ListThreads lists all threads for the current user
func (h *ChatHandler) ListThreads(c *gin.Context) {
	claims := h.getClaims(c)
	if claims == nil {
		return
	}

	// Parse query parameters
	limit := 20
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	lastKey := c.Query("last_key")

	query := chatApp.ListThreadsQuery{
		UserID:  claims.UserID,
		Limit:   limit,
		LastKey: lastKey,
	}

	dto, err := h.listThreadsUC.Execute(c.Request.Context(), query)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    dto,
	})
}

// GetThread gets a thread with its messages
func (h *ChatHandler) GetThread(c *gin.Context) {
	claims := h.getClaims(c)
	if claims == nil {
		return
	}

	threadID := c.Param("id")
	if threadID == "" {
		_ = c.Error(shared.NewBadRequestError("INVALID_REQUEST", "Thread ID is required", nil))
		return
	}

	limit := 100
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	query := chatApp.GetThreadQuery{
		UserID:   claims.UserID,
		ThreadID: threadID,
		Limit:    limit,
	}

	dto, err := h.getThreadUC.Execute(c.Request.Context(), query)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    dto,
	})
}

// DeleteThread deletes a thread
func (h *ChatHandler) DeleteThread(c *gin.Context) {
	claims := h.getClaims(c)
	if claims == nil {
		return
	}

	threadID := c.Param("id")
	if threadID == "" {
		_ = c.Error(shared.NewBadRequestError("INVALID_REQUEST", "Thread ID is required", nil))
		return
	}

	cmd := chatApp.DeleteThreadCommand{
		UserID:   claims.UserID,
		ThreadID: threadID,
	}

	if err := h.deleteThreadUC.Execute(c.Request.Context(), cmd); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"message": "Thread deleted successfully",
		},
	})
}

// SendMessage sends a message to a thread
func (h *ChatHandler) SendMessage(c *gin.Context) {
	claims := h.getClaims(c)
	if claims == nil {
		return
	}

	threadID := c.Param("id")
	if threadID == "" {
		_ = c.Error(shared.NewBadRequestError("INVALID_REQUEST", "Thread ID is required", nil))
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(shared.NewBadRequestError("INVALID_REQUEST", "Message content is required", err))
		return
	}

	cmd := chatApp.SendMessageCommand{
		UserID:   claims.UserID,
		ThreadID: threadID,
		Content:  req.Content,
	}

	dto, err := h.sendMessageUC.Execute(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    dto,
	})
}

// StreamMessage sends a message and streams the response using SSE
func (h *ChatHandler) StreamMessage(c *gin.Context) {
	claims := h.getClaims(c)
	if claims == nil {
		return
	}

	threadID := c.Param("id")
	if threadID == "" {
		_ = c.Error(shared.NewBadRequestError("INVALID_REQUEST", "Thread ID is required", nil))
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(shared.NewBadRequestError("INVALID_REQUEST", "Message content is required", err))
		return
	}

	cmd := chatApp.SendMessageCommand{
		UserID:   claims.UserID,
		ThreadID: threadID,
		Content:  req.Content,
	}

	// Set SSE headers
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")

	// Get flusher interface
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		_ = c.Error(shared.NewInternalError("STREAM_ERROR", "Streaming not supported", nil))
		return
	}

	// Write status and send initial ping to establish the stream
	c.Writer.WriteHeader(http.StatusOK)
	fmt.Fprintf(c.Writer, ": ping\n\n")
	flusher.Flush()

	log.Printf("[Stream] Starting stream for thread %s", threadID)

	// Stream the response
	chunkCount := 0
	dto, err := h.sendMessageUC.ExecuteStream(c.Request.Context(), cmd, func(event ports.StreamEvent) error {
		switch event.Type {
		case ports.StreamEventChunk:
			chunkCount++
			log.Printf("[Stream] Sending chunk %d: %d bytes", chunkCount, len(event.Content))
			encoded, _ := json.Marshal(event.Content)
			_, writeErr := fmt.Fprintf(c.Writer, "data: %s\n\n", encoded)
			if writeErr != nil {
				return writeErr
			}
		case ports.StreamEventToolStart:
			toolData, _ := json.Marshal(map[string]string{"tool": event.ToolCall.Name})
			writeSSEEvent(c.Writer, "tool_start", toolData)
		case ports.StreamEventToolEnd:
			toolData, _ := json.Marshal(map[string]string{"tool": event.ToolCall.Name})
			writeSSEEvent(c.Writer, "tool_end", toolData)
		}
		flusher.Flush()
		return nil
	})

	if err != nil {
		// Send error event with JSON-encoded error message to prevent XSS
		sanitizedErr, _ := json.Marshal(err.Error())
		writeSSEEvent(c.Writer, "error", sanitizedErr)
		flusher.Flush()
		return
	}

	// Send done event with JSON-marshaled response to prevent XSS
	doneData, _ := json.Marshal(map[string]string{
		"message_id":  dto.Message.MessageID,
		"response_id": dto.Response.MessageID,
	})
	writeSSEEvent(c.Writer, "done", doneData)
	flusher.Flush()
}

// TestStream is a debug endpoint to test SSE streaming without AI
func (h *ChatHandler) TestStream(c *gin.Context) {
	// Set SSE headers
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Streaming not supported"})
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
	fmt.Fprintf(c.Writer, ": ping\n\n")
	flusher.Flush()

	// Send test chunks with delays
	words := []string{"Hello", " this", " is", " a", " streaming", " test", "!"}
	for i, word := range words {
		log.Printf("[TestStream] Sending chunk %d: %s", i+1, word)
		fmt.Fprintf(c.Writer, "data: %s\n\n", word)
		flusher.Flush()
		// Delay between chunks to demonstrate streaming
		select {
		case <-c.Request.Context().Done():
			return
		case <-time.After(200 * time.Millisecond):
			// Continue to next chunk
		}
	}

	fmt.Fprintf(c.Writer, "event: done\ndata: complete\n\n")
	flusher.Flush()
}

// writeSSEEvent writes a named SSE event with pre-sanitized JSON data to the writer.
// The data parameter must already be JSON-encoded to prevent XSS.
func writeSSEEvent(w io.Writer, event string, data []byte) {
	var buf []byte
	buf = append(buf, "event: "...)
	buf = append(buf, event...)
	buf = append(buf, "\ndata: "...)
	buf = append(buf, data...)
	buf = append(buf, "\n\n"...)
	_, _ = w.Write(buf)
}

// getClaims is a helper function to extract claims from context
func (h *ChatHandler) getClaims(c *gin.Context) *ports.TokenClaims {
	claimsInterface, exists := c.Get("claims")
	if !exists {
		_ = c.Error(shared.NewUnauthorizedError("UNAUTHORIZED", "User not authenticated", nil))
		return nil
	}

	claims, ok := claimsInterface.(*ports.TokenClaims)
	if !ok {
		_ = c.Error(shared.NewUnauthorizedError("UNAUTHORIZED", "Invalid authentication", nil))
		return nil
	}

	return claims
}
