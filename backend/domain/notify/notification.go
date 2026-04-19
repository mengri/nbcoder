package notify

import "time"

type ChannelType string

const (
	ChannelWebSocket ChannelType = "WEBSOCKET"
	ChannelEmail     ChannelType = "EMAIL"
	ChannelSystem    ChannelType = "SYSTEM"
)

type Notification struct {
	ID        string      `json:"id"`
	Title     string      `json:"title"`
	Content   string      `json:"content"`
	EventType string      `json:"event_type"`
	Channel   ChannelType `json:"channel"`
	Recipient string      `json:"recipient"`
	Read      bool        `json:"read"`
	CreatedAt time.Time   `json:"created_at"`
}

func NewNotification(id, title, content, eventType string, channel ChannelType, recipient string) *Notification {
	return &Notification{
		ID:        id,
		Title:     title,
		Content:   content,
		EventType: eventType,
		Channel:   channel,
		Recipient: recipient,
		Read:      false,
		CreatedAt: time.Now().UTC(),
	}
}

func (n *Notification) MarkRead() {
	n.Read = true
}
