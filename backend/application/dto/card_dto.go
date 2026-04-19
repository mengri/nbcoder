package dto

type CreateCardRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Original    string `json:"original"`
	ProjectID   string `json:"project_id" binding:"required"`
}

type CardResponse struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Status       string `json:"status"`
	ProjectID    string `json:"project_id"`
	SupersededBy string `json:"superseded_by,omitempty"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type AddDependencyRequest struct {
	DependsOnID string `json:"depends_on_id" binding:"required"`
	Type        string `json:"type" binding:"required"`
}
