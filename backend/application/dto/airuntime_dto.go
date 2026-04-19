package dto

type CreateProviderRequest struct {
	Name      string `json:"name" binding:"required"`
	APIKeyRef string `json:"api_key_ref"`
}

type CreateModelRequest struct {
	Name       string `json:"name" binding:"required"`
	ProviderID string `json:"provider_id" binding:"required"`
	ModelType  string `json:"model_type"`
}

type CallModelRequest struct {
	ProviderID string `json:"provider_id" binding:"required"`
	ModelID    string `json:"model_id" binding:"required"`
	AgentID    string `json:"agent_id" binding:"required"`
	Messages   []struct {
		Role    string `json:"role" binding:"required"`
		Content string `json:"content" binding:"required"`
	} `json:"messages" binding:"required"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int    `json:"max_tokens"`
}

type ProviderResponse struct {
	ID     string            `json:"id"`
	Name   string            `json:"name"`
	Models []ModelResponse   `json:"models"`
}

type ModelResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	ProviderID string `json:"provider_id"`
	ModelType  string `json:"model_type"`
}
