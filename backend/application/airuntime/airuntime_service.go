package airuntime

import (
	"context"
	"fmt"
	"time"

	"github.com/mengri/nbcoder/domain/airuntime"
	"github.com/mengri/nbcoder/domain/event"
	"github.com/mengri/nbcoder/infrastructure/ai"
	"github.com/mengri/nbcoder/pkg/uid"
)

type AIRuntimeService struct {
	providerRepo   airuntime.ProviderRepo
	chainRepo      airuntime.ChainRepo
	callLogRepo    airuntime.CallLogRepo
	registry       *airuntime.ProviderRegistry
	eventBus       event.EventBus
	modelCaller    *airuntime.ModelCaller
	clientFactory  *ai.ClientFactory
	apiKeyResolver func(providerID string) (string, error)
}

func NewAIRuntimeService(
	providerRepo airuntime.ProviderRepo,
	chainRepo airuntime.ChainRepo,
	callLogRepo airuntime.CallLogRepo,
	registry *airuntime.ProviderRegistry,
	eventBus event.EventBus,
	clientFactory *ai.ClientFactory,
	apiKeyResolver func(providerID string) (string, error),
) *AIRuntimeService {
	modelClient := airuntime.NewDefaultModelClient(clientFactory, apiKeyResolver)
	modelCaller := airuntime.NewModelCaller(registry, modelClient)

	return &AIRuntimeService{
		providerRepo:   providerRepo,
		chainRepo:      chainRepo,
		callLogRepo:    callLogRepo,
		registry:       registry,
		eventBus:       eventBus,
		modelCaller:    modelCaller,
		clientFactory:  clientFactory,
		apiKeyResolver: apiKeyResolver,
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
	evt := event.NewAIRuntimeEvent(uid.NewID(), log.ModelID, event.ModelCalledEvent)
	_ = s.eventBus.Publish(evt)
	return nil
}

func (s *AIRuntimeService) CallModel(ctx context.Context, providerID, modelID string, messages []airuntime.Message, agentID string, opts ...airuntime.CallOption) (*airuntime.ModelResponse, error) {
	startTime := time.Now()

	response, err := s.modelCaller.CallModel(ctx, providerID, modelID, messages, opts...)
	if err != nil {
		callLog := &airuntime.CallLog{
			ID:        uid.NewID(),
			ModelID:   modelID,
			AgentID:   agentID,
			Input:     fmt.Sprintf("%+v", messages),
			Output:    "",
			Tokens:    0,
			Timestamp: startTime,
		}
		_ = s.callLogRepo.Save(callLog)

		evt := event.NewAIRuntimeEvent(uid.NewID(), modelID, event.ModelFailedEvent)
		evt.Payload["error"] = err.Error()
		evt.Payload["duration_ms"] = time.Since(startTime).Milliseconds()
		_ = s.eventBus.Publish(evt)

		return nil, err
	}

	callLog := &airuntime.CallLog{
		ID:        uid.NewID(),
		ModelID:   modelID,
		AgentID:   agentID,
		Input:     fmt.Sprintf("%+v", messages),
		Output:    response.Content,
		Tokens:    response.TotalTokens,
		Timestamp: startTime,
	}

	if err := s.callLogRepo.Save(callLog); err != nil {
		return response, nil
	}

	evt := event.NewAIRuntimeEvent(uid.NewID(), modelID, event.ModelCalledEvent)
	evt.Payload["prompt_tokens"] = response.PromptTokens
	evt.Payload["completion_tokens"] = response.CompletionTokens
	evt.Payload["total_tokens"] = response.TotalTokens
	evt.Payload["duration_ms"] = time.Since(startTime).Milliseconds()
	_ = s.eventBus.Publish(evt)

	return response, nil
}
