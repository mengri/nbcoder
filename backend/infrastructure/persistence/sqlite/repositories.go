package sqlite

import (
	"github.com/mengri/nbcoder/domain/agent"
	"github.com/mengri/nbcoder/domain/airuntime"
	"github.com/mengri/nbcoder/domain/clonepool"
	"github.com/mengri/nbcoder/domain/git"
	"github.com/mengri/nbcoder/domain/knowledge"
	"github.com/mengri/nbcoder/domain/notify"
	"github.com/mengri/nbcoder/domain/pipeline"
	"github.com/mengri/nbcoder/domain/project"
	"github.com/mengri/nbcoder/domain/requirement"
)

type Repositories struct {
	Project              project.ProjectRepo
	ProjectConfig        project.ProjectConfigRepo
	Standards            project.StandardsRepo
	DevStandard          project.DevStandardRepo
	BranchPolicyConfig    project.BranchPolicyConfigRepo
	ProjectLifecycle     project.ProjectLifecycleRepo
	ConfigChangeLog      project.ConfigChangeLogRepo
	Card                 requirement.CardRepo
	CardDependency       requirement.CardDependencyRepo
	Pipeline             pipeline.PipelineRepo
	StageRecord          pipeline.StageRecordRepo
	Task                 agent.TaskRepo
	AgentExecution       agent.AgentExecutionRepo
	Skill                agent.SkillRepo
	Provider             airuntime.ProviderRepo
	Model                *ModelRepo
	ModelChain           airuntime.ChainRepo
	CallLog              airuntime.CallLogRepo
	Document             knowledge.DocumentRepo
	DocumentChunk        *DocumentChunkRepo
	DocumentIndex        *DocumentIndexRepo
	Directory            knowledge.DirectoryRepo
	CloneInstance        clonepool.CloneInstanceRepo

	PullRequest          git.PullRequestRepo
	Notification         notify.NotificationRepo
	Subscription         notify.SubscriptionRepo
	SubscriptionPreference notify.SubscriptionPreferenceRepo
	NotificationTemplate notify.NotificationTemplateRepo
	NotificationHistory  notify.NotificationHistoryRepo
}

func NewRepositories(dbProvider DBProvider, projectBaseDir string) *Repositories {
	return &Repositories{
		Project:               NewProjectRepo(dbProvider, projectBaseDir),
		ProjectConfig:         NewProjectConfigRepo(dbProvider),
		Standards:             NewStandardsRepo(dbProvider),
		DevStandard:           NewDevStandardRepo(dbProvider),
		BranchPolicyConfig:    NewBranchPolicyConfigRepo(dbProvider),
		ProjectLifecycle:      NewProjectLifecycleRepo(dbProvider),
		ConfigChangeLog:       NewConfigChangeLogRepo(dbProvider),
		Card:                  NewCardRepo(dbProvider),
		CardDependency:        NewCardDependencyRepo(dbProvider),
		Pipeline:              NewPipelineRepo(dbProvider),
		StageRecord:           NewStageRecordRepo(dbProvider),
		Task:                  NewTaskRepo(dbProvider),
		AgentExecution:        NewAgentExecutionRepo(dbProvider),
		Skill:                 NewSkillRepo(dbProvider),
		Provider:              NewProviderRepo(dbProvider),
		Model:                 NewModelRepo(dbProvider),
		ModelChain:            NewModelChainRepo(dbProvider),
		CallLog:               NewCallLogRepo(dbProvider),
		Document:              NewDocumentRepo(dbProvider),
		DocumentChunk:         NewDocumentChunkRepo(dbProvider),
		DocumentIndex:         NewDocumentIndexRepo(dbProvider),
		Directory:             NewDirectoryRepo(dbProvider),
		CloneInstance:         NewCloneInstanceRepo(dbProvider),
		PullRequest:           NewPullRequestRepo(dbProvider),
		Notification:          NewNotificationRepo(dbProvider),
		Subscription:          NewSubscriptionRepo(dbProvider),
		SubscriptionPreference: NewSubscriptionPreferenceRepo(dbProvider),
		NotificationTemplate:  NewNotificationTemplateRepo(dbProvider),
		NotificationHistory:   NewNotificationHistoryRepo(dbProvider),
	}
}
