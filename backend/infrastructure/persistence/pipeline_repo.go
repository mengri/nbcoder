package persistence

import (
	"sync"

	"github.com/mengri/nbcoder/domain/pipeline"
)

type InMemoryPipelineRepo struct {
	pipelines map[string]*pipeline.Pipeline
	mu        sync.RWMutex
}

func NewInMemoryPipelineRepo() *InMemoryPipelineRepo {
	return &InMemoryPipelineRepo{
		pipelines: make(map[string]*pipeline.Pipeline),
	}
}

func (r *InMemoryPipelineRepo) Save(p *pipeline.Pipeline) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.pipelines[p.ID] = p
	return nil
}

func (r *InMemoryPipelineRepo) FindByID(id string) (*pipeline.Pipeline, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.pipelines[id]
	if !ok {
		return nil, nil
	}
	return p, nil
}

func (r *InMemoryPipelineRepo) FindByCardID(cardID string) (*pipeline.Pipeline, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.pipelines {
		if p.CardID == cardID {
			return p, nil
		}
	}
	return nil, nil
}

func (r *InMemoryPipelineRepo) Update(p *pipeline.Pipeline) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.pipelines[p.ID] = p
	return nil
}
