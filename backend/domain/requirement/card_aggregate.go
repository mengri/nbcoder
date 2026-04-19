package requirement

import "fmt"

type CardAggregate struct {
	Card         *Card
	Dependencies []*CardDependency
}

func NewCardAggregate(card *Card) *CardAggregate {
	return &CardAggregate{
		Card:         card,
		Dependencies: []*CardDependency{},
	}
}

func (ca *CardAggregate) AddDependency(dep *CardDependency) {
	ca.Dependencies = append(ca.Dependencies, dep)
}

func (ca *CardAggregate) IsBlocked() bool {
	for _, dep := range ca.Dependencies {
		if dep.Type == DependencyDependsOn {
			return true
		}
	}
	return false
}

func (ca *CardAggregate) Confirm() error {
	if ca.IsBlocked() {
		return fmt.Errorf("card is blocked by dependencies")
	}
	return ca.Card.Confirm()
}
