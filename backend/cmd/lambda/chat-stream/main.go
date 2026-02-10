package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	chatApp "github.com/yuki5155/go-google-auth/internal/application/chat"
	"github.com/yuki5155/go-google-auth/internal/infrastructure/config"
	"github.com/yuki5155/go-google-auth/internal/infrastructure/container"
)

var c *container.Container

func init() {
	cfg := config.Load()
	c = container.NewContainer(cfg)
}

// handler processes the streaming chat request and returns a streaming response
func handler(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyStreamingResponse, error) {
	log.Printf("[SSE] Received request: %s %s", req.HTTPMethod, req.Path)

	// Set CORS origin
	origin := req.Headers["origin"]
	if origin == "" {
		origin = req.Headers["Origin"]
	}
	if origin == "" {
		origin = "*"
	}

	// Create response headers
	headers := map[string]string{
		"Content-Type":                     "text/event-stream",
		"Cache-Control":                    "no-cache",
		"Connection":                       "keep-alive",
		"Access-Control-Allow-Origin":      origin,
		"Access-Control-Allow-Credentials": "true",
	}

	// Extract thread ID from path: /api/chat/threads/{id}/stream
	threadID := req.PathParameters["id"]
	if threadID == "" {
		log.Printf("[SSE] ERROR: No thread ID")
		return createSSEResponse(200, headers, "error", "Thread ID is required"), nil
	}
	log.Printf("[SSE] Thread ID: %s", threadID)

	// Validate auth token from cookie
	claims, err := validateAuth(req)
	if err != nil {
		log.Printf("[SSE] ERROR: Auth failed: %v", err)
		return createSSEResponse(200, headers, "error", err.Error()), nil
	}
	log.Printf("[SSE] Auth OK, user: %s, role: %s", claims.UserID, claims.Role)

	// Check subscriber role
	if !isAllowedRole(claims.Role) {
		log.Printf("[SSE] ERROR: Role not allowed: %s", claims.Role)
		return createSSEResponse(200, headers, "error", "subscriber access required"), nil
	}

	// Parse request body
	var body struct {
		Content string `json:"content"`
	}
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		log.Printf("[SSE] ERROR: Invalid body: %v", err)
		return createSSEResponse(200, headers, "error", "Invalid request body"), nil
	}
	if body.Content == "" {
		log.Printf("[SSE] ERROR: Empty content")
		return createSSEResponse(200, headers, "error", "Message content is required"), nil
	}
	log.Printf("[SSE] Message: %s", body.Content)

	// Create pipe for streaming
	pr, pw := io.Pipe()

	// Execute streaming in goroutine
	go func() {
		defer pw.Close()

		// Write SSE ping
		fmt.Fprintf(pw, ": ping\n\n")
		log.Printf("[SSE] Sent ping")

		// Execute streaming
		cmd := chatApp.SendMessageCommand{
			UserID:   claims.UserID,
			ThreadID: threadID,
			Content:  body.Content,
		}

		log.Printf("[SSE] Starting stream execution")
		dto, err := c.SendMessageUseCase.ExecuteStream(ctx, cmd, func(chunk string) error {
			log.Printf("[SSE] Writing chunk: %d bytes", len(chunk))
			_, writeErr := fmt.Fprintf(pw, "data: %s\n\n", chunk)
			return writeErr
		})

		if err != nil {
			log.Printf("[SSE] ERROR: Stream failed: %v", err)
			writeSSE(pw, "error", err.Error())
			return
		}

		// Send done event
		log.Printf("[SSE] Stream completed")
		doneData := fmt.Sprintf(`{"message_id":"%s","response_id":"%s"}`, dto.Message.MessageID, dto.Response.MessageID)
		writeSSE(pw, "done", doneData)
	}()

	return &events.APIGatewayProxyStreamingResponse{
		StatusCode: 200,
		Headers:    headers,
		Body:       pr,
	}, nil
}

// createSSEResponse creates a simple SSE response with a single event
func createSSEResponse(statusCode int, headers map[string]string, event, data string) *events.APIGatewayProxyStreamingResponse {
	var content string
	if event != "" && event != "message" {
		content = fmt.Sprintf("event: %s\ndata: %s\n\n", event, data)
	} else {
		content = fmt.Sprintf("data: %s\n\n", data)
	}
	return &events.APIGatewayProxyStreamingResponse{
		StatusCode: statusCode,
		Headers:    headers,
		Body:       strings.NewReader(content),
	}
}

func writeSSE(w io.Writer, event, data string) {
	if event != "" && event != "message" {
		fmt.Fprintf(w, "event: %s\n", event)
	}
	fmt.Fprintf(w, "data: %s\n\n", data)
}

func validateAuth(req events.APIGatewayProxyRequest) (*tokenClaims, error) {
	// Get token from cookie
	cookieHeader := req.Headers["cookie"]
	if cookieHeader == "" {
		cookieHeader = req.Headers["Cookie"]
	}
	if cookieHeader == "" {
		return nil, fmt.Errorf("no auth cookie found")
	}

	// Parse cookies
	var accessToken string
	cookies := strings.Split(cookieHeader, ";")
	for _, cookie := range cookies {
		cookie = strings.TrimSpace(cookie)
		if strings.HasPrefix(cookie, "access_token=") {
			accessToken = strings.TrimPrefix(cookie, "access_token=")
			break
		}
	}

	if accessToken == "" {
		return nil, fmt.Errorf("access token not found")
	}

	// Validate token
	claims, err := c.TokenGenerator.ValidateAccessToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return &tokenClaims{
		UserID: claims.UserID,
		Email:  claims.Email,
		Role:   claims.Role,
	}, nil
}

type tokenClaims struct {
	UserID string
	Email  string
	Role   string
}

func isAllowedRole(role string) bool {
	return role == "subscriber" || role == "admin" || role == "root"
}

func main() {
	log.Println("[SSE] Starting Lambda with API Gateway response streaming")
	lambda.Start(handler)
}
