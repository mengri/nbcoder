package project

import (
	"fmt"
	"regexp"
	"time"
)

type ProjectStatus string

const (
	ProjectActive   ProjectStatus = "ACTIVE"
	ProjectArchived ProjectStatus = "ARCHIVED"
)

type Project struct {
	Name        string        `json:"name"`
	Description string        `json:"description,omitempty"`
	RepoURL     string        `json:"repo_url,omitempty"`
	Status      ProjectStatus `json:"status"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

func NewProject(name, description, repoURL string) *Project {
	now := time.Now().UTC()
	return &Project{
		Name:        name,
		Description: description,
		RepoURL:     repoURL,
		Status:      ProjectActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (p *Project) Update(name, description, repoURL string) {
	if name != "" {
		p.Name = name
	}
	if description != "" {
		p.Description = description
	}
	if repoURL != "" {
		p.RepoURL = repoURL
	}
	p.UpdatedAt = time.Now().UTC()
}

func (p *Project) Archive() error {
	if p.Status == ProjectArchived {
		return fmt.Errorf("project already archived")
	}
	p.Status = ProjectArchived
	p.UpdatedAt = time.Now().UTC()
	return nil
}

func (p *Project) Activate() error {
	if p.Status == ProjectActive {
		return fmt.Errorf("project already active")
	}
	p.Status = ProjectActive
	p.UpdatedAt = time.Now().UTC()
	return nil
}

func (p *Project) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("project name is required")
	}
	
	matched, err := regexp.MatchString(`^[a-z0-9_]+$`, p.Name)
	if err != nil {
		return fmt.Errorf("invalid project name: %w", err)
	}
	if !matched {
		return fmt.Errorf("project name can only contain lowercase letters, numbers, and underscores")
	}
	
	return nil
}
