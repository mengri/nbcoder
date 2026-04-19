package models

import (
	"time"

	"gorm.io/gorm"
)

type Provider struct {
	ID        string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null;uniqueIndex:idx_provider_name" json:"name"`
	APIKeyRef string         `gorm:"type:varchar(500)" json:"api_key_ref"`
	BaseURL   string         `gorm:"type:varchar(500)" json:"base_url"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Models []Model `gorm:"foreignKey:ProviderID" json:"models,omitempty"`
}

func (Provider) TableName() string {
	return "providers"
}

type Model struct {
	ID         string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name       string         `gorm:"type:varchar(255);not null;index" json:"name"`
	ProviderID string         `gorm:"type:varchar(36);not null;index:idx_model_provider" json:"provider_id"`
	ModelType  string         `gorm:"type:varchar(100);index" json:"model_type"`
	Meta       JSONMap        `gorm:"type:json" json:"meta,omitempty"`
	IsActive   bool           `gorm:"default:true" json:"is_active"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Provider     *Provider     `gorm:"foreignKey:ProviderID" json:"provider,omitempty"`
	ModelChains  []ModelChain  `gorm:"foreignKey:ModelID" json:"model_chains,omitempty"`
	CallLogs     []CallLog     `gorm:"foreignKey:ModelID" json:"call_logs,omitempty"`
}

func (Model) TableName() string {
	return "models"
}

type ModelChain struct {
	ID         string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name       string         `gorm:"type:varchar(255);not null" json:"name"`
	ModelID    string         `gorm:"type:varchar(36);not null;index:idx_chain_model" json:"model_id"`
	ChainType  string         `gorm:"type:varchar(100);index" json:"chain_type"`
	Config     JSONMap        `gorm:"type:json" json:"config,omitempty"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Model *Model `gorm:"foreignKey:ModelID" json:"model,omitempty"`
}

func (ModelChain) TableName() string {
	return "model_chains"
}

type CallLog struct {
	ID          string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	AgentID     string         `gorm:"type:varchar(36);index" json:"agent_id"`
	ModelID     string         `gorm:"type:varchar(36);not null;index:idx_call_log_model" json:"model_id"`
	CallType    string         `gorm:"type:varchar(100);index" json:"call_type"`
	Input       string         `gorm:"type:text" json:"input"`
	Output      string         `gorm:"type:text" json:"output"`
	TokensUsed  int            `gorm:"default:0" json:"tokens_used"`
	Cost        float64        `gorm:"default:0" json:"cost"`
	LatencyMs   int            `gorm:"default:0" json:"latency_ms"`
	Status      string         `gorm:"type:varchar(50);not null;default:'SUCCESS';index" json:"status"`
	Error       string         `gorm:"type:text" json:"error"`
	CreatedAt   time.Time      `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Model *Model `gorm:"foreignKey:ModelID" json:"model,omitempty"`
}

func (CallLog) TableName() string {
	return "call_logs"
}
