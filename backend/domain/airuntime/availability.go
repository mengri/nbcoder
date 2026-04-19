package airuntime
// availability.go
// Provider/模型可用性检测与降级
package airuntime

import "time"

type ModelStatus string

const (
	ModelAvailable   ModelStatus = "AVAILABLE"
	ModelUnavailable ModelStatus = "UNAVAILABLE"
	ModelDegraded    ModelStatus = "DEGRADED"
)

type ModelHealth struct {
	ModelID    string      `json:"model_id"`
	CheckedAt  time.Time   `json:"checked_at"`
	Status     ModelStatus `json:"status"`
	Reason     string      `json:"reason,omitempty"`
}

type AvailabilityChecker interface {
	Check(model *Model) ModelHealth
}

type Degrader interface {
	Degrade(model *Model) error
}

// 示例：简单可用性检测与降级
func SimpleCheck(model *Model) ModelHealth {
	// 这里可扩展为实际探测逻辑
	return ModelHealth{
		ModelID:   model.ID,
		CheckedAt: time.Now().UTC(),
		Status:    ModelAvailable,
	}
}

func SimpleDegrade(model *Model) error {
	// 降级逻辑，如切换备用模型
	return nil
}
