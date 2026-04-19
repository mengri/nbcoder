package notify

import (
	"fmt"
	"time"
)

type ChannelType string

const (
	ChannelWebSocket ChannelType = "WEBSOCKET"
	ChannelEmail     ChannelType = "EMAIL"
	ChannelSystem    ChannelType = "SYSTEM"
)

func (ct ChannelType) IsValid() bool {
	switch ct {
	case ChannelWebSocket, ChannelEmail, ChannelSystem:
		return true
	}
	return false
}

type NotificationStatus string

const (
	NotificationPending   NotificationStatus = "PENDING"
	NotificationSent      NotificationStatus = "SENT"
	NotificationFailed    NotificationStatus = "FAILED"
	NotificationDelivered NotificationStatus = "DELIVERED"
)

func (ns NotificationStatus) IsValid() bool {
	switch ns {
	case NotificationPending, NotificationSent, NotificationFailed, NotificationDelivered:
		return true
	}
	return false
}

type Notification struct {
	ID        string             `json:"id"`
	Title     string             `json:"title"`
	Content   string             `json:"content"`
	EventType string             `json:"event_type"`
	Channel   ChannelType        `json:"channel"`
	Recipient string             `json:"recipient"`
	Status    NotificationStatus `json:"status"`
	Read      bool               `json:"read"`
	CreatedAt time.Time          `json:"created_at"`
	SentAt    *time.Time         `json:"sent_at,omitempty"`
}

func NewNotification(id, title, content, eventType string, channel ChannelType, recipient string) *Notification {
	return &Notification{
		ID:        id,
		Title:     title,
		Content:   content,
		EventType: eventType,
		Channel:   channel,
		Recipient: recipient,
		Status:    NotificationPending,
		Read:      false,
		CreatedAt: time.Now().UTC(),
	}
}

func (n *Notification) MarkSent() {
	now := time.Now().UTC()
	n.Status = NotificationSent
	n.SentAt = &now
}

func (n *Notification) MarkFailed() {
	n.Status = NotificationFailed
}

func (n *Notification) MarkDelivered() {
	n.Status = NotificationDelivered
}

func (n *Notification) MarkRead() {
	n.Read = true
}

func (n *Notification) Validate() error {
	if n.Title == "" {
		return fmt.Errorf("notification title is required")
	}
	if n.Recipient == "" {
		return fmt.Errorf("notification recipient is required")
	}
	if !n.Channel.IsValid() {
		return fmt.Errorf("invalid channel type: %s", n.Channel)
	}
	return nil
}
