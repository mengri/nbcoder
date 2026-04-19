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

type InMemoryDevStandardRepo struct {
	standards map[string]*project.DevStandard
	mu        sync.RWMutex
}

func NewInMemoryDevStandardRepo() *InMemoryDevStandardRepo {
	return &InMemoryDevStandardRepo{
		standards: make(map[string]*project.DevStandard),
	}
}

func (r *InMemoryDevStandardRepo) Save(standard *project.DevStandard) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.standards[standard.ID] = standard
	return nil
}

func (r *InMemoryDevStandardRepo) FindByProjectID(projectID string) ([]*project.DevStandard, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*project.DevStandard
	for _, s := range r.standards {
		if s.ProjectID == projectID {
			result = append(result, s)
		}
	}
	return result, nil
}

func (r *InMemoryDevStandardRepo) FindByID(id string) (*project.DevStandard, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.standards[id]
	if !ok {
		return nil, nil
	}
	return s, nil
}

func (r *InMemoryDevStandardRepo) Update(standard *project.DevStandard) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.standards[standard.ID] = standard
	return nil
}

func (r *InMemoryDevStandardRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.standards, id)
	return nil
}

type InMemoryBranchPolicyConfigRepo struct {
	configs map[string]*project.BranchPolicyConfig
	mu      sync.RWMutex
}

func NewInMemoryBranchPolicyConfigRepo() *InMemoryBranchPolicyConfigRepo {
	return &InMemoryBranchPolicyConfigRepo{
		configs: make(map[string]*project.BranchPolicyConfig),
	}
}

func (r *InMemoryBranchPolicyConfigRepo) Save(config *project.BranchPolicyConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.configs[config.ID] = config
	return nil
}

func (r *InMemoryBranchPolicyConfigRepo) FindByProjectID(projectID string) ([]*project.BranchPolicyConfig, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*project.BranchPolicyConfig
	for _, c := range r.configs {
		if c.ProjectID == projectID {
			result = append(result, c)
		}
	}
	return result, nil
}

func (r *InMemoryBranchPolicyConfigRepo) FindByID(id string) (*project.BranchPolicyConfig, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.configs[id]
	if !ok {
		return nil, nil
	}
	return c, nil
}

func (r *InMemoryBranchPolicyConfigRepo) Update(config *project.BranchPolicyConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.configs[config.ID] = config
	return nil
}

func (r *InMemoryBranchPolicyConfigRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.configs, id)
	return nil
}

type InMemoryProjectLifecycleRepo struct {
	lifecycles map[string]*project.ProjectLifecycle
	mu         sync.RWMutex
}

func NewInMemoryProjectLifecycleRepo() *InMemoryProjectLifecycleRepo {
	return &InMemoryProjectLifecycleRepo{
		lifecycles: make(map[string]*project.ProjectLifecycle),
	}
}

func (r *InMemoryProjectLifecycleRepo) Save(lifecycle *project.ProjectLifecycle) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lifecycles[lifecycle.ID] = lifecycle
	return nil
}

func (r *InMemoryProjectLifecycleRepo) FindByProjectID(projectID string) (*project.ProjectLifecycle, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, l := range r.lifecycles {
		if l.ProjectID == projectID {
			return l, nil
		}
	}
	return nil, nil
}

func (r *InMemoryProjectLifecycleRepo) Update(lifecycle *project.ProjectLifecycle) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lifecycles[lifecycle.ID] = lifecycle
	return nil
}
