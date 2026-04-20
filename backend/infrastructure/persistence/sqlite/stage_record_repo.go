package sqlite

import (
	"fmt"
	"time"

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

// FindByTimeRange implements [pipeline.StageRecordRepo].
func (r *StageRecordRepo) FindByTimeRange(start time.Time, end time.Time) ([]*pipeline.StageRecord, error) {
	panic("unimplemented")
}
func (r *StageRecordRepo) Save(record *pipeline.StageRecord) error {
	var startedAt *time.Time
	if !record.StartedAt.IsZero() {
		startedAt = &record.StartedAt
	}

	var completedAt *time.Time
	if !record.EndedAt.IsZero() {
		completedAt = &record.EndedAt
	}

	createdAt := time.Now()
	if !record.StartedAt.IsZero() {
		createdAt = record.StartedAt
	}

	model := &models.StageRecord{
		ID:          record.ID,
		StageName:   record.StageID,
		Status:      string(record.Status),
		StartedAt:   startedAt,
		CompletedAt: completedAt,
		Output:      record.Output,
		CreatedAt:   createdAt,
		UpdatedAt:   time.Now(),
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
	var startedAt *time.Time
	if !record.StartedAt.IsZero() {
		startedAt = &record.StartedAt
	}

	var completedAt *time.Time
	if !record.EndedAt.IsZero() {
		completedAt = &record.EndedAt
	}

	model := &models.StageRecord{
		Status:      string(record.Status),
		StartedAt:   startedAt,
		CompletedAt: completedAt,
		Output:      record.Output,
		UpdatedAt:   time.Now(),
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

func (r *StageRecordRepo) FindByStageID(stageID string) ([]*pipeline.StageRecord, error) {
	return []*pipeline.StageRecord{}, nil
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
	var startedAt time.Time
	if m.StartedAt != nil {
		startedAt = *m.StartedAt
	}

	var endedAt time.Time
	if m.CompletedAt != nil {
		endedAt = *m.CompletedAt
	}

	return &pipeline.StageRecord{
		ID:        m.ID,
		StageID:   m.StageName,
		Status:    pipeline.StageStatus(m.Status),
		StartedAt: startedAt,
		EndedAt:   endedAt,
		Output:    m.Output,
	}
}

func (r *StageRecordRepo) modelsToDomain(models []models.StageRecord) []*pipeline.StageRecord {
	result := make([]*pipeline.StageRecord, len(models))
	for i, m := range models {
		result[i] = r.modelToDomain(&m)
	}
	return result
}
