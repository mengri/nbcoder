package sqlite

import (
	"encoding/json"
	"fmt"

	"github.com/mengri/nbcoder/domain/agent"
	"github.com/mengri/nbcoder/infrastructure/database/models"
	"gorm.io/gorm"
)

type TaskRepo struct {
	dbProvider DBProvider
}

func NewTaskRepo(dbProvider DBProvider) agent.TaskRepo {
	return &TaskRepo{dbProvider: dbProvider}
}

func (r *TaskRepo) getDB(projectName string) (*gorm.DB, error) {
	if projectName == "" {
		return r.dbProvider.GetGlobalDB(), nil
	}
	return r.dbProvider.GetProjectDB(projectName)
}

func (r *TaskRepo) Save(task *agent.Task) error {
	db, err := r.getDB(task.ProjectName)
	if err != nil {
		return err
	}

	_, err = json.Marshal(task.Context)
	if err != nil {
		return fmt.Errorf("failed to marshal task context: %w", err)
	}

	model := &models.Task{
		ID:          task.ID,
		Name:        task.Name,
		Description: task.Description,
		AgentType:   string(task.AgentType),
		TaskType:    task.TaskType,
		Status:      string(task.Status),
		Priority:    task.Priority,
		AssignedTo:  task.AssignedTo,
		PipelineID:  task.PipelineID,
		ProjectName: task.ProjectName,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
		StartedAt:   task.StartedAt,
		CompletedAt: task.CompletedAt,
		Context:     models.JSONMap(task.Context),
	}

	result := db.Save(model)
	if result.Error != nil {
		return fmt.Errorf("failed to save task: %w", result.Error)
	}
	return nil
}

func (r *TaskRepo) FindByID(id string, projectName string) (*agent.Task, error) {
	db, err := r.getDB(projectName)
	if err != nil {
		return nil, err
	}

	var model models.Task
	result := db.First(&model, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find task by id: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *TaskRepo) FindByProjectName(projectName string) ([]*agent.Task, error) {
	db, err := r.getDB(projectName)
	if err != nil {
		return nil, err
	}

	var models []models.Task
	result := db.Where("project_name = ?", projectName).Order("created_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find tasks by project name: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *TaskRepo) FindByStatus(status agent.TaskStatus) ([]*agent.Task, error) {
	db, err := r.getDB("")
	if err != nil {
		return nil, err
	}

	var models []models.Task
	result := db.Where("status = ?", status).Order("created_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find tasks by status: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *TaskRepo) FindByAgentID(agentID string) ([]*agent.Task, error) {
	db, err := r.getDB("")
	if err != nil {
		return nil, err
	}

	var models []models.Task
	result := db.Where("assigned_to = ?", agentID).Order("created_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find tasks by agent id: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *TaskRepo) FindByPipelineID(pipelineID string) ([]*agent.Task, error) {
	db, err := r.getDB("")
	if err != nil {
		return nil, err
	}

	var models []models.Task
	result := db.Where("pipeline_id = ?", pipelineID).Order("created_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find tasks by pipeline id: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *TaskRepo) FindAll() ([]*agent.Task, error) {
	db, err := r.getDB("")
	if err != nil {
		return nil, err
	}

	var models []models.Task
	result := db.Order("created_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find all tasks: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *TaskRepo) Update(task *agent.Task) error {
	db, err := r.getDB(task.ProjectName)
	if err != nil {
		return err
	}

	_, err = json.Marshal(task.Context)
	if err != nil {
		return fmt.Errorf("failed to marshal task context: %w", err)
	}

	model := &models.Task{
		ID:          task.ID,
		Name:        task.Name,
		Description: task.Description,
		AgentType:   string(task.AgentType),
		TaskType:    task.TaskType,
		Status:      string(task.Status),
		Priority:    task.Priority,
		AssignedTo:  task.AssignedTo,
		PipelineID:  task.PipelineID,
		ProjectName: task.ProjectName,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
		StartedAt:   task.StartedAt,
		CompletedAt: task.CompletedAt,
		Context:     models.JSONMap(task.Context),
	}

	result := db.Model(&models.Task{}).Where("id = ?", task.ID).Updates(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update task: %w", result.Error)
	}
	return nil
}

func (r *TaskRepo) Delete(id string, projectName string) error {
	db, err := r.getDB(projectName)
	if err != nil {
		return err
	}

	result := db.Delete(&models.Task{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete task: %w", result.Error)
	}
	return nil
}

func (r *TaskRepo) FindByPriority(minPriority int) ([]*agent.Task, error) {
	db, err := r.getDB("")
	if err != nil {
		return nil, err
	}

	var models []models.Task
	result := db.Where("priority >= ?", minPriority).Order("priority DESC, created_at ASC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find tasks by priority: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *TaskRepo) FindPending() ([]*agent.Task, error) {
	db, err := r.getDB("")
	if err != nil {
		return nil, err
	}

	var models []models.Task
	result := db.Where("status = ?", "PENDING").Order("priority DESC, created_at ASC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find pending tasks: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *TaskRepo) modelToDomain(m *models.Task) *agent.Task {
	return &agent.Task{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		AgentType:   agent.AgentType(m.AgentType),
		TaskType:    m.TaskType,
		Status:      agent.TaskStatus(m.Status),
		Priority:    m.Priority,
		AssignedTo:  m.AssignedTo,
		PipelineID:  m.PipelineID,
		ProjectName: m.ProjectName,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		StartedAt:   m.StartedAt,
		CompletedAt: m.CompletedAt,
		Context:     map[string]interface{}(m.Context),
	}
}

func (r *TaskRepo) modelsToDomain(models []models.Task) []*agent.Task {
	result := make([]*agent.Task, len(models))
	for i, m := range models {
		result[i] = r.modelToDomain(&m)
	}
	return result
}
