package sqlite

import (
	"fmt"

	"github.com/mengri/nbcoder/domain/project"
	"github.com/mengri/nbcoder/infrastructure/database/models"
	"gorm.io/gorm"
)

type ConfigChangeLogRepo struct {
	dbProvider DBProvider
}

func NewConfigChangeLogRepo(dbProvider DBProvider) project.ConfigChangeLogRepo {
	return &ConfigChangeLogRepo{dbProvider: dbProvider}
}

func (r *ConfigChangeLogRepo) getDB(projectName string) (*gorm.DB, error) {
	if projectName == "" {
		return r.dbProvider.GetGlobalDB(), nil
	}
	return r.dbProvider.GetProjectDB(projectName)
}

func (r *ConfigChangeLogRepo) Save(log *project.ConfigChangeLog) error {
	db, err := r.getDB(log.ProjectName)
	if err != nil {
		return err
	}

	model := &models.ConfigChangeLog{
		ID:          log.ID,
		ProjectName: log.ProjectName,
		ConfigKey:   log.ConfigKey,
		OldValue:    log.OldValue,
		NewValue:    log.NewValue,
		ChangedAt:   log.ChangedAt,
		ChangedBy:   log.ChangedBy,
	}

	result := db.Save(model)
	if result.Error != nil {
		return fmt.Errorf("failed to save config change log: %w", result.Error)
	}
	return nil
}

func (r *ConfigChangeLogRepo) FindByProjectName(projectName string) ([]*project.ConfigChangeLog, error) {
	db, err := r.getDB(projectName)
	if err != nil {
		return nil, err
	}

	var models []models.ConfigChangeLog
	result := db.Where("project_name = ?", projectName).Order("changed_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find config change logs by project name: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *ConfigChangeLogRepo) FindByConfigKey(projectName, configKey string) ([]*project.ConfigChangeLog, error) {
	db, err := r.getDB(projectName)
	if err != nil {
		return nil, err
	}

	var models []models.ConfigChangeLog
	result := db.Where("project_name = ? AND config_key = ?", projectName, configKey).Order("changed_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find config change logs by project name and config key: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *ConfigChangeLogRepo) FindByChangedBy(changedBy string) ([]*project.ConfigChangeLog, error) {
	db, err := r.getDB("")
	if err != nil {
		return nil, err
	}

	var models []models.ConfigChangeLog
	result := db.Where("changed_by = ?", changedBy).Order("changed_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find config change logs by changed by: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *ConfigChangeLogRepo) FindByTimeRange(start, end int64) ([]*project.ConfigChangeLog, error) {
	db, err := r.getDB("")
	if err != nil {
		return nil, err
	}

	var models []models.ConfigChangeLog
	result := db.Where("changed_at BETWEEN ? AND ?", start, end).Order("changed_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find config change logs by time range: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *ConfigChangeLogRepo) modelToDomain(m *models.ConfigChangeLog) *project.ConfigChangeLog {
	return &project.ConfigChangeLog{
		ID:          m.ID,
		ProjectName: m.ProjectName,
		ConfigKey:   m.ConfigKey,
		OldValue:    m.OldValue,
		NewValue:    m.NewValue,
		ChangedAt:   m.ChangedAt,
		ChangedBy:   m.ChangedBy,
	}
}

func (r *ConfigChangeLogRepo) modelsToDomain(models []models.ConfigChangeLog) []*project.ConfigChangeLog {
	result := make([]*project.ConfigChangeLog, len(models))
	for i, m := range models {
		result[i] = r.modelToDomain(&m)
	}
	return result
}
