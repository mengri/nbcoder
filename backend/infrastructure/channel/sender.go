package channel

import (
	"log"

	"github.com/mengri/nbcoder/domain/notify"
)

type SystemSender struct{}

func NewSystemSender() *SystemSender {
	return &SystemSender{}
}

func (s *SystemSender) Send(n *notify.Notification) error {
	log.Printf("[SYSTEM] To: %s | Title: %s | Content: %s", n.Recipient, n.Title, n.Content)
	n.MarkDelivered()
	return nil
}

func (s *SystemSender) Supports(channelType notify.ChannelType) bool {
	return channelType == notify.ChannelSystem
}

type WebSocketSender struct{}

func NewWebSocketSender() *WebSocketSender {
	return &WebSocketSender{}
}

func (s *WebSocketSender) Send(n *notify.Notification) error {
	log.Printf("[WEBSOCKET] To: %s | Title: %s | Content: %s", n.Recipient, n.Title, n.Content)
	n.MarkDelivered()
	return nil
}

func (s *WebSocketSender) Supports(channelType notify.ChannelType) bool {
	return channelType == notify.ChannelWebSocket
}

type EmailSender struct{}

func NewEmailSender() *EmailSender {
	return &EmailSender{}
}

func (s *EmailSender) Send(n *notify.Notification) error {
	log.Printf("[EMAIL] To: %s | Title: %s | Content: %s", n.Recipient, n.Title, n.Content)
	n.MarkDelivered()
	return nil
}

func (s *EmailSender) Supports(channelType notify.ChannelType) bool {
	return channelType == notify.ChannelEmail
}
