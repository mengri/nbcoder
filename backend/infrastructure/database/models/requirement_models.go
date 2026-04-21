package models

import (
	"time"

	"gorm.io/gorm"
)

type Card struct {
	ID               string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	Title            string         `gorm:"type:varchar(500);not null;index" json:"title"`
	Description      string         `gorm:"type:text" json:"description"`
	Original         string         `gorm:"type:text" json:"original"`
	Status           string         `gorm:"type:varchar(50);not null;default:'DRAFT';index" json:"status"`
	Priority         string         `gorm:"type:varchar(50);not null;default:'MEDIUM';index" json:"priority"`
	StructuredOutput string         `gorm:"type:text" json:"structured_output"`
	PipelineID       string         `gorm:"type:varchar(36);index" json:"pipeline_id"`
	ProjectName      string         `gorm:"type:varchar(255);not null;index:idx_card_project" json:"project_name"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	SupersededBy     string         `gorm:"type:varchar(36);index" json:"superseded_by"`

	Pipeline    *Pipeline          `gorm:"foreignKey:PipelineID" json:"pipeline,omitempty"`
	Dependencies []CardDependency  `gorm:"foreignKey:CardID" json:"dependencies,omitempty"`
	Dependents  []CardDependency   `gorm:"foreignKey:DependsOnCardID" json:"dependents,omitempty"`
}

func (Card) TableName() string {
	return "cards"
}

type CardDependency struct {
	ID               string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	CardID           string         `gorm:"type:varchar(36);not null;index:idx_dependency_card" json:"card_id"`
	DependsOnCardID  string         `gorm:"type:varchar(36);not null;index:idx_dependency_depends_on" json:"depends_on_card_id"`
	DependencyType   string         `gorm:"type:varchar(50);not null;default:'BLOCKING'" json:"dependency_type"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Card          *Card `gorm:"foreignKey:CardID" json:"card,omitempty"`
	DependsOnCard *Card `gorm:"foreignKey:DependsOnCardID" json:"depends_on_card,omitempty"`
}

func (CardDependency) TableName() string {
	return "card_dependencies"
}
