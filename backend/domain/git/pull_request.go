package git

import (
	"fmt"
	"time"
)

type PullRequestStatus string

const (
	PROpen   PullRequestStatus = "OPEN"
	PRClosed PullRequestStatus = "CLOSED"
	PRMerged PullRequestStatus = "MERGED"
)

type PullRequest struct {
	ID             string            `json:"id"`
	Title          string            `json:"title"`
	Description    string            `json:"description,omitempty"`
	SourceBranch   string            `json:"source_branch"`
	TargetBranch   string            `json:"target_branch"`
	Status         PullRequestStatus `json:"status"`
	ProjectName    string            `json:"project_name,omitempty"`
	Author         string            `json:"author,omitempty"`
	GeneratedDesc  string            `json:"generated_desc,omitempty"`
	SquashCommitMsg string           `json:"squash_commit_msg,omitempty"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

func NewPullRequest(id, title, source, target string) *PullRequest {
	now := time.Now().UTC()
	return &PullRequest{
		ID:           id,
		Title:        title,
		SourceBranch: source,
		TargetBranch: target,
		Status:       PROpen,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func (pr *PullRequest) Close() error {
	if pr.Status != PROpen {
		return fmt.Errorf("cannot close PR in status %s", pr.Status)
	}
	pr.Status = PRClosed
	pr.UpdatedAt = time.Now().UTC()
	return nil
}

func (pr *PullRequest) Merge() error {
	if pr.Status != PROpen {
		return fmt.Errorf("cannot merge PR in status %s", pr.Status)
	}
	pr.Status = PRMerged
	pr.UpdatedAt = time.Now().UTC()
	return nil
}

func (pr *PullRequest) SquashMerge(commitMsg string, reviews []*Review) error {
	if pr.Status != PROpen {
		return fmt.Errorf("cannot squash merge PR in status %s", pr.Status)
	}
	if commitMsg == "" {
		return fmt.Errorf("commit message is required for squash merge")
	}
	if err := AllReviewsApproved(reviews); err != nil {
		return fmt.Errorf("squash merge blocked: %w", err)
	}
	pr.Status = PRMerged
	pr.SquashCommitMsg = commitMsg
	pr.UpdatedAt = time.Now().UTC()
	return nil
}

func (pr *PullRequest) SetDescription(desc string) {
	pr.Description = desc
	pr.UpdatedAt = time.Now().UTC()
}

func (pr *PullRequest) SetGeneratedDesc(desc string) {
	pr.GeneratedDesc = desc
	pr.UpdatedAt = time.Now().UTC()
}

type DescriptionGenerator interface {
	Generate(sourceBranch, targetBranch string, commits []*Commit) (string, error)
}

type DefaultDescriptionGenerator struct{}

func (g *DefaultDescriptionGenerator) Generate(sourceBranch, targetBranch string, commits []*Commit) (string, error) {
	desc := fmt.Sprintf("## PR: %s → %s\n\n", sourceBranch, targetBranch)
	if len(commits) == 0 {
		desc += "No commits found.\n"
		return desc, nil
	}
	desc += "### Changes\n\n"
	for _, c := range commits {
		desc += fmt.Sprintf("- %s\n", c.Message)
	}
	return desc, nil
}
