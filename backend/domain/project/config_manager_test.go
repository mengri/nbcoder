package project

import (
	"errors"
	"testing"
)

type MockGlobalConfigRepo struct {
	configs map[string]*GlobalConfig
}

func NewMockGlobalConfigRepo() *MockGlobalConfigRepo {
	return &MockGlobalConfigRepo{
		configs: make(map[string]*GlobalConfig),
	}
}

func (m *MockGlobalConfigRepo) Save(config *GlobalConfig) error {
	m.configs[config.ID] = config
	return nil
}

func (m *MockGlobalConfigRepo) FindByID(id string) (*GlobalConfig, error) {
	config, ok := m.configs[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return config, nil
}

func (m *MockGlobalConfigRepo) FindByKey(key string) (*GlobalConfig, error) {
	for _, config := range m.configs {
		if config.Key == key {
			return config, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *MockGlobalConfigRepo) FindAll() ([]*GlobalConfig, error) {
	var configs []*GlobalConfig
	for _, config := range m.configs {
		configs = append(configs, config)
	}
	return configs, nil
}

func (m *MockGlobalConfigRepo) Update(config *GlobalConfig) error {
	m.configs[config.ID] = config
	return nil
}

func (m *MockGlobalConfigRepo) Delete(id string) error {
	delete(m.configs, id)
	return nil
}

type MockProjectConfigRepo struct {
	configs map[string]*ProjectConfig
}

func NewMockProjectConfigRepo() *MockProjectConfigRepo {
	return &MockProjectConfigRepo{
		configs: make(map[string]*ProjectConfig),
	}
}

func (m *MockProjectConfigRepo) Save(config *ProjectConfig) error {
	m.configs[config.ID] = config
	return nil
}

func (m *MockProjectConfigRepo) FindByID(id string) (*ProjectConfig, error) {
	config, ok := m.configs[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return config, nil
}

func (m *MockProjectConfigRepo) FindByProjectID(projectID string) ([]*ProjectConfig, error) {
	var configs []*ProjectConfig
	for _, config := range m.configs {
		if config.ProjectID == projectID {
			configs = append(configs, config)
		}
	}
	return configs, nil
}

func (m *MockProjectConfigRepo) FindByKey(projectID, key string) (*ProjectConfig, error) {
	for _, config := range m.configs {
		if config.ProjectID == projectID && config.Key == key {
			return config, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *MockProjectConfigRepo) FindAll() ([]*ProjectConfig, error) {
	var configs []*ProjectConfig
	for _, config := range m.configs {
		configs = append(configs, config)
	}
	return configs, nil
}

func (m *MockProjectConfigRepo) Update(config *ProjectConfig) error {
	m.configs[config.ID] = config
	return nil
}

func (m *MockProjectConfigRepo) Delete(id string) error {
	delete(m.configs, id)
	return nil
}

func setupTestConfigs() (*MockGlobalConfigRepo, *MockProjectConfigRepo) {
	globalRepo := NewMockGlobalConfigRepo()
	projectRepo := NewMockProjectConfigRepo()

	globalRepo.Save(NewGlobalConfig("global-1", "max_connections", "100"))
	globalRepo.Save(NewGlobalConfig("global-2", "timeout", "30"))

	projectRepo.Save(NewProjectConfig("proj-1", "project-1", "timeout", "60"))
	projectRepo.Save(NewProjectConfig("proj-2", "project-1", "feature_flag", "true"))

	return globalRepo, projectRepo
}

func TestConfigManager_GetConfig(t *testing.T) {
	globalRepo, projectRepo := setupTestConfigs()
	cm := NewConfigManager(globalRepo, projectRepo)

	value, err := cm.GetConfig("project-1", "timeout")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if value != "60" {
		t.Errorf("expected timeout 60, got %s", value)
	}

	value, err = cm.GetConfig("project-1", "max_connections")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if value != "100" {
		t.Errorf("expected max_connections 100, got %s", value)
	}
}

func TestConfigManager_GetConfig_NotFound(t *testing.T) {
	globalRepo, projectRepo := setupTestConfigs()
	cm := NewConfigManager(globalRepo, projectRepo)

	_, err := cm.GetConfig("project-1", "nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent config")
	}
}

func TestConfigManager_GetGlobalConfig(t *testing.T) {
	globalRepo, projectRepo := setupTestConfigs()
	cm := NewConfigManager(globalRepo, projectRepo)

	config, err := cm.GetGlobalConfig("max_connections")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if config.Value != "100" {
		t.Errorf("expected value 100, got %s", config.Value)
	}
}

func TestConfigManager_GetGlobalConfig_NotFound(t *testing.T) {
	globalRepo, projectRepo := setupTestConfigs()
	cm := NewConfigManager(globalRepo, projectRepo)

	_, err := cm.GetGlobalConfig("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent config")
	}
}

func TestConfigManager_GetProjectConfig(t *testing.T) {
	globalRepo, projectRepo := setupTestConfigs()
	cm := NewConfigManager(globalRepo, projectRepo)

	config, err := cm.GetProjectConfig("project-1", "timeout")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if config.Value != "60" {
		t.Errorf("expected value 60, got %s", config.Value)
	}
}

func TestConfigManager_GetProjectConfig_NotFound(t *testing.T) {
	globalRepo, projectRepo := setupTestConfigs()
	cm := NewConfigManager(globalRepo, projectRepo)

	_, err := cm.GetProjectConfig("project-1", "nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent config")
	}
}

func TestConfigManager_GetAllGlobalConfigs(t *testing.T) {
	globalRepo, projectRepo := setupTestConfigs()
	cm := NewConfigManager(globalRepo, projectRepo)

	configs, err := cm.GetAllGlobalConfigs()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(configs) != 2 {
		t.Errorf("expected 2 global configs, got %d", len(configs))
	}
}

func TestConfigManager_GetAllProjectConfigs(t *testing.T) {
	globalRepo, projectRepo := setupTestConfigs()
	cm := NewConfigManager(globalRepo, projectRepo)

	configs, err := cm.GetAllProjectConfigs("project-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(configs) != 2 {
		t.Errorf("expected 2 project configs, got %d", len(configs))
	}
}

func TestConfigManager_GetAllConfigs(t *testing.T) {
	globalRepo, projectRepo := setupTestConfigs()
	cm := NewConfigManager(globalRepo, projectRepo)

	configs, err := cm.GetAllConfigs("project-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(configs) != 4 {
		t.Errorf("expected 4 total configs, got %d", len(configs))
	}

	var globalCount, projectCount int
	for _, config := range configs {
		switch config.Scope {
		case ConfigScopeGlobal:
			globalCount++
		case ConfigScopeProject:
			projectCount++
		}
	}

	if globalCount != 2 {
		t.Errorf("expected 2 global configs, got %d", globalCount)
	}

	if projectCount != 2 {
		t.Errorf("expected 2 project configs, got %d", projectCount)
	}
}

func TestConfigManager_SetGlobalConfig(t *testing.T) {
	globalRepo := NewMockGlobalConfigRepo()
	projectRepo := NewMockProjectConfigRepo()
	cm := NewConfigManager(globalRepo, projectRepo)

	err := cm.SetGlobalConfig("global-3", "new_key", "new_value")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	config, err := cm.GetGlobalConfig("new_key")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if config.Value != "new_value" {
		t.Errorf("expected value new_value, got %s", config.Value)
	}
}

func TestConfigManager_SetProjectConfig(t *testing.T) {
	globalRepo := NewMockGlobalConfigRepo()
	projectRepo := NewMockProjectConfigRepo()
	cm := NewConfigManager(globalRepo, projectRepo)

	err := cm.SetProjectConfig("proj-3", "project-2", "new_key", "new_value")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	config, err := cm.GetProjectConfig("project-2", "new_key")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if config.Value != "new_value" {
		t.Errorf("expected value new_value, got %s", config.Value)
	}
}

func TestConfigManager_UpdateGlobalConfig(t *testing.T) {
	globalRepo, projectRepo := setupTestConfigs()
	cm := NewConfigManager(globalRepo, projectRepo)

	err := cm.UpdateGlobalConfig("max_connections", "200")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	config, err := cm.GetGlobalConfig("max_connections")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if config.Value != "200" {
		t.Errorf("expected value 200, got %s", config.Value)
	}
}

func TestConfigManager_UpdateProjectConfig(t *testing.T) {
	globalRepo, projectRepo := setupTestConfigs()
	cm := NewConfigManager(globalRepo, projectRepo)

	err := cm.UpdateProjectConfig("proj-1", "90")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	config, err := cm.GetProjectConfig("project-1", "timeout")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if config.Value != "90" {
		t.Errorf("expected value 90, got %s", config.Value)
	}
}

func TestConfigManager_DeleteGlobalConfig(t *testing.T) {
	globalRepo, projectRepo := setupTestConfigs()
	cm := NewConfigManager(globalRepo, projectRepo)

	err := cm.DeleteGlobalConfig("global-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	_, err = cm.GetGlobalConfig("max_connections")
	if err == nil {
		t.Error("expected error after deletion")
	}
}

func TestConfigManager_DeleteProjectConfig(t *testing.T) {
	globalRepo, projectRepo := setupTestConfigs()
	cm := NewConfigManager(globalRepo, projectRepo)

	err := cm.DeleteProjectConfig("proj-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	_, err = cm.GetProjectConfig("project-1", "timeout")
	if err == nil {
		t.Error("expected error after deletion")
	}
}

func TestConfigManager_CopyGlobalToProject(t *testing.T) {
	globalRepo, projectRepo := setupTestConfigs()
	cm := NewConfigManager(globalRepo, projectRepo)

	err := cm.CopyGlobalToProject("project-2", "max_connections")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	config, err := cm.GetProjectConfig("project-2", "max_connections")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if config.Value != "100" {
		t.Errorf("expected value 100, got %s", config.Value)
	}
}

func TestConfigManager_GetConfigHierarchy(t *testing.T) {
	globalRepo, projectRepo := setupTestConfigs()
	cm := NewConfigManager(globalRepo, projectRepo)

	globalItem, projectItem, err := cm.GetConfigHierarchy("project-1", "timeout")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if globalItem == nil || globalItem.Key != "timeout" {
		t.Error("expected global item to exist")
	}

	if projectItem == nil || projectItem.Key != "timeout" {
		t.Error("expected project item to exist")
	}
}

func TestConfigManager_GetConfigHierarchy_GlobalOnly(t *testing.T) {
	globalRepo, projectRepo := setupTestConfigs()
	cm := NewConfigManager(globalRepo, projectRepo)

	globalItem, projectItem, err := cm.GetConfigHierarchy("project-1", "max_connections")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if globalItem == nil || globalItem.Key != "max_connections" {
		t.Error("expected global item to exist")
	}

	if projectItem != nil {
		t.Error("expected no project item")
	}
}

func TestConfigManager_ValidateConfigKey(t *testing.T) {
	globalRepo := NewMockGlobalConfigRepo()
	projectRepo := NewMockProjectConfigRepo()
	cm := NewConfigManager(globalRepo, projectRepo)

	err := cm.SetGlobalConfig("global-1", "", "value")
	if err == nil {
		t.Error("expected error for empty key")
	}
}
