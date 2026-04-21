package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/yuki5155/go-google-auth/internal/presentation/http/handlers"
	"github.com/yuki5155/go-google-auth/internal/presentation/http/middleware"
	"github.com/yuki5155/go-google-auth/internal/presentation/lambda/common"
)

var ginLambda *ginadapter.GinLambda

func init() {
	r, c := common.Bootstrap()

	chatHandler := handlers.NewChatHandler(
		c.CreateThreadUseCase,
		c.ListThreadsUseCase,
		c.GetThreadUseCase,
		c.SendMessageUseCase,
		c.DeleteThreadUseCase,
		c.ListModelsUseCase,
	)

	r.Use(middleware.Auth(c.TokenGenerator))
	r.Use(middleware.RequireSubscriber(c.CheckUserRoleUseCase))

	// Non-streaming message endpoint (fallback)
	r.POST("/api/chat/threads/:id/message", chatHandler.SendMessage)

	ginLambda = ginadapter.New(r)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
