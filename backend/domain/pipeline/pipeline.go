package pipeline

import (
	"time"
)

type Pipeline struct {
	ID        string   `json:"id"`
	CardID    string   `json:"card_id"`
	Stages    []*Stage `json:"stages"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewPipeline(id, cardID string) *Pipeline {
	now := time.Now().UTC()
	return &Pipeline{
		ID:        id,
		CardID:    cardID,
		Stages:    []*Stage{},
		CreatedAt: now,
		UpdatedAt: now,
	}
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
