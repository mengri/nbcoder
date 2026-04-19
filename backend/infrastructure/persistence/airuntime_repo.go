package persistence

import (
	"sync"
	"time"

	"github.com/mengri/nbcoder/domain/airuntime"
)

type InMemoryProviderRepo struct {
	providers map[string]*airuntime.Provider
	mu        sync.RWMutex
}

func NewInMemoryProviderRepo() *InMemoryProviderRepo {
	return &InMemoryProviderRepo{
		providers: make(map[string]*airuntime.Provider),
	}
}

func (r *InMemoryProviderRepo) Save(provider *airuntime.Provider) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[provider.ID] = provider
	return nil
}

func (r *InMemoryProviderRepo) FindByID(id string) (*airuntime.Provider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.providers[id]
	if !ok {
		return nil, nil
	}
	return p, nil
}

func (r *InMemoryProviderRepo) FindAll() ([]*airuntime.Provider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*airuntime.Provider, 0, len(r.providers))
	for _, p := range r.providers {
		result = append(result, p)
	}
	return result, nil
}

func (r *InMemoryProviderRepo) Update(provider *airuntime.Provider) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[provider.ID] = provider
	return nil
}

type InMemoryChainRepo struct {
	chains map[string]*airuntime.Chain
	mu     sync.RWMutex
}

func NewInMemoryChainRepo() *InMemoryChainRepo {
	return &InMemoryChainRepo{
		chains: make(map[string]*airuntime.Chain),
	}
}

func (r *InMemoryChainRepo) Save(chain *airuntime.Chain) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.chains[chain.ID] = chain
	return nil
}

func (r *InMemoryChainRepo) FindByID(id string) (*airuntime.Chain, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.chains[id]
	if !ok {
		return nil, nil
	}
	return c, nil
}

func (r *InMemoryChainRepo) FindAll() ([]*airuntime.Chain, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*airuntime.Chain, 0, len(r.chains))
	for _, c := range r.chains {
		result = append(result, c)
	}
	return result, nil
}

type InMemoryCallLogRepo struct {
	logs map[string]*airuntime.CallLog
	mu   sync.RWMutex
}

func NewInMemoryCallLogRepo() *InMemoryCallLogRepo {
	return &InMemoryCallLogRepo{
		logs: make(map[string]*airuntime.CallLog),
	}
}

func (r *InMemoryCallLogRepo) Save(log *airuntime.CallLog) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.logs[log.ID] = log
	return nil
}

func (r *InMemoryCallLogRepo) FindByAgentID(agentID string) ([]*airuntime.CallLog, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*airuntime.CallLog
	for _, l := range r.logs {
		if l.AgentID == agentID {
			result = append(result, l)
		}
	}
	return result, nil
}

func (r *InMemoryCallLogRepo) FindByTimeRange(start, end time.Time) ([]*airuntime.CallLog, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*airuntime.CallLog
	for _, l := range r.logs {
		if l.Timestamp.After(start) && l.Timestamp.Before(end) {
			result = append(result, l)
		}
	}
	return result, nil
}
