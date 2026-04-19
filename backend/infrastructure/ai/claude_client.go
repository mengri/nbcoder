package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ClaudeClient struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewClaudeClient(apiKey string) *ClaudeClient {
	return &ClaudeClient{
		apiKey:  apiKey,
		baseURL: "https://api.anthropic.com/v1",
		client:  &http.Client{},
	}
}

type ClaudeRequest struct {
	Model       string    `json:"model"`
	MaxTokens   int       `json:"max_tokens"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
}

type ClaudeResponse struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Usage struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

func (c *ClaudeClient) Complete(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	claudeReq := ClaudeRequest{
		Model:       req.Model,
		MaxTokens:   req.MaxTokens,
		Messages:    req.Messages,
		Temperature: req.Temperature,
	}

	body, err := json.Marshal(claudeReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/messages", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", c.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var claudeResp ClaudeResponse
	if err := json.NewDecoder(resp.Body).Decode(&claudeResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	result := &CompletionResponse{}
	if len(claudeResp.Content) > 0 {
		result.Choices = []struct {
			Message Message `json:"message"`
		}{
			{
				Message: Message{
					Role:    "assistant",
					Content: claudeResp.Content[0].Text,
				},
			},
		}
	}
	result.Usage.PromptTokens = claudeResp.Usage.InputTokens
	result.Usage.CompletionTokens = claudeResp.Usage.OutputTokens
	result.Usage.TotalTokens = claudeResp.Usage.InputTokens + claudeResp.Usage.OutputTokens

	return result, nil
}

func (c *ClaudeClient) GetProviderName() string {
	return "Claude"
}
