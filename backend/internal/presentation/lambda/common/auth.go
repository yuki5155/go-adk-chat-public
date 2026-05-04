package common

import (
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/yuki5155/go-google-auth/internal/application/ports"
	"github.com/yuki5155/go-google-auth/internal/infrastructure/container"
)

// ValidateAuth extracts and validates the access token from the Lambda Function URL request cookies.
func ValidateAuth(req events.LambdaFunctionURLRequest, c *container.Container) (*ports.TokenClaims, error) {
	cookieHeader := req.Headers["cookie"]
	if cookieHeader == "" {
		cookieHeader = req.Headers["Cookie"]
	}
	if cookieHeader == "" {
		return nil, fmt.Errorf("no auth cookie found")
	}

	var accessToken string
	for _, cookie := range strings.Split(cookieHeader, ";") {
		cookie = strings.TrimSpace(cookie)
		if strings.HasPrefix(cookie, "access_token=") {
			accessToken = strings.TrimPrefix(cookie, "access_token=")
			break
		}
	}

	if accessToken == "" {
		return nil, fmt.Errorf("access token not found")
	}

	claims, err := c.TokenGenerator.ValidateAccessToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return claims, nil
}
