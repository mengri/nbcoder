#!/bin/bash

# 仓储文件重构脚本
# 这个脚本会批量重构所有剩余的仓储文件

FILES=(
"task_repo.go:requirement.TaskRepo:agent"
"pipeline_repo.go:pipeline.PipelineRepo:pipeline" 
"agent_execution_repo.go:agent.AgentExecutionRepo:agent"
"knowledge_repo.go:knowledge.DocumentRepo:knowledge"
"knowledge_repo.go:knowledge.DirectoryRepo:knowledge"
"knowledge_chunk_repo.go:*DocumentChunkRepo:knowledge"
"knowledge_chunk_repo.go:*DocumentIndexRepo:knowledge"
"git_clone_repo.go:clonepool.CloneInstanceRepo:clonepool"
"git_clone_repo.go:git.PullRequestRepo:git"
"notify_repo.go:notify.NotificationRepo:notify"
"notify_repo.go:notify.SubscriptionRepo:notify"
"notify_pref_repo.go:notify.SubscriptionPreferenceRepo:notify"
"notify_pref_repo.go:notify.NotificationTemplateRepo:notify"
"notify_pref_repo.go:notify.NotificationHistoryRepo:notify"
"skill_repo.go:agent.SkillRepo:agent"
"airuntime_repo.go:airuntime.ProviderRepo:airuntime"
"airuntime_repo.go:*ModelRepo:airuntime"
"airuntime_chain_repo.go:airuntime.ChainRepo:airuntime"
"airuntime_chain_repo.go:airuntime.CallLogRepo:airuntime"
"stage_record_repo.go:pipeline.StageRecordRepo:pipeline"
)

for entry in "${FILES[@]}"; do
    IFS=':' read -ra PARTS <<< "$entry"
    FILE="${PARTS[0]}"
    REPO_TYPE="${PARTS[1]}"
    DOMAIN="${PARTS[2]}"
    
    echo "Processing: $FILE - $REPO_TYPE"
done
