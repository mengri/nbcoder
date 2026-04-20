package project

type ProjectRepo interface {
	Save(project *Project) error
	FindByID(id string) (*Project, error)
	FindAll() ([]*Project, error)
	Update(project *Project) error
	Delete(id string) error
	FindByStatus(active ProjectStatus) ([]*Project, error)
}

type GlobalConfigRepo interface {
	Save(config *GlobalConfig) error
	FindByID(id string) (*GlobalConfig, error)
	FindByKey(key string) (*GlobalConfig, error)
	FindAll() ([]*GlobalConfig, error)
	Update(config *GlobalConfig) error
	Delete(id string) error
}

type ProjectConfigRepo interface {
	Save(config *ProjectConfig) error
	FindByID(id string) (*ProjectConfig, error)
	FindByProjectID(projectID string) ([]*ProjectConfig, error)
	FindByKey(projectID, key string) (*ProjectConfig, error)
	FindAll() ([]*ProjectConfig, error)
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

type ConfigChangeLogRepo interface {
	Save(log *ConfigChangeLog) error
	FindByProjectID(projectID string) ([]*ConfigChangeLog, error)
	FindByConfigKey(projectID, configKey string) ([]*ConfigChangeLog, error)
}
