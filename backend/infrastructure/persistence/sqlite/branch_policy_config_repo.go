package sqlite

import (
	"fmt"

	"github.com/mengri/nbcoder/domain/project"
	"github.com/mengri/nbcoder/infrastructure/database/models"
	"gorm.io/gorm"
)

type BranchPolicyConfigRepo struct {
	dbProvider DBProvider
}

func NewBranchPolicyConfigRepo(dbProvider DBProvider) project.BranchPolicyConfigRepo {
	return &BranchPolicyConfigRepo{dbProvider: dbProvider}
}

func (r *BranchPolicyConfigRepo) getDB(projectName string) (*gorm.DB, error) {
	if projectName == "" {
		return r.dbProvider.GetGlobalDB(), nil
	}
	return r.dbProvider.GetProjectDB(projectName)
}

func (r *BranchPolicyConfigRepo) Save(c *project.BranchPolicyConfig) error {
	db, err := r.getDB(c.ProjectName)
	if err != nil {
		return err
	}

	model := &models.BranchPolicyConfig{
		ID:               c.ID,
		ProjectName:      c.ProjectName,
		PolicyName:       c.Pattern,
		PolicyConfig:     c.Description,
		RequireReviews:   false,
		MinReviewers:     1,
		RequireTests:     false,
		AutoMergeEnabled: false,
		CreatedAt:        c.CreatedAt,
		UpdatedAt:        c.UpdatedAt,
	}

	result := db.Save(model)
	if result.Error != nil {
		return fmt.Errorf("failed to save branch policy config: %w", result.Error)
	}
	return nil
}

func (r *BranchPolicyConfigRepo) FindByProjectName(projectName string) ([]*project.BranchPolicyConfig, error) {
	db, err := r.getDB(projectName)
	if err != nil {
		return nil, err
	}

	var models []models.BranchPolicyConfig
	result := db.Where("project_name = ?", projectName).Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find branch policy configs by project name: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *BranchPolicyConfigRepo) FindByID(id string, projectName string) (*project.BranchPolicyConfig, error) {
	db, err := r.getDB(projectName)
	if err != nil {
		return nil, err
	}

	var model models.BranchPolicyConfig
	result := db.First(&model, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find branch policy config by id: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *BranchPolicyConfigRepo) Update(c *project.BranchPolicyConfig) error {
	db, err := r.getDB(c.ProjectName)
	if err != nil {
		return err
	}

	model := &models.BranchPolicyConfig{
		ProjectName:  c.ProjectName,
		PolicyName:   c.Pattern,
		PolicyConfig: c.Description,
	}

	result := db.Model(&models.BranchPolicyConfig{}).Where("id = ?", c.ID).Updates(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update branch policy config: %w", result.Error)
	}
	return nil
}

func (r *BranchPolicyConfigRepo) Delete(id string, projectName string) error {
	db, err := r.getDB(projectName)
	if err != nil {
		return err
	}

	result := db.Delete(&models.BranchPolicyConfig{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete branch policy config: %w", result.Error)
	}
	return nil
}

func (r *BranchPolicyConfigRepo) FindByPolicyName(projectName, policyName string) (*project.BranchPolicyConfig, error) {
	db, err := r.getDB(projectName)
	if err != nil {
		return nil, err
	}

	var model models.BranchPolicyConfig
	result := db.Where("project_name = ? AND policy_name = ?", projectName, policyName).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find branch policy config by name: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *BranchPolicyConfigRepo) modelToDomain(m *models.BranchPolicyConfig) *project.BranchPolicyConfig {
	return &project.BranchPolicyConfig{
		ID:          m.ID,
		ProjectName: m.ProjectName,
		Pattern:     m.PolicyName,
		Description: m.PolicyConfig,
		IsDefault:   false,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func (r *BranchPolicyConfigRepo) modelsToDomain(models []models.BranchPolicyConfig) []*project.BranchPolicyConfig {
	result := make([]*project.BranchPolicyConfig, len(models))
	for i, m := range models {
		result[i] = r.modelToDomain(&m)
	}
	return result
}
