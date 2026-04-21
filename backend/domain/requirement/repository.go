package requirement

type CardRepo interface {
	Save(card *Card) error
	FindByID(id string, projectName string) (*Card, error)
	FindByProjectName(projectName string) ([]*Card, error)
	FindByStatus(status CardStatus) ([]*Card, error)
	FindAll() ([]*Card, error)
	Update(card *Card) error
	Delete(id string, projectName string) error
}

type CardDependencyRepo interface {
	Save(dep *CardDependency) error
	FindByCardID(cardID string) ([]*CardDependency, error)
	FindByDependsOnID(dependsOnID string) ([]*CardDependency, error)
	FindByID(id string) (*CardDependency, error)
	Delete(id string) error
	DeleteByCardID(cardID string) error
}
