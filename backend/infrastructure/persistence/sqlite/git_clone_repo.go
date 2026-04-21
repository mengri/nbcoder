package sqlite

import (
	"fmt"

	"github.com/mengri/nbcoder/domain/clonepool"
	"github.com/mengri/nbcoder/domain/git"
	"github.com/mengri/nbcoder/infrastructure/database/models"
	"gorm.io/gorm"
)

type CloneInstanceRepo struct {
	dbProvider DBProvider
}

func NewCloneInstanceRepo(dbProvider DBProvider) clonepool.CloneInstanceRepo {
	return &CloneInstanceRepo{dbProvider: dbProvider}
}

func (r *CloneInstanceRepo) getDB() (*gorm.DB, error) {
	return r.dbProvider.GetGlobalDB(), nil
}

func (r *CloneInstanceRepo) Save(instance *clonepool.CloneInstance) error {
	db, err := r.getDB()
	if err != nil {
		return err
	}

	model := &models.CloneInstance{
		ID:         instance.ID,
		Name:       "Clone Instance",
		ProjectName: "",
		SourcePath: "",
		TargetPath: "",
		Status:     string(instance.Status),
		IsHealthy:  true,
		LastUsedAt: nil,
		CreatedAt:  instance.UpdatedAt,
		UpdatedAt:  instance.UpdatedAt,
	}

	result := db.Save(model)
	if result.Error != nil {
		return fmt.Errorf("failed to save clone instance: %w", result.Error)
	}
	return nil
}

func (r *CloneInstanceRepo) FindByID(id string) (*clonepool.CloneInstance, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var model models.CloneInstance
	result := db.First(&model, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find clone instance by id: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *CloneInstanceRepo) FindByRepositoryID(repositoryID string) ([]*clonepool.CloneInstance, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var models []models.CloneInstance
	result := db.Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find clone instances: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *CloneInstanceRepo) FindByStatus(status clonepool.CloneInstanceStatus) ([]*clonepool.CloneInstance, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var models []models.CloneInstance
	result := db.Where("status = ?", string(status)).Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find clone instances by status: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *CloneInstanceRepo) Update(instance *clonepool.CloneInstance) error {
	db, err := r.getDB()
	if err != nil {
		return err
	}

	model := &models.CloneInstance{
		ID:         instance.ID,
		Name:       "Clone Instance",
		ProjectName: "",
		SourcePath: "",
		TargetPath: "",
		Status:     string(instance.Status),
		IsHealthy:  true,
		LastUsedAt: nil,
		CreatedAt:  instance.UpdatedAt,
		UpdatedAt:  instance.UpdatedAt,
	}

	result := db.Model(&models.CloneInstance{}).Where("id = ?", instance.ID).Updates(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update clone instance: %w", result.Error)
	}
	return nil
}

func (r *CloneInstanceRepo) modelToDomain(m *models.CloneInstance) *clonepool.CloneInstance {
	return &clonepool.CloneInstance{
		ID:           m.ID,
		RepositoryID: "",
		Status:       clonepool.CloneInstanceStatus(m.Status),
		AssignedTask: "",
		UpdatedAt:    m.UpdatedAt,
	}
}

func (r *CloneInstanceRepo) modelsToDomain(models []models.CloneInstance) []*clonepool.CloneInstance {
	result := make([]*clonepool.CloneInstance, len(models))
	for i, m := range models {
		result[i] = r.modelToDomain(&m)
	}
	return result
}

type PullRequestRepo struct {
	dbProvider DBProvider
}

func NewPullRequestRepo(dbProvider DBProvider) git.PullRequestRepo {
	return &PullRequestRepo{dbProvider: dbProvider}
}

func (r *PullRequestRepo) getDB(projectName string) (*gorm.DB, error) {
	if projectName == "" {
		return r.dbProvider.GetGlobalDB(), nil
	}
	return r.dbProvider.GetProjectDB(projectName)
}

func (r *PullRequestRepo) Save(pr *git.PullRequest) error {
	db, err := r.getDB(pr.ProjectName)
	if err != nil {
		return err
	}

	model := &models.PullRequest{
		ID:             pr.ID,
		Title:          pr.Title,
		Description:    pr.Description,
		SourceBranch:   pr.SourceBranch,
		TargetBranch:   pr.TargetBranch,
		Status:         string(pr.Status),
		ProjectName:    pr.ProjectName,
		Author:         pr.Author,
		GeneratedDesc:  pr.GeneratedDesc,
		SquashCommitMsg: pr.SquashCommitMsg,
		CreatedAt:      pr.CreatedAt,
		UpdatedAt:      pr.UpdatedAt,
	}

	result := db.Save(model)
	if result.Error != nil {
		return fmt.Errorf("failed to save pull request: %w", result.Error)
	}
	return nil
}

func (r *PullRequestRepo) FindByID(id string, projectName string) (*git.PullRequest, error) {
	db, err := r.getDB(projectName)
	if err != nil {
		return nil, err
	}

	var model models.PullRequest
	result := db.First(&model, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find pull request by id: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *PullRequestRepo) FindByProjectName(projectName string) ([]*git.PullRequest, error) {
	db, err := r.getDB(projectName)
	if err != nil {
		return nil, err
	}

	var models []models.PullRequest
	result := db.Where("project_name = ?", projectName).Order("created_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find pull requests by project name: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *PullRequestRepo) FindByBranch(sourceBranch string) ([]*git.PullRequest, error) {
	db, err := r.getDB("")
	if err != nil {
		return nil, err
	}

	var models []models.PullRequest
	result := db.Where("source_branch = ?", sourceBranch).Order("created_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find pull requests by source branch: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *PullRequestRepo) Update(pr *git.PullRequest) error {
	db, err := r.getDB(pr.ProjectName)
	if err != nil {
		return err
	}

	model := &models.PullRequest{
		ID:             pr.ID,
		Title:          pr.Title,
		Description:    pr.Description,
		SourceBranch:   pr.SourceBranch,
		TargetBranch:   pr.TargetBranch,
		Status:         string(pr.Status),
		ProjectName:    pr.ProjectName,
		Author:         pr.Author,
		GeneratedDesc:  pr.GeneratedDesc,
		SquashCommitMsg: pr.SquashCommitMsg,
		CreatedAt:      pr.CreatedAt,
		UpdatedAt:      pr.UpdatedAt,
	}

	result := db.Model(&models.PullRequest{}).Where("id = ?", pr.ID).Updates(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update pull request: %w", result.Error)
	}
	return nil
}

func (r *PullRequestRepo) modelToDomain(m *models.PullRequest) *git.PullRequest {
	return &git.PullRequest{
		ID:             m.ID,
		Title:          m.Title,
		Description:    m.Description,
		SourceBranch:   m.SourceBranch,
		TargetBranch:   m.TargetBranch,
		Status:         git.PullRequestStatus(m.Status),
		ProjectName:    m.ProjectName,
		Author:         m.Author,
		GeneratedDesc:  m.GeneratedDesc,
		SquashCommitMsg: m.SquashCommitMsg,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

func (r *PullRequestRepo) modelsToDomain(models []models.PullRequest) []*git.PullRequest {
	result := make([]*git.PullRequest, len(models))
	for i, m := range models {
		result[i] = r.modelToDomain(&m)
	}
	return result
}
