package agent

import (
	"fmt"
	"time"
)

type Skill struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	AgentType   AgentType              `json:"agent_type"`
	Config      map[string]interface{} `json:"config"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

func NewSkill(id, name, description string, agentType AgentType) *Skill {
	now := time.Now().UTC()
	return &Skill{
		ID:          id,
		Name:        name,
		Description: description,
		AgentType:   agentType,
		Config:      make(map[string]interface{}),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (s *Skill) SetConfig(key string, value interface{}) {
	if s.Config == nil {
		s.Config = make(map[string]interface{})
	}
	s.Config[key] = value
	s.UpdatedAt = time.Now().UTC()
}

func (s *Skill) GetConfig(key string) (interface{}, bool) {
	if s.Config == nil {
		return nil, false
	}
	value, exists := s.Config[key]
	return value, exists
}

func (s *Skill) Invoke(args ...interface{}) (interface{}, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("skill %s: no arguments provided", s.Name)
	}
	s.UpdatedAt = time.Now().UTC()
	return fmt.Sprintf("skill %s executed with %d args", s.Name, len(args)), nil
}
