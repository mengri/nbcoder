package agent

import (
	"testing"

	"github.com/mengri/nbcoder/domain/event"
	"github.com/mengri/nbcoder/infrastructure/eventbus"
)

func TestNewTaskAggregate(t *testing.T) {
	task := NewAgentTask("task-1", "Test Task", "Description", "development", AgentTypeTechStack, "proj-1")
	aggregate := NewTaskAggregate(task)

	if aggregate.Task.ID != "task-1" {
		t.Errorf("expected task ID task-1, got %s", aggregate.Task.ID)
	}
	if len(aggregate.Executions) != 0 {
		t.Errorf("expected 0 executions, got %d", len(aggregate.Executions))
	}
	if len(aggregate.Skills) != 0 {
		t.Errorf("expected 0 skills, got %d", len(aggregate.Skills))
	}
}

func TestTaskAggregate_Assign(t *testing.T) {
	task := NewAgentTask("task-1", "Test Task", "Description", "development", AgentTypeTechStack, "proj-1")
	aggregate := NewTaskAggregate(task)
	eventBus := eventbus.NewInMemoryEventBus()
	handler := NewTestEventHandler()

	if err := eventBus.Subscribe(string(event.TaskAssignedEvent), handler.Handle); err != nil {
		t.Fatalf("failed to subscribe to event: %v", err)
	}

	if err := aggregate.Assign("agent-1", eventBus); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if aggregate.Task.Status != TaskInProgress {
		t.Errorf("expected status IN_PROGRESS, got %s", aggregate.Task.Status)
	}

	if handler.GetEventCount() != 1 {
		t.Errorf("expected 1 event, got %d", handler.GetEventCount())
	}

	evt := handler.GetEvents()[0].(*event.AgentTaskEvent)
	if evt.Type != event.TaskAssignedEvent {
		t.Errorf("expected event type TaskAssigned, got %s", evt.Type)
	}

	if evt.Payload["agent_id"] != "agent-1" {
		t.Errorf("expected agent_id agent-1, got %v", evt.Payload["agent_id"])
	}
}

func TestTaskAggregate_Assign_InvalidStatus(t *testing.T) {
	task := NewAgentTask("task-1", "Test Task", "Description", "development", AgentTypeTechStack, "proj-1")
	task.Status = TaskCompleted
	aggregate := NewTaskAggregate(task)
	eventBus := eventbus.NewInMemoryEventBus()

	if err := aggregate.Assign("agent-1", eventBus); err == nil {
		t.Error("expected error assigning completed task")
	}
}

func TestTaskAggregate_Complete(t *testing.T) {
	task := NewAgentTask("task-1", "Test Task", "Description", "development", AgentTypeTechStack, "proj-1")
	aggregate := NewTaskAggregate(task)
	eventBus := eventbus.NewInMemoryEventBus()
	handler := NewTestEventHandler()

	if err := eventBus.Subscribe(string(event.TaskCompletedEvent), handler.Handle); err != nil {
		t.Fatalf("failed to subscribe to event: %v", err)
	}

	aggregate.Assign("agent-1", eventBus)

	if err := aggregate.Complete(eventBus); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if aggregate.Task.Status != TaskCompleted {
		t.Errorf("expected status COMPLETED, got %s", aggregate.Task.Status)
	}

	if handler.GetEventCount() != 1 {
		t.Errorf("expected 1 event, got %d", handler.GetEventCount())
	}

	evt := handler.GetEvents()[0].(*event.AgentTaskEvent)
	if evt.Type != event.TaskCompletedEvent {
		t.Errorf("expected event type TaskCompleted, got %s", evt.Type)
	}
}

func TestTaskAggregate_Fail(t *testing.T) {
	task := NewAgentTask("task-1", "Test Task", "Description", "development", AgentTypeTechStack, "proj-1")
	aggregate := NewTaskAggregate(task)
	eventBus := eventbus.NewInMemoryEventBus()
	handler := NewTestEventHandler()

	if err := eventBus.Subscribe(string(event.TaskFailedEvent), handler.Handle); err != nil {
		t.Fatalf("failed to subscribe to event: %v", err)
	}

	aggregate.Assign("agent-1", eventBus)

	if err := aggregate.Fail("timeout", eventBus); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if aggregate.Task.Status != TaskFailed {
		t.Errorf("expected status FAILED, got %s", aggregate.Task.Status)
	}

	if handler.GetEventCount() != 1 {
		t.Errorf("expected 1 event, got %d", handler.GetEventCount())
	}

	evt := handler.GetEvents()[0].(*event.AgentTaskEvent)
	if evt.Type != event.TaskFailedEvent {
		t.Errorf("expected event type TaskFailed, got %s", evt.Type)
	}

	if evt.Payload["reason"] != "timeout" {
		t.Errorf("expected reason timeout, got %v", evt.Payload["reason"])
	}
}

func TestTaskAggregate_Interrupt(t *testing.T) {
	task := NewAgentTask("task-1", "Test Task", "Description", "development", AgentTypeTechStack, "proj-1")
	aggregate := NewTaskAggregate(task)
	eventBus := eventbus.NewInMemoryEventBus()
	handler := NewTestEventHandler()

	if err := eventBus.Subscribe(string(event.TaskInterruptedEvent), handler.Handle); err != nil {
		t.Fatalf("failed to subscribe to event: %v", err)
	}

	aggregate.Assign("agent-1", eventBus)

	if err := aggregate.Interrupt("user cancelled", eventBus); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if aggregate.Task.Status != TaskInterrupted {
		t.Errorf("expected status INTERRUPTED, got %s", aggregate.Task.Status)
	}

	if handler.GetEventCount() != 1 {
		t.Errorf("expected 1 event, got %d", handler.GetEventCount())
	}

	evt := handler.GetEvents()[0].(*event.AgentTaskEvent)
	if evt.Type != event.TaskInterruptedEvent {
		t.Errorf("expected event type TaskInterrupted, got %s", evt.Type)
	}

	if evt.Payload["reason"] != "user cancelled" {
		t.Errorf("expected reason user cancelled, got %v", evt.Payload["reason"])
	}
}

func TestTaskAggregate_Archive(t *testing.T) {
	task := NewAgentTask("task-1", "Test Task", "Description", "development", AgentTypeTechStack, "proj-1")
	aggregate := NewTaskAggregate(task)
	eventBus := eventbus.NewInMemoryEventBus()
	handler := NewTestEventHandler()

	if err := eventBus.Subscribe(string(event.TaskArchivedEvent), handler.Handle); err != nil {
		t.Fatalf("failed to subscribe to event: %v", err)
	}

	aggregate.Assign("agent-1", eventBus)
	aggregate.Complete(eventBus)

	if err := aggregate.Archive(eventBus); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if aggregate.Task.Status != TaskArchived {
		t.Errorf("expected status ARCHIVED, got %s", aggregate.Task.Status)
	}

	if handler.GetEventCount() != 2 {
		t.Errorf("expected 2 events, got %d", handler.GetEventCount())
	}

	if handler.GetEvents()[1].(*event.AgentTaskEvent).Type != event.TaskArchivedEvent {
		t.Errorf("expected last event to be TaskArchived, got %s", handler.GetEvents()[1].(*event.AgentTaskEvent).Type)
	}
}

func TestTaskAggregate_AddExecution(t *testing.T) {
	task := NewAgentTask("task-1", "Test Task", "Description", "development", AgentTypeTechStack, "proj-1")
	aggregate := NewTaskAggregate(task)

	execution := NewAgentExecution("exec-1", "agent-1", "task-1")
	aggregate.AddExecution(execution)

	if len(aggregate.Executions) != 1 {
		t.Errorf("expected 1 execution, got %d", len(aggregate.Executions))
	}
}

func TestTaskAggregate_AddSkill(t *testing.T) {
	task := NewAgentTask("task-1", "Test Task", "Description", "development", AgentTypeTechStack, "proj-1")
	aggregate := NewTaskAggregate(task)

	skill := NewSkill("skill-1", "Code Generation", "Generate code", AgentTypeTechStack)
	aggregate.AddSkill(skill)

	if len(aggregate.Skills) != 1 {
		t.Errorf("expected 1 skill, got %d", len(aggregate.Skills))
	}
}

func TestTaskAggregate_GetExecutionByID(t *testing.T) {
	task := NewAgentTask("task-1", "Test Task", "Description", "development", AgentTypeTechStack, "proj-1")
	aggregate := NewTaskAggregate(task)

	execution := NewAgentExecution("exec-1", "agent-1", "task-1")
	aggregate.AddExecution(execution)

	found, err := aggregate.GetExecutionByID("exec-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if found.ID != "exec-1" {
		t.Errorf("expected execution ID exec-1, got %s", found.ID)
	}

	_, err = aggregate.GetExecutionByID("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent execution")
	}
}

func TestTaskAggregate_GetSkillsByAgentType(t *testing.T) {
	task := NewAgentTask("task-1", "Test Task", "Description", "development", AgentTypeTechStack, "proj-1")
	aggregate := NewTaskAggregate(task)

	skill1 := NewSkill("skill-1", "Code Gen", "Generate code", AgentTypeTechStack)
	skill2 := NewSkill("skill-2", "Design", "Create design", AgentTypeArchitecture)
	skill3 := NewSkill("skill-3", "Another Code Gen", "More code", AgentTypeTechStack)

	aggregate.AddSkill(skill1)
	aggregate.AddSkill(skill2)
	aggregate.AddSkill(skill3)

	techStackSkills := aggregate.GetSkillsByAgentType(AgentTypeTechStack)
	if len(techStackSkills) != 2 {
		t.Errorf("expected 2 tech stack skills, got %d", len(techStackSkills))
	}

	architectureSkills := aggregate.GetSkillsByAgentType(AgentTypeArchitecture)
	if len(architectureSkills) != 1 {
		t.Errorf("expected 1 architecture skill, got %d", len(architectureSkills))
	}
}

type TestEventHandler struct {
	events []event.DomainEvent
}

func NewTestEventHandler() *TestEventHandler {
	return &TestEventHandler{
		events: []event.DomainEvent{},
	}
}

func (h *TestEventHandler) Handle(evt event.DomainEvent) {
	h.events = append(h.events, evt)
}

func (h *TestEventHandler) GetEventCount() int {
	return len(h.events)
}

func (h *TestEventHandler) GetEvents() []event.DomainEvent {
	return h.events
}