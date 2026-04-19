package event

import "time"

type DomainEvent interface {
	EventType() string
	OccurredAt() time.Time
	AggregateID() string
}

type BaseEvent struct {
	ID        string    `json:"id"`
	Occurred  time.Time `json:"occurred"`
	Payload   map[string]interface{} `json:"payload,omitempty"`
}

func (e *BaseEvent) OccurredAt() time.Time {
	return e.Occurred
}
