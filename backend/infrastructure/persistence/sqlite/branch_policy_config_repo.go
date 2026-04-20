package sqlite

import (
	"fmt"

	"github.com/mengri/nbcoder/domain/project"
	"github.com/mengri/nbcoder/infrastructure/database/models"
	"gorm.io/gorm"
)

type BranchPolicyConfigRepo struct {
	db *gorm.DB
}

func NewBranchPolicyConfigRepo(db *gorm.DB) project.BranchPolicyConfigRepo {
	return &BranchPolicyConfigRepo{db: db}
}

func (r *BranchPolicyConfigRepo) Save(c *project.BranchPolicyConfig) error {
	model := &models.BranchPolicyConfig{
		ID:          c.ID,
		ProjectID:   c.ProjectID,
		PolicyName:  "",
		PolicyConfig: "",
		RequireReviews: false,
		MinReviewers:   0,
		RequireTests:   false,
		AutoMergeEnabled: false,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}

	result := r.db.Save(model)
	if result.Error != nil {
		return fmt.Errorf("failed to save branch policy config: %w", result.Error)
	}
	return nil
}

func (r *BranchPolicyConfigRepo) FindByProjectID(projectID string) ([]*project.BranchPolicyConfig, error) {
	var models []models.BranchPolicyConfig
	result := r.db.Where("project_id = ?", projectID).Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find branch policy configs by project id: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *BranchPolicyConfigRepo) FindByID(id string) (*project.BranchPolicyConfig, error) {
	var model models.BranchPolicyConfig
	result := r.db.First(&model, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find branch policy config by id: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *BranchPolicyConfigRepo) Update(c *project.BranchPolicyConfig) error {
	model := &models.BranchPolicyConfig{
		ProjectID: c.ProjectID,
	}

	result := r.db.Model(&models.BranchPolicyConfig{}).Where("id = ?", c.ID).Updates(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update branch policy config: %w", result.Error)
	}
	return nil
}

func (r *BranchPolicyConfigRepo) Delete(id string) error {
	result := r.db.Delete(&models.BranchPolicyConfig{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete branch policy config: %w", result.Error)
	}
	return nil
}

func (r *BranchPolicyConfigRepo) FindByPolicyName(projectID, policyName string) (*project.BranchPolicyConfig, error) {
	var model models.BranchPolicyConfig
	result := r.db.Where("project_id = ? AND policy_name = ?", projectID, policyName).First(&model)
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
		ProjectID:   m.ProjectID,
		Pattern:     "",
		Description: "",
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
