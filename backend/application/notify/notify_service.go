package notify

import (
	"github.com/mengri/nbcoder/domain/event"
	"github.com/mengri/nbcoder/domain/notify"
)

type NotifyService struct {
	notificationRepo notify.NotificationRepo
	subscriptionRepo notify.SubscriptionRepo
	eventBus        event.EventBus
}

func NewNotifyService(
	notificationRepo notify.NotificationRepo,
	subscriptionRepo notify.SubscriptionRepo,
	eventBus event.EventBus,
) *NotifyService {
	return &NotifyService{
		notificationRepo: notificationRepo,
		subscriptionRepo: subscriptionRepo,
		eventBus:        eventBus,
	}
}

func (s *NotifyService) Send(notification *notify.Notification) error {
	return s.notificationRepo.Save(notification)
}

func (s *NotifyService) GetByRecipient(recipient string) ([]*notify.Notification, error) {
	return s.notificationRepo.FindByRecipient(recipient)
}

func (s *NotifyService) Subscribe(sub *notify.Subscription) error {
	return s.subscriptionRepo.Save(sub)
}
