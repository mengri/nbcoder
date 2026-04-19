package project

type ProjectRepo interface {
	Save(project *Project) error
	FindByID(id string) (*Project, error)
	FindAll() ([]*Project, error)
	Update(project *Project) error
	Delete(id string) error
}

type ProjectConfigRepo interface {
	Save(config *ProjectConfig) error
	FindByProjectID(projectID string) ([]*ProjectConfig, error)
	Update(config *ProjectConfig) error
	Delete(id string) error
}

type StandardsRepo interface {
	Save(standards *Standards) error
	FindByProjectID(projectID string) (*Standards, error)
	Update(standards *Standards) error
}
