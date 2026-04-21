package sqlite

import (
	"fmt"

	"github.com/mengri/nbcoder/domain/project"
	"github.com/mengri/nbcoder/infrastructure/database/models"
	"gorm.io/gorm"
)

type ProjectLifecycleRepo struct {
	dbProvider DBProvider
}

func NewProjectLifecycleRepo(dbProvider DBProvider) project.ProjectLifecycleRepo {
	return &ProjectLifecycleRepo{dbProvider: dbProvider}
}

func (r *ProjectLifecycleRepo) getDB(projectName string) (*gorm.DB, error) {
	if projectName == "" {
		return r.dbProvider.GetGlobalDB(), nil
	}
	return r.dbProvider.GetProjectDB(projectName)
}

func (r *ProjectLifecycleRepo) Save(l *project.ProjectLifecycle) error {
	db, err := r.getDB(l.ProjectName)
	if err != nil {
		return err
	}

	model := &models.ProjectLifecycle{
		ID:          l.ID,
		ProjectName: l.ProjectName,
		Status:      string(l.Status),
		ActivatedAt: l.ActivatedAt,
		SuspendedAt: l.SuspendedAt,
		ArchivedAt:  l.ArchivedAt,
		CreatedAt:   l.CreatedAt,
		UpdatedAt:   l.UpdatedAt,
	}

	result := db.Save(model)
	if result.Error != nil {
		return fmt.Errorf("failed to save project lifecycle: %w", result.Error)
	}
	return nil
}

func (r *ProjectLifecycleRepo) FindByProjectName(projectName string) (*project.ProjectLifecycle, error) {
	db, err := r.getDB(projectName)
	if err != nil {
		return nil, err
	}

	var model models.ProjectLifecycle
	result := db.Where("project_name = ?", projectName).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find project lifecycle by project name: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *ProjectLifecycleRepo) Update(l *project.ProjectLifecycle) error {
	db, err := r.getDB(l.ProjectName)
	if err != nil {
		return err
	}

	model := &models.ProjectLifecycle{
		ID:          l.ID,
		ProjectName: l.ProjectName,
		Status:      string(l.Status),
		ActivatedAt: l.ActivatedAt,
		SuspendedAt: l.SuspendedAt,
		ArchivedAt:  l.ArchivedAt,
		CreatedAt:   l.CreatedAt,
		UpdatedAt:   l.UpdatedAt,
	}

	result := db.Model(&models.ProjectLifecycle{}).Where("id = ?", l.ID).Updates(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update project lifecycle: %w", result.Error)
	}
	return nil
}

func (r *ProjectLifecycleRepo) FindByStatus(status string) ([]*project.ProjectLifecycle, error) {
	db, err := r.getDB("")
	if err != nil {
		return nil, err
	}

	var models []models.ProjectLifecycle
	result := db.Where("status = ?", status).Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find project lifecycles by status: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *ProjectLifecycleRepo) modelToDomain(m *models.ProjectLifecycle) *project.ProjectLifecycle {
	return &project.ProjectLifecycle{
		ID:          m.ID,
		ProjectName: m.ProjectName,
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
