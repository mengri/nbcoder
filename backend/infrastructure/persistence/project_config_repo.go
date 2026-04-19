package persistence

import (
	"sync"

	"github.com/mengri/nbcoder/domain/project"
)

type InMemoryProjectConfigRepo struct {
	configs map[string]*project.ProjectConfig
	mu      sync.RWMutex
}

func NewInMemoryProjectConfigRepo() *InMemoryProjectConfigRepo {
	return &InMemoryProjectConfigRepo{
		configs: make(map[string]*project.ProjectConfig),
	}
}

func (r *InMemoryProjectConfigRepo) Save(config *project.ProjectConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.configs[config.ID] = config
	return nil
}

func (r *InMemoryProjectConfigRepo) FindByProjectID(projectID string) ([]*project.ProjectConfig, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*project.ProjectConfig
	for _, c := range r.configs {
		if c.ProjectID == projectID {
			result = append(result, c)
		}
	}
	return result, nil
}

func (r *InMemoryProjectConfigRepo) Update(config *project.ProjectConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.configs[config.ID] = config
	return nil
}

func (r *InMemoryProjectConfigRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.configs, id)
	return nil
}

type InMemoryStandardsRepo struct {
	standards map[string]*project.Standards
	mu        sync.RWMutex
}

func NewInMemoryStandardsRepo() *InMemoryStandardsRepo {
	return &InMemoryStandardsRepo{
		standards: make(map[string]*project.Standards),
	}
}

func (r *InMemoryStandardsRepo) Save(standards *project.Standards) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.standards[standards.ID] = standards
	return nil
}

func (r *InMemoryStandardsRepo) FindByProjectID(projectID string) (*project.Standards, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, s := range r.standards {
		if s.ProjectID == projectID {
			return s, nil
		}
	}
	return nil, nil
}

func (r *InMemoryStandardsRepo) Update(standards *project.Standards) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.standards[standards.ID] = standards
	return nil
}
