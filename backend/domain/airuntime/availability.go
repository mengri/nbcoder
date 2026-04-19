package airuntime

import "time"

type ModelStatus string

const (
	ModelAvailable   ModelStatus = "AVAILABLE"
	ModelUnavailable ModelStatus = "UNAVAILABLE"
	ModelDegraded    ModelStatus = "DEGRADED"
)

type ModelHealth struct {
	ModelID   string      `json:"model_id"`
	CheckedAt time.Time   `json:"checked_at"`
	Status    ModelStatus `json:"status"`
	Reason    string      `json:"reason,omitempty"`
}

type AvailabilityChecker interface {
	Check(model *Model) ModelHealth
}

type Degrader interface {
	Degrade(model *Model) error
}
