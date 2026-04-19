package agent

type TaskStatus string

const (
	TaskPending     TaskStatus = "PENDING"
	TaskInProgress  TaskStatus = "IN_PROGRESS"
	TaskCompleted   TaskStatus = "COMPLETED"
	TaskFailed      TaskStatus = "FAILED"
	TaskInterrupted TaskStatus = "INTERRUPTED"
	TaskArchived    TaskStatus = "ARCHIVED"
)

func (s TaskStatus) IsValid() bool {
	switch s {
	case TaskPending, TaskInProgress, TaskCompleted, TaskFailed, TaskInterrupted, TaskArchived:
		return true
	}
	return false
}
