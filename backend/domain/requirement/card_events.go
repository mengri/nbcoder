package requirement

import (
	"github.com/mengri/nbcoder/domain/event"
	"github.com/mengri/nbcoder/pkg/uid"
)

type EventPublishMode string

const (
	EventPublishModeSync  EventPublishMode = "SYNC"
	EventPublishModeAsync EventPublishMode = "ASYNC"
)

func (ca *CardAggregate) CreateAndPublish(eventBus event.EventBus, mode EventPublishMode) error {
	evt := event.NewRequirementEvent(uid.NewID(), ca.Card.ID, event.CardCreatedEvent)
	evt.Payload["title"] = ca.Card.Title
	evt.Payload["description"] = ca.Card.Description
	evt.Payload["original"] = ca.Card.Original
	evt.Payload["priority"] = string(ca.Card.Priority)
	evt.Payload["project_name"] = ca.Card.ProjectName
	evt.Payload["pipeline_id"] = ca.Card.PipelineID
	evt.Payload["structured_output"] = ca.Card.StructuredOutput

	return ca.publishEvent(evt, eventBus, mode)
}

func (ca *CardAggregate) ConfirmAndPublish(eventBus event.EventBus, mode EventPublishMode) error {
	if err := ca.Confirm(); err != nil {
		return err
	}

	evt := event.NewRequirementEvent(uid.NewID(), ca.Card.ID, event.CardConfirmedEvent)
	evt.Payload["old_status"] = string(CardDraft)
	evt.Payload["new_status"] = string(ca.Card.Status)
	evt.Payload["project_name"] = ca.Card.ProjectName
	evt.Payload["title"] = ca.Card.Title

	return ca.publishEvent(evt, eventBus, mode)
}

func (ca *CardAggregate) SupersedeAndPublish(newCardID string, eventBus event.EventBus, mode EventPublishMode) error {
	if err := ca.Supersede(newCardID); err != nil {
		return err
	}

	evt := event.NewRequirementEvent(uid.NewID(), ca.Card.ID, event.CardSupersededEvent)
	evt.Payload["superseded_by"] = newCardID
	evt.Payload["old_status"] = string(ca.Card.Status)
	evt.Payload["project_name"] = ca.Card.ProjectName

	return ca.publishEvent(evt, eventBus, mode)
}

func (ca *CardAggregate) AbandonAndPublish(eventBus event.EventBus, mode EventPublishMode) error {
	if err := ca.Abandon(); err != nil {
		return err
	}

	evt := event.NewRequirementEvent(uid.NewID(), ca.Card.ID, event.CardAbandonedEvent)
	evt.Payload["old_status"] = string(ca.Card.Status)
	evt.Payload["project_name"] = ca.Card.ProjectName
	evt.Payload["had_dependencies"] = ca.HasDependencies()

	return ca.publishEvent(evt, eventBus, mode)
}

func (ca *CardAggregate) publishEvent(evt *event.RequirementEvent, eventBus event.EventBus, mode EventPublishMode) error {
	if mode == EventPublishModeAsync {
		go func() {
			_ = eventBus.Publish(evt)
		}()
		return nil
	}
	return eventBus.Publish(evt)
}

func (ca *CardAggregate) GetLineageInfo() map[string]interface{} {
	lineage := map[string]interface{}{
		"card_id":         ca.Card.ID,
		"title":           ca.Card.Title,
		"status":          string(ca.Card.Status),
		"superseded_by":   ca.Card.SupersededBy,
		"is_superseded":   ca.IsSuperseded(),
		"is_abandoned":    ca.IsAbandoned(),
		"has_dependencies": ca.HasDependencies(),
		"dependency_count": len(ca.Dependencies),
	}
	return lineage
}