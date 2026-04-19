package clonepool

import (
	"fmt"
	"sync"
)

type ClonePool struct {
	RepositoryID string
	instances    map[string]*CloneInstance
	mu           sync.RWMutex
}

func NewClonePool(repositoryID string) *ClonePool {
	return &ClonePool{
		RepositoryID: repositoryID,
		instances:    make(map[string]*CloneInstance),
	}
}

func (p *ClonePool) Add(instance *CloneInstance) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.instances[instance.ID] = instance
}

func (p *ClonePool) Acquire(taskID string) (*CloneInstance, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, inst := range p.instances {
		if inst.Status == InstanceIdle {
			if err := inst.Acquire(taskID); err != nil {
				return nil, err
			}
			return inst, nil
		}
	}
	return nil, fmt.Errorf("no idle clone instance available")
}

func (p *ClonePool) Release(instanceID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	inst, ok := p.instances[instanceID]
	if !ok {
		return fmt.Errorf("instance %s not found", instanceID)
	}
	return inst.Release()
}

func (p *ClonePool) MarkDirty(instanceID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	inst, ok := p.instances[instanceID]
	if !ok {
		return fmt.Errorf("instance %s not found", instanceID)
	}
	return inst.MarkDirty()
}

func (p *ClonePool) List() []*CloneInstance {
	p.mu.RLock()
	defer p.mu.RUnlock()
	result := make([]*CloneInstance, 0, len(p.instances))
	for _, inst := range p.instances {
		result = append(result, inst)
	}
	return result
}

func (p *ClonePool) FindByID(instanceID string) (*CloneInstance, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	inst, ok := p.instances[instanceID]
	return inst, ok
}
