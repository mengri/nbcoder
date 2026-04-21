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

func (s *RequirementService) CreateCard(id, title, description, original, projectName string, priority requirement.Priority) (*requirement.Card, error) {
	card := requirement.NewCard(id, title, description, original, projectName)
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
	cards, _ := s.cardRepo.FindAll()
	projectName := ""
	if len(cards) > 0 {
		projectName = cards[0].ProjectName
	}
	return s.cardRepo.FindByID(cardID, projectName)
}

func (s *RequirementService) ListCards(projectName string) ([]*requirement.Card, error) {
	if projectName != "" {
		return s.cardRepo.FindByProjectName(projectName)
	}
	return s.cardRepo.FindAll()
}

func (s *RequirementService) ListCardsByStatus(status requirement.CardStatus) ([]*requirement.Card, error) {
	return s.cardRepo.FindByStatus(status)
}

func (s *RequirementService) UpdateCard(cardID string, title, description *string, priority requirement.Priority, structuredOutput, pipelineID *string) (*requirement.Card, error) {
	cards, _ := s.cardRepo.FindAll()
	projectName := ""
	if len(cards) > 0 {
		projectName = cards[0].ProjectName
	}
	card, err := s.cardRepo.FindByID(cardID, projectName)
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
	cards, _ := s.cardRepo.FindAll()
	projectName := ""
	if len(cards) > 0 {
		projectName = cards[0].ProjectName
	}
	card, err := s.cardRepo.FindByID(cardID, projectName)
	if err != nil {
		return err
	}
	if card == nil {
		return fmt.Errorf("card not found: %s", cardID)
	}
	_ = s.depRepo.DeleteByCardID(cardID)
	return s.cardRepo.Delete(cardID, projectName)
}

func (s *RequirementService) ConfirmCard(cardID string) error {
	cards, _ := s.cardRepo.FindAll()
	projectName := ""
	if len(cards) > 0 {
		projectName = cards[0].ProjectName
	}
	card, err := s.cardRepo.FindByID(cardID, projectName)
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
	if err := s.eventBus.Publish(evt); err != nil {
		return err
	}
	if aggregate.ChangeType == requirement.DependencyChangeRemoved {
		evt2 := event.NewRequirementEvent(uid.NewID(), cardID, event.CardDependenciesChangedEvent)
		evt2.Payload["change_type"] = string(aggregate.ChangeType)
		evt2.Payload["dependency_count"] = len(deps) - 1
		_ = s.eventBus.Publish(evt2)
	}
	return nil
}

func (s *RequirementService) StartCard(cardID string) error {
	cards, _ := s.cardRepo.FindAll()
	projectName := ""
	if len(cards) > 0 {
		projectName = cards[0].ProjectName
	}
	card, err := s.cardRepo.FindByID(cardID, projectName)
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
	cards, _ := s.cardRepo.FindAll()
	projectName := ""
	if len(cards) > 0 {
		projectName = cards[0].ProjectName
	}
	card, err := s.cardRepo.FindByID(cardID, projectName)
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
	cards, _ := s.cardRepo.FindAll()
	projectName := ""
	if len(cards) > 0 {
		projectName = cards[0].ProjectName
	}
	card, err := s.cardRepo.FindByID(cardID, projectName)
	if err != nil {
		return nil, err
	}
	if card == nil {
		return nil, fmt.Errorf("card not found: %s", cardID)
	}
	target, err := s.cardRepo.FindByID(dependsOnID, projectName)
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
	if err := s.depRepo.Delete(depID); err != nil {
		return err
	}
	return nil
}

func (s *RequirementService) GetDependencies(cardID string) ([]*requirement.CardDependency, error) {
	return s.depRepo.FindByCardID(cardID)
}

func (s *RequirementService) UpdateDependency(id, cardID, dependsOnID string, depType requirement.DependencyType) (*requirement.CardDependency, error) {
	cards, _ := s.cardRepo.FindAll()
	projectName := ""
	if len(cards) > 0 {
		projectName = cards[0].ProjectName
	}
	card, err := s.cardRepo.FindByID(cardID, projectName)
	if err != nil {
		return nil, err
	}
	if card == nil {
		return nil, fmt.Errorf("card not found: %s", cardID)
	}
	oldDeps, _ := s.depRepo.FindByCardID(cardID)
	aggregate := requirement.NewCardAggregate(card)
	for _, dep := range oldDeps {
		aggregate.AddDependency(dep)
	}

	dep := requirement.NewCardDependency(id, cardID, dependsOnID, depType)
	if err := s.depRepo.Save(dep); err != nil {
		return nil, err
	}

	if err := aggregate.PublishDependencyChangeEvent(s.eventBus); err != nil {
		return nil, err
	}

	return dep, nil
}
