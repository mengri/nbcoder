package requirement

import (
	"fmt"
	"time"
)

type Priority string

const (
	PriorityLow      Priority = "LOW"
	PriorityMedium   Priority = "MEDIUM"
	PriorityHigh     Priority = "HIGH"
	PriorityCritical Priority = "CRITICAL"
)

func (p Priority) IsValid() bool {
	switch p {
	case PriorityLow, PriorityMedium, PriorityHigh, PriorityCritical:
		return true
	}
	return false
}

type Card struct {
	ID               string     `json:"id"`
	Title            string     `json:"title"`
	Description      string     `json:"description"`
	Original         string     `json:"original"`
	Status           CardStatus `json:"status"`
	Priority         Priority   `json:"priority"`
	StructuredOutput string     `json:"structured_output,omitempty"`
	PipelineID       string     `json:"pipeline_id,omitempty"`
	ProjectName      string     `json:"project_name"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	SupersededBy     string     `json:"superseded_by,omitempty"`
}

func NewCard(id, title, description, original, projectName string) *Card {
	now := time.Now().UTC()
	return &Card{
		ID:           id,
		Title:        title,
		Description:  description,
		Original:     original,
		Status:       CardDraft,
		Priority:     PriorityMedium,
		ProjectName:  projectName,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func (c *Card) SetPriority(p Priority) error {
	if !p.IsValid() {
		return fmt.Errorf("invalid priority: %s", p)
	}
	c.Priority = p
	c.UpdatedAt = time.Now().UTC()
	return nil
}

func (c *Card) SetStructuredOutput(output string) {
	c.StructuredOutput = output
	c.UpdatedAt = time.Now().UTC()
}

func (c *Card) SetPipelineID(pipelineID string) {
	c.PipelineID = pipelineID
	c.UpdatedAt = time.Now().UTC()
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
