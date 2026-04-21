package dto

type CreateProjectRequest struct {
	Name             string `json:"name" binding:"required"`
	Description      string `json:"description"`
	RepoURL          string `json:"repoUrl"`
	BranchStrategy   string `json:"branchStrategy"`
	TechStack        string `json:"techStack"`
	CodingConventions string `json:"codingConventions"`
}

type ProjectResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	RepoURL     string `json:"repoUrl"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type InitProjectResponse struct {
	Project     ProjectResponse    `json:"project"`
	Directories []string           `json:"directories"`
	Configs     []ConfigResponse   `json:"configs"`
	Standards   *StandardsResponse `json:"standards,omitempty"`
}

type ConfigResponse struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type StandardsResponse struct {
	ID                string `json:"id"`
	BranchStrategy    string `json:"branchStrategy"`
	TechStack         string `json:"techStack"`
	CodingConventions string `json:"codingConventions"`
}

type ConfigChangeLogResponse struct {
	ID          string `json:"id"`
	ProjectName string `json:"projectName"`
	ConfigKey   string `json:"configKey"`
	OldValue    string `json:"oldValue"`
	NewValue    string `json:"newValue"`
	ChangedAt   string `json:"changedAt"`
	ChangedBy   string `json:"changedBy"`
}
