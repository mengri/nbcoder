package project

import (
	"fmt"
)

type ConfigManager struct {
	globalRepo  GlobalConfigRepo
	projectRepo ProjectConfigRepo
}

func NewConfigManager(globalRepo GlobalConfigRepo, projectRepo ProjectConfigRepo) *ConfigManager {
	return &ConfigManager{
		globalRepo:  globalRepo,
		projectRepo: projectRepo,
	}
}

func (cm *ConfigManager) GetConfig(projectID, key string) (string, error) {
	projectConfig, err := cm.projectRepo.FindByKey(projectID, key)
	if err == nil && projectConfig != nil {
		return projectConfig.Value, nil
	}

	globalConfig, err := cm.globalRepo.FindByKey(key)
	if err != nil {
		return "", fmt.Errorf("config key %s not found in project or global config", key)
	}

	return globalConfig.Value, nil
}

func (cm *ConfigManager) GetGlobalConfig(key string) (*GlobalConfig, error) {
	config, err := cm.globalRepo.FindByKey(key)
	if err != nil {
		return nil, fmt.Errorf("global config key %s not found", key)
	}
	return config, nil
}

func (cm *ConfigManager) GetProjectConfig(projectID, key string) (*ProjectConfig, error) {
	config, err := cm.projectRepo.FindByKey(projectID, key)
	if err != nil {
		return nil, fmt.Errorf("project config key %s not found", key)
	}
	return config, nil
}

func (cm *ConfigManager) GetAllGlobalConfigs() ([]*GlobalConfig, error) {
	return cm.globalRepo.FindAll()
}

func (cm *ConfigManager) GetAllProjectConfigs(projectID string) ([]*ProjectConfig, error) {
	return cm.projectRepo.FindByProjectID(projectID)
}

func (cm *ConfigManager) GetAllConfigs(projectID string) ([]*ConfigItem, error) {
	var items []*ConfigItem

	globalConfigs, err := cm.globalRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch global configs: %w", err)
	}

	for _, config := range globalConfigs {
		item := NewConfigItem(ConfigScopeGlobal, "", config.Key, config.Value)
		items = append(items, item)
	}

	projectConfigs, err := cm.projectRepo.FindByProjectID(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch project configs: %w", err)
	}

	for _, config := range projectConfigs {
		item := NewConfigItem(ConfigScopeProject, projectID, config.Key, config.Value)
		items = append(items, item)
	}

	return items, nil
}

func (cm *ConfigManager) SetGlobalConfig(id, key, value string) error {
	if err := cm.validateConfigKey(key); err != nil {
		return err
	}

	config := NewGlobalConfig(id, key, value)
	if err := config.Validate(); err != nil {
		return err
	}

	return cm.globalRepo.Save(config)
}

func (cm *ConfigManager) SetProjectConfig(id, projectID, key, value string) error {
	if err := cm.validateConfigKey(key); err != nil {
		return err
	}

	config := NewProjectConfig(id, projectID, key, value)
	if err := config.Validate(); err != nil {
		return err
	}

	return cm.projectRepo.Save(config)
}

func (cm *ConfigManager) UpdateGlobalConfig(key, value string) error {
	config, err := cm.globalRepo.FindByKey(key)
	if err != nil {
		return fmt.Errorf("global config with key %s not found", key)
	}

	config.Update(value)
	return cm.globalRepo.Update(config)
}

func (cm *ConfigManager) UpdateProjectConfig(id, value string) error {
	configs, err := cm.projectRepo.FindAll()
	if err != nil {
		return fmt.Errorf("failed to fetch project configs: %w", err)
	}

	var targetConfig *ProjectConfig
	for _, config := range configs {
		if config.ID == id {
			targetConfig = config
			break
		}
	}

	if targetConfig == nil {
		return fmt.Errorf("project config with id %s not found", id)
	}

	targetConfig.Update(value)
	return cm.projectRepo.Update(targetConfig)
}

func (cm *ConfigManager) DeleteGlobalConfig(id string) error {
	return cm.globalRepo.Delete(id)
}

func (cm *ConfigManager) DeleteProjectConfig(id string) error {
	return cm.projectRepo.Delete(id)
}

func (cm *ConfigManager) CopyGlobalToProject(projectID, key string) error {
	globalConfig, err := cm.globalRepo.FindByKey(key)
	if err != nil {
		return fmt.Errorf("global config key %s not found", key)
	}

	projectConfig := NewProjectConfig(globalConfig.ID, projectID, key, globalConfig.Value)
	return cm.projectRepo.Save(projectConfig)
}

func (cm *ConfigManager) GetConfigHierarchy(projectID, key string) (*ConfigItem, *ConfigItem, error) {
	globalConfig, globalErr := cm.globalRepo.FindByKey(key)
	projectConfig, projectErr := cm.projectRepo.FindByKey(projectID, key)

	if globalErr != nil && projectErr != nil {
		return nil, nil, fmt.Errorf("config key %s not found in project or global config", key)
	}

	var globalItem, projectItem *ConfigItem

	if globalConfig != nil {
		globalItem = NewConfigItem(ConfigScopeGlobal, "", globalConfig.Key, globalConfig.Value)
	}

	if projectConfig != nil {
		projectItem = NewConfigItem(ConfigScopeProject, projectID, projectConfig.Key, projectConfig.Value)
	}

	return globalItem, projectItem, nil
}

func (cm *ConfigManager) validateConfigKey(key string) error {
	if key == "" {
		return fmt.Errorf("config key cannot be empty")
	}
	return nil
}