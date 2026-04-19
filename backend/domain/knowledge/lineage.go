package knowledge

import (
	"time"

	"github.com/google/uuid"
)

type ChangeType string

const (
	ChangeTypeCreated  ChangeType = "CREATED"
	ChangeTypeModified ChangeType = "MODIFIED"
	ChangeTypeRenamed  ChangeType = "RENAMED"
	ChangeTypeDeleted  ChangeType = "DELETED"
)

type DocumentLineage struct {
	ID               string     `json:"id"`
	DocumentID       string     `json:"document_id"`
	ParentDocumentID string     `json:"parent_document_id,omitempty"`
	ChangeType       ChangeType `json:"change_type"`
	Description      string     `json:"description,omitempty"`
	ChangedAt        time.Time  `json:"changed_at"`
}

func NewDocumentLineage(documentID, parentDocumentID string, changeType ChangeType, description string) *DocumentLineage {
	return &DocumentLineage{
		ID:               uuid.New().String(),
		DocumentID:       documentID,
		ParentDocumentID: parentDocumentID,
		ChangeType:       changeType,
		Description:      description,
		ChangedAt:        time.Now().UTC(),
	}
}
