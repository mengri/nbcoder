package sqlite

import (
	"fmt"

	"github.com/mengri/nbcoder/domain/project"
	"github.com/mengri/nbcoder/infrastructure/database/models"
	"gorm.io/gorm"
)

type ProjectConfigRepo struct {
	dbProvider DBProvider
}

func NewProjectConfigRepo(dbProvider DBProvider) project.ProjectConfigRepo {
	return &ProjectConfigRepo{dbProvider: dbProvider}
}

func (r *ProjectConfigRepo) getDB(projectName string) (*gorm.DB, error) {
	if projectName == "" {
		return r.dbProvider.GetGlobalDB(), nil
	}
	return r.dbProvider.GetProjectDB(projectName)
}

func (r *ProjectConfigRepo) Save(c *project.ProjectConfig) error {
	db, err := r.getDB(c.ProjectName)
	if err != nil {
		return err
	}

	model := &models.ProjectConfig{
		ID:          c.ID,
		ProjectName: c.ProjectName,
		Key:         c.Key,
		Value:       c.Value,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}

	result := db.Save(model)
	if result.Error != nil {
		return fmt.Errorf("failed to save project config: %w", result.Error)
	}
	return nil
}

func (r *ProjectConfigRepo) FindByID(id string, projectName string) (*project.ProjectConfig, error) {
	db, err := r.getDB(projectName)
	if err != nil {
		return nil, err
	}

	var model models.ProjectConfig
	result := db.First(&model, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find project config by id: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *ProjectConfigRepo) FindByProjectName(projectName string) ([]*project.ProjectConfig, error) {
	db, err := r.getDB(projectName)
	if err != nil {
		return nil, err
	}

	var models []models.ProjectConfig
	result := db.Where("project_name = ?", projectName).Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find configs by project name: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *ProjectConfigRepo) FindByKey(projectName, key string) (*project.ProjectConfig, error) {
	db, err := r.getDB(projectName)
	if err != nil {
		return nil, err
	}

	var model models.ProjectConfig
	result := db.Where("project_name = ? AND key = ?", projectName, key).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find config by project name and key: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *ProjectConfigRepo) FindAll() ([]*project.ProjectConfig, error) {
	db, err := r.getDB("")
	if err != nil {
		return nil, err
	}

	var models []models.ProjectConfig
	result := db.Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find all configs: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *ProjectConfigRepo) Update(c *project.ProjectConfig) error {
	db, err := r.getDB(c.ProjectName)
	if err != nil {
		return err
	}

	model := &models.ProjectConfig{
		ID:          c.ID,
		ProjectName: c.ProjectName,
		Key:         c.Key,
		Value:       c.Value,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}

	result := db.Model(&models.ProjectConfig{}).Where("id = ?", c.ID).Updates(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update project config: %w", result.Error)
	}
	return nil
}

func (r *ProjectConfigRepo) Delete(id string, projectName string) error {
	db, err := r.getDB(projectName)
	if err != nil {
		return err
	}

	result := db.Delete(&models.ProjectConfig{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete project config: %w", result.Error)
	}
	return nil
}

func (r *ProjectConfigRepo) modelToDomain(m *models.ProjectConfig) *project.ProjectConfig {
	return &project.ProjectConfig{
		ID:          m.ID,
		ProjectName: m.ProjectName,
		Key:         m.Key,
		Value:       m.Value,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func (r *ProjectConfigRepo) modelsToDomain(models []models.ProjectConfig) []*project.ProjectConfig {
	result := make([]*project.ProjectConfig, len(models))
	for i, m := range models {
		result[i] = r.modelToDomain(&m)
	}
	return result
}
