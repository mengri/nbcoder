package event

type EventHandler func(event DomainEvent)

type EventBus interface {
	Publish(event DomainEvent) error
	Subscribe(eventType string, handler EventHandler) error
}

type EventStore interface {
	Store(event DomainEvent) error
	Load(aggregateID string) ([]DomainEvent, error)
}
