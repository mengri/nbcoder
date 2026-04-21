package git

type PullRequestRepo interface {
	Save(pr *PullRequest) error
	FindByID(id string, projectName string) (*PullRequest, error)
	FindByProjectName(projectName string) ([]*PullRequest, error)
	FindByBranch(sourceBranch string) ([]*PullRequest, error)
	Update(pr *PullRequest) error
}

type ReviewRepo interface {
	Save(review *Review) error
	FindByID(id string) (*Review, error)
	FindByPullRequestID(pullRequestID string) ([]*Review, error)
	Update(review *Review) error
}
