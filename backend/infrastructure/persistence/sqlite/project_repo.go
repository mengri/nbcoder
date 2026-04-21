package sqlite

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mengri/nbcoder/domain/project"
)

type ProjectRepo struct {
	dbProvider     DBProvider
	projectBaseDir string
}

func NewProjectRepo(dbProvider DBProvider, projectBaseDir string) project.ProjectRepo {
	return &ProjectRepo{
		dbProvider:     dbProvider,
		projectBaseDir: projectBaseDir,
	}
}

type ProjectConfigFile struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	RepoURL     string          `json:"repo_url,omitempty"`
	Status      string          `json:"status"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
}

func (r *ProjectRepo) Save(p *project.Project) error {
	projectDir := filepath.Join(r.projectBaseDir, p.Name)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	configPath := filepath.Join(projectDir, ".nbcoder", "project.json")
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	config := ProjectConfigFile{
		Name:        p.Name,
		Description: p.Description,
		RepoURL:     p.RepoURL,
		Status:      string(p.Status),
		CreatedAt:   p.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   p.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z"),
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal project config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write project config: %w", err)
	}

	return nil
}

func (r *ProjectRepo) FindByName(name string) (*project.Project, error) {
	configPath := filepath.Join(r.projectBaseDir, name, ".nbcoder", "project.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read project config: %w", err)
	}

	var config ProjectConfigFile
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal project config: %w", err)
	}

	return r.configToDomain(&config), nil
}

func (r *ProjectRepo) FindAll() ([]*project.Project, error) {
	entries, err := os.ReadDir(r.projectBaseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read projects directory: %w", err)
	}

	var projects []*project.Project
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		p, err := r.FindByName(entry.Name())
		if err != nil {
			continue
		}
		if p != nil {
			projects = append(projects, p)
		}
	}

	return projects, nil
}

func (r *ProjectRepo) Update(p *project.Project) error {
	return r.Save(p)
}

func (r *ProjectRepo) Delete(name string) error {
	projectDir := filepath.Join(r.projectBaseDir, name)
	if err := os.RemoveAll(projectDir); err != nil {
		return fmt.Errorf("failed to delete project directory: %w", err)
	}
	return nil
}

func (r *ProjectRepo) FindByStatus(status project.ProjectStatus) ([]*project.Project, error) {
	projects, err := r.FindAll()
	if err != nil {
		return nil, err
	}

	var result []*project.Project
	for _, p := range projects {
		if p.Status == status {
			result = append(result, p)
		}
	}

	return result, nil
}

func (r *ProjectRepo) configToDomain(c *ProjectConfigFile) *project.Project {
	return &project.Project{
		Name:        c.Name,
		Description: c.Description,
		RepoURL:     c.RepoURL,
		Status:      project.ProjectStatus(c.Status),
		CreatedAt:   parseTime(c.CreatedAt),
		UpdatedAt:   parseTime(c.UpdatedAt),
	}
}

func parseTime(s string) time.Time {
	t, err := time.Parse("2006-01-02T15:04:05Z", s)
	if err != nil {
		return time.Time{}
	}
	return t
}
