package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	chatApp "github.com/yuki5155/go-google-auth/internal/application/chat"
	"github.com/yuki5155/go-google-auth/internal/application/dto"
	"github.com/yuki5155/go-google-auth/internal/application/ports"
	"github.com/yuki5155/go-google-auth/internal/infrastructure/config"
	"github.com/yuki5155/go-google-auth/internal/infrastructure/container"
	lambdacommon "github.com/yuki5155/go-google-auth/internal/presentation/lambda/common"
)

var c *container.Container

func init() {
	cfg := config.Load()
	c = container.NewContainer(cfg)
}

// handler processes the streaming chat request and returns a streaming response
func handler(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyStreamingResponse, error) {
	log.Printf("[SSE] Received request: %s %s", req.HTTPMethod, req.Path)

	headers := lambdacommon.NewSSEHeaders(req)

	threadID := req.PathParameters["id"]
	if threadID == "" {
		log.Printf("[SSE] ERROR: No thread ID")
		return lambdacommon.CreateSSEResponse(200, headers, "error", "Thread ID is required"), nil
	}
	log.Printf("[SSE] Thread ID: %s", threadID)

	claims, err := lambdacommon.ValidateSubscriberAuth(req, c)
	if err != nil {
		log.Printf("[SSE] ERROR: Auth failed: %v", err)
		return lambdacommon.CreateSSEResponse(200, headers, "error", err.Error()), nil
	}
	log.Printf("[SSE] Auth OK, user: %s, role: %s", claims.UserID, claims.Role)

	var body dto.SendMessageRequest
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		log.Printf("[SSE] ERROR: Invalid body: %v", err)
		return lambdacommon.CreateSSEResponse(200, headers, "error", "Invalid request body"), nil
	}
	if body.Content == "" {
		log.Printf("[SSE] ERROR: Empty content")
		return lambdacommon.CreateSSEResponse(200, headers, "error", "Message content is required"), nil
	}
	log.Printf("[SSE] Message: %s", body.Content)

	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()

		fmt.Fprintf(pw, ": ping\n\n")
		log.Printf("[SSE] Sent ping")

		cmd := chatApp.SendMessageCommand{
			UserID:   claims.UserID,
			ThreadID: threadID,
			Content:  body.Content,
		}

		log.Printf("[SSE] Starting stream execution")
		result, err := c.SendMessageUseCase.ExecuteStream(ctx, cmd, func(event ports.StreamEvent) error {
			switch event.Type {
			case ports.StreamEventChunk:
				log.Printf("[SSE] Writing chunk: %d bytes", len(event.Content))
				encoded, _ := json.Marshal(event.Content)
				_, writeErr := fmt.Fprintf(pw, "data: %s\n\n", encoded)
				return writeErr
			case ports.StreamEventToolStart:
				log.Printf("[SSE] Tool start: %s", event.ToolCall.Name)
				lambdacommon.WriteSSE(pw, "tool_start", fmt.Sprintf(`{"tool":"%s"}`, event.ToolCall.Name))
			case ports.StreamEventToolEnd:
				log.Printf("[SSE] Tool end: %s", event.ToolCall.Name)
				lambdacommon.WriteSSE(pw, "tool_end", fmt.Sprintf(`{"tool":"%s"}`, event.ToolCall.Name))
			}
			return nil
		})

		if err != nil {
			log.Printf("[SSE] ERROR: Stream failed: %v", err)
			lambdacommon.WriteSSE(pw, "error", err.Error())
			return
		}

		log.Printf("[SSE] Stream completed")
		doneData := fmt.Sprintf(`{"message_id":"%s","response_id":"%s"}`, result.Message.MessageID, result.Response.MessageID)
		lambdacommon.WriteSSE(pw, "done", doneData)
	}()

	return &events.APIGatewayProxyStreamingResponse{
		StatusCode: 200,
		Headers:    headers,
		Body:       pr,
	}, nil
}

func main() {
	log.Println("[SSE] Starting Lambda with API Gateway response streaming")
	lambda.Start(handler)
}
