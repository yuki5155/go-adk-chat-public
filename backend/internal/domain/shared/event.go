package shared

import "time"

// DomainEvent represents a domain event that occurs in the system
type DomainEvent interface {
	// EventType returns the type of the event
	EventType() string
	// OccurredAt returns when the event occurred
	OccurredAt() time.Time
	// AggregateID returns the ID of the aggregate that generated this event
	AggregateID() string
}

// BaseDomainEvent provides common fields for all domain events
type BaseDomainEvent struct {
	eventType   string
	occurredAt  time.Time
	aggregateID string
}

// NewBaseDomainEvent creates a new base domain event
func NewBaseDomainEvent(eventType, aggregateID string) BaseDomainEvent {
	return BaseDomainEvent{
		eventType:   eventType,
		occurredAt:  time.Now(),
		aggregateID: aggregateID,
	}
}

// EventType returns the type of the event
func (e BaseDomainEvent) EventType() string {
	return e.eventType
}

// OccurredAt returns when the event occurred
func (e BaseDomainEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// AggregateID returns the ID of the aggregate that generated this event
func (e BaseDomainEvent) AggregateID() string {
	return e.aggregateID
}
