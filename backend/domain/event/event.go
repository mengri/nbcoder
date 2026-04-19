package event
// event.go
// 领域事件类型与发布订阅接口（无依赖，供各聚合使用）
package event

import "time"

type ClonePoolEventType string

const (
	CloneAcquired ClonePoolEventType = "CloneAcquired"
	CloneReleased ClonePoolEventType = "CloneReleased"
)

type ClonePoolEvent struct {
	ID         string                 `json:"id"`
	InstanceID string                 `json:"instance_id"`
	Occurred   time.Time              `json:"occurred"`
	Type       ClonePoolEventType     `json:"type"`
	Payload    map[string]interface{} `json:"payload,omitempty"`
}

type ClonePoolEventPublisher interface {
	Publish(event *ClonePoolEvent) error
}

type ClonePoolEventSubscriber interface {
	OnEvent(event *ClonePoolEvent)
}
