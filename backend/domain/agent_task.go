// agent_task.go
// AgentTask 及其生命周期状态的领域模型定义（Golang 实现）
package domain

import (
	"time"
)

type AgentTaskStatus string

const (
	TaskPending     AgentTaskStatus = "PENDING"     // 待分配
	TaskInProgress  AgentTaskStatus = "IN_PROGRESS" // 进行中
	TaskCompleted   AgentTaskStatus = "COMPLETED"   // 已完成
	TaskFailed      AgentTaskStatus = "FAILED"      // 失败
	TaskInterrupted AgentTaskStatus = "INTERRUPTED" // 中断
	TaskArchived    AgentTaskStatus = "ARCHIVED"    // 已归档
)

type AgentTask struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Status      AgentTaskStatus `json:"status"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	AssignedTo  string          `json:"assigned_to,omitempty"`
	Logs        []string        `json:"logs"`
}

func NewAgentTask(id, name, desc string) *AgentTask {
	now := time.Now().UTC()
	return &AgentTask{
		ID:          id,
		Name:        name,
		Description: desc,
		Status:      TaskPending,
		CreatedAt:   now,
		UpdatedAt:   now,
		Logs:        []string{},
	}
}

func (t *AgentTask) Assign(agentID string) {
	t.AssignedTo = agentID
	t.Status = TaskInProgress
	t.UpdatedAt = time.Now().UTC()
	t.Logs = append(t.Logs, "Assigned to "+agentID+" at "+t.UpdatedAt.String())
}

func (t *AgentTask) Complete() {
	t.Status = TaskCompleted
	t.UpdatedAt = time.Now().UTC()
	t.Logs = append(t.Logs, "Completed at "+t.UpdatedAt.String())
}

func (t *AgentTask) Fail(reason string) {
	t.Status = TaskFailed
	t.UpdatedAt = time.Now().UTC()
	t.Logs = append(t.Logs, "Failed at "+t.UpdatedAt.String()+": "+reason)
}

func (t *AgentTask) Interrupt(reason string) {
	t.Status = TaskInterrupted
	t.UpdatedAt = time.Now().UTC()
	t.Logs = append(t.Logs, "Interrupted at "+t.UpdatedAt.String()+": "+reason)
}

func (t *AgentTask) Archive() {
	t.Status = TaskArchived
	t.UpdatedAt = time.Now().UTC()
	t.Logs = append(t.Logs, "Archived at "+t.UpdatedAt.String())
}
