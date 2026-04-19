package sqlite

import (
	"fmt"

	"github.com/mengri/nbcoder/domain/requirement"
	"github.com/mengri/nbcoder/infrastructure/database/models"
	"gorm.io/gorm"
)

type CardRepo struct {
	db *gorm.DB
}

func NewCardRepo(db *gorm.DB) requirement.CardRepo {
	return &CardRepo{db: db}
}

func (r *CardRepo) Save(c *requirement.Card) error {
	model := &models.Card{
		ID:               c.ID,
		Title:            c.Title,
		Description:      c.Description,
		Original:         c.Original,
		Status:           string(c.Status),
		Priority:         string(c.Priority),
		StructuredOutput: c.StructuredOutput,
		PipelineID:       c.PipelineID,
		ProjectID:        c.ProjectID,
		CreatedAt:        c.CreatedAt,
		UpdatedAt:        c.UpdatedAt,
		SupersededBy:     c.SupersededBy,
	}

	result := r.db.Save(model)
	if result.Error != nil {
		return fmt.Errorf("failed to save card: %w", result.Error)
	}
	return nil
}

func (r *CardRepo) FindByID(id string) (*requirement.Card, error) {
	var model models.Card
	result := r.db.First(&model, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find card by id: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *CardRepo) FindAll() ([]*requirement.Card, error) {
	var models []models.Card
	result := r.db.Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find all cards: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *CardRepo) Update(c *requirement.Card) error {
	model := &models.Card{
		ID:               c.ID,
		Title:            c.Title,
		Description:      c.Description,
		Original:         c.Original,
		Status:           string(c.Status),
		Priority:         string(c.Priority),
		StructuredOutput: c.StructuredOutput,
		PipelineID:       c.PipelineID,
		ProjectID:        c.ProjectID,
		CreatedAt:        c.CreatedAt,
		UpdatedAt:        c.UpdatedAt,
		SupersededBy:     c.SupersededBy,
	}

	result := r.db.Model(&models.Card{}).Where("id = ?", c.ID).Updates(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update card: %w", result.Error)
	}
	return nil
}

func (r *CardRepo) Delete(id string) error {
	result := r.db.Delete(&models.Card{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete card: %w", result.Error)
	}
	return nil
}

func (r *CardRepo) FindByProjectID(projectID string) ([]*requirement.Card, error) {
	var models []models.Card
	result := r.db.Where("project_id = ?", projectID).Order("created_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find cards by project id: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *CardRepo) FindByStatus(status string) ([]*requirement.Card, error) {
	var models []models.Card
	result := r.db.Where("status = ?", status).Order("created_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find cards by status: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *CardRepo) FindByPriority(priority string) ([]*requirement.Card, error) {
	var models []models.Card
	result := r.db.Where("priority = ?", priority).Order("created_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find cards by priority: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *CardRepo) FindByPipelineID(pipelineID string) (*requirement.Card, error) {
	var model models.Card
	result := r.db.Where("pipeline_id = ?", pipelineID).First(&model)
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
		ProjectID:        m.ProjectID,
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
