package requirement

import (
	"github.com/mengri/nbcoder/domain/event"
	"github.com/mengri/nbcoder/domain/requirement"
)

type RequirementService struct {
	cardRepo     requirement.CardRepo
	depRepo      requirement.CardDependencyRepo
	eventBus     event.EventBus
}

func NewRequirementService(
	cardRepo requirement.CardRepo,
	depRepo requirement.CardDependencyRepo,
	eventBus event.EventBus,
) *RequirementService {
	return &RequirementService{
		cardRepo: cardRepo,
		depRepo:  depRepo,
		eventBus: eventBus,
	}
}

func (s *RequirementService) CreateCard(id, title, description, original, projectID string) (*requirement.Card, error) {
	card := requirement.NewCard(id, title, description, original, projectID)
	if err := s.cardRepo.Save(card); err != nil {
		return nil, err
	}
	evt := event.NewRequirementEvent(generateID(), id, event.CardCreatedEvent)
	_ = s.eventBus.Publish(evt)
	return card, nil
}

func (s *RequirementService) ConfirmCard(cardID string) error {
	card, err := s.cardRepo.FindByID(cardID)
	if err != nil {
		return err
	}
	deps, _ := s.depRepo.FindByCardID(cardID)
	aggregate := requirement.NewCardAggregate(card)
	for _, dep := range deps {
		aggregate.AddDependency(dep)
	}
	if err := aggregate.Confirm(); err != nil {
		return err
	}
	if err := s.cardRepo.Update(card); err != nil {
		return err
	}
	evt := event.NewRequirementEvent(generateID(), cardID, event.CardConfirmedEvent)
	return s.eventBus.Publish(evt)
}

func (s *RequirementService) GetCard(cardID string) (*requirement.Card, error) {
	return s.cardRepo.FindByID(cardID)
}

func generateID() string {
	return ""
}
