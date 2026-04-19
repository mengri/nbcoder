package git

import (
	"fmt"
	"time"
)

type ReviewStatus string

const (
	ReviewPending  ReviewStatus = "PENDING"
	ReviewApproved ReviewStatus = "APPROVED"
	ReviewRejected ReviewStatus = "REJECTED"
)

type Review struct {
	ID            string       `json:"id"`
	PullRequestID string       `json:"pull_request_id"`
	Reviewer      string       `json:"reviewer"`
	Status        ReviewStatus `json:"status"`
	Comment       string       `json:"comment,omitempty"`
	ReviewedAt    *time.Time   `json:"reviewed_at,omitempty"`
	CreatedAt     time.Time    `json:"created_at"`
}

func NewReview(id, pullRequestID, reviewer string) *Review {
	now := time.Now().UTC()
	return &Review{
		ID:            id,
		PullRequestID: pullRequestID,
		Reviewer:      reviewer,
		Status:        ReviewPending,
		CreatedAt:     now,
	}
}

func (r *Review) Approve(comment string) {
	if r.Status != ReviewPending {
		return
	}
	r.Status = ReviewApproved
	r.Comment = comment
	now := time.Now().UTC()
	r.ReviewedAt = &now
}

func (r *Review) Reject(comment string) {
	if r.Status != ReviewPending {
		return
	}
	r.Status = ReviewRejected
	r.Comment = comment
	now := time.Now().UTC()
	r.ReviewedAt = &now
}

func AllReviewsApproved(reviews []*Review) error {
	if len(reviews) == 0 {
		return fmt.Errorf("no reviews found, at least one approval required")
	}
	for _, r := range reviews {
		if r.Status == ReviewRejected {
			return fmt.Errorf("review by %s was rejected", r.Reviewer)
		}
		if r.Status != ReviewApproved {
			return fmt.Errorf("review by %s is still pending", r.Reviewer)
		}
	}
	return nil
}
