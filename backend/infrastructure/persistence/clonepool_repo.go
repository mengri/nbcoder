package persistence

import (
	"sync"

	"github.com/mengri/nbcoder/domain/clonepool"
)

type InMemoryCloneInstanceRepo struct {
	instances map[string]*clonepool.CloneInstance
	mu        sync.RWMutex
}

func NewInMemoryCloneInstanceRepo() *InMemoryCloneInstanceRepo {
	return &InMemoryCloneInstanceRepo{
		instances: make(map[string]*clonepool.CloneInstance),
	}
}

func (r *InMemoryCloneInstanceRepo) Save(instance *clonepool.CloneInstance) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.instances[instance.ID] = instance
	return nil
}

func (r *InMemoryCloneInstanceRepo) FindByID(id string) (*clonepool.CloneInstance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	inst, ok := r.instances[id]
	if !ok {
		return nil, nil
	}
	return inst, nil
}

func (r *InMemoryCloneInstanceRepo) FindByRepositoryID(repositoryID string) ([]*clonepool.CloneInstance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*clonepool.CloneInstance
	for _, inst := range r.instances {
		if inst.RepositoryID == repositoryID {
			result = append(result, inst)
		}
	}
	return result, nil
}

func (r *InMemoryCloneInstanceRepo) FindByStatus(status clonepool.CloneInstanceStatus) ([]*clonepool.CloneInstance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*clonepool.CloneInstance
	for _, inst := range r.instances {
		if inst.Status == status {
			result = append(result, inst)
		}
	}
	return result, nil
}

func (r *InMemoryCloneInstanceRepo) Update(instance *clonepool.CloneInstance) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.instances[instance.ID] = instance
	return nil
}

type InMemoryRepositoryRepo struct {
	repos map[string]*clonepool.Repository
	mu    sync.RWMutex
}

func NewInMemoryRepositoryRepo() *InMemoryRepositoryRepo {
	return &InMemoryRepositoryRepo{
		repos: make(map[string]*clonepool.Repository),
	}
}

func (r *InMemoryRepositoryRepo) Save(repo *clonepool.Repository) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.repos[repo.ID] = repo
	return nil
}

func (r *InMemoryRepositoryRepo) FindByID(id string) (*clonepool.Repository, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	repo, ok := r.repos[id]
	if !ok {
		return nil, nil
	}
	return repo, nil
}

func (r *InMemoryRepositoryRepo) FindAll() ([]*clonepool.Repository, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*clonepool.Repository, 0, len(r.repos))
	for _, repo := range r.repos {
		result = append(result, repo)
	}
	return result, nil
}
