package requirement

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

type MockCardRepo struct {
	cards map[string]*Card
}

func NewMockCardRepo() *MockCardRepo {
	return &MockCardRepo{
		cards: make(map[string]*Card),
	}
}

func (m *MockCardRepo) Save(card *Card) error {
	m.cards[card.ID] = card
	return nil
}

func (m *MockCardRepo) FindByID(id string) (*Card, error) {
	card, ok := m.cards[id]
	if !ok {
		return nil, errors.New("card not found")
	}
	return card, nil
}

func (m *MockCardRepo) FindByProjectID(projectID string) ([]*Card, error) {
	var cards []*Card
	for _, card := range m.cards {
		if card.ProjectID == projectID {
			cards = append(cards, card)
		}
	}
	return cards, nil
}

func (m *MockCardRepo) FindByStatus(status CardStatus) ([]*Card, error) {
	var cards []*Card
	for _, card := range m.cards {
		if card.Status == status {
			cards = append(cards, card)
		}
	}
	return cards, nil
}

func (m *MockCardRepo) FindAll() ([]*Card, error) {
	var cards []*Card
	for _, card := range m.cards {
		cards = append(cards, card)
	}
	return cards, nil
}

func (m *MockCardRepo) Update(card *Card) error {
	m.cards[card.ID] = card
	return nil
}

func (m *MockCardRepo) Delete(id string) error {
	delete(m.cards, id)
	return nil
}

func setupTestCards() (*MockCardRepo, []string) {
	repo := NewMockCardRepo()
	cardIDs := []string{}

	now := time.Now().UTC()

	for i := 1; i <= 5; i++ {
		id := fmt.Sprintf("card-%d", i)
		cardIDs = append(cardIDs, id)

		var status CardStatus
		var priority Priority
		switch i {
		case 1:
			status = CardDraft
			priority = PriorityMedium
		case 2:
			status = CardConfirmed
			priority = PriorityHigh
		case 3:
			status = CardInProgress
			priority = PriorityLow
		case 4:
			status = CardCompleted
			priority = PriorityCritical
		case 5:
			status = CardDraft
			priority = PriorityHigh
		}

		card := &Card{
			ID:        id,
			Title:     fmt.Sprintf("Test Card %d", i),
			Status:    status,
			Priority:  priority,
			ProjectID: "proj-1",
			CreatedAt: now.Add(time.Duration(-i) * time.Hour),
			UpdatedAt: now.Add(time.Duration(-i) * time.Hour),
		}
		repo.Save(card)
	}

	return repo, cardIDs
}

func TestNewCardPool(t *testing.T) {
	repo := NewMockCardRepo()
	pool := NewCardPool(repo)

	if pool == nil {
		t.Error("expected card pool to be created")
	}
}

func TestCardPool_QueryCards_NoFilter(t *testing.T) {
	repo, cardIDs := setupTestCards()
	pool := NewCardPool(repo)

	query := &CardPoolQuery{}
	view, err := pool.QueryCards(query)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if view.TotalCount != len(cardIDs) {
		t.Errorf("expected total count %d, got %d", len(cardIDs), view.TotalCount)
	}

	if len(view.Cards) != len(cardIDs) {
		t.Errorf("expected %d cards, got %d", len(cardIDs), len(view.Cards))
	}
}

func TestCardPool_QueryCards_FilterByStatus(t *testing.T) {
	repo, _ := setupTestCards()
	pool := NewCardPool(repo)

	query := &CardPoolQuery{
		Filter: &CardFilter{
			Statuses: []CardStatus{CardDraft},
		},
	}

	view, err := pool.QueryCards(query)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expectedCount := 2
	if view.TotalCount != expectedCount {
		t.Errorf("expected total count %d, got %d", expectedCount, view.TotalCount)
	}

	for _, card := range view.Cards {
		if card.Status != CardDraft {
			t.Errorf("expected card status DRAFT, got %s", card.Status)
		}
	}
}

func TestCardPool_QueryCards_FilterByPriority(t *testing.T) {
	repo, _ := setupTestCards()
	pool := NewCardPool(repo)

	highPriority := PriorityHigh
	query := &CardPoolQuery{
		Filter: &CardFilter{
			Priority: &highPriority,
		},
	}

	view, err := pool.QueryCards(query)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expectedCount := 2
	if view.TotalCount != expectedCount {
		t.Errorf("expected total count %d, got %d", expectedCount, view.TotalCount)
	}

	for _, card := range view.Cards {
		if card.Priority != PriorityHigh {
			t.Errorf("expected card priority HIGH, got %s", card.Priority)
		}
	}
}

func TestCardPool_QueryCards_FilterByProjectID(t *testing.T) {
	repo, _ := setupTestCards()

	now := time.Now().UTC()
	project2Card := &Card{
		ID:        "card-6",
		Title:     "Project 2 Card",
		Status:    CardDraft,
		Priority:  PriorityMedium,
		ProjectID: "proj-2",
		CreatedAt: now,
		UpdatedAt: now,
	}
	repo.Save(project2Card)

	pool := NewCardPool(repo)

	query := &CardPoolQuery{
		Filter: &CardFilter{
			ProjectID: "proj-1",
		},
	}

	view, err := pool.QueryCards(query)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expectedCount := 5
	if view.TotalCount != expectedCount {
		t.Errorf("expected total count %d, got %d", expectedCount, view.TotalCount)
	}

	for _, card := range view.Cards {
		if card.ProjectID != "proj-1" {
			t.Errorf("expected project ID proj-1, got %s", card.ProjectID)
		}
	}
}

func TestCardPool_QueryCards_SortByPriorityAsc(t *testing.T) {
	repo, _ := setupTestCards()
	pool := NewCardPool(repo)

	query := &CardPoolQuery{
		SortField: SortByPriority,
		SortOrder: SortAsc,
	}

	view, err := pool.QueryCards(query)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(view.Cards) < 2 {
		t.Fatal("expected at least 2 cards for sorting test")
	}

	expected := []Priority{PriorityLow, PriorityMedium, PriorityHigh, PriorityHigh, PriorityCritical}
	for i, card := range view.Cards {
		if card.Priority != expected[i] {
			t.Errorf("expected priority %s at position %d, got %s", expected[i], i, card.Priority)
		}
	}
}

func TestCardPool_QueryCards_SortByPriorityDesc(t *testing.T) {
	repo, _ := setupTestCards()
	pool := NewCardPool(repo)

	query := &CardPoolQuery{
		SortField: SortByPriority,
		SortOrder: SortDesc,
	}

	view, err := pool.QueryCards(query)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(view.Cards) < 2 {
		t.Fatal("expected at least 2 cards for sorting test")
	}

	expected := []Priority{PriorityCritical, PriorityHigh, PriorityHigh, PriorityMedium, PriorityLow}
	for i, card := range view.Cards {
		if card.Priority != expected[i] {
			t.Errorf("expected priority %s at position %d, got %s", expected[i], i, card.Priority)
		}
	}
}

func TestCardPool_QueryCards_SortByUpdatedAtDesc(t *testing.T) {
	repo, _ := setupTestCards()
	pool := NewCardPool(repo)

	query := &CardPoolQuery{
		SortField: SortByUpdatedAt,
		SortOrder: SortDesc,
	}

	view, err := pool.QueryCards(query)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(view.Cards) < 2 {
		t.Fatal("expected at least 2 cards for sorting test")
	}

	for i := 1; i < len(view.Cards); i++ {
		if view.Cards[i-1].UpdatedAt.Before(view.Cards[i].UpdatedAt) {
			t.Errorf("expected cards to be sorted by updated_at desc")
		}
	}
}

func TestCardPool_QueryCards_Pagination(t *testing.T) {
	repo, _ := setupTestCards()
	pool := NewCardPool(repo)

	query := &CardPoolQuery{
		Limit:  2,
		Offset: 1,
	}

	view, err := pool.QueryCards(query)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(view.Cards) != 2 {
		t.Errorf("expected 2 cards with limit, got %d", len(view.Cards))
	}

	if view.TotalCount != 5 {
		t.Errorf("expected total count 5, got %d", view.TotalCount)
	}
}

func TestCardPool_BatchOperation_BatchConfirm(t *testing.T) {
	repo, _ := setupTestCards()
	pool := NewCardPool(repo)

	cardIDs := []string{"card-1", "card-5"}
	result, err := pool.BatchOperation(cardIDs, BatchConfirm, nil)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if result.SuccessCount != 2 {
		t.Errorf("expected 2 successful operations, got %d", result.SuccessCount)
	}

	if result.FailureCount != 0 {
		t.Errorf("expected 0 failed operations, got %d", result.FailureCount)
	}

	for _, cardID := range cardIDs {
		card, _ := repo.FindByID(cardID)
		if card.Status != CardConfirmed {
			t.Errorf("expected card %s to be confirmed, got status %s", cardID, card.Status)
		}
	}
}

func TestCardPool_BatchOperation_BatchAbandon(t *testing.T) {
	repo, _ := setupTestCards()
	pool := NewCardPool(repo)

	cardIDs := []string{"card-1", "card-3"}
	result, err := pool.BatchOperation(cardIDs, BatchAbandon, nil)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if result.SuccessCount != 2 {
		t.Errorf("expected 2 successful operations, got %d", result.SuccessCount)
	}

	if result.FailureCount != 0 {
		t.Errorf("expected 0 failed operations, got %d", result.FailureCount)
	}

	for _, cardID := range cardIDs {
		card, _ := repo.FindByID(cardID)
		if card.Status != CardAbandoned {
			t.Errorf("expected card %s to be abandoned, got status %s", cardID, card.Status)
		}
	}
}

func TestCardPool_BatchOperation_BatchSupersede(t *testing.T) {
	repo, _ := setupTestCards()
	pool := NewCardPool(repo)

	newCardID := "card-new"
	cardIDs := []string{"card-1", "card-3"}
	params := map[string]interface{}{
		"new_card_id": newCardID,
	}

	result, err := pool.BatchOperation(cardIDs, BatchSupersede, params)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if result.SuccessCount != 2 {
		t.Errorf("expected 2 successful operations, got %d", result.SuccessCount)
	}

	if result.FailureCount != 0 {
		t.Errorf("expected 0 failed operations, got %d", result.FailureCount)
	}

	for _, cardID := range cardIDs {
		card, _ := repo.FindByID(cardID)
		if card.Status != CardSuperseded {
			t.Errorf("expected card %s to be superseded, got status %s", cardID, card.Status)
		}
		if card.SupersededBy != newCardID {
			t.Errorf("expected card %s superseded by %s, got %s", cardID, newCardID, card.SupersededBy)
		}
	}
}

func TestCardPool_BatchOperation_InvalidCardIDs(t *testing.T) {
	repo, _ := setupTestCards()
	pool := NewCardPool(repo)

	cardIDs := []string{"invalid-card-1", "invalid-card-2"}
	result, err := pool.BatchOperation(cardIDs, BatchConfirm, nil)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if result.SuccessCount != 0 {
		t.Errorf("expected 0 successful operations, got %d", result.SuccessCount)
	}

	if result.FailureCount != 2 {
		t.Errorf("expected 2 failed operations, got %d", result.FailureCount)
	}

	if len(result.Errors) != 2 {
		t.Errorf("expected 2 errors, got %d", len(result.Errors))
	}
}

func TestCardPool_BatchOperation_EmptyCardIDs(t *testing.T) {
	repo := NewMockCardRepo()
	pool := NewCardPool(repo)

	cardIDs := []string{}
	_, err := pool.BatchOperation(cardIDs, BatchConfirm, nil)

	if err == nil {
		t.Error("expected error for empty card IDs")
	}
}

func TestCardPool_BatchOperation_UnknownOperation(t *testing.T) {
	repo, _ := setupTestCards()
	pool := NewCardPool(repo)

	cardIDs := []string{"card-1"}
	result, err := pool.BatchOperation(cardIDs, BatchOperation("unknown"), nil)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if result.SuccessCount != 0 {
		t.Errorf("expected 0 successful operations, got %d", result.SuccessCount)
	}

	if result.FailureCount != 1 {
		t.Errorf("expected 1 failed operation, got %d", result.FailureCount)
	}
}