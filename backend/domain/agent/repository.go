package agent

type AgentTaskRepo interface {
	Save(task *AgentTask) error
	FindByID(id string) (*AgentTask, error)
	FindByProjectID(projectID string) ([]*AgentTask, error)
	FindByStatus(status AgentTaskStatus) ([]*AgentTask, error)
	FindByAgentID(agentID string) ([]*AgentTask, error)
	FindByPipelineID(pipelineID string) ([]*AgentTask, error)
	FindAll() ([]*AgentTask, error)
	Update(task *AgentTask) error
	Delete(id string) error
}

type SkillRepo interface {
	Save(skill *Skill) error
	FindByID(id string) (*Skill, error)
	FindByAgentType(agentType AgentType) ([]*Skill, error)
	FindAll() ([]*Skill, error)
	Update(skill *Skill) error
	Delete(id string) error
}

type AgentExecutionRepo interface {
	Save(execution *AgentExecution) error
	FindByID(id string) (*AgentExecution, error)
	FindByTaskID(taskID string) ([]*AgentExecution, error)
	FindByAgentID(agentID string) ([]*AgentExecution, error)
	FindAll() ([]*AgentExecution, error)
	Update(execution *AgentExecution) error
	Delete(id string) error
}
