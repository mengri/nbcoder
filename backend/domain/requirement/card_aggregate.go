package requirement

import (
	"fmt"

	"github.com/mengri/nbcoder/domain/event"
	"github.com/mengri/nbcoder/pkg/uid"
)

type DependencyChangeType string

const (
	DependencyChangeAdded     DependencyChangeType = "ADDED"
	DependencyChangeRemoved   DependencyChangeType = "REMOVED"
	DependencyChangeTypeChange DependencyChangeType = "TYPE_CHANGED"
)

type CardAggregate struct {
	Card         *Card
	Dependencies []*CardDependency
	ChangeType   DependencyChangeType
}

func NewCardAggregate(card *Card) *CardAggregate {
	return &CardAggregate{
		Card:         card,
		Dependencies: []*CardDependency{},
	}
}

func (ca *CardAggregate) AddDependency(dep *CardDependency) *CardDependency {
	ca.Dependencies = append(ca.Dependencies, dep)
	ca.ChangeType = DependencyChangeAdded
	return dep
}

func (ca *CardAggregate) RemoveDependency(depID string) error {
	for i, dep := range ca.Dependencies {
		if dep.ID == depID {
			ca.Dependencies = append(ca.Dependencies[:i], ca.Dependencies[i+1:]...)
			ca.ChangeType = DependencyChangeRemoved
			return nil
		}
	}
	return fmt.Errorf("dependency %s not found", depID)
}

func (ca *CardAggregate) IsBlocked() bool {
	return ca.IsDependencySatisfied()
}

func (ca *CardAggregate) IsDependencySatisfied() bool {
	for _, dep := range ca.Dependencies {
		if dep.Type == DependencyDependsOn {
			return true
		}
	}
	return false
}

func (ca *CardAggregate) GetBlockingDependencies() []*CardDependency {
	var blocking []*CardDependency
	for _, dep := range ca.Dependencies {
		if dep.Type == DependencyDependsOn {
			blocking = append(blocking, dep)
		}
	}
	return blocking
}

func (ca *CardAggregate) CheckAutoUnblock() bool {
	return !ca.IsDependencySatisfied()
}

func (ca *CardAggregate) Confirm() error {
	if ca.IsBlocked() {
		return fmt.Errorf("card is blocked by dependencies")
	}
	return ca.Card.Confirm()
}

func (ca *CardAggregate) GetDependencies() []*CardDependency {
	return ca.Dependencies
}

func (ca *CardAggregate) PublishDependencyChangeEvent(eventBus event.EventBus) error {
	evt := event.NewRequirementEvent(uid.NewID(), ca.Card.ID, CardDependenciesChangedEvent)
	evt.Payload["change_type"] = string(ca.ChangeType)
	evt.Payload["dependency_count"] = len(ca.Dependencies)
	return eventBus.Publish(evt)
}
