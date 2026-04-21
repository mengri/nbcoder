package sqlite

import (
	"fmt"

	"github.com/mengri/nbcoder/domain/project"
	"github.com/mengri/nbcoder/infrastructure/database/models"
	"gorm.io/gorm"
)

type DevStandardRepo struct {
	dbProvider DBProvider
}

func NewDevStandardRepo(dbProvider DBProvider) project.DevStandardRepo {
	return &DevStandardRepo{dbProvider: dbProvider}
}

func (r *DevStandardRepo) getDB(projectName string) (*gorm.DB, error) {
	if projectName == "" {
		return r.dbProvider.GetGlobalDB(), nil
	}
	return r.dbProvider.GetProjectDB(projectName)
}

func (r *DevStandardRepo) Save(s *project.DevStandard) error {
	db, err := r.getDB(s.ProjectName)
	if err != nil {
		return err
	}

	model := &models.DevStandard{
		ID:          s.ID,
		ProjectName: s.ProjectName,
		Name:        s.Name,
		Description: s.Description,
		StandardType: string(s.RuleType),
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}

	result := db.Save(model)
	if result.Error != nil {
		return fmt.Errorf("failed to save dev standard: %w", result.Error)
	}
	return nil
}

func (r *DevStandardRepo) FindByProjectName(projectName string) ([]*project.DevStandard, error) {
	db, err := r.getDB(projectName)
	if err != nil {
		return nil, err
	}

	var models []models.DevStandard
	result := db.Where("project_name = ?", projectName).Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find dev standards by project name: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *DevStandardRepo) FindByID(id string, projectName string) (*project.DevStandard, error) {
	db, err := r.getDB(projectName)
	if err != nil {
		return nil, err
	}

	var model models.DevStandard
	result := db.First(&model, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find dev standard by id: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *DevStandardRepo) Update(s *project.DevStandard) error {
	db, err := r.getDB(s.ProjectName)
	if err != nil {
		return err
	}

	model := &models.DevStandard{
		Name:        s.Name,
		Description: s.Description,
		StandardType: string(s.RuleType),
	}

	result := db.Model(&models.DevStandard{}).Where("id = ?", s.ID).Updates(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update dev standard: %w", result.Error)
	}
	return nil
}

func (r *DevStandardRepo) Delete(id string, projectName string) error {
	db, err := r.getDB(projectName)
	if err != nil {
		return err
	}

	result := db.Delete(&models.DevStandard{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete dev standard: %w", result.Error)
	}
	return nil
}

func (r *DevStandardRepo) FindByType(projectName, standardType string) ([]*project.DevStandard, error) {
	db, err := r.getDB(projectName)
	if err != nil {
		return nil, err
	}

	var models []models.DevStandard
	result := db.Where("project_name = ? AND standard_type = ?", projectName, standardType).Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find dev standards by type: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *DevStandardRepo) modelToDomain(m *models.DevStandard) *project.DevStandard {
	return &project.DevStandard{
		ID:          m.ID,
		ProjectName: m.ProjectName,
		Name:        m.Name,
		Description: m.Description,
		RuleType:    project.RuleType(m.StandardType),
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func (r *DevStandardRepo) modelsToDomain(models []models.DevStandard) []*project.DevStandard {
	result := make([]*project.DevStandard, len(models))
	for i, m := range models {
		result[i] = r.modelToDomain(&m)
	}
	return result
}
