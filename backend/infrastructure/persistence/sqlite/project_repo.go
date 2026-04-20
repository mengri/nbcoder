package sqlite

import (
	"fmt"

	"github.com/mengri/nbcoder/domain/project"
	"github.com/mengri/nbcoder/infrastructure/database/models"
	"gorm.io/gorm"
)

type ProjectRepo struct {
	db *gorm.DB
}

func NewProjectRepo(db *gorm.DB) project.ProjectRepo {
	return &ProjectRepo{db: db}
}

func (r *ProjectRepo) Save(p *project.Project) error {
	model := &models.Project{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		RepoURL:     p.RepoURL,
		Status:      string(p.Status),
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}

	result := r.db.Save(model)
	if result.Error != nil {
		return fmt.Errorf("failed to save project: %w", result.Error)
	}
	return nil
}

func (r *ProjectRepo) FindByID(id string) (*project.Project, error) {
	var model models.Project
	result := r.db.First(&model, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find project by id: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *ProjectRepo) FindAll() ([]*project.Project, error) {
	var models []models.Project
	result := r.db.Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find all projects: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *ProjectRepo) Update(p *project.Project) error {
	model := &models.Project{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		RepoURL:     p.RepoURL,
		Status:      string(p.Status),
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}

	result := r.db.Model(&models.Project{}).Where("id = ?", p.ID).Updates(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update project: %w", result.Error)
	}
	return nil
}

func (r *ProjectRepo) Delete(id string) error {
	result := r.db.Delete(&models.Project{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete project: %w", result.Error)
	}
	return nil
}

func (r *ProjectRepo) FindByStatus(status project.ProjectStatus) ([]*project.Project, error) {
	var models []models.Project
	result := r.db.Where("status = ?", status).Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find projects by status: %w", result.Error)
	}
	return r.modelsToDomain(models), nil
}

func (r *ProjectRepo) FindByName(name string) (*project.Project, error) {
	var model models.Project
	result := r.db.Where("name = ?", name).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find project by name: %w", result.Error)
	}
	return r.modelToDomain(&model), nil
}

func (r *ProjectRepo) modelToDomain(m *models.Project) *project.Project {
	return &project.Project{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		RepoURL:     m.RepoURL,
		Status:      project.ProjectStatus(m.Status),
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func (r *ProjectRepo) modelsToDomain(models []models.Project) []*project.Project {
	result := make([]*project.Project, len(models))
	for i, m := range models {
		result[i] = r.modelToDomain(&m)
	}
	return result
}
