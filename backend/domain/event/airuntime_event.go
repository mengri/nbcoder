package event

import "time"

type AIRuntimeEventType string

const (
	ModelCalledEvent   AIRuntimeEventType = "ModelCalled"
	ModelFailedEvent   AIRuntimeEventType = "ModelFailed"
	ModelSwitchedEvent AIRuntimeEventType = "ModelSwitched"
	ModelDegradedEvent AIRuntimeEventType = "ModelDegraded"
)

type AIRuntimeEvent struct {
	BaseEvent
	ModelID string            `json:"model_id"`
	Type    AIRuntimeEventType `json:"type"`
}

func (e *AIRuntimeEvent) EventType() string {
	return string(e.Type)
}

func (e *AIRuntimeEvent) AggregateID() string {
	return e.ModelID
}

func NewAIRuntimeEvent(id, modelID string, eventType AIRuntimeEventType) *AIRuntimeEvent {
	return &AIRuntimeEvent{
		BaseEvent: BaseEvent{
			ID:       id,
			Occurred: time.Now().UTC(),
			Payload:  make(map[string]interface{}),
		},
		ModelID: modelID,
		Type:    eventType,
	}
}
