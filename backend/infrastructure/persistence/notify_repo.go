package persistence

import (
	"sync"
	"time"

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

type InMemoryNotificationTemplateRepo struct {
	templates map[string]*notify.NotificationTemplate
	mu        sync.RWMutex
}

func NewInMemoryNotificationTemplateRepo() *InMemoryNotificationTemplateRepo {
	return &InMemoryNotificationTemplateRepo{
		templates: make(map[string]*notify.NotificationTemplate),
	}
}

func (r *InMemoryNotificationTemplateRepo) Save(template *notify.NotificationTemplate) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.templates[template.ID] = template
	return nil
}

func (r *InMemoryNotificationTemplateRepo) FindByID(id string) (*notify.NotificationTemplate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.templates[id]
	if !ok {
		return nil, nil
	}
	return t, nil
}

func (r *InMemoryNotificationTemplateRepo) FindByEventType(eventType string) ([]*notify.NotificationTemplate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*notify.NotificationTemplate
	for _, t := range r.templates {
		if t.EventType == eventType {
			result = append(result, t)
		}
	}
	return result, nil
}

func (r *InMemoryNotificationTemplateRepo) Update(template *notify.NotificationTemplate) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.templates[template.ID] = template
	return nil
}

func (r *InMemoryNotificationTemplateRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.templates, id)
	return nil
}

type InMemoryNotificationHistoryRepo struct {
	histories map[string]*notify.NotificationHistory
	mu        sync.RWMutex
}

func NewInMemoryNotificationHistoryRepo() *InMemoryNotificationHistoryRepo {
	return &InMemoryNotificationHistoryRepo{
		histories: make(map[string]*notify.NotificationHistory),
	}
}

func (r *InMemoryNotificationHistoryRepo) Save(history *notify.NotificationHistory) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.histories[history.ID] = history
	return nil
}

func (r *InMemoryNotificationHistoryRepo) FindByNotificationID(notificationID string) ([]*notify.NotificationHistory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*notify.NotificationHistory
	for _, h := range r.histories {
		if h.NotificationID == notificationID {
			result = append(result, h)
		}
	}
	return result, nil
}

func (r *InMemoryNotificationHistoryRepo) FindByRecipient(recipient string) ([]*notify.NotificationHistory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*notify.NotificationHistory
	for _, h := range r.histories {
		if h.Recipient == recipient {
			result = append(result, h)
		}
	}
	return result, nil
}

func (r *InMemoryNotificationHistoryRepo) FindByTimeRange(start, end time.Time) ([]*notify.NotificationHistory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*notify.NotificationHistory
	for _, h := range r.histories {
		if !h.SentAt.Before(start) && !h.SentAt.After(end) {
			result = append(result, h)
		}
	}
	return result, nil
}
