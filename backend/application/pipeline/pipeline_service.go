package pipeline

import (
	"github.com/mengri/nbcoder/domain/event"
	"github.com/mengri/nbcoder/domain/pipeline"
)

type PipelineService struct {
	pipelineRepo   pipeline.PipelineRepo
	stageRecordRepo pipeline.StageRecordRepo
	eventBus       event.EventBus
}

func NewPipelineService(
	pipelineRepo pipeline.PipelineRepo,
	stageRecordRepo pipeline.StageRecordRepo,
	eventBus event.EventBus,
) *PipelineService {
	return &PipelineService{
		pipelineRepo:   pipelineRepo,
		stageRecordRepo: stageRecordRepo,
		eventBus:       eventBus,
	}
}

func (s *PipelineService) CreatePipeline(id, cardID string) (*pipeline.PipelineAggregate, error) {
	pl := pipeline.NewPipeline(id, cardID)
	aggregate := pipeline.NewPipelineAggregate(pl)
	if err := s.pipelineRepo.Save(pl); err != nil {
		return nil, err
	}
	return aggregate, nil
}

func (s *PipelineService) StartNextStage(pipelineID string) error {
	pl, err := s.pipelineRepo.FindByID(pipelineID)
	if err != nil {
		return err
	}
	aggregate := pipeline.NewPipelineAggregate(pl)
	if err := aggregate.StartNextStage(); err != nil {
		return err
	}
	return s.pipelineRepo.Update(pl)
}

func (s *PipelineService) CompleteStage(pipelineID string) error {
	pl, err := s.pipelineRepo.FindByID(pipelineID)
	if err != nil {
		return err
	}
	aggregate := pipeline.NewPipelineAggregate(pl)
	if err := aggregate.CompleteCurrentStage(); err != nil {
		return err
	}
	return s.pipelineRepo.Update(pl)
}

func (s *PipelineService) GetPipeline(pipelineID string) (*pipeline.Pipeline, error) {
	return s.pipelineRepo.FindByID(pipelineID)
}
