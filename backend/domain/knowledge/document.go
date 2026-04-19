package knowledge

import "time"

type Document struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	ProjectID string    `json:"project_id"`
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
