package project

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mengri/nbcoder/domain/project"
	"github.com/mengri/nbcoder/infrastructure/database"
	"github.com/mengri/nbcoder/pkg/uid"
)

type ProjectService struct {
	projectRepo          project.ProjectRepo
	configRepo           project.ProjectConfigRepo
	standardsRepo        project.StandardsRepo
	configChangeLogRepo  project.ConfigChangeLogRepo
	projectBaseDir       string
	dbManager            database.DBManager
}

func NewProjectService(
	projectRepo project.ProjectRepo,
	configRepo project.ProjectConfigRepo,
	standardsRepo project.StandardsRepo,
	configChangeLogRepo project.ConfigChangeLogRepo,
	projectBaseDir string,
	dbManager database.DBManager,
) *ProjectService {
	return &ProjectService{
		projectRepo:         projectRepo,
		configRepo:          configRepo,
		standardsRepo:       standardsRepo,
		configChangeLogRepo: configChangeLogRepo,
		projectBaseDir:      projectBaseDir,
		dbManager:           dbManager,
	}
}

type InitProjectResult struct {
	Project     *project.Project
	Configs     []*project.ProjectConfig
	Standards   *project.Standards
	Directories *project.ProjectDirectory
}

func (s *ProjectService) InitProject(name, description, repoURL, branchStrategy, techStack, codingConventions string) (*InitProjectResult, error) {
	p := project.NewProject(name, description, repoURL)
	if err := p.Validate(); err != nil {
		return nil, err
	}
	
	projectDir := filepath.Join(s.projectBaseDir, name)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create project directory: %w", err)
	}
	
	if err := s.projectRepo.Save(p); err != nil {
		return nil, err
	}

	projectDB, err := s.dbManager.GetProjectDB(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get project database: %w", err)
	}

	if err := database.InitProjectSchema(projectDB); err != nil {
		return nil, fmt.Errorf("failed to init project schema: %w", err)
	}

	defaultConfigs := []struct{ Key, Value string }{
		{"initialized", "true"},
		{"branch_strategy", branchStrategy},
		{"tech_stack", techStack},
	}
	var configs []*project.ProjectConfig
	for _, c := range defaultConfigs {
		if c.Value != "" {
			cfg := project.NewProjectConfig(uid.NewID(), p.Name, c.Key, c.Value)
			_ = s.configRepo.Save(cfg)
			configs = append(configs, cfg)
		}
	}

	standards := project.NewStandards(uid.NewID(), p.Name, branchStrategy, techStack, codingConventions)
	_ = s.standardsRepo.Save(standards)

	dirs := project.DefaultProjectDirectory(p.Name)

	return &InitProjectResult{
		Project:     p,
		Configs:     configs,
		Standards:   standards,
		Directories: dirs,
	}, nil
}

func (s *ProjectService) CreateProject(name, description, repoURL string) (*project.Project, error) {
	p := project.NewProject(name, description, repoURL)
	if err := p.Validate(); err != nil {
		return nil, err
	}
	
	projectDir := filepath.Join(s.projectBaseDir, name)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create project directory: %w", err)
	}
	
	if err := s.projectRepo.Save(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *ProjectService) GetProject(name string) (*project.Project, error) {
	return s.projectRepo.FindByName(name)
}

func (s *ProjectService) ListProjects() ([]*project.Project, error) {
	return s.projectRepo.FindAll()
}

func (s *ProjectService) UpdateProject(name, newName, description, repoURL string) (*project.Project, error) {
	p, err := s.projectRepo.FindByName(name)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, fmt.Errorf("project not found: %s", name)
	}
	if newName != "" {
		p.Name = newName
	}
	if description != "" {
		p.Description = description
	}
	if repoURL != "" {
		p.RepoURL = repoURL
	}
	if err := s.projectRepo.Update(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *ProjectService) DeleteProject(name string) error {
	p, err := s.projectRepo.FindByName(name)
	if err != nil {
		return err
	}
	if p == nil {
		return fmt.Errorf("project not found: %s", name)
	}
	projectDir := filepath.Join(s.projectBaseDir, name)
	if err := os.RemoveAll(projectDir); err != nil {
		return fmt.Errorf("failed to delete project directory: %w", err)
	}
	return s.projectRepo.Delete(name)
}

func (s *ProjectService) ArchiveProject(name string) error {
	p, err := s.projectRepo.FindByName(name)
	if err != nil {
		return err
	}
	if p == nil {
		return fmt.Errorf("project not found: %s", name)
	}
	return p.Archive()
}

func (s *ProjectService) ActivateProject(name string) error {
	p, err := s.projectRepo.FindByName(name)
	if err != nil {
		return err
	}
	if p == nil {
		return fmt.Errorf("project not found: %s", name)
	}
	return p.Activate()
}

func (s *ProjectService) GetConfigs(projectName string) ([]*project.ProjectConfig, error) {
	return s.configRepo.FindByProjectName(projectName)
}

func (s *ProjectService) SetConfig(projectName, key, value string) (*project.ProjectConfig, error) {
	configs, _ := s.configRepo.FindByProjectName(projectName)
	for _, c := range configs {
		if c.Key == key {
			oldValue := c.Value
			c.Update(value)
			_ = s.configRepo.Update(c)
			changeLog := project.NewConfigChangeLog(uid.NewID(), projectName, key, oldValue, value, "")
			_ = s.configChangeLogRepo.Save(changeLog)
			return c, nil
		}
	}
	cfg := project.NewProjectConfig(uid.NewID(), projectName, key, value)
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if err := s.configRepo.Save(cfg); err != nil {
		return nil, err
	}
	changeLog := project.NewConfigChangeLog(uid.NewID(), projectName, key, "", value, "")
	_ = s.configChangeLogRepo.Save(changeLog)
	return cfg, nil
}

func (s *ProjectService) GetConfigHistory(projectName string) ([]*project.ConfigChangeLog, error) {
	return s.configChangeLogRepo.FindByProjectName(projectName)
}

func (s *ProjectService) GetConfig(projectName, key string) (*project.ProjectConfig, error) {
	configs, err := s.configRepo.FindByProjectName(projectName)
	if err != nil {
		return nil, err
	}
	for _, c := range configs {
		if c.Key == key {
			return c, nil
		}
	}
	return nil, fmt.Errorf("config not found: %s", key)
}

func (s *ProjectService) GetStandards(projectName string) (*project.Standards, error) {
	return s.standardsRepo.FindByProjectName(projectName)
}

func (s *ProjectService) UpdateStandards(projectName, branchStrategy, techStack, codingConventions string) (*project.Standards, error) {
	std, err := s.standardsRepo.FindByProjectName(projectName)
	if err != nil {
		return nil, err
	}
	if std == nil {
		std = project.NewStandards(uid.NewID(), projectName, branchStrategy, techStack, codingConventions)
		if err := s.standardsRepo.Save(std); err != nil {
			return nil, err
		}
		return std, nil
	}
	std.Update(branchStrategy, techStack, codingConventions)
	_ = s.standardsRepo.Update(std)
	return std, nil
}
