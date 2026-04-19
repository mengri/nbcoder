package event

import "time"

type PipelineEventType string

const (
	StageStartedEvent       PipelineEventType = "StageStarted"
	StageCompletedEvent     PipelineEventType = "StageCompleted"
	StageFailedEvent        PipelineEventType = "StageFailed"
	StageReviewRequiredEvent PipelineEventType = "StageReviewRequired"
)

type PipelineEvent struct {
	BaseEvent
	StageID string            `json:"stage_id"`
	Type    PipelineEventType `json:"type"`
}

func (e *PipelineEvent) EventType() string {
	return string(e.Type)
}

func (e *PipelineEvent) AggregateID() string {
	return e.StageID
}

func NewPipelineEvent(id, stageID string, eventType PipelineEventType) *PipelineEvent {
	return &PipelineEvent{
		BaseEvent: BaseEvent{
			ID:       id,
			Occurred: time.Now().UTC(),
			Payload:  make(map[string]interface{}),
		},
		StageID: stageID,
		Type:    eventType,
	}
}
