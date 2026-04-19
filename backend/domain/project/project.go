package project

import (
	"fmt"
	"time"
)

type Project struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	RepoURL     string    `json:"repo_url,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewProject(id, name, description, repoURL string) *Project {
	now := time.Now().UTC()
	return &Project{
		ID:          id,
		Name:        name,
		Description: description,
		RepoURL:     repoURL,
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

func (p *Project) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("project name is required")
	}
	return nil
}
