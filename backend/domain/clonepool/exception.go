package clonepool

import "time"

type ExceptionType string

const (
	ExceptionNone    ExceptionType = "NONE"
	ExceptionTimeout ExceptionType = "TIMEOUT"
	ExceptionCrash   ExceptionType = "CRASH"
	ExceptionUnknown ExceptionType = "UNKNOWN"
)

type CloneInstanceException struct {
	InstanceID string
	Type       ExceptionType
	DetectedAt time.Time
	Recovered  bool
	Note       string
}

func DetectException(inst *CloneInstance) *CloneInstanceException {
	if inst.Status == InstanceBusy {
		return &CloneInstanceException{
			InstanceID: inst.ID,
			Type:       ExceptionTimeout,
			DetectedAt: time.Now().UTC(),
			Recovered:  false,
			Note:       "Timeout detected",
		}
	}
	return nil
}

func RecoverInstance(inst *CloneInstance, ex *CloneInstanceException) {
	if ex != nil && !ex.Recovered {
		inst.Recycle()
		ex.Recovered = true
	}
}
