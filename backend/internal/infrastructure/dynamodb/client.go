package dynamodb

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// NewClient creates a new DynamoDB client configured for the environment
func NewClient(ctx context.Context) (*dynamodb.Client, error) {
	// Check if running in local development mode
	goEnv := os.Getenv("GO_ENV")
	if (goEnv == "development" || goEnv == "dev") && os.Getenv("DYNAMODB_ENDPOINT") != "" {
		return newLocalClient(ctx)
	}

	// Production configuration with AWS
	return newAWSClient(ctx)
}

// newLocalClient creates a DynamoDB client for local development
func newLocalClient(ctx context.Context) (*dynamodb.Client, error) {
	endpoint := os.Getenv("DYNAMODB_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:8000"
	}

	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "ap-northeast-1"
	}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithBaseEndpoint(endpoint),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummy", "dummy", "")),
	)
	if err != nil {
		return nil, err
	}

	return dynamodb.NewFromConfig(cfg), nil
}

// newAWSClient creates a DynamoDB client for AWS (production/staging)
func newAWSClient(ctx context.Context) (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	return dynamodb.NewFromConfig(cfg), nil
}
