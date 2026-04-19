package project

import (
	"fmt"
	"time"
)

type LifecycleStatus string

const (
	LifecycleCreating  LifecycleStatus = "CREATING"
	LifecycleActive    LifecycleStatus = "ACTIVE"
	LifecycleSuspended LifecycleStatus = "SUSPENDED"
	LifecycleArchived  LifecycleStatus = "ARCHIVED"
	LifecycleDeleted   LifecycleStatus = "DELETED"
)

var validTransitions = map[LifecycleStatus][]LifecycleStatus{
	LifecycleCreating:  {LifecycleActive, LifecycleDeleted},
	LifecycleActive:    {LifecycleSuspended, LifecycleArchived, LifecycleDeleted},
	LifecycleSuspended: {LifecycleActive, LifecycleArchived, LifecycleDeleted},
	LifecycleArchived:  {LifecycleDeleted},
	LifecycleDeleted:   {},
}

type ProjectLifecycle struct {
	ID          string          `json:"id"`
	ProjectID   string          `json:"project_id"`
	Status      LifecycleStatus `json:"status"`
	ActivatedAt *time.Time      `json:"activated_at,omitempty"`
	SuspendedAt *time.Time      `json:"suspended_at,omitempty"`
	ArchivedAt  *time.Time      `json:"archived_at,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

func NewProjectLifecycle(id, projectID string) *ProjectLifecycle {
	now := time.Now().UTC()
	return &ProjectLifecycle{
		ID:        id,
		ProjectID: projectID,
		Status:    LifecycleCreating,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (l *ProjectLifecycle) CanTransitionTo(target LifecycleStatus) bool {
	allowed, ok := validTransitions[l.Status]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == target {
			return true
		}
	}
	return false
}

func (l *ProjectLifecycle) Activate() error {
	if !l.CanTransitionTo(LifecycleActive) {
		return fmt.Errorf("cannot transition from %s to ACTIVE", l.Status)
	}
	now := time.Now().UTC()
	l.Status = LifecycleActive
	l.ActivatedAt = &now
	l.UpdatedAt = now
	return nil
}

func (l *ProjectLifecycle) Suspend() error {
	if !l.CanTransitionTo(LifecycleSuspended) {
		return fmt.Errorf("cannot transition from %s to SUSPENDED", l.Status)
	}
	now := time.Now().UTC()
	l.Status = LifecycleSuspended
	l.SuspendedAt = &now
	l.UpdatedAt = now
	return nil
}

func (l *ProjectLifecycle) Archive() error {
	if !l.CanTransitionTo(LifecycleArchived) {
		return fmt.Errorf("cannot transition from %s to ARCHIVED", l.Status)
	}
	now := time.Now().UTC()
	l.Status = LifecycleArchived
	l.ArchivedAt = &now
	l.UpdatedAt = now
	return nil
}

func (l *ProjectLifecycle) Delete() error {
	if !l.CanTransitionTo(LifecycleDeleted) {
		return fmt.Errorf("cannot transition from %s to DELETED", l.Status)
	}
	now := time.Now().UTC()
	l.Status = LifecycleDeleted
	l.UpdatedAt = now
	return nil
}
