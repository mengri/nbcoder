package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	agentApp "github.com/mengri/nbcoder/application/agent"
	airuntimeApp "github.com/mengri/nbcoder/application/airuntime"
	clonepoolApp "github.com/mengri/nbcoder/application/clonepool"
	gitApp "github.com/mengri/nbcoder/application/git"
	knowledgeApp "github.com/mengri/nbcoder/application/knowledge"
	notifyApp "github.com/mengri/nbcoder/application/notify"
	pipelineApp "github.com/mengri/nbcoder/application/pipeline"
	projectApp "github.com/mengri/nbcoder/application/project"
	requirementApp "github.com/mengri/nbcoder/application/requirement"
	"github.com/mengri/nbcoder/domain/agent"
	"github.com/mengri/nbcoder/domain/airuntime"
	"github.com/mengri/nbcoder/domain/notify"
	"github.com/mengri/nbcoder/infrastructure/ai"
	"github.com/mengri/nbcoder/infrastructure/channel"
	"github.com/mengri/nbcoder/infrastructure/eventbus"
	"github.com/mengri/nbcoder/infrastructure/git"
	"github.com/mengri/nbcoder/infrastructure/persistence"
	"github.com/mengri/nbcoder/interfaces/api"
)

func main() {
	eventBus := eventbus.NewInMemoryEventBus()

	taskRepo := persistence.NewInMemoryTaskRepo()
	executionRepo := persistence.NewInMemoryAgentExecutionRepo()
	agentRegistry := agent.NewAgentRegistry()
	agentService := agentApp.NewAgentService(taskRepo, executionRepo, agentRegistry, eventBus)

	cardRepo := persistence.NewInMemoryCardRepo()
	cardDepRepo := persistence.NewInMemoryCardDependencyRepo()
	requirementService := requirementApp.NewRequirementService(cardRepo, cardDepRepo, eventBus)

	pipelineRepo := persistence.NewInMemoryPipelineRepo()
	stageRecordRepo := persistence.NewInMemoryStageRecordRepo()
	pipelineService := pipelineApp.NewPipelineService(pipelineRepo, stageRecordRepo, eventBus)

	projectRepo := persistence.NewInMemoryProjectRepo()
	projectConfigRepo := persistence.NewInMemoryProjectConfigRepo()
	standardsRepo := persistence.NewInMemoryStandardsRepo()
	configChangeLogRepo := persistence.NewInMemoryConfigChangeLogRepo()
	projectService := projectApp.NewProjectService(projectRepo, projectConfigRepo, standardsRepo, configChangeLogRepo)

	cloneInstanceRepo := persistence.NewInMemoryCloneInstanceRepo()
	repositoryRepo := persistence.NewInMemoryRepositoryRepo()
	gitClient := git.NewShellGitClient("/tmp/nbcoder/clones")
	clonePoolService := clonepoolApp.NewClonePoolService(
		cloneInstanceRepo, repositoryRepo, eventBus, gitClient, "/tmp/nbcoder/clones",
	)

	providerRepo := persistence.NewInMemoryProviderRepo()
	chainRepo := persistence.NewInMemoryChainRepo()
	callLogRepo := persistence.NewInMemoryCallLogRepo()
	providerRegistry := airuntime.NewProviderRegistry()
	clientFactory := ai.NewClientFactory()

	apiKeyResolver := func(providerID string) (string, error) {
		provider, ok := providerRegistry.Get(providerID)
		if !ok {
			return "", fmt.Errorf("provider not found: %s", providerID)
		}
		return provider.APIKeyRef, nil
	}

	aiRuntimeService := airuntimeApp.NewAIRuntimeService(
		providerRepo, chainRepo, callLogRepo, providerRegistry, eventBus,
		clientFactory, apiKeyResolver,
	)

	documentRepo := persistence.NewInMemoryDocumentRepo()
	directoryRepo := persistence.NewInMemoryDirectoryRepo()
	chunkRepo := persistence.NewInMemoryChunkRepo()
	documentIndexRepo := persistence.NewInMemoryDocumentIndexRepo()
	knowledgeService := knowledgeApp.NewKnowledgeService(documentRepo, directoryRepo, chunkRepo, documentIndexRepo)

	notificationRepo := persistence.NewInMemoryNotificationRepo()
	subscriptionRepo := persistence.NewInMemorySubscriptionRepo()
	prefRepo := persistence.NewInMemorySubscriptionPreferenceRepo()
	channelRepo := persistence.NewInMemoryChannelRepo()
	dispatcher := notify.NewChannelDispatcher()
	dispatcher.Register(channel.NewSystemSender())
	dispatcher.Register(channel.NewWebSocketSender())
	dispatcher.Register(channel.NewEmailSender())
	notifyService := notifyApp.NewNotifyService(notificationRepo, subscriptionRepo, prefRepo, channelRepo, dispatcher, eventBus)

	prRepo := persistence.NewInMemoryPullRequestRepo()
	gitService := gitApp.NewGitService(prRepo)

	router := gin.Default()
	apiGroup := router.Group("/api/v1")

	api.NewAgentHandler(agentService).RegisterRoutes(apiGroup)
	api.NewRequirementHandler(requirementService).RegisterRoutes(apiGroup)
	api.NewPipelineHandler(pipelineService).RegisterRoutes(apiGroup)
	api.NewProjectHandler(projectService).RegisterRoutes(apiGroup)
	api.NewClonePoolHandler(clonePoolService).RegisterRoutes(apiGroup)
	api.NewAIRuntimeHandler(aiRuntimeService).RegisterRoutes(apiGroup)
	api.NewKnowledgeHandler(knowledgeService).RegisterRoutes(apiGroup)
	api.NewNotifyHandler(notifyService).RegisterRoutes(apiGroup)
	api.NewGitHandler(gitService).RegisterRoutes(apiGroup)

	log.Println("Starting NBCoder server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
