package persistence

import (
	"sync"
	"time"

	"github.com/mengri/nbcoder/domain/pipeline"
)

type InMemoryStageRecordRepo struct {
	records map[string]*pipeline.StageRecord
	mu      sync.RWMutex
}

func NewInMemoryStageRecordRepo() *InMemoryStageRecordRepo {
	return &InMemoryStageRecordRepo{
		records: make(map[string]*pipeline.StageRecord),
	}
}

func (r *InMemoryStageRecordRepo) Save(record *pipeline.StageRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.records[record.ID] = record
	return nil
}

func (r *InMemoryStageRecordRepo) FindByStageID(stageID string) ([]*pipeline.StageRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*pipeline.StageRecord
	for _, rec := range r.records {
		if rec.StageID == stageID {
			result = append(result, rec)
		}
	}
	return result, nil
}

func (r *InMemoryStageRecordRepo) FindByTimeRange(start, end time.Time) ([]*pipeline.StageRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*pipeline.StageRecord
	for _, rec := range r.records {
		if rec.StartedAt.After(start) && rec.StartedAt.Before(end) {
			result = append(result, rec)
		}
	}
	return result, nil
}
