package pipeline
// stage.go
// Pipeline 阶段状态机与执行控制
package pipeline

import "time"

type StageStatus string

const (
	StageNotStarted   StageStatus = "NOT_STARTED"
	StageInProgress   StageStatus = "IN_PROGRESS"
	StageCompleted    StageStatus = "COMPLETED"
	StageFailed       StageStatus = "FAILED"
	StageReviewNeeded StageStatus = "REVIEW_NEEDED"
)

type Stage struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Status    StageStatus `json:"status"`
	StartedAt time.Time   `json:"started_at,omitempty"`
	EndedAt   time.Time   `json:"ended_at,omitempty"`
	Logs      []string    `json:"logs"`
}

func (s *Stage) Start() {
	s.Status = StageInProgress
	s.StartedAt = time.Now().UTC()
	s.Logs = append(s.Logs, "Stage started at "+s.StartedAt.String())
	// TODO: 发布 StageStarted 领域事件
}

func (s *Stage) Complete() {
	s.Status = StageCompleted
	s.EndedAt = time.Now().UTC()
	s.Logs = append(s.Logs, "Stage completed at "+s.EndedAt.String())
	// TODO: 发布 StageCompleted 领域事件
}

func (s *Stage) Fail(reason string) {
	s.Status = StageFailed
	s.EndedAt = time.Now().UTC()
	s.Logs = append(s.Logs, "Stage failed at "+s.EndedAt.String()+": "+reason)
	// TODO: 发布 StageFailed 领域事件
}

func (s *Stage) RequireReview() {
	s.Status = StageReviewNeeded
	s.Logs = append(s.Logs, "Stage requires review at "+time.Now().UTC().String())
	// TODO: 发布 StageReviewRequired 领域事件
}
