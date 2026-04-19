package project

import (
	"fmt"
	"time"
)

type RuleType string

const (
	RuleTypeNaming  RuleType = "NAMING"
	RuleTypeCommit  RuleType = "COMMIT"
	RuleTypeReview  RuleType = "REVIEW"
	RuleTypeTest    RuleType = "TEST"
)

func (rt RuleType) IsValid() bool {
	switch rt {
	case RuleTypeNaming, RuleTypeCommit, RuleTypeReview, RuleTypeTest:
		return true
	}
	return false
}

type DevStandard struct {
	ID          string    `json:"id"`
	ProjectID   string    `json:"project_id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	RuleType    RuleType  `json:"rule_type"`
	Rule        string    `json:"rule"`
	Enabled     bool      `json:"enabled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewDevStandard(id, projectID, name, description string, ruleType RuleType, rule string, enabled bool) *DevStandard {
	now := time.Now().UTC()
	return &DevStandard{
		ID:          id,
		ProjectID:   projectID,
		Name:        name,
		Description: description,
		RuleType:    ruleType,
		Rule:        rule,
		Enabled:     enabled,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (ds *DevStandard) Validate() error {
	if ds.ProjectID == "" {
		return fmt.Errorf("project_id is required")
	}
	if ds.Name == "" {
		return fmt.Errorf("name is required")
	}
	if !ds.RuleType.IsValid() {
		return fmt.Errorf("invalid rule_type: %s", ds.RuleType)
	}
	if ds.Rule == "" {
		return fmt.Errorf("rule is required")
	}
	return nil
}

type BranchPolicyConfig struct {
	ID          string    `json:"id"`
	ProjectID   string    `json:"project_id"`
	Pattern     string    `json:"pattern"`
	Description string    `json:"description,omitempty"`
	IsDefault   bool      `json:"is_default"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewBranchPolicyConfig(id, projectID, pattern, description string, isDefault bool) *BranchPolicyConfig {
	now := time.Now().UTC()
	return &BranchPolicyConfig{
		ID:          id,
		ProjectID:   projectID,
		Pattern:     pattern,
		Description: description,
		IsDefault:   isDefault,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (bpc *BranchPolicyConfig) Validate() error {
	if bpc.ProjectID == "" {
		return fmt.Errorf("project_id is required")
	}
	if bpc.Pattern == "" {
		return fmt.Errorf("pattern is required")
	}
	return nil
}
