package airuntime

import (
	"context"
	"fmt"
)

type ModelCaller struct {
	registry    *ProviderRegistry
	modelClient ModelClient
}

type ModelClient interface {
	Call(ctx context.Context, providerID, modelID string, messages []Message, temperature float64, maxTokens int) (*ModelResponse, error)
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ModelResponse struct {
	Content       string `json:"content"`
	PromptTokens  int    `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens   int    `json:"total_tokens"`
	Model         string `json:"model"`
}

func NewModelCaller(registry *ProviderRegistry, client ModelClient) *ModelCaller {
	return &ModelCaller{
		registry:    registry,
		modelClient: client,
	}
}

func (c *ModelCaller) CallModel(ctx context.Context, providerID, modelID string, messages []Message, opts ...CallOption) (*ModelResponse, error) {
	_, ok := c.registry.Get(providerID)
	if !ok {
		return nil, fmt.Errorf("provider %s not found", providerID)
	}

	options := &CallOptions{
		Temperature: 0.7,
		MaxTokens:   1000,
	}

	for _, opt := range opts {
		opt(options)
	}

	model, err := c.registry.FindModelByID(modelID)
	if err != nil {
		return nil, fmt.Errorf("model %s not found: %w", modelID, err)
	}

	if model.ProviderID != providerID {
		return nil, fmt.Errorf("model %s does not belong to provider %s", modelID, providerID)
	}

	response, err := c.modelClient.Call(ctx, providerID, modelID, messages, options.Temperature, options.MaxTokens)
	if err != nil {
		return nil, fmt.Errorf("failed to call model: %w", err)
	}

	response.Model = model.ID
	return response, nil
}

type CallOptions struct {
	Temperature float64
	MaxTokens   int
}

type CallOption func(*CallOptions)

func WithTemperature(temp float64) CallOption {
	return func(opts *CallOptions) {
		opts.Temperature = temp
	}
}

func WithMaxTokens(tokens int) CallOption {
	return func(opts *CallOptions) {
		opts.MaxTokens = tokens
	}
}
