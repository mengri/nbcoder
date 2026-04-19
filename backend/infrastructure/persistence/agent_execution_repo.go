package persistence

import (
	"sync"
	"time"

	"github.com/mengri/nbcoder/domain/agent"
)

type InMemoryAgentExecutionRepo struct {
	executions map[string]*agent.AgentExecution
	mu         sync.RWMutex
}

func NewInMemoryAgentExecutionRepo() *InMemoryAgentExecutionRepo {
	return &InMemoryAgentExecutionRepo{
		executions: make(map[string]*agent.AgentExecution),
	}
}

func (r *InMemoryAgentExecutionRepo) Save(execution *agent.AgentExecution) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.executions[execution.ID] = execution
	return nil
}

func (r *InMemoryAgentExecutionRepo) QueryByTask(taskID string) ([]*agent.AgentExecution, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*agent.AgentExecution
	for _, e := range r.executions {
		if e.TaskID == taskID {
			result = append(result, e)
		}
	}
	return result, nil
}

func (r *InMemoryAgentExecutionRepo) QueryByAgent(agentID string) ([]*agent.AgentExecution, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*agent.AgentExecution
	for _, e := range r.executions {
		if e.AgentID == agentID {
			result = append(result, e)
		}
	}
	return result, nil
}

func (r *InMemoryAgentExecutionRepo) QueryByTimeRange(start, end time.Time) ([]*agent.AgentExecution, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*agent.AgentExecution
	for _, e := range r.executions {
		if e.Timestamp.After(start) && e.Timestamp.Before(end) {
			result = append(result, e)
		}
	}
	return result, nil
}
