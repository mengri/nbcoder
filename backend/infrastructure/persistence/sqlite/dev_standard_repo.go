package sqlite

import (
	"fmt"

	"github.com/mengri/nbcoder/domain/project"
	"github.com/mengri/nbcoder/infrastructure/database/models"
	"gorm.io/gorm"
)

type DevStandardRepo struct {
	db *gorm.DB
}

func NewDevStandardRepo(db *gorm.DB) project.DevStandardRepo {
	return &DevStandardRepo{db: db}
}

func (r *DevStandardRepo) Save(s *project.DevStandard) error {
	model := &models.DevStandard{
		ID:          s.ID,
		ProjectID:   s.ProjectID,
		Name:        s.Name,
		Description: s.Description,
		StandardType: string(s.RuleType),
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}

	result := r.db.Save(model)
	if result.Error != nil {
		return fmt.Errorf("failed to save dev standard: %w", result.Error)
	}
	return nil
}

func (r *DevStandardRepo) FindByProjectID(projectID string) ([]*project.DevStandard, error) {
	var models []models.DevStandard
	result := r.db.Where("project_id = ?", projectID).Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find dev standards by project id: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *DevStandardRepo) FindByID(id string) (*project.DevStandard, error) {
	var model models.DevStandard
	result := r.db.First(&model, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find dev standard by id: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *DevStandardRepo) Update(s *project.DevStandard) error {
	model := &models.DevStandard{
		Name:        s.Name,
		Description: s.Description,
		StandardType: string(s.RuleType),
	}

	result := r.db.Model(&models.DevStandard{}).Where("id = ?", s.ID).Updates(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update dev standard: %w", result.Error)
	}
	return nil
}

func (r *DevStandardRepo) Delete(id string) error {
	result := r.db.Delete(&models.DevStandard{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete dev standard: %w", result.Error)
	}
	return nil
}

func (r *DevStandardRepo) FindByType(projectID, standardType string) ([]*project.DevStandard, error) {
	var models []models.DevStandard
	result := r.db.Where("project_id = ? AND standard_type = ?", projectID, standardType).Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find dev standards by type: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *DevStandardRepo) modelToDomain(m *models.DevStandard) *project.DevStandard {
	return &project.DevStandard{
		ID:          m.ID,
		ProjectID:   m.ProjectID,
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
