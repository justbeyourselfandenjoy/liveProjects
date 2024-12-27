package kafka_helpers

import (
	"time"

	"github.com/google/uuid"
)

// BaseEvent represents common properties of an event
type BaseEvent struct {
	EventID        uuid.UUID
	EventTimestamp time.Time
	EventName      string
	EventBody      string
}

func BuildBaseEvent(name string, body string) *BaseEvent {
	return &BaseEvent{
		EventID:        uuid.New(),
		EventTimestamp: time.Now(),
		EventName:      name,
		EventBody:      body,
	}
}
