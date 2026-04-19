package dto

type CreatePipelineRequest struct {
	CardID string `json:"card_id" binding:"required"`
}

type StageConfigDTO struct {
	Enabled    bool   `json:"enabled"`
	ReviewMode string `json:"review_mode"`
	MaxRetries int    `json:"max_retries"`
}

type PipelineResponse struct {
	ID        string            `json:"id"`
	CardID    string            `json:"card_id"`
	Stages    []StageResponse   `json:"stages"`
	CreatedAt string            `json:"created_at"`
}

type StageResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}
