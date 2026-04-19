# Token 计费实施计划

> 任务编号: 11. Token 计费 [P2]
> 预估工作量: 3-5 天
> 优先级: P2

---

## 1. 概述

实现 AI Runtime 的 Token 计费系统，支持按量计费、分层计费、包月计费等多种计费策略，提供准确的费用计算和账单生成功能。

---

## 2. 架构设计

### 2.1 计费架构
```
┌─────────────────────────────────────────────────────────┐
│              Cost Calculation Engine                     │
├─────────────────────────────────────────────────────────┤
│  ├── Token Counter (Token Counting)                    │
│  ├── Cost Calculator (Cost Calculation)                 │
│  ├── Pricing Engine (Pricing Rules)                    │
│  └── Budget Controller (Budget Control)                │
└─────────────────────────────────────────────────────────┘
                           │
┌─────────────────────────────────────────────────────────┐
│              Billing Engine                             │
├─────────────────────────────────────────────────────────┤
│  ├── Bill Generator (Bill Generation)                   │
│  ├── Bill Aggregator (Bill Aggregation)                 │
│  ├── Bill Formatter (Bill Formatting)                   │
│  └── Bill Exporter (Bill Export)                       │
└─────────────────────────────────────────────────────────┘
```

---

## 3. 实施步骤

### 3.1 Phase 1: 费用计算引擎 (1-2 天)

#### 1.1 费用模型配置
**文件**: `backend/domain/billing/pricing_model.go`

```go
type PricingModel struct {
    ID              string
    Name            string
    Provider        string
    Model           string
    InputPrice      float64  // per 1K tokens
    OutputPrice     float64  // per 1K tokens
    Currency        string
    EffectiveFrom   time.Time
    EffectiveTo     *time.Time
}

type ModelPricing struct {
    models map[string]*PricingModel  // model_id -> pricing
}

func NewModelPricing() *ModelPricing
func (p *ModelPricing) AddPricing(model *PricingModel) error
func (p *ModelPricing) GetPricing(modelID string) (*PricingModel, error)
func (p *ModelPricing) GetPrice(modelID string, inputTokens, outputTokens int) float64
```

#### 1.2 模型价格表
**文件**: `backend/domain/billing/price_table.go`

```go
type PriceTable struct {
    id       string
    currency string
    prices   map[string]*ModelPrice
}

type ModelPrice struct {
    ModelID    string
    InputPrice float64
    OutputPrice float64
    Unit       string  // tokens, characters, requests
}

func NewPriceTable(currency string) *PriceTable
func (t *PriceTable) SetPrice(modelID string, inputPrice, outputPrice float64) error
func (t *PriceTable) GetPrice(modelID string) (*ModelPrice, error)
func (t *PriceTable) CalculateCost(modelID string, inputTokens, outputTokens int) float64
func (t *PriceTable) GetCurrency() string
```

#### 1.3 用量计算
**文件**: `backend/application/billing/usage_calculator.go`

```go
type UsageCalculator struct {
    tokenizer   *tokenizer.Tokenizer
    priceTable  *PriceTable
}

type Usage struct {
    ModelID       string
    InputTokens   int
    OutputTokens  int
    TotalTokens   int
    InputCost     float64
    OutputCost    float64
    TotalCost     float64
    Timestamp     time.Time
}

func NewUsageCalculator(tokenizer *tokenizer.Tokenizer, priceTable *PriceTable) *UsageCalculator
func (c *UsageCalculator) Calculate(modelID string, input, output string) (*Usage, error)
func (c *UsageCalculator) CalculateFromTokens(modelID string, inputTokens, outputTokens int) (*Usage, error)
func (c *UsageCalculator) CountTokens(text string) (int, error)
```

#### 1.4 费用合计
**文件**: `backend/application/billing/cost_aggregator.go`

```go
type CostAggregator struct {
    usageRepo billing.UsageRepository
}

type CostSummary struct {
    UserID       string
    ModelID      string
    StartTime    time.Time
    EndTime      time.Time
    TotalCost    float64
    TotalTokens  int
    RequestCount int
    Breakdown    []*CostBreakdown
}

type CostBreakdown struct {
    Date         string
    ModelID      string
    Cost         float64
    Tokens       int
    Requests     int
}

func NewCostAggregator(usageRepo billing.UsageRepository) *CostAggregator
func (a *CostAggregator) Aggregate(userID string, start, end time.Time) (*CostSummary, error)
func (a *CostAggregator) AggregateByModel(userID string, start, end time.Time) []*CostBreakdown
func (a *CostAggregator) AggregateByDate(userID string, start, end time.Time) []*CostBreakdown
```

#### 1.5 费用预估
**文件**: `backend/application/billing/cost_estimator.go`

```go
type CostEstimator struct {
    priceTable *PriceTable
    tokenizer  *tokenizer.Tokenizer
}

type CostEstimate struct {
    ModelID      string
    InputTokens  int
    OutputTokens int
    MinCost      float64
    MaxCost      float64
    AvgCost      float64
    Confidence   float64
}

func NewCostEstimator(priceTable *PriceTable, tokenizer *tokenizer.Tokenizer) *CostEstimator
func (e *CostEstimator) Estimate(modelID string, input string) (*CostEstimate, error)
func (e *CostEstimator) EstimateFromTokens(modelID string, inputTokens, estimatedOutputTokens int) (*CostEstimate, error)
func (e *CostEstimator) EstimateBatch(estimates []*CostEstimate) (*CostSummary, error)
```

---

### 3.2 Phase 2: 计费策略 (1 天)

#### 2.1 按量计费
**文件**: `backend/application/billing/usage_based_billing.go`

```go
type UsageBasedBilling struct {
    priceTable *PriceTable
}

type UsageBasedConfig struct {
    Unit       string  // tokens, requests, characters
    Price      float64
    Currency   string
}

func NewUsageBasedBilling(priceTable *PriceTable) *UsageBasedBilling
func (b *UsageBasedBilling) Calculate(modelID string, inputTokens, outputTokens int) float64
func (b *UsageBasedBilling) CalculateRequest(modelID string, request *ModelRequest) float64
```

#### 2.2 分层计费
**文件**: `backend/application/billing/tiered_billing.go`

```go
type TieredBilling struct {
    tiers []*BillingTier
}

type BillingTier struct {
    Level      int
    MinUsage   int
    MaxUsage   *int
    Price      float64
    Currency   string
}

type TieredConfig struct {
    TierType   string  // tokens, requests, cost
    Tiers      []*BillingTier
}

func NewTieredBilling(tiers []*BillingTier) *TieredBilling
func (b *TieredBilling) Calculate(usage int) float64
func (b *TieredBilling) GetTier(usage int) *BillingTier
func (b *TieredBilling) AddTier(tier *BillingTier) error
```

#### 2.3 包月计费
**文件**: `backend/application/billing/monthly_billing.go`

```go
type MonthlyBilling struct {
    plans []*BillingPlan
}

type BillingPlan struct {
    ID          string
    Name        string
    Price       float64
    Currency    string
    Duration    time.Duration
    Features    []string
    Limitations map[string]interface{}
}

type Subscription struct {
    UserID      string
    PlanID      string
    StartTime   time.Time
    EndTime     time.Time
    Status      string
}

func NewMonthlyBilling(plans []*BillingPlan) *MonthlyBilling
func (b *MonthlyBilling) CreateSubscription(userID, planID string) (*Subscription, error)
func (b *MonthlyBilling) GetBill(subscription *Subscription) float64
func (b *MonthlyBilling) CheckLimit(subscription *Subscription, limitType string, value interface{}) bool
```

#### 2.4 免费额度
**文件**: `backend/application/billing/free_quota.go`

```go
type FreeQuota struct {
    quotas map[string]*QuotaConfig
}

type QuotaConfig struct {
    ModelID      string
    FreeTokens   int
    FreeRequests int
    Period       time.Duration
    ResetTime    time.Time
}

type QuotaUsage struct {
    ModelID      string
    UsedTokens   int
    UsedRequests int
    ResetTime    time.Time
}

func NewFreeQuota() *FreeQuota
func (q *FreeQuota) SetQuota(modelID string, config *QuotaConfig) error
func (q *FreeQuota) GetQuota(modelID string) (*QuotaConfig, error)
func (q *FreeQuota) CheckQuota(modelID string, usage *QuotaUsage) (bool, error)
func (q *FreeQuota) ResetQuota(modelID string) error
```

#### 2.5 优惠规则
**文件**: `backend/application/billing/discount_rule.go`

```go
type DiscountRule struct {
    ID          string
    Name        string
    Type        string  // percentage, fixed, tiered
    Value       float64
    Conditions  []*DiscountCondition
    ValidFrom   time.Time
    ValidTo     *time.Time
}

type DiscountCondition struct {
    Type      string  // user, model, usage, time
    Operator  string
    Value     interface{}
}

type DiscountApplier struct {
    rules []*DiscountRule
}

func NewDiscountApplier() *DiscountApplier
func (a *DiscountApplier) AddRule(rule *DiscountRule) error
func (a *DiscountApplier) Apply(cost float64, context map[string]interface{}) float64
func (a *DiscountApplier) GetApplicableRules(context map[string]interface{}) []*DiscountRule
```

---

### 3.3 Phase 3: 账单生成 (1 天)

#### 3.1 账单周期
**文件**: `backend/domain/billing/bill_cycle.go`

```go
type BillCycle struct {
    ID         string
    Name       string
    Type       string  // daily, weekly, monthly, yearly
    StartDate  time.Time
    EndDate    time.Time
    Status     string
}

func NewBillCycle(cycleType string) (*BillCycle, error)
func (c *BillCycle) GetCurrentPeriod() (time.Time, time.Time)
func (c *BillCycle) GetNextPeriod() (time.Time, time.Time)
func (c *BillCycle) Contains(date time.Time) bool
```

#### 3.2 账单明细
**文件**: `backend/domain/billing/bill_item.go`

```go
type BillItem struct {
    ID           string
    BillID       string
    ModelID      string
    Description  string
    Quantity     int
    UnitPrice    float64
    TotalPrice   float64
    Currency     string
    Timestamp    time.Time
    Metadata     map[string]interface{}
}

type BillItemList struct {
    items []*BillItem
}

func NewBillItem(billID, modelID, description string, quantity int, unitPrice float64) *BillItem
func (l *BillItemList) AddItem(item *BillItem)
func (l *BillItemList) GetTotal() float64
func (l *BillItemList) GetItemsByModel(modelID string) []*BillItem
```

#### 3.3 账单统计
**文件**: `backend/application/billing/bill_statistics.go`

```go
type BillStatistics struct {
    BillID       string
    TotalCost    float64
    TotalTokens  int
    TotalRequests int
    ModelBreakdown  map[string]*ModelStatistics
    TimeBreakdown   []*TimeStatistics
}

type ModelStatistics struct {
    ModelID      string
    Cost         float64
    Tokens       int
    Requests     int
}

type TimeStatistics struct {
    Period       string
    Cost         float64
    Tokens       int
    Requests     int
}

func NewBillStatistics() *BillStatistics
func (s *BillStatistics) AddUsage(usage *Usage)
func (s *BillStatistics) Calculate()
```

#### 3.4 账单导出
**文件**: `backend/application/billing/bill_exporter.go`

```go
type BillExporter struct {
    formatter BillFormatter
}

type BillFormat string
const (
    FormatJSON  BillFormat = "json"
    FormatCSV   BillFormat = "csv"
    FormatPDF   BillFormat = "pdf"
    FormatHTML  BillFormat = "html"
)

type BillFormatter interface {
    Format(bill *Bill) ([]byte, error)
}

func NewBillExporter(formatter BillFormatter) *BillExporter
func (e *BillExporter) Export(bill *Bill, format BillFormat) ([]byte, error)
```

#### 3.5 账单通知
**文件**: `backend/application/billing/bill_notifier.go`

```go
type BillNotifier struct {
    notificationService notification.Service
    templateManager     *template.Manager
}

type BillNotification struct {
    BillID      string
    UserID      string
    TotalCost   float64
    DueDate     time.Time
    Attachments []string
}

func NewBillNotifier(service notification.Service, templateMgr *template.Manager) *BillNotifier
func (n *BillNotifier) Notify(bill *Bill) error
func (n *BillNotifier) NotifyOverdue(bill *Bill) error
func (n *BillNotifier) NotifyPayment(bill *Bill) error
```

---

## 4. 时间安排

| 阶段 | 工作内容 | 预计时间 |
|------|----------|----------|
| Phase 1 | 费用计算引擎 | 1-2 天 |
| Phase 2 | 计费策略 | 1 天 |
| Phase 3 | 账单生成 | 1 天 |

**总计**: 3-5 天

---

## 5. 验收标准

### 5.1 功能验收
- [ ] 费用模型配置正常
- [ ] 用量计算准确
- [ ] 费用合计正确
- [ ] 费用预估合理
- [ ] 按量计费正常
- [ ] 分层计费正常
- [ ] 包月计费正常
- [ ] 免费额度正确
- [ ] 优惠规则生效
- [ ] 账单生成正常
- [ ] 账单导出正常
- [ ] 账单通知及时

### 5.2 性能验收
- [ ] Token 计算速度 ≥ 10000 tokens/s
- [ ] 费用计算响应时间 < 100ms
- [ ] 账单生成时间 < 5 秒

### 5.3 质量验收
- [ ] 单元测试覆盖率 ≥ 80%
- [ ] 集成测试覆盖主要场景
- [ ] 代码审查通过
- [ ] 文档完整

---

## 6. 风险与缓解

### 6.1 技术风险
**风险**: Token 计算不准确
**缓解**:
- 使用成熟的 Token 计算库
- 提供人工校对机制
- 实现差异检测

### 6.2 计费风险
**风险**: 计费规则复杂
**缓解**:
- 提供计费规则测试工具
- 实现计费规则可视化
- 支持计费规则版本管理

### 6.3 数据风险
**风险**: 计费数据丢失
**缓解**:
- 实现计费数据备份
- 使用事务保证一致性
- 提供数据恢复机制
