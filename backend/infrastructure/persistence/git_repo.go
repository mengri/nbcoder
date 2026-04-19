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

type InMemoryReviewRepo struct {
	reviews map[string]*git.Review
	mu      sync.RWMutex
}

func NewInMemoryReviewRepo() *InMemoryReviewRepo {
	return &InMemoryReviewRepo{
		reviews: make(map[string]*git.Review),
	}
}

func (r *InMemoryReviewRepo) Save(review *git.Review) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.reviews[review.ID] = review
	return nil
}

func (r *InMemoryReviewRepo) FindByID(id string) (*git.Review, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	review, ok := r.reviews[id]
	if !ok {
		return nil, nil
	}
	return review, nil
}

func (r *InMemoryReviewRepo) FindByPullRequestID(pullRequestID string) ([]*git.Review, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*git.Review
	for _, review := range r.reviews {
		if review.PullRequestID == pullRequestID {
			result = append(result, review)
		}
	}
	return result, nil
}

func (r *InMemoryReviewRepo) Update(review *git.Review) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.reviews[review.ID] = review
	return nil
}
