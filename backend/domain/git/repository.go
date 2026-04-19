package git

type PullRequestRepo interface {
	Save(pr *PullRequest) error
	FindByID(id string) (*PullRequest, error)
	FindByBranch(sourceBranch string) ([]*PullRequest, error)
	Update(pr *PullRequest) error
}
