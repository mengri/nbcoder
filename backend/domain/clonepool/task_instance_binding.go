package clonepool
// task_instance_binding.go
// 任务自动获取与归还克隆实例
package clonepool

import (
	"github.com/mengri/nbcoder/backend/domain"
)

type TaskInstanceBinding struct {
	TaskID      string
	InstanceID  string
}

// 获取可用实例并绑定到任务
func AssignInstanceToTask(pool *ClonePool, task *domain.AgentTask) *CloneInstance {
	inst := pool.Allocate()
	if inst != nil {
		return inst
	}
	return nil
}

// 归还实例并解绑
func ReleaseInstanceFromTask(pool *ClonePool, instanceID string) {
	pool.Release(instanceID)
}
