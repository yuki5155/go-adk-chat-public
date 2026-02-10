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

// ChatSessionRepository implements chat.SessionRepository using DynamoDB
type ChatSessionRepository struct {
	client    *dynamodb.Client
	tableName string
}

// NewChatSessionRepository creates a new DynamoDB-based session repository
func NewChatSessionRepository(client *dynamodb.Client) *ChatSessionRepository {
	tables := dynamodbInfra.GetTableNames()
	return &ChatSessionRepository{
		client:    client,
		tableName: tables.ChatSessions,
	}
}

// sessionModel represents the DynamoDB schema for sessions
type sessionModel struct {
	ThreadID   string `dynamodbav:"thread_id"`
	SessionID  string `dynamodbav:"session_id"`
	UserID     string `dynamodbav:"user_id"`
	AppName    string `dynamodbav:"app_name"`
	State      string `dynamodbav:"state"`
	EventCount int    `dynamodbav:"event_count"`
	CreatedAt  int64  `dynamodbav:"created_at"`
	UpdatedAt  int64  `dynamodbav:"updated_at"`
}

// Save creates or updates a session
func (r *ChatSessionRepository) Save(ctx context.Context, session *chat.Session) error {
	model := r.sessionDomainToModel(session)

	item, err := attributevalue.MarshalMap(model)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	}

	_, err = r.client.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	return nil
}

// FindByID finds a session by thread ID and session ID
func (r *ChatSessionRepository) FindByID(ctx context.Context, threadID, sessionID string) (*chat.Session, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"thread_id":  &types.AttributeValueMemberS{Value: threadID},
			"session_id": &types.AttributeValueMemberS{Value: sessionID},
		},
	}

	result, err := r.client.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	if result.Item == nil {
		return nil, chat.ErrSessionNotFound
	}

	var model sessionModel
	if err := attributevalue.UnmarshalMap(result.Item, &model); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	return r.sessionModelToDomain(&model), nil
}

// FindActiveByThread finds the active session for a thread
func (r *ChatSessionRepository) FindActiveByThread(ctx context.Context, threadID string) (*chat.Session, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("thread_id = :thread_id"),
		FilterExpression:       aws.String("#state = :state"),
		ExpressionAttributeNames: map[string]string{
			"#state": "state",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":thread_id": &types.AttributeValueMemberS{Value: threadID},
			":state":     &types.AttributeValueMemberS{Value: string(chat.SessionStateActive)},
		},
		ScanIndexForward: aws.Bool(false), // Get most recent first
		Limit:            aws.Int32(1),
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to find active session: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, nil // No active session found
	}

	var model sessionModel
	if err := attributevalue.UnmarshalMap(result.Items[0], &model); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	return r.sessionModelToDomain(&model), nil
}

// ListByThread lists sessions for a thread
func (r *ChatSessionRepository) ListByThread(ctx context.Context, threadID string) ([]*chat.Session, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("thread_id = :thread_id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":thread_id": &types.AttributeValueMemberS{Value: threadID},
		},
		ScanIndexForward: aws.Bool(false), // Descending order
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}

	sessions := make([]*chat.Session, 0, len(result.Items))
	for _, item := range result.Items {
		var model sessionModel
		if err := attributevalue.UnmarshalMap(item, &model); err != nil {
			return nil, fmt.Errorf("failed to unmarshal session: %w", err)
		}
		sessions = append(sessions, r.sessionModelToDomain(&model))
	}

	return sessions, nil
}

// Model conversion helpers
func (r *ChatSessionRepository) sessionModelToDomain(model *sessionModel) *chat.Session {
	return chat.ReconstructSession(
		model.SessionID,
		model.ThreadID,
		model.UserID,
		model.AppName,
		chat.SessionState(model.State),
		model.EventCount,
		time.Unix(model.CreatedAt, 0),
		time.Unix(model.UpdatedAt, 0),
	)
}

func (r *ChatSessionRepository) sessionDomainToModel(session *chat.Session) *sessionModel {
	return &sessionModel{
		ThreadID:   session.ThreadID(),
		SessionID:  session.SessionID(),
		UserID:     session.UserID(),
		AppName:    session.AppName(),
		State:      string(session.State()),
		EventCount: session.EventCount(),
		CreatedAt:  session.CreatedAt().Unix(),
		UpdatedAt:  session.UpdatedAt().Unix(),
	}
}
