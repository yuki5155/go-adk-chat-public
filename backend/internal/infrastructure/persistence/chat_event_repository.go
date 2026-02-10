package persistence

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/yuki5155/go-google-auth/internal/domain/chat"
	dynamodbInfra "github.com/yuki5155/go-google-auth/internal/infrastructure/dynamodb"
)

// safeInt32 safely converts an int to int32, capping at MaxInt32
func safeInt32(n int) int32 {
	if n > math.MaxInt32 {
		return math.MaxInt32
	}
	if n < math.MinInt32 {
		return math.MinInt32
	}
	return int32(n)
}

// ChatEventRepository implements chat.EventRepository using DynamoDB
type ChatEventRepository struct {
	client    *dynamodb.Client
	tableName string
}

// NewChatEventRepository creates a new DynamoDB-based event repository
func NewChatEventRepository(client *dynamodb.Client) *ChatEventRepository {
	tables := dynamodbInfra.GetTableNames()
	return &ChatEventRepository{
		client:    client,
		tableName: tables.ChatEvents,
	}
}

// eventModel represents the DynamoDB schema for events
type eventModel struct {
	SessionID    string `dynamodbav:"session_id"`
	EventID      string `dynamodbav:"event_id"`
	ThreadID     string `dynamodbav:"thread_id"`
	Role         string `dynamodbav:"role"`
	Content      string `dynamodbav:"content"`
	Author       string `dynamodbav:"author"`
	InvocationID string `dynamodbav:"invocation_id"`
	Timestamp    int64  `dynamodbav:"timestamp"`
}

// Save saves an event
func (r *ChatEventRepository) Save(ctx context.Context, event *chat.Event) error {
	model := r.eventDomainToModel(event)

	item, err := attributevalue.MarshalMap(model)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	}

	_, err = r.client.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to save event: %w", err)
	}

	return nil
}

// ListBySession lists events for a session, ordered by event_id
func (r *ChatEventRepository) ListBySession(ctx context.Context, sessionID string, limit int) ([]*chat.Event, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("session_id = :session_id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":session_id": &types.AttributeValueMemberS{Value: sessionID},
		},
		ScanIndexForward: aws.Bool(true), // Ascending order (oldest first)
	}

	if limit > 0 {
		input.Limit = aws.Int32(safeInt32(limit))
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}

	events := make([]*chat.Event, 0, len(result.Items))
	for _, item := range result.Items {
		var model eventModel
		if err := attributevalue.UnmarshalMap(item, &model); err != nil {
			return nil, fmt.Errorf("failed to unmarshal event: %w", err)
		}
		events = append(events, r.eventModelToDomain(&model))
	}

	return events, nil
}

// DeleteBySession deletes all events for a session
func (r *ChatEventRepository) DeleteBySession(ctx context.Context, sessionID string) error {
	// First, query all events for this session
	events, err := r.ListBySession(ctx, sessionID, 0)
	if err != nil {
		return fmt.Errorf("failed to list events for deletion: %w", err)
	}

	// Batch delete (DynamoDB limits to 25 items per batch)
	for i := 0; i < len(events); i += 25 {
		end := i + 25
		if end > len(events) {
			end = len(events)
		}

		writeRequests := make([]types.WriteRequest, 0, end-i)
		for _, event := range events[i:end] {
			writeRequests = append(writeRequests, types.WriteRequest{
				DeleteRequest: &types.DeleteRequest{
					Key: map[string]types.AttributeValue{
						"session_id": &types.AttributeValueMemberS{Value: sessionID},
						"event_id":   &types.AttributeValueMemberS{Value: event.EventID()},
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
			return fmt.Errorf("failed to batch delete events: %w", err)
		}
	}

	return nil
}

// Model conversion helpers
func (r *ChatEventRepository) eventModelToDomain(model *eventModel) *chat.Event {
	return chat.ReconstructEvent(
		model.EventID,
		model.SessionID,
		model.ThreadID,
		chat.EventRole(model.Role),
		model.Content,
		model.Author,
		model.InvocationID,
		time.Unix(0, model.Timestamp),
	)
}

func (r *ChatEventRepository) eventDomainToModel(event *chat.Event) *eventModel {
	return &eventModel{
		SessionID:    event.SessionID(),
		EventID:      event.EventID(),
		ThreadID:     event.ThreadID(),
		Role:         string(event.Role()),
		Content:      event.Content(),
		Author:       event.Author(),
		InvocationID: event.InvocationID(),
		Timestamp:    event.Timestamp().UnixNano(),
	}
}
