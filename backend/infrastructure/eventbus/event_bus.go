package eventbus

import (
	"fmt"
	"sync"

	"github.com/mengri/nbcoder/domain/event"
)

type InMemoryEventBus struct {
	handlers map[string][]event.EventHandler
	mu       sync.RWMutex
}

func NewInMemoryEventBus() *InMemoryEventBus {
	return &InMemoryEventBus{
		handlers: make(map[string][]event.EventHandler),
	}
}

func (b *InMemoryEventBus) Publish(evt event.DomainEvent) error {
	b.mu.RLock()
	handlers, ok := b.handlers[evt.EventType()]
	b.mu.RUnlock()
	if !ok {
		return nil
	}
	for _, handler := range handlers {
		handler(evt)
	}
	return nil
}

func (b *InMemoryEventBus) Subscribe(eventType string, handler event.EventHandler) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventType] = append(b.handlers[eventType], handler)
	return nil
}

type InMemoryEventStore struct {
	events map[string][]event.DomainEvent
	mu     sync.RWMutex
}

func NewInMemoryEventStore() *InMemoryEventStore {
	return &InMemoryEventStore{
		events: make(map[string][]event.DomainEvent),
	}
}

func (s *InMemoryEventStore) Store(evt event.DomainEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events[evt.AggregateID()] = append(s.events[evt.AggregateID()], evt)
	return nil
}

func (s *InMemoryEventStore) Load(aggregateID string) ([]event.DomainEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	events, ok := s.events[aggregateID]
	if !ok {
		return nil, fmt.Errorf("no events found for aggregate %s", aggregateID)
	}
	return events, nil
}
