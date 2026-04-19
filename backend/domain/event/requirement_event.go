package event

import "time"

type RequirementEventType string

const (
	CardCreatedEvent    RequirementEventType = "CardCreated"
	CardConfirmedEvent  RequirementEventType = "CardConfirmed"
	CardSupersededEvent RequirementEventType = "CardSuperseded"
	CardAbandonedEvent  RequirementEventType = "CardAbandoned"
)

type RequirementEvent struct {
	BaseEvent
	CardID string               `json:"card_id"`
	Type   RequirementEventType `json:"type"`
}

func (e *RequirementEvent) EventType() string {
	return string(e.Type)
}

func (e *RequirementEvent) AggregateID() string {
	return e.CardID
}

func NewRequirementEvent(id, cardID string, eventType RequirementEventType) *RequirementEvent {
	return &RequirementEvent{
		BaseEvent: BaseEvent{
			ID:       id,
			Occurred: time.Now().UTC(),
			Payload:  make(map[string]interface{}),
		},
		CardID: cardID,
		Type:   eventType,
	}
}
