package agent

import (
	"fmt"
	"time"
)

type Task struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	AgentType   AgentType       `json:"agent_type"`
	TaskType    string          `json:"task_type"`
	Status      TaskStatus `json:"status"`
	Priority    int             `json:"priority"`
	AssignedTo  string          `json:"assigned_to,omitempty"`
	PipelineID  string          `json:"pipeline_id,omitempty"`
	ProjectName string          `json:"projectName"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	StartedAt   *time.Time      `json:"started_at,omitempty"`
	CompletedAt *time.Time      `json:"completed_at,omitempty"`
	Context     map[string]interface{} `json:"context,omitempty"`
}

func NewTask(id, name, description, taskType string, agentType AgentType, projectName string) *Task {
	now := time.Now().UTC()
	return &Task{
		ID:          id,
		Name:        name,
		Description: description,
		AgentType:   agentType,
		TaskType:    taskType,
		Status:      TaskPending,
		Priority:    5,
		ProjectName: projectName,
		CreatedAt:   now,
		UpdatedAt:   now,
		Context:     make(map[string]interface{}),
	}
}

func (t *Task) Assign(agentID string) error {
	if !t.Status.CanTransitionTo(TaskInProgress) {
		return fmt.Errorf("cannot assign task in status %s", t.Status)
	}

	now := time.Now().UTC()
	t.AssignedTo = agentID
	t.Status = TaskInProgress
	t.StartedAt = &now
	t.UpdatedAt = now
	return nil
}

func (t *Task) Complete() error {
	if !t.Status.CanTransitionTo(TaskCompleted) {
		return fmt.Errorf("cannot complete task in status %s", t.Status)
	}

	now := time.Now().UTC()
	t.Status = TaskCompleted
	t.CompletedAt = &now
	t.UpdatedAt = now
	return nil
}

func (t *Task) Fail(reason string) error {
	if !t.Status.CanTransitionTo(TaskFailed) {
		return fmt.Errorf("cannot fail task in status %s", t.Status)
	}

	now := time.Now().UTC()
	t.Status = TaskFailed
	t.CompletedAt = &now
	t.UpdatedAt = now
	return nil
}

func (t *Task) Interrupt(reason string) error {
	if !t.Status.CanTransitionTo(TaskInterrupted) {
		return fmt.Errorf("cannot interrupt task in status %s", t.Status)
	}

	now := time.Now().UTC()
	t.Status = TaskInterrupted
	t.CompletedAt = &now
	t.UpdatedAt = now
	return nil
}

func (t *Task) Archive() error {
	if !t.Status.CanTransitionTo(TaskArchived) {
		return fmt.Errorf("cannot archive task in status %s", t.Status)
	}

	now := time.Now().UTC()
	t.Status = TaskArchived
	t.UpdatedAt = now
	return nil
}

func (t *Task) UpdateContext(key string, value interface{}) {
	if t.Context == nil {
		t.Context = make(map[string]interface{})
	}
	t.Context[key] = value
	t.UpdatedAt = time.Now().UTC()
}

func (t *Task) SetPriority(priority int) {
	t.Priority = priority
	t.UpdatedAt = time.Now().UTC()
}

func (t *Task) SetPipelineID(pipelineID string) {
	t.PipelineID = pipelineID
	t.UpdatedAt = time.Now().UTC()
}

func (t *Task) GetDuration() time.Duration {
	if t.StartedAt == nil {
		return 0
	}

	endTime := time.Now().UTC()
	if t.CompletedAt != nil {
		endTime = *t.CompletedAt
	}

	return endTime.Sub(*t.StartedAt)
}
