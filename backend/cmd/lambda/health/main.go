package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/yuki5155/go-google-auth/internal/handlers"
	"github.com/yuki5155/go-google-auth/internal/presentation/lambda/common"
)

var ginLambda *ginadapter.GinLambda

func init() {
	// Bootstrap with shared initialization
	r, _ := common.Bootstrap()

	// Register health endpoints
	healthHandler := handlers.NewHealthHandler()
	r.GET("/health", healthHandler.Handle)
	r.GET("/health/ready", healthHandler.Handle)

	// Wrap Gin router with Lambda adapter
	ginLambda = ginadapter.New(r)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
