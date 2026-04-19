package sqlite

import (
	"fmt"

	"github.com/mengri/nbcoder/domain/project"
	"github.com/mengri/nbcoder/infrastructure/database/models"
	"gorm.io/gorm"
)

type ProjectConfigRepo struct {
	db *gorm.DB
}

func NewProjectConfigRepo(db *gorm.DB) project.ProjectConfigRepo {
	return &ProjectConfigRepo{db: db}
}

func (r *ProjectConfigRepo) Save(c *project.ProjectConfig) error {
	model := &models.ProjectConfig{
		ID:        c.ID,
		ProjectID: c.ProjectID,
		Key:       c.Key,
		Value:     c.Value,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}

	result := r.db.Save(model)
	if result.Error != nil {
		return fmt.Errorf("failed to save project config: %w", result.Error)
	}
	return nil
}

func (r *ProjectConfigRepo) FindByID(id string) (*project.ProjectConfig, error) {
	var model models.ProjectConfig
	result := r.db.First(&model, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find project config by id: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *ProjectConfigRepo) FindByProjectID(projectID string) ([]*project.ProjectConfig, error) {
	var models []models.ProjectConfig
	result := r.db.Where("project_id = ?", projectID).Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find configs by project id: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *ProjectConfigRepo) FindByKey(projectID, key string) (*project.ProjectConfig, error) {
	var model models.ProjectConfig
	result := r.db.Where("project_id = ? AND key = ?", projectID, key).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find config by project id and key: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *ProjectConfigRepo) FindAll() ([]*project.ProjectConfig, error) {
	var models []models.ProjectConfig
	result := r.db.Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find all configs: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *ProjectConfigRepo) Update(c *project.ProjectConfig) error {
	model := &models.ProjectConfig{
		ID:        c.ID,
		ProjectID: c.ProjectID,
		Key:       c.Key,
		Value:     c.Value,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}

	result := r.db.Model(&models.ProjectConfig{}).Where("id = ?", c.ID).Updates(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update project config: %w", result.Error)
	}
	return nil
}

func (r *ProjectConfigRepo) Delete(id string) error {
	result := r.db.Delete(&models.ProjectConfig{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete project config: %w", result.Error)
	}
	return nil
}

func (r *ProjectConfigRepo) modelToDomain(m *models.ProjectConfig) *project.ProjectConfig {
	return &project.ProjectConfig{
		ID:        m.ID,
		ProjectID: m.ProjectID,
		Key:       m.Key,
		Value:     m.Value,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (r *ProjectConfigRepo) modelsToDomain(models []models.ProjectConfig) []*project.ProjectConfig {
	result := make([]*project.ProjectConfig, len(models))
	for i, m := range models {
		result[i] = r.modelToDomain(&m)
	}
	return result
}
