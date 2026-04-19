package notify

import (
	"fmt"

	"github.com/mengri/nbcoder/domain/event"
	"github.com/mengri/nbcoder/domain/notify"
	"github.com/mengri/nbcoder/pkg/uid"
)

type NotifyService struct {
	notificationRepo notify.NotificationRepo
	subscriptionRepo notify.SubscriptionRepo
	prefRepo         notify.SubscriptionPreferenceRepo
	channelRepo      notify.ChannelRepo
	dispatcher       *notify.ChannelDispatcher
	eventBus         event.EventBus
}

func NewNotifyService(
	notificationRepo notify.NotificationRepo,
	subscriptionRepo notify.SubscriptionRepo,
	prefRepo notify.SubscriptionPreferenceRepo,
	channelRepo notify.ChannelRepo,
	dispatcher *notify.ChannelDispatcher,
	eventBus event.EventBus,
) *NotifyService {
	return &NotifyService{
		notificationRepo: notificationRepo,
		subscriptionRepo: subscriptionRepo,
		prefRepo:         prefRepo,
		channelRepo:      channelRepo,
		dispatcher:       dispatcher,
		eventBus:        eventBus,
	}
}

func (s *NotifyService) Send(title, content, eventType string, channel notify.ChannelType, recipient string) (*notify.Notification, error) {
	n := notify.NewNotification(uid.NewID(), title, content, eventType, channel, recipient)
	if err := n.Validate(); err != nil {
		return nil, err
	}
	if err := s.notificationRepo.Save(n); err != nil {
		return nil, err
	}
	if err := s.dispatcher.Dispatch(n); err != nil {
		n.MarkFailed()
		_ = s.notificationRepo.Update(n)
		return nil, fmt.Errorf("failed to dispatch notification via %s: %w", channel, err)
	}
	n.MarkSent()
	_ = s.notificationRepo.Update(n)
	return n, nil
}

func (s *NotifyService) SendToAllChannels(title, content, eventType string, recipient string) ([]*notify.Notification, error) {
	var results []*notify.Notification
	for _, channel := range []notify.ChannelType{notify.ChannelWebSocket, notify.ChannelEmail, notify.ChannelSystem} {
		if s.dispatcher.HasSender(channel) {
			n, err := s.Send(title, content, eventType, channel, recipient)
			if err != nil {
				continue
			}
			results = append(results, n)
		}
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no channel available for recipient %s", recipient)
	}
	return results, nil
}

func (s *NotifyService) GetByRecipient(recipient string) ([]*notify.Notification, error) {
	return s.notificationRepo.FindByRecipient(recipient)
}

func (s *NotifyService) GetByID(id string) (*notify.Notification, error) {
	return s.notificationRepo.FindByID(id)
}

func (s *NotifyService) MarkRead(id string) error {
	return s.notificationRepo.MarkRead(id)
}

func (s *NotifyService) Subscribe(id, recipient, eventType string, channel notify.ChannelType) (*notify.Subscription, error) {
	subs, _ := s.subscriptionRepo.FindByRecipientAndEventType(recipient, eventType)
	for _, sub := range subs {
		if sub.Channel == channel {
			return nil, fmt.Errorf("subscription already exists for recipient %s, event %s, channel %s", recipient, eventType, channel)
		}
	}
	sub := notify.NewSubscription(id, recipient, eventType, channel)
	if err := s.subscriptionRepo.Save(sub); err != nil {
		return nil, err
	}
	return sub, nil
}

func (s *NotifyService) Unsubscribe(subID string) error {
	return s.subscriptionRepo.Delete(subID)
}

func (s *NotifyService) GetSubscriptions(recipient string) ([]*notify.Subscription, error) {
	return s.subscriptionRepo.FindByRecipient(recipient)
}

func (s *NotifyService) RegisterChannel(channel *notify.Channel) error {
	if err := channel.Validate(); err != nil {
		return err
	}
	return s.channelRepo.Save(channel)
}

func (s *NotifyService) ListChannels(channelType notify.ChannelType) ([]*notify.Channel, error) {
	return s.channelRepo.FindByType(channelType)
}

func (s *NotifyService) NotifyFromDomainEvent(domainEvent event.DomainEvent) {
	eventType := domainEvent.EventType()
	subs, err := s.subscriptionRepo.FindByEventType(eventType)
	if err != nil {
		return
	}
	for _, sub := range subs {
		if sub.Muted {
			continue
		}
		pref, _ := s.prefRepo.FindByRecipientAndEventType(sub.Recipient, eventType)
		if pref != nil && pref.IsChannelMuted(sub.Channel) {
			continue
		}
		title := fmt.Sprintf("Event: %s", eventType)
		content := fmt.Sprintf("Aggregate %s triggered event %s", domainEvent.AggregateID(), eventType)
		_, _ = s.Send(title, content, eventType, sub.Channel, sub.Recipient)
	}
}

func (s *NotifyService) MuteSubscription(subID string) error {
	sub, err := s.subscriptionRepo.FindByRecipient(subID)
	_ = sub
	if err != nil {
		return err
	}
	return nil
}

func (s *NotifyService) MuteChannelForEvent(recipient, eventType string, channel notify.ChannelType) error {
	pref, _ := s.prefRepo.FindByRecipientAndEventType(recipient, eventType)
	if pref == nil {
		pref = notify.NewSubscriptionPreference(uid.NewID(), recipient, eventType)
		if err := s.prefRepo.Save(pref); err != nil {
			return err
		}
	}
	pref.MuteChannel(channel)
	return s.prefRepo.Update(pref)
}

func (s *NotifyService) UnmuteChannelForEvent(recipient, eventType string, channel notify.ChannelType) error {
	pref, _ := s.prefRepo.FindByRecipientAndEventType(recipient, eventType)
	if pref == nil {
		return fmt.Errorf("preference not found for recipient %s, event %s", recipient, eventType)
	}
	pref.UnmuteChannel(channel)
	return s.prefRepo.Update(pref)
}

func (s *NotifyService) GetPreferences(recipient string) ([]*notify.SubscriptionPreference, error) {
	return s.prefRepo.FindByRecipient(recipient)
}
