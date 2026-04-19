package agent

import (
	"time"
)

type AgentExecution struct {
	ID            string                 `json:"id"`
	AgentID       string                 `json:"agent_id"`
	TaskID        string                 `json:"task_id"`
	StartTime     time.Time              `json:"start_time"`
	EndTime       *time.Time             `json:"end_time,omitempty"`
	Status        string                 `json:"status"`
	Input         map[string]interface{} `json:"input,omitempty"`
	Output        map[string]interface{} `json:"output,omitempty"`
	Error         string                 `json:"error,omitempty"`
	ModelUsed     string                 `json:"model_used,omitempty"`
	TokensUsed    int                    `json:"tokens_used,omitempty"`
	Timestamp     time.Time              `json:"timestamp"`
}

func NewAgentExecution(id, agentID, taskID string) *AgentExecution {
	now := time.Now().UTC()
	return &AgentExecution{
		ID:        id,
		AgentID:   agentID,
		TaskID:    taskID,
		StartTime: now,
		Status:    "IN_PROGRESS",
		Input:     make(map[string]interface{}),
		Output:    make(map[string]interface{}),
		Timestamp: now,
	}
}

func (e *AgentExecution) Complete(output map[string]interface{}, modelUsed string, tokensUsed int) {
	now := time.Now().UTC()
	e.EndTime = &now
	e.Status = "COMPLETED"
	e.Output = output
	e.ModelUsed = modelUsed
	e.TokensUsed = tokensUsed
}

func (e *AgentExecution) Fail(errorMessage string) {
	now := time.Now().UTC()
	e.EndTime = &now
	e.Status = "FAILED"
	e.Error = errorMessage
}

func (e *AgentExecution) GetDuration() time.Duration {
	if e.EndTime == nil {
		return 0
	}
	return e.EndTime.Sub(e.StartTime)
}
