package domain

// agent_integration.go
// Agent与Pipeline/ClonePool/AI Runtime集成机制

import (
	"github.com/mengri/nbcoder/domain/airuntime"
	"github.com/mengri/nbcoder/domain/clonepool"
	"github.com/mengri/nbcoder/domain/pipeline"
)

type AgentIntegration struct {
	PipelineRef *pipeline.Stage
	CloneRef    *clonepool.CloneInstance
	ModelChain  *airuntime.ModelChain
}

// 任务自动流转与资源调度示例
func (ai *AgentIntegration) AutoDispatch() error {
	if ai.PipelineRef != nil {
		ai.PipelineRef.Start()
	}
	if ai.CloneRef != nil && ai.CloneRef.Status == clonepool.InstanceIdle {
		ai.CloneRef.Allocate()
	}
	// 路由模型链
	if ai.ModelChain != nil {
		_ = ai.ModelChain.SelectRoute("")
	}
	return nil
}

// 可扩展异常处理与回退机制
