package models

import (
	"time"

	"gorm.io/gorm"
)

type Project struct {
	ID          string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null;index" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	RepoURL     string    `gorm:"type:varchar(500);index" json:"repo_url"`
	Status      string    `gorm:"type:varchar(50);not null;default:'ACTIVE';index" json:"status"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Configs        []ProjectConfig        `gorm:"foreignKey:ProjectID" json:"configs,omitempty"`
	Lifecycle      *ProjectLifecycle      `gorm:"foreignKey:ProjectID" json:"lifecycle,omitempty"`
	Standards      *Standards             `gorm:"foreignKey:ProjectID" json:"standards,omitempty"`
	ChangeLogs     []ConfigChangeLog      `gorm:"foreignKey:ProjectID" json:"change_logs,omitempty"`
	Cards          []Card                 `gorm:"foreignKey:ProjectID" json:"cards,omitempty"`
	Tasks          []Task                 `gorm:"foreignKey:ProjectID" json:"tasks,omitempty"`
	Documents      []Document             `gorm:"foreignKey:ProjectID" json:"documents,omitempty"`
	Directories    []Directory            `gorm:"foreignKey:ProjectID" json:"directories,omitempty"`
	PullRequests   []PullRequest          `gorm:"foreignKey:ProjectID" json:"pull_requests,omitempty"`
	CloneInstances []CloneInstance        `gorm:"foreignKey:ProjectID" json:"clone_instances,omitempty"`
}

func (Project) TableName() string {
	return "projects"
}

type ProjectConfig struct {
	ID        string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	ProjectID string    `gorm:"type:varchar(36);not null;index:idx_project_config_project" json:"project_id"`
	Key       string    `gorm:"type:varchar(255);not null;index:idx_project_config_key" json:"key"`
	Value     string    `gorm:"type:text" json:"value"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty"`

	Project *Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
}

func (ProjectConfig) TableName() string {
	return "project_configs"
}

type Standards struct {
	ID                string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	ProjectID         string    `gorm:"type:varchar(36);not null;uniqueIndex:idx_standards_project" json:"project_id"`
	BranchStrategy    string    `gorm:"type:text" json:"branch_strategy"`
	TechStack         string    `gorm:"type:text" json:"tech_stack"`
	CodingConventions string    `gorm:"type:text" json:"coding_conventions"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Project      *Project      `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	DevStandards []DevStandard `gorm:"foreignKey:ProjectID" json:"dev_standards,omitempty"`
}

func (Standards) TableName() string {
	return "standards"
}

type DevStandard struct {
	ID              string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	ProjectID       string    `gorm:"type:varchar(36);not null;index:idx_dev_standard_project" json:"project_id"`
	Name            string    `gorm:"type:varchar(255);not null" json:"name"`
	Description     string    `gorm:"type:text" json:"description"`
	StandardType    string    `gorm:"type:varchar(100);index" json:"standard_type"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Project   *Project   `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	Standards *Standards `gorm:"foreignKey:ProjectID" json:"standards,omitempty"`
}

func (DevStandard) TableName() string {
	return "dev_standards"
}

type BranchPolicyConfig struct {
	ID               string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	ProjectID        string    `gorm:"type:varchar(36);not null;index:idx_branch_policy_project" json:"project_id"`
	PolicyName       string    `gorm:"type:varchar(255);not null" json:"policy_name"`
	PolicyConfig     string    `gorm:"type:text" json:"policy_config"`
	RequireReviews   bool      `gorm:"default:false" json:"require_reviews"`
	MinReviewers     int       `gorm:"default:1" json:"min_reviewers"`
	RequireTests     bool      `gorm:"default:false" json:"require_tests"`
	AutoMergeEnabled bool      `gorm:"default:false" json:"auto_merge_enabled"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Project *Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
}

func (BranchPolicyConfig) TableName() string {
	return "branch_policy_configs"
}

type ProjectLifecycle struct {
	ID          string     `gorm:"type:varchar(36);primaryKey" json:"id"`
	ProjectID   string     `gorm:"type:varchar(36);not null;uniqueIndex:idx_lifecycle_project" json:"project_id"`
	Status      string     `gorm:"type:varchar(50);not null;default:'CREATING';index" json:"status"`
	ActivatedAt *time.Time `json:"activated_at"`
	SuspendedAt *time.Time `json:"suspended_at"`
	ArchivedAt  *time.Time `json:"archived_at"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Project *Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
}

func (ProjectLifecycle) TableName() string {
	return "project_lifecycles"
}

type ConfigChangeLog struct {
	ID         string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	ProjectID  string    `gorm:"type:varchar(36);not null;index:idx_change_log_project" json:"project_id"`
	ConfigKey  string    `gorm:"type:varchar(255);not null;index:idx_change_log_key" json:"config_key"`
	OldValue   string    `gorm:"type:text" json:"old_value"`
	NewValue   string    `gorm:"type:text" json:"new_value"`
	ChangedAt  time.Time `gorm:"autoCreateTime;index" json:"changed_at"`
	ChangedBy  string    `gorm:"type:varchar(255)" json:"changed_by"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Project *Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
}

func (ConfigChangeLog) TableName() string {
	return "config_change_logs"
}
