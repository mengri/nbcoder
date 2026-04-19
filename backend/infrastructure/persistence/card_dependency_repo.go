package persistence

import (
	"sync"

	"github.com/mengri/nbcoder/domain/requirement"
)

type InMemoryCardDependencyRepo struct {
	deps map[string]*requirement.CardDependency
	mu   sync.RWMutex
}

func NewInMemoryCardDependencyRepo() *InMemoryCardDependencyRepo {
	return &InMemoryCardDependencyRepo{
		deps: make(map[string]*requirement.CardDependency),
	}
}

func (r *InMemoryCardDependencyRepo) Save(dep *requirement.CardDependency) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.deps[dep.ID] = dep
	return nil
}

func (r *InMemoryCardDependencyRepo) FindByCardID(cardID string) ([]*requirement.CardDependency, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*requirement.CardDependency
	for _, d := range r.deps {
		if d.CardID == cardID {
			result = append(result, d)
		}
	}
	return result, nil
}

func (r *InMemoryCardDependencyRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.deps, id)
	return nil
}
