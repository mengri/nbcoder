package project

import (
	"github.com/mengri/nbcoder/domain/project"
)

type ProjectService struct {
	projectRepo project.ProjectRepo
	configRepo  project.ProjectConfigRepo
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
