package models

import (
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	ID          string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	ProjectName string         `gorm:"type:varchar(255);not null;uniqueIndex:idx_repo_project" json:"project_name"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	URL         string         `gorm:"type:varchar(500);not null" json:"url"`
	Branch      string         `gorm:"type:varchar(255);default:'main'" json:"branch"`
	LocalPath   string         `gorm:"type:varchar(500)" json:"local_path"`
	IsCloned    bool           `gorm:"default:false" json:"is_cloned"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	PullRequests []PullRequest  `gorm:"foreignKey:ProjectName" json:"pull_requests,omitempty"`
	Commits      []Commit       `gorm:"foreignKey:ProjectName" json:"commits,omitempty"`
}

func (Repository) TableName() string {
	return "repositories"
}

type PullRequest struct {
	ID             string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	Title          string         `gorm:"type:varchar(500);not null" json:"title"`
	Description    string         `gorm:"type:text" json:"description"`
	SourceBranch   string         `gorm:"type:varchar(255);not null;index" json:"source_branch"`
	TargetBranch   string         `gorm:"type:varchar(255);not null" json:"target_branch"`
	Status         string         `gorm:"type:varchar(50);not null;default:'OPEN';index" json:"status"`
	ProjectName    string         `gorm:"type:varchar(255);not null;index:idx_pr_project" json:"project_name"`
	Author         string         `gorm:"type:varchar(255)" json:"author"`
	GeneratedDesc  string         `gorm:"type:text" json:"generated_desc"`
	SquashCommitMsg string         `gorm:"type:text" json:"squash_commit_msg"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (PullRequest) TableName() string {
	return "pull_requests"
}

type Commit struct {
	ID          string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	ProjectName string         `gorm:"type:varchar(255);not null;index:idx_commit_project" json:"project_name"`
	Hash        string         `gorm:"type:varchar(100);not null;index" json:"hash"`
	Message     string         `gorm:"type:text;not null" json:"message"`
	Author      string         `gorm:"type:varchar(255);not null" json:"author"`
	Committer   string         `gorm:"type:varchar(255)" json:"committer"`
	CommitTime  time.Time      `gorm:"not null;index" json:"commit_time"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Commit) TableName() string {
	return "commits"
}
