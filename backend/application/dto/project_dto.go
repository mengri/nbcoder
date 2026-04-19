package dto

type CreateProjectRequest struct {
	Name             string `json:"name" binding:"required"`
	Description      string `json:"description"`
	RepoURL          string `json:"repo_url"`
	BranchStrategy   string `json:"branch_strategy"`
	TechStack        string `json:"tech_stack"`
	CodingConventions string `json:"coding_conventions"`
}

type ProjectResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	RepoURL     string `json:"repo_url"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type InitProjectResponse struct {
	Project     ProjectResponse   `json:"project"`
	Directories []string          `json:"directories"`
	Configs     []ConfigResponse  `json:"configs"`
	Standards   *StandardsResponse `json:"standards,omitempty"`
}

type ConfigResponse struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type StandardsResponse struct {
	ID                string `json:"id"`
	BranchStrategy    string `json:"branch_strategy"`
	TechStack         string `json:"tech_stack"`
	CodingConventions string `json:"coding_conventions"`
}
