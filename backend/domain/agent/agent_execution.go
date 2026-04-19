package agent

import "time"

type AgentExecution struct {
	ID        string    `json:"id"`
	TaskID    string    `json:"task_id"`
	AgentID   string    `json:"agent_id"`
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
}

type AgentExecutionRepo interface {
	Save(execution *AgentExecution) error
	QueryByTask(taskID string) ([]*AgentExecution, error)
	QueryByAgent(agentID string) ([]*AgentExecution, error)
	QueryByTimeRange(start, end time.Time) ([]*AgentExecution, error)
}
