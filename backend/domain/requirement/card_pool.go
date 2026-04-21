package requirement

import (
	"errors"
	"fmt"
	"sort"
	"time"
)

type SortField string

const (
	SortByPriority  SortField = "priority"
	SortByCreatedAt SortField = "created_at"
	SortByUpdatedAt SortField = "updated_at"
)

type SortOrder string

const (
	SortAsc  SortOrder = "ASC"
	SortDesc SortOrder = "DESC"
)

type CardFilter struct {
	Statuses    []CardStatus
	Priority    *Priority
	ProjectName string
}

type CardPoolQuery struct {
	Filter    *CardFilter
	SortField SortField
	SortOrder SortOrder
	Limit     int
	Offset    int
}

type CardPoolView struct {
	Cards      []*Card
	TotalCount int
	Filter     *CardFilter
	SortField  SortField
	SortOrder  SortOrder
}

type BatchOperation string

const (
	BatchConfirm   BatchOperation = "confirm"
	BatchAbandon   BatchOperation = "abandon"
	BatchSupersede BatchOperation = "supersede"
)

type BatchOperationResult struct {
	SuccessCount int
	FailureCount int
	Errors       []error
}

type CardPool struct {
	cardRepo CardRepo
}

func NewCardPool(cardRepo CardRepo) *CardPool {
	return &CardPool{
		cardRepo: cardRepo,
	}
}

func (cp *CardPool) QueryCards(query *CardPoolQuery) (*CardPoolView, error) {
	if query == nil {
		query = &CardPoolQuery{}
	}

	var cards []*Card

	if query.Filter != nil {
		if len(query.Filter.Statuses) > 0 {
			var allCards []*Card
			for _, status := range query.Filter.Statuses {
				statusCards, err := cp.cardRepo.FindByStatus(status)
				if err != nil {
					return nil, fmt.Errorf("failed to find cards by status %s: %w", status, err)
				}
				allCards = append(allCards, statusCards...)
			}
			cards = cp.filterByProjectName(allCards, query.Filter.ProjectName)
		} else if query.Filter.ProjectName != "" {
			projectCards, err := cp.cardRepo.FindByProjectName(query.Filter.ProjectName)
			if err != nil {
				return nil, fmt.Errorf("failed to find cards by project %s: %w", query.Filter.ProjectName, err)
			}
			cards = projectCards
		} else {
			allCards, err := cp.cardRepo.FindAll()
			if err != nil {
				return nil, fmt.Errorf("failed to find all cards: %w", err)
			}
			cards = allCards
		}

		if query.Filter.Priority != nil {
			cards = cp.filterByPriority(cards, *query.Filter.Priority)
		}
	} else {
		allCards, err := cp.cardRepo.FindAll()
		if err != nil {
			return nil, fmt.Errorf("failed to find all cards: %w", err)
		}
		cards = allCards
	}

	totalCount := len(cards)

	cards = cp.sortCards(cards, query.SortField, query.SortOrder)

	if query.Offset > 0 {
		if query.Offset >= len(cards) {
			return &CardPoolView{
				Cards:      []*Card{},
				TotalCount: totalCount,
				Filter:     query.Filter,
				SortField:  query.SortField,
				SortOrder:  query.SortOrder,
			}, nil
		}
		cards = cards[query.Offset:]
	}

	if query.Limit > 0 && query.Limit < len(cards) {
		cards = cards[:query.Limit]
	}

	return &CardPoolView{
		Cards:      cards,
		TotalCount: totalCount,
		Filter:     query.Filter,
		SortField:  query.SortField,
		SortOrder:  query.SortOrder,
	}, nil
}

func (cp *CardPool) filterByProjectName(cards []*Card, projectName string) []*Card {
	if projectName == "" {
		return cards
	}

	var filtered []*Card
	for _, card := range cards {
		if card.ProjectName == projectName {
			filtered = append(filtered, card)
		}
	}
	return filtered
}

func (cp *CardPool) filterByPriority(cards []*Card, priority Priority) []*Card {
	var filtered []*Card
	for _, card := range cards {
		if card.Priority == priority {
			filtered = append(filtered, card)
		}
	}
	return filtered
}

func (cp *CardPool) sortCards(cards []*Card, field SortField, order SortOrder) []*Card {
	if field == "" {
		field = SortByUpdatedAt
	}
	if order == "" {
		order = SortDesc
	}

	sorted := make([]*Card, len(cards))
	copy(sorted, cards)

	switch field {
	case SortByPriority:
		sort.Slice(sorted, func(i, j int) bool {
			return cp.comparePriority(sorted[i].Priority, sorted[j].Priority, order)
		})
	case SortByCreatedAt:
		sort.Slice(sorted, func(i, j int) bool {
			return cp.compareTime(sorted[i].CreatedAt, sorted[j].CreatedAt, order)
		})
	case SortByUpdatedAt:
		sort.Slice(sorted, func(i, j int) bool {
			return cp.compareTime(sorted[i].UpdatedAt, sorted[j].UpdatedAt, order)
		})
	default:
		sort.Slice(sorted, func(i, j int) bool {
			return cp.compareTime(sorted[i].UpdatedAt, sorted[j].UpdatedAt, order)
		})
	}

	return sorted
}

func (cp *CardPool) comparePriority(p1, p2 Priority, order SortOrder) bool {
	priorityOrder := map[Priority]int{
		PriorityLow:      0,
		PriorityMedium:   1,
		PriorityHigh:     2,
		PriorityCritical: 3,
	}

	if order == SortAsc {
		return priorityOrder[p1] < priorityOrder[p2]
	}
	return priorityOrder[p1] > priorityOrder[p2]
}

func (cp *CardPool) compareTime(t1, t2 time.Time, order SortOrder) bool {
	if order == SortAsc {
		return t1.Before(t2)
	}
	return t1.After(t2)
}

func (cp *CardPool) BatchOperation(cardIDs []string, operation BatchOperation, params map[string]interface{}) (*BatchOperationResult, error) {
	if len(cardIDs) == 0 {
		return nil, errors.New("no card IDs provided for batch operation")
	}

	result := &BatchOperationResult{
		Errors: []error{},
	}

	projectName := ""
	for _, cardID := range cardIDs {
		card, err := cp.cardRepo.FindByID(cardID, "")
		if err != nil {
			result.FailureCount++
			result.Errors = append(result.Errors, fmt.Errorf("card %s not found: %w", cardID, err))
			continue
		}
		if card == nil {
			result.FailureCount++
			result.Errors = append(result.Errors, fmt.Errorf("card %s not found", cardID))
			continue
		}
		if projectName == "" {
			projectName = card.ProjectName
		}
		break
	}

	for _, cardID := range cardIDs {
		card, err := cp.cardRepo.FindByID(cardID, projectName)
		if err != nil {
			result.FailureCount++
			result.Errors = append(result.Errors, fmt.Errorf("card %s not found: %w", cardID, err))
			continue
		}

		var operationErr error
		switch operation {
		case BatchConfirm:
			operationErr = card.Confirm()
		case BatchAbandon:
			operationErr = card.Abandon()
		case BatchSupersede:
			newCardID, ok := params["new_card_id"].(string)
			if !ok {
				operationErr = errors.New("new_card_id parameter required for supersede operation")
			} else {
				operationErr = card.Supersede(newCardID)
			}
		default:
			operationErr = fmt.Errorf("unknown batch operation: %s", operation)
		}

		if operationErr != nil {
			result.FailureCount++
			result.Errors = append(result.Errors, fmt.Errorf("failed to %s card %s: %w", operation, cardID, operationErr))
			continue
		}

		if err := cp.cardRepo.Update(card); err != nil {
			result.FailureCount++
			result.Errors = append(result.Errors, fmt.Errorf("failed to save card %s after operation: %w", cardID, err))
			continue
		}

		result.SuccessCount++
	}

	return result, nil
}