package sqlite

import (
	"fmt"

	"github.com/mengri/nbcoder/domain/project"
	"github.com/mengri/nbcoder/infrastructure/database/models"
	"gorm.io/gorm"
)

type StandardsRepo struct {
	db *gorm.DB
}

func NewStandardsRepo(db *gorm.DB) project.StandardsRepo {
	return &StandardsRepo{db: db}
}

func (r *StandardsRepo) Save(s *project.Standards) error {
	model := &models.Standards{
		ID:                s.ID,
		ProjectID:         s.ProjectID,
		BranchStrategy:    s.BranchStrategy,
		TechStack:         s.TechStack,
		CodingConventions: s.CodingConventions,
		CreatedAt:         s.CreatedAt,
		UpdatedAt:         s.UpdatedAt,
	}

	result := r.db.Save(model)
	if result.Error != nil {
		return fmt.Errorf("failed to save standards: %w", result.Error)
	}
	return nil
}

func (r *StandardsRepo) FindByProjectID(projectID string) (*project.Standards, error) {
	var model models.Standards
	result := r.db.Where("project_id = ?", projectID).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find standards by project id: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *StandardsRepo) Update(s *project.Standards) error {
	model := &models.Standards{
		ID:                s.ID,
		ProjectID:         s.ProjectID,
		BranchStrategy:    s.BranchStrategy,
		TechStack:         s.TechStack,
		CodingConventions: s.CodingConventions,
		CreatedAt:         s.CreatedAt,
		UpdatedAt:         s.UpdatedAt,
	}

	result := r.db.Model(&models.Standards{}).Where("id = ?", s.ID).Updates(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update standards: %w", result.Error)
	}
	return nil
}

func (r *StandardsRepo) modelToDomain(m *models.Standards) *project.Standards {
	return &project.Standards{
		ID:                m.ID,
		ProjectID:         m.ProjectID,
		BranchStrategy:    m.BranchStrategy,
		TechStack:         m.TechStack,
		CodingConventions: m.CodingConventions,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
}
