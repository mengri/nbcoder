package airuntime
// event.go
// AI运行时领域事件发布与订阅
package airuntime

import "time"

type AIRuntimeEventType string

const (
	ModelCalled    AIRuntimeEventType = "ModelCalled"
	ModelSwitched  AIRuntimeEventType = "ModelSwitched"
	ModelDegraded  AIRuntimeEventType = "ModelDegraded"
)

type AIRuntimeEvent struct {
	ID        string            `json:"id"`
	ModelID   string            `json:"model_id"`
	Occurred  time.Time         `json:"occurred"`
	Type      AIRuntimeEventType `json:"type"`
	Payload   map[string]interface{} `json:"payload,omitempty"`
}

type AIRuntimeEventPublisher interface {
	Publish(event *AIRuntimeEvent) error
}

type AIRuntimeEventSubscriber interface {
	OnEvent(event *AIRuntimeEvent)
}
