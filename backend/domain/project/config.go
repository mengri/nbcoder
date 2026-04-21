package project

import (
	"fmt"
	"time"
)

type GlobalConfig struct {
	ID        string    `json:"id"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewGlobalConfig(id, key, value string) *GlobalConfig {
	now := time.Now().UTC()
	return &GlobalConfig{
		ID:        id,
		Key:       key,
		Value:     value,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (c *GlobalConfig) Validate() error {
	if c.Key == "" {
		return fmt.Errorf("global config key is required")
	}
	return nil
}

func (c *GlobalConfig) Update(value string) {
	c.Value = value
	c.UpdatedAt = time.Now().UTC()
}

type ConfigScope string

const (
	ConfigScopeGlobal   ConfigScope = "GLOBAL"
	ConfigScopeProject ConfigScope = "PROJECT"
)

type ConfigItem struct {
	Scope      ConfigScope `json:"scope"`
	ProjectName string     `json:"project_name,omitempty"`
	Key        string     `json:"key"`
	Value      string     `json:"value"`
}

func NewConfigItem(scope ConfigScope, projectName, key, value string) *ConfigItem {
	return &ConfigItem{
		Scope:      scope,
		ProjectName: projectName,
		Key:        key,
		Value:      value,
	}
}

type ProjectConfig struct {
	ID           string    `json:"id"`
	ProjectName  string    `json:"project_name"`
	Key          string    `json:"key"`
	Value        string    `json:"value"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func NewProjectConfig(id, projectName, key, value string) *ProjectConfig {
	now := time.Now().UTC()
	return &ProjectConfig{
		ID:          id,
		ProjectName: projectName,
		Key:         key,
		Value:       value,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (c *ProjectConfig) Validate() error {
	if c.Key == "" {
		return fmt.Errorf("config key is required")
	}
	return nil
}

func (c *ProjectConfig) Update(value string) {
	c.Value = value
	c.UpdatedAt = time.Now().UTC()
}

type ConfigChangeLog struct {
	ID           string    `json:"id"`
	ProjectName  string    `json:"project_name"`
	ConfigKey    string    `json:"config_key"`
	OldValue     string    `json:"old_value"`
	NewValue     string    `json:"new_value"`
	ChangedAt    time.Time `json:"changed_at"`
	ChangedBy    string    `json:"changed_by"`
}

func NewConfigChangeLog(id, projectName, configKey, oldValue, newValue, changedBy string) *ConfigChangeLog {
	return &ConfigChangeLog{
		ID:          id,
		ProjectName: projectName,
		ConfigKey:   configKey,
		OldValue:    oldValue,
		NewValue:    newValue,
		ChangedAt:   time.Now().UTC(),
		ChangedBy:   changedBy,
	}
}

type Standards struct {
	ID                string    `json:"id"`
	ProjectName       string    `json:"project_name"`
	BranchStrategy    string    `json:"branch_strategy,omitempty"`
	TechStack         string    `json:"tech_stack,omitempty"`
	CodingConventions string    `json:"coding_conventions,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func NewStandards(id, projectName, branchStrategy, techStack, codingConventions string) *Standards {
	now := time.Now().UTC()
	return &Standards{
		ID:                id,
		ProjectName:       projectName,
		BranchStrategy:    branchStrategy,
		TechStack:         techStack,
		CodingConventions: codingConventions,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

func (s *Standards) Update(branchStrategy, techStack, codingConventions string) {
	if branchStrategy != "" {
		s.BranchStrategy = branchStrategy
	}
	if techStack != "" {
		s.TechStack = techStack
	}
	if codingConventions != "" {
		s.CodingConventions = codingConventions
	}
	s.UpdatedAt = time.Now().UTC()
}

type ProjectDirectory struct {
	ProjectName string
	Dirs        []string
}

func DefaultProjectDirectory(projectName string) *ProjectDirectory {
	return &ProjectDirectory{
		ProjectName: projectName,
		Dirs: []string{
			".NBCoder",
			".NBCoder/config",
			"knowledge-base",
			"knowledge-base/docs",
			"cards",
		},
	}
}
