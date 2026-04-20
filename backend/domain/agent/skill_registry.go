package agent

import (
	"sync"
	"time"
)

// ExecutionContext Skill 执行上下文
type ExecutionContext struct {
	TaskID      string
	ProjectID   string
	StartTime   time.Time
	Timeout     time.Duration
	Environment map[string]interface{}
}

// ExecutionResult Skill 执行结果
type ExecutionResult struct {
	Success    bool
	Output     interface{}
	Error      error
	Duration   time.Duration
	Metrics    map[string]interface{}
}

// SkillExecutor 定义 Skill 执行接口
type SkillExecutor interface {
	Execute(ctx *ExecutionContext, params map[string]interface{}) (*ExecutionResult, error)
	Validate(params map[string]interface{}) error
	GetName() string
	GetDescription() string
}

// SkillRegistry Skill 注册表，负责 Skill 的注册、发现和管理
type SkillRegistry struct {
	skills map[string]SkillExecutor
	mu     sync.RWMutex
}

// NewSkillRegistry 创建新的 Skill 注册表
func NewSkillRegistry() *SkillRegistry {
	return &SkillRegistry{
		skills: make(map[string]SkillExecutor),
	}
}

// Register 注册 Skill
func (r *SkillRegistry) Register(name string, skill SkillExecutor) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.skills[name] = skill
}

// Get 根据 name 获取 Skill
func (r *SkillRegistry) Get(name string) (SkillExecutor, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	skill, ok := r.skills[name]
	return skill, ok
}

// List 列出所有 Skill 名称
func (r *SkillRegistry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.skills))
	for name := range r.skills {
		names = append(names, name)
	}
	return names
}

// GetSkillInfo 获取 Skill 信息
func (r *SkillRegistry) GetSkillInfo(name string) *SkillInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()
	skill, ok := r.skills[name]
	if !ok {
		return nil
	}
	return &SkillInfo{
		Name:        skill.GetName(),
		Description: skill.GetDescription(),
	}
}

// SkillInfo Skill 基本信息
type SkillInfo struct {
	Name        string
	Description string
}
