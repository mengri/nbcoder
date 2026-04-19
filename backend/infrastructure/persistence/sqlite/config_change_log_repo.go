package sqlite

import (
	"fmt"

	"github.com/mengri/nbcoder/domain/project"
	"github.com/mengri/nbcoder/infrastructure/database/models"
	"gorm.io/gorm"
)

type ConfigChangeLogRepo struct {
	db *gorm.DB
}

func NewConfigChangeLogRepo(db *gorm.DB) project.ConfigChangeLogRepo {
	return &ConfigChangeLogRepo{db: db}
}

func (r *ConfigChangeLogRepo) Save(log *project.ConfigChangeLog) error {
	model := &models.ConfigChangeLog{
		ID:         log.ID,
		ProjectID:  log.ProjectID,
		ConfigKey:  log.ConfigKey,
		OldValue:   log.OldValue,
		NewValue:   log.NewValue,
		ChangedAt:  log.ChangedAt,
		ChangedBy:  log.ChangedBy,
	}

	result := r.db.Save(model)
	if result.Error != nil {
		return fmt.Errorf("failed to save config change log: %w", result.Error)
	}
	return nil
}

func (r *ConfigChangeLogRepo) FindByProjectID(projectID string) ([]*project.ConfigChangeLog, error) {
	var models []models.ConfigChangeLog
	result := r.db.Where("project_id = ?", projectID).Order("changed_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find config change logs by project id: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *ConfigChangeLogRepo) FindByConfigKey(projectID, configKey string) ([]*project.ConfigChangeLog, error) {
	var models []models.ConfigChangeLog
	result := r.db.Where("project_id = ? AND config_key = ?", projectID, configKey).Order("changed_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find config change logs by project id and config key: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *ConfigChangeLogRepo) FindByChangedBy(changedBy string) ([]*project.ConfigChangeLog, error) {
	var models []models.ConfigChangeLog
	result := r.db.Where("changed_by = ?", changedBy).Order("changed_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find config change logs by changed by: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *ConfigChangeLogRepo) FindByTimeRange(start, end int64) ([]*project.ConfigChangeLog, error) {
	var models []models.ConfigChangeLog
	result := r.db.Where("changed_at BETWEEN ? AND ?", start, end).Order("changed_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find config change logs by time range: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *ConfigChangeLogRepo) modelToDomain(m *models.ConfigChangeLog) *project.ConfigChangeLog {
	return &project.ConfigChangeLog{
		ID:         m.ID,
		ProjectID:  m.ProjectID,
		ConfigKey:  m.ConfigKey,
		OldValue:   m.OldValue,
		NewValue:   m.NewValue,
		ChangedAt:  m.ChangedAt,
		ChangedBy:  m.ChangedBy,
	}
}

func (r *ConfigChangeLogRepo) modelsToDomain(models []models.ConfigChangeLog) []*project.ConfigChangeLog {
	result := make([]*project.ConfigChangeLog, len(models))
	for i, m := range models {
		result[i] = r.modelToDomain(&m)
	}
	return result
}
