package git

import (
	"testing"
)

func TestNewReview(t *testing.T) {
	r := NewReview("rev-1", "pr-1", "alice")
	if r.ID != "rev-1" {
		t.Errorf("expected ID rev-1, got %s", r.ID)
	}
	if r.PullRequestID != "pr-1" {
		t.Errorf("expected PullRequestID pr-1, got %s", r.PullRequestID)
	}
	if r.Reviewer != "alice" {
		t.Errorf("expected Reviewer alice, got %s", r.Reviewer)
	}
	if r.Status != ReviewPending {
		t.Errorf("expected status PENDING, got %s", r.Status)
	}
	if r.ReviewedAt != nil {
		t.Error("expected ReviewedAt to be nil for new review")
	}
}

func TestReview_Approve(t *testing.T) {
	r := NewReview("rev-1", "pr-1", "alice")
	r.Approve("looks good")
	if r.Status != ReviewApproved {
		t.Errorf("expected status APPROVED, got %s", r.Status)
	}
	if r.Comment != "looks good" {
		t.Errorf("expected comment 'looks good', got %s", r.Comment)
	}
	if r.ReviewedAt == nil {
		t.Error("expected ReviewedAt to be set")
	}
}

func TestReview_Reject(t *testing.T) {
	r := NewReview("rev-1", "pr-1", "alice")
	r.Reject("needs work")
	if r.Status != ReviewRejected {
		t.Errorf("expected status REJECTED, got %s", r.Status)
	}
	if r.Comment != "needs work" {
		t.Errorf("expected comment 'needs work', got %s", r.Comment)
	}
	if r.ReviewedAt == nil {
		t.Error("expected ReviewedAt to be set")
	}
}

func TestReview_Approve_AlreadyApproved(t *testing.T) {
	r := NewReview("rev-1", "pr-1", "alice")
	r.Approve("first")
	r.Approve("second")
	if r.Comment != "first" {
		t.Errorf("expected comment to remain 'first', got %s", r.Comment)
	}
}

func TestReview_Reject_AlreadyRejected(t *testing.T) {
	r := NewReview("rev-1", "pr-1", "alice")
	r.Reject("first")
	r.Reject("second")
	if r.Comment != "first" {
		t.Errorf("expected comment to remain 'first', got %s", r.Comment)
	}
}

func TestReview_Approve_AlreadyRejected(t *testing.T) {
	r := NewReview("rev-1", "pr-1", "alice")
	r.Reject("bad")
	r.Approve("ok now")
	if r.Status != ReviewRejected {
		t.Errorf("expected status REJECTED, got %s", r.Status)
	}
}

func TestAllReviewsApproved_NoReviews(t *testing.T) {
	err := AllReviewsApproved(nil)
	if err == nil {
		t.Error("expected error for no reviews")
	}
}

func TestAllReviewsApproved_AllApproved(t *testing.T) {
	r1 := NewReview("rev-1", "pr-1", "alice")
	r1.Approve("ok")
	r2 := NewReview("rev-2", "pr-1", "bob")
	r2.Approve("lgtm")
	err := AllReviewsApproved([]*Review{r1, r2})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestAllReviewsApproved_OneRejected(t *testing.T) {
	r1 := NewReview("rev-1", "pr-1", "alice")
	r1.Approve("ok")
	r2 := NewReview("rev-2", "pr-1", "bob")
	r2.Reject("no")
	err := AllReviewsApproved([]*Review{r1, r2})
	if err == nil {
		t.Error("expected error when a review is rejected")
	}
}

func TestAllReviewsApproved_OnePending(t *testing.T) {
	r1 := NewReview("rev-1", "pr-1", "alice")
	r1.Approve("ok")
	r2 := NewReview("rev-2", "pr-1", "bob")
	err := AllReviewsApproved([]*Review{r1, r2})
	if err == nil {
		t.Error("expected error when a review is pending")
	}
}

func TestPullRequest_SquashMerge(t *testing.T) {
	pr := NewPullRequest("pr-1", "Test", "feature/x", "main")
	r1 := NewReview("rev-1", "pr-1", "alice")
	r1.Approve("ok")
	err := pr.SquashMerge("squash: all changes", []*Review{r1})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if pr.Status != PRMerged {
		t.Errorf("expected status MERGED, got %s", pr.Status)
	}
	if pr.SquashCommitMsg != "squash: all changes" {
		t.Errorf("expected squash commit msg, got %s", pr.SquashCommitMsg)
	}
}

func TestPullRequest_SquashMerge_NoReviews(t *testing.T) {
	pr := NewPullRequest("pr-1", "Test", "feature/x", "main")
	err := pr.SquashMerge("squash: all changes", nil)
	if err == nil {
		t.Error("expected error with no reviews")
	}
}

func TestPullRequest_SquashMerge_RejectedReview(t *testing.T) {
	pr := NewPullRequest("pr-1", "Test", "feature/x", "main")
	r1 := NewReview("rev-1", "pr-1", "alice")
	r1.Reject("no")
	err := pr.SquashMerge("squash: all changes", []*Review{r1})
	if err == nil {
		t.Error("expected error with rejected review")
	}
	if pr.Status != PROpen {
		t.Errorf("expected status OPEN, got %s", pr.Status)
	}
}

func TestPullRequest_SquashMerge_PendingReview(t *testing.T) {
	pr := NewPullRequest("pr-1", "Test", "feature/x", "main")
	r1 := NewReview("rev-1", "pr-1", "alice")
	err := pr.SquashMerge("squash: all changes", []*Review{r1})
	if err == nil {
		t.Error("expected error with pending review")
	}
}

func TestPullRequest_SquashMerge_EmptyCommitMsg(t *testing.T) {
	pr := NewPullRequest("pr-1", "Test", "feature/x", "main")
	r1 := NewReview("rev-1", "pr-1", "alice")
	r1.Approve("ok")
	err := pr.SquashMerge("", []*Review{r1})
	if err == nil {
		t.Error("expected error with empty commit message")
	}
}

func TestPullRequest_SquashMerge_NotOpen(t *testing.T) {
	pr := NewPullRequest("pr-1", "Test", "feature/x", "main")
	_ = pr.Close()
	r1 := NewReview("rev-1", "pr-1", "alice")
	r1.Approve("ok")
	err := pr.SquashMerge("squash: all changes", []*Review{r1})
	if err == nil {
		t.Error("expected error squash merging closed PR")
	}
}
