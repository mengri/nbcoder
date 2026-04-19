package requirement

import (
	"testing"
)

func TestCardAggregate_AddDependency(t *testing.T) {
	card := NewCard("id-1", "Test", "desc", "orig", "proj-1")
	ca := NewCardAggregate(card)
	dep := NewCardDependency("dep-1", "id-1", "id-2", DependencyBlocks)
	result := ca.AddDependency(dep)
	if result.ID != "dep-1" {
		t.Errorf("expected dep-1, got %s", result.ID)
	}
	if len(ca.Dependencies) != 1 {
		t.Errorf("expected 1 dependency, got %d", len(ca.Dependencies))
	}
	if ca.ChangeType != DependencyChangeAdded {
		t.Errorf("expected change type ADDED, got %s", ca.ChangeType)
	}
}

func TestCardAggregate_RemoveDependency(t *testing.T) {
	card := NewCard("id-1", "Test", "desc", "orig", "proj-1")
	ca := NewCardAggregate(card)
	dep1 := NewCardDependency("dep-1", "id-1", "id-2", DependencyDependsOn)
	ca.AddDependency(dep1)
	dep2 := NewCardDependency("dep-2", "id-1", "id-3", DependencyBlocks)
	ca.AddDependency(dep2)

	err := ca.RemoveDependency("dep-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(ca.Dependencies) != 1 {
		t.Errorf("expected 1 dependency after removal, got %d", len(ca.Dependencies))
	}
	if ca.ChangeType != DependencyChangeRemoved {
		t.Errorf("expected change type REMOVED, got %s", ca.ChangeType)
	}
}

func TestCardAggregate_IsBlocked(t *testing.T) {
	card := NewCard("id-1", "Test", "desc", "orig", "proj-1")
	ca := NewCardAggregate(card)
	if ca.IsBlocked() {
		t.Error("expected card not to be blocked initially")
	}

	dep := NewCardDependency("dep-1", "id-1", "id-2", DependencyDependsOn)
	ca.AddDependency(dep)
	if !ca.IsBlocked() {
		t.Error("expected card to be blocked after adding DEPENDS_ON")
	}
}

func TestCardAggregate_IsDependencySatisfied(t *testing.T) {
	card := NewCard("id-1", "Test", "desc", "orig", "proj-1")
	ca := NewCardAggregate(card)
	if ca.IsDependencySatisfied() {
		t.Error("expected card to not have satisfied dependencies initially")
	}

	dep := NewCardDependency("dep-1", "id-1", "id-2", DependencyDependsOn)
	ca.AddDependency(dep)
	if ca.IsDependencySatisfied() {
		t.Error("expected card to still have dependencies unsatisfied")
	}

	dep2 := NewCardDependency("dep-2", "id-1", "id-3", DependencyBlocks)
	ca.AddDependency(dep2)
	if ca.IsDependencySatisfied() {
		t.Error("expected card to still have unsatisfied dependencies")
	}
}

func TestCardAggregate_GetBlockingDependencies(t *testing.T) {
	card := NewCard("id-1", "Test", "desc", "orig", "proj-1")
	ca := NewCardAggregate(card)
	dep1 := NewCardDependency("dep-1", "id-1", "id-2", DependencyDependsOn)
	ca.AddDependency(dep1)
	dep2 := NewCardDependency("dep-2", "id-1", "id-3", DependencyBlocks)
	ca.AddDependency(dep2)
	dep3 := NewCardDependency("dep-3", "id-1", "id-4", DependencyDependsOn)
	ca.AddDependency(dep3)

	blocking := ca.GetBlockingDependencies()
	if len(blocking) != 2 {
		t.Errorf("expected 2 blocking dependencies, got %d", len(blocking))
	}
	if blocking[0].DependsOnID != "id-2" && blocking[0].DependsOnID != "id-4" {
		t.Errorf("unexpected blocking dependency")
	}
	if blocking[1].DependsOnID != "id-2" && blocking[1].DependsOnID != "id-4" {
		t.Errorf("unexpected blocking dependency")
	}
}

func TestCardAggregate_CheckAutoUnblock(t *testing.T) {
	card := NewCard("id-1", "Test", "desc", "orig", "proj-1")
	ca := NewCardAggregate(card)
	if !ca.CheckAutoUnblock() {
		t.Error("expected card to be unblocked initially")
	}

	dep := NewCardDependency("dep-1", "id-1", "id-2", DependencyDependsOn)
	ca.AddDependency(dep)
	if ca.CheckAutoUnblock() {
		t.Error("expected card to still be blocked")
	}

	ca.RemoveDependency("dep-1")
	if !ca.CheckAutoUnblock() {
		t.Error("expected card to be unblocked after removing blocking dependency")
	}
}

func TestCardDependency_Validate(t *testing.T) {
	dep := NewCardDependency("dep-1", "card-1", "card-2", DependencyDependsOn)
	if err := dep.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	invalid := NewCardDependency("dep-2", "", "", DependencyDependsOn)
	if err := invalid.Validate(); err == nil {
		t.Error("expected error for empty card ID")
	}
}
