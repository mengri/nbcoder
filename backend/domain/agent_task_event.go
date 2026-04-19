package domain
// agent_task_event.go
// AgentTaskEvent 领域事件发布与订阅


import "time"

type AgentTaskEventType string

const (
	TaskAssigned   AgentTaskEventType = "TaskAssigned"
	TaskStarted    AgentTaskEventType = "TaskStarted"
	TaskCompleted  AgentTaskEventType = "TaskCompleted"
	TaskFailed     AgentTaskEventType = "TaskFailed"
	TaskInterrupted AgentTaskEventType = "TaskInterrupted"
)

type AgentTaskEvent struct {
	ID        string             `json:"id"`
	TaskID    string             `json:"task_id"`
	AgentID   string             `json:"agent_id"`
	Occurred  time.Time          `json:"occurred"`
	Type      AgentTaskEventType `json:"type"`
	Payload   map[string]interface{} `json:"payload,omitempty"`
}

type AgentTaskEventPublisher interface {
	Publish(event *AgentTaskEvent) error
}

type AgentTaskEventSubscriber interface {
	OnEvent(event *AgentTaskEvent)
}
