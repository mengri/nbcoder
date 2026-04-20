package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"

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
	"github.com/mengri/nbcoder/infrastructure/database"
	"github.com/mengri/nbcoder/infrastructure/eventbus"
	"github.com/mengri/nbcoder/infrastructure/git"
	"github.com/mengri/nbcoder/infrastructure/persistence/sqlite"
	"github.com/mengri/nbcoder/interfaces/api"
	embeddedWeb "github.com/mengri/nbcoder/web/embedded"
)

func main() {
	eventBus := eventbus.NewInMemoryEventBus()

	dbConfig := database.DefaultDatabaseConfig()
	dbManager, err := database.NewDatabaseManager(dbConfig)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbManager.Close()

	err = database.InitSchema(dbManager.GetDB())
	if err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}

	repos := sqlite.NewRepositories(dbManager.GetDB())
	agentRegistry := agent.NewAgentRegistry()
	agentService := agentApp.NewAgentService(repos.Task, repos.AgentExecution, agentRegistry, eventBus)

	requirementService := requirementApp.NewRequirementService(repos.Card, repos.CardDependency, eventBus)

	pipelineService := pipelineApp.NewPipelineService(repos.Pipeline, repos.StageRecord, eventBus)

	projectService := projectApp.NewProjectService(repos.Project, repos.ProjectConfig, repos.Standards, repos.ConfigChangeLog)

	gitClient := git.NewShellGitClient("/tmp/nbcoder/clones")
	clonePoolService := clonepoolApp.NewClonePoolService(
		repos.CloneInstance, repos.Repository, eventBus, gitClient, "/tmp/nbcoder/clones",
	)

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
		repos.Provider, repos.ModelChain, repos.CallLog, providerRegistry, eventBus,
		clientFactory, apiKeyResolver,
	)

	knowledgeService := knowledgeApp.NewKnowledgeService(repos.Document, repos.Directory, repos.DocumentChunk, repos.DocumentIndex)

	dispatcher := notify.NewChannelDispatcher()
	dispatcher.Register(channel.NewSystemSender())
	dispatcher.Register(channel.NewWebSocketSender())
	dispatcher.Register(channel.NewEmailSender())
	notifyService := notifyApp.NewNotifyService(repos.Notification, repos.Subscription, repos.SubscriptionPreference, nil, dispatcher, eventBus)

	gitService := gitApp.NewGitService(repos.PullRequest)

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

	embedFS := embeddedWeb.WebFS()
	fileServer := http.FileServer(http.FS(embedFS))
	router.NoRoute(func(c *gin.Context) {
		if c.Request.Method == "GET" && c.Request.URL.Path != "/api/" {
			fileServer.ServeHTTP(c.Writer, c.Request)
		}
	})
	log.Println("Embedded web assets loaded")

	log.Println("Starting NBCoder server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
