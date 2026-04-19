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

type DevStandardRepo interface {
	Save(standard *DevStandard) error
	FindByProjectID(projectID string) ([]*DevStandard, error)
	FindByID(id string) (*DevStandard, error)
	Update(standard *DevStandard) error
	Delete(id string) error
}

type BranchPolicyConfigRepo interface {
	Save(config *BranchPolicyConfig) error
	FindByProjectID(projectID string) ([]*BranchPolicyConfig, error)
	FindByID(id string) (*BranchPolicyConfig, error)
	Update(config *BranchPolicyConfig) error
	Delete(id string) error
}

type ProjectLifecycleRepo interface {
	Save(lifecycle *ProjectLifecycle) error
	FindByProjectID(projectID string) (*ProjectLifecycle, error)
	Update(lifecycle *ProjectLifecycle) error
}
