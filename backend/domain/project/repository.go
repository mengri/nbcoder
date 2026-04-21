package project

type ProjectRepo interface {
	Save(project *Project) error
	FindByName(name string) (*Project, error)
	FindAll() ([]*Project, error)
	Update(project *Project) error
	Delete(name string) error
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
	FindByID(id string, projectName string) (*ProjectConfig, error)
	FindByProjectName(projectName string) ([]*ProjectConfig, error)
	FindByKey(projectName, key string) (*ProjectConfig, error)
	FindAll() ([]*ProjectConfig, error)
	Update(config *ProjectConfig) error
	Delete(id string, projectName string) error
}

type StandardsRepo interface {
	Save(standards *Standards) error
	FindByProjectName(projectName string) (*Standards, error)
	Update(standards *Standards) error
}

type DevStandardRepo interface {
	Save(standard *DevStandard) error
	FindByProjectName(projectName string) ([]*DevStandard, error)
	FindByID(id string, projectName string) (*DevStandard, error)
	Update(standard *DevStandard) error
	Delete(id string, projectName string) error
}

type BranchPolicyConfigRepo interface {
	Save(config *BranchPolicyConfig) error
	FindByProjectName(projectName string) ([]*BranchPolicyConfig, error)
	FindByID(id string, projectName string) (*BranchPolicyConfig, error)
	Update(config *BranchPolicyConfig) error
	Delete(id string, projectName string) error
}

type ProjectLifecycleRepo interface {
	Save(lifecycle *ProjectLifecycle) error
	FindByProjectName(projectName string) (*ProjectLifecycle, error)
	Update(lifecycle *ProjectLifecycle) error
}

type ConfigChangeLogRepo interface {
	Save(log *ConfigChangeLog) error
	FindByProjectName(projectName string) ([]*ConfigChangeLog, error)
	FindByConfigKey(projectName, configKey string) ([]*ConfigChangeLog, error)
}
