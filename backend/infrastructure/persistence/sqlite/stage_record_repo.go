package sqlite

import (
	"fmt"

	"github.com/mengri/nbcoder/domain/pipeline"
	"github.com/mengri/nbcoder/infrastructure/database/models"
	"gorm.io/gorm"
)

type StageRecordRepo struct {
	db *gorm.DB
}

func NewStageRecordRepo(db *gorm.DB) pipeline.StageRecordRepo {
	return &StageRecordRepo{db: db}
}

func (r *StageRecordRepo) Save(record *pipeline.StageRecord) error {
	model := &models.StageRecord{
		ID:         record.ID,
		PipelineID: record.PipelineID,
		StageName:  record.StageName,
		Status:     string(record.Status),
		StartedAt:  record.StartedAt,
		CompletedAt: record.CompletedAt,
		Output:     record.Output,
		Error:      record.Error,
		CreatedAt:  record.CreatedAt,
		UpdatedAt:  record.UpdatedAt,
	}

	result := r.db.Save(model)
	if result.Error != nil {
		return fmt.Errorf("failed to save stage record: %w", result.Error)
	}
	return nil
}

func (r *StageRecordRepo) FindByID(id string) (*pipeline.StageRecord, error) {
	var model models.StageRecord
	result := r.db.First(&model, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find stage record by id: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *StageRecordRepo) FindAll() ([]*pipeline.StageRecord, error) {
	var models []models.StageRecord
	result := r.db.Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find all stage records: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *StageRecordRepo) Update(record *pipeline.StageRecord) error {
	model := &models.StageRecord{
		ID:         record.ID,
		PipelineID: record.PipelineID,
		StageName:  record.StageName,
		Status:     string(record.Status),
		StartedAt:  record.StartedAt,
		CompletedAt: record.CompletedAt,
		Output:     record.Output,
		Error:      record.Error,
		CreatedAt:  record.CreatedAt,
		UpdatedAt:  record.UpdatedAt,
	}

	result := r.db.Model(&models.StageRecord{}).Where("id = ?", record.ID).Updates(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update stage record: %w", result.Error)
	}
	return nil
}

func (r *StageRecordRepo) Delete(id string) error {
	result := r.db.Delete(&models.StageRecord{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete stage record: %w", result.Error)
	}
	return nil
}

func (r *StageRecordRepo) FindByPipelineID(pipelineID string) ([]*pipeline.StageRecord, error) {
	var models []models.StageRecord
	result := r.db.Where("pipeline_id = ?", pipelineID).Order("created_at ASC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find stage records by pipeline id: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *StageRecordRepo) FindByPipelineIDAndStageName(pipelineID, stageName string) (*pipeline.StageRecord, error) {
	var model models.StageRecord
	result := r.db.Where("pipeline_id = ? AND stage_name = ?", pipelineID, stageName).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find stage record by pipeline id and stage name: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *StageRecordRepo) FindByStatus(status string) ([]*pipeline.StageRecord, error) {
	var models []models.StageRecord
	result := r.db.Where("status = ?", status).Order("created_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find stage records by status: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *StageRecordRepo) modelToDomain(m *models.StageRecord) *pipeline.StageRecord {
	return &pipeline.StageRecord{
		ID:         m.ID,
		PipelineID: m.PipelineID,
		StageName:  m.StageName,
		Status:     pipeline.StageStatus(m.Status),
		StartedAt:  m.StartedAt,
		CompletedAt: m.CompletedAt,
		Output:     m.Output,
		Error:      m.Error,
	}
}

func (r *StageRecordRepo) modelsToDomain(models []models.StageRecord) []*pipeline.StageRecord {
	result := make([]*pipeline.StageRecord, len(models))
	for i, m := range models {
		result[i] = r.modelToDomain(&m)
	}
	return result
}
