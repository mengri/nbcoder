package persistence

import (
	"sync"

	"github.com/mengri/nbcoder/domain/agent"
)

type InMemoryTaskRepo struct {
	tasks map[string]*agent.Task
	mu    sync.RWMutex
}

func NewInMemoryTaskRepo() *InMemoryTaskRepo {
	return &InMemoryTaskRepo{
		tasks: make(map[string]*agent.Task),
	}
}

func (r *InMemoryTaskRepo) Save(task *agent.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tasks[task.ID] = task
	return nil
}

func (r *InMemoryTaskRepo) FindByID(id string, projectName string) (*agent.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	task, ok := r.tasks[id]
	if !ok {
		return nil, nil
	}
	return task, nil
}

func (r *InMemoryTaskRepo) FindByStatus(status agent.TaskStatus) ([]*agent.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*agent.Task
	for _, t := range r.tasks {
		if t.Status == status {
			result = append(result, t)
		}
	}
	return result, nil
}

func (r *InMemoryTaskRepo) Update(task *agent.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tasks[task.ID] = task
	return nil
}

func (r *InMemoryTaskRepo) FindByProjectName(projectName string) ([]*agent.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*agent.Task
	for _, t := range r.tasks {
		if t.ProjectName == projectName {
			result = append(result, t)
		}
	}
	return result, nil
}

func (r *InMemoryTaskRepo) FindByAgentID(agentID string) ([]*agent.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*agent.Task
	for _, t := range r.tasks {
		if t.AssignedTo == agentID {
			result = append(result, t)
		}
	}
	return result, nil
}

func (r *InMemoryTaskRepo) FindByPipelineID(pipelineID string) ([]*agent.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*agent.Task
	for _, t := range r.tasks {
		if t.PipelineID == pipelineID {
			result = append(result, t)
		}
	}
	return result, nil
}

func (r *InMemoryTaskRepo) FindAll() ([]*agent.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*agent.Task, 0, len(r.tasks))
	for _, t := range r.tasks {
		result = append(result, t)
	}
	return result, nil
}

func (r *InMemoryTaskRepo) Delete(id string, projectName string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.tasks, id)
	return nil
}
