package git

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type ShellGitClient struct {
	baseDir string
}

func NewShellGitClient(baseDir string) *ShellGitClient {
	return &ShellGitClient{baseDir: baseDir}
}

func (c *ShellGitClient) runCommand(ctx context.Context, dir string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git %s failed: %w\nOutput: %s", strings.Join(args, " "), err, string(output))
	}

	return string(output), nil
}

func (c *ShellGitClient) Clone(ctx context.Context, repoURL, targetDir string) error {
	parentDir := filepath.Dir(targetDir)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	_, err := c.runCommand(ctx, parentDir, "clone", repoURL, targetDir)
	return err
}

func (c *ShellGitClient) Commit(ctx context.Context, dir, message string) error {
	output, err := c.runCommand(ctx, dir, "status", "--porcelain")
	if err != nil {
		return fmt.Errorf("failed to check status: %w", err)
	}

	if strings.TrimSpace(output) == "" {
		return fmt.Errorf("no changes to commit")
	}

	if _, err := c.runCommand(ctx, dir, "add", "."); err != nil {
		return fmt.Errorf("failed to stage files: %w", err)
	}

	if _, err := c.runCommand(ctx, dir, "commit", "-m", message); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

func (c *ShellGitClient) Push(ctx context.Context, dir string) error {
	_, err := c.runCommand(ctx, dir, "push", "origin", "HEAD")
	return err
}

func (c *ShellGitClient) Pull(ctx context.Context, dir string) error {
	_, err := c.runCommand(ctx, dir, "pull", "origin", "HEAD")
	return err
}

func (c *ShellGitClient) Branch(ctx context.Context, dir, branchName string) error {
	_, err := c.runCommand(ctx, dir, "checkout", "-b", branchName)
	return err
}

func (c *ShellGitClient) Checkout(ctx context.Context, dir, branch string) error {
	_, err := c.runCommand(ctx, dir, "checkout", branch)
	return err
}

func (c *ShellGitClient) Status(ctx context.Context, dir string) (*RepoStatus, error) {
	output, err := c.runCommand(ctx, dir, "branch", "--show-current")
	if err != nil {
		return nil, fmt.Errorf("failed to get current branch: %w", err)
	}
	branch := strings.TrimSpace(output)

	output, err = c.runCommand(ctx, dir, "status", "--porcelain")
	isDirty := strings.TrimSpace(output) != ""

	staged, unstaged, untracked := c.parseStatus(output)

	commits, err := c.getRecentCommits(ctx, dir, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get commits: %w", err)
	}

	return &RepoStatus{
		Branch:    branch,
		IsDirty:   isDirty,
		Staged:    staged,
		Unstaged:  unstaged,
		Untracked: untracked,
		Commits:   commits,
	}, nil
}

func (c *ShellGitClient) parseStatus(output string) ([]string, []string, []string) {
	staged := make([]string, 0)
	unstaged := make([]string, 0)
	untracked := make([]string, 0)

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		if len(line) < 2 {
			continue
		}

		status := line[:2]
		filename := strings.TrimSpace(line[3:])

		switch status {
		case "??":
			untracked = append(untracked, filename)
		case "A ", "M ", "D ", "R ", "C ":
			staged = append(staged, filename)
		case " M", " D", " R", " C", "MM", "AM", "AD":
			unstaged = append(unstaged, filename)
		default:
			if status[0] != ' ' && status[1] == ' ' {
				staged = append(staged, filename)
			} else if status[0] == ' ' && status[1] != ' ' {
				unstaged = append(unstaged, filename)
			}
		}
	}

	return staged, unstaged, untracked
}

func (c *ShellGitClient) getRecentCommits(ctx context.Context, dir string, limit int) ([]*Commit, error) {
	args := []string{"log", "-n", fmt.Sprintf("%d", limit), "--pretty=format:%H|%an|%ai|%s"}
	output, err := c.runCommand(ctx, dir, args...)
	if err != nil {
		return nil, err
	}

	commits := make([]*Commit, 0)
	lines := strings.Split(strings.TrimSpace(output), "\n")

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.SplitN(line, "|", 4)
		if len(parts) != 4 {
			continue
		}

		parsedTime, err := time.Parse("2006-01-02 15:04:05 -0700", parts[2])
		if err != nil {
			parsedTime = time.Now()
		}

		commits = append(commits, &Commit{
			ID:        parts[0],
			Author:    parts[1],
			Timestamp: parsedTime.Format(time.RFC3339),
			Message:   parts[3],
		})
	}

	return commits, nil
}
