package clonepool

import (
	"fmt"
	"time"
)

type CloneInstance struct {
	ID           string              `json:"id"`
	RepositoryID string              `json:"repository_id"`
	Status       CloneInstanceStatus `json:"status"`
	AssignedTask string              `json:"assigned_task,omitempty"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

func NewCloneInstance(id, repositoryID string) *CloneInstance {
	return &CloneInstance{
		ID:           id,
		RepositoryID: repositoryID,
		Status:       InstanceIdle,
		UpdatedAt:    time.Now().UTC(),
	}
}

func (ci *CloneInstance) Acquire(taskID string) error {
	if ci.Status != InstanceIdle {
		return fmt.Errorf("cannot acquire instance in status %s", ci.Status)
	}
	ci.Status = InstanceBusy
	ci.AssignedTask = taskID
	ci.UpdatedAt = time.Now().UTC()
	return nil
}

func (ci *CloneInstance) Release() error {
	if ci.Status == InstanceDirty {
		ci.Recycle()
		return nil
	}
	if ci.Status != InstanceBusy {
		return fmt.Errorf("cannot release instance in status %s", ci.Status)
	}
	ci.Status = InstanceIdle
	ci.AssignedTask = ""
	ci.UpdatedAt = time.Now().UTC()
	return nil
}

func (ci *CloneInstance) MarkDirty() error {
	if ci.Status != InstanceBusy {
		return fmt.Errorf("cannot mark dirty instance in status %s", ci.Status)
	}
	ci.Status = InstanceDirty
	ci.UpdatedAt = time.Now().UTC()
	return nil
}

func (ci *CloneInstance) Recycle() {
	ci.Status = InstanceIdle
	ci.AssignedTask = ""
	ci.UpdatedAt = time.Now().UTC()
}
