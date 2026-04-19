# 需求卡片状态机实施计划

> 任务编号: 10. 需求卡片状态机 [P1]
> 预估工作量: 3-5 天
> 优先级: P1

---

## 1. 概述

实现完整的需求卡片状态机，管理需求卡片从创建到完成的全生命周期，包括状态转换规则、验证器、权限控制和通知机制。

---

## 2. 架构设计

### 2.1 状态机架构
```
┌─────────────────────────────────────────────────────────┐
│              Card State Machine                         │
├─────────────────────────────────────────────────────────┤
│  ├── States (All Card States)                          │
│  ├── Transitions (State Transitions)                   │
│  ├── Validator (Transition Validation)                 │
│  ├── Permission Checker (Permission Control)           │
│  └── Notifier (State Change Notification)              │
└─────────────────────────────────────────────────────────┘
                           │
┌─────────────────────────────────────────────────────────┐
│              Business Rules                             │
│  ├── Transition Conditions                             │
│  ├── Side Effects                                      │
│  ├── Permission Rules                                  │
│  └── Notification Rules                                │
└─────────────────────────────────────────────────────────┘
```

---

## 3. 实施步骤

### 3.1 Phase 1: 完整状态机 (1-2 天)

#### 1.1 状态定义
**文件**: `backend/domain/requirement/card_state.go`

```go
type CardState string

const (
    CardDraft      CardState = "draft"
    CardConfirmed  CardState = "confirmed"
    CardPlanned    CardState = "planned"
    CardInProgress CardState = "in_progress"
    CardReview     CardState = "review"
    CardCompleted  CardState = "completed"
    CardRejected   CardState = "rejected"
    CardBlocked    CardState = "blocked"
    CardAbandoned  CardState = "abandoned"
    CardSuperseded CardState = "superseded"
)

type StateDefinition struct {
    State       CardState
    Label       string
    Description string
    AllowedFrom []CardState
    AllowedTo   []CardState
    Actions     []StateAction
}
```

#### 1.2 状态转换表
**文件**: `backend/domain/requirement/transitions.go`

```go
var CardTransitions = map[CardState][]CardState{
    CardDraft:      {CardConfirmed, CardRejected, CardAbandoned},
    CardConfirmed:  {CardPlanned, CardDraft, CardRejected},
    CardPlanned:    {CardInProgress, CardDraft, CardConfirmed},
    CardInProgress: {CardReview, CardBlocked, CardPlanned},
    CardReview:     {CardCompleted, CardInProgress, CardRejected},
    CardCompleted:  {},  // terminal state
    CardRejected:   {CardDraft, CardAbandoned},
    CardBlocked:    {CardInProgress, CardAbandoned},
    CardAbandoned:  {},  // terminal state
    CardSuperseded: {},  // terminal state
}

func CanTransition(from, to CardState) bool {
    allowed, ok := CardTransitions[from]
    if !ok {
        return false
    }
    for _, state := range allowed {
        if state == to {
            return true
        }
    }
    return false
}
```

#### 1.3 状态转换校验
**文件**: `backend/application/requirement/transition_validator.go`

```go
type CardTransitionValidator struct {
    rules []ValidationRule
}

type ValidationRule interface {
    Validate(card *requirement.Card, to CardState) error
    Name() string
}

func NewCardTransitionValidator() *CardTransitionValidator
func (v *CardTransitionValidator) AddRule(rule ValidationRule) error
func (v *CardTransitionValidator) Validate(card *requirement.Card, to CardState) error
```

#### 1.4 状态转换历史
**文件**: `backend/domain/requirement/transition_history.go`

```go
type CardTransitionHistory struct {
    ID         string
    CardID     string
    From       CardState
    To         CardState
    Timestamp  time.Time
    Reason     string
    UserID     string
    Context    map[string]interface{}
}

func NewTransitionHistory(cardID, userID string, from, to CardState, reason string) *CardTransitionHistory
```

#### 1.5 状态转换审计
**文件**: `backend/application/requirement/transition_auditor.go`

```go
type TransitionAuditor struct {
    repo      requirement.CardRepository
    historyRepo CardTransitionHistoryRepository
    logger    *log.Logger
}

type AuditEvent struct {
    EventType   string
    CardID      string
    UserID      string
    StateFrom   CardState
    StateTo     CardState
    Timestamp   time.Time
    Details     map[string]interface{}
}

func NewTransitionAuditor(repo requirement.CardRepository, historyRepo CardTransitionHistoryRepository) *TransitionAuditor
func (a *TransitionAuditor) Audit(card *requirement.Card, to CardState, userID, reason string) error
func (a *TransitionAuditor) GetAuditTrail(cardID string) ([]*AuditEvent, error)
```

---

### 3.2 Phase 2: 状态转换规则 (1 天)

#### 2.1 状态转换条件
**文件**: `backend/application/requirement/transition_conditions.go`

```go
type TransitionCondition interface {
    Evaluate(card *requirement.Card, to CardState) bool
    Description() string
}

type DependencyCondition struct {
    dependencyRepo requirement.CardDependencyRepository
}

func NewDependencyCondition(repo requirement.CardDependencyRepository) *DependencyCondition
func (c *DependencyCondition) Evaluate(card *requirement.Card, to CardState) bool
func (c *DependencyCondition) Description() string

type PipelineCondition struct {
    pipelineRepo pipeline.PipelineRepository
}

func NewPipelineCondition(repo pipeline.PipelineRepository) *PipelineCondition
func (c *PipelineCondition) Evaluate(card *requirement.Card, to CardState) bool
```

#### 2.2 状态转换副作用
**文件**: `backend/application/requirement/transition_side_effects.go`

```go
type TransitionSideEffect interface {
    Execute(card *requirement.Card, from, to CardState) error
    Name() string
}

type NotificationSideEffect struct {
    notifier *notification.Notifier
}

func NewNotificationSideEffect(notifier *notification.Notifier) *NotificationSideEffect
func (e *NotificationSideEffect) Execute(card *requirement.Card, from, to CardState) error

type PipelineTriggerSideEffect struct {
    pipelineService pipeline.Service
}

func NewPipelineTriggerSideEffect(service pipeline.Service) *PipelineTriggerSideEffect
func (e *PipelineTriggerSideEffect) Execute(card *requirement.Card, from, to CardState) error

type DependencyUpdateSideEffect struct {
    dependencyRepo requirement.CardDependencyRepository
}

func NewDependencyUpdateSideEffect(repo requirement.CardDependencyRepository) *DependencyUpdateSideEffect
func (e *DependencyUpdateSideEffect) Execute(card *requirement.Card, from, to CardState) error
```

#### 2.3 状态转换权限
**文件**: `backend/application/requirement/transition_permission.go`

```go
type TransitionPermissionChecker struct {
    roleMgr *rbac.Manager
}

type PermissionRule struct {
    From      CardState
    To        CardState
    Roles     []string
    AllowAll  bool
}

func NewTransitionPermissionChecker(roleMgr *rbac.Manager) *TransitionPermissionChecker
func (c *TransitionPermissionChecker) Check(card *requirement.Card, to CardState, userID string) error
func (c *TransitionPermissionChecker) AddRule(rule *PermissionRule) error
func (c *TransitionPermissionChecker) GetRequiredRoles(from, to CardState) []string
```

#### 2.4 状态转换通知
**文件**: `backend/application/requirement/transition_notifier.go`

```go
type TransitionNotifier struct {
    notificationService notification.Service
    templateManager    *template.Manager
}

type NotificationConfig struct {
    To         []string
    Cc         []string
    Template   string
    Priority   string
    Channels   []string
}

func NewTransitionNotifier(service notification.Service, templateMgr *template.Manager) *TransitionNotifier
func (n *TransitionNotifier) Notify(card *requirement.Card, from, to CardState, config *NotificationConfig) error
func (n *TransitionNotifier) GetConfig(from, to CardState) (*NotificationConfig, error)
```

#### 2.5 状态转换验证
**文件**: `backend/application/requirement/transition_verification.go`

```go
type TransitionVerifier struct {
    preCheckers  []TransitionChecker
    postCheckers []TransitionChecker
}

type TransitionChecker interface {
    Check(card *requirement.Card, from, to CardState) error
    Name() string
}

func NewTransitionVerifier() *TransitionVerifier
func (v *TransitionVerifier) AddPreChecker(checker TransitionChecker)
func (v *TransitionVerifier) AddPostChecker(checker TransitionChecker)
func (v *TransitionVerifier) VerifyBefore(card *requirement.Card, to CardState) error
func (v *TransitionVerifier) VerifyAfter(card *requirement.Card, from, to CardState) error
```

---

### 3.3 Phase 3: 状态验证器 (1 天)

#### 3.1 转换前验证
**文件**: `backend/application/requirement/pre_transition_validator.go`

```go
type PreTransitionValidator struct {
    checkers []PreTransitionChecker
}

type PreTransitionChecker interface {
    Check(card *requirement.Card, to CardState) error
}

type ContentCompletenessChecker struct{}

func NewContentCompletenessChecker() *ContentCompletenessChecker
func (c *ContentCompletenessChecker) Check(card *requirement.Card, to CardState) error

type DependencySatisfiedChecker struct {
    repo requirement.CardDependencyRepository
}

func NewDependencySatisfiedChecker(repo requirement.CardDependencyRepository) *DependencySatisfiedChecker
func (c *DependencySatisfiedChecker) Check(card *requirement.Card, to CardState) error
```

#### 3.2 转换后检查
**文件**: `backend/application/requirement/post_transition_validator.go`

```go
type PostTransitionValidator struct {
    checkers []PostTransitionChecker
}

type PostTransitionChecker interface {
    Check(card *requirement.Card, from, to CardState) error
}

type StateConsistencyChecker struct{}

func NewStateConsistencyChecker() *StateConsistencyChecker
func (c *StateConsistencyChecker) Check(card *requirement.Card, from, to CardState) error

type PipelineStateChecker struct {
    pipelineRepo pipeline.PipelineRepository
}

func NewPipelineStateChecker(repo pipeline.PipelineRepository) *PipelineStateChecker
func (c *PipelineStateChecker) Check(card *requirement.Card, from, to CardState) error
```

#### 3.3 一致性验证
**文件**: `backend/application/requirement/consistency_validator.go`

```go
type ConsistencyValidator struct {
    cardRepo      requirement.CardRepository
    dependencyRepo requirement.CardDependencyRepository
    pipelineRepo  pipeline.PipelineRepository
}

func NewConsistencyValidator(cardRepo requirement.CardRepository, dependencyRepo requirement.CardDependencyRepository, pipelineRepo pipeline.PipelineRepository) *ConsistencyValidator
func (v *ConsistencyValidator) ValidateCard(card *requirement.Card) error
func (v *ConsistencyValidator) ValidateDependencies(card *requirement.Card) error
func (v *ConsistencyValidator) ValidatePipeline(card *requirement.Card) error
func (v *ConsistencyValidator) ValidateProject(projectID string) error
```

#### 3.4 业务规则验证
**文件**: `backend/application/requirement/business_rule_validator.go`

```go
type BusinessRuleValidator struct {
    rules []BusinessRule
}

type BusinessRule interface {
    Validate(card *requirement.Card, to CardState) error
    Name() string
    Description() string
}

type NoUnblockedDependentsRule struct {
    dependencyRepo requirement.CardDependencyRepository
}

func NewNoUnblockedDependentsRule(repo requirement.CardDependencyRepository) *NoUnblockedDependentsRule
func (r *NoUnblockedDependentsRule) Validate(card *requirement.Card, to CardState) error

type RequiredFieldsRule struct {
    requiredFields map[CardState][]string
}

func NewRequiredFieldsRule() *RequiredFieldsRule
func (r *RequiredFieldsRule) Validate(card *requirement.Card, to CardState) error
```

#### 3.5 异常处理
**文件**: `backend/application/requirement/transition_error_handler.go`

```go
type TransitionErrorHandler struct {
    logger *log.Logger
    notifier *notification.Notifier
}

type TransitionError struct {
    CardID    string
    From      CardState
    To        CardState
    Error     error
    Timestamp time.Time
    Context   map[string]interface{}
}

func NewTransitionErrorHandler(logger *log.Logger, notifier *notification.Notifier) *TransitionErrorHandler
func (h *TransitionErrorHandler) Handle(err error, card *requirement.Card, to CardState) error
func (h *TransitionErrorHandler) Log(error *TransitionError) error
func (h *TransitionErrorHandler) Notify(error *TransitionError) error
func (h *TransitionErrorHandler) Recover(card *requirement.Card) error
```

---

## 4. 时间安排

| 阶段 | 工作内容 | 预计时间 |
|------|----------|----------|
| Phase 1 | 完整状态机 | 1-2 天 |
| Phase 2 | 状态转换规则 | 1 天 |
| Phase 3 | 状态验证器 | 1 天 |

**总计**: 3-5 天

---

## 5. 验收标准

### 5.1 功能验收
- [ ] 所有状态定义完整
- [ ] 所有转换规则实现
- [ ] 状态转换校验正确
- [ ] 状态转换历史完整
- [ ] 状态转换审计正常
- [ ] 转换条件正确
- [ ] 转换副作用正常
- [ ] 权限控制正确
- [ ] 通知机制正常
- [ ] 转换前后验证正常
- [ ] 一致性验证正确
- [ ] 业务规则验证正常
- [ ] 异常处理完善

### 5.2 性能验收
- [ ] 状态转换响应时间 < 100ms
- [ ] 验证检查响应时间 < 200ms
- [ ] 状态机性能稳定

### 5.3 质量验收
- [ ] 单元测试覆盖率 ≥ 80%
- [ ] 集成测试覆盖主要场景
- [ ] 代码审查通过
- [ ] 文档完整

---

## 6. 风险与缓解

### 6.1 技术风险
**风险**: 状态转换冲突
**缓解**:
- 实现乐观锁
- 提供状态转换重试
- 记录详细转换日志

### 6.2 业务风险
**风险**: 状态转换规则复杂
**缓解**:
- 提供规则测试工具
- 实现规则可视化
- 支持规则版本管理

### 6.3 数据风险
**风险**: 状态不一致
**缓解**:
- 使用事务保证一致性
- 实现状态修复工具
- 定期一致性检查
