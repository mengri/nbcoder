package event

import "time"

type AgentTaskEventType string

const (
	TaskAssignedEvent  AgentTaskEventType = "TaskAssigned"
	TaskStartedEvent   AgentTaskEventType = "TaskStarted"
	TaskCompletedEvent AgentTaskEventType = "TaskCompleted"
	TaskFailedEvent    AgentTaskEventType = "TaskFailed"
	TaskInterruptedEvent AgentTaskEventType = "TaskInterrupted"
	TaskArchivedEvent  AgentTaskEventType = "TaskArchived"
)

type AgentTaskEvent struct {
	BaseEvent
	TaskID string             `json:"task_id"`
	Type   AgentTaskEventType `json:"type"`
}

func (e *AgentTaskEvent) EventType() string {
	return string(e.Type)
}

func (e *AgentTaskEvent) AggregateID() string {
	return e.TaskID
}

func NewAgentTaskEvent(id, taskID string, eventType AgentTaskEventType) *AgentTaskEvent {
	return &AgentTaskEvent{
		BaseEvent: BaseEvent{
			ID:       id,
			Occurred: time.Now().UTC(),
			Payload:  make(map[string]interface{}),
		},
		TaskID: taskID,
		Type:   eventType,
	}
}