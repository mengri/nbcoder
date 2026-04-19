# Agent 调度与 Skill 封装实施计划

> 任务编号: 7. Agent 调度与 Skill 封装 [P1]
> 预估工作量: 5-7 天
> 优先级: P1

---

## 1. 概述

实现 Agent 任务调度系统和 Skill 执行封装，提供高效的任务分配、执行和监控能力，支持优先级调度、依赖处理和资源管理。

---

## 2. 架构设计

### 2.1 调度架构
```
┌─────────────────────────────────────────────────────────┐
│              Task Queue Layer                           │
│  Task Queue -> Priority Queue -> Dependency Queue       │
└─────────────────────────────────────────────────────────┘
                           │
┌─────────────────────────────────────────────────────────┐
│              Scheduler Layer                           │
├─────────────────────────────────────────────────────────┤
│  ├── Task Scheduler (Priority Scheduling)              │
│  ├── Dependency Resolver (Dependency Handling)         │
│  ├── Resource Allocator (Resource Management)          │
│  └── State Tracker (State Tracking)                    │
└─────────────────────────────────────────────────────────┘
                           │
┌─────────────────────────────────────────────────────────┐
│              Executor Layer                            │
│  Agent -> Skill Engine -> Skill Execution               │
└─────────────────────────────────────────────────────────┘
                           │
┌─────────────────────────────────────────────────────────┐
│              Monitoring Layer                          │
│  Task Monitor -> Skill Monitor -> Performance Monitor   │
└─────────────────────────────────────────────────────────┘
```

---

## 3. 实施步骤

### 3.1 Phase 1: 任务调度器 (2 天)

#### 1.1 任务队列管理
**文件**: `backend/application/agent/task_queue.go`

```go
type TaskQueue struct {
    priorityQueue *PriorityQueue
    waitingQueue  chan *agent.Task
    runningMap    map[string]*agent.Task
    completedMap  map[string]*agent.Task
    maxRunning    int
}

type PriorityQueue struct {
    tasks []*agent.Task
    cmp   func(a, b *agent.Task) bool
}

func NewTaskQueue(maxRunning int) *TaskQueue
func (q *TaskQueue) Push(task *agent.Task) error
func (q *TaskQueue) Pop() (*agent.Task, error)
func (q *TaskQueue) Peek() (*agent.Task, error)
func (q *TaskQueue) Update(task *agent.Task) error
func (q *TaskQueue) Remove(taskID string) error
func (q *TaskQueue) GetStatus() *QueueStatus
```

#### 1.2 优先级调度
**文件**: `backend/application/agent/priority_scheduler.go`

```go
type PriorityScheduler struct {
    queue       *TaskQueue
    policy      SchedulingPolicy
    workers     []*Worker
}

type SchedulingPolicy string
const (
    PolicyFIFO        SchedulingPolicy = "fifo"
    PolicyPriority    SchedulingPolicy = "priority"
    PolicyRoundRobin  SchedulingPolicy = "round-robin"
    PolicyWeighted    SchedulingPolicy = "weighted"
)

type Worker struct {
    ID        string
    Task      *agent.Task
    Status    WorkerStatus
    StartTime time.Time
}

type WorkerStatus string
const (
    WorkerIdle    WorkerStatus = "idle"
    WorkerBusy    WorkerStatus = "busy"
    WorkerPaused  WorkerStatus = "paused"
)

func NewPriorityScheduler(queue *TaskQueue, policy SchedulingPolicy, numWorkers int) *PriorityScheduler
func (s *PriorityScheduler) Start(ctx context.Context)
func (s *PriorityScheduler) Stop()
func (s *PriorityScheduler) Schedule(task *agent.Task) error
func (s *PriorityScheduler) GetNextTask() (*agent.Task, error)
```

#### 1.3 资源分配
**文件**: `backend/application/agent/resource_allocator.go`

```go
type ResourceAllocator struct {
    totalResources *Resources
    allocatedMap   map[string]*Resources
}

type Resources struct {
    CPU       int
    Memory    int64
    GPU       int
    Network   int
}

type AllocationRequest struct {
    TaskID    string
    Required  *Resources
    Priority  int
}

func NewResourceAllocator(total *Resources) *ResourceAllocator
func (a *ResourceAllocator) Allocate(req *AllocationRequest) (*Resources, error)
func (a *ResourceAllocator) Release(taskID string) error
func (a *ResourceAllocator) CanAllocate(required *Resources) bool
func (a *ResourceAllocator) GetAvailable() *Resources
```

#### 1.4 任务依赖处理
**文件**: `backend/application/agent/dependency_resolver.go`

```go
type DependencyResolver struct {
    taskRepo agent.TaskRepository
    graph    *TaskGraph
}

type TaskGraph struct {
    nodes map[string]*TaskNode
    edges map[string][]string  // task_id -> dependent_task_ids
}

type TaskNode struct {
    Task       *agent.Task
    Dependents []string
    DependedBy []string
}

func NewDependencyResolver(taskRepo agent.TaskRepository) *DependencyResolver
func (r *DependencyResolver) Resolve(task *agent.Task) ([]*agent.Task, error)
func (r *DependencyResolver) CheckDependencies(task *agent.Task) (bool, error)
func (r *DependencyResolver) BuildGraph(tasks []*agent.Task) *TaskGraph
func (r *DependencyResolver) TopologicalSort() ([]string, error)
```

#### 1.5 任务状态跟踪
**文件**: `backend/application/agent/state_tracker.go`

```go
type StateTracker struct {
    taskStates    map[string]TaskState
    stateHistory  map[string][]StateTransition
    listeners     []StateChangeListener
}

type TaskState string
const (
    StatePending    TaskState = "pending"
    StateReady      TaskState = "ready"
    StateRunning    TaskState = "running"
    StateCompleted  TaskState = "completed"
    StateFailed     TaskState = "failed"
    StateCancelled  TaskState = "cancelled"
)

type StateTransition struct {
    From      TaskState
    To        TaskState
    Timestamp time.Time
    Reason    string
}

type StateChangeListener func(taskID string, from, to TaskState)

func NewStateTracker() *StateTracker
func (t *StateTracker) Transition(taskID string, to TaskState, reason string) error
func (t *StateTracker) GetState(taskID string) TaskState
func (t *StateTracker) GetHistory(taskID string) []StateTransition
func (t *StateTracker) AddListener(listener StateChangeListener)
```

---

### 3.2 Phase 2: Skill 执行引擎 (2 天)

#### 2.1 Skill 注册和发现
**文件**: `backend/domain/agent/skill_registry.go`

```go
type SkillRegistry struct {
    skills map[string]Skill
    groups map[string][]string  // group -> skill_names
}

type Skill interface {
    Name() string
    Group() string
    Version() string
    Description() string
    InputSchema() *Schema
    OutputSchema() *Schema
    Validate(input interface{}) error
    Execute(ctx context.Context, input interface{}) (*SkillResult, error)
}

func NewSkillRegistry() *SkillRegistry
func (r *SkillRegistry) Register(skill Skill) error
func (r *SkillRegistry) Unregister(name string) error
func (r *SkillRegistry) Get(name string) (Skill, bool)
func (r *SkillRegistry) List() []Skill
func (r *SkillRegistry) ListByGroup(group string) []Skill
func (r *SkillRegistry) Discover(taskType string) []Skill
```

#### 2.2 Skill 参数映射
**文件**: `backend/application/agent/skill_mapper.go`

```go
type SkillMapper struct {
    mappings map[string]*MappingRule
}

type MappingRule struct {
    SourceType   string
    TargetSkill  string
    FieldMap     map[string]string  // source_field -> target_field
    Transforms   map[string]Transform
}

type Transform struct {
    Type      string
    Params    map[string]interface{}
}

func NewSkillMapper() *SkillMapper
func (m *SkillMapper) AddMapping(rule *MappingRule) error
func (m *SkillMapper) Map(source interface{}, targetSkill string) (interface{}, error)
func (m *SkillMapper) ApplyTransform(value interface{}, transform Transform) (interface{}, error)
```

#### 2.3 Skill 执行编排
**文件**: `backend/application/agent/skill_orchestrator.go`

```go
type SkillOrchestrator struct {
    registry  *SkillRegistry
    mapper    *SkillMapper
    executor  *SkillExecutor
}

type ExecutionPlan struct {
    Steps      []*ExecutionStep
    Variables  map[string]interface{}
}

type ExecutionStep struct {
    Skill      string
    Input      interface{}
    OutputVar  string
    DependsOn  []string
    Retry      *RetryConfig
}

type RetryConfig struct {
    MaxAttempts int
    Delay       time.Duration
    Backoff     float64
}

func NewSkillOrchestrator(registry *SkillRegistry, mapper *SkillMapper, executor *SkillExecutor) *SkillOrchestrator
func (o *SkillOrchestrator) CreatePlan(task *agent.Task) (*ExecutionPlan, error)
func (o *SkillOrchestrator) ExecutePlan(ctx context.Context, plan *ExecutionPlan) (*ExecutionResult, error)
func (o *SkillOrchestrator) ExecuteStep(ctx context.Context, step *ExecutionStep, vars map[string]interface{}) (interface{}, error)
```

#### 2.4 结果聚合
**文件**: `backend/application/agent/result_aggregator.go`

```go
type ResultAggregator struct{}

type AggregationConfig struct {
    Type        string  // merge, concat, select, transform
    Source      string  // variable name
    Transform   Transform
    Filter      func(interface{}) bool
}

func NewResultAggregator() *ResultAggregator
func (a *ResultAggregator) Aggregate(results []interface{}, config *AggregationConfig) (interface{}, error)
func (a *ResultAggregator) Merge(results []interface{}) (interface{}, error)
func (a *ResultAggregator) Concat(results []interface{}) (interface{}, error)
func (a *ResultAggregator) Select(results []interface{}, selector func(interface{}) bool) (interface{}, error)
```

#### 2.5 错误处理
**文件**: `backend/application/agent/error_handler.go`

```go
type SkillErrorHandler struct {
    strategies map[string]ErrorStrategy
}

type ErrorStrategy string
const (
    StrategyRetry      ErrorStrategy = "retry"
    StrategySkip       ErrorStrategy = "skip"
    StrategyFail       ErrorStrategy = "fail"
    StrategyFallback   ErrorStrategy = "fallback"
)

type ErrorConfig struct {
    Strategy    ErrorStrategy
    MaxRetries  int
    Fallback    interface{}
    OnError     func(error)
}

func NewSkillErrorHandler() *SkillErrorHandler
func (h *SkillErrorHandler) Handle(err error, config *ErrorConfig) error
func (h *SkillErrorHandler) ShouldRetry(err error, attempt int, config *ErrorConfig) bool
func (h *SkillErrorHandler) GetFallback(err error, config *ErrorConfig) (interface{}, error)
```

---

### 3.3 Phase 3: 任务队列 (1 天)

#### 3.1 队列持久化
**文件**: `backend/infrastructure/agent/persistent_queue.go`

```go
type PersistentQueue struct {
    db        *gorm.DB
    tableName string
}

type QueueEntry struct {
    ID        string
    TaskID    string
    Priority  int
    Status    string
    Data      []byte
    CreatedAt time.Time
    UpdatedAt time.Time
}

func NewPersistentQueue(db *gorm.DB, tableName string) *PersistentQueue
func (q *PersistentQueue) Push(task *agent.Task) error
func (q *PersistentQueue) Pop() (*agent.Task, error)
func (q *PersistentQueue) UpdateStatus(taskID, status string) error
func (q *PersistentQueue) GetTasks(status string) ([]*agent.Task, error)
func (q *PersistentQueue) CleanOldTasks(olderThan time.Duration) error
```

#### 3.2 队列监控
**文件**: `backend/application/agent/queue_monitor.go`

```go
type QueueMonitor struct {
    queue      *TaskQueue
    metrics    *QueueMetrics
    collectors []MetricCollector
}

type QueueMetrics struct {
    TotalTasks      int
    PendingTasks    int
    RunningTasks    int
    CompletedTasks  int
    FailedTasks     int
    AvgWaitTime     time.Duration
    AvgExecTime     time.Duration
    Throughput      float64
}

type MetricCollector interface {
    Collect() *QueueMetrics
}

func NewQueueMonitor(queue *TaskQueue) *QueueMonitor
func (m *QueueMonitor) Start(ctx context.Context)
func (m *QueueMonitor) GetMetrics() *QueueMetrics
func (m *QueueMonitor) GetMetricsHistory(duration time.Duration) []*QueueMetrics
```

#### 3.3 队列管理 API
**文件**: `backend/interfaces/http/queue_api.go`

```go
type QueueAPI struct {
    queue      *TaskQueue
    monitor    *QueueMonitor
    scheduler  *PriorityScheduler
}

type QueueStatusResponse struct {
    Total      int           `json:"total"`
    Pending    int           `json:"pending"`
    Running    int           `json:"running"`
    Completed  int           `json:"completed"`
    Failed     int           `json:"failed"`
    Metrics    *QueueMetrics `json:"metrics"`
}

func NewQueueAPI(queue *TaskQueue, monitor *QueueMonitor, scheduler *PriorityScheduler) *QueueAPI
func (api *QueueAPI) GetStatus(c *gin.Context)
func (api *QueueAPI) Pause(c *gin.Context)
func (api *QueueAPI) Resume(c *gin.Context)
func (api *QueueAPI) Clear(c *gin.Context)
func (api *QueueAPI) GetTasks(c *gin.Context)
```

#### 3.4 死信队列
**文件**: `backend/application/agent/dead_letter_queue.go`

```go
type DeadLetterQueue struct {
    tasks  map[string]*DeadLetterTask
    db     *gorm.DB
}

type DeadLetterTask struct {
    Task       *agent.Task
    Error      error
    Retries    int
    FirstFailed time.Time
    LastFailed  time.Time
}

func NewDeadLetterQueue(db *gorm.DB) *DeadLetterQueue
func (q *DeadLetterQueue) Add(task *agent.Task, err error, retries int) error
func (q *DeadLetterQueue) Get(taskID string) (*DeadLetterTask, error)
func (q *DeadLetterQueue) List() []*DeadLetterTask
func (q *DeadLetterQueue) Retry(taskID string) error
func (q *DeadLetterQueue) Delete(taskID string) error
func (q *DeadLetterQueue) CleanOldTasks(olderThan time.Duration) error
```

#### 3.5 队列优先级
**文件**: `backend/application/agent/priority_manager.go`

```go
type PriorityManager struct {
    priorities map[string]int  // task_type -> priority
    rules      []*PriorityRule
}

type PriorityRule struct {
    Condition  func(*agent.Task) bool
    Priority   int
    Reason     string
}

func NewPriorityManager() *PriorityManager
func (m *PriorityManager) SetPriority(taskType string, priority int) error
func (m *PriorityManager) GetPriority(task *agent.Task) int
func (m *PriorityManager) AddRule(rule *PriorityRule) error
func (m *PriorityManager) ApplyRules(task *agent.Task) int
func (m *PriorityManager) BoostPriority(taskID string, amount int) error
```

---

## 4. 时间安排

| 阶段 | 工作内容 | 预计时间 |
|------|----------|----------|
| Phase 1 | 任务调度器 | 2 天 |
| Phase 2 | Skill 执行引擎 | 2 天 |
| Phase 3 | 任务队列 | 1 天 |

**总计**: 5-7 天

---

## 5. 验收标准

### 5.1 功能验收
- [ ] 任务队列正常工作
- [ ] 优先级调度准确
- [ ] 资源分配合理
- [ ] 任务依赖处理正确
- [ ] 任务状态跟踪完整
- [ ] Skill 注册和发现正常
- [ ] Skill 参数映射正确
- [ ] Skill 执行编排有效
- [ ] 结果聚合准确
- [ ] 错误处理完善
- [ ] 队列持久化正常
- [ ] 队列监控准确
- [ ] 死信队列正常
- [ ] 队列优先级正确

### 5.2 性能验收
- [ ] 任务调度延迟 < 100ms
- [ ] Skill 执行响应时间合理
- [ ] 队列吞吐量 ≥ 100 任务/分钟
- [ ] 支持并发执行

### 5.3 质量验收
- [ ] 单元测试覆盖率 ≥ 80%
- [ ] 集成测试覆盖主要场景
- [ ] 代码审查通过
- [ ] 文档完整

---

## 6. 风险与缓解

### 6.1 技术风险
**风险**: 任务调度冲突
**缓解**:
- 实现优先级机制
- 提供任务依赖检查
- 支持任务重试

### 6.2 性能风险
**风险**: 队列积压
**缓解**:
- 实现队列监控
- 动态调整 Worker 数量
- 优化任务执行逻辑

### 6.3 资源风险
**风险**: 资源耗尽
**缓解**:
- 实现资源限制
- 提供资源回收机制
- 监控资源使用情况
