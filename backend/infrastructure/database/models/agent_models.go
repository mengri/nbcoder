package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID           string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name         string         `gorm:"type:varchar(255);not null;index" json:"name"`
	Description  string         `gorm:"type:text" json:"description"`
	AgentType    string         `gorm:"type:varchar(50);not null;index" json:"agent_type"`
	TaskType     string         `gorm:"type:varchar(100);not null" json:"task_type"`
	Status       string         `gorm:"type:varchar(50);not null;default:'PENDING';index" json:"status"`
	Priority     int            `gorm:"not null;default:5;index" json:"priority"`
	AssignedTo   string         `gorm:"type:varchar(36);index" json:"assigned_to"`
	PipelineID   string         `gorm:"type:varchar(36);index" json:"pipeline_id"`
	ProjectName  string         `gorm:"type:varchar(255);not null;index:idx_task_project" json:"project_name"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	StartedAt    *time.Time     `json:"started_at"`
	CompletedAt  *time.Time     `json:"completed_at"`
	Context      JSONMap        `gorm:"type:json" json:"context,omitempty"`

	Pipeline        *Pipeline        `gorm:"foreignKey:PipelineID" json:"pipeline,omitempty"`
	AgentExecutions []AgentExecution `gorm:"foreignKey:TaskID" json:"agent_executions,omitempty"`
}

func (Task) TableName() string {
	return "tasks"
}

type AgentExecution struct {
	ID          string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	TaskID      string         `gorm:"type:varchar(36);not null;index:idx_agent_execution_task" json:"task_id"`
	AgentID     string         `gorm:"type:varchar(36);not null;index" json:"agent_id"`
	AgentType   string         `gorm:"type:varchar(50);not null" json:"agent_type"`
	SkillName   string         `gorm:"type:varchar(255);not null" json:"skill_name"`
	Status      string         `gorm:"type:varchar(50);not null;default:'PENDING';index" json:"status"`
	Input       string         `gorm:"type:text" json:"input"`
	Output      string         `gorm:"type:text" json:"output"`
	Error       string         `gorm:"type:text" json:"error"`
	StartedAt   *time.Time     `json:"started_at"`
	CompletedAt *time.Time     `json:"completed_at"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	ExecCount   int            `gorm:"default:0" json:"exec_count"`

	Task *Task `gorm:"foreignKey:TaskID" json:"task,omitempty"`
}

func (AgentExecution) TableName() string {
	return "agent_executions"
}

type Skill struct {
	ID          string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null;uniqueIndex:idx_skill_name" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	AgentType   string         `gorm:"type:varchar(50);not null;index" json:"agent_type"`
	Config      JSONMap        `gorm:"type:json" json:"config,omitempty"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
