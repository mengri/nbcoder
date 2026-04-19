# 流水线阶段状态机实施计划

> 任务编号: 8. 流水线阶段状态机 [P1]
> 预估工作量: 3-5 天
> 优先级: P1

---

## 1. 概述

实现完整的流水线阶段状态机，管理流水线各阶段的状态转换，包括阶段依赖、审核规则、超时处理和失败重试等功能。

---

## 2. 架构设计

### 2.1 状态机架构
```
┌─────────────────────────────────────────────────────────┐
│              State Machine Engine                        │
├─────────────────────────────────────────────────────────┤
│  ├── State Definition (State Configuration)             │
│  ├── Transition Engine (Transition Logic)               │
│  ├── Validation (Transition Validation)                │
│  ├── Interceptor (Transition Interceptor)               │
│  └── Logger (Transition Logging)                        │
└─────────────────────────────────────────────────────────┘
                           │
┌─────────────────────────────────────────────────────────┐
│              Rule Engine                                │
│  ├── Dependency Rules                                  │
│  ├── Approval Rules                                    │
│  ├── Timeout Rules                                     │
│  ├── Retry Rules                                       │
│  └── Skip Rules                                         │
└─────────────────────────────────────────────────────────┘
```

---

## 3. 实施步骤

### 3.1 Phase 1: 状态机引擎 (1-2 天)

#### 1.1 状态定义
**文件**: `backend/domain/pipeline/stage_state.go`

```go
type StageState string

const (
    StagePending     StageState = "pending"
    StageReady       StageState = "ready"
    StageRunning     StageState = "running"
    StageWaiting     StageState = "waiting"      // waiting for approval
    StageCompleted   StageState = "completed"
    StageFailed      StageState = "failed"
    StageSkipped     StageState = "skipped"
    StageCancelled   StageState = "cancelled"
)

type StateDefinition struct {
    State      StageState
    Label      string
    Description string
    Transitions []Transition
    Actions    []StateAction
}

type Transition struct {
    To         StageState
    Condition  TransitionCondition
    Actions    []TransitionAction
}

type StateAction struct {
    Type    string
    Handler interface{}
}

type TransitionCondition func(stage *pipeline.StageRecord) bool

type TransitionAction func(stage *pipeline.StageRecord) error
```

#### 1.2 状态转换引擎
**文件**: `backend/application/pipeline/transition_engine.go`

```go
type TransitionEngine struct {
    stateMap    map[StageState]*StateDefinition
    validator   *TransitionValidator
    interceptor *TransitionInterceptor
    logger      *TransitionLogger
}

type TransitionRequest struct {
    StageID    string
    From       StageState
    To         StageState
    Reason     string
    Context    map[string]interface{}
}

type TransitionResult struct {
    Success    bool
    Stage      *pipeline.StageRecord
    Error      error
    Actions    []string
    Duration   time.Duration
}

func NewTransitionEngine() *TransitionEngine
func (e *TransitionEngine) AddDefinition(def *StateDefinition) error
func (e *TransitionEngine) Transition(ctx context.Context, req *TransitionRequest) (*TransitionResult, error)
func (e *TransitionEngine) CanTransition(from, to StageState) bool
func (e *TransitionEngine) GetValidTransitions(from StageState) []StageState
```

#### 1.3 状态转换校验
**文件**: `backend/application/pipeline/transition_validator.go`

```go
type TransitionValidator struct {
    rules []ValidationRule
}

type ValidationRule interface {
    Validate(stage *pipeline.StageRecord, to StageState) error
    Name() string
}

func NewTransitionValidator() *TransitionValidator
func (v *TransitionValidator) AddRule(rule ValidationRule) error
func (v *TransitionValidator) Validate(stage *pipeline.StageRecord, to StageState) error
func (v *TransitionValidator) RemoveRule(name string) error
```

#### 1.4 状态转换拦截器
**文件**: `backend/application/pipeline/transition_interceptor.go`

```go
type TransitionInterceptor struct {
    beforeHooks []TransitionHook
    afterHooks  []TransitionHook
}

type TransitionHook func(ctx context.Context, stage *pipeline.StageRecord, from, to StageState) error

type InterceptorConfig struct {
    BeforeHooks []TransitionHook
    AfterHooks  []TransitionHook
}

func NewTransitionInterceptor() *TransitionInterceptor
func (i *TransitionInterceptor) AddBeforeHook(hook TransitionHook)
func (i *TransitionInterceptor) AddAfterHook(hook TransitionHook)
func (i *TransitionInterceptor) ExecuteBefore(ctx context.Context, stage *pipeline.StageRecord, from, to StageState) error
func (i *TransitionInterceptor) ExecuteAfter(ctx context.Context, stage *pipeline.StageRecord, from, to StageState) error
```

#### 1.5 状态转换日志
**文件**: `backend/application/pipeline/transition_logger.go`

```go
type TransitionLogger struct {
    repo   pipeline.StageRecordRepository
    logger *log.Logger
}

type TransitionLog struct {
    ID         string
    StageID    string
    From       StageState
    To         StageState
    Timestamp  time.Time
    Duration   time.Duration
    Reason     string
    Success    bool
    Error      string
    Context    map[string]interface{}
}

func NewTransitionLogger(repo pipeline.StageRecordRepository, logger *log.Logger) *TransitionLogger
func (l *TransitionLogger) Log(log *TransitionLog) error
func (l *TransitionLogger) GetHistory(stageID string, limit int) ([]*TransitionLog, error)
func (l *TransitionLogger) GetHistoryByTime(stageID string, start, end time.Time) ([]*TransitionLog, error)
```

---

### 3.2 Phase 2: 状态转换规则 (1 天)

#### 2.1 阶段依赖规则
**文件**: `backend/application/pipeline/dependency_rules.go`

```go
type DependencyRule struct {
    name     string
    required []string  // required stage names
}

func NewDependencyRule(name string, required []string) *DependencyRule
func (r *DependencyRule) Validate(stage *pipeline.StageRecord, to StageState) error
func (r *DependencyRule) Name() string
func (r *DependencyRule) CheckDependencies(stage *pipeline.StageRecord) (bool, error)
```

#### 2.2 审核规则
**文件**: `backend/application/pipeline/approval_rules.go`

```go
type ApprovalRule struct {
    name          string
    requiredApprovers int
    approvers     []string
}

func NewApprovalRule(name string, requiredApprovers int) *ApprovalRule
func (r *ApprovalRule) Validate(stage *pipeline.StageRecord, to StageState) error
func (r *ApprovalRule) Name() string
func (r *ApprovalRule) CheckApproval(stage *pipeline.StageRecord) (bool, error)
func (r *ApprovalRule) AddApprover(approver string) error
func (r *ApprovalRule) RemoveApprover(approver string) error
```

#### 2.3 超时规则
**文件**: `backend/application/pipeline/timeout_rules.go`

```go
type TimeoutRule struct {
    name         string
    timeout      time.Duration
    action       TimeoutAction
}

type TimeoutAction string
const (
    TimeoutActionFail    TimeoutAction = "fail"
    TimeoutActionSkip    TimeoutAction = "skip"
    TimeoutActionRetry   TimeoutAction = "retry"
    TimeoutActionNotify TimeoutAction = "notify"
)

func NewTimeoutRule(name string, timeout time.Duration, action TimeoutAction) *TimeoutRule
func (r *TimeoutRule) Validate(stage *pipeline.StageRecord, to StageState) error
func (r *TimeoutRule) Name() string
func (r *TimeoutRule) CheckTimeout(stage *pipeline.StageRecord) (bool, error)
func (r *TimeoutRule) ExecuteTimeoutAction(stage *pipeline.StageRecord) error
```

#### 2.4 失败重试规则
**文件**: `backend/application/pipeline/retry_rules.go`

```go
type RetryRule struct {
    name         string
    maxRetries   int
    retryDelay   time.Duration
    backoff      float64
    retryableErrors []string
}

func NewRetryRule(name string, maxRetries int, retryDelay time.Duration) *RetryRule
func (r *RetryRule) Validate(stage *pipeline.StageRecord, to StageState) error
func (r *RetryRule) Name() string
func (r *RetryRule) ShouldRetry(stage *pipeline.StageRecord) bool
func (r *RetryRule) GetRetryDelay(attempt int) time.Duration
func (r *RetryRule) IsRetryableError(err error) bool
```

#### 2.5 跳过规则
**文件**: `backend/application/pipeline/skip_rules.go`

```go
type SkipRule struct {
    name      string
    condition SkipCondition
}

type SkipCondition func(stage *pipeline.StageRecord) bool

func NewSkipRule(name string, condition SkipCondition) *SkipRule
func (r *SkipRule) Validate(stage *pipeline.StageRecord, to StageState) error
func (r *SkipRule) Name() string
func (r *SkipRule) ShouldSkip(stage *pipeline.StageRecord) bool
```

---

### 3.3 Phase 3: 阶段依赖管理 (1 天)

#### 3.1 依赖关系定义
**文件**: `backend/domain/pipeline/dependency.go`

```go
type StageDependency struct {
    StageID       string
    DependsOn     string
    Type          DependencyType
    Condition     DependencyCondition
}

type DependencyType string
const (
    DependencyTypeFinish   DependencyType = "finish"
    DependencyTypeSuccess DependencyType = "success"
    DependencyTypeOutput  DependencyType = "output"
)

type DependencyCondition map[string]interface{}

func NewStageDependency(stageID, dependsOn string, depType DependencyType) *StageDependency
func (d *StageDependency) IsSatisfied(stage *pipeline.StageRecord) (bool, error)
```

#### 3.2 依赖检查机制
**文件**: `backend/application/pipeline/dependency_checker.go`

```go
type DependencyChecker struct {
    repo pipeline.StageRecordRepository
}

type DependencyCheckResult struct {
    Satisfied bool
    Reason    string
    BlockedBy []string
}

func NewDependencyChecker(repo pipeline.StageRecordRepository) *DependencyChecker
func (c *DependencyChecker) Check(stage *pipeline.StageRecord) (*DependencyCheckResult, error)
func (c *DependencyChecker) CheckAll(stage *pipeline.StageRecord) []*DependencyCheckResult
func (c *DependencyChecker) GetBlockedStages(pipelineID string) ([]string, error)
```

#### 3.3 依赖可视化
**文件**: `backend/application/pipeline/dependency_visualizer.go`

```go
type DependencyVisualizer struct {
    repo pipeline.StageRecordRepository
}

type DependencyGraph struct {
    Nodes []*GraphVertex
    Edges []*GraphEdge
}

type GraphVertex struct {
    ID     string
    Label  string
    State  StageState
}

type GraphEdge struct {
    From string
    To   string
    Type DependencyType
}

func NewDependencyVisualizer(repo pipeline.StageRecordRepository) *DependencyVisualizer
func (v *DependencyVisualizer) BuildGraph(pipelineID string) (*DependencyGraph, error)
func (v *DependencyVisualizer) GenerateDot(graph *DependencyGraph) string
func (v *DependencyVisualizer) GenerateMermaid(graph *DependencyGraph) string
```

#### 3.4 依赖冲突处理
**文件**: `backend/application/pipeline/dependency_resolver.go`

```go
type DependencyResolver struct {
    repo pipeline.StageRecordRepository
}

type ConflictResolution string
const (
    ResolutionFail    ConflictResolution = "fail"
    ResolutionSkip    ConflictResolution = "skip"
    ResolutionRetry   ConflictResolution = "retry"
    ResolutionMerge   ConflictResolution = "merge"
)

func NewDependencyResolver(repo pipeline.StageRecordRepository) *DependencyResolver
func (r *DependencyResolver) ResolveConflicts(stage *pipeline.StageRecord, resolution ConflictResolution) error
func (r *DependencyResolver) DetectConflicts(stage *pipeline.StageRecord) ([]string, error)
func (r *DependencyResolver) SuggestResolution(stage *pipeline.StageRecord) (ConflictResolution, error)
```

#### 3.5 依赖优化建议
**文件**: `backend/application/pipeline/dependency_optimizer.go`

```go
type DependencyOptimizer struct {
    repo pipeline.StageRecordRepository
}

type OptimizationSuggestion struct {
    Type        string
    Description string
    Impact      string
    Priority    int
}

func NewDependencyOptimizer(repo pipeline.StageRecordRepository) *DependencyOptimizer
func (o *DependencyOptimizer) Analyze(pipelineID string) ([]*OptimizationSuggestion, error)
func (o *DependencyOptimizer) DetectCircularDependencies(pipelineID string) ([]string, error)
func (o *DependencyOptimizer) DetectRedundantDependencies(pipelineID string) ([]string, error)
func (o *DependencyOptimizer) SuggestParallelization(pipelineID string) ([]string, error)
```

---

## 4. 时间安排

| 阶段 | 工作内容 | 预计时间 |
|------|----------|----------|
| Phase 1 | 状态机引擎 | 1-2 天 |
| Phase 2 | 状态转换规则 | 1 天 |
| Phase 3 | 阶段依赖管理 | 1 天 |

**总计**: 3-5 天

---

## 5. 验收标准

### 5.1 功能验收
- [ ] 状态定义完整
- [ ] 状态转换引擎正常
- [ ] 状态转换校验正确
- [ ] 拦截器机制正常
- [ ] 状态转换日志完整
- [ ] 阶段依赖规则正常
- [ ] 审核规则正常
- [ ] 超时规则正常
- [ ] 失败重试规则正常
- [ ] 跳过规则正常
- [ ] 依赖检查准确
- [ ] 依赖可视化正常
- [ ] 依赖冲突处理正确
- [ ] 依赖优化建议合理

### 5.2 性能验收
- [ ] 状态转换响应时间 < 100ms
- [ ] 依赖检查响应时间 < 500ms
- [ ] 状态机性能稳定

### 5.3 质量验收
- [ ] 单元测试覆盖率 ≥ 80%
- [ ] 集成测试覆盖主要场景
- [ ] 代码审查通过
- [ ] 文档完整

---

## 6. 风险与缓解

### 6.1 技术风险
**风险**: 状态转换死锁
**缓解**:
- 实现状态转换超时
- 提供状态转换回滚
- 监控状态转换时间

### 6.2 依赖风险
**风险**: 循环依赖
**缓解**:
- 实现循环依赖检测
- 提供依赖冲突解决
- 优化依赖结构

### 6.3 数据风险
**风险**: 状态不一致
**缓解**:
- 使用事务保证一致性
- 实现状态恢复机制
- 记录详细转换日志
