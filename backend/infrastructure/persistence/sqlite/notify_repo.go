package sqlite

import (
	"fmt"
	"time"

	"github.com/mengri/nbcoder/domain/notify"
	"github.com/mengri/nbcoder/infrastructure/database/models"
	"gorm.io/gorm"
)

type NotificationRepo struct {
	dbProvider DBProvider
}

func NewNotificationRepo(dbProvider DBProvider) notify.NotificationRepo {
	return &NotificationRepo{dbProvider: dbProvider}
}

func (r *NotificationRepo) getDB() (*gorm.DB, error) {
	return r.dbProvider.GetGlobalDB(), nil
}

func (r *NotificationRepo) Save(notification *notify.Notification) error {
	db, err := r.getDB()
	if err != nil {
		return err
	}

	var readAt *time.Time
	if notification.SentAt != nil {
		readAt = notification.SentAt
	}

	status := "UNREAD"
	if notification.Read {
		status = "READ"
	}

	model := &models.Notification{
		ID:         notification.ID,
		Title:      notification.Title,
		Message:    notification.Content,
		EventType:  notification.EventType,
		Recipient:  notification.Recipient,
		Channel:    string(notification.Channel),
		Status:     status,
		Priority:   "NORMAL",
		ReadAt:     readAt,
		CreatedAt:  notification.CreatedAt,
		UpdatedAt:  time.Now(),
	}

	result := db.Save(model)
	if result.Error != nil {
		return fmt.Errorf("failed to save notification: %w", result.Error)
	}
	return nil
}

func (r *NotificationRepo) FindByID(id string) (*notify.Notification, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var model models.Notification
	result := db.First(&model, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find notification by id: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *NotificationRepo) FindByRecipient(recipient string) ([]*notify.Notification, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var models []models.Notification
	result := db.Where("recipient = ?", recipient).Order("created_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find notifications by recipient: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *NotificationRepo) FindByEventType(eventType string) ([]*notify.Notification, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var models []models.Notification
	result := db.Where("event_type = ?", eventType).Order("created_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find notifications by event type: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *NotificationRepo) FindByRecipientAndChannel(recipient string, channel notify.ChannelType) ([]*notify.Notification, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var models []models.Notification
	result := db.Where("recipient = ? AND channel = ?", recipient, string(channel)).Order("created_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find notifications by recipient and channel: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *NotificationRepo) Update(notification *notify.Notification) error {
	db, err := r.getDB()
	if err != nil {
		return err
	}

	var readAt *time.Time
	if notification.SentAt != nil {
		readAt = notification.SentAt
	}

	status := "UNREAD"
	if notification.Read {
		status = "READ"
	}

	model := &models.Notification{
		Title:      notification.Title,
		Message:    notification.Content,
		EventType:  notification.EventType,
		Recipient:  notification.Recipient,
		Channel:    string(notification.Channel),
		Status:     status,
		ReadAt:     readAt,
	}

	result := db.Model(&models.Notification{}).Where("id = ?", notification.ID).Updates(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update notification: %w", result.Error)
	}
	return nil
}

func (r *NotificationRepo) MarkRead(id string) error {
	db, err := r.getDB()
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	result := db.Model(&models.Notification{}).Where("id = ?", id).Update("read_at", now)
	if result.Error != nil {
		return fmt.Errorf("failed to mark notification as read: %w", result.Error)
	}
	return nil
}

func (r *NotificationRepo) modelToDomain(m *models.Notification) *notify.Notification {
	var sentAt *time.Time
	if m.ReadAt != nil {
		sentAt = m.ReadAt
	}

	return &notify.Notification{
		ID:        m.ID,
		Title:     m.Title,
		Content:   m.Message,
		EventType: m.EventType,
		Channel:   notify.ChannelType(m.Channel),
		Recipient: m.Recipient,
		Status:    notify.NotificationStatus(m.Status),
		Read:      m.Status == "READ",
		CreatedAt: m.CreatedAt,
		SentAt:    sentAt,
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
	dbProvider DBProvider
}

func NewSubscriptionRepo(dbProvider DBProvider) notify.SubscriptionRepo {
	return &SubscriptionRepo{dbProvider: dbProvider}
}

func (r *SubscriptionRepo) getDB() (*gorm.DB, error) {
	return r.dbProvider.GetGlobalDB(), nil
}

func (r *SubscriptionRepo) Save(sub *notify.Subscription) error {
	db, err := r.getDB()
	if err != nil {
		return err
	}

	model := &models.Subscription{
		ID:         sub.ID,
		Recipient:  sub.Recipient,
		EventType:  sub.EventType,
		Channel:    string(sub.Channel),
		IsActive:   !sub.Muted,
		CreatedAt:  sub.CreatedAt,
		UpdatedAt:  time.Now(),
	}

	result := db.Save(model)
	if result.Error != nil {
		return fmt.Errorf("failed to save subscription: %w", result.Error)
	}
	return nil
}

func (r *SubscriptionRepo) FindByRecipient(recipient string) ([]*notify.Subscription, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var models []models.Subscription
	result := db.Where("recipient = ?", recipient).Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find subscriptions by recipient: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *SubscriptionRepo) FindByEventType(eventType string) ([]*notify.Subscription, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var models []models.Subscription
	result := db.Where("event_type = ?", eventType).Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find subscriptions by event type: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *SubscriptionRepo) FindByRecipientAndEventType(recipient, eventType string) ([]*notify.Subscription, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var models []models.Subscription
	result := db.Where("recipient = ? AND event_type = ?", recipient, eventType).Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find subscriptions by recipient and event type: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *SubscriptionRepo) Update(sub *notify.Subscription) error {
	db, err := r.getDB()
	if err != nil {
		return err
	}

	model := &models.Subscription{
		Recipient: sub.Recipient,
		EventType: sub.EventType,
		Channel:    string(sub.Channel),
		IsActive:  !sub.Muted,
		UpdatedAt:  time.Now(),
	}

	result := db.Model(&models.Subscription{}).Where("id = ?", sub.ID).Updates(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update subscription: %w", result.Error)
	}
	return nil
}

func (r *SubscriptionRepo) Delete(id string) error {
	db, err := r.getDB()
	if err != nil {
		return err
	}

	result := db.Delete(&models.Subscription{}, "id = ?", id)
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
		Muted:     !m.IsActive,
		CreatedAt: m.CreatedAt,
	}
}

func (r *SubscriptionRepo) modelsToDomain(models []models.Subscription) []*notify.Subscription {
	result := make([]*notify.Subscription, len(models))
	for i, m := range models {
		result[i] = r.modelToDomain(&m)
	}
	return result
}
