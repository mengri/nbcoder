package pipeline

import (
	"testing"
)

func TestDefaultStageNames(t *testing.T) {
	if len(DefaultStageNames) != 7 {
		t.Fatalf("expected 7 default stage names, got %d", len(DefaultStageNames))
	}
	expected := []string{"需求分析", "方案设计", "任务拆解", "测试用例", "代码开发", "测试验证", "评审合并"}
	for i, name := range expected {
		if DefaultStageNames[i] != name {
			t.Errorf("stage name[%d]: expected %s, got %s", i, name, DefaultStageNames[i])
		}
	}
}

func TestDefaultStageConfig(t *testing.T) {
	cfg := DefaultStageConfig()
	if !cfg.Enabled {
		t.Error("default config should be enabled")
	}
	if cfg.ReviewMode != ReviewAI {
		t.Errorf("default review mode should be AI, got %s", cfg.ReviewMode)
	}
	if cfg.MaxRetries != 3 {
		t.Errorf("default max retries should be 3, got %d", cfg.MaxRetries)
	}
}

func TestNewDefaultPipeline(t *testing.T) {
	p := NewDefaultPipeline("pipe-1", "card-1")
	if p.ID != "pipe-1" {
		t.Errorf("expected ID pipe-1, got %s", p.ID)
	}
	if p.CardID != "card-1" {
		t.Errorf("expected CardID card-1, got %s", p.CardID)
	}
	if len(p.Stages) != 7 {
		t.Fatalf("expected 7 stages, got %d", len(p.Stages))
	}
	if p.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
	if p.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should not be zero")
	}

	for i, s := range p.Stages {
		if s.Name != DefaultStageNames[i] {
			t.Errorf("stage[%d]: expected name %s, got %s", i, DefaultStageNames[i], s.Name)
		}
		if s.Status != StageNotStarted {
			t.Errorf("stage[%d]: expected status NOT_STARTED, got %s", i, s.Status)
		}
		if !s.Config.Enabled {
			t.Errorf("stage[%d]: should be enabled", i)
		}
		if s.Config.ReviewMode != ReviewAI {
			t.Errorf("stage[%d]: review mode should be AI, got %s", i, s.Config.ReviewMode)
		}
		if s.Config.MaxRetries != 3 {
			t.Errorf("stage[%d]: max retries should be 3, got %d", i, s.Config.MaxRetries)
		}
		expectedID := "pipe-1-stage-1"
		if i == 0 && s.ID != expectedID {
			t.Errorf("stage[0]: expected ID %s, got %s", expectedID, s.ID)
		}
	}
}

func TestGetStageByName(t *testing.T) {
	p := NewDefaultPipeline("pipe-2", "card-2")

	s := p.GetStageByName("需求分析")
	if s == nil {
		t.Fatal("expected to find stage '需求分析'")
	}
	if s.Name != "需求分析" {
		t.Errorf("expected name '需求分析', got %s", s.Name)
	}

	s = p.GetStageByName("评审合并")
	if s == nil {
		t.Fatal("expected to find stage '评审合并'")
	}

	s = p.GetStageByName("不存在")
	if s != nil {
		t.Error("expected nil for non-existent stage name")
	}
}

func TestGetStageRecords(t *testing.T) {
	p := NewDefaultPipeline("pipe-3", "card-3")

	records := p.GetStageRecords("stage-1")
	if len(records) != 0 {
		t.Errorf("expected 0 records, got %d", len(records))
	}

	p.Records = append(p.Records, &StageRecord{ID: "r1", StageID: "stage-1"})
	p.Records = append(p.Records, &StageRecord{ID: "r2", StageID: "stage-1"})
	p.Records = append(p.Records, &StageRecord{ID: "r3", StageID: "stage-2"})

	records = p.GetStageRecords("stage-1")
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
	if records[0].ID != "r1" || records[1].ID != "r2" {
		t.Errorf("unexpected record IDs: %s, %s", records[0].ID, records[1].ID)
	}

	records = p.GetStageRecords("stage-2")
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
}

func TestPipelineIsCompleted(t *testing.T) {
	p := NewDefaultPipeline("pipe-4", "card-4")
	if p.IsCompleted() {
		t.Error("new pipeline should not be completed")
	}
	for _, s := range p.Stages {
		s.Status = StageCompleted
	}
	if !p.IsCompleted() {
		t.Error("pipeline with all completed stages should be completed")
	}
}

func TestPipelineCurrentStage(t *testing.T) {
	p := NewDefaultPipeline("pipe-5", "card-5")
	if p.CurrentStage() != nil {
		t.Error("new pipeline should have no current stage")
	}
	p.Stages[0].Status = StageInProgress
	if p.CurrentStage() != p.Stages[0] {
		t.Error("expected first stage as current")
	}
}

func TestPipelineNextPendingStage(t *testing.T) {
	p := NewDefaultPipeline("pipe-6", "card-6")
	s := p.NextPendingStage()
	if s != p.Stages[0] {
		t.Error("expected first stage as next pending")
	}
	p.Stages[0].Status = StageCompleted
	s = p.NextPendingStage()
	if s != p.Stages[1] {
		t.Error("expected second stage as next pending")
	}
}

func TestReviewResultConstants(t *testing.T) {
	if ReviewApproved != "APPROVED" {
		t.Errorf("expected APPROVED, got %s", ReviewApproved)
	}
	if ReviewRejected != "REJECTED" {
		t.Errorf("expected REJECTED, got %s", ReviewRejected)
	}
}

func TestStageRecordWithReviewFields(t *testing.T) {
	r := &StageRecord{
		ID:           "r1",
		StageID:      "s1",
		Status:       StageCompleted,
		Output:       "some output",
		ReviewResult: ReviewApproved,
		Reviewer:     "reviewer-1",
	}
	if r.ReviewResult != ReviewApproved {
		t.Errorf("expected APPROVED, got %s", r.ReviewResult)
	}
	if r.Reviewer != "reviewer-1" {
		t.Errorf("expected reviewer-1, got %s", r.Reviewer)
	}
}
