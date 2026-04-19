package git

import (
	"context"
)

type GitClient interface {
	Clone(ctx context.Context, repoURL, targetDir string) error
	Commit(ctx context.Context, dir, message string) error
	Push(ctx context.Context, dir string) error
	Pull(ctx context.Context, dir string) error
	Branch(ctx context.Context, dir, branchName string) error
	Checkout(ctx context.Context, dir, branch string) error
	Status(ctx context.Context, dir string) (*RepoStatus, error)
}

type RepoStatus struct {
	Branch     string      `json:"branch"`
	IsDirty    bool        `json:"is_dirty"`
	Staged     []string    `json:"staged"`
	Unstaged   []string    `json:"unstaged"`
	Untracked  []string    `json:"untracked"`
	Ahead      int         `json:"ahead"`
	Behind     int         `json:"behind"`
	Commits    []*Commit   `json:"commits"`
}

type Commit struct {
	ID        string `json:"id"`
	Message   string `json:"message"`
	Author    string `json:"author"`
	Timestamp string `json:"timestamp"`
}
