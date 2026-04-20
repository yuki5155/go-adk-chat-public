package common

import (
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/yuki5155/go-google-auth/internal/application/ports"
	"github.com/yuki5155/go-google-auth/internal/infrastructure/container"
)

// NewSSEHeaders returns standard SSE response headers with CORS origin resolved from the request.
func NewSSEHeaders(req events.APIGatewayProxyRequest) map[string]string {
	origin := req.Headers["origin"]
	if origin == "" {
		origin = req.Headers["Origin"]
	}
	if origin == "" {
		origin = "*"
	}
	return map[string]string{
		"Content-Type":                     "text/event-stream",
		"Cache-Control":                    "no-cache",
		"Connection":                       "keep-alive",
		"Access-Control-Allow-Origin":      origin,
		"Access-Control-Allow-Credentials": "true",
	}
}

// ValidateSubscriberAuth validates the access token and checks for subscriber-level access.
// Returns claims on success, or an error if auth fails or the role is insufficient.
func ValidateSubscriberAuth(req events.APIGatewayProxyRequest, c *container.Container) (*ports.TokenClaims, error) {
	claims, err := ValidateAuth(req, c)
	if err != nil {
		return nil, err
	}
	if !ports.HasSubscriberAccess(claims.Role) {
		return nil, fmt.Errorf("subscriber access required")
	}
	return claims, nil
}

// CreateSSEResponse creates a streaming Lambda response containing a single SSE event.
func CreateSSEResponse(statusCode int, headers map[string]string, event, data string) *events.APIGatewayProxyStreamingResponse {
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

// WriteSSE writes a named SSE event to w.
func WriteSSE(w io.Writer, event, data string) {
	if event != "" && event != "message" {
		fmt.Fprintf(w, "event: %s\n", event)
	}
	fmt.Fprintf(w, "data: %s\n\n", data)
}
