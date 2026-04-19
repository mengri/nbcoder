package git

type PullRequestStatus string

const (
	PROpen   PullRequestStatus = "OPEN"
	PRClosed PullRequestStatus = "CLOSED"
	PRMerged PullRequestStatus = "MERGED"
)

type PullRequest struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description,omitempty"`
	SourceBranch string           `json:"source_branch"`
	TargetBranch string           `json:"target_branch"`
	Status      PullRequestStatus `json:"status"`
}

func NewPullRequest(id, title, source, target string) *PullRequest {
	return &PullRequest{
		ID:           id,
		Title:        title,
		SourceBranch: source,
		TargetBranch: target,
		Status:       PROpen,
	}
}

func (pr *PullRequest) Close() {
	pr.Status = PRClosed
}

func (pr *PullRequest) Merge() {
	pr.Status = PRMerged
}
