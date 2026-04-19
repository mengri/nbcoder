package pipeline

import "time"

type StageRecord struct {
	ID        string      `json:"id"`
	StageID   string      `json:"stage_id"`
	Status    StageStatus `json:"status"`
	StartedAt time.Time   `json:"started_at"`
	EndedAt   time.Time   `json:"ended_at,omitempty"`
	Output    string      `json:"output,omitempty"`
	Reviewer  string      `json:"reviewer,omitempty"`
}

type StageRecordRepo interface {
	Save(record *StageRecord) error
	FindByStageID(stageID string) ([]*StageRecord, error)
	FindByTimeRange(start, end time.Time) ([]*StageRecord, error)
}
