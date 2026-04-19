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

func TestCardAggregate_HasBlockingDependencies(t *testing.T) {
	card := NewCard("id-1", "Test", "desc", "orig", "proj-1")
	ca := NewCardAggregate(card)
	if ca.HasBlockingDependencies() {
		t.Error("expected card to not have blocking dependencies initially")
	}

	dep := NewCardDependency("dep-1", "id-1", "id-2", DependencyDependsOn)
	ca.AddDependency(dep)
	if !ca.HasBlockingDependencies() {
		t.Error("expected card to have blocking dependencies after adding DEPENDS_ON")
	}

	ca.RemoveDependency("dep-1")
	if ca.HasBlockingDependencies() {
		t.Error("expected card to not have blocking dependencies after removing DEPENDS_ON")
	}

	dep2 := NewCardDependency("dep-2", "id-1", "id-3", DependencyBlocks)
	ca.AddDependency(dep2)
	if ca.HasBlockingDependencies() {
		t.Error("expected card to not have blocking dependencies with BLOCKS dependency")
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

func TestCardAggregate_Supersede(t *testing.T) {
	card := NewCard("id-1", "Test Card", "desc", "orig", "proj-1")
	ca := NewCardAggregate(card)

	if err := ca.Supersede("id-2"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !ca.IsSuperseded() {
		t.Error("expected card to be superseded")
	}

	if ca.GetSupersededBy() != "id-2" {
		t.Errorf("expected superseded by id-2, got %s", ca.GetSupersededBy())
	}
}

func TestCardAggregate_Supersede_InvalidStatus(t *testing.T) {
	card := NewCard("id-1", "Test Card", "desc", "orig", "proj-1")
	ca := NewCardAggregate(card)
	ca.Abandon()

	if err := ca.Supersede("id-2"); err == nil {
		t.Error("expected error superseding abandoned card")
	}
}

func TestCardAggregate_Abandon(t *testing.T) {
	card := NewCard("id-1", "Test Card", "desc", "orig", "proj-1")
	ca := NewCardAggregate(card)

	if err := ca.Abandon(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !ca.IsAbandoned() {
		t.Error("expected card to be abandoned")
	}
}

func TestCardAggregate_Abandon_InvalidStatus(t *testing.T) {
	card := NewCard("id-1", "Test Card", "desc", "orig", "proj-1")
	ca := NewCardAggregate(card)
	card.Confirm()
	card.Start()
	card.Complete()

	if err := ca.Abandon(); err == nil {
		t.Error("expected error abandoning completed card")
	}
}

func TestCardAggregate_ClearDependencies(t *testing.T) {
	card := NewCard("id-1", "Test Card", "desc", "orig", "proj-1")
	ca := NewCardAggregate(card)

	dep := NewCardDependency("dep-1", "id-1", "id-2", DependencyDependsOn)
	ca.AddDependency(dep)

	if len(ca.Dependencies) != 1 {
		t.Error("expected 1 dependency")
	}

	ca.ClearDependencies()

	if len(ca.Dependencies) != 0 {
		t.Error("expected no dependencies after clearing")
	}

	if ca.ChangeType != DependencyChangeRemoved {
		t.Error("expected change type REMOVED")
	}
}

func TestCardAggregate_HasDependencies(t *testing.T) {
	card := NewCard("id-1", "Test Card", "desc", "orig", "proj-1")
	ca := NewCardAggregate(card)

	if ca.HasDependencies() {
		t.Error("expected no dependencies initially")
	}

	dep := NewCardDependency("dep-1", "id-1", "id-2", DependencyDependsOn)
	ca.AddDependency(dep)

	if !ca.HasDependencies() {
		t.Error("expected to have dependencies")
	}
}

func TestCardAggregate_GetSupersededBy(t *testing.T) {
	card := NewCard("id-1", "Test Card", "desc", "orig", "proj-1")
	ca := NewCardAggregate(card)

	ca.Supersede("id-2")

	if ca.GetSupersededBy() != "id-2" {
		t.Errorf("expected superseded by id-2, got %s", ca.GetSupersededBy())
	}
}

func TestCardAggregate_IsSuperseded(t *testing.T) {
	card := NewCard("id-1", "Test Card", "desc", "orig", "proj-1")
	ca := NewCardAggregate(card)

	if ca.IsSuperseded() {
		t.Error("expected card not to be superseded initially")
	}

	ca.Supersede("id-2")

	if !ca.IsSuperseded() {
		t.Error("expected card to be superseded")
	}
}

func TestCardAggregate_IsAbandoned(t *testing.T) {
	card := NewCard("id-1", "Test Card", "desc", "orig", "proj-1")
	ca := NewCardAggregate(card)

	if ca.IsAbandoned() {
		t.Error("expected card not to be abandoned initially")
	}

	ca.Abandon()

	if !ca.IsAbandoned() {
		t.Error("expected card to be abandoned")
	}
}

func TestCardAggregate_GetLineageInfo(t *testing.T) {
	card := NewCard("id-1", "Test Card", "desc", "orig", "proj-1")
	ca := NewCardAggregate(card)

	dep := NewCardDependency("dep-1", "id-1", "id-2", DependencyDependsOn)
	ca.AddDependency(dep)

	ca.Supersede("id-3")

	lineage := ca.GetLineageInfo()

	if lineage["card_id"] != "id-1" {
		t.Errorf("expected card_id id-1, got %v", lineage["card_id"])
	}

	if lineage["title"] != "Test Card" {
		t.Errorf("expected title Test Card, got %v", lineage["title"])
	}

	if lineage["is_superseded"] != true {
		t.Error("expected is_superseded to be true")
	}

	if lineage["superseded_by"] != "id-3" {
		t.Errorf("expected superseded_by id-3, got %v", lineage["superseded_by"])
	}

	if lineage["has_dependencies"] != true {
		t.Error("expected has_dependencies to be true")
	}

	if lineage["dependency_count"] != 1 {
		t.Errorf("expected dependency_count 1, got %v", lineage["dependency_count"])
	}
}
