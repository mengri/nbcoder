package sqlite

import (
	"fmt"

	"github.com/mengri/nbcoder/domain/requirement"
	"github.com/mengri/nbcoder/infrastructure/database/models"
	"gorm.io/gorm"
)

type CardRepo struct {
	dbProvider DBProvider
}

func NewCardRepo(dbProvider DBProvider) requirement.CardRepo {
	return &CardRepo{dbProvider: dbProvider}
}

func (r *CardRepo) getDB(projectName string) (*gorm.DB, error) {
	if projectName == "" {
		return r.dbProvider.GetGlobalDB(), nil
	}
	return r.dbProvider.GetProjectDB(projectName)
}

func (r *CardRepo) Save(c *requirement.Card) error {
	db, err := r.getDB(c.ProjectName)
	if err != nil {
		return err
	}

	model := &models.Card{
		ID:               c.ID,
		Title:            c.Title,
		Description:      c.Description,
		Original:         c.Original,
		Status:           string(c.Status),
		Priority:         string(c.Priority),
		StructuredOutput: c.StructuredOutput,
		PipelineID:       c.PipelineID,
		ProjectName:      c.ProjectName,
		CreatedAt:        c.CreatedAt,
		UpdatedAt:        c.UpdatedAt,
		SupersededBy:     c.SupersededBy,
	}

	result := db.Save(model)
	if result.Error != nil {
		return fmt.Errorf("failed to save card: %w", result.Error)
	}
	return nil
}

func (r *CardRepo) FindByID(id string, projectName string) (*requirement.Card, error) {
	db, err := r.getDB(projectName)
	if err != nil {
		return nil, err
	}

	var model models.Card
	result := db.First(&model, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find card by id: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *CardRepo) FindAll() ([]*requirement.Card, error) {
	db, err := r.getDB("")
	if err != nil {
		return nil, err
	}

	var models []models.Card
	result := db.Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find all cards: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *CardRepo) Update(c *requirement.Card) error {
	db, err := r.getDB(c.ProjectName)
	if err != nil {
		return err
	}

	model := &models.Card{
		ID:               c.ID,
		Title:            c.Title,
		Description:      c.Description,
		Original:         c.Original,
		Status:           string(c.Status),
		Priority:         string(c.Priority),
		StructuredOutput: c.StructuredOutput,
		PipelineID:       c.PipelineID,
		ProjectName:      c.ProjectName,
		CreatedAt:        c.CreatedAt,
		UpdatedAt:        c.UpdatedAt,
		SupersededBy:     c.SupersededBy,
	}

	result := db.Model(&models.Card{}).Where("id = ?", c.ID).Updates(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update card: %w", result.Error)
	}
	return nil
}

func (r *CardRepo) Delete(id string, projectName string) error {
	db, err := r.getDB(projectName)
	if err != nil {
		return err
	}

	result := db.Delete(&models.Card{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete card: %w", result.Error)
	}
	return nil
}

func (r *CardRepo) FindByProjectName(projectName string) ([]*requirement.Card, error) {
	db, err := r.getDB(projectName)
	if err != nil {
		return nil, err
	}

	var models []models.Card
	result := db.Where("project_name = ?", projectName).Order("created_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find cards by project name: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *CardRepo) FindByStatus(status requirement.CardStatus) ([]*requirement.Card, error) {
	db, err := r.getDB("")
	if err != nil {
		return nil, err
	}

	var models []models.Card
	result := db.Where("status = ?", string(status)).Order("created_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find cards by status: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *CardRepo) FindByPriority(priority string) ([]*requirement.Card, error) {
	db, err := r.getDB("")
	if err != nil {
		return nil, err
	}

	var models []models.Card
	result := db.Where("priority = ?", priority).Order("created_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find cards by priority: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *CardRepo) FindByPipelineID(pipelineID string) (*requirement.Card, error) {
	db, err := r.getDB("")
	if err != nil {
		return nil, err
	}

	var model models.Card
	result := db.Where("pipeline_id = ?", pipelineID).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find card by pipeline id: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *CardRepo) modelToDomain(m *models.Card) *requirement.Card {
	return &requirement.Card{
		ID:               m.ID,
		Title:            m.Title,
		Description:      m.Description,
		Original:         m.Original,
		Status:           requirement.CardStatus(m.Status),
		Priority:         requirement.Priority(m.Priority),
		StructuredOutput: m.StructuredOutput,
		PipelineID:       m.PipelineID,
		ProjectName:      m.ProjectName,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
		SupersededBy:     m.SupersededBy,
	}
}

func (r *CardRepo) modelsToDomain(models []models.Card) []*requirement.Card {
	result := make([]*requirement.Card, len(models))
	for i, m := range models {
		result[i] = r.modelToDomain(&m)
	}
	return result
}
