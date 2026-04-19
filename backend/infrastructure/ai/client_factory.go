package ai

import (
	"fmt"
	"sync"
)

type ClientFactory struct {
	clients map[string]AIClient
	mu      sync.RWMutex
}

func NewClientFactory() *ClientFactory {
	return &ClientFactory{
		clients: make(map[string]AIClient),
	}
}

func (f *ClientFactory) CreateClient(providerType, apiKey string) (AIClient, error) {
	switch providerType {
	case "openai":
		return NewOpenAIClient(apiKey), nil
	case "anthropic":
		return NewClaudeClient(apiKey), nil
	default:
		return nil, fmt.Errorf("unsupported provider type: %s", providerType)
	}
}

func (f *ClientFactory) GetOrCreateClient(providerID, providerType, apiKey string) (AIClient, error) {
	f.mu.RLock()
	client, ok := f.clients[providerID]
	f.mu.RUnlock()

	if ok {
		return client, nil
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	if client, ok := f.clients[providerID]; ok {
		return client, nil
	}

	client, err := f.CreateClient(providerType, apiKey)
	if err != nil {
		return nil, err
	}

	f.clients[providerID] = client
	return client, nil
}

func (f *ClientFactory) RemoveClient(providerID string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.clients, providerID)
}
