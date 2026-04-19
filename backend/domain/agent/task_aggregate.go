package agent

import (
	"fmt"
	"time"
)

type TaskAggregate struct {
	Task       *Task
	Logs       []*AgentExecution
	AssignedTo string
}

func NewTaskAggregate(task *Task) *TaskAggregate {
	return &TaskAggregate{
		Task: task,
		Logs: []*AgentExecution{},
	}
}

func (ta *TaskAggregate) Assign(agentID string) error {
	if err := ta.Task.Assign(agentID); err != nil {
		return err
	}
	ta.AssignedTo = agentID
	ta.Logs = append(ta.Logs, &AgentExecution{
		TaskID:    ta.Task.ID,
		AgentID:   agentID,
		Level:     "info",
		Message:   fmt.Sprintf("Task assigned to agent %s", agentID),
		Timestamp: time.Now().UTC(),
	})
	return nil
}

func (ta *TaskAggregate) Complete() error {
	if err := ta.Task.Complete(); err != nil {
		return err
	}
	ta.Logs = append(ta.Logs, &AgentExecution{
		TaskID:    ta.Task.ID,
		AgentID:   ta.AssignedTo,
		Level:     "info",
		Message:   "Task completed",
		Timestamp: time.Now().UTC(),
	})
	return nil
}

func (ta *TaskAggregate) Fail(reason string) error {
	if err := ta.Task.Fail(reason); err != nil {
		return err
	}
	ta.Logs = append(ta.Logs, &AgentExecution{
		TaskID:    ta.Task.ID,
		AgentID:   ta.AssignedTo,
		Level:     "error",
		Message:   fmt.Sprintf("Task failed: %s", reason),
		Timestamp: time.Now().UTC(),
	})
	return nil
}

func (ta *TaskAggregate) Interrupt(reason string) error {
	if err := ta.Task.Interrupt(reason); err != nil {
		return err
	}
	ta.Logs = append(ta.Logs, &AgentExecution{
		TaskID:    ta.Task.ID,
		AgentID:   ta.AssignedTo,
		Level:     "warn",
		Message:   fmt.Sprintf("Task interrupted: %s", reason),
		Timestamp: time.Now().UTC(),
	})
	return nil
}
