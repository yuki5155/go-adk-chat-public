//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	region = "ap-northeast-1"
)

func getEndpoint() string {
	if endpoint := os.Getenv("DYNAMODB_ENDPOINT"); endpoint != "" {
		return endpoint
	}
	return "http://localhost:8000"
}

func getProjectName() string {
	if name := os.Getenv("PROJECT_NAME"); name != "" {
		return name
	}
	return "go-adk-chat"
}

func getEnvironment() string {
	if env := os.Getenv("GO_ENV"); env != "" {
		return env
	}
	return "dev"
}

func main() {
	ctx := context.Background()
	endpoint := getEndpoint()
	projectName := getProjectName()
	environment := getEnvironment()

	// Configure AWS SDK for DynamoDB Local
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           endpoint,
					SigningRegion: region,
				}, nil
			}),
		),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummy", "dummy", "")),
	)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	client := dynamodb.NewFromConfig(cfg)

	fmt.Println("=== Creating DynamoDB Local Tables ===")
	fmt.Printf("Endpoint: %s\n", endpoint)
	fmt.Printf("Project: %s\n", projectName)
	fmt.Printf("Environment: %s\n\n", environment)

	// Create user_roles table
	if err := createUserRolesTable(ctx, client, projectName, environment); err != nil {
		log.Printf("Warning: Could not create user_roles table: %v", err)
	}

	// Create role_requests table
	if err := createRoleRequestsTable(ctx, client, projectName, environment); err != nil {
		log.Printf("Warning: Could not create role_requests table: %v", err)
	}

	// List tables to verify
	fmt.Println("\n=== Listing Tables ===")
	listTables(ctx, client)

	fmt.Println("\n✅ Table creation completed!")
}

func createUserRolesTable(ctx context.Context, client *dynamodb.Client, projectName, environment string) error {
	tableName := fmt.Sprintf("%s-%s-user-roles", projectName, environment)

	fmt.Printf("Creating table: %s\n", tableName)

	input := &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("user_id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("email"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("role"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("user_id"),
				KeyType:       types.KeyTypeHash,
			},
		},
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String("email-index"),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("email"),
						KeyType:       types.KeyTypeHash,
					},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll,
				},
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(5),
					WriteCapacityUnits: aws.Int64(5),
				},
			},
			{
				IndexName: aws.String("role-index"),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("role"),
						KeyType:       types.KeyTypeHash,
					},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll,
				},
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(5),
					WriteCapacityUnits: aws.Int64(5),
				},
			},
		},
		BillingMode: types.BillingModeProvisioned,
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		StreamSpecification: &types.StreamSpecification{
			StreamEnabled:  aws.Bool(true),
			StreamViewType: types.StreamViewTypeNewAndOldImages,
		},
	}

	_, err := client.CreateTable(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	fmt.Printf("✓ Created table: %s\n", tableName)
	return nil
}

func createRoleRequestsTable(ctx context.Context, client *dynamodb.Client, projectName, environment string) error {
	tableName := fmt.Sprintf("%s-%s-role-requests", projectName, environment)

	fmt.Printf("Creating table: %s\n", tableName)

	input := &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("request_id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("user_id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("status"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("requested_at"),
				AttributeType: types.ScalarAttributeTypeN,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("request_id"),
				KeyType:       types.KeyTypeHash,
			},
		},
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String("user_id-index"),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("user_id"),
						KeyType:       types.KeyTypeHash,
					},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll,
				},
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(5),
					WriteCapacityUnits: aws.Int64(5),
				},
			},
			{
				IndexName: aws.String("status-requested_at-index"),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("status"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("requested_at"),
						KeyType:       types.KeyTypeRange,
					},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll,
				},
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(5),
					WriteCapacityUnits: aws.Int64(5),
				},
			},
		},
		BillingMode: types.BillingModeProvisioned,
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		StreamSpecification: &types.StreamSpecification{
			StreamEnabled:  aws.Bool(true),
			StreamViewType: types.StreamViewTypeNewAndOldImages,
		},
	}

	_, err := client.CreateTable(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	fmt.Printf("✓ Created table: %s\n", tableName)
	return nil
}

func listTables(ctx context.Context, client *dynamodb.Client) {
	output, err := client.ListTables(ctx, &dynamodb.ListTablesInput{})
	if err != nil {
		log.Printf("Failed to list tables: %v", err)
		return
	}

	for _, table := range output.TableNames {
		fmt.Printf("  - %s\n", table)
	}
}
