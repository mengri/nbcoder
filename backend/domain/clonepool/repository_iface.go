package clonepool

type CloneInstanceRepo interface {
	Save(instance *CloneInstance) error
	FindByID(id string) (*CloneInstance, error)
	FindByRepositoryID(repositoryID string) ([]*CloneInstance, error)
	FindByStatus(status CloneInstanceStatus) ([]*CloneInstance, error)
	Update(instance *CloneInstance) error
}

type RepositoryRepo interface {
	Save(repo *Repository) error
	FindByID(id string) (*Repository, error)
	FindAll() ([]*Repository, error)
}
