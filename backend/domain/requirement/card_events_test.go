package requirement

import (
	"sync"
	"testing"

	"github.com/mengri/nbcoder/domain/event"
	"github.com/mengri/nbcoder/infrastructure/eventbus"
)

type TestEventHandler struct {
	events []event.DomainEvent
	mu     sync.Mutex
}

func NewTestEventHandler() *TestEventHandler {
	return &TestEventHandler{
		events: []event.DomainEvent{},
	}
}

func (h *TestEventHandler) Handle(evt event.DomainEvent) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.events = append(h.events, evt)
}

func (h *TestEventHandler) GetEvents() []event.DomainEvent {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.events
}

func (h *TestEventHandler) GetEventCount() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.events)
}

func (h *TestEventHandler) Clear() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.events = []event.DomainEvent{}
}

func TestCardAggregate_CreateAndPublish_Sync(t *testing.T) {
	card := NewCard("id-1", "Test Card", "desc", "orig", "proj-1")
	ca := NewCardAggregate(card)

	eventBus := eventbus.NewInMemoryEventBus()
	handler := NewTestEventHandler()

	if err := eventBus.Subscribe(string(event.CardCreatedEvent), handler.Handle); err != nil {
		t.Fatalf("failed to subscribe to event: %v", err)
	}

	if err := ca.CreateAndPublish(eventBus, EventPublishModeSync); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if handler.GetEventCount() != 1 {
		t.Errorf("expected 1 event, got %d", handler.GetEventCount())
	}

	evt := handler.GetEvents()[0].(*event.RequirementEvent)
	if evt.Type != event.CardCreatedEvent {
		t.Errorf("expected event type CardCreatedEvent, got %s", evt.Type)
	}

	if evt.CardID != "id-1" {
		t.Errorf("expected card id id-1, got %s", evt.CardID)
	}

	if evt.Payload["title"] != "Test Card" {
		t.Errorf("expected title Test Card, got %v", evt.Payload["title"])
	}
}

func TestCardAggregate_CreateAndPublish_Async(t *testing.T) {
	card := NewCard("id-1", "Test Card", "desc", "orig", "proj-1")
	ca := NewCardAggregate(card)

	eventBus := eventbus.NewInMemoryEventBus()
	handler := NewTestEventHandler()

	if err := eventBus.Subscribe(string(event.CardCreatedEvent), handler.Handle); err != nil {
		t.Fatalf("failed to subscribe to event: %v", err)
	}

	if err := ca.CreateAndPublish(eventBus, EventPublishModeAsync); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if handler.GetEventCount() != 0 {
		t.Error("expected 0 events immediately after async publish")
	}
}

func TestCardAggregate_ConfirmAndPublish_Sync(t *testing.T) {
	card := NewCard("id-1", "Test Card", "desc", "orig", "proj-1")
	ca := NewCardAggregate(card)

	eventBus := eventbus.NewInMemoryEventBus()
	handler := NewTestEventHandler()

	if err := eventBus.Subscribe(string(event.CardConfirmedEvent), handler.Handle); err != nil {
		t.Fatalf("failed to subscribe to event: %v", err)
	}

	if err := ca.ConfirmAndPublish(eventBus, EventPublishModeSync); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if card.Status != CardConfirmed {
		t.Errorf("expected status CardConfirmed, got %s", card.Status)
	}

	if handler.GetEventCount() != 1 {
		t.Errorf("expected 1 event, got %d", handler.GetEventCount())
	}

	evt := handler.GetEvents()[0].(*event.RequirementEvent)
	if evt.Type != event.CardConfirmedEvent {
		t.Errorf("expected event type CardConfirmedEvent, got %s", evt.Type)
	}

	if evt.Payload["new_status"] != string(CardConfirmed) {
		t.Errorf("expected new_status CONFIRMED, got %v", evt.Payload["new_status"])
	}
}

func TestCardAggregate_ConfirmAndPublish_Blocked(t *testing.T) {
	card := NewCard("id-1", "Test Card", "desc", "orig", "proj-1")
	ca := NewCardAggregate(card)

	dep := NewCardDependency("dep-1", "id-1", "id-2", DependencyDependsOn)
	ca.AddDependency(dep)

	eventBus := eventbus.NewInMemoryEventBus()
	handler := NewTestEventHandler()

	if err := eventBus.Subscribe(string(event.CardConfirmedEvent), handler.Handle); err != nil {
		t.Fatalf("failed to subscribe to event: %v", err)
	}

	if err := ca.ConfirmAndPublish(eventBus, EventPublishModeSync); err == nil {
		t.Error("expected error confirming blocked card")
	}

	if handler.GetEventCount() != 0 {
		t.Error("expected 0 events when confirmation fails")
	}
}

func TestCardAggregate_SupersedeAndPublish_Sync(t *testing.T) {
	card := NewCard("id-1", "Test Card", "desc", "orig", "proj-1")
	ca := NewCardAggregate(card)

	eventBus := eventbus.NewInMemoryEventBus()
	handler := NewTestEventHandler()

	if err := eventBus.Subscribe(string(event.CardSupersededEvent), handler.Handle); err != nil {
		t.Fatalf("failed to subscribe to event: %v", err)
	}

	if err := ca.SupersedeAndPublish("id-2", eventBus, EventPublishModeSync); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if card.Status != CardSuperseded {
		t.Errorf("expected status CardSuperseded, got %s", card.Status)
	}

	if handler.GetEventCount() != 1 {
		t.Errorf("expected 1 event, got %d", handler.GetEventCount())
	}

	evt := handler.GetEvents()[0].(*event.RequirementEvent)
	if evt.Type != event.CardSupersededEvent {
		t.Errorf("expected event type CardSupersededEvent, got %s", evt.Type)
	}

	if evt.Payload["superseded_by"] != "id-2" {
		t.Errorf("expected superseded_by id-2, got %v", evt.Payload["superseded_by"])
	}
}

func TestCardAggregate_AbandonAndPublish_Sync(t *testing.T) {
	card := NewCard("id-1", "Test Card", "desc", "orig", "proj-1")
	ca := NewCardAggregate(card)

	eventBus := eventbus.NewInMemoryEventBus()
	handler := NewTestEventHandler()

	if err := eventBus.Subscribe(string(event.CardAbandonedEvent), handler.Handle); err != nil {
		t.Fatalf("failed to subscribe to event: %v", err)
	}

	if err := ca.AbandonAndPublish(eventBus, EventPublishModeSync); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if card.Status != CardAbandoned {
		t.Errorf("expected status CardAbandoned, got %s", card.Status)
	}

	if handler.GetEventCount() != 1 {
		t.Errorf("expected 1 event, got %d", handler.GetEventCount())
	}

	evt := handler.GetEvents()[0].(*event.RequirementEvent)
	if evt.Type != event.CardAbandonedEvent {
		t.Errorf("expected event type CardAbandonedEvent, got %s", evt.Type)
	}

	if evt.Payload["project_id"] != "proj-1" {
		t.Errorf("expected project_id proj-1, got %v", evt.Payload["project_id"])
	}
}

func TestCardAggregate_MultipleEvents(t *testing.T) {
	card := NewCard("id-1", "Test Card", "desc", "orig", "proj-1")
	ca := NewCardAggregate(card)

	eventBus := eventbus.NewInMemoryEventBus()
	handler := NewTestEventHandler()

	if err := eventBus.Subscribe(string(event.CardCreatedEvent), handler.Handle); err != nil {
		t.Fatalf("failed to subscribe to event: %v", err)
	}
	if err := eventBus.Subscribe(string(event.CardConfirmedEvent), handler.Handle); err != nil {
		t.Fatalf("failed to subscribe to event: %v", err)
	}

	if err := ca.CreateAndPublish(eventBus, EventPublishModeSync); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := ca.ConfirmAndPublish(eventBus, EventPublishModeSync); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if handler.GetEventCount() != 2 {
		t.Errorf("expected 2 events, got %d", handler.GetEventCount())
	}

	if handler.GetEvents()[0].(*event.RequirementEvent).Type != event.CardCreatedEvent {
		t.Error("expected first event to be CardCreatedEvent")
	}

	if handler.GetEvents()[1].(*event.RequirementEvent).Type != event.CardConfirmedEvent {
		t.Error("expected second event to be CardConfirmedEvent")
	}
}

func TestCardAggregate_MultipleSubscribers(t *testing.T) {
	card := NewCard("id-1", "Test Card", "desc", "orig", "proj-1")
	ca := NewCardAggregate(card)

	eventBus := eventbus.NewInMemoryEventBus()
	handler1 := NewTestEventHandler()
	handler2 := NewTestEventHandler()

	if err := eventBus.Subscribe(string(event.CardCreatedEvent), handler1.Handle); err != nil {
		t.Fatalf("failed to subscribe to event: %v", err)
	}
	if err := eventBus.Subscribe(string(event.CardCreatedEvent), handler2.Handle); err != nil {
		t.Fatalf("failed to subscribe to event: %v", err)
	}

	if err := ca.CreateAndPublish(eventBus, EventPublishModeSync); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if handler1.GetEventCount() != 1 {
		t.Errorf("expected handler1 to receive 1 event, got %d", handler1.GetEventCount())
	}

	if handler2.GetEventCount() != 1 {
		t.Errorf("expected handler2 to receive 1 event, got %d", handler2.GetEventCount())
	}
}