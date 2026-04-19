package sqlite

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mengri/nbcoder/domain/agent"
	"github.com/mengri/nbcoder/infrastructure/database/models"
	"gorm.io/gorm"
)

type AgentExecutionRepo struct {
	db *gorm.DB
}

func NewAgentExecutionRepo(db *gorm.DB) agent.AgentExecutionRepo {
	return &AgentExecutionRepo{db: db}
}

func (r *AgentExecutionRepo) Save(execution *agent.AgentExecution) error {
	inputJSON, _ := json.Marshal(execution.Input)
	outputJSON, _ := json.Marshal(execution.Output)

	model := &models.AgentExecution{
		ID:          execution.ID,
		TaskID:      execution.TaskID,
		AgentID:     execution.AgentID,
		AgentType:   "",
		SkillName:   "",
		Status:      execution.Status,
		Input:       string(inputJSON),
		Output:      string(outputJSON),
		Error:       execution.Error,
		StartedAt:   &execution.StartTime,
		CompletedAt: execution.EndTime,
		CreatedAt:   execution.Timestamp,
		UpdatedAt:   execution.Timestamp,
		ExecCount:   0,
	}

	result := r.db.Save(model)
	if result.Error != nil {
		return fmt.Errorf("failed to save agent execution: %w", result.Error)
	}
	return nil
}

func (r *AgentExecutionRepo) FindByID(id string) (*agent.AgentExecution, error) {
	var model models.AgentExecution
	result := r.db.First(&model, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find agent execution by id: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *AgentExecutionRepo) FindByTaskID(taskID string) ([]*agent.AgentExecution, error) {
	var models []models.AgentExecution
	result := r.db.Where("task_id = ?", taskID).Order("created_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find agent executions by task id: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *AgentExecutionRepo) FindByAgentID(agentID string) ([]*agent.AgentExecution, error) {
	var models []models.AgentExecution
	result := r.db.Where("agent_id = ?", agentID).Order("created_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find agent executions by agent id: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *AgentExecutionRepo) FindAll() ([]*agent.AgentExecution, error) {
	var models []models.AgentExecution
	result := r.db.Order("created_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find all agent executions: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *AgentExecutionRepo) Update(execution *agent.AgentExecution) error {
	inputJSON, _ := json.Marshal(execution.Input)
	outputJSON, _ := json.Marshal(execution.Output)

	model := &models.AgentExecution{
		ID:          execution.ID,
		TaskID:      execution.TaskID,
		AgentID:     execution.AgentID,
		AgentType:   "",
		SkillName:   "",
		Status:      execution.Status,
		Input:       string(inputJSON),
		Output:      string(outputJSON),
		Error:       execution.Error,
		StartedAt:   &execution.StartTime,
		CompletedAt: execution.EndTime,
		CreatedAt:   execution.Timestamp,
		UpdatedAt:   execution.Timestamp,
		ExecCount:   0,
	}

	result := r.db.Model(&models.AgentExecution{}).Where("id = ?", execution.ID).Updates(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update agent execution: %w", result.Error)
	}
	return nil
}

func (r *AgentExecutionRepo) Delete(id string) error {
	result := r.db.Delete(&models.AgentExecution{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete agent execution: %w", result.Error)
	}
	return nil
}

func (r *AgentExecutionRepo) FindByStatus(status string) ([]*agent.AgentExecution, error) {
	var models []models.AgentExecution
	result := r.db.Where("status = ?", status).Order("created_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find agent executions by status: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *AgentExecutionRepo) FindByTaskIDAndAgentType(taskID string, agentType string) ([]*agent.AgentExecution, error) {
	var models []models.AgentExecution
	result := r.db.Where("task_id = ? AND agent_type = ?", taskID, agentType).Order("created_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find agent executions by task id and agent type: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *AgentExecutionRepo) modelToDomain(m *models.AgentExecution) *agent.AgentExecution {
	var input, output map[string]interface{}
	json.Unmarshal([]byte(m.Input), &input)
	json.Unmarshal([]byte(m.Output), &output)

	var startedTime time.Time
	if m.StartedAt != nil {
		startedTime = *m.StartedAt
	}

	return &agent.AgentExecution{
		ID:         m.ID,
		AgentID:    m.AgentID,
		TaskID:     m.TaskID,
		StartTime:  startedTime,
		EndTime:    m.CompletedAt,
		Status:     m.Status,
		Input:      input,
		Output:     output,
		Error:      m.Error,
		ModelUsed:  "",
		TokensUsed: 0,
		Timestamp:  m.CreatedAt,
	}
}

func (r *AgentExecutionRepo) modelsToDomain(models []models.AgentExecution) []*agent.AgentExecution {
	result := make([]*agent.AgentExecution, len(models))
	for i, m := range models {
		result[i] = r.modelToDomain(&m)
	}
	return result
}
