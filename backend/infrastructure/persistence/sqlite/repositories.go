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
	"gorm.io/gorm"
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

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Project:               NewProjectRepo(db),
		ProjectConfig:         NewProjectConfigRepo(db),
		Standards:             NewStandardsRepo(db),
		DevStandard:           NewDevStandardRepo(db),
		BranchPolicyConfig:    NewBranchPolicyConfigRepo(db),
		ProjectLifecycle:      NewProjectLifecycleRepo(db),
		ConfigChangeLog:       NewConfigChangeLogRepo(db),
		Card:                  NewCardRepo(db),
		CardDependency:        NewCardDependencyRepo(db),
		Pipeline:              NewPipelineRepo(db),
		StageRecord:           NewStageRecordRepo(db),
		Task:                  NewTaskRepo(db),
		AgentExecution:        NewAgentExecutionRepo(db),
		Skill:                 NewSkillRepo(db),
		Provider:              NewProviderRepo(db),
		Model:                 NewModelRepo(db),
		ModelChain:            NewModelChainRepo(db),
		CallLog:               NewCallLogRepo(db),
		Document:              NewDocumentRepo(db),
		DocumentChunk:         NewDocumentChunkRepo(db),
		DocumentIndex:         NewDocumentIndexRepo(db),
		Directory:             NewDirectoryRepo(db),
		CloneInstance:         NewCloneInstanceRepo(db),
		PullRequest:           NewPullRequestRepo(db),
		Notification:          NewNotificationRepo(db),
		Subscription:          NewSubscriptionRepo(db),
		SubscriptionPreference: NewSubscriptionPreferenceRepo(db),
		NotificationTemplate:  NewNotificationTemplateRepo(db),
		NotificationHistory:   NewNotificationHistoryRepo(db),
	}
}
