package clonepool

import (
	"github.com/mengri/nbcoder/domain/clonepool"
	"github.com/mengri/nbcoder/domain/event"
	"github.com/mengri/nbcoder/pkg/uid"
)

type ClonePoolService struct {
	instanceRepo clonepool.CloneInstanceRepo
	repoRepo     clonepool.RepositoryRepo
	eventBus     event.EventBus
}

func NewClonePoolService(
	instanceRepo clonepool.CloneInstanceRepo,
	repoRepo clonepool.RepositoryRepo,
	eventBus event.EventBus,
) *ClonePoolService {
	return &ClonePoolService{
		instanceRepo: instanceRepo,
		repoRepo:     repoRepo,
		eventBus:     eventBus,
	}
}

func (s *ClonePoolService) AcquireInstance(repositoryID, taskID string) (*clonepool.CloneInstance, error) {
	instances, err := s.instanceRepo.FindByRepositoryID(repositoryID)
	if err != nil {
		return nil, err
	}
	pool := clonepool.NewClonePool(repositoryID)
	for _, inst := range instances {
		pool.Add(inst)
	}
	inst, err := pool.Acquire(taskID)
	if err != nil {
		return nil, err
	}
	if err := s.instanceRepo.Update(inst); err != nil {
		return nil, err
	}
	evt := event.NewClonePoolEvent(uid.NewID(), inst.ID, event.CloneAcquiredEvent)
	_ = s.eventBus.Publish(evt)
	return inst, nil
}

func (s *ClonePoolService) ReleaseInstance(instanceID string) error {
	inst, err := s.instanceRepo.FindByID(instanceID)
	if err != nil {
		return err
	}
	if err := inst.Release(); err != nil {
		return err
	}
	if err := s.instanceRepo.Update(inst); err != nil {
		return err
	}
	evt := event.NewClonePoolEvent(uid.NewID(), instanceID, event.CloneReleasedEvent)
	return s.eventBus.Publish(evt)
}
