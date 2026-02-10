package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/yuki5155/go-google-auth/internal/domain/chat"
	dynamodbInfra "github.com/yuki5155/go-google-auth/internal/infrastructure/dynamodb"
)

// ChatThreadRepository implements chat.ThreadRepository using DynamoDB
type ChatThreadRepository struct {
	client    *dynamodb.Client
	tableName string
}

// NewChatThreadRepository creates a new DynamoDB-based thread repository
func NewChatThreadRepository(client *dynamodb.Client) *ChatThreadRepository {
	tables := dynamodbInfra.GetTableNames()
	return &ChatThreadRepository{
		client:    client,
		tableName: tables.ChatThreads,
	}
}

// threadModel represents the DynamoDB schema for threads
type threadModel struct {
	UserID       string `dynamodbav:"user_id"`
	ThreadID     string `dynamodbav:"thread_id"`
	Title        string `dynamodbav:"title"`
	Model        string `dynamodbav:"model"`
	Status       string `dynamodbav:"status"`
	MessageCount int    `dynamodbav:"message_count"`
	LastMessage  string `dynamodbav:"last_message"`
	CreatedAt    int64  `dynamodbav:"created_at"`
	UpdatedAt    int64  `dynamodbav:"updated_at"`
}

// Save creates or updates a thread
func (r *ChatThreadRepository) Save(ctx context.Context, thread *chat.Thread) error {
	model := r.threadDomainToModel(thread)

	item, err := attributevalue.MarshalMap(model)
	if err != nil {
		return fmt.Errorf("failed to marshal thread: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	}

	_, err = r.client.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to save thread: %w", err)
	}

	return nil
}

// FindByID finds a thread by user ID and thread ID
func (r *ChatThreadRepository) FindByID(ctx context.Context, userID, threadID string) (*chat.Thread, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"user_id":   &types.AttributeValueMemberS{Value: userID},
			"thread_id": &types.AttributeValueMemberS{Value: threadID},
		},
	}

	result, err := r.client.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get thread: %w", err)
	}

	if result.Item == nil {
		return nil, chat.ErrThreadNotFound
	}

	var model threadModel
	if err := attributevalue.UnmarshalMap(result.Item, &model); err != nil {
		return nil, fmt.Errorf("failed to unmarshal thread: %w", err)
	}

	return r.threadModelToDomain(&model), nil
}

// ListByUser lists threads for a user, ordered by updated_at desc
func (r *ChatThreadRepository) ListByUser(ctx context.Context, userID string, limit int, lastKey string) ([]*chat.Thread, string, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("thread-updated-index"),
		KeyConditionExpression: aws.String("user_id = :user_id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":user_id": &types.AttributeValueMemberS{Value: userID},
		},
		ScanIndexForward: aws.Bool(false), // Descending order (newest first)
		Limit:            aws.Int32(safeInt32(limit)),
	}

	// Handle pagination
	if lastKey != "" {
		input.ExclusiveStartKey = map[string]types.AttributeValue{
			"user_id":    &types.AttributeValueMemberS{Value: userID},
			"updated_at": &types.AttributeValueMemberN{Value: lastKey},
		}
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		return nil, "", fmt.Errorf("failed to list threads: %w", err)
	}

	threads := make([]*chat.Thread, 0, len(result.Items))
	for _, item := range result.Items {
		var model threadModel
		if err := attributevalue.UnmarshalMap(item, &model); err != nil {
			return nil, "", fmt.Errorf("failed to unmarshal thread: %w", err)
		}
		threads = append(threads, r.threadModelToDomain(&model))
	}

	// Get next page key
	var nextKey string
	if result.LastEvaluatedKey != nil {
		if updatedAt, ok := result.LastEvaluatedKey["updated_at"].(*types.AttributeValueMemberN); ok {
			nextKey = updatedAt.Value
		}
	}

	return threads, nextKey, nil
}

// Delete deletes a thread
func (r *ChatThreadRepository) Delete(ctx context.Context, userID, threadID string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"user_id":   &types.AttributeValueMemberS{Value: userID},
			"thread_id": &types.AttributeValueMemberS{Value: threadID},
		},
	}

	_, err := r.client.DeleteItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete thread: %w", err)
	}

	return nil
}

// Model conversion helpers
func (r *ChatThreadRepository) threadModelToDomain(model *threadModel) *chat.Thread {
	return chat.ReconstructThread(
		model.ThreadID,
		model.UserID,
		model.Title,
		model.Model,
		chat.ThreadStatus(model.Status),
		model.MessageCount,
		model.LastMessage,
		time.Unix(model.CreatedAt, 0),
		time.Unix(model.UpdatedAt, 0),
	)
}

func (r *ChatThreadRepository) threadDomainToModel(thread *chat.Thread) *threadModel {
	return &threadModel{
		UserID:       thread.UserID(),
		ThreadID:     thread.ThreadID(),
		Title:        thread.Title(),
		Model:        thread.Model(),
		Status:       string(thread.Status()),
		MessageCount: thread.MessageCount(),
		LastMessage:  thread.LastMessage(),
		CreatedAt:    thread.CreatedAt().Unix(),
		UpdatedAt:    thread.UpdatedAt().Unix(),
	}
}
