package clonepool
// exception.go
// 克隆实例异常检测、恢复与状态标记
import "time"

type ExceptionType string

const (
	ExceptionNone      ExceptionType = "NONE"
	ExceptionTimeout   ExceptionType = "TIMEOUT"
	ExceptionCrash     ExceptionType = "CRASH"
	ExceptionUnknown   ExceptionType = "UNKNOWN"
)

type CloneInstanceException struct {
	InstanceID string
	Type       ExceptionType
	DetectedAt time.Time
	Recovered  bool
	Note       string
}

// 检测异常
func DetectException(inst *CloneInstance) *CloneInstanceException {
	// 示例：假定超时为异常
	if inst.Status == InstanceInUse {
		// 这里可扩展为实际检测逻辑
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

// 自动恢复
func RecoverInstance(inst *CloneInstance, ex *CloneInstanceException) {
	if ex != nil && !ex.Recovered {
		inst.Recycle()
		ex.Recovered = true
		inst.Logs = append(inst.Logs, "Exception recovered at "+time.Now().UTC().String())
	}
}
