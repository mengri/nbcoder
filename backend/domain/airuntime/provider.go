package airuntime
// provider.go
// Provider与模型管理

type Provider struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Models  []*Model `json:"models"`
}

type Model struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	ProviderID string `json:"provider_id"`
	Meta       map[string]interface{} `json:"meta,omitempty"`
}

type ProviderRegistry struct {
	providers map[string]*Provider
}

func NewProviderRegistry() *ProviderRegistry {
	return &ProviderRegistry{
		providers: make(map[string]*Provider),
	}
}

func (r *ProviderRegistry) Register(p *Provider) {
	r.providers[p.ID] = p
}

func (r *ProviderRegistry) Get(id string) (*Provider, bool) {
	p, ok := r.providers[id]
	return p, ok
}

func (r *ProviderRegistry) SwitchModel(providerID, modelID string) (*Model, bool) {
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
