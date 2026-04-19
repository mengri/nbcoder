package dto

type CreateTaskRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	ProjectID   string `json:"project_id"`
}

type AssignTaskRequest struct {
	AgentID string `json:"agent_id" binding:"required"`
}

type TaskResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	AssignedTo  string `json:"assigned_to,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
