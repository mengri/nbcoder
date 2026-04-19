package agent

import (
	"github.com/mengri/nbcoder/domain/agent"
	"github.com/mengri/nbcoder/domain/event"
	"github.com/mengri/nbcoder/pkg/uid"
)

type AgentService struct {
	taskRepo       agent.TaskRepo
	executionRepo  agent.AgentExecutionRepo
	agentRegistry  *agent.AgentRegistry
	eventBus       event.EventBus
}

func NewAgentService(
	taskRepo agent.TaskRepo,
	executionRepo agent.AgentExecutionRepo,
	agentRegistry *agent.AgentRegistry,
	eventBus event.EventBus,
) *AgentService {
	return &AgentService{
		taskRepo:      taskRepo,
		executionRepo: executionRepo,
		agentRegistry: agentRegistry,
		eventBus:      eventBus,
	}
}

func (s *AgentService) CreateTask(id, name, desc string) (*agent.TaskAggregate, error) {
	task := agent.NewTask(id, name, desc)
	aggregate := agent.NewTaskAggregate(task)
	if err := s.taskRepo.Save(task); err != nil {
		return nil, err
	}
	return aggregate, nil
}

func (s *AgentService) AssignTask(taskID, agentID string) error {
	task, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return err
	}
	aggregate := agent.NewTaskAggregate(task)
	if err := aggregate.Assign(agentID); err != nil {
		return err
	}
	if err := s.taskRepo.Update(task); err != nil {
		return err
	}
	evt := event.NewAgentTaskEvent(
		uid.NewID(), taskID, agentID, event.TaskAssignedEvent,
	)
	return s.eventBus.Publish(evt)
}

func (s *AgentService) CompleteTask(taskID string) error {
	task, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return err
	}
	aggregate := agent.NewTaskAggregate(task)
	if err := aggregate.Complete(); err != nil {
		return err
	}
	if err := s.taskRepo.Update(task); err != nil {
		return err
	}
	evt := event.NewAgentTaskEvent(
		uid.NewID(), taskID, task.AssignedTo, event.TaskCompletedEvent,
	)
	return s.eventBus.Publish(evt)
}

func (s *AgentService) FailTask(taskID, reason string) error {
	task, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return err
	}
	aggregate := agent.NewTaskAggregate(task)
	if err := aggregate.Fail(reason); err != nil {
		return err
	}
	if err := s.taskRepo.Update(task); err != nil {
		return err
	}
	evt := event.NewAgentTaskEvent(
		uid.NewID(), taskID, task.AssignedTo, event.TaskFailedEvent,
	)
	return s.eventBus.Publish(evt)
}

func (s *AgentService) GetTask(taskID string) (*agent.Task, error) {
	return s.taskRepo.FindByID(taskID)
}
