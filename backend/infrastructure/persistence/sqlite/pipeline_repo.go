package sqlite

import (
	"fmt"

	"github.com/mengri/nbcoder/domain/pipeline"
	"github.com/mengri/nbcoder/infrastructure/database/models"
	"gorm.io/gorm"
)

type PipelineRepo struct {
	db *gorm.DB
}

func NewPipelineRepo(db *gorm.DB) pipeline.PipelineRepo {
	return &PipelineRepo{db: db}
}

func (r *PipelineRepo) Save(p *pipeline.Pipeline) error {
	model := &models.Pipeline{
		ID:        p.ID,
		CardID:    p.CardID,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}

	result := r.db.Save(model)
	if result.Error != nil {
		return fmt.Errorf("failed to save pipeline: %w", result.Error)
	}
	return nil
}

func (r *PipelineRepo) FindByID(id string) (*pipeline.Pipeline, error) {
	var model models.Pipeline
	result := r.db.Preload("StageRecords").First(&model, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find pipeline by id: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *PipelineRepo) FindByCardID(cardID string) (*pipeline.Pipeline, error) {
	var model models.Pipeline
	result := r.db.Preload("StageRecords").Where("card_id = ?", cardID).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find pipeline by card id: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *PipelineRepo) Update(p *pipeline.Pipeline) error {
	model := &models.Pipeline{
		ID:        p.ID,
		CardID:    p.CardID,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}

	result := r.db.Model(&models.Pipeline{}).Where("id = ?", p.ID).Updates(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update pipeline: %w", result.Error)
	}
	return nil
}

func (r *PipelineRepo) FindAll() ([]*pipeline.Pipeline, error) {
	var models []models.Pipeline
	result := r.db.Preload("StageRecords").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find all pipelines: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *PipelineRepo) modelToDomain(m *models.Pipeline) *pipeline.Pipeline {
	domainPipeline := &pipeline.Pipeline{
		ID:        m.ID,
		CardID:    m.CardID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}

	if len(m.StageRecords) > 0 {
		domainPipeline.Records = make([]*pipeline.StageRecord, len(m.StageRecords))
		for i, r := range m.StageRecords {
			domainPipeline.Records[i] = &pipeline.StageRecord{
				ID:          r.ID,
				PipelineID:  r.PipelineID,
				StageName:   r.StageName,
				Status:      pipeline.StageStatus(r.Status),
				StartedAt:   r.StartedAt,
				CompletedAt: r.CompletedAt,
				Output:      r.Output,
				Error:       r.Error,
			}
		}
	}

	return domainPipeline
}

func (r *PipelineRepo) modelsToDomain(models []models.Pipeline) []*pipeline.Pipeline {
	result := make([]*pipeline.Pipeline, len(models))
	for i, m := range models {
		result[i] = r.modelToDomain(&m)
	}
	return result
}
