package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/yuki5155/go-google-auth/internal/presentation/http/handlers"
	"github.com/yuki5155/go-google-auth/internal/presentation/lambda/common"
)

var ginLambda *ginadapter.GinLambda

func init() {
	// Bootstrap with shared initialization
	r, c := common.Bootstrap()

	// Create auth handler using use cases from container
	authHandler := handlers.NewAuthHandler(
		c.GoogleLoginUseCase,
		c.RefreshTokenUseCase,
		c.GetCurrentUserUseCase,
		c.LogoutUseCase,
		c.TokenGenerator,
		c.Config,
	)

	// Register this Lambda's specific endpoint
	r.POST("/auth/refresh", authHandler.RefreshToken)

	// Wrap Gin router with Lambda adapter
	ginLambda = ginadapter.New(r)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
