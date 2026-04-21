package sqlite

import (
	"fmt"
	"time"

	"github.com/mengri/nbcoder/domain/notify"
	"github.com/mengri/nbcoder/infrastructure/database/models"
	"gorm.io/gorm"
)

type SubscriptionPreferenceRepo struct {
	dbProvider DBProvider
}

func NewSubscriptionPreferenceRepo(dbProvider DBProvider) notify.SubscriptionPreferenceRepo {
	return &SubscriptionPreferenceRepo{dbProvider: dbProvider}
}

func (r *SubscriptionPreferenceRepo) getDB() (*gorm.DB, error) {
	return r.dbProvider.GetGlobalDB(), nil
}

func (r *SubscriptionPreferenceRepo) Save(pref *notify.SubscriptionPreference) error {
	db, err := r.getDB()
	if err != nil {
		return err
	}

	model := &models.SubscriptionPreference{
		ID:              pref.ID,
		Recipient:       pref.Recipient,
		EventType:       pref.EventType,
		EnabledChannels: "",
		MinPriority:     "",
		DigestEnabled:   false,
		DigestFrequency: "",
	}

	result := db.Save(model)
	if result.Error != nil {
		return fmt.Errorf("failed to save subscription preference: %w", result.Error)
	}
	return nil
}

func (r *SubscriptionPreferenceRepo) FindByRecipient(recipient string) ([]*notify.SubscriptionPreference, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var models []models.SubscriptionPreference
	result := db.Where("recipient = ?", recipient).Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find subscription preferences by recipient: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *SubscriptionPreferenceRepo) FindByEventType(eventType string) ([]*notify.SubscriptionPreference, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var models []models.SubscriptionPreference
	result := db.Where("event_type = ?", eventType).Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find subscription preferences by event type: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *SubscriptionPreferenceRepo) FindByRecipientAndEventType(recipient, eventType string) (*notify.SubscriptionPreference, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var model models.SubscriptionPreference
	result := db.Where("recipient = ? AND event_type = ?", recipient, eventType).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find subscription preference: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *SubscriptionPreferenceRepo) Update(pref *notify.SubscriptionPreference) error {
	db, err := r.getDB()
	if err != nil {
		return err
	}

	model := &models.SubscriptionPreference{
		Recipient: pref.Recipient,
		EventType: pref.EventType,
	}

	result := db.Model(&models.SubscriptionPreference{}).Where("id = ?", pref.ID).Updates(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update subscription preference: %w", result.Error)
	}
	return nil
}

func (r *SubscriptionPreferenceRepo) Delete(id string) error {
	db, err := r.getDB()
	if err != nil {
		return err
	}

	result := db.Delete(&models.SubscriptionPreference{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete subscription preference: %w", result.Error)
	}
	return nil
}

func (r *SubscriptionPreferenceRepo) modelToDomain(m *models.SubscriptionPreference) *notify.SubscriptionPreference {
	return &notify.SubscriptionPreference{
		ID:            m.ID,
		Recipient:     m.Recipient,
		EventType:     m.EventType,
		MutedChannels: []notify.ChannelType{},
	}
}

func (r *SubscriptionPreferenceRepo) modelsToDomain(models []models.SubscriptionPreference) []*notify.SubscriptionPreference {
	result := make([]*notify.SubscriptionPreference, len(models))
	for i, m := range models {
		result[i] = r.modelToDomain(&m)
	}
	return result
}

type NotificationTemplateRepo struct {
	dbProvider DBProvider
}

func NewNotificationTemplateRepo(dbProvider DBProvider) notify.NotificationTemplateRepo {
	return &NotificationTemplateRepo{dbProvider: dbProvider}
}

func (r *NotificationTemplateRepo) getDB() (*gorm.DB, error) {
	return r.dbProvider.GetGlobalDB(), nil
}

func (r *NotificationTemplateRepo) Save(template *notify.NotificationTemplate) error {
	db, err := r.getDB()
	if err != nil {
		return err
	}

	model := &models.NotificationTemplate{
		ID:        template.ID,
		Name:      template.Name,
		EventType: template.EventType,
		Subject:   template.SubjectTemplate,
		Body:      template.BodyTemplate,
		Channel:   "",
		Variables: "",
		IsActive:  true,
		CreatedAt: template.CreatedAt,
		UpdatedAt: template.UpdatedAt,
	}

	result := db.Save(model)
	if result.Error != nil {
		return fmt.Errorf("failed to save notification template: %w", result.Error)
	}
	return nil
}

func (r *NotificationTemplateRepo) FindByID(id string) (*notify.NotificationTemplate, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var model models.NotificationTemplate
	result := db.First(&model, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find notification template by id: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *NotificationTemplateRepo) FindByEventType(eventType string) ([]*notify.NotificationTemplate, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var models []models.NotificationTemplate
	result := db.Where("event_type = ? AND is_active = ?", eventType, true).Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find notification templates by event type: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *NotificationTemplateRepo) Update(template *notify.NotificationTemplate) error {
	db, err := r.getDB()
	if err != nil {
		return err
	}

	model := &models.NotificationTemplate{
		Name:      template.Name,
		EventType: template.EventType,
		Subject:   template.SubjectTemplate,
		Body:      template.BodyTemplate,
	}

	result := db.Model(&models.NotificationTemplate{}).Where("id = ?", template.ID).Updates(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update notification template: %w", result.Error)
	}
	return nil
}

func (r *NotificationTemplateRepo) Delete(id string) error {
	db, err := r.getDB()
	if err != nil {
		return err
	}

	result := db.Delete(&models.NotificationTemplate{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete notification template: %w", result.Error)
	}
	return nil
}

func (r *NotificationTemplateRepo) FindByName(name string) (*notify.NotificationTemplate, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var model models.NotificationTemplate
	result := db.Where("name = ?", name).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find notification template by name: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *NotificationTemplateRepo) modelToDomain(m *models.NotificationTemplate) *notify.NotificationTemplate {
	return &notify.NotificationTemplate{
		ID:              m.ID,
		Name:            m.Name,
		EventType:       m.EventType,
		SubjectTemplate: m.Subject,
		BodyTemplate:    m.Body,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

func (r *NotificationTemplateRepo) modelsToDomain(models []models.NotificationTemplate) []*notify.NotificationTemplate {
	result := make([]*notify.NotificationTemplate, len(models))
	for i, m := range models {
		result[i] = r.modelToDomain(&m)
	}
	return result
}

type NotificationHistoryRepo struct {
	dbProvider DBProvider
}

func NewNotificationHistoryRepo(dbProvider DBProvider) notify.NotificationHistoryRepo {
	return &NotificationHistoryRepo{dbProvider: dbProvider}
}

func (r *NotificationHistoryRepo) getDB() (*gorm.DB, error) {
	return r.dbProvider.GetGlobalDB(), nil
}

func (r *NotificationHistoryRepo) Save(history *notify.NotificationHistory) error {
	db, err := r.getDB()
	if err != nil {
		return err
	}

	model := &models.NotificationHistory{
		ID:             history.ID,
		NotificationID: history.NotificationID,
		Channel:        string(history.Channel),
		Recipient:      history.Recipient,
		Status:         string(history.Status),
		SentAt:         &history.SentAt,
		Error:          history.Error,
	}

	result := db.Save(model)
	if result.Error != nil {
		return fmt.Errorf("failed to save notification history: %w", result.Error)
	}
	return nil
}

func (r *NotificationHistoryRepo) FindByNotificationID(notificationID string) ([]*notify.NotificationHistory, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var models []models.NotificationHistory
	result := db.Where("notification_id = ?", notificationID).Order("created_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find notification history by notification id: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *NotificationHistoryRepo) FindByRecipient(recipient string) ([]*notify.NotificationHistory, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var models []models.NotificationHistory
	result := db.Where("recipient = ?", recipient).Order("created_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find notification history by recipient: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *NotificationHistoryRepo) FindByTimeRange(start, end time.Time) ([]*notify.NotificationHistory, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var models []models.NotificationHistory
	result := db.Where("created_at BETWEEN ? AND ?", start, end).Order("created_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find notification history by time range: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *NotificationHistoryRepo) modelToDomain(m *models.NotificationHistory) *notify.NotificationHistory {
	return &notify.NotificationHistory{
		ID:             m.ID,
		NotificationID: m.NotificationID,
		Channel:        notify.ChannelType(m.Channel),
		Recipient:      m.Recipient,
		Status:         notify.HistoryStatus(m.Status),
		SentAt:         *m.SentAt,
		Error:          m.Error,
	}
}

func (r *NotificationHistoryRepo) modelsToDomain(models []models.NotificationHistory) []*notify.NotificationHistory {
	result := make([]*notify.NotificationHistory, len(models))
	for i, m := range models {
		result[i] = r.modelToDomain(&m)
	}
	return result
}
