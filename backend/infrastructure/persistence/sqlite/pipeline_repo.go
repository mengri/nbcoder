package sqlite

import (
	"fmt"
	"time"

	"github.com/mengri/nbcoder/domain/pipeline"
	"github.com/mengri/nbcoder/infrastructure/database/models"
	"gorm.io/gorm"
)

type PipelineRepo struct {
	dbProvider DBProvider
}

func NewPipelineRepo(dbProvider DBProvider) pipeline.PipelineRepo {
	return &PipelineRepo{dbProvider: dbProvider}
}

func (r *PipelineRepo) getDB() (*gorm.DB, error) {
	return r.dbProvider.GetGlobalDB(), nil
}

func (r *PipelineRepo) Save(p *pipeline.Pipeline) error {
	db, err := r.getDB()
	if err != nil {
		return err
	}

	model := &models.Pipeline{
		ID:        p.ID,
		CardID:    p.CardID,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}

	result := db.Save(model)
	if result.Error != nil {
		return fmt.Errorf("failed to save pipeline: %w", result.Error)
	}
	return nil
}

func (r *PipelineRepo) FindByID(id string) (*pipeline.Pipeline, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var model models.Pipeline
	result := db.Preload("StageRecords").First(&model, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find pipeline by id: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *PipelineRepo) FindByCardID(cardID string) (*pipeline.Pipeline, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var model models.Pipeline
	result := db.Preload("StageRecords").Where("card_id = ?", cardID).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find pipeline by card id: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *PipelineRepo) Update(p *pipeline.Pipeline) error {
	db, err := r.getDB()
	if err != nil {
		return err
	}

	model := &models.Pipeline{
		ID:        p.ID,
		CardID:    p.CardID,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}

	result := db.Model(&models.Pipeline{}).Where("id = ?", p.ID).Updates(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update pipeline: %w", result.Error)
	}
	return nil
}

func (r *PipelineRepo) FindAll() ([]*pipeline.Pipeline, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var models []models.Pipeline
	result := db.Preload("StageRecords").Find(&models)
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
			var startedAt time.Time
			if r.StartedAt != nil {
				startedAt = *r.StartedAt
			}

			var endedAt time.Time
			if r.CompletedAt != nil {
				endedAt = *r.CompletedAt
			}

			domainPipeline.Records[i] = &pipeline.StageRecord{
				ID:           r.ID,
				StageID:      "",
				Status:       pipeline.StageStatus(r.Status),
				StartedAt:    startedAt,
				EndedAt:      endedAt,
				Output:       r.Output,
				ReviewResult: "",
				Reviewer:     "",
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
