package pipeline

import (
	"fmt"
	"time"
)

type Stage struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Status    StageStatus `json:"status"`
	Config    StageConfig `json:"config"`
	StartedAt time.Time   `json:"started_at,omitempty"`
	EndedAt   time.Time   `json:"ended_at,omitempty"`
}

func NewStage(id, name string, config StageConfig) *Stage {
	return &Stage{
		ID:     id,
		Name:   name,
		Status: StageNotStarted,
		Config: config,
	}
}

func (s *Stage) Start() error {
	if s.Status != StageNotStarted {
		return fmt.Errorf("cannot start stage in status %s", s.Status)
	}
	s.Status = StageInProgress
	s.StartedAt = time.Now().UTC()
	return nil
}

func (s *Stage) Complete() error {
	if s.Status != StageInProgress {
		return fmt.Errorf("cannot complete stage in status %s", s.Status)
	}
	s.Status = StageCompleted
	s.EndedAt = time.Now().UTC()
	return nil
}

func (s *Stage) Fail(reason string) error {
	if s.Status != StageInProgress {
		return fmt.Errorf("cannot fail stage in status %s", s.Status)
	}
	s.Status = StageFailed
	s.EndedAt = time.Now().UTC()
	return nil
}

func (s *Stage) RequireReview() error {
	if s.Status != StageInProgress {
		return fmt.Errorf("cannot require review for stage in status %s", s.Status)
	}
	s.Status = StageReviewNeeded
	return nil
}

func (s *Stage) Retry() error {
	if s.Status != StageFailed {
		return fmt.Errorf("cannot retry stage in status %s", s.Status)
	}
	s.Status = StageNotStarted
	s.StartedAt = time.Time{}
	s.EndedAt = time.Time{}
	return nil
}
