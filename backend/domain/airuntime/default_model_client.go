package airuntime

import (
	"context"
	"fmt"

	"github.com/mengri/nbcoder/infrastructure/ai"
)

type DefaultModelClient struct {
	clientFactory *ai.ClientFactory
	apiKeyGetter  func(providerID string) (string, error)
}

func NewDefaultModelClient(factory *ai.ClientFactory, apiKeyGetter func(providerID string) (string, error)) *DefaultModelClient {
	return &DefaultModelClient{
		clientFactory: factory,
		apiKeyGetter:  apiKeyGetter,
	}
}

func (c *DefaultModelClient) Call(ctx context.Context, providerID, modelID string, messages []Message, temperature float64, maxTokens int) (*ModelResponse, error) {
	apiKey, err := c.apiKeyGetter(providerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get API key: %w", err)
	}

	client, err := c.clientFactory.GetOrCreateClient(providerID, "openai", apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	aiMessages := make([]ai.Message, len(messages))
	for i, msg := range messages {
		aiMessages[i] = ai.Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	req := &ai.CompletionRequest{
		Model:       modelID,
		Messages:    aiMessages,
		Temperature: temperature,
		MaxTokens:   maxTokens,
	}

	resp, err := client.Complete(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to complete: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned from model")
	}

	return &ModelResponse{
		Content:          resp.Choices[0].Message.Content,
		PromptTokens:     resp.Usage.PromptTokens,
		CompletionTokens: resp.Usage.CompletionTokens,
		TotalTokens:      resp.Usage.TotalTokens,
	}, nil
}
