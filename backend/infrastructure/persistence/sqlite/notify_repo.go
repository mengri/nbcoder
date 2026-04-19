package sqlite

import (
	"fmt"
	"time"

	"github.com/mengri/nbcoder/domain/notify"
	"github.com/mengri/nbcoder/infrastructure/database/models"
	"gorm.io/gorm"
)

type NotificationRepo struct {
	db *gorm.DB
}

func NewNotificationRepo(db *gorm.DB) notify.NotificationRepo {
	return &NotificationRepo{db: db}
}

func (r *NotificationRepo) Save(notification *notify.Notification) error {
	model := &models.Notification{
		ID:         notification.ID,
		Title:      notification.Title,
		Message:    notification.Message,
		EventType:  notification.EventType,
		Recipient:  notification.Recipient,
		Channel:    string(notification.Channel),
		Status:     string(notification.Status),
		Priority:   string(notification.Priority),
		ReadAt:     notification.ReadAt,
		CreatedAt:  notification.CreatedAt,
		UpdatedAt:  notification.UpdatedAt,
	}

	result := r.db.Save(model)
	if result.Error != nil {
		return fmt.Errorf("failed to save notification: %w", result.Error)
	}
	return nil
}

func (r *NotificationRepo) FindByID(id string) (*notify.Notification, error) {
	var model models.Notification
	result := r.db.First(&model, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find notification by id: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *NotificationRepo) FindByRecipient(recipient string) ([]*notify.Notification, error) {
	var models []models.Notification
	result := r.db.Where("recipient = ?", recipient).Order("created_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find notifications by recipient: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *NotificationRepo) FindByEventType(eventType string) ([]*notify.Notification, error) {
	var models []models.Notification
	result := r.db.Where("event_type = ?", eventType).Order("created_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find notifications by event type: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *NotificationRepo) FindByRecipientAndChannel(recipient string, channel notify.ChannelType) ([]*notify.Notification, error) {
	var models []models.Notification
	result := r.db.Where("recipient = ? AND channel = ?", recipient, string(channel)).Order("created_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find notifications by recipient and channel: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *NotificationRepo) Update(notification *notify.Notification) error {
	model := &models.Notification{
		ID:         notification.ID,
		Title:      notification.Title,
		Message:    notification.Message,
		EventType:  notification.EventType,
		Recipient:  notification.Recipient,
		Channel:    string(notification.Channel),
		Status:     string(notification.Status),
		Priority:   string(notification.Priority),
		ReadAt:     notification.ReadAt,
		CreatedAt:  notification.CreatedAt,
		UpdatedAt:  notification.UpdatedAt,
	}

	result := r.db.Model(&models.Notification{}).Where("id = ?", notification.ID).Updates(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update notification: %w", result.Error)
	}
	return nil
}

func (r *NotificationRepo) MarkRead(id string) error {
	now := time.Now().UTC()
	result := r.db.Model(&models.Notification{}).Where("id = ?", id).Update("read_at", now)
	if result.Error != nil {
		return fmt.Errorf("failed to mark notification as read: %w", result.Error)
	}
	return nil
}

func (r *NotificationRepo) modelToDomain(m *models.Notification) *notify.Notification {
	return &notify.Notification{
		ID:         m.ID,
		Title:      m.Title,
		Message:    m.Message,
		EventType:  m.EventType,
		Recipient:  m.Recipient,
		Channel:    notify.ChannelType(m.Channel),
		Status:     notify.NotificationStatus(m.Status),
		Priority:   notify.NotificationPriority(m.Priority),
		ReadAt:     m.ReadAt,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

func (r *NotificationRepo) modelsToDomain(models []models.Notification) []*notify.Notification {
	result := make([]*notify.Notification, len(models))
	for i, m := range models {
		result[i] = r.modelToDomain(&m)
	}
	return result
}

type SubscriptionRepo struct {
	db *gorm.DB
}

func NewSubscriptionRepo(db *gorm.DB) notify.SubscriptionRepo {
	return &SubscriptionRepo{db: db}
}

func (r *SubscriptionRepo) Save(sub *notify.Subscription) error {
	model := &models.Subscription{
		ID:         sub.ID,
		Recipient:  sub.Recipient,
		EventType:  sub.EventType,
		Channel:    string(sub.Channel),
		IsActive:   sub.IsActive,
		CreatedAt:  sub.CreatedAt,
		UpdatedAt:  sub.UpdatedAt,
	}

	result := r.db.Save(model)
	if result.Error != nil {
		return fmt.Errorf("failed to save subscription: %w", result.Error)
	}
	return nil
}

func (r *SubscriptionRepo) FindByRecipient(recipient string) ([]*notify.Subscription, error) {
	var models []models.Subscription
	result := r.db.Where("recipient = ?", recipient).Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find subscriptions by recipient: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *SubscriptionRepo) FindByEventType(eventType string) ([]*notify.Subscription, error) {
	var models []models.Subscription
	result := r.db.Where("event_type = ?", eventType).Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find subscriptions by event type: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *SubscriptionRepo) FindByRecipientAndEventType(recipient, eventType string) ([]*notify.Subscription, error) {
	var models []models.Subscription
	result := r.db.Where("recipient = ? AND event_type = ?", recipient, eventType).Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find subscriptions by recipient and event type: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *SubscriptionRepo) Update(sub *notify.Subscription) error {
	model := &models.Subscription{
		ID:         sub.ID,
		Recipient:  sub.Recipient,
		EventType:  sub.EventType,
		Channel:    string(sub.Channel),
		IsActive:   sub.IsActive,
		CreatedAt:  sub.CreatedAt,
		UpdatedAt:  sub.UpdatedAt,
	}

	result := r.db.Model(&models.Subscription{}).Where("id = ?", sub.ID).Updates(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update subscription: %w", result.Error)
	}
	return nil
}

func (r *SubscriptionRepo) Delete(id string) error {
	result := r.db.Delete(&models.Subscription{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete subscription: %w", result.Error)
	}
	return nil
}

func (r *SubscriptionRepo) modelToDomain(m *models.Subscription) *notify.Subscription {
	return &notify.Subscription{
		ID:        m.ID,
		Recipient: m.Recipient,
		EventType: m.EventType,
		Channel:   notify.ChannelType(m.Channel),
		IsActive:  m.IsActive,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (r *SubscriptionRepo) modelsToDomain(models []models.Subscription) []*notify.Subscription {
	result := make([]*notify.Subscription, len(models))
	for i, m := range models {
		result[i] = r.modelToDomain(&m)
	}
	return result
}
