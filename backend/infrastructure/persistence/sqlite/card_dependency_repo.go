package sqlite

import (
	"fmt"

	"github.com/mengri/nbcoder/domain/requirement"
	"github.com/mengri/nbcoder/infrastructure/database/models"
	"gorm.io/gorm"
)

type CardDependencyRepo struct {
	db *gorm.DB
}

func NewCardDependencyRepo(db *gorm.DB) requirement.CardDependencyRepo {
	return &CardDependencyRepo{db: db}
}

func (r *CardDependencyRepo) Save(d *requirement.CardDependency) error {
	model := &models.CardDependency{
		ID:              d.ID,
		CardID:          d.CardID,
		DependsOnCardID: d.DependsOnID,
		DependencyType:  string(d.Type),
	}

	result := r.db.Save(model)
	if result.Error != nil {
		return fmt.Errorf("failed to save card dependency: %w", result.Error)
	}
	return nil
}

func (r *CardDependencyRepo) FindByID(id string) (*requirement.CardDependency, error) {
	var model models.CardDependency
	result := r.db.First(&model, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find card dependency by id: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *CardDependencyRepo) FindAll() ([]*requirement.CardDependency, error) {
	var models []models.CardDependency
	result := r.db.Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find all card dependencies: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *CardDependencyRepo) Update(d *requirement.CardDependency) error {
	model := &models.CardDependency{
		CardID:          d.CardID,
		DependsOnCardID: d.DependsOnID,
		DependencyType:  string(d.Type),
	}

	result := r.db.Model(&models.CardDependency{}).Where("id = ?", d.ID).Updates(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update card dependency: %w", result.Error)
	}
	return nil
}

func (r *CardDependencyRepo) Delete(id string) error {
	result := r.db.Delete(&models.CardDependency{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete card dependency: %w", result.Error)
	}
	return nil
}

func (r *CardDependencyRepo) DeleteByCardID(cardID string) error {
	result := r.db.Where("card_id = ?", cardID).Delete(&models.CardDependency{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete card dependencies by card id: %w", result.Error)
	}
	return nil
}

func (r *CardDependencyRepo) FindByCardID(cardID string) ([]*requirement.CardDependency, error) {
	var models []models.CardDependency
	result := r.db.Where("card_id = ?", cardID).Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find dependencies by card id: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *CardDependencyRepo) FindByDependsOnCardID(dependsOnCardID string) ([]*requirement.CardDependency, error) {
	var models []models.CardDependency
	result := r.db.Where("depends_on_card_id = ?", dependsOnCardID).Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find dependencies by depends on card id: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *CardDependencyRepo) FindByDependsOnID(dependsOnID string) ([]*requirement.CardDependency, error) {
	return r.FindByDependsOnCardID(dependsOnID)
}

func (r *CardDependencyRepo) FindByCardIDAndDependsOnCardID(cardID, dependsOnCardID string) (*requirement.CardDependency, error) {
	var model models.CardDependency
	result := r.db.Where("card_id = ? AND depends_on_card_id = ?", cardID, dependsOnCardID).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find dependency by card id and depends on card id: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *CardDependencyRepo) modelToDomain(m *models.CardDependency) *requirement.CardDependency {
	return &requirement.CardDependency{
		ID:          m.ID,
		CardID:      m.CardID,
		DependsOnID: m.DependsOnCardID,
		Type:        requirement.DependencyType(m.DependencyType),
	}
}

func (r *CardDependencyRepo) modelsToDomain(models []models.CardDependency) []*requirement.CardDependency {
	result := make([]*requirement.CardDependency, len(models))
	for i, m := range models {
		result[i] = r.modelToDomain(&m)
	}
	return result
}
