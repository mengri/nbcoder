package dto

type CreateCardRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Original    string `json:"original"`
	Priority    string `json:"priority"`
	ProjectName string `json:"projectName" binding:"required"`
}

type CardResponse struct {
	ID               string `json:"id"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	Original         string `json:"original"`
	Status           string `json:"status"`
	Priority         string `json:"priority"`
	StructuredOutput string `json:"structuredOutput,omitempty"`
	PipelineID       string `json:"pipelineId,omitempty"`
	ProjectName      string `json:"projectName"`
	SupersededBy     string `json:"supersededBy,omitempty"`
	CreatedAt        string `json:"createdAt"`
	UpdatedAt        string `json:"updatedAt"`
}

type AddDependencyRequest struct {
	DependsOnID string `json:"dependsOnId" binding:"required"`
	Type        string `json:"type" binding:"required"`
}

type UpdateCardRequest struct {
	Title            *string `json:"title,omitempty"`
	Description      *string `json:"description,omitempty"`
	Priority         *string `json:"priority,omitempty"`
	StructuredOutput *string `json:"structuredOutput,omitempty"`
	PipelineID       *string `json:"pipelineId,omitempty"`
}
