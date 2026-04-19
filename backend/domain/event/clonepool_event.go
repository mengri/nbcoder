package event

import "time"

type ClonePoolEventType string

const (
	CloneAcquiredEvent  ClonePoolEventType = "CloneAcquired"
	CloneReleasedEvent  ClonePoolEventType = "CloneReleased"
	CloneBecameDirtyEvent ClonePoolEventType = "CloneBecameDirty"
)

type ClonePoolEvent struct {
	BaseEvent
	InstanceID string            `json:"instance_id"`
	Type       ClonePoolEventType `json:"type"`
}

func (e *ClonePoolEvent) EventType() string {
	return string(e.Type)
}

func (e *ClonePoolEvent) AggregateID() string {
	return e.InstanceID
}

func NewClonePoolEvent(id, instanceID string, eventType ClonePoolEventType) *ClonePoolEvent {
	return &ClonePoolEvent{
		BaseEvent: BaseEvent{
			ID:       id,
			Occurred: time.Now().UTC(),
			Payload:  make(map[string]interface{}),
		},
		InstanceID: instanceID,
		Type:       eventType,
	}
}
