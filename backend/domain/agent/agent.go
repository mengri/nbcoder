package agent

import (
	"fmt"
	"sync"
)

type Agent struct {
	ID     string
	Type   AgentType
	Skills map[string]*Skill
}

func NewAgent(id string, agentType AgentType) *Agent {
	return &Agent{
		ID:     id,
		Type:   agentType,
		Skills: make(map[string]*Skill),
	}
}

func (a *Agent) AddSkill(skill *Skill) {
	a.Skills[skill.Name] = skill
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
	return nil, fmt.Errorf("skill %s not found for agent type %s", skillName, agentType)
}

func (r *AgentRegistry) List() []*Agent {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*Agent, 0, len(r.agents))
	for _, agent := range r.agents {
		result = append(result, agent)
	}
	return result
}
