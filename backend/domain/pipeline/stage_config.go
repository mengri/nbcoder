package pipeline

type ReviewMode string

const (
	ReviewAI     ReviewMode = "AI"
	ReviewManual ReviewMode = "MANUAL"
	ReviewSkip   ReviewMode = "SKIP"
)

type StageConfig struct {
	Enabled    bool       `json:"enabled"`
	ReviewMode ReviewMode `json:"review_mode"`
	MaxRetries int        `json:"max_retries"`
}
