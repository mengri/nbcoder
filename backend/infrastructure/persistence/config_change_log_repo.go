package persistence

import (
	"sync"

	"github.com/mengri/nbcoder/domain/project"
)

type InMemoryConfigChangeLogRepo struct {
	logs map[string]*project.ConfigChangeLog
	mu   sync.RWMutex
}

func NewInMemoryConfigChangeLogRepo() *InMemoryConfigChangeLogRepo {
	return &InMemoryConfigChangeLogRepo{
		logs: make(map[string]*project.ConfigChangeLog),
	}
}

func (r *InMemoryConfigChangeLogRepo) Save(log *project.ConfigChangeLog) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.logs[log.ID] = log
	return nil
}

func (r *InMemoryConfigChangeLogRepo) FindByProjectID(projectID string) ([]*project.ConfigChangeLog, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*project.ConfigChangeLog
	for _, l := range r.logs {
		if l.ProjectID == projectID {
			result = append(result, l)
		}
	}
	return result, nil
}

func (r *InMemoryConfigChangeLogRepo) FindByConfigKey(projectID, configKey string) ([]*project.ConfigChangeLog, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*project.ConfigChangeLog
	for _, l := range r.logs {
		if l.ProjectID == projectID && l.ConfigKey == configKey {
			result = append(result, l)
		}
	}
	return result, nil
}
