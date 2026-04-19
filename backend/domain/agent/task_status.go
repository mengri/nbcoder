package agent

type AgentTaskStatus string

const (
	TaskPending     AgentTaskStatus = "PENDING"
	TaskInProgress  AgentTaskStatus = "IN_PROGRESS"
	TaskCompleted   AgentTaskStatus = "COMPLETED"
	TaskFailed      AgentTaskStatus = "FAILED"
	TaskInterrupted  AgentTaskStatus = "INTERRUPTED"
	TaskArchived    AgentTaskStatus = "ARCHIVED"
)

func (s AgentTaskStatus) IsValid() bool {
	switch s {
	case TaskPending, TaskInProgress, TaskCompleted, TaskFailed, TaskInterrupted, TaskArchived:
		return true
	}
	return false
}

func (s AgentTaskStatus) CanTransitionTo(newStatus AgentTaskStatus) bool {
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