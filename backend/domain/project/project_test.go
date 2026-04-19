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
