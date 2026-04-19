package requirement

type CardRepo interface {
	Save(card *Card) error
	FindByID(id string) (*Card, error)
	FindByProjectID(projectID string) ([]*Card, error)
	FindByStatus(status CardStatus) ([]*Card, error)
	Update(card *Card) error
}

type CardDependencyRepo interface {
	Save(dep *CardDependency) error
	FindByCardID(cardID string) ([]*CardDependency, error)
	Delete(id string) error
}
