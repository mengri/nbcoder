package project

import (
	"testing"
)

func TestRuleType_IsValid(t *testing.T) {
	validTypes := []RuleType{RuleTypeNaming, RuleTypeCommit, RuleTypeReview, RuleTypeTest}
	for _, rt := range validTypes {
		if !rt.IsValid() {
			t.Errorf("expected %s to be valid", rt)
		}
	}
	invalid := RuleType("INVALID")
	if invalid.IsValid() {
		t.Error("expected INVALID to be invalid")
	}
}

func TestNewDevStandard(t *testing.T) {
	ds := NewDevStandard("ds-1", "proj-1", "Naming Rule", "desc", RuleTypeNaming, "snake_case", true)
	if ds.ID != "ds-1" {
		t.Errorf("expected ID ds-1, got %s", ds.ID)
	}
	if ds.ProjectID != "proj-1" {
		t.Errorf("expected ProjectID proj-1, got %s", ds.ProjectID)
	}
	if ds.RuleType != RuleTypeNaming {
		t.Errorf("expected NAMING, got %s", ds.RuleType)
	}
	if !ds.Enabled {
		t.Error("expected Enabled true")
	}
}

func TestDevStandard_Validate(t *testing.T) {
	ds := NewDevStandard("ds-1", "proj-1", "Naming Rule", "desc", RuleTypeNaming, "snake_case", true)
	if err := ds.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	noProject := NewDevStandard("ds-2", "", "Name", "", RuleTypeNaming, "rule", true)
	if err := noProject.Validate(); err == nil {
		t.Error("expected error for empty project_id")
	}

	noName := NewDevStandard("ds-3", "proj-1", "", "", RuleTypeNaming, "rule", true)
	if err := noName.Validate(); err == nil {
		t.Error("expected error for empty name")
	}

	invalidType := NewDevStandard("ds-4", "proj-1", "Name", "", RuleType("INVALID"), "rule", true)
	if err := invalidType.Validate(); err == nil {
		t.Error("expected error for invalid rule_type")
	}

	noRule := NewDevStandard("ds-5", "proj-1", "Name", "", RuleTypeNaming, "", true)
	if err := noRule.Validate(); err == nil {
		t.Error("expected error for empty rule")
	}
}

func TestNewBranchPolicyConfig(t *testing.T) {
	bpc := NewBranchPolicyConfig("bpc-1", "proj-1", "feature/*", "feature branch", true)
	if bpc.ID != "bpc-1" {
		t.Errorf("expected ID bpc-1, got %s", bpc.ID)
	}
	if bpc.Pattern != "feature/*" {
		t.Errorf("expected feature/*, got %s", bpc.Pattern)
	}
	if !bpc.IsDefault {
		t.Error("expected IsDefault true")
	}
}

func TestBranchPolicyConfig_Validate(t *testing.T) {
	bpc := NewBranchPolicyConfig("bpc-1", "proj-1", "feature/*", "desc", true)
	if err := bpc.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	noProject := NewBranchPolicyConfig("bpc-2", "", "feature/*", "desc", false)
	if err := noProject.Validate(); err == nil {
		t.Error("expected error for empty project_id")
	}

	noPattern := NewBranchPolicyConfig("bpc-3", "proj-1", "", "desc", false)
	if err := noPattern.Validate(); err == nil {
		t.Error("expected error for empty pattern")
	}
}
