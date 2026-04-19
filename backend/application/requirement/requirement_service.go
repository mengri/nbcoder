package requirement

import (
	"fmt"

	"github.com/mengri/nbcoder/domain/event"
	"github.com/mengri/nbcoder/domain/requirement"
	"github.com/mengri/nbcoder/pkg/uid"
)

type RequirementService struct {
	cardRepo requirement.CardRepo
	depRepo  requirement.CardDependencyRepo
	eventBus event.EventBus
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

func (s *RequirementService) CreateCard(id, title, description, original, projectID string, priority requirement.Priority) (*requirement.Card, error) {
	card := requirement.NewCard(id, title, description, original, projectID)
	if priority != "" {
		if err := card.SetPriority(priority); err != nil {
			return nil, err
		}
	}
	if err := s.cardRepo.Save(card); err != nil {
		return nil, err
	}
	evt := event.NewRequirementEvent(uid.NewID(), id, event.CardCreatedEvent)
	_ = s.eventBus.Publish(evt)
	return card, nil
}

func (s *RequirementService) GetCard(cardID string) (*requirement.Card, error) {
	return s.cardRepo.FindByID(cardID)
}

func (s *RequirementService) ListCards(projectID string) ([]*requirement.Card, error) {
	if projectID != "" {
		return s.cardRepo.FindByProjectID(projectID)
	}
	return s.cardRepo.FindAll()
}

func (s *RequirementService) ListCardsByStatus(status requirement.CardStatus) ([]*requirement.Card, error) {
	return s.cardRepo.FindByStatus(status)
}

func (s *RequirementService) UpdateCard(cardID string, title, description *string, priority requirement.Priority, structuredOutput, pipelineID *string) (*requirement.Card, error) {
	card, err := s.cardRepo.FindByID(cardID)
	if err != nil {
		return nil, err
	}
	if card == nil {
		return nil, fmt.Errorf("card not found: %s", cardID)
	}
	if title != nil {
		card.Title = *title
		card.UpdatedAt = card.UpdatedAt.UTC()
	}
	if description != nil {
		card.Description = *description
		card.UpdatedAt = card.UpdatedAt.UTC()
	}
	if priority != "" {
		if err := card.SetPriority(priority); err != nil {
			return nil, err
		}
	}
	if structuredOutput != nil {
		card.SetStructuredOutput(*structuredOutput)
	}
	if pipelineID != nil {
		card.SetPipelineID(*pipelineID)
	}
	if err := s.cardRepo.Update(card); err != nil {
		return nil, err
	}
	return card, nil
}

func (s *RequirementService) DeleteCard(cardID string) error {
	card, err := s.cardRepo.FindByID(cardID)
	if err != nil {
		return err
	}
	if card == nil {
		return fmt.Errorf("card not found: %s", cardID)
	}
	_ = s.depRepo.DeleteByCardID(cardID)
	return s.cardRepo.Delete(cardID)
}

func (s *RequirementService) ConfirmCard(cardID string) error {
	card, err := s.cardRepo.FindByID(cardID)
	if err != nil {
		return err
	}
	if card == nil {
		return fmt.Errorf("card not found: %s", cardID)
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
	evt := event.NewRequirementEvent(uid.NewID(), cardID, event.CardConfirmedEvent)
	return s.eventBus.Publish(evt)
}

func (s *RequirementService) StartCard(cardID string) error {
	card, err := s.cardRepo.FindByID(cardID)
	if err != nil {
		return err
	}
	if card == nil {
		return fmt.Errorf("card not found: %s", cardID)
	}
	if err := card.Start(); err != nil {
		return err
	}
	return s.cardRepo.Update(card)
}

func (s *RequirementService) CompleteCard(cardID string) error {
	card, err := s.cardRepo.FindByID(cardID)
	if err != nil {
		return err
	}
	if card == nil {
		return fmt.Errorf("card not found: %s", cardID)
	}
	if err := card.Complete(); err != nil {
		return err
	}
	return s.cardRepo.Update(card)
}

func (s *RequirementService) AddDependency(id, cardID, dependsOnID string, depType requirement.DependencyType) (*requirement.CardDependency, error) {
	card, err := s.cardRepo.FindByID(cardID)
	if err != nil {
		return nil, err
	}
	if card == nil {
		return nil, fmt.Errorf("card not found: %s", cardID)
	}
	target, err := s.cardRepo.FindByID(dependsOnID)
	if err != nil {
		return nil, err
	}
	if target == nil {
		return nil, fmt.Errorf("depends-on card not found: %s", dependsOnID)
	}
	dep := requirement.NewCardDependency(id, cardID, dependsOnID, depType)
	if err := s.depRepo.Save(dep); err != nil {
		return nil, err
	}
	return dep, nil
}

func (s *RequirementService) RemoveDependency(depID string) error {
	return s.depRepo.Delete(depID)
}

func (s *RequirementService) GetDependencies(cardID string) ([]*requirement.CardDependency, error) {
	return s.depRepo.FindByCardID(cardID)
}
