# Provider 可用性检测实施计划

> 任务编号: 6. Provider 可用性检测 [P1]
> 预估工作量: 3-5 天
> 优先级: P1

---

## 1. 概述

实现 AI Provider 的可用性检测和自动恢复机制，确保 AI Runtime 的高可用性，及时发现和处理 Provider 故障，提供稳定的服务。

---

## 2. 架构设计

### 2.1 健康检查架构
```
┌─────────────────────────────────────────────────────────┐
│              Health Check Engine                        │
├─────────────────────────────────────────────────────────┤
│  ├── Scheduler (Periodic Check)                        │
│  ├── Checker (Health Check)                            │
│  ├── Analyzer (Failure Analysis)                       │
│  └── Notifier (Alert & Notification)                   │
└─────────────────────────────────────────────────────────┘
                           │
┌─────────────────────────────────────────────────────────┐
│              Provider Layer                            │
│  Provider -> Model -> LLM                              │
└─────────────────────────────────────────────────────────┘
                           │
┌─────────────────────────────────────────────────────────┐
│              Recovery Layer                            │
│  ├── Auto Reconnect                                    │
│  ├── Failover                                           │
│  ├── Circuit Breaker                                   │
│  └── Recovery History                                   │
└─────────────────────────────────────────────────────────┘
```

---

## 3. 实施步骤

### 3.1 Phase 1: 健康检查机制 (1-2 天)

#### 1.1 健康检查调度器
**文件**: `backend/application/ai_runtime/health_scheduler.go`

```go
type HealthCheckScheduler struct {
    providers map[string]*ProviderConfig
    checker   *HealthChecker
    interval  time.Duration
    ctx       context.Context
    cancel    context.CancelFunc
}

type ProviderConfig struct {
    ProviderID    string
    CheckInterval time.Duration
    Timeout       time.Duration
    Enabled       bool
}

func NewHealthCheckScheduler(interval time.Duration) *HealthCheckScheduler
func (s *HealthCheckScheduler) AddProvider(config *ProviderConfig) error
func (s *HealthCheckScheduler) RemoveProvider(providerID string) error
func (s *HealthCheckScheduler) Start(ctx context.Context)
func (s *HealthCheckScheduler) Stop()
func (s *HealthCheckScheduler) CheckProvider(providerID string) (*HealthStatus, error)
```

#### 1.2 健康检查器
**文件**: `backend/application/ai_runtime/health_checker.go`

```go
type HealthChecker struct {
    client *http.Client
}

type HealthCheckConfig struct {
    Endpoint    string
    Method      string
    Headers     map[string]string
    Body        string
    Timeout     time.Duration
    Expected    int
}

type HealthStatus struct {
    ProviderID   string
    Status       string  // healthy, unhealthy, unknown
    ResponseTime time.Duration
    StatusCode   int
    Error        error
    Timestamp    time.Time
    Details      map[string]interface{}
}

func NewHealthChecker(timeout time.Duration) *HealthChecker
func (c *HealthChecker) Check(config *HealthCheckConfig) (*HealthStatus, error)
func (c *HealthChecker) Ping(endpoint string, timeout time.Duration) (*HealthStatus, error)
func (c *HealthChecker) Probe(providerID string) (*HealthStatus, error)
```

#### 1.3 响应时间监控
**文件**: `backend/application/ai_runtime/response_monitor.go`

```go
type ResponseMonitor struct {
    metrics map[string]*ResponseMetrics
    window  int
}

type ResponseMetrics struct {
    Samples      []time.Duration
    Avg          time.Duration
    P50          time.Duration
    P95          time.Duration
    P99          time.Duration
    Min          time.Duration
    Max          time.Duration
    LastUpdated  time.Time
}

func NewResponseMonitor(window int) *ResponseMonitor
func (m *ResponseMonitor) Record(providerID string, duration time.Duration) error
func (m *ResponseMonitor) GetMetrics(providerID string) (*ResponseMetrics, error)
func (m *ResponseMonitor) CalculateMetrics(samples []time.Duration) *ResponseMetrics
func (m *ResponseMonitor) IsSlow(providerID string, threshold time.Duration) bool
```

#### 1.4 错误率统计
**文件**: `backend/application/ai_runtime/error_tracker.go`

```go
type ErrorTracker struct {
    errors map[string]*ErrorStats
    window int
}

type ErrorStats struct {
    Total      int
    Errors     int
    ErrorRate  float64
    ErrorTypes map[string]int
    LastError  time.Time
}

func NewErrorTracker(window int) *ErrorTracker
func (m *ErrorTracker) Record(providerID string, err error) error
func (m *ErrorTracker) RecordSuccess(providerID string) error
func (m *ErrorTracker) GetStats(providerID string) (*ErrorStats, error)
func (m *ErrorTracker) IsHighErrorRate(providerID string, threshold float64) bool
```

#### 1.5 健康状态评估
**文件**: `backend/application/ai_runtime/health_evaluator.go`

```go
type HealthEvaluator struct {
    monitor       *ResponseMonitor
    errorTracker  *ErrorTracker
    config        *EvaluationConfig
}

type EvaluationConfig struct {
    MaxResponseTime    time.Duration
    MaxErrorRate       float64
    MinSuccessRate     float64
    ConsecutiveFailures int
}

func NewHealthEvaluator(config *EvaluationConfig) *HealthEvaluator
func (e *HealthEvaluator) Evaluate(providerID string) (*ProviderHealth, error)
func (e *HealthEvaluator) CheckResponseTime(providerID string) bool
func (e *HealthEvaluator) CheckErrorRate(providerID string) bool
```

---

### 3.2 Phase 2: 故障检测 (1 天)

#### 2.1 连接失败检测
**文件**: `backend/application/ai_runtime/connection_detector.go`

```go
type ConnectionDetector struct {
    timeout time.Duration
    retries int
}

type ConnectionResult struct {
    Success      bool
    Latency      time.Duration
    Error        error
    Timestamp    time.Time
}

func NewConnectionDetector(timeout time.Duration, retries int) *ConnectionDetector
func (d *ConnectionDetector) Detect(endpoint string) (*ConnectionResult, error)
func (d *ConnectionDetector) TestConnection(endpoint string) (*ConnectionResult, error)
```

#### 2.2 超时检测
**文件**: `backend/application/ai_runtime/timeout_detector.go`

```go
type TimeoutDetector struct {
    threshold time.Duration
    window    int
}

func NewTimeoutDetector(threshold time.Duration, window int) *TimeoutDetector
func (d *TimeoutDetector) Detect(providerID string, metrics *ResponseMetrics) bool
func (d *TimeoutDetector) GetTimeoutRate(providerID string) float64
```

#### 2.3 异常响应检测
**文件**: `backend/application/ai_runtime/anomaly_detector.go`

```go
type AnomalyDetector struct {
    baseline  *ResponseMetrics
    threshold float64
}

type Anomaly struct {
    Type       string
    Severity   string
    Value      interface{}
    Expected   interface{}
    Timestamp  time.Time
}

func NewAnomalyDetector(baseline *ResponseMetrics, threshold float64) *AnomalyDetector
func (d *AnomalyDetector) Detect(metrics *ResponseMetrics) []*Anomaly
func (d *AnomalyDetector) DetectLatencyAnomaly(latency time.Duration) bool
func (d *AnomalyDetector) DetectErrorRateAnomaly(errorRate float64) bool
```

#### 2.4 限流检测
**文件**: `backend/application/ai_runtime/rate_limit_detector.go`

```go
type RateLimitDetector struct {
    statusCode int
    patterns   []string
}

func NewRateLimitDetector() *RateLimitDetector
func (d *RateLimitDetector) Detect(statusCode int, body string) bool
func (d *RateLimitDetector) ParseRetryAfter(headers http.Header) time.Duration
```

#### 2.5 配置错误检测
**文件**: `backend/application/ai_runtime/config_detector.go`

```go
type ConfigDetector struct {
    validator *validator.Validator
}

func NewConfigDetector() *ConfigDetector
func (d *ConfigDetector) Detect(config *ProviderConfig) []*ConfigError
func (d *ConfigDetector) ValidateEndpoint(endpoint string) error
func (d *ConfigDetector) ValidateAPIKey(apiKey string) error
```

---

### 3.3 Phase 3: 自动恢复 (1 天)

#### 3.1 自动重连机制
**文件**: `backend/application/ai_runtime/auto_reconnect.go`

```go
type AutoReconnector struct {
    maxRetries      int
    retryDelay      time.Duration
    maxRetryDelay   time.Duration
    backoffMultiplier float64
}

type ReconnectResult struct {
    Success    bool
    Attempts   int
    Duration   time.Duration
    Error      error
}

func NewAutoReconnector(maxRetries int, retryDelay time.Duration) *AutoReconnector
func (r *AutoReconnector) Reconnect(providerID string) (*ReconnectResult, error)
func (r *AutoReconnector) ShouldRetry(attempt int, err error) bool
func (r *AutoReconnector) GetRetryDelay(attempt int) time.Duration
```

#### 3.2 备用模型切换
**文件**: `backend/application/ai_runtime/failover.go`

```go
type FailoverManager struct {
    modelMgr        *model.Manager
    fallbackMap     map[string][]string  // model_id -> fallback_model_ids
    currentFallback map[string]string    // model_id -> current_fallback
}

type FailoverResult struct {
    OriginalModel string
    FallbackModel string
    Success       bool
    Reason        string
}

func NewFailoverManager(modelMgr *model.Manager) *FailoverManager
func (m *FailoverManager) AddFallback(modelID, fallbackID string) error
func (m *FailoverManager) ExecuteFailover(modelID string) (*FailoverResult, error)
func (m *FailoverManager) SelectFallback(modelID string) (string, error)
func (m *FailoverManager) Restore(modelID string) error
```

#### 3.3 故障隔离
**文件**: `backend/application/ai_runtime/circuit_breaker.go`

```go
type CircuitBreaker struct {
    state           CircuitState
    failureCount    int
    successCount    int
    lastFailureTime time.Time
    config          *CircuitBreakerConfig
}

type CircuitState string
const (
    StateClosed   CircuitState = "closed"
    StateOpen     CircuitState = "open"
    StateHalfOpen CircuitState = "half-open"
)

type CircuitBreakerConfig struct {
    FailureThreshold   int
    SuccessThreshold   int
    Timeout            time.Duration
    HalfOpenMaxCalls   int
}

func NewCircuitBreaker(config *CircuitBreakerConfig) *CircuitBreaker
func (cb *CircuitBreaker) Execute(fn func() error) error
func (cb *CircuitBreaker) RecordSuccess()
func (cb *CircuitBreaker) RecordFailure()
func (cb *CircuitBreaker) GetState() CircuitState
func (cb *CircuitBreaker) Reset()
```

#### 3.4 恢复通知
**文件**: `backend/application/ai_runtime/recovery_notifier.go`

```go
type RecoveryNotifier struct {
    alertMgr    *alert.Manager
    templateMgr *template.Manager
}

type RecoveryEvent struct {
    ProviderID  string
    EventType   string  // failure, recovery, degraded
    Timestamp   time.Time
    Details     map[string]interface{}
}

func NewRecoveryNotifier(alertMgr *alert.Manager) *RecoveryNotifier
func (n *RecoveryNotifier) NotifyFailure(providerID string, err error) error
func (n *RecoveryNotifier) NotifyRecovery(providerID string) error
func (n *RecoveryNotifier) NotifyDegraded(providerID string, metrics *ResponseMetrics) error
```

#### 3.5 恢复历史记录
**文件**: `backend/application/ai_runtime/recovery_history.go`

```go
type RecoveryHistory struct {
    events map[string][]*RecoveryEvent
    maxAge time.Duration
}

func NewRecoveryHistory(maxAge time.Duration) *RecoveryHistory
func (h *RecoveryHistory) Record(event *RecoveryEvent) error
func (h *RecoveryHistory) GetHistory(providerID string, duration time.Duration) []*RecoveryEvent
func (h *RecoveryHistory) GetFailureCount(providerID string, duration time.Duration) int
func (h *RecoveryHistory) GetRecoveryTime(providerID string) time.Duration
func (h *RecoveryHistory) CleanOldEvents() error
```

---

## 4. 时间安排

| 阶段 | 工作内容 | 预计时间 |
|------|----------|----------|
| Phase 1 | 健康检查机制 | 1-2 天 |
| Phase 2 | 故障检测 | 1 天 |
| Phase 3 | 自动恢复 | 1 天 |

**总计**: 3-5 天

---

## 5. 验收标准

### 5.1 功能验收
- [ ] 定期健康检查正常
- [ ] 响应时间监控准确
- [ ] 错误率统计正确
- [ ] 健康状态评估合理
- [ ] 故障检测及时
- [ ] 自动重连机制正常
- [ ] 备用模型切换有效
- [ ] 故障隔离正常
- [ ] 恢复通知及时
- [ ] 恢复历史记录完整

### 5.2 性能验收
- [ ] 健康检查响应时间 < 5 秒
- [ ] 故障检测延迟 < 30 秒
- [ ] 自动恢复时间 < 1 分钟
- [ ] 支持并发检查

### 5.3 质量验收
- [ ] 单元测试覆盖率 ≥ 80%
- [ ] 集成测试覆盖主要场景
- [ ] 代码审查通过
- [ ] 文档完整

---

## 6. 风险与缓解

### 6.1 技术风险
**风险**: 健康检查误报
**缓解**:
- 实现多次检查确认
- 使用滑动窗口统计
- 调整检测阈值

### 6.2 性能风险
**风险**: 健康检查开销大
**缓解**:
- 使用异步检查
- 合理设置检查间隔
- 优化检查逻辑

### 6.3 恢复风险
**风险**: 自动恢复失败
**缓解**:
- 实现多重恢复策略
- 提供人工干预机制
- 记录详细恢复日志
