package project

type ProjectConfig struct {
	ID        string `json:"id"`
	ProjectID string `json:"project_id"`
	Key       string `json:"key"`
	Value     string `json:"value"`
}

type Standards struct {
	ID                 string `json:"id"`
	ProjectID          string `json:"project_id"`
	BranchStrategy     string `json:"branch_strategy,omitempty"`
	TechStack          string `json:"tech_stack,omitempty"`
	CodingConventions  string `json:"coding_conventions,omitempty"`
}
