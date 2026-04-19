package project

import (
	"testing"
)

func TestNewProjectLifecycle(t *testing.T) {
	lc := NewProjectLifecycle("lc-1", "proj-1")
	if lc.ID != "lc-1" {
		t.Errorf("expected ID lc-1, got %s", lc.ID)
	}
	if lc.ProjectID != "proj-1" {
		t.Errorf("expected ProjectID proj-1, got %s", lc.ProjectID)
	}
	if lc.Status != LifecycleCreating {
		t.Errorf("expected status CREATING, got %s", lc.Status)
	}
	if lc.ActivatedAt != nil {
		t.Error("expected ActivatedAt to be nil")
	}
}

func TestProjectLifecycle_CanTransitionTo(t *testing.T) {
	lc := NewProjectLifecycle("lc-1", "proj-1")
	if !lc.CanTransitionTo(LifecycleActive) {
		t.Error("CREATING should transition to ACTIVE")
	}
	if !lc.CanTransitionTo(LifecycleDeleted) {
		t.Error("CREATING should transition to DELETED")
	}
	if lc.CanTransitionTo(LifecycleSuspended) {
		t.Error("CREATING should not transition to SUSPENDED")
	}
}

func TestProjectLifecycle_Activate(t *testing.T) {
	lc := NewProjectLifecycle("lc-1", "proj-1")
	if err := lc.Activate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if lc.Status != LifecycleActive {
		t.Errorf("expected ACTIVE, got %s", lc.Status)
	}
	if lc.ActivatedAt == nil {
		t.Error("expected ActivatedAt to be set")
	}
	if err := lc.Activate(); err == nil {
		t.Error("expected error activating already active lifecycle")
	}
}

func TestProjectLifecycle_Suspend(t *testing.T) {
	lc := NewProjectLifecycle("lc-1", "proj-1")
	if err := lc.Suspend(); err == nil {
		t.Error("expected error suspending from CREATING")
	}
	_ = lc.Activate()
	if err := lc.Suspend(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if lc.Status != LifecycleSuspended {
		t.Errorf("expected SUSPENDED, got %s", lc.Status)
	}
	if lc.SuspendedAt == nil {
		t.Error("expected SuspendedAt to be set")
	}
}

func TestProjectLifecycle_Archive(t *testing.T) {
	lc := NewProjectLifecycle("lc-1", "proj-1")
	_ = lc.Activate()
	if err := lc.Archive(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if lc.Status != LifecycleArchived {
		t.Errorf("expected ARCHIVED, got %s", lc.Status)
	}
	if lc.ArchivedAt == nil {
		t.Error("expected ArchivedAt to be set")
	}
}

func TestProjectLifecycle_Delete(t *testing.T) {
	lc := NewProjectLifecycle("lc-1", "proj-1")
	_ = lc.Activate()
	if err := lc.Delete(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if lc.Status != LifecycleDeleted {
		t.Errorf("expected DELETED, got %s", lc.Status)
	}
}

func TestProjectLifecycle_FullLifecycle(t *testing.T) {
	lc := NewProjectLifecycle("lc-1", "proj-1")
	_ = lc.Activate()
	_ = lc.Suspend()
	_ = lc.Activate()
	_ = lc.Archive()
	_ = lc.Delete()
	if lc.Status != LifecycleDeleted {
		t.Errorf("expected DELETED, got %s", lc.Status)
	}
}

func TestProjectLifecycle_ArchiveFromSuspended(t *testing.T) {
	lc := NewProjectLifecycle("lc-1", "proj-1")
	_ = lc.Activate()
	_ = lc.Suspend()
	if err := lc.Archive(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if lc.Status != LifecycleArchived {
		t.Errorf("expected ARCHIVED, got %s", lc.Status)
	}
}

func TestProjectLifecycle_DeletedCannotTransition(t *testing.T) {
	lc := NewProjectLifecycle("lc-1", "proj-1")
	_ = lc.Activate()
	_ = lc.Delete()
	if lc.CanTransitionTo(LifecycleActive) {
		t.Error("DELETED should not transition to any state")
	}
	if err := lc.Activate(); err == nil {
		t.Error("expected error activating deleted lifecycle")
	}
}

func TestProjectLifecycle_CreatingCannotSuspendOrArchive(t *testing.T) {
	lc := NewProjectLifecycle("lc-1", "proj-1")
	if err := lc.Suspend(); err == nil {
		t.Error("expected error suspending from CREATING")
	}
	if err := lc.Archive(); err == nil {
		t.Error("expected error archiving from CREATING")
	}
}
