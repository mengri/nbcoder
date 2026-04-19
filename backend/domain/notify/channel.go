package notify

type Channel struct {
	ID     string      `json:"id"`
	Type   ChannelType `json:"type"`
	Config string      `json:"config,omitempty"`
}

type Subscription struct {
	ID        string `json:"id"`
	Recipient string `json:"recipient"`
	EventType string `json:"event_type"`
	Channel   ChannelType `json:"channel"`
	Muted     bool   `json:"muted"`
}

func NewSubscription(id, recipient, eventType string, channel ChannelType) *Subscription {
	return &Subscription{
		ID:        id,
		Recipient: recipient,
		EventType: eventType,
		Channel:   channel,
		Muted:     false,
	}
}

func (s *Subscription) Mute() {
	s.Muted = true
}

func (s *Subscription) Unmute() {
	s.Muted = false
}
