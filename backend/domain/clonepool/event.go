package clonepool
// event.go
// 克隆池领域事件发布与订阅
package clonepool

import "time"

type ClonePoolEventType string

const (
	CloneAcquired ClonePoolEventType = "CloneAcquired"
	CloneReleased ClonePoolEventType = "CloneReleased"
)

type ClonePoolEvent struct {
	ID         string            `json:"id"`
	InstanceID string            `json:"instance_id"`
	Occurred   time.Time         `json:"occurred"`
	Type       ClonePoolEventType `json:"type"`
	Payload    map[string]interface{} `json:"payload,omitempty"`
}

type ClonePoolEventPublisher interface {
	Publish(event *ClonePoolEvent) error
}

type ClonePoolEventSubscriber interface {
	OnEvent(event *ClonePoolEvent)
}
