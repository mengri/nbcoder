package sqlite

import (
	"fmt"

	"github.com/mengri/nbcoder/domain/airuntime"
	"github.com/mengri/nbcoder/infrastructure/database/models"
	"gorm.io/gorm"
)

type ProviderRepo struct {
	dbProvider DBProvider
}

func NewProviderRepo(dbProvider DBProvider) airuntime.ProviderRepo {
	return &ProviderRepo{dbProvider: dbProvider}
}

func (r *ProviderRepo) getDB() (*gorm.DB, error) {
	return r.dbProvider.GetGlobalDB(), nil
}

func (r *ProviderRepo) Save(provider *airuntime.Provider) error {
	db, err := r.getDB()
	if err != nil {
		return err
	}

	model := &models.Provider{
		ID:        provider.ID,
		Name:      provider.Name,
		APIKeyRef: provider.APIKeyRef,
		BaseURL:   "",
		IsActive:  true,
	}

	result := db.Save(model)
	if result.Error != nil {
		return fmt.Errorf("failed to save provider: %w", result.Error)
	}
	return nil
}

func (r *ProviderRepo) FindByID(id string) (*airuntime.Provider, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var model models.Provider
	result := db.Preload("Models").First(&model, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find provider by id: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *ProviderRepo) FindAll() ([]*airuntime.Provider, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var models []models.Provider
	result := db.Preload("Models").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find all providers: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *ProviderRepo) Update(provider *airuntime.Provider) error {
	db, err := r.getDB()
	if err != nil {
		return err
	}

	model := &models.Provider{
		Name:      provider.Name,
		APIKeyRef: provider.APIKeyRef,
	}

	result := db.Model(&models.Provider{}).Where("id = ?", provider.ID).Updates(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update provider: %w", result.Error)
	}
	return nil
}

func (r *ProviderRepo) FindByName(name string) (*airuntime.Provider, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var model models.Provider
	result := db.Preload("Models").Where("name = ?", name).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find provider by name: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *ProviderRepo) FindActive() ([]*airuntime.Provider, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var models []models.Provider
	result := db.Preload("Models").Where("is_active = ?", true).Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find active providers: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *ProviderRepo) modelToDomain(m *models.Provider) *airuntime.Provider {
	provider := &airuntime.Provider{
		ID:        m.ID,
		Name:      m.Name,
		APIKeyRef: m.APIKeyRef,
		Models:    make([]*airuntime.Model, len(m.Models)),
	}

	for i, model := range m.Models {
		provider.Models[i] = &airuntime.Model{
			ID:         model.ID,
			Name:       model.Name,
			ProviderID: model.ProviderID,
			ModelType:  model.ModelType,
			Meta:       map[string]interface{}(model.Meta),
		}
	}

	return provider
}

func (r *ProviderRepo) modelsToDomain(models []models.Provider) []*airuntime.Provider {
	result := make([]*airuntime.Provider, len(models))
	for i, m := range models {
		result[i] = r.modelToDomain(&m)
	}
	return result
}

type ModelRepo struct {
	dbProvider DBProvider
}

func NewModelRepo(dbProvider DBProvider) *ModelRepo {
	return &ModelRepo{dbProvider: dbProvider}
}

func (r *ModelRepo) getDB() (*gorm.DB, error) {
	return r.dbProvider.GetGlobalDB(), nil
}

func (r *ModelRepo) Save(model *airuntime.Model) error {
	db, err := r.getDB()
	if err != nil {
		return err
	}

	dbModel := &models.Model{
		ID:         model.ID,
		Name:       model.Name,
		ProviderID: model.ProviderID,
		ModelType:  model.ModelType,
		Meta:       models.JSONMap(model.Meta),
		IsActive:   true,
	}

	result := db.Save(dbModel)
	if result.Error != nil {
		return fmt.Errorf("failed to save model: %w", result.Error)
	}
	return nil
}

func (r *ModelRepo) FindByID(id string) (*airuntime.Model, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var model models.Model
	result := db.First(&model, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find model by id: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *ModelRepo) FindByProviderID(providerID string) ([]*airuntime.Model, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var models []models.Model
	result := db.Where("provider_id = ?", providerID).Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find models by provider id: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *ModelRepo) Update(model *airuntime.Model) error {
	db, err := r.getDB()
	if err != nil {
		return err
	}

	dbModel := &models.Model{
		Name:      model.Name,
		ModelType: model.ModelType,
		Meta:      models.JSONMap(model.Meta),
	}

	result := db.Model(&models.Model{}).Where("id = ?", model.ID).Updates(dbModel)
	if result.Error != nil {
		return fmt.Errorf("failed to update model: %w", result.Error)
	}
	return nil
}

func (r *ModelRepo) modelToDomain(m *models.Model) *airuntime.Model {
	return &airuntime.Model{
		ID:         m.ID,
		Name:       m.Name,
		ProviderID: m.ProviderID,
		ModelType:  m.ModelType,
		Meta:       map[string]interface{}(m.Meta),
	}
}

func (r *ModelRepo) modelsToDomain(models []models.Model) []*airuntime.Model {
	result := make([]*airuntime.Model, len(models))
	for i, m := range models {
		result[i] = r.modelToDomain(&m)
	}
	return result
}
