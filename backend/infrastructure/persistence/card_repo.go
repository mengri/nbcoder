package persistence

import (
	"sync"

	"github.com/mengri/nbcoder/domain/requirement"
)

type InMemoryCardRepo struct {
	cards map[string]*requirement.Card
	mu    sync.RWMutex
}

func NewInMemoryCardRepo() *InMemoryCardRepo {
	return &InMemoryCardRepo{
		cards: make(map[string]*requirement.Card),
	}
}

func (r *InMemoryCardRepo) Save(card *requirement.Card) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cards[card.ID] = card
	return nil
}

func (r *InMemoryCardRepo) FindByID(id string) (*requirement.Card, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	card, ok := r.cards[id]
	if !ok {
		return nil, nil
	}
	return card, nil
}

func (r *InMemoryCardRepo) FindByProjectID(projectID string) ([]*requirement.Card, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*requirement.Card
	for _, c := range r.cards {
		if c.ProjectName == projectID {
			result = append(result, c)
		}
	}
	return result, nil
}

func (r *InMemoryCardRepo) FindByStatus(status requirement.CardStatus) ([]*requirement.Card, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*requirement.Card
	for _, c := range r.cards {
		if c.Status == status {
			result = append(result, c)
		}
	}
	return result, nil
}

func (r *InMemoryCardRepo) FindAll() ([]*requirement.Card, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*requirement.Card, 0, len(r.cards))
	for _, c := range r.cards {
		result = append(result, c)
	}
	return result, nil
}

func (r *InMemoryCardRepo) Update(card *requirement.Card) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cards[card.ID] = card
	return nil
}

func (r *InMemoryCardRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.cards, id)
	return nil
}
