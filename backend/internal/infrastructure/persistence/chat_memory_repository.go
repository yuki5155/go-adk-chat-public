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

// ChatMemoryRepository implements chat.MemoryRepository using DynamoDB
type ChatMemoryRepository struct {
	client    *dynamodb.Client
	tableName string
}

// NewChatMemoryRepository creates a new DynamoDB-based memory repository
func NewChatMemoryRepository(client *dynamodb.Client) *ChatMemoryRepository {
	tables := dynamodbInfra.GetTableNames()
	return &ChatMemoryRepository{
		client:    client,
		tableName: tables.ChatMemories,
	}
}

// memoryModel represents the DynamoDB schema for memories
type memoryModel struct {
	ThreadID        string   `dynamodbav:"thread_id"`
	MemoryID        string   `dynamodbav:"memory_id"`
	UserID          string   `dynamodbav:"user_id"`
	Content         string   `dynamodbav:"content"`
	Keywords        []string `dynamodbav:"keywords"`
	SourceSessionID string   `dynamodbav:"source_session_id"`
	Importance      int      `dynamodbav:"importance"`
	Timestamp       int64    `dynamodbav:"timestamp"`
}

// Save saves a memory
func (r *ChatMemoryRepository) Save(ctx context.Context, memory *chat.Memory) error {
	model := r.memoryDomainToModel(memory)

	item, err := attributevalue.MarshalMap(model)
	if err != nil {
		return fmt.Errorf("failed to marshal memory: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	}

	_, err = r.client.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to save memory: %w", err)
	}

	return nil
}

// FindByID finds a memory by thread ID and memory ID
func (r *ChatMemoryRepository) FindByID(ctx context.Context, threadID, memoryID string) (*chat.Memory, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"thread_id": &types.AttributeValueMemberS{Value: threadID},
			"memory_id": &types.AttributeValueMemberS{Value: memoryID},
		},
	}

	result, err := r.client.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get memory: %w", err)
	}

	if result.Item == nil {
		return nil, chat.ErrMemoryNotFound
	}

	var model memoryModel
	if err := attributevalue.UnmarshalMap(result.Item, &model); err != nil {
		return nil, fmt.Errorf("failed to unmarshal memory: %w", err)
	}

	return r.memoryModelToDomain(&model), nil
}

// ListByThread lists memories for a thread
func (r *ChatMemoryRepository) ListByThread(ctx context.Context, threadID string, limit int) ([]*chat.Memory, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("thread_id = :thread_id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":thread_id": &types.AttributeValueMemberS{Value: threadID},
		},
		ScanIndexForward: aws.Bool(false), // Descending order (newest first)
	}

	if limit > 0 {
		input.Limit = aws.Int32(safeInt32(limit))
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list memories: %w", err)
	}

	memories := make([]*chat.Memory, 0, len(result.Items))
	for _, item := range result.Items {
		var model memoryModel
		if err := attributevalue.UnmarshalMap(item, &model); err != nil {
			return nil, fmt.Errorf("failed to unmarshal memory: %w", err)
		}
		memories = append(memories, r.memoryModelToDomain(&model))
	}

	return memories, nil
}

// ListByUser lists memories for a user using GSI
func (r *ChatMemoryRepository) ListByUser(ctx context.Context, userID string, limit int) ([]*chat.Memory, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("user-memories-index"),
		KeyConditionExpression: aws.String("user_id = :user_id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":user_id": &types.AttributeValueMemberS{Value: userID},
		},
		ScanIndexForward: aws.Bool(false), // Descending order (newest first)
	}

	if limit > 0 {
		input.Limit = aws.Int32(safeInt32(limit))
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list user memories: %w", err)
	}

	memories := make([]*chat.Memory, 0, len(result.Items))
	for _, item := range result.Items {
		var model memoryModel
		if err := attributevalue.UnmarshalMap(item, &model); err != nil {
			return nil, fmt.Errorf("failed to unmarshal memory: %w", err)
		}
		memories = append(memories, r.memoryModelToDomain(&model))
	}

	return memories, nil
}

// DeleteByThread deletes all memories for a thread
func (r *ChatMemoryRepository) DeleteByThread(ctx context.Context, threadID string) error {
	// First, query all memories for this thread
	memories, err := r.ListByThread(ctx, threadID, 0)
	if err != nil {
		return fmt.Errorf("failed to list memories for deletion: %w", err)
	}

	// Batch delete (DynamoDB limits to 25 items per batch)
	for i := 0; i < len(memories); i += 25 {
		end := i + 25
		if end > len(memories) {
			end = len(memories)
		}

		writeRequests := make([]types.WriteRequest, 0, end-i)
		for _, memory := range memories[i:end] {
			writeRequests = append(writeRequests, types.WriteRequest{
				DeleteRequest: &types.DeleteRequest{
					Key: map[string]types.AttributeValue{
						"thread_id": &types.AttributeValueMemberS{Value: threadID},
						"memory_id": &types.AttributeValueMemberS{Value: memory.MemoryID()},
					},
				},
			})
		}

		input := &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				r.tableName: writeRequests,
			},
		}

		_, err := r.client.BatchWriteItem(ctx, input)
		if err != nil {
			return fmt.Errorf("failed to batch delete memories: %w", err)
		}
	}

	return nil
}

// Model conversion helpers
func (r *ChatMemoryRepository) memoryModelToDomain(model *memoryModel) *chat.Memory {
	return chat.ReconstructMemory(
		model.MemoryID,
		model.ThreadID,
		model.UserID,
		model.Content,
		model.Keywords,
		model.SourceSessionID,
		chat.ImportanceLevel(model.Importance),
		time.Unix(0, model.Timestamp),
	)
}

func (r *ChatMemoryRepository) memoryDomainToModel(memory *chat.Memory) *memoryModel {
	return &memoryModel{
		ThreadID:        memory.ThreadID(),
		MemoryID:        memory.MemoryID(),
		UserID:          memory.UserID(),
		Content:         memory.Content(),
		Keywords:        memory.Keywords(),
		SourceSessionID: memory.SourceSessionID(),
		Importance:      int(memory.Importance()),
		Timestamp:       memory.Timestamp().UnixNano(),
	}
}
