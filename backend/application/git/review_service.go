package git

import (
	"fmt"

	"github.com/mengri/nbcoder/domain/git"
	"github.com/mengri/nbcoder/pkg/uid"
)

type ReviewService struct {
	reviewRepo git.ReviewRepo
	prRepo     git.PullRequestRepo
}

func NewReviewService(reviewRepo git.ReviewRepo, prRepo git.PullRequestRepo) *ReviewService {
	return &ReviewService{
		reviewRepo: reviewRepo,
		prRepo:     prRepo,
	}
}

func (s *ReviewService) CreateReview(pullRequestID, reviewer, projectName string) (*git.Review, error) {
	pr, err := s.prRepo.FindByID(pullRequestID, projectName)
	if err != nil {
		return nil, err
	}
	if pr == nil {
		return nil, fmt.Errorf("pull request not found: %s", pullRequestID)
	}
	if pr.Status != git.PROpen {
		return nil, fmt.Errorf("cannot review PR in status %s", pr.Status)
	}
	review := git.NewReview(uid.NewID(), pullRequestID, reviewer)
	if err := s.reviewRepo.Save(review); err != nil {
		return nil, err
	}
	return review, nil
}

func (s *ReviewService) ApproveReview(pullRequestID, reviewID, comment string) (*git.Review, error) {
	review, err := s.reviewRepo.FindByID(reviewID)
	if err != nil {
		return nil, err
	}
	if review == nil {
		return nil, fmt.Errorf("review not found: %s", reviewID)
	}
	if review.PullRequestID != pullRequestID {
		return nil, fmt.Errorf("review %s does not belong to PR %s", reviewID, pullRequestID)
	}
	review.Approve(comment)
	if err := s.reviewRepo.Update(review); err != nil {
		return nil, err
	}
	return review, nil
}

func (s *ReviewService) RejectReview(pullRequestID, reviewID, comment string) (*git.Review, error) {
	review, err := s.reviewRepo.FindByID(reviewID)
	if err != nil {
		return nil, err
	}
	if review == nil {
		return nil, fmt.Errorf("review not found: %s", reviewID)
	}
	if review.PullRequestID != pullRequestID {
		return nil, fmt.Errorf("review %s does not belong to PR %s", reviewID, pullRequestID)
	}
	review.Reject(comment)
	if err := s.reviewRepo.Update(review); err != nil {
		return nil, err
	}
	return review, nil
}

func (s *ReviewService) GetReviews(pullRequestID string) ([]*git.Review, error) {
	return s.reviewRepo.FindByPullRequestID(pullRequestID)
}
