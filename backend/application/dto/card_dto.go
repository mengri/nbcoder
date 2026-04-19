package dto

type CreateCardRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Original    string `json:"original"`
	Priority    string `json:"priority"`
	ProjectID   string `json:"project_id" binding:"required"`
}

type CardResponse struct {
	ID               string `json:"id"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	Original         string `json:"original"`
	Status           string `json:"status"`
	Priority         string `json:"priority"`
	StructuredOutput string `json:"structured_output,omitempty"`
	PipelineID       string `json:"pipeline_id,omitempty"`
	ProjectID        string `json:"project_id"`
	SupersededBy     string `json:"superseded_by,omitempty"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

type AddDependencyRequest struct {
	DependsOnID string `json:"depends_on_id" binding:"required"`
	Type        string `json:"type" binding:"required"`
}

type UpdateCardRequest struct {
	Title            *string `json:"title,omitempty"`
	Description      *string `json:"description,omitempty"`
	Priority         *string `json:"priority,omitempty"`
	StructuredOutput *string `json:"structured_output,omitempty"`
	PipelineID       *string `json:"pipeline_id,omitempty"`
}
