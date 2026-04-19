package requirement

import (
	"fmt"
)

type DependencyType string

const (
	DependencyBlocks   DependencyType = "BLOCKS"
	DependencyDependsOn DependencyType = "DEPENDS_ON"
)

type CardDependency struct {
	ID           string         `json:"id"`
	CardID       string         `json:"card_id"`
	DependsOnID  string         `json:"depends_on_id"`
	Type         DependencyType `json:"type"`
}

func NewCardDependency(id, cardID, dependsOnID string, depType DependencyType) *CardDependency {
	return &CardDependency{
		ID:          id,
		CardID:      cardID,
		DependsOnID: dependsOnID,
		Type:        depType,
	}
}

func (d *CardDependency) Validate() error {
	if d.ID == "" {
		return fmt.Errorf("dependency ID cannot be empty")
	}
	if d.CardID == "" {
		return fmt.Errorf("card ID cannot be empty")
	}
	if d.DependsOnID == "" {
		return fmt.Errorf("depends on ID cannot be empty")
	}
	if d.Type != DependencyBlocks && d.Type != DependencyDependsOn {
		return fmt.Errorf("invalid dependency type: %s", d.Type)
	}
	return nil
}
