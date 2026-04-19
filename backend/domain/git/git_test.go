package git

import (
	"strings"
	"testing"
)

func TestNewPullRequest(t *testing.T) {
	pr := NewPullRequest("pr-1", "Test PR", "feature/x", "main")
	if pr.ID != "pr-1" {
		t.Errorf("expected ID pr-1, got %s", pr.ID)
	}
	if pr.Status != PROpen {
		t.Errorf("expected status OPEN, got %s", pr.Status)
	}
}

func TestPullRequest_Merge(t *testing.T) {
	pr := NewPullRequest("pr-1", "Test PR", "feature/x", "main")
	if err := pr.Merge(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if pr.Status != PRMerged {
		t.Errorf("expected status MERGED, got %s", pr.Status)
	}
}

func TestPullRequest_Merge_NotOpen(t *testing.T) {
	pr := NewPullRequest("pr-1", "Test PR", "feature/x", "main")
	_ = pr.Close()
	if err := pr.Merge(); err == nil {
		t.Error("expected error merging closed PR")
	}
}

func TestPullRequest_Close(t *testing.T) {
	pr := NewPullRequest("pr-1", "Test PR", "feature/x", "main")
	if err := pr.Close(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if pr.Status != PRClosed {
		t.Errorf("expected status CLOSED, got %s", pr.Status)
	}
}

func TestPullRequest_Close_NotOpen(t *testing.T) {
	pr := NewPullRequest("pr-1", "Test PR", "feature/x", "main")
	_ = pr.Merge()
	if err := pr.Close(); err == nil {
		t.Error("expected error closing merged PR")
	}
}

func TestPullRequest_SetGeneratedDesc(t *testing.T) {
	pr := NewPullRequest("pr-1", "Test PR", "feature/x", "main")
	pr.SetGeneratedDesc("auto generated")
	if pr.GeneratedDesc != "auto generated" {
		t.Errorf("expected auto generated, got %s", pr.GeneratedDesc)
	}
}

func TestDefaultDescriptionGenerator(t *testing.T) {
	gen := &DefaultDescriptionGenerator{}
	commits := []*Commit{
		{Hash: "abc123", Message: "feat: add feature"},
		{Hash: "def456", Message: "fix: fix bug"},
	}
	desc, err := gen.Generate("feature/x", "main", commits)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !strings.Contains(desc, "feature/x") || !strings.Contains(desc, "main") {
		t.Errorf("expected description to contain branch names, got %s", desc)
	}
	if !strings.Contains(desc, "add feature") {
		t.Errorf("expected description to contain commit messages, got %s", desc)
	}
}

func TestDefaultDescriptionGenerator_NoCommits(t *testing.T) {
	gen := &DefaultDescriptionGenerator{}
	desc, err := gen.Generate("feature/x", "main", nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !strings.Contains(desc, "No commits") {
		t.Errorf("expected no commits message, got %s", desc)
	}
}

func TestCommitHistory(t *testing.T) {
	ch := &CommitHistory{}
	c1 := &Commit{Hash: "abc", Message: "first"}
	c2 := &Commit{Hash: "def", Message: "second"}
	ch.Add(c1)
	ch.Add(c2)
	if len(ch.Commits) != 2 {
		t.Errorf("expected 2 commits, got %d", len(ch.Commits))
	}
	if ch.Latest().Hash != "def" {
		t.Errorf("expected latest hash def, got %s", ch.Latest().Hash)
	}
}

func TestCommitHistory_Latest_Empty(t *testing.T) {
	ch := &CommitHistory{}
	if ch.Latest() != nil {
		t.Error("expected nil for empty history")
	}
}

func TestBranchPolicy_Validate(t *testing.T) {
	policy := &BranchPolicy{AllowedPattern: `^feature/`}
	if !policy.Validate("feature/A1.1-test") {
		t.Error("expected feature branch to be valid")
	}
	if policy.Validate("main") {
		t.Error("expected main to be invalid for feature-only pattern")
	}
}
