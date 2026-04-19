package persistence

import (
	"sync"

	"github.com/mengri/nbcoder/domain/git"
)

type InMemoryPullRequestRepo struct {
	prs map[string]*git.PullRequest
	mu  sync.RWMutex
}

func NewInMemoryPullRequestRepo() *InMemoryPullRequestRepo {
	return &InMemoryPullRequestRepo{
		prs: make(map[string]*git.PullRequest),
	}
}

func (r *InMemoryPullRequestRepo) Save(pr *git.PullRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.prs[pr.ID] = pr
	return nil
}

func (r *InMemoryPullRequestRepo) FindByID(id string) (*git.PullRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	pr, ok := r.prs[id]
	if !ok {
		return nil, nil
	}
	return pr, nil
}

func (r *InMemoryPullRequestRepo) FindByBranch(sourceBranch string) ([]*git.PullRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*git.PullRequest
	for _, pr := range r.prs {
		if pr.SourceBranch == sourceBranch {
			result = append(result, pr)
		}
	}
	return result, nil
}

func (r *InMemoryPullRequestRepo) Update(pr *git.PullRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.prs[pr.ID] = pr
	return nil
}
