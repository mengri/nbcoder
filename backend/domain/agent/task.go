package agent

import (
	"fmt"
	"time"
)

type Task struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	AssignedTo  string     `json:"assigned_to,omitempty"`
}

func NewTask(id, name, desc string) *Task {
	now := time.Now().UTC()
	return &Task{
		ID:          id,
		Name:        name,
		Description: desc,
		Status:      TaskPending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (t *Task) Assign(agentID string) error {
	if t.Status != TaskPending {
		return fmt.Errorf("cannot assign task in status %s", t.Status)
	}
	t.AssignedTo = agentID
	t.Status = TaskInProgress
	t.UpdatedAt = time.Now().UTC()
	return nil
}

func (t *Task) Complete() error {
	if t.Status != TaskInProgress {
		return fmt.Errorf("cannot complete task in status %s", t.Status)
	}
	t.Status = TaskCompleted
	t.UpdatedAt = time.Now().UTC()
	return nil
}

func (t *Task) Fail(reason string) error {
	if t.Status != TaskInProgress {
		return fmt.Errorf("cannot fail task in status %s", t.Status)
	}
	t.Status = TaskFailed
	t.UpdatedAt = time.Now().UTC()
	return nil
}

func (t *Task) Interrupt(reason string) error {
	if t.Status != TaskInProgress {
		return fmt.Errorf("cannot interrupt task in status %s", t.Status)
	}
	t.Status = TaskInterrupted
	t.UpdatedAt = time.Now().UTC()
	return nil
}

func (t *Task) Archive() error {
	if t.Status != TaskCompleted && t.Status != TaskFailed && t.Status != TaskInterrupted {
		return fmt.Errorf("cannot archive task in status %s", t.Status)
	}
	t.Status = TaskArchived
	t.UpdatedAt = time.Now().UTC()
	return nil
}
