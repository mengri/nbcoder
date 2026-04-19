package requirement

import (
	"testing"
)

func TestNewCard(t *testing.T) {
	card := NewCard("id-1", "Test Card", "desc", "original input", "proj-1")
	if card.ID != "id-1" {
		t.Errorf("expected ID id-1, got %s", card.ID)
	}
	if card.Status != CardDraft {
		t.Errorf("expected status DRAFT, got %s", card.Status)
	}
	if card.Priority != PriorityMedium {
		t.Errorf("expected default priority MEDIUM, got %s", card.Priority)
	}
	if card.ProjectID != "proj-1" {
		t.Errorf("expected project proj-1, got %s", card.ProjectID)
	}
	if card.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestCard_SetPriority(t *testing.T) {
	card := NewCard("id-1", "Test", "desc", "orig", "proj-1")

	if err := card.SetPriority(PriorityHigh); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if card.Priority != PriorityHigh {
		t.Errorf("expected priority HIGH, got %s", card.Priority)
	}

	if err := card.SetPriority(Priority("INVALID")); err == nil {
		t.Error("expected error for invalid priority")
	}
}

func TestCard_SetStructuredOutput(t *testing.T) {
	card := NewCard("id-1", "Test", "desc", "orig", "proj-1")
	card.SetStructuredOutput("structured content")
	if card.StructuredOutput != "structured content" {
		t.Errorf("expected structured output, got %s", card.StructuredOutput)
	}
}

func TestCard_SetPipelineID(t *testing.T) {
	card := NewCard("id-1", "Test", "desc", "orig", "proj-1")
	card.SetPipelineID("pipeline-1")
	if card.PipelineID != "pipeline-1" {
		t.Errorf("expected pipeline-1, got %s", card.PipelineID)
	}
}

func TestCard_Confirm(t *testing.T) {
	card := NewCard("id-1", "Test", "desc", "orig", "proj-1")
	if err := card.Confirm(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if card.Status != CardConfirmed {
		t.Errorf("expected status CONFIRMED, got %s", card.Status)
	}
}

func TestCard_Confirm_InvalidStatus(t *testing.T) {
	card := NewCard("id-1", "Test", "desc", "orig", "proj-1")
	card.Status = CardCompleted
	if err := card.Confirm(); err == nil {
		t.Error("expected error confirming non-draft card")
	}
}

func TestCard_Start(t *testing.T) {
	card := NewCard("id-1", "Test", "desc", "orig", "proj-1")
	_ = card.Confirm()
	if err := card.Start(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if card.Status != CardInProgress {
		t.Errorf("expected status IN_PROGRESS, got %s", card.Status)
	}
}

func TestCard_Start_InvalidStatus(t *testing.T) {
	card := NewCard("id-1", "Test", "desc", "orig", "proj-1")
	if err := card.Start(); err == nil {
		t.Error("expected error starting non-confirmed card")
	}
}

func TestCard_Complete(t *testing.T) {
	card := NewCard("id-1", "Test", "desc", "orig", "proj-1")
	_ = card.Confirm()
	_ = card.Start()
	if err := card.Complete(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if card.Status != CardCompleted {
		t.Errorf("expected status COMPLETED, got %s", card.Status)
	}
}

func TestCard_Complete_InvalidStatus(t *testing.T) {
	card := NewCard("id-1", "Test", "desc", "orig", "proj-1")
	if err := card.Complete(); err == nil {
		t.Error("expected error completing non-in-progress card")
	}
}

func TestCard_Supersede(t *testing.T) {
	card := NewCard("id-1", "Test", "desc", "orig", "proj-1")
	if err := card.Supersede("id-2"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if card.Status != CardSuperseded {
		t.Errorf("expected status SUPERSEDED, got %s", card.Status)
	}
	if card.SupersededBy != "id-2" {
		t.Errorf("expected superseded_by id-2, got %s", card.SupersededBy)
	}
}

func TestCard_Supersede_InvalidStatus(t *testing.T) {
	card := NewCard("id-1", "Test", "desc", "orig", "proj-1")
	_ = card.Abandon()
	if err := card.Supersede("id-2"); err == nil {
		t.Error("expected error superseding abandoned card")
	}
}

func TestCard_Abandon(t *testing.T) {
	card := NewCard("id-1", "Test", "desc", "orig", "proj-1")
	if err := card.Abandon(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if card.Status != CardAbandoned {
		t.Errorf("expected status ABANDONED, got %s", card.Status)
	}
}

func TestCard_Abandon_InvalidStatus(t *testing.T) {
	card := NewCard("id-1", "Test", "desc", "orig", "proj-1")
	_ = card.Confirm()
	_ = card.Start()
	_ = card.Complete()
	if err := card.Abandon(); err == nil {
		t.Error("expected error abandoning completed card")
	}
}

func TestCardStatus_IsValid(t *testing.T) {
	statuses := []CardStatus{CardDraft, CardConfirmed, CardInProgress, CardCompleted, CardSuperseded, CardAbandoned}
	for _, s := range statuses {
		if !s.IsValid() {
			t.Errorf("expected status %s to be valid", s)
		}
	}
	invalid := CardStatus("INVALID")
	if invalid.IsValid() {
		t.Error("expected INVALID status to be invalid")
	}
}

func TestPriority_IsValid(t *testing.T) {
	priorities := []Priority{PriorityLow, PriorityMedium, PriorityHigh, PriorityCritical}
	for _, p := range priorities {
		if !p.IsValid() {
			t.Errorf("expected priority %s to be valid", p)
		}
	}
	invalid := Priority("INVALID")
	if invalid.IsValid() {
		t.Error("expected INVALID priority to be invalid")
	}
}

func TestCardDependency(t *testing.T) {
	dep := NewCardDependency("dep-1", "card-1", "card-2", DependencyDependsOn)
	if dep.ID != "dep-1" {
		t.Errorf("expected ID dep-1, got %s", dep.ID)
	}
	if dep.Type != DependencyDependsOn {
		t.Errorf("expected type DEPENDS_ON, got %s", dep.Type)
	}
}

func TestCardAggregate_IsBlocked(t *testing.T) {
	card := NewCard("id-1", "Test", "desc", "orig", "proj-1")
	aggregate := NewCardAggregate(card)

	if aggregate.IsBlocked() {
		t.Error("expected aggregate not to be blocked without dependencies")
	}

	dep := NewCardDependency("dep-1", "id-1", "id-2", DependencyDependsOn)
	aggregate.AddDependency(dep)
	if !aggregate.IsBlocked() {
		t.Error("expected aggregate to be blocked with DEPENDS_ON dependency")
	}
}

func TestCardAggregate_Confirm_Blocked(t *testing.T) {
	card := NewCard("id-1", "Test", "desc", "orig", "proj-1")
	aggregate := NewCardAggregate(card)
	dep := NewCardDependency("dep-1", "id-1", "id-2", DependencyDependsOn)
	aggregate.AddDependency(dep)

	if err := aggregate.Confirm(); err == nil {
		t.Error("expected error confirming blocked card")
	}
}

func TestCardAggregate_Confirm_Unblocked(t *testing.T) {
	card := NewCard("id-1", "Test", "desc", "orig", "proj-1")
	aggregate := NewCardAggregate(card)
	dep := NewCardDependency("dep-1", "id-1", "id-2", DependencyBlocks)
	aggregate.AddDependency(dep)

	if err := aggregate.Confirm(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if card.Status != CardConfirmed {
		t.Errorf("expected status CONFIRMED, got %s", card.Status)
	}
}
