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
	// Bootstrap with shared initialization
	r, c := common.Bootstrap()

	// Add error handler middleware
	r.Use(middleware.ErrorHandler())

	// Create role handler using use cases from container
	roleHandler := handlers.NewRoleHandler(
		nil, // Not needed for this endpoint
		nil, // Not needed for this endpoint
		nil, // Not needed for this endpoint
		nil, // Not needed for this endpoint
		c.ListUsersByRoleUseCase,
	)

	// Register protected admin route with auth and admin middleware
	r.GET("/api/admin/users",
		middleware.Auth(c.TokenGenerator),
		middleware.RequireAdmin(c.RoleRepository),
		roleHandler.ListUsers)

	// Wrap Gin router with Lambda adapter
	ginLambda = ginadapter.New(r)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
