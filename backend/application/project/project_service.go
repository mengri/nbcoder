package project

import (
	"fmt"

	"github.com/mengri/nbcoder/domain/project"
	"github.com/mengri/nbcoder/pkg/uid"
)

type ProjectService struct {
	projectRepo   project.ProjectRepo
	configRepo    project.ProjectConfigRepo
	standardsRepo project.StandardsRepo
}

func NewProjectService(
	projectRepo project.ProjectRepo,
	configRepo project.ProjectConfigRepo,
	standardsRepo project.StandardsRepo,
) *ProjectService {
	return &ProjectService{
		projectRepo:   projectRepo,
		configRepo:    configRepo,
		standardsRepo: standardsRepo,
	}
}

type InitProjectResult struct {
	Project     *project.Project
	Configs     []*project.ProjectConfig
	Standards   *project.Standards
	Directories *project.ProjectDirectory
}

func (s *ProjectService) InitProject(name, description, repoURL, branchStrategy, techStack, codingConventions string) (*InitProjectResult, error) {
	p := project.NewProject(uid.NewID(), name, description, repoURL)
	if err := p.Validate(); err != nil {
		return nil, err
	}
	if err := s.projectRepo.Save(p); err != nil {
		return nil, err
	}

	defaultConfigs := []struct{ Key, Value string }{
		{"initialized", "true"},
		{"branch_strategy", branchStrategy},
		{"tech_stack", techStack},
	}
	var configs []*project.ProjectConfig
	for _, c := range defaultConfigs {
		if c.Value != "" {
			cfg := project.NewProjectConfig(uid.NewID(), p.ID, c.Key, c.Value)
			_ = s.configRepo.Save(cfg)
			configs = append(configs, cfg)
		}
	}

	standards := project.NewStandards(uid.NewID(), p.ID, branchStrategy, techStack, codingConventions)
	_ = s.standardsRepo.Save(standards)

	dirs := project.DefaultProjectDirectory(p.ID)

	return &InitProjectResult{
		Project:     p,
		Configs:     configs,
		Standards:   standards,
		Directories: dirs,
	}, nil
}

func (s *ProjectService) CreateProject(id, name, description, repoURL string) (*project.Project, error) {
	p := project.NewProject(id, name, description, repoURL)
	if err := p.Validate(); err != nil {
		return nil, err
	}
	if err := s.projectRepo.Save(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *ProjectService) GetProject(id string) (*project.Project, error) {
	return s.projectRepo.FindByID(id)
}

func (s *ProjectService) ListProjects() ([]*project.Project, error) {
	return s.projectRepo.FindAll()
}

func (s *ProjectService) UpdateProject(id, name, description, repoURL string) (*project.Project, error) {
	p, err := s.projectRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, fmt.Errorf("project not found: %s", id)
	}
	p.Update(name, description, repoURL)
	if err := s.projectRepo.Update(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *ProjectService) DeleteProject(id string) error {
	p, err := s.projectRepo.FindByID(id)
	if err != nil {
		return err
	}
	if p == nil {
		return fmt.Errorf("project not found: %s", id)
	}
	return s.projectRepo.Delete(id)
}

func (s *ProjectService) ArchiveProject(id string) error {
	p, err := s.projectRepo.FindByID(id)
	if err != nil {
		return err
	}
	if p == nil {
		return fmt.Errorf("project not found: %s", id)
	}
	return p.Archive()
}

func (s *ProjectService) ActivateProject(id string) error {
	p, err := s.projectRepo.FindByID(id)
	if err != nil {
		return err
	}
	if p == nil {
		return fmt.Errorf("project not found: %s", id)
	}
	return p.Activate()
}

func (s *ProjectService) GetConfigs(projectID string) ([]*project.ProjectConfig, error) {
	return s.configRepo.FindByProjectID(projectID)
}

func (s *ProjectService) SetConfig(projectID, key, value string) (*project.ProjectConfig, error) {
	configs, _ := s.configRepo.FindByProjectID(projectID)
	for _, c := range configs {
		if c.Key == key {
			c.Value = value
			c.UpdatedAt = c.UpdatedAt.UTC()
			_ = s.configRepo.Update(c)
			return c, nil
		}
	}
	cfg := project.NewProjectConfig(uid.NewID(), projectID, key, value)
	if err := s.configRepo.Save(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (s *ProjectService) GetStandards(projectID string) (*project.Standards, error) {
	return s.standardsRepo.FindByProjectID(projectID)
}

func (s *ProjectService) UpdateStandards(projectID, branchStrategy, techStack, codingConventions string) (*project.Standards, error) {
	std, err := s.standardsRepo.FindByProjectID(projectID)
	if err != nil {
		return nil, err
	}
	if std == nil {
		std = project.NewStandards(uid.NewID(), projectID, branchStrategy, techStack, codingConventions)
		if err := s.standardsRepo.Save(std); err != nil {
			return nil, err
		}
		return std, nil
	}
	std.Update(branchStrategy, techStack, codingConventions)
	_ = s.standardsRepo.Update(std)
	return std, nil
}
