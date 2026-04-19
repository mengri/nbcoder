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

func (r *InMemoryNotificationRepo) FindByRecipientAndChannel(recipient string, channel notify.ChannelType) ([]*notify.Notification, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*notify.Notification
	for _, n := range r.notifications {
		if n.Recipient == recipient && n.Channel == channel {
			result = append(result, n)
		}
	}
	return result, nil
}

func (r *InMemoryNotificationRepo) Update(notification *notify.Notification) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.notifications[notification.ID] = notification
	return nil
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

func (r *InMemorySubscriptionRepo) FindByRecipientAndEventType(recipient, eventType string) ([]*notify.Subscription, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*notify.Subscription
	for _, s := range r.subs {
		if s.Recipient == recipient && s.EventType == eventType {
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

type InMemoryChannelRepo struct {
	channels map[string]*notify.Channel
	mu       sync.RWMutex
}

func NewInMemoryChannelRepo() *InMemoryChannelRepo {
	return &InMemoryChannelRepo{
		channels: make(map[string]*notify.Channel),
	}
}

func (r *InMemoryChannelRepo) Save(channel *notify.Channel) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.channels[channel.ID] = channel
	return nil
}

func (r *InMemoryChannelRepo) FindByID(id string) (*notify.Channel, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ch, ok := r.channels[id]
	if !ok {
		return nil, nil
	}
	return ch, nil
}

func (r *InMemoryChannelRepo) FindByType(channelType notify.ChannelType) ([]*notify.Channel, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*notify.Channel
	for _, ch := range r.channels {
		if ch.Type == channelType {
			result = append(result, ch)
		}
	}
	return result, nil
}

func (r *InMemoryChannelRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.channels, id)
	return nil
}

type InMemorySubscriptionPreferenceRepo struct {
	prefs map[string]*notify.SubscriptionPreference
	mu    sync.RWMutex
}

func NewInMemorySubscriptionPreferenceRepo() *InMemorySubscriptionPreferenceRepo {
	return &InMemorySubscriptionPreferenceRepo{
		prefs: make(map[string]*notify.SubscriptionPreference),
	}
}

func (r *InMemorySubscriptionPreferenceRepo) Save(pref *notify.SubscriptionPreference) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.prefs[pref.ID] = pref
	return nil
}

func (r *InMemorySubscriptionPreferenceRepo) FindByRecipient(recipient string) ([]*notify.SubscriptionPreference, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*notify.SubscriptionPreference
	for _, p := range r.prefs {
		if p.Recipient == recipient {
			result = append(result, p)
		}
	}
	return result, nil
}

func (r *InMemorySubscriptionPreferenceRepo) FindByEventType(eventType string) ([]*notify.SubscriptionPreference, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*notify.SubscriptionPreference
	for _, p := range r.prefs {
		if p.EventType == eventType {
			result = append(result, p)
		}
	}
	return result, nil
}

func (r *InMemorySubscriptionPreferenceRepo) FindByRecipientAndEventType(recipient, eventType string) (*notify.SubscriptionPreference, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.prefs {
		if p.Recipient == recipient && p.EventType == eventType {
			return p, nil
		}
	}
	return nil, nil
}

func (r *InMemorySubscriptionPreferenceRepo) Update(pref *notify.SubscriptionPreference) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.prefs[pref.ID] = pref
	return nil
}

func (r *InMemorySubscriptionPreferenceRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.prefs, id)
	return nil
}
