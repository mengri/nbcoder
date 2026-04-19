// taskbinding.go
// 任务与克隆实例绑定协调（领域层，解耦clonepool包依赖）
package domain

import (
	"github.com/mengri/nbcoder/domain/clonepool"
)

type TaskInstanceBinding struct {
	TaskID     string
	InstanceID string
}

// 获取可用实例并绑定到任务
func AssignInstanceToTask(pool *clonepool.ClonePool, task *AgentTask) *clonepool.CloneInstance {
	inst := pool.Allocate()
	if inst != nil {
		return inst
	}
	return nil
}

// 归还实例并解绑
func ReleaseInstanceFromTask(pool *clonepool.ClonePool, instanceID string) {
	pool.Release(instanceID)
}
