package requirement

import (
	"fmt"
	"time"
)

type Card struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Original    string     `json:"original"`
	Status      CardStatus `json:"status"`
	ProjectID   string     `json:"project_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	SupersededBy string   `json:"superseded_by,omitempty"`
}

func NewCard(id, title, description, original, projectID string) *Card {
	now := time.Now().UTC()
	return &Card{
		ID:          id,
		Title:       title,
		Description: description,
		Original:    original,
		Status:      CardDraft,
		ProjectID:   projectID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (c *Card) Confirm() error {
	if c.Status != CardDraft {
		return fmt.Errorf("cannot confirm card in status %s", c.Status)
	}
	c.Status = CardConfirmed
	c.UpdatedAt = time.Now().UTC()
	return nil
}

func (c *Card) Start() error {
	if c.Status != CardConfirmed {
		return fmt.Errorf("cannot start card in status %s", c.Status)
	}
	c.Status = CardInProgress
	c.UpdatedAt = time.Now().UTC()
	return nil
}

func (c *Card) Complete() error {
	if c.Status != CardInProgress {
		return fmt.Errorf("cannot complete card in status %s", c.Status)
	}
	c.Status = CardCompleted
	c.UpdatedAt = time.Now().UTC()
	return nil
}

func (c *Card) Supersede(newCardID string) error {
	if c.Status == CardAbandoned || c.Status == CardSuperseded {
		return fmt.Errorf("cannot supersede card in status %s", c.Status)
	}
	c.Status = CardSuperseded
	c.SupersededBy = newCardID
	c.UpdatedAt = time.Now().UTC()
	return nil
}

func (c *Card) Abandon() error {
	if c.Status == CardCompleted || c.Status == CardAbandoned {
		return fmt.Errorf("cannot abandon card in status %s", c.Status)
	}
	c.Status = CardAbandoned
	c.UpdatedAt = time.Now().UTC()
	return nil
}
