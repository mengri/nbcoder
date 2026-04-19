package notify

import (
	"fmt"
	"time"
)

type Channel struct {
	ID     string      `json:"id"`
	Type   ChannelType `json:"type"`
	Config string      `json:"config,omitempty"`
}

func NewChannel(id string, channelType ChannelType, config string) *Channel {
	return &Channel{
		ID:     id,
		Type:   channelType,
		Config: config,
	}
}

func (c *Channel) Validate() error {
	switch c.Type {
	case ChannelWebSocket, ChannelEmail, ChannelSystem:
		return nil
	}
	return fmt.Errorf("invalid channel type: %s", c.Type)
}

type ChannelSender interface {
	Send(notification *Notification) error
	Supports(channelType ChannelType) bool
}

type ChannelDispatcher struct {
	senders map[ChannelType]ChannelSender
}

func NewChannelDispatcher() *ChannelDispatcher {
	return &ChannelDispatcher{
		senders: make(map[ChannelType]ChannelSender),
	}
}

func (d *ChannelDispatcher) Register(sender ChannelSender) {
	for _, ct := range []ChannelType{ChannelWebSocket, ChannelEmail, ChannelSystem} {
		if sender.Supports(ct) {
			d.senders[ct] = sender
		}
	}
}

func (d *ChannelDispatcher) Dispatch(notification *Notification) error {
	sender, ok := d.senders[notification.Channel]
	if !ok {
		return fmt.Errorf("no sender registered for channel type: %s", notification.Channel)
	}
	return sender.Send(notification)
}

func (d *ChannelDispatcher) HasSender(channelType ChannelType) bool {
	_, ok := d.senders[channelType]
	return ok
}

type Subscription struct {
	ID        string      `json:"id"`
	Recipient string      `json:"recipient"`
	EventType string      `json:"event_type"`
	Channel   ChannelType `json:"channel"`
	Muted     bool        `json:"muted"`
	CreatedAt time.Time   `json:"created_at"`
}

func NewSubscription(id, recipient, eventType string, channel ChannelType) *Subscription {
	return &Subscription{
		ID:        id,
		Recipient: recipient,
		EventType: eventType,
		Channel:   channel,
		Muted:     false,
		CreatedAt: time.Now().UTC(),
	}
}

func (s *Subscription) Mute() {
	s.Muted = true
}

func (s *Subscription) Unmute() {
	s.Muted = false
}

func (s *Subscription) IsMutedForEvent(eventType string) bool {
	if s.Muted {
		return true
	}
	return false
}

type SubscriptionPreference struct {
	ID           string      `json:"id"`
	Recipient    string      `json:"recipient"`
	EventType    string      `json:"event_type"`
	MutedChannels []ChannelType `json:"muted_channels"`
}

func NewSubscriptionPreference(id, recipient, eventType string) *SubscriptionPreference {
	return &SubscriptionPreference{
		ID:            id,
		Recipient:     recipient,
		EventType:     eventType,
		MutedChannels: []ChannelType{},
	}
}

func (sp *SubscriptionPreference) MuteChannel(channel ChannelType) {
	for _, c := range sp.MutedChannels {
		if c == channel {
			return
		}
	}
	sp.MutedChannels = append(sp.MutedChannels, channel)
}

func (sp *SubscriptionPreference) UnmuteChannel(channel ChannelType) {
	for i, c := range sp.MutedChannels {
		if c == channel {
			sp.MutedChannels = append(sp.MutedChannels[:i], sp.MutedChannels[i+1:]...)
			return
		}
	}
}

func (sp *SubscriptionPreference) IsChannelMuted(channel ChannelType) bool {
	for _, c := range sp.MutedChannels {
		if c == channel {
			return true
		}
	}
	return false
}
