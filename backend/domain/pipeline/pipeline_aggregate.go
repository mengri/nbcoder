package pipeline

import "fmt"

type PipelineAggregate struct {
	Pipeline *Pipeline
}

func NewPipelineAggregate(pipeline *Pipeline) *PipelineAggregate {
	return &PipelineAggregate{
		Pipeline: pipeline,
	}
}

func (pa *PipelineAggregate) StartNextStage() error {
	stage := pa.Pipeline.NextPendingStage()
	if stage == nil {
		return fmt.Errorf("no pending stage to start")
	}
	if !stage.Config.Enabled {
		stage.Status = StageCompleted
		return pa.StartNextStage()
	}
	return stage.Start()
}

func (pa *PipelineAggregate) CompleteCurrentStage() error {
	stage := pa.Pipeline.CurrentStage()
	if stage == nil {
		return fmt.Errorf("no active stage to complete")
	}
	return stage.Complete()
}

func (pa *PipelineAggregate) FailCurrentStage(reason string) error {
	stage := pa.Pipeline.CurrentStage()
	if stage == nil {
		return fmt.Errorf("no active stage to fail")
	}
	return stage.Fail(reason)
}

func (pa *PipelineAggregate) RequestReview() error {
	stage := pa.Pipeline.CurrentStage()
	if stage == nil {
		return fmt.Errorf("no active stage to request review")
	}
	return stage.RequireReview()
}
