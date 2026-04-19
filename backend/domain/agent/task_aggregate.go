package agent

import (
	"fmt"

	"github.com/mengri/nbcoder/domain/event"
	"github.com/mengri/nbcoder/pkg/uid"
)

type TaskAggregate struct {
	Task       *Task
	Executions []*AgentExecution
	Skills     []*Skill
}

func NewTaskAggregate(task *Task) *TaskAggregate {
	return &TaskAggregate{
		Task:       task,
		Executions: []*AgentExecution{},
		Skills:     []*Skill{},
	}
}

func (ta *TaskAggregate) Assign(agentID string, eventBus event.EventBus) error {
	if err := ta.Task.Assign(agentID); err != nil {
		return err
	}

	evt := event.NewAgentTaskEvent(uid.NewID(), ta.Task.ID, event.TaskAssignedEvent)
	evt.Payload["agent_id"] = agentID
	evt.Payload["agent_type"] = string(ta.Task.AgentType)
	evt.Payload["project_id"] = ta.Task.ProjectID

	return eventBus.Publish(evt)
}

func (ta *TaskAggregate) Complete(eventBus event.EventBus) error {
	if err := ta.Task.Complete(); err != nil {
		return err
	}

	evt := event.NewAgentTaskEvent(uid.NewID(), ta.Task.ID, event.TaskCompletedEvent)
	evt.Payload["duration_ms"] = ta.Task.GetDuration().Milliseconds()
	evt.Payload["executions_count"] = len(ta.Executions)

	return eventBus.Publish(evt)
}

func (ta *TaskAggregate) Fail(reason string, eventBus event.EventBus) error {
	if err := ta.Task.Fail(reason); err != nil {
		return err
	}

	evt := event.NewAgentTaskEvent(uid.NewID(), ta.Task.ID, event.TaskFailedEvent)
	evt.Payload["reason"] = reason
	evt.Payload["duration_ms"] = ta.Task.GetDuration().Milliseconds()

	return eventBus.Publish(evt)
}

func (ta *TaskAggregate) Interrupt(reason string, eventBus event.EventBus) error {
	if err := ta.Task.Interrupt(reason); err != nil {
		return err
	}

	evt := event.NewAgentTaskEvent(uid.NewID(), ta.Task.ID, event.TaskInterruptedEvent)
	evt.Payload["reason"] = reason
	evt.Payload["duration_ms"] = ta.Task.GetDuration().Milliseconds()

	return eventBus.Publish(evt)
}

func (ta *TaskAggregate) Archive(eventBus event.EventBus) error {
	if err := ta.Task.Archive(); err != nil {
		return err
	}

	evt := event.NewAgentTaskEvent(uid.NewID(), ta.Task.ID, event.TaskArchivedEvent)
	evt.Payload["project_id"] = ta.Task.ProjectID

	return eventBus.Publish(evt)
}

func (ta *TaskAggregate) AddExecution(execution *AgentExecution) {
	ta.Executions = append(ta.Executions, execution)
}

func (ta *TaskAggregate) AddSkill(skill *Skill) {
	ta.Skills = append(ta.Skills, skill)
}

func (ta *TaskAggregate) GetExecutionByID(executionID string) (*AgentExecution, error) {
	for _, execution := range ta.Executions {
		if execution.ID == executionID {
			return execution, nil
		}
	}
	return nil, fmt.Errorf("execution %s not found", executionID)
}

func (ta *TaskAggregate) GetSkillsByAgentType(agentType AgentType) []*Skill {
	var skills []*Skill
	for _, skill := range ta.Skills {
		if skill.AgentType == agentType {
			skills = append(skills, skill)
		}
	}
	return skills
}
