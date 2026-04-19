package clonepool
// pool.go
// 克隆实例池管理


import "sync"

type ClonePool struct {
	instances map[string]*CloneInstance
	mu        sync.RWMutex
}

func NewClonePool() *ClonePool {
	return &ClonePool{
		instances: make(map[string]*CloneInstance),
	}
}

func (p *ClonePool) Add(instance *CloneInstance) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.instances[instance.ID] = instance
}

func (p *ClonePool) Allocate() *CloneInstance {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, inst := range p.instances {
		if inst.Status == InstanceIdle {
			inst.Allocate()
			return inst
		}
	}
	return nil
}

func (p *ClonePool) Release(id string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if inst, ok := p.instances[id]; ok {
		inst.Release()
	}
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
