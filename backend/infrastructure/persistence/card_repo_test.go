package persistence

import (
	"testing"

	"github.com/mengri/nbcoder/domain/requirement"
)

func TestInMemoryCardRepo_Save(t *testing.T) {
	repo := NewInMemoryCardRepo()
	card := requirement.NewCard("id-1", "Test", "desc", "orig", "proj-1")
	if err := repo.Save(card); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestInMemoryCardRepo_FindByID(t *testing.T) {
	repo := NewInMemoryCardRepo()
	card := requirement.NewCard("id-1", "Test", "desc", "orig", "proj-1")
	_ = repo.Save(card)

	found, err := repo.FindByID("id-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if found == nil {
		t.Fatal("expected card to be found")
	}
	if found.ID != "id-1" {
		t.Errorf("expected ID id-1, got %s", found.ID)
	}
}

func TestInMemoryCardRepo_FindByID_NotFound(t *testing.T) {
	repo := NewInMemoryCardRepo()
	found, err := repo.FindByID("nonexistent")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if found != nil {
		t.Error("expected nil for nonexistent card")
	}
}

func TestInMemoryCardRepo_FindByProjectID(t *testing.T) {
	repo := NewInMemoryCardRepo()
	_ = repo.Save(requirement.NewCard("id-1", "Card1", "desc", "orig", "proj-1"))
	_ = repo.Save(requirement.NewCard("id-2", "Card2", "desc", "orig", "proj-1"))
	_ = repo.Save(requirement.NewCard("id-3", "Card3", "desc", "orig", "proj-2"))

	cards, err := repo.FindByProjectID("proj-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(cards) != 2 {
		t.Errorf("expected 2 cards, got %d", len(cards))
	}
}

func TestInMemoryCardRepo_FindByStatus(t *testing.T) {
	repo := NewInMemoryCardRepo()
	card1 := requirement.NewCard("id-1", "Card1", "desc", "orig", "proj-1")
	_ = repo.Save(card1)
	card2 := requirement.NewCard("id-2", "Card2", "desc", "orig", "proj-1")
	_ = card2.Confirm()
	_ = repo.Save(card2)

	cards, err := repo.FindByStatus(requirement.CardConfirmed)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(cards) != 1 {
		t.Errorf("expected 1 card, got %d", len(cards))
	}
}

func TestInMemoryCardRepo_FindAll(t *testing.T) {
	repo := NewInMemoryCardRepo()
	_ = repo.Save(requirement.NewCard("id-1", "Card1", "desc", "orig", "proj-1"))
	_ = repo.Save(requirement.NewCard("id-2", "Card2", "desc", "orig", "proj-2"))

	cards, err := repo.FindAll()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(cards) != 2 {
		t.Errorf("expected 2 cards, got %d", len(cards))
	}
}

func TestInMemoryCardRepo_Update(t *testing.T) {
	repo := NewInMemoryCardRepo()
	card := requirement.NewCard("id-1", "Test", "desc", "orig", "proj-1")
	_ = repo.Save(card)

	card.Title = "Updated"
	_ = repo.Update(card)

	found, _ := repo.FindByID("id-1")
	if found.Title != "Updated" {
		t.Errorf("expected title Updated, got %s", found.Title)
	}
}

func TestInMemoryCardRepo_Delete(t *testing.T) {
	repo := NewInMemoryCardRepo()
	card := requirement.NewCard("id-1", "Test", "desc", "orig", "proj-1")
	_ = repo.Save(card)

	if err := repo.Delete("id-1"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	found, _ := repo.FindByID("id-1")
	if found != nil {
		t.Error("expected card to be deleted")
	}
}

func TestInMemoryCardDependencyRepo(t *testing.T) {
	repo := NewInMemoryCardDependencyRepo()
	dep := requirement.NewCardDependency("dep-1", "card-1", "card-2", requirement.DependencyDependsOn)

	if err := repo.Save(dep); err != nil {
		t.Errorf("unexpected error saving: %v", err)
	}

	deps, err := repo.FindByCardID("card-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(deps) != 1 {
		t.Errorf("expected 1 dependency, got %d", len(deps))
	}

	found, err := repo.FindByID("dep-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if found == nil || found.ID != "dep-1" {
		t.Error("expected dependency to be found by ID")
	}

	depsByDep, err := repo.FindByDependsOnID("card-2")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(depsByDep) != 1 {
		t.Errorf("expected 1 dependency, got %d", len(depsByDep))
	}

	if err := repo.DeleteByCardID("card-1"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	deps, _ = repo.FindByCardID("card-1")
	if len(deps) != 0 {
		t.Errorf("expected 0 dependencies after delete, got %d", len(deps))
	}
}
