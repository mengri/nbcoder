package agent

type TaskRepo interface {
	Save(task *Task) error
	FindByID(id string) (*Task, error)
	FindByProjectID(projectID string) ([]*Task, error)
	FindByStatus(status TaskStatus) ([]*Task, error)
	FindByAgentID(agentID string) ([]*Task, error)
	FindByPipelineID(pipelineID string) ([]*Task, error)
	FindAll() ([]*Task, error)
	Update(task *Task) error
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
