package knowledge

import (
	"testing"
)

func TestNewDocument(t *testing.T) {
	doc := NewDocument("id-1", "test.md", "/docs/test.md", "proj-1")
	if doc.ID != "id-1" {
		t.Errorf("expected ID id-1, got %s", doc.ID)
	}
	if doc.Version != 1 {
		t.Errorf("expected version 1, got %d", doc.Version)
	}
	if doc.DirectoryID != "" {
		t.Errorf("expected empty directory ID, got %s", doc.DirectoryID)
	}
}

func TestDocument_SetDirectory(t *testing.T) {
	doc := NewDocument("id-1", "test.md", "/docs/test.md", "proj-1")
	doc.SetDirectory("dir-1")
	if doc.DirectoryID != "dir-1" {
		t.Errorf("expected directory ID dir-1, got %s", doc.DirectoryID)
	}
}

func TestDocument_IncrementVersion(t *testing.T) {
	doc := NewDocument("id-1", "test.md", "/docs/test.md", "proj-1")
	doc.IncrementVersion()
	if doc.Version != 2 {
		t.Errorf("expected version 2, got %d", doc.Version)
	}
}

func TestDocument_Validate(t *testing.T) {
	doc := NewDocument("id-1", "test.md", "/docs/test.md", "proj-1")
	if err := doc.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	noName := NewDocument("id-2", "", "/docs/test.md", "proj-1")
	if err := noName.Validate(); err == nil {
		t.Error("expected error for empty name")
	}

	noProject := NewDocument("id-3", "test.md", "/docs/test.md", "")
	if err := noProject.Validate(); err == nil {
		t.Error("expected error for empty project ID")
	}
}

func TestNewDirectory(t *testing.T) {
	dir := NewDirectory("dir-1", "docs", "", "proj-1")
	if dir.ID != "dir-1" {
		t.Errorf("expected ID dir-1, got %s", dir.ID)
	}
	if !dir.IsRoot() {
		t.Error("expected directory to be root")
	}
	if dir.Path != "docs" {
		t.Errorf("expected path docs, got %s", dir.Path)
	}
}

func TestNewDirectory_WithParent(t *testing.T) {
	dir := NewDirectory("dir-2", "sub", "parent-1", "proj-1")
	if dir.IsRoot() {
		t.Error("expected directory to not be root")
	}
	if dir.ParentID != "parent-1" {
		t.Errorf("expected parent ID parent-1, got %s", dir.ParentID)
	}
}

func TestDirectory_Validate(t *testing.T) {
	dir := NewDirectory("dir-1", "docs", "", "proj-1")
	if err := dir.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	noName := NewDirectory("dir-2", "", "", "proj-1")
	if err := noName.Validate(); err == nil {
		t.Error("expected error for empty name")
	}

	noProject := NewDirectory("dir-3", "docs", "", "")
	if err := noProject.Validate(); err == nil {
		t.Error("expected error for empty project ID")
	}
}

func TestChunk(t *testing.T) {
	chunk := NewChunk("chunk-1", "doc-1", "content", 0)
	if chunk.ID != "chunk-1" {
		t.Errorf("expected ID chunk-1, got %s", chunk.ID)
	}
	if chunk.Index != 0 {
		t.Errorf("expected index 0, got %d", chunk.Index)
	}
}
