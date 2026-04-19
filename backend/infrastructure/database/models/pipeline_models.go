package models

import (
	"time"

	"gorm.io/gorm"
)

type Pipeline struct {
	ID        string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	CardID    string         `gorm:"type:varchar(36);not null;uniqueIndex:idx_pipeline_card" json:"card_id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Card         *Card         `gorm:"foreignKey:CardID" json:"card,omitempty"`
	StageRecords []StageRecord `gorm:"foreignKey:PipelineID" json:"stage_records,omitempty"`
}

func (Pipeline) TableName() string {
	return "pipelines"
}

type StageRecord struct {
	ID            string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	PipelineID    string         `gorm:"type:varchar(36);not null;index:idx_stage_record_pipeline" json:"pipeline_id"`
	StageName     string         `gorm:"type:varchar(255);not null" json:"stage_name"`
	Status        string         `gorm:"type:varchar(50);not null;default:'NOT_STARTED';index" json:"status"`
	StartedAt     *time.Time     `json:"started_at"`
	CompletedAt   *time.Time     `json:"completed_at"`
	Output        string         `gorm:"type:text" json:"output"`
	Error         string         `gorm:"type:text" json:"error"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Pipeline *Pipeline `gorm:"foreignKey:PipelineID" json:"pipeline,omitempty"`
}

func (StageRecord) TableName() string {
	return "stage_records"
}
