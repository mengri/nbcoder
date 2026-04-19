
// agent.go
// Agent、Skill 及调度机制领域模型定义
package domain

import (
	"sync"
)

type AgentType string

type Skill struct {
	Name        string
	Description string
	Invoke      func(args ...interface{}) (interface{}, error)
}

type Agent struct {
	ID          string
	Type        AgentType
	Skills      map[string]*Skill
}

type AgentRegistry struct {
	agents map[string]*Agent
	mu     sync.RWMutex
}

func NewAgentRegistry() *AgentRegistry {
	return &AgentRegistry{
		agents: make(map[string]*Agent),
	}
}

func (r *AgentRegistry) Register(agent *Agent) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.agents[agent.ID] = agent
}

func (r *AgentRegistry) Get(id string) (*Agent, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	agent, ok := r.agents[id]
	return agent, ok
}

func (r *AgentRegistry) Dispatch(agentType AgentType, skillName string, args ...interface{}) (interface{}, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, agent := range r.agents {
		if agent.Type == agentType {
			skill, ok := agent.Skills[skillName]
			if ok {
				return skill.Invoke(args...)
			}
		}
	}
	return nil, ErrSkillNotFound
}

var ErrSkillNotFound = &SkillError{"Skill not found"}

type SkillError struct {
	Msg string
}

func (e *SkillError) Error() string {
	return e.Msg
}
