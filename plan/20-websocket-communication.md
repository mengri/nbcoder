# WebSocket 实时通信实施计划

> 任务编号: 20. WebSocket 实时通信 [P2]
> 预估工作量: 3-5 天
> 优先级: P2

---

## 1. 概述

实现 WebSocket 实时通信系统，支持服务器推送、消息广播、房间管理和心跳检测等功能，为前端提供实时数据更新能力。

---

## 2. 架构设计

### 2.1 WebSocket 架构
```
┌─────────────────────────────────────────────────────────┐
│              WebSocket Server                           │
├─────────────────────────────────────────────────────────┤
│  ├── Connection Manager (Connection Management)        │
│  ├── Message Router (Message Routing)                  │
│  ├── Room Manager (Room Management)                    │
│  ├── Heartbeat (Heartbeat Detection)                   │
│  └── Rate Limiter (Rate Limiting)                      │
└─────────────────────────────────────────────────────────┘
                           │
┌─────────────────────────────────────────────────────────┐
│              Message Queue                             │
│  └── Persistent Queue (Message Persistence)             │
└─────────────────────────────────────────────────────────┘
                           │
┌─────────────────────────────────────────────────────────┐
│              Event Publisher                           │
│  └── Event Dispatcher (Event Distribution)              │
└─────────────────────────────────────────────────────────┘
```

---

## 3. 实施步骤

### 3.1 Phase 1: WebSocket 服务器 (1 天)

#### 1.1 连接管理
**文件**: `backend/interfaces/websocket/connection_manager.go`

```go
type ConnectionManager struct {
    connections map[string]*Connection
    mu          sync.RWMutex
}

type Connection struct {
    ID         string
    Conn       *websocket.Conn
    UserID     string
    Rooms      []string
    CreatedAt  time.Time
    LastActive time.Time
    Send       chan []byte
}

func NewConnectionManager() *ConnectionManager
func (m *ConnectionManager) Add(conn *websocket.Conn, userID string) (*Connection, error)
func (m *ConnectionManager) Remove(connectionID string) error
func (m *ConnectionManager) Get(connectionID string) (*Connection, bool)
func (m *ConnectionManager) GetByUser(userID string) []*Connection
func (m *ConnectionManager) Count() int
func (m *ConnectionManager) CloseAll() error
```

#### 1.2 消息广播
**文件**: `backend/interfaces/websocket/broadcaster.go`

```go
type Broadcaster struct {
    connectionManager *ConnectionManager
}

type Message struct {
    Type      string                 `json:"type"`
    Data      interface{}            `json:"data"`
    Timestamp time.Time              `json:"timestamp"`
    Meta      map[string]interface{} `json:"meta"`
}

func NewBroadcaster(cm *ConnectionManager) *Broadcaster
func (b *Broadcaster) Broadcast(message *Message) error
func (b *Broadcaster) BroadcastToUser(userID string, message *Message) error
func (b *Broadcaster) BroadcastToRoom(roomID string, message *Message) error
func (b *Broadcaster) BroadcastToUsers(userIDs []string, message *Message) error
func (b *Broadcaster) BroadcastExcept(excludeID string, message *Message) error
```

#### 1.3 房间管理
**文件**: `backend/interfaces/websocket/room_manager.go`

```go
type RoomManager struct {
    rooms map[string]*Room
    mu    sync.RWMutex
}

type Room struct {
    ID          string
    Name        string
    Connections map[string]bool  // connection_id -> true
    MaxSize     int
    CreatedAt   time.Time
    Metadata    map[string]interface{}
}

func NewRoomManager() *RoomManager
func (m *RoomManager) CreateRoom(id, name string, maxSize int) (*Room, error)
func (m *RoomManager) GetRoom(id string) (*Room, bool)
func (m *RoomManager) JoinRoom(roomID, connectionID string) error
func (m *RoomManager) LeaveRoom(roomID, connectionID string) error
func (m *RoomManager) DeleteRoom(id string) error
func (m *RoomManager) GetRoomConnections(roomID string) []*Connection
func (m *RoomManager) GetUserRooms(connectionID string) []*Room
```

#### 1.4 心跳检测
**文件**: `backend/interfaces/websocket/heartbeat.go`

```go
type Heartbeat struct {
    interval       time.Duration
    timeout        time.Duration
    connectionMgr  *ConnectionManager
    ticker         *time.Ticker
}

type HeartbeatMessage struct {
    Type      string    `json:"type"`
    Timestamp time.Time `json:"timestamp"`
}

func NewHeartbeat(interval, timeout time.Duration, cm *ConnectionManager) *Heartbeat
func (h *Heartbeat) Start(ctx context.Context)
func (h *Heartbeat) Stop()
func (h *Heartbeat) CheckConnections() error
func (h *Heartbeat) IsAlive(conn *Connection) bool
func (h *Heartbeat) SendHeartbeat(conn *Connection) error
```

#### 1.5 连接限制
**文件**: `backend/interfaces/websocket/limiter.go`

```go
type ConnectionLimiter struct {
    maxConnections  int
    maxConnectionsPerUser int
    connectionMgr   *ConnectionManager
    userConnections map[string]int
    mu              sync.RWMutex
}

func NewConnectionLimiter(maxConns, maxConnsPerUser int, cm *ConnectionManager) *ConnectionLimiter
func (l *ConnectionLimiter) CanConnect(userID string) (bool, error)
func (l *ConnectionLimiter) RecordConnection(userID string) error
func (l *ConnectionLimiter) RemoveConnection(userID string) error
func (l *ConnectionLimiter) GetUserConnectionCount(userID string) int
func (l *ConnectionLimiter) GetTotalConnectionCount() int
```

---

### 3.2 Phase 2: 实时推送机制 (1 天)

#### 2.1 消息队列
**文件**: `backend/interfaces/websocket/message_queue.go`

```go
type MessageQueue struct {
    queue     chan *QueuedMessage
    batchSize int
    flushTime time.Duration
    storage   QueueStorage
}

type QueuedMessage struct {
    ID        string
    Target    string  // user_id, room_id, or "all"
    Message   *Message
    Priority  int
    CreatedAt time.Time
}

type QueueStorage interface {
    Save(message *QueuedMessage) error
    Get(limit int) ([]*QueuedMessage, error)
    Delete(ids []string) error
}

func NewMessageQueue(batchSize int, flushTime time.Duration, storage QueueStorage) *MessageQueue
func (q *MessageQueue) Enqueue(target string, message *Message, priority int) error
func (q *MessageQueue) Start(ctx context.Context)
func (q *MessageQueue) Stop()
func (q *MessageQueue) Process(broadcaster *Broadcaster) error
```

#### 2.2 消息路由
**文件**: `backend/interfaces/websocket/message_router.go`

```go
type MessageRouter struct {
    handlers map[string]MessageHandler
}

type MessageHandler func(ctx context.Context, conn *Connection, msg *Message) error

func NewMessageRouter() *MessageRouter
func (r *MessageRouter) RegisterHandler(messageType string, handler MessageHandler) error
func (r *MessageRouter) UnregisterHandler(messageType string) error
func (r *MessageRouter) Route(ctx context.Context, conn *Connection, msg *Message) error
func (r *MessageRouter) HasHandler(messageType string) bool
```

#### 2.3 消息持久化
**文件**: `backend/interfaces/websocket/message_persistence.go`

```go
type MessagePersistence struct {
    db    *gorm.DB
    table string
}

type PersistentMessage struct {
    ID        string    `gorm:"primaryKey"`
    Type      string
    Data      []byte
    Target    string
    CreatedAt time.Time
    SentAt    *time.Time
}

func NewMessagePersistence(db *gorm.DB, table string) *MessagePersistence
func (p *MessagePersistence) Save(message *QueuedMessage) error
func (p *MessagePersistence) MarkAsSent(id string) error
func (p *MessagePersistence) GetUnsent(limit int) ([]*QueuedMessage, error)
func (p *MessagePersistence) CleanOldMessages(olderThan time.Duration) error
```

#### 2.4 消息确认
**文件**: `backend/interfaces/websocket/message_ack.go`

```go
type MessageAck struct {
    MessageID  string
    ConnectionID string
    AckedAt    time.Time
    Status     string
}

type AckManager struct {
    pending map[string]*PendingMessage
    timeout time.Duration
}

type PendingMessage struct {
    Message    *Message
    SentAt     time.Time
    Retries    int
}

func NewAckManager(timeout time.Duration) *AckManager
func (m *AckManager) SendWithAck(conn *Connection, message *Message) error
func (m *AckManager) ReceiveAck(ack *MessageAck) error
func (m *AckManager) CheckTimeouts() []*PendingMessage
func (m *AckManager) Retry(pending *PendingMessage) error
```

#### 2.5 消息重发
**文件**: `backend/interfaces/websocket/message_retry.go`

```go
type RetryManager struct {
    pending  map[string]*RetryMessage
    maxRetries int
    backoff   time.Duration
}

type RetryMessage struct {
    Message    *Message
    Target     string
    SentAt     time.Time
    Retries    int
    NextRetry  time.Time
}

func NewRetryManager(maxRetries int, backoff time.Duration) *RetryManager
func (m *RetryManager) ScheduleRetry(target string, message *Message) error
func (m *RetryManager) ProcessRetries(broadcaster *Broadcaster) error
func (m *RetryManager) Remove(id string) error
func (m *RetryManager) GetRetryCount(id string) int
```

---

### 3.3 Phase 3: 连接管理 (1 天)

#### 3.1 连接认证
**文件**: `backend/interfaces/websocket/auth.go`

```go
type WebSocketAuth struct {
    jwtVerifier *jwt.Verifier
}

type AuthMessage struct {
    Type     string `json:"type"`
    Token    string `json:"token"`
}

type AuthResult struct {
    Success bool
    UserID  string
    Error   error
}

func NewWebSocketAuth(jwtVerifier *jwt.Verifier) *WebSocketAuth
func (a *WebSocketAuth) Authenticate(token string) (*AuthResult, error)
func (a *WebSocketAuth) ValidateConnection(conn *websocket.Conn) (*AuthResult, error)
func (a *WebSocketAuth) HandleAuthMessage(conn *Connection, msg *AuthMessage) (*AuthResult, error)
```

#### 3.2 连接授权
**文件**: `backend/interfaces/websocket/authorization.go`

```go
type WebSocketAuthorization struct {
    roleMgr *rbac.Manager
}

type Permission struct {
    Resource string
    Action   string
}

func NewWebSocketAuthorization(roleMgr *rbac.Manager) *WebSocketAuthorization
func (a *WebSocketAuthorization) CheckPermission(userID string, permission Permission) (bool, error)
func (a *WebSocketAuthorization) CanJoinRoom(userID, roomID string) (bool, error)
func (a *WebSocketAuthorization) CanSubscribe(userID, channel string) (bool, error)
func (a *WebSocketAuthorization) CanPublish(userID, channel string) (bool, error)
```

#### 3.3 连接监控
**文件**: `backend/interfaces/websocket/connection_monitor.go`

```go
type ConnectionMonitor struct {
    connectionMgr *ConnectionManager
    metrics       *ConnectionMetrics
    collectors    []MetricCollector
}

type ConnectionMetrics struct {
    TotalConnections     int
    ActiveConnections    int
    MessagesSent         int64
    MessagesReceived     int64
    BytesSent            int64
    BytesReceived        int64
    AvgMessageLatency    time.Duration
    ErrorRate            float64
}

func NewConnectionMonitor(cm *ConnectionManager) *ConnectionMonitor
func (m *ConnectionMonitor) Start(ctx context.Context)
func (m *ConnectionMonitor) GetMetrics() *ConnectionMetrics
func (m *ConnectionMonitor) GetMetricsHistory(duration time.Duration) []*ConnectionMetrics
func (m *ConnectionMonitor) RecordMessageSent(size int)
func (m *ConnectionMonitor) RecordMessageReceived(size int)
func (m *ConnectionMonitor) RecordError()
```

#### 3.4 连接限流
**文件**: `backend/interfaces/websocket/rate_limiter.go`

```go
type ConnectionRateLimiter struct {
    limits map[string]*RateLimit
    mu     sync.RWMutex
}

type RateLimit struct {
    MaxRequests int
    Window      time.Duration
    Requests    []time.Time
}

func NewConnectionRateLimiter() *ConnectionRateLimiter
func (l *ConnectionRateLimiter) SetLimit(connectionID string, maxRequests int, window time.Duration) error
func (l *ConnectionRateLimiter) CheckLimit(connectionID string) (bool, error)
func (l *ConnectionRateLimiter) RecordRequest(connectionID string) error
func (l *ConnectionRateLimiter) GetRemaining(connectionID string) (int, error)
func (l *ConnectionRateLimiter) Reset(connectionID string) error
```

#### 3.5 连接清理
**文件**: `backend/interfaces/websocket/connection_cleaner.go`

```go
type ConnectionCleaner struct {
    connectionMgr *ConnectionManager
    roomMgr       *RoomManager
    checkInterval time.Duration
    idleTimeout   time.Duration
}

func NewConnectionCleaner(cm *ConnectionManager, rm *RoomManager, checkInterval, idleTimeout time.Duration) *ConnectionCleaner
func (c *ConnectionCleaner) Start(ctx context.Context)
func (c *ConnectionCleaner) CleanIdleConnections() error
func (c *ConnectionCleaner) CleanDeadConnections() error
func (c *ConnectionCleaner) CleanEmptyRooms() error
func (c *ConnectionCleaner) IsIdle(conn *Connection) bool
func (c *ConnectionCleaner) IsDead(conn *Connection) bool
```

---

### 3.4 Phase 4: 消息队列实现 (可选，1 天)

#### 4.1 队列实现
**文件**: `backend/interfaces/websocket/queue_impl.go`

```go
type InMemoryQueue struct {
    messages []*QueuedMessage
    mu       sync.RWMutex
}

func NewInMemoryQueue() *InMemoryQueue
func (q *InMemoryQueue) Save(message *QueuedMessage) error
func (q *InMemoryQueue) Get(limit int) ([]*QueuedMessage, error)
func (q *InMemoryQueue) Delete(ids []string) error
func (q *InMemoryQueue) Size() int
```

#### 4.2 队列监控
**文件**: `backend/interfaces/websocket/queue_monitor.go`

```go
type QueueMonitor struct {
    queue    MessageQueue
    metrics  *QueueMetrics
}

type QueueMetrics struct {
    TotalMessages    int64
    PendingMessages  int64
    ProcessedMessages int64
    FailedMessages   int64
    AvgProcessTime   time.Duration
    Throughput       float64
}

func NewQueueMonitor(queue MessageQueue) *QueueMonitor
func (m *QueueMonitor) Start(ctx context.Context)
func (m *QueueMonitor) GetMetrics() *QueueMetrics
func (m *QueueMonitor) GetMetricsHistory(duration time.Duration) []*QueueMetrics
```

#### 4.3 队列管理
**文件**: `backend/interfaces/websocket/queue_manager.go`

```go
type QueueManager struct {
    queue    MessageQueue
    monitor  *QueueMonitor
}

type QueueStatus struct {
    Size       int
    Throughput float64
    Metrics    *QueueMetrics
}

func NewQueueManager(queue MessageQueue) *QueueManager
func (m *QueueManager) GetStatus() *QueueStatus
func (m *QueueManager) Pause() error
func (m *QueueManager) Resume() error
func (m *QueueManager) Clear() error
func (m *QueueManager) GetPriorityStats() map[int]int
```

---

## 4. 时间安排

| 阶段 | 工作内容 | 预计时间 |
|------|----------|----------|
| Phase 1 | WebSocket 服务器 | 1 天 |
| Phase 2 | 实时推送机制 | 1 天 |
| Phase 3 | 连接管理 | 1 天 |
| Phase 4 | 消息队列实现 | 1 天 |

**总计**: 3-5 天

---

## 5. 验收标准

### 5.1 功能验收
- [ ] 连接管理正常
- [ ] 消息广播正常
- [ ] 房间管理正常
- [ ] 心跳检测正常
- [ ] 连接限制有效
- [ ] 消息队列正常
- [ ] 消息路由正确
- [ ] 消息持久化正常
- [ ] 消息确认机制正常
- [ ] 消息重发机制正常
- [ ] 连接认证正常
- [ ] 连接授权正确
- [ ] 连接监控准确
- [ ] 连接限流有效
- [ ] 连接清理正常

### 5.2 性能验收
- [ ] 连接建立时间 < 100ms
- [ ] 消息发送延迟 < 50ms
- [ ] 支持并发连接数 ≥ 1000
- [ ] 消息吞吐量 ≥ 10000 msg/s

### 5.3 质量验收
- [ ] 单元测试覆盖率 ≥ 80%
- [ ] 集成测试覆盖主要场景
- [ ] 代码审查通过
- [ ] 文档完整

---

## 6. 风险与缓解

### 6.1 技术风险
**风险**: 连接断开频繁
**缓解**:
- 实现自动重连机制
- 优化心跳检测策略
- 提供连接状态监控

### 6.2 性能风险
**风险**: 大量连接导致性能下降
**缓解**:
- 实现连接池管理
- 使用消息队列缓冲
- 优化消息处理逻辑

### 6.3 安全风险
**风险**: 未授权连接
**缓解**:
- 实现强认证机制
- 实现连接授权
- 监控异常连接行为
