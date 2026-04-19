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
	Tokens    int       `json:"tokens"`
	Timestamp time.Time `json:"timestamp"`
}

func CalculateCost(tokens int, rate float64) float64 {
	return float64(tokens) * rate
}
