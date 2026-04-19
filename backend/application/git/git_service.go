package git

import (
	"fmt"

	"github.com/mengri/nbcoder/domain/git"
	"github.com/mengri/nbcoder/pkg/uid"
)

type GitService struct {
	prRepo    git.PullRequestRepo
	reviewRepo git.ReviewRepo
	descGen   git.DescriptionGenerator
}

func NewGitService(prRepo git.PullRequestRepo) *GitService {
	return &GitService{
		prRepo:  prRepo,
		descGen: &git.DefaultDescriptionGenerator{},
	}
}

func NewGitServiceWithGenerator(prRepo git.PullRequestRepo, descGen git.DescriptionGenerator) *GitService {
	return &GitService{
		prRepo:  prRepo,
		descGen: descGen,
	}
}

func NewGitServiceWithRepos(prRepo git.PullRequestRepo, reviewRepo git.ReviewRepo) *GitService {
	return &GitService{
		prRepo:     prRepo,
		reviewRepo: reviewRepo,
		descGen:    &git.DefaultDescriptionGenerator{},
	}
}

func (s *GitService) CreatePullRequest(title, source, target, projectID, author string) (*git.PullRequest, error) {
	pr := git.NewPullRequest(uid.NewID(), title, source, target)
	pr.ProjectID = projectID
	pr.Author = author
	if err := s.prRepo.Save(pr); err != nil {
		return nil, err
	}
	return pr, nil
}

func (s *GitService) CreatePullRequestWithDescription(title, source, target, projectID, author string, commits []*git.Commit) (*git.PullRequest, error) {
	pr, err := s.CreatePullRequest(title, source, target, projectID, author)
	if err != nil {
		return nil, err
	}
	desc, err := s.descGen.Generate(source, target, commits)
	if err != nil {
		return nil, fmt.Errorf("failed to generate description: %w", err)
	}
	pr.SetGeneratedDesc(desc)
	_ = s.prRepo.Update(pr)
	return pr, nil
}

func (s *GitService) GetPullRequest(id string) (*git.PullRequest, error) {
	return s.prRepo.FindByID(id)
}

func (s *GitService) MergePullRequest(id string) error {
	pr, err := s.prRepo.FindByID(id)
	if err != nil {
		return err
	}
	if pr == nil {
		return fmt.Errorf("pull request not found: %s", id)
	}
	if err := pr.Merge(); err != nil {
		return err
	}
	return s.prRepo.Update(pr)
}

func (s *GitService) ClosePullRequest(id string) error {
	pr, err := s.prRepo.FindByID(id)
	if err != nil {
		return err
	}
	if pr == nil {
		return fmt.Errorf("pull request not found: %s", id)
	}
	if err := pr.Close(); err != nil {
		return err
	}
	return s.prRepo.Update(pr)
}

func (s *GitService) ValidateBranch(policy *git.BranchPolicy, branch string) bool {
	return policy.Validate(branch)
}

func (s *GitService) SquashMergePullRequest(id, commitMsg string) error {
	pr, err := s.prRepo.FindByID(id)
	if err != nil {
		return err
	}
	if pr == nil {
		return fmt.Errorf("pull request not found: %s", id)
	}
	reviews, err := s.reviewRepo.FindByPullRequestID(id)
	if err != nil {
		return err
	}
	if err := pr.SquashMerge(commitMsg, reviews); err != nil {
		return err
	}
	return s.prRepo.Update(pr)
}
