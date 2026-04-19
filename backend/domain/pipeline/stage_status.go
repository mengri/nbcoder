package pipeline

type StageStatus string

const (
	StageNotStarted   StageStatus = "NOT_STARTED"
	StageInProgress   StageStatus = "IN_PROGRESS"
	StageCompleted    StageStatus = "COMPLETED"
	StageFailed       StageStatus = "FAILED"
	StageReviewNeeded StageStatus = "REVIEW_NEEDED"
)

func (s StageStatus) IsValid() bool {
	switch s {
	case StageNotStarted, StageInProgress, StageCompleted, StageFailed, StageReviewNeeded:
		return true
	}
	return false
}
