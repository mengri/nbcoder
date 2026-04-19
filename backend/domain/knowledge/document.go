package knowledge

import (
	"fmt"
	"time"
)

type Document struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	ProjectID   string    `json:"project_id"`
	DirectoryID string    `json:"directory_id,omitempty"`
	Content     string    `json:"content,omitempty"`
	Version     int       `json:"version"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewDocument(id, name, path, projectID string) *Document {
	now := time.Now().UTC()
	return &Document{
		ID:        id,
		Name:      name,
		Path:      path,
		ProjectID: projectID,
		Version:   1,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (d *Document) IncrementVersion() {
	d.Version++
	d.UpdatedAt = time.Now().UTC()
}

func (d *Document) SetDirectory(directoryID string) {
	d.DirectoryID = directoryID
	d.UpdatedAt = time.Now().UTC()
}

func (d *Document) SetContent(content string) {
	d.Content = content
	d.UpdatedAt = time.Now().UTC()
}

func (d *Document) Validate() error {
	if d.Name == "" {
		return fmt.Errorf("document name is required")
	}
	if d.ProjectID == "" {
		return fmt.Errorf("project ID is required")
	}
	return nil
}

type Directory struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	ParentID  string    `json:"parent_id,omitempty"`
	ProjectID string    `json:"project_id"`
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewDirectory(id, name, parentID, projectID string) *Directory {
	now := time.Now().UTC()
	path := name
	if parentID != "" {
		path = parentID + "/" + name
	}
	return &Directory{
		ID:        id,
		Name:      name,
		ParentID:  parentID,
		ProjectID: projectID,
		Path:      path,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (d *Directory) Validate() error {
	if d.Name == "" {
		return fmt.Errorf("directory name is required")
	}
	if d.ProjectID == "" {
		return fmt.Errorf("project ID is required")
	}
	return nil
}

func (d *Directory) IsRoot() bool {
	return d.ParentID == ""
}
