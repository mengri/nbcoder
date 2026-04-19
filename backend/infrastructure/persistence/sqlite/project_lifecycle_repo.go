package sqlite

import (
	"fmt"

	"github.com/mengri/nbcoder/domain/project"
	"github.com/mengri/nbcoder/infrastructure/database/models"
	"gorm.io/gorm"
)

type ProjectLifecycleRepo struct {
	db *gorm.DB
}

func NewProjectLifecycleRepo(db *gorm.DB) project.ProjectLifecycleRepo {
	return &ProjectLifecycleRepo{db: db}
}

func (r *ProjectLifecycleRepo) Save(l *project.ProjectLifecycle) error {
	model := &models.ProjectLifecycle{
		ID:          l.ID,
		ProjectID:   l.ProjectID,
		Status:      string(l.Status),
		ActivatedAt: l.ActivatedAt,
		SuspendedAt: l.SuspendedAt,
		ArchivedAt:  l.ArchivedAt,
		CreatedAt:   l.CreatedAt,
		UpdatedAt:   l.UpdatedAt,
	}

	result := r.db.Save(model)
	if result.Error != nil {
		return fmt.Errorf("failed to save project lifecycle: %w", result.Error)
	}
	return nil
}

func (r *ProjectLifecycleRepo) FindByProjectID(projectID string) (*project.ProjectLifecycle, error) {
	var model models.ProjectLifecycle
	result := r.db.Where("project_id = ?", projectID).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find project lifecycle by project id: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *ProjectLifecycleRepo) Update(l *project.ProjectLifecycle) error {
	model := &models.ProjectLifecycle{
		ID:          l.ID,
		ProjectID:   l.ProjectID,
		Status:      string(l.Status),
		ActivatedAt: l.ActivatedAt,
		SuspendedAt: l.SuspendedAt,
		ArchivedAt:  l.ArchivedAt,
		CreatedAt:   l.CreatedAt,
		UpdatedAt:   l.UpdatedAt,
	}

	result := r.db.Model(&models.ProjectLifecycle{}).Where("id = ?", l.ID).Updates(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update project lifecycle: %w", result.Error)
	}
	return nil
}

func (r *ProjectLifecycleRepo) FindByStatus(status string) ([]*project.ProjectLifecycle, error) {
	var models []models.ProjectLifecycle
	result := r.db.Where("status = ?", status).Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find project lifecycles by status: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *ProjectLifecycleRepo) modelToDomain(m *models.ProjectLifecycle) *project.ProjectLifecycle {
	return &project.ProjectLifecycle{
		ID:          m.ID,
		ProjectID:   m.ProjectID,
		Status:      project.LifecycleStatus(m.Status),
		ActivatedAt: m.ActivatedAt,
		SuspendedAt: m.SuspendedAt,
		ArchivedAt:  m.ArchivedAt,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func (r *ProjectLifecycleRepo) modelsToDomain(models []models.ProjectLifecycle) []*project.ProjectLifecycle {
	result := make([]*project.ProjectLifecycle, len(models))
	for i, m := range models {
		result[i] = r.modelToDomain(&m)
	}
	return result
}
