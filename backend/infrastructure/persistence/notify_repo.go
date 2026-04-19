package persistence

import (
	"sync"

	"github.com/mengri/nbcoder/domain/notify"
)

type InMemoryNotificationRepo struct {
	notifications map[string]*notify.Notification
	mu            sync.RWMutex
}

func NewInMemoryNotificationRepo() *InMemoryNotificationRepo {
	return &InMemoryNotificationRepo{
		notifications: make(map[string]*notify.Notification),
	}
}

func (r *InMemoryNotificationRepo) Save(notification *notify.Notification) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.notifications[notification.ID] = notification
	return nil
}

func (r *InMemoryNotificationRepo) FindByID(id string) (*notify.Notification, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	n, ok := r.notifications[id]
	if !ok {
		return nil, nil
	}
	return n, nil
}

func (r *InMemoryNotificationRepo) FindByRecipient(recipient string) ([]*notify.Notification, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*notify.Notification
	for _, n := range r.notifications {
		if n.Recipient == recipient {
			result = append(result, n)
		}
	}
	return result, nil
}

func (r *InMemoryNotificationRepo) FindByEventType(eventType string) ([]*notify.Notification, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*notify.Notification
	for _, n := range r.notifications {
		if n.EventType == eventType {
			result = append(result, n)
		}
	}
	return result, nil
}

func (r *InMemoryNotificationRepo) MarkRead(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	n, ok := r.notifications[id]
	if !ok {
		return nil
	}
	n.MarkRead()
	return nil
}

type InMemorySubscriptionRepo struct {
	subs map[string]*notify.Subscription
	mu   sync.RWMutex
}

func NewInMemorySubscriptionRepo() *InMemorySubscriptionRepo {
	return &InMemorySubscriptionRepo{
		subs: make(map[string]*notify.Subscription),
	}
}

func (r *InMemorySubscriptionRepo) Save(sub *notify.Subscription) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.subs[sub.ID] = sub
	return nil
}

func (r *InMemorySubscriptionRepo) FindByRecipient(recipient string) ([]*notify.Subscription, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*notify.Subscription
	for _, s := range r.subs {
		if s.Recipient == recipient {
			result = append(result, s)
		}
	}
	return result, nil
}

func (r *InMemorySubscriptionRepo) FindByEventType(eventType string) ([]*notify.Subscription, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*notify.Subscription
	for _, s := range r.subs {
		if s.EventType == eventType {
			result = append(result, s)
		}
	}
	return result, nil
}

func (r *InMemorySubscriptionRepo) Update(sub *notify.Subscription) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.subs[sub.ID] = sub
	return nil
}

func (r *InMemorySubscriptionRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.subs, id)
	return nil
}
