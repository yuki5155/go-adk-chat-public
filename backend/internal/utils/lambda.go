package utils

import (
	"os"

	"github.com/aws/aws-lambda-go/events"
)

// GetCORSHeaders returns the common CORS headers for Lambda responses
func GetCORSHeaders() map[string]string {
	// Allowed origin should match the frontend URL
	allowedOrigin := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigin == "" {
		allowedOrigin = "*"
	}

	return map[string]string{
		"Access-Control-Allow-Origin":      allowedOrigin,
		"Access-Control-Allow-Credentials": "true",
		"Access-Control-Allow-Headers":     "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
		"Access-Control-Allow-Methods":     "GET,POST,PUT,DELETE,OPTIONS",
	}
}

// ApiResponse creates a formatted APIGatewayProxyResponse with CORS headers
func ApiResponse(statusCode int, body string) events.APIGatewayProxyResponse {
	headers := GetCORSHeaders()
	headers["Content-Type"] = "application/json"

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       body,
		Headers:    headers,
	}
}
