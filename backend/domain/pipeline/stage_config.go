package pipeline

type ReviewMode string

const (
	ReviewAI     ReviewMode = "AI"
	ReviewManual ReviewMode = "MANUAL"
	ReviewSkip   ReviewMode = "SKIP"
)

var DefaultStageNames = []string{
	"需求分析",
	"方案设计",
	"任务拆解",
	"测试用例",
	"代码开发",
	"测试验证",
	"评审合并",
}

type StageConfig struct {
	Enabled    bool       `json:"enabled"`
	ReviewMode ReviewMode `json:"review_mode"`
	MaxRetries int        `json:"max_retries"`
}

func DefaultStageConfig() StageConfig {
	return StageConfig{
		Enabled:    true,
		ReviewMode: ReviewAI,
		MaxRetries: 3,
	}
}
