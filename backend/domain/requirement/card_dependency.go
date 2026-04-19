package requirement

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
