package airuntime

import (
	"fmt"
	"sync"
)

type Provider struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	APIKeyRef  string   `json:"api_key_ref"`
	Models     []*Model `json:"models"`
}

type Model struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	ProviderID string `json:"provider_id"`
	ModelType  string `json:"model_type,omitempty"`
	Meta       map[string]interface{} `json:"meta,omitempty"`
}

type ProviderRegistry struct {
	providers map[string]*Provider
	mu        sync.RWMutex
}

func NewProviderRegistry() *ProviderRegistry {
	return &ProviderRegistry{
		providers: make(map[string]*Provider),
	}
}

func (r *ProviderRegistry) Register(p *Provider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[p.ID] = p
}

func (r *ProviderRegistry) Get(id string) (*Provider, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.providers[id]
	return p, ok
}

func (r *ProviderRegistry) SwitchModel(providerID, modelID string) (*Model, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.providers[providerID]
	if !ok {
		return nil, false
	}
	for _, m := range p.Models {
		if m.ID == modelID {
			return m, true
		}
	}
	return nil, false
}

func (r *ProviderRegistry) List() []*Provider {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*Provider, 0, len(r.providers))
	for _, p := range r.providers {
		result = append(result, p)
	}
	return result
}

func (r *ProviderRegistry) FindModelByID(modelID string) (*Model, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.providers {
		for _, m := range p.Models {
			if m.ID == modelID {
				return m, nil
			}
		}
	}
	return nil, fmt.Errorf("model %s not found", modelID)
}
