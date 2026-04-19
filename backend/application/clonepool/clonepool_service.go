package clonepool

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/mengri/nbcoder/domain/clonepool"
	"github.com/mengri/nbcoder/domain/event"
	"github.com/mengri/nbcoder/infrastructure/git"
	"github.com/mengri/nbcoder/pkg/uid"
)

type ClonePoolService struct {
	instanceRepo clonepool.CloneInstanceRepo
	repoRepo     clonepool.RepositoryRepo
	eventBus     event.EventBus
	gitClient    git.GitClient
	baseDir      string
}

func NewClonePoolService(
	instanceRepo clonepool.CloneInstanceRepo,
	repoRepo clonepool.RepositoryRepo,
	eventBus event.EventBus,
	gitClient git.GitClient,
	baseDir string,
) *ClonePoolService {
	return &ClonePoolService{
		instanceRepo: instanceRepo,
		repoRepo:     repoRepo,
		eventBus:     eventBus,
		gitClient:    gitClient,
		baseDir:      baseDir,
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

func (s *ClonePoolService) CreateCloneInstance(ctx context.Context, repositoryID, repoURL string) (*clonepool.CloneInstance, error) {
	inst := clonepool.NewCloneInstance(uid.NewID(), repositoryID)

	targetDir := filepath.Join(s.baseDir, repositoryID, inst.ID)
	if err := s.gitClient.Clone(ctx, repoURL, targetDir); err != nil {
		return nil, fmt.Errorf("failed to clone repository: %w", err)
	}

	if err := s.instanceRepo.Save(inst); err != nil {
		return nil, fmt.Errorf("failed to save instance: %w", err)
	}

	return inst, nil
}

func (s *ClonePoolService) CommitChanges(ctx context.Context, instanceID, message string) error {
	inst, err := s.instanceRepo.FindByID(instanceID)
	if err != nil {
		return err
	}

	targetDir := filepath.Join(s.baseDir, inst.RepositoryID, inst.ID)
	if err := s.gitClient.Commit(ctx, targetDir, message); err != nil {
		return err
	}

	return nil
}

func (s *ClonePoolService) PullLatest(ctx context.Context, instanceID string) error {
	inst, err := s.instanceRepo.FindByID(instanceID)
	if err != nil {
		return err
	}

	targetDir := filepath.Join(s.baseDir, inst.RepositoryID, inst.ID)
	return s.gitClient.Pull(ctx, targetDir)
}

func (s *ClonePoolService) GetStatus(ctx context.Context, instanceID string) (*git.RepoStatus, error) {
	inst, err := s.instanceRepo.FindByID(instanceID)
	if err != nil {
		return nil, err
	}

	targetDir := filepath.Join(s.baseDir, inst.RepositoryID, inst.ID)
	return s.gitClient.Status(ctx, targetDir)
}
