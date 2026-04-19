package notify

import "time"

type HistoryStatus string

const (
	HistoryStatusSent     HistoryStatus = "SENT"
	HistoryStatusFailed   HistoryStatus = "FAILED"
	HistoryStatusDelivered HistoryStatus = "DELIVERED"
)

type NotificationHistory struct {
	ID             string        `json:"id"`
	NotificationID string        `json:"notification_id"`
	Recipient      string        `json:"recipient"`
	Channel        ChannelType   `json:"channel"`
	Status         HistoryStatus `json:"status"`
	SentAt         time.Time     `json:"sent_at"`
	Error          string        `json:"error,omitempty"`
}

func NewNotificationHistory(id, notificationID, recipient string, channel ChannelType, status HistoryStatus, sentAt time.Time, errMsg string) *NotificationHistory {
	return &NotificationHistory{
		ID:             id,
		NotificationID: notificationID,
		Recipient:      recipient,
		Channel:        channel,
		Status:         status,
		SentAt:         sentAt,
		Error:          errMsg,
	}
}
