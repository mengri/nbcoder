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

func (s *GitService) CreatePullRequest(title, source, target, projectName, author string) (*git.PullRequest, error) {
	pr := git.NewPullRequest(uid.NewID(), title, source, target)
	pr.ProjectName = projectName
	pr.Author = author
	if err := s.prRepo.Save(pr); err != nil {
		return nil, err
	}
	return pr, nil
}

func (s *GitService) CreatePullRequestWithDescription(title, source, target, projectName, author string, commits []*git.Commit) (*git.PullRequest, error) {
	pr, err := s.CreatePullRequest(title, source, target, projectName, author)
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

func (s *GitService) GetPullRequest(id, projectName string) (*git.PullRequest, error) {
	return s.prRepo.FindByID(id, projectName)
}

func (s *GitService) MergePullRequest(id, projectName string) error {
	pr, err := s.prRepo.FindByID(id, projectName)
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

func (s *GitService) ClosePullRequest(id, projectName string) error {
	pr, err := s.prRepo.FindByID(id, projectName)
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

func (s *GitService) SquashMergePullRequest(id, commitMsg, projectName string) error {
	pr, err := s.prRepo.FindByID(id, projectName)
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
