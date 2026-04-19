package project

import "time"

type ProjectConfig struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewProjectConfig(id, projectID, key, value string) *ProjectConfig {
	now := time.Now().UTC()
	return &ProjectConfig{
		ID:        id,
		ProjectID: projectID,
		Key:       key,
		Value:     value,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

type Standards struct {
	ID                string    `json:"id"`
	ProjectID         string    `json:"project_id"`
	BranchStrategy    string    `json:"branch_strategy,omitempty"`
	TechStack         string    `json:"tech_stack,omitempty"`
	CodingConventions string    `json:"coding_conventions,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func NewStandards(id, projectID, branchStrategy, techStack, codingConventions string) *Standards {
	now := time.Now().UTC()
	return &Standards{
		ID:                id,
		ProjectID:         projectID,
		BranchStrategy:    branchStrategy,
		TechStack:         techStack,
		CodingConventions: codingConventions,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

func (s *Standards) Update(branchStrategy, techStack, codingConventions string) {
	if branchStrategy != "" {
		s.BranchStrategy = branchStrategy
	}
	if techStack != "" {
		s.TechStack = techStack
	}
	if codingConventions != "" {
		s.CodingConventions = codingConventions
	}
	s.UpdatedAt = time.Now().UTC()
}

type ProjectDirectory struct {
	ProjectID string
	Dirs      []string
}

func DefaultProjectDirectory(projectID string) *ProjectDirectory {
	return &ProjectDirectory{
		ProjectID: projectID,
		Dirs: []string{
			".NBCoder",
			".NBCoder/config",
			"knowledge-base",
			"knowledge-base/docs",
			"cards",
		},
	}
}
