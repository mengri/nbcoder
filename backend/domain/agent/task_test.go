package agent

import (
	"testing"
	"time"
)

func TestAgentTaskStatus_IsValid(t *testing.T) {
	validStatuses := []AgentTaskStatus{TaskPending, TaskInProgress, TaskCompleted, TaskFailed, TaskInterrupted, TaskArchived}

	for _, status := range validStatuses {
		if !status.IsValid() {
			t.Errorf("expected status %s to be valid", status)
		}
	}

	invalidStatus := AgentTaskStatus("INVALID")
	if invalidStatus.IsValid() {
		t.Error("expected INVALID status to be invalid")
	}
}

func TestAgentTaskStatus_CanTransitionTo(t *testing.T) {
	tests := []struct {
		currentStatus   AgentTaskStatus
		targetStatus   AgentTaskStatus
		expectedResult bool
	}{
		{TaskPending, TaskInProgress, true},
		{TaskPending, TaskArchived, true},
		{TaskPending, TaskCompleted, false},
		{TaskInProgress, TaskCompleted, true},
		{TaskInProgress, TaskFailed, true},
		{TaskInProgress, TaskInterrupted, true},
		{TaskInProgress, TaskPending, false},
		{TaskCompleted, TaskArchived, true},
		{TaskFailed, TaskArchived, true},
		{TaskInterrupted, TaskArchived, true},
		{TaskArchived, TaskCompleted, false},
	}

	for _, test := range tests {
		result := test.currentStatus.CanTransitionTo(test.targetStatus)
		if result != test.expectedResult {
			t.Errorf("expected %s->%s transition to be %v, got %v", 
				test.currentStatus, test.targetStatus, test.expectedResult, result)
		}
	}
}

func TestAgentType_IsValid(t *testing.T) {
	validTypes := []AgentType{AgentTypeProduct, AgentTypeArchitecture, AgentTypeManagement, AgentTypeTechStack}

	for _, agentType := range validTypes {
		if !agentType.IsValid() {
			t.Errorf("expected agent type %s to be valid", agentType)
		}
	}

	invalidType := AgentType("INVALID")
	if invalidType.IsValid() {
		t.Error("expected INVALID agent type to be invalid")
	}
}

func TestNewAgentTask(t *testing.T) {
	task := NewAgentTask("task-1", "Test Task", "Description", "development", AgentTypeTechStack, "proj-1")

	if task.ID != "task-1" {
		t.Errorf("expected ID task-1, got %s", task.ID)
	}
	if task.Name != "Test Task" {
		t.Errorf("expected name Test Task, got %s", task.Name)
	}
	if task.Status != TaskPending {
		t.Errorf("expected status PENDING, got %s", task.Status)
	}
	if task.AgentType != AgentTypeTechStack {
		t.Errorf("expected agent type TECH_STACK, got %s", task.AgentType)
	}
	if task.ProjectID != "proj-1" {
		t.Errorf("expected project ID proj-1, got %s", task.ProjectID)
	}
	if task.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
	if task.Priority != 5 {
		t.Errorf("expected default priority 5, got %d", task.Priority)
	}
}

func TestAgentTask_Assign(t *testing.T) {
	task := NewAgentTask("task-1", "Test Task", "Description", "development", AgentTypeTechStack, "proj-1")

	err := task.Assign("agent-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if task.Status != TaskInProgress {
		t.Errorf("expected status IN_PROGRESS, got %s", task.Status)
	}
	if task.AssignedTo != "agent-1" {
		t.Errorf("expected assigned to agent-1, got %s", task.AssignedTo)
	}
	if task.StartedAt == nil {
		t.Error("expected StartedAt to be set")
	}
}

func TestAgentTask_Assign_InvalidStatus(t *testing.T) {
	task := NewAgentTask("task-1", "Test Task", "Description", "development", AgentTypeTechStack, "proj-1")
	task.Status = TaskCompleted

	err := task.Assign("agent-1")
	if err == nil {
		t.Error("expected error assigning completed task")
	}
}

func TestAgentTask_Complete(t *testing.T) {
	task := NewAgentTask("task-1", "Test Task", "Description", "development", AgentTypeTechStack, "proj-1")
	task.Assign("agent-1")

	err := task.Complete()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if task.Status != TaskCompleted {
		t.Errorf("expected status COMPLETED, got %s", task.Status)
	}
	if task.CompletedAt == nil {
		t.Error("expected CompletedAt to be set")
	}
}

func TestAgentTask_Complete_InvalidStatus(t *testing.T) {
	task := NewAgentTask("task-1", "Test Task", "Description", "development", AgentTypeTechStack, "proj-1")

	err := task.Complete()
	if err == nil {
		t.Error("expected error completing pending task")
	}
}

func TestAgentTask_Fail(t *testing.T) {
	task := NewAgentTask("task-1", "Test Task", "Description", "development", AgentTypeTechStack, "proj-1")
	task.Assign("agent-1")

	err := task.Fail("timeout error")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if task.Status != TaskFailed {
		t.Errorf("expected status FAILED, got %s", task.Status)
	}
	if task.CompletedAt == nil {
		t.Error("expected CompletedAt to be set")
	}
}

func TestAgentTask_Interrupt(t *testing.T) {
	task := NewAgentTask("task-1", "Test Task", "Description", "development", AgentTypeTechStack, "proj-1")
	task.Assign("agent-1")

	err := task.Interrupt("user cancelled")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if task.Status != TaskInterrupted {
		t.Errorf("expected status INTERRUPTED, got %s", task.Status)
	}
	if task.CompletedAt == nil {
		t.Error("expected CompletedAt to be set")
	}
}

func TestAgentTask_Archive(t *testing.T) {
	task := NewAgentTask("task-1", "Test Task", "Description", "development", AgentTypeTechStack, "proj-1")
	task.Assign("agent-1")
	task.Complete()

	err := task.Archive()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if task.Status != TaskArchived {
		t.Errorf("expected status ARCHIVED, got %s", task.Status)
	}
}

func TestAgentTask_Archive_InvalidStatus(t *testing.T) {
	task := NewAgentTask("task-1", "Test Task", "Description", "development", AgentTypeTechStack, "proj-1")

	err := task.Archive()
	if err == nil {
		t.Error("expected error archiving pending task")
	}
}

func TestAgentTask_UpdateContext(t *testing.T) {
	task := NewAgentTask("task-1", "Test Task", "Description", "development", AgentTypeTechStack, "proj-1")

	task.UpdateContext("key1", "value1")
	task.UpdateContext("key2", 123)

	if task.Context["key1"] != "value1" {
		t.Errorf("expected key1 value1, got %v", task.Context["key1"])
	}
	if task.Context["key2"] != 123 {
		t.Errorf("expected key2 123, got %v", task.Context["key2"])
	}
}

func TestAgentTask_SetPriority(t *testing.T) {
	task := NewAgentTask("task-1", "Test Task", "Description", "development", AgentTypeTechStack, "proj-1")

	task.SetPriority(10)

	if task.Priority != 10 {
		t.Errorf("expected priority 10, got %d", task.Priority)
	}
}

func TestAgentTask_SetPipelineID(t *testing.T) {
	task := NewAgentTask("task-1", "Test Task", "Description", "development", AgentTypeTechStack, "proj-1")

	task.SetPipelineID("pipeline-1")

	if task.PipelineID != "pipeline-1" {
		t.Errorf("expected pipeline ID pipeline-1, got %s", task.PipelineID)
	}
}

func TestAgentTask_GetDuration(t *testing.T) {
	task := NewAgentTask("task-1", "Test Task", "Description", "development", AgentTypeTechStack, "proj-1")

	duration := task.GetDuration()
	if duration != 0 {
		t.Errorf("expected duration 0 for unstarted task, got %v", duration)
	}

	task.Assign("agent-1")
	time.Sleep(10 * time.Millisecond)
	duration = task.GetDuration()

	if duration < 10*time.Millisecond {
		t.Errorf("expected duration > 10ms, got %v", duration)
	}

	task.Complete()
	duration = task.GetDuration()

	if duration < 10*time.Millisecond {
		t.Errorf("expected duration > 10ms for completed task, got %v", duration)
	}
}

func TestNewSkill(t *testing.T) {
	skill := NewSkill("skill-1", "Code Generation", "Generate code snippets", AgentTypeTechStack)

	if skill.ID != "skill-1" {
		t.Errorf("expected ID skill-1, got %s", skill.ID)
	}
	if skill.Name != "Code Generation" {
		t.Errorf("expected name Code Generation, got %s", skill.Name)
	}
	if skill.AgentType != AgentTypeTechStack {
		t.Errorf("expected agent type TECH_STACK, got %s", skill.AgentType)
	}
	if skill.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestSkill_SetConfig(t *testing.T) {
	skill := NewSkill("skill-1", "Code Generation", "Generate code snippets", AgentTypeTechStack)

	skill.SetConfig("model", "gpt-4")
	skill.SetConfig("max_tokens", 2000)

	if skill.Config["model"] != "gpt-4" {
		t.Errorf("expected model gpt-4, got %v", skill.Config["model"])
	}
	if skill.Config["max_tokens"] != 2000 {
		t.Errorf("expected max_tokens 2000, got %v", skill.Config["max_tokens"])
	}
}

func TestSkill_GetConfig(t *testing.T) {
	skill := NewSkill("skill-1", "Code Generation", "Generate code snippets", AgentTypeTechStack)
	skill.SetConfig("model", "gpt-4")

	value, exists := skill.GetConfig("model")
	if !exists || value != "gpt-4" {
		t.Errorf("expected model config to exist with value gpt-4")
	}

	value, exists = skill.GetConfig("nonexistent")
	if exists {
		t.Error("expected nonexistent config to not exist")
	}
}

func TestNewAgentExecution(t *testing.T) {
	execution := NewAgentExecution("exec-1", "agent-1", "task-1")

	if execution.ID != "exec-1" {
		t.Errorf("expected ID exec-1, got %s", execution.ID)
	}
	if execution.AgentID != "agent-1" {
		t.Errorf("expected agent ID agent-1, got %s", execution.AgentID)
	}
	if execution.TaskID != "task-1" {
		t.Errorf("expected task ID task-1, got %s", execution.TaskID)
	}
	if execution.Status != "IN_PROGRESS" {
		t.Errorf("expected status IN_PROGRESS, got %s", execution.Status)
	}
	if execution.StartTime.IsZero() {
		t.Error("expected StartTime to be set")
	}
}

func TestAgentExecution_Complete(t *testing.T) {
	execution := NewAgentExecution("exec-1", "agent-1", "task-1")
	time.Sleep(10 * time.Millisecond)

	output := map[string]interface{}{"result": "success"}
	execution.Complete(output, "gpt-4", 150)

	if execution.Status != "COMPLETED" {
		t.Errorf("expected status COMPLETED, got %s", execution.Status)
	}
	if execution.EndTime == nil {
		t.Error("expected EndTime to be set")
	}
	if execution.Output["result"] != "success" {
		t.Errorf("expected output result success, got %v", execution.Output["result"])
	}
	if execution.ModelUsed != "gpt-4" {
		t.Errorf("expected model gpt-4, got %s", execution.ModelUsed)
	}
	if execution.TokensUsed != 150 {
		t.Errorf("expected tokens 150, got %d", execution.TokensUsed)
	}
}

func TestAgentExecution_Fail(t *testing.T) {
	execution := NewAgentExecution("exec-1", "agent-1", "task-1")

	execution.Fail("API timeout")

	if execution.Status != "FAILED" {
		t.Errorf("expected status FAILED, got %s", execution.Status)
	}
	if execution.EndTime == nil {
		t.Error("expected EndTime to be set")
	}
	if execution.Error != "API timeout" {
		t.Errorf("expected error API timeout, got %s", execution.Error)
	}
}

func TestAgentExecution_GetDuration(t *testing.T) {
	execution := NewAgentExecution("exec-1", "agent-1", "task-1")

	duration := execution.GetDuration()
	if duration > 0 {
		t.Errorf("expected duration 0 for in-progress execution, got %v", duration)
	}

	time.Sleep(10 * time.Millisecond)
	execution.Complete(map[string]interface{}{}, "gpt-4", 100)

	duration = execution.GetDuration()
	if duration < 10*time.Millisecond {
		t.Errorf("expected duration > 10ms, got %v", duration)
	}
}