package knowledge

import (
	"fmt"
	"time"
)

type Document struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	ProjectName string    `json:"project_name"`
	DirectoryID string    `json:"directory_id,omitempty"`
	Content     string    `json:"content,omitempty"`
	Version     int       `json:"version"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewDocument(id, name, path, projectName string) *Document {
	now := time.Now().UTC()
	return &Document{
		ID:           id,
		Name:         name,
		Path:         path,
		ProjectName:  projectName,
		Version:      1,
		CreatedAt:    now,
		UpdatedAt:    now,
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
	if d.ProjectName == "" {
		return fmt.Errorf("project name is required")
	}
	return nil
}

type Directory struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	ParentID    string    `json:"parent_id,omitempty"`
	ProjectName string    `json:"project_name"`
	Path        string    `json:"path"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewDirectory(id, name, parentID, projectName string) *Directory {
	now := time.Now().UTC()
	path := name
	if parentID != "" {
		path = parentID + "/" + name
	}
	return &Directory{
		ID:          id,
		Name:        name,
		ParentID:    parentID,
		ProjectName: projectName,
		Path:        path,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (d *Directory) Validate() error {
	if d.Name == "" {
		return fmt.Errorf("directory name is required")
	}
	if d.ProjectName == "" {
		return fmt.Errorf("project name is required")
	}
	return nil
}

func (d *Directory) IsRoot() bool {
	return d.ParentID == ""
}
