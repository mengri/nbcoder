package pipeline

import (
	"fmt"
	"time"
)

type Pipeline struct {
	ID        string         `json:"id"`
	CardID    string         `json:"card_id"`
	Stages    []*Stage       `json:"stages"`
	Records   []*StageRecord `json:"records,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func NewPipeline(id, cardID string) *Pipeline {
	now := time.Now().UTC()
	return &Pipeline{
		ID:        id,
		CardID:    cardID,
		Stages:    []*Stage{},
		Records:   []*StageRecord{},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func NewDefaultPipeline(id, cardID string) *Pipeline {
	p := NewPipeline(id, cardID)
	for i, name := range DefaultStageNames {
		stageID := fmt.Sprintf("%s-stage-%d", id, i+1)
		stage := NewStage(stageID, name, DefaultStageConfig())
		p.Stages = append(p.Stages, stage)
	}
	p.UpdatedAt = time.Now().UTC()
	return p
}

func (p *Pipeline) AddStage(stage *Stage) {
	p.Stages = append(p.Stages, stage)
	p.UpdatedAt = time.Now().UTC()
}

func (p *Pipeline) CurrentStage() *Stage {
	for _, s := range p.Stages {
		if s.Status == StageInProgress || s.Status == StageReviewNeeded {
			return s
		}
	}
	return nil
}

func (p *Pipeline) NextPendingStage() *Stage {
	for _, s := range p.Stages {
		if s.Status == StageNotStarted {
			return s
		}
	}
	return nil
}

func (p *Pipeline) IsCompleted() bool {
	for _, s := range p.Stages {
		if s.Status != StageCompleted {
			return false
		}
	}
	return len(p.Stages) > 0
}

func (p *Pipeline) GetStageByName(name string) *Stage {
	for _, s := range p.Stages {
		if s.Name == name {
			return s
		}
	}
	return nil
}

func (p *Pipeline) GetStageRecords(stageID string) []*StageRecord {
	var result []*StageRecord
	for _, r := range p.Records {
		if r.StageID == stageID {
			result = append(result, r)
		}
	}
	return result
}
