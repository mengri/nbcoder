package git

import (
	"github.com/mengri/nbcoder/domain/git"
)

type GitService struct {
	prRepo git.PullRequestRepo
}

func NewGitService(prRepo git.PullRequestRepo) *GitService {
	return &GitService{
		prRepo: prRepo,
	}
}

func (s *GitService) CreatePullRequest(id, title, source, target string) (*git.PullRequest, error) {
	pr := git.NewPullRequest(id, title, source, target)
	if err := s.prRepo.Save(pr); err != nil {
		return nil, err
	}
	return pr, nil
}

func (s *GitService) MergePullRequest(id string) error {
	pr, err := s.prRepo.FindByID(id)
	if err != nil {
		return err
	}
	pr.Merge()
	return s.prRepo.Update(pr)
}

func (s *GitService) ValidateBranch(policy *git.BranchPolicy, branch string) bool {
	return policy.Validate(branch)
}
