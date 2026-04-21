package sqlite

import (
	"fmt"

	"github.com/mengri/nbcoder/domain/project"
	"github.com/mengri/nbcoder/infrastructure/database/models"
	"gorm.io/gorm"
)

type StandardsRepo struct {
	dbProvider DBProvider
}

func NewStandardsRepo(dbProvider DBProvider) project.StandardsRepo {
	return &StandardsRepo{dbProvider: dbProvider}
}

func (r *StandardsRepo) getDB(projectName string) (*gorm.DB, error) {
	if projectName == "" {
		return r.dbProvider.GetGlobalDB(), nil
	}
	return r.dbProvider.GetProjectDB(projectName)
}

func (r *StandardsRepo) Save(s *project.Standards) error {
	db, err := r.getDB(s.ProjectName)
	if err != nil {
		return err
	}

	model := &models.Standards{
		ID:                s.ID,
		ProjectName:       s.ProjectName,
		BranchStrategy:    s.BranchStrategy,
		TechStack:         s.TechStack,
		CodingConventions: s.CodingConventions,
		CreatedAt:         s.CreatedAt,
		UpdatedAt:         s.UpdatedAt,
	}

	result := db.Save(model)
	if result.Error != nil {
		return fmt.Errorf("failed to save standards: %w", result.Error)
	}
	return nil
}

func (r *StandardsRepo) FindByProjectName(projectName string) (*project.Standards, error) {
	db, err := r.getDB(projectName)
	if err != nil {
		return nil, err
	}

	var model models.Standards
	result := db.Where("project_name = ?", projectName).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find standards by project name: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *StandardsRepo) Update(s *project.Standards) error {
	db, err := r.getDB(s.ProjectName)
	if err != nil {
		return err
	}

	model := &models.Standards{
		ID:                s.ID,
		ProjectName:       s.ProjectName,
		BranchStrategy:    s.BranchStrategy,
		TechStack:         s.TechStack,
		CodingConventions: s.CodingConventions,
		CreatedAt:         s.CreatedAt,
		UpdatedAt:         s.UpdatedAt,
	}

	result := db.Model(&models.Standards{}).Where("id = ?", s.ID).Updates(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update standards: %w", result.Error)
	}
	return nil
}

func (r *StandardsRepo) modelToDomain(m *models.Standards) *project.Standards {
	return &project.Standards{
		ID:                m.ID,
		ProjectName:       m.ProjectName,
		BranchStrategy:    m.BranchStrategy,
		TechStack:         m.TechStack,
		CodingConventions: m.CodingConventions,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
}
