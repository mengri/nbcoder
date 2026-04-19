package requirement

type CardStatus string

const (
	CardDraft      CardStatus = "DRAFT"
	CardConfirmed  CardStatus = "CONFIRMED"
	CardInProgress CardStatus = "IN_PROGRESS"
	CardCompleted  CardStatus = "COMPLETED"
	CardSuperseded CardStatus = "SUPERSEDED"
	CardAbandoned  CardStatus = "ABANDONED"
)

func (s CardStatus) IsValid() bool {
	switch s {
	case CardDraft, CardConfirmed, CardInProgress, CardCompleted, CardSuperseded, CardAbandoned:
		return true
	}
	return false
}
