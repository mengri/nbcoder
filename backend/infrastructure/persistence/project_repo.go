package persistence

import (
	"sync"

	"github.com/mengri/nbcoder/domain/project"
)

type InMemoryProjectRepo struct {
	projects map[string]*project.Project
	mu       sync.RWMutex
}

func NewInMemoryProjectRepo() *InMemoryProjectRepo {
	return &InMemoryProjectRepo{
		projects: make(map[string]*project.Project),
	}
}

func (r *InMemoryProjectRepo) Save(p *project.Project) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.projects[p.ID] = p
	return nil
}

func (r *InMemoryProjectRepo) FindByID(id string) (*project.Project, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.projects[id]
	if !ok {
		return nil, nil
	}
	return p, nil
}

func (r *InMemoryProjectRepo) FindAll() ([]*project.Project, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*project.Project, 0, len(r.projects))
	for _, p := range r.projects {
		result = append(result, p)
	}
	return result, nil
}

func (r *InMemoryProjectRepo) Update(p *project.Project) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.projects[p.ID] = p
	return nil
}

func (r *InMemoryProjectRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.projects, id)
	return nil
}
