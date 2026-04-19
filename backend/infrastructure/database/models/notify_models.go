package models

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	ID         string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	Title      string         `gorm:"type:varchar(255);not null" json:"title"`
	Message    string         `gorm:"type:text;not null" json:"message"`
	EventType  string         `gorm:"type:varchar(100);not null;index" json:"event_type"`
	Recipient  string         `gorm:"type:varchar(255);not null;index" json:"recipient"`
	Channel    string         `gorm:"type:varchar(50);not null" json:"channel"`
	Status     string         `gorm:"type:varchar(50);not null;default:'UNREAD';index" json:"status"`
	Priority   string         `gorm:"type:varchar(50);default:'NORMAL'" json:"priority"`
	ReadAt     *time.Time     `json:"read_at"`
	CreatedAt  time.Time      `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Histories []NotificationHistory `gorm:"foreignKey:NotificationID" json:"histories,omitempty"`
}

func (Notification) TableName() string {
	return "notifications"
}

type Subscription struct {
	ID         string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	Recipient  string         `gorm:"type:varchar(255);not null;index:idx_sub_recipient" json:"recipient"`
	EventType  string         `gorm:"type:varchar(100);not null;index:idx_sub_event_type" json:"event_type"`
	Channel    string         `gorm:"type:varchar(50);not null" json:"channel"`
	IsActive   bool           `gorm:"default:true" json:"is_active"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Subscription) TableName() string {
	return "subscriptions"
}

type SubscriptionPreference struct {
	ID              string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	Recipient       string         `gorm:"type:varchar(255);not null;index:idx_pref_recipient" json:"recipient"`
	EventType       string         `gorm:"type:varchar(100);not null;index:idx_pref_event_type" json:"event_type"`
	EnabledChannels string         `gorm:"type:varchar(255)" json:"enabled_channels"`
	MinPriority     string         `gorm:"type:varchar(50);default:'LOW'" json:"min_priority"`
	DigestEnabled   bool           `gorm:"default:false" json:"digest_enabled"`
	DigestFrequency string         `gorm:"type:varchar(50)" json:"digest_frequency"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (SubscriptionPreference) TableName() string {
	return "subscription_preferences"
}

type NotificationTemplate struct {
	ID         string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name       string         `gorm:"type:varchar(255);not null;uniqueIndex:idx_template_name" json:"name"`
	EventType  string         `gorm:"type:varchar(100);not null;index" json:"event_type"`
	Subject    string         `gorm:"type:varchar(500)" json:"subject"`
	Body       string         `gorm:"type:text;not null" json:"body"`
	Channel    string         `gorm:"type:varchar(50);not null" json:"channel"`
	Variables  string         `gorm:"type:text" json:"variables"`
	IsActive   bool           `gorm:"default:true" json:"is_active"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (NotificationTemplate) TableName() string {
	return "notification_templates"
}

type NotificationHistory struct {
	ID             string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	NotificationID string         `gorm:"type:varchar(36);not null;index:idx_hist_notification" json:"notification_id"`
	Channel        string         `gorm:"type:varchar(50);not null" json:"channel"`
	Recipient      string         `gorm:"type:varchar(255);not null;index:idx_hist_recipient" json:"recipient"`
	Status         string         `gorm:"type:varchar(50);not null" json:"status"`
	SentAt         *time.Time     `json:"sent_at"`
	Error          string         `gorm:"type:text" json:"error"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Notification *Notification `gorm:"foreignKey:NotificationID" json:"notification,omitempty"`
}

func (NotificationHistory) TableName() string {
	return "notification_histories"
}
