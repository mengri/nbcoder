package sqlite

import (
	"fmt"
	"time"

	"github.com/mengri/nbcoder/domain/airuntime"
	"github.com/mengri/nbcoder/infrastructure/database/models"
	"gorm.io/gorm"
)

type ModelChainRepo struct {
	db *gorm.DB
}

func NewModelChainRepo(db *gorm.DB) airuntime.ChainRepo {
	return &ModelChainRepo{db: db}
}

func (r *ModelChainRepo) Save(chain *airuntime.Chain) error {
	model := &models.ModelChain{
		ID:        chain.ID,
		Name:      chain.Name,
		ModelID:   "",
		ChainType: "",
		Config:    map[string]interface{}{"models": chain.Models, "routes": chain.Routes},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result := r.db.Save(model)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *ModelChainRepo) FindByID(id string) (*airuntime.Chain, error) {
	var model models.ModelChain
	result := r.db.First(&model, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return r.modelToDomain(&model), nil
}

func (r *ModelChainRepo) FindAll() ([]*airuntime.Chain, error) {
	var models []models.ModelChain
	result := r.db.Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	return r.modelsToDomain(models), nil
}

func (r *ModelChainRepo) modelToDomain(m *models.ModelChain) *airuntime.Chain {
	return &airuntime.Chain{
		ID:     m.ID,
		Name:   m.Name,
		Models: []string{},
		Routes: []airuntime.ModelRoute{},
	}
}

func (r *ModelChainRepo) modelsToDomain(models []models.ModelChain) []*airuntime.Chain {
	result := make([]*airuntime.Chain, len(models))
	for i, m := range models {
		result[i] = r.modelToDomain(&m)
	}
	return result
}

type CallLogRepo struct {
	db *gorm.DB
}

func NewCallLogRepo(db *gorm.DB) airuntime.CallLogRepo {
	return &CallLogRepo{db: db}
}

func (r *CallLogRepo) Save(log *airuntime.CallLog) error {
	model := &models.CallLog{
		ID:         log.ID,
		AgentID:    log.AgentID,
		ModelID:    log.ModelID,
		CallType:   "API_CALL",
		Input:      log.Input,
		Output:     log.Output,
		TokensUsed: log.Tokens,
		Cost:       0,
		LatencyMs:  0,
		Status:     "SUCCESS",
		Error:      "",
		CreatedAt:  log.Timestamp,
		UpdatedAt:  time.Now(),
	}

	result := r.db.Save(model)
	if result.Error != nil {
		return fmt.Errorf("failed to save call log: %w", result.Error)
	}
	return nil
}

func (r *CallLogRepo) FindByAgentID(agentID string) ([]*airuntime.CallLog, error) {
	var models []models.CallLog
	result := r.db.Where("agent_id = ?", agentID).Order("created_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find call logs by agent id: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *CallLogRepo) FindByTimeRange(start, end time.Time) ([]*airuntime.CallLog, error) {
	var models []models.CallLog
	result := r.db.Where("created_at BETWEEN ? AND ?", start, end).Order("created_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find call logs by time range: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *CallLogRepo) FindByModelID(modelID string) ([]*airuntime.CallLog, error) {
	var models []models.CallLog
	result := r.db.Where("model_id = ?", modelID).Order("created_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find call logs by model id: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *CallLogRepo) FindByStatus(status string) ([]*airuntime.CallLog, error) {
	var models []models.CallLog
	result := r.db.Where("status = ?", status).Order("created_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find call logs by status: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *CallLogRepo) modelToDomain(m *models.CallLog) *airuntime.CallLog {
	return &airuntime.CallLog{
		ID:        m.ID,
		AgentID:   m.AgentID,
		ModelID:   m.ModelID,
		Input:     m.Input,
		Output:    m.Output,
		Tokens:    m.TokensUsed,
		Timestamp: m.CreatedAt,
	}
}

func (r *CallLogRepo) modelsToDomain(models []models.CallLog) []*airuntime.CallLog {
	result := make([]*airuntime.CallLog, len(models))
	for i, m := range models {
		result[i] = r.modelToDomain(&m)
	}
	return result
}
