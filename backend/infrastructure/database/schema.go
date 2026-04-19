package database

import (
	"fmt"

	"github.com/mengri/nbcoder/infrastructure/database/models"
	"gorm.io/gorm"
)

func InitSchema(db *gorm.DB) error {
	modelsList := []interface{}{
		&models.Project{},
		&models.ProjectConfig{},
		&models.Standards{},
		&models.DevStandard{},
		&models.BranchPolicyConfig{},
		&models.ProjectLifecycle{},
		&models.ConfigChangeLog{},
		&models.Skill{},
		&models.Directory{},
		&models.Provider{},
		&models.Model{},
		&models.ModelChain{},
		&models.CallLog{},
		&models.Repository{},
		&models.PullRequest{},
		&models.Commit{},
		&models.Card{},
		&models.CardDependency{},
		&models.Pipeline{},
		&models.StageRecord{},
		&models.Task{},
		&models.AgentExecution{},
		&models.Document{},
		&models.DocumentChunk{},
		&models.DocumentIndex{},
		&models.CloneInstance{},
		&models.Notification{},
		&models.Subscription{},
		&models.SubscriptionPreference{},
		&models.NotificationTemplate{},
		&models.NotificationHistory{},
	}

	for _, model := range modelsList {
		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to auto migrate model %T: %w", model, err)
		}
	}

	return nil
}

func DropSchema(db *gorm.DB, models []interface{}) error {
	for _, model := range models {
		err := db.Migrator().DropTable(model)
		if err != nil {
			return fmt.Errorf("failed to drop table for model %T: %w", model, err)
		}
	}
	return nil
}

func GetAllModels() []interface{} {
	return []interface{}{
		&models.Project{},
		&models.ProjectConfig{},
		&models.Standards{},
		&models.DevStandard{},
		&models.BranchPolicyConfig{},
		&models.ProjectLifecycle{},
		&models.ConfigChangeLog{},
		&models.Card{},
		&models.CardDependency{},
		&models.Pipeline{},
		&models.StageRecord{},
		&models.Task{},
		&models.AgentExecution{},
		&models.Skill{},
		&models.Provider{},
		&models.Model{},
		&models.ModelChain{},
		&models.CallLog{},
		&models.Document{},
		&models.DocumentChunk{},
		&models.DocumentIndex{},
		&models.Directory{},
		&models.Repository{},
		&models.PullRequest{},
		&models.Commit{},
		&models.Notification{},
		&models.Subscription{},
		&models.SubscriptionPreference{},
		&models.NotificationTemplate{},
		&models.NotificationHistory{},
		&models.CloneInstance{},
	}
}
