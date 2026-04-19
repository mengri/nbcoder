package airuntime
// token_usage.go
// Token计费与调用日志
package airuntime

import "time"

type TokenUsage struct {
	ID        string    `json:"id"`
	ModelID   string    `json:"model_id"`
	AgentID   string    `json:"agent_id"`
	Tokens    int       `json:"tokens"`
	Cost      float64   `json:"cost"`
	Timestamp time.Time `json:"timestamp"`
}

type CallLog struct {
	ID        string    `json:"id"`
	ModelID   string    `json:"model_id"`
	AgentID   string    `json:"agent_id"`
	Input     string    `json:"input"`
	Output    string    `json:"output"`
	Timestamp time.Time `json:"timestamp"`
}

// 计费方式可扩展
func CalculateCost(tokens int, rate float64) float64 {
	return float64(tokens) * rate
}
