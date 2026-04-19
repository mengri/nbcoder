package project

import (
	"testing"
)

func TestNewProject(t *testing.T) {
	p := NewProject("id-1", "Test Project", "description", "https://github.com/test/repo")
	if p.ID != "id-1" {
		t.Errorf("expected ID id-1, got %s", p.ID)
	}
	if p.Status != ProjectActive {
		t.Errorf("expected status ACTIVE, got %s", p.Status)
	}
}

func TestProject_Validate(t *testing.T) {
	p := NewProject("id-1", "Test", "desc", "")
	if err := p.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	noName := NewProject("id-2", "", "desc", "")
	if err := noName.Validate(); err == nil {
		t.Error("expected error for empty name")
	}
}

func TestProject_Archive(t *testing.T) {
	p := NewProject("id-1", "Test", "desc", "")
	if err := p.Archive(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if p.Status != ProjectArchived {
		t.Errorf("expected status ARCHIVED, got %s", p.Status)
	}
	if err := p.Archive(); err == nil {
		t.Error("expected error archiving already archived project")
	}
}

func TestProject_Activate(t *testing.T) {
	p := NewProject("id-1", "Test", "desc", "")
	_ = p.Archive()
	if err := p.Activate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if p.Status != ProjectActive {
		t.Errorf("expected status ACTIVE, got %s", p.Status)
	}
	if err := p.Activate(); err == nil {
		t.Error("expected error activating already active project")
	}
}

func TestProject_Update(t *testing.T) {
	p := NewProject("id-1", "Test", "desc", "")
	p.Update("New Name", "new desc", "https://github.com/new/repo")
	if p.Name != "New Name" {
		t.Errorf("expected name New Name, got %s", p.Name)
	}
}

func TestNewProjectConfig(t *testing.T) {
	cfg := NewProjectConfig("cfg-1", "proj-1", "key", "value")
	if cfg.Key != "key" {
		t.Errorf("expected key, got %s", cfg.Key)
	}
}

func TestNewStandards(t *testing.T) {
	std := NewStandards("std-1", "proj-1", "gitflow", "go", "snake_case")
	if std.BranchStrategy != "gitflow" {
		t.Errorf("expected gitflow, got %s", std.BranchStrategy)
	}
}

func TestStandards_Update(t *testing.T) {
	std := NewStandards("std-1", "proj-1", "gitflow", "go", "snake_case")
	std.Update("trunk", "", "")
	if std.BranchStrategy != "trunk" {
		t.Errorf("expected trunk, got %s", std.BranchStrategy)
	}
	if std.TechStack != "go" {
		t.Errorf("expected go, got %s", std.TechStack)
	}
}

func TestDefaultProjectDirectory(t *testing.T) {
	dir := DefaultProjectDirectory("proj-1")
	if len(dir.Dirs) == 0 {
		t.Error("expected default directories")
	}
	hasNBCoder := false
	for _, d := range dir.Dirs {
		if d == ".NBCoder" {
			hasNBCoder = true
		}
	}
	if !hasNBCoder {
		t.Error("expected .NBCoder directory")
	}
}

func TestProjectConfig_Validate(t *testing.T) {
	cfg := NewProjectConfig("cfg-1", "proj-1", "key1", "value1")
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	noKey := NewProjectConfig("cfg-2", "proj-1", "", "value1")
	if err := noKey.Validate(); err == nil {
		t.Error("expected error for empty key")
	}
}

func TestProjectConfig_Update(t *testing.T) {
	cfg := NewProjectConfig("cfg-1", "proj-1", "key1", "old")
	oldUpdatedAt := cfg.UpdatedAt
	cfg.Update("new")
	if cfg.Value != "new" {
		t.Errorf("expected value new, got %s", cfg.Value)
	}
	if !cfg.UpdatedAt.After(oldUpdatedAt) {
		t.Error("expected UpdatedAt to be after old value")
	}
}

func TestNewConfigChangeLog(t *testing.T) {
	log := NewConfigChangeLog("log-1", "proj-1", "key1", "old", "new", "user1")
	if log.ID != "log-1" {
		t.Errorf("expected ID log-1, got %s", log.ID)
	}
	if log.ProjectID != "proj-1" {
		t.Errorf("expected ProjectID proj-1, got %s", log.ProjectID)
	}
	if log.ConfigKey != "key1" {
		t.Errorf("expected ConfigKey key1, got %s", log.ConfigKey)
	}
	if log.OldValue != "old" {
		t.Errorf("expected OldValue old, got %s", log.OldValue)
	}
	if log.NewValue != "new" {
		t.Errorf("expected NewValue new, got %s", log.NewValue)
	}
	if log.ChangedBy != "user1" {
		t.Errorf("expected ChangedBy user1, got %s", log.ChangedBy)
	}
	if log.ChangedAt.IsZero() {
		t.Error("expected ChangedAt to be set")
	}
}
