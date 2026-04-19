package requirement

import (
	"github.com/mengri/nbcoder/domain/event"
	"github.com/mengri/nbcoder/pkg/uid"
)

func (ca *CardAggregate) SupersedeAndPublish(newCardID string, eventBus event.EventBus) error {
	if err := ca.Supersede(newCardID); err != nil {
		return err
	}

	evt := event.NewRequirementEvent(uid.NewID(), ca.Card.ID, event.CardSupersededEvent)
	evt.Payload["superseded_by"] = newCardID
	evt.Payload["old_status"] = string(ca.Card.Status)
	evt.Payload["project_id"] = ca.Card.ProjectID

	return eventBus.Publish(evt)
}

func (ca *CardAggregate) AbandonAndPublish(eventBus event.EventBus) error {
	if err := ca.Abandon(); err != nil {
		return err
	}

	evt := event.NewRequirementEvent(uid.NewID(), ca.Card.ID, event.CardAbandonedEvent)
	evt.Payload["old_status"] = string(ca.Card.Status)
	evt.Payload["project_id"] = ca.Card.ProjectID
	evt.Payload["had_dependencies"] = ca.HasDependencies()

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