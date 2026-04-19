package knowledge

import (
	"testing"
)

func TestNewDocumentLineage(t *testing.T) {
	l := NewDocumentLineage("doc-1", "doc-parent", ChangeTypeModified, "updated content")
	if l.ID == "" {
		t.Error("expected non-empty ID")
	}
	if l.DocumentID != "doc-1" {
		t.Errorf("expected DocumentID doc-1, got %s", l.DocumentID)
	}
	if l.ParentDocumentID != "doc-parent" {
		t.Errorf("expected ParentDocumentID doc-parent, got %s", l.ParentDocumentID)
	}
	if l.ChangeType != ChangeTypeModified {
		t.Errorf("expected ChangeType MODIFIED, got %s", l.ChangeType)
	}
	if l.Description != "updated content" {
		t.Errorf("expected Description 'updated content', got %s", l.Description)
	}
	if l.ChangedAt.IsZero() {
		t.Error("expected non-zero ChangedAt")
	}
}

func TestNewDocumentLineage_NoParent(t *testing.T) {
	l := NewDocumentLineage("doc-1", "", ChangeTypeCreated, "initial creation")
	if l.ParentDocumentID != "" {
		t.Errorf("expected empty ParentDocumentID, got %s", l.ParentDocumentID)
	}
	if l.ChangeType != ChangeTypeCreated {
		t.Errorf("expected ChangeType CREATED, got %s", l.ChangeType)
	}
}

func TestChangeTypeConstants(t *testing.T) {
	tests := []struct {
		ct       ChangeType
		expected string
	}{
		{ChangeTypeCreated, "CREATED"},
		{ChangeTypeModified, "MODIFIED"},
		{ChangeTypeRenamed, "RENAMED"},
		{ChangeTypeDeleted, "DELETED"},
	}
	for _, tt := range tests {
		if string(tt.ct) != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, tt.ct)
		}
	}
}
