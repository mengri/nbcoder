package domain

// agent_task_log.go
// AgentTaskLog 任务执行日志与追溯领域模型

import (
	"time"
)

type AgentTaskLog struct {
	ID        string    `json:"id"`
	TaskID    string    `json:"task_id"`
	AgentID   string    `json:"agent_id"`
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"` // info, warn, error, debug
	Message   string    `json:"message"`
}

type AgentTaskLogRepo interface {
	Save(log *AgentTaskLog) error
	QueryByTask(taskID string) ([]*AgentTaskLog, error)
	QueryByAgent(agentID string) ([]*AgentTaskLog, error)
	QueryByTimeRange(start, end time.Time) ([]*AgentTaskLog, error)
}
