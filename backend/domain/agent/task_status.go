package agent

type TaskStatus string

const (
	TaskPending    TaskStatus = "PENDING"
	TaskInProgress TaskStatus = "IN_PROGRESS"
	TaskCompleted  TaskStatus = "COMPLETED"
	TaskFailed     TaskStatus = "FAILED"
	TaskInterrupted TaskStatus = "INTERRUPTED"
	TaskArchived   TaskStatus = "ARCHIVED"
)

func (s TaskStatus) IsValid() bool {
	switch s {
	case TaskPending, TaskInProgress, TaskCompleted, TaskFailed, TaskInterrupted, TaskArchived:
		return true
	}
	return false
}

func (s TaskStatus) CanTransitionTo(newStatus TaskStatus) bool {
	switch s {
	case TaskPending:
		return newStatus == TaskInProgress || newStatus == TaskArchived
	case TaskInProgress:
		return newStatus == TaskCompleted || newStatus == TaskFailed || newStatus == TaskInterrupted
	case TaskCompleted, TaskFailed, TaskInterrupted:
		return newStatus == TaskArchived
	case TaskArchived:
		return false
	}
	return false
}
