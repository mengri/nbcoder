package airuntime

import (
	"github.com/mengri/nbcoder/domain/airuntime"
	"github.com/mengri/nbcoder/domain/event"
)

type AIRuntimeService struct {
	providerRepo airuntime.ProviderRepo
	chainRepo    airuntime.ChainRepo
	callLogRepo  airuntime.CallLogRepo
	registry     *airuntime.ProviderRegistry
	eventBus     event.EventBus
}

func NewAIRuntimeService(
	providerRepo airuntime.ProviderRepo,
	chainRepo airuntime.ChainRepo,
	callLogRepo airuntime.CallLogRepo,
	registry *airuntime.ProviderRegistry,
	eventBus event.EventBus,
) *AIRuntimeService {
	return &AIRuntimeService{
		providerRepo: providerRepo,
		chainRepo:    chainRepo,
		callLogRepo:  callLogRepo,
		registry:     registry,
		eventBus:     eventBus,
	}
}

func (s *AIRuntimeService) RegisterProvider(provider *airuntime.Provider) error {
	s.registry.Register(provider)
	return s.providerRepo.Save(provider)
}

func (s *AIRuntimeService) GetProvider(id string) (*airuntime.Provider, bool) {
	return s.registry.Get(id)
}

func (s *AIRuntimeService) RecordCall(log *airuntime.CallLog) error {
	if err := s.callLogRepo.Save(log); err != nil {
		return err
	}
	evt := event.NewAIRuntimeEvent(generateID(), log.ModelID, event.ModelCalledEvent)
	_ = s.eventBus.Publish(evt)
	return nil
}

func generateID() string {
	return ""
}
