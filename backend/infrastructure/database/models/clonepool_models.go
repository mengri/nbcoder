package models

import (
	"time"

	"gorm.io/gorm"
)

type CloneInstance struct {
	ID         string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name       string         `gorm:"type:varchar(255);not null" json:"name"`
	ProjectID  string         `gorm:"type:varchar(36);not null;index:idx_clone_project" json:"project_id"`
	SourcePath string         `gorm:"type:varchar(500);not null" json:"source_path"`
	TargetPath string         `gorm:"type:varchar(500);not null" json:"target_path"`
	Status     string         `gorm:"type:varchar(50);not null;default:'IDLE';index" json:"status"`
	IsHealthy  bool           `gorm:"default:true" json:"is_healthy"`
	LastUsedAt *time.Time     `json:"last_used_at"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Project *Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
}

func (CloneInstance) TableName() string {
	return "clone_instances"
}
