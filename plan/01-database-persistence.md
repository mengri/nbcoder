# 数据库持久化实现实施计划

> 任务编号: 1. 数据库持久化实现 [P0]
> 预估工作量: 1-2 周
> 阻塞影响: 所有功能无法持久化

---

## 1. 概述

将 NBCoder 系统从 InMemory 内存存储迁移到 SQLite 持久化存储，确保数据在服务重启后不会丢失，支持生产环境部署。

---

## 2. 技术选型

### 2.1 数据库选型
- **主数据库**: SQLite 3.x
- **ORM 框架**: GORM v2 (已有依赖)
- **迁移工具**: gormigrate (可选) 或自建迁移机制
- **连接池**: SQLite 内置连接池 (无需额外配置)

### 2.2 选型理由
- SQLite 零配置，适合私有化部署
- 单文件存储，便于备份和迁移
- GORM 成熟稳定，支持丰富的数据库操作
- 支持 ACID 事务，保证数据一致性
- 兼容性强，支持多平台

---

## 3. 架构设计

### 3.1 数据库架构
```
┌─────────────────────────────────────────────────────────┐
│                    数据库层                              │
├─────────────────────────────────────────────────────────┤
│  SQLite Connection Pool                               │
│  ├── Connection Manager                               │
│  ├── Transaction Manager                               │
│  └── Migration Engine                                 │
└─────────────────────────────────────────────────────────┘
                          │
┌─────────────────────────────────────────────────────────┐
│                    ORM 层 (GORM)                         │
├─────────────────────────────────────────────────────────┤
│  ├── Model Definitions                                │
│  ├── Repository Implementations                       │
│  ├── Query Builders                                    │
│  └── Transaction Wrappers                              │
└─────────────────────────────────────────────────────────┘
                          │
┌─────────────────────────────────────────────────────────┐
│                  Domain Layer                           │
│  (保持不变，通过接口访问数据)                          │
└─────────────────────────────────────────────────────────┘
```

### 3.2 数据库目录结构
```
{workdir}/.NBCoder-global/
├── nbcoder.db              # 主数据库文件
├── nbcoder.db-shm           # SQLite 共享内存文件
├── nbcoder.db-wal           # SQLite 写前日志
├── migrations/              # 数据库迁移脚本
│   ├── 000001_init.down.sql
│   ├── 000001_init.up.sql
│   └── ...
└── backups/                 # 数据库备份目录
    ├── nbcoder_20260420_120000.db
    └── ...
```

---

## 4. 实施步骤

### 4.1 Phase 1: 基础设施搭建 (2-3 天)

#### 1.1 数据库连接管理器
**文件**: `backend/infrastructure/database/manager.go`

```go
type DatabaseManager struct {
    db     *gorm.DB
    config DatabaseConfig
}

type DatabaseConfig struct {
    DSN          string
    MaxOpenConns int
    MaxIdleConns int
    MaxLifetime  time.Duration
}

func NewDatabaseManager(config DatabaseConfig) (*DatabaseManager, error)
func (dm *DatabaseManager) GetDB() *gorm.DB
func (dm *DatabaseManager) Close() error
func (dm *DatabaseManager) Ping() error
```

#### 1.2 数据库初始化脚本
**文件**: `backend/infrastructure/database/schema.go`

```go
func InitSchema(db *gorm.DB) error
func DropSchema(db *gorm.DB) error
```

#### 1.3 事务管理器
**文件**: `backend/infrastructure/database/transaction.go`

```go
type TransactionManager struct {
    db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) *TransactionManager
func (tm *TransactionManager) Execute(fn func(*gorm.DB) error) error
func (tm *TransactionManager) ExecuteInTx(fn func(*gorm.DB) error) error
```

---

### 4.2 Phase 2: 数据模型定义 (3-4 天)

#### 2.1 创建 GORM 模型
**目录**: `backend/infrastructure/database/models/`

需要为以下实体创建 GORM 模型：

**项目相关**:
- `project.go` - 项目元数据
- `project_config.go` - 项目配置
- `standards.go` - 开发规范
- `config_change_log.go` - 配置变更日志

**需求相关**:
- `card.go` - 需求卡片
- `card_dependency.go` - 卡片依赖

**流水线相关**:
- `pipeline.go` - 流水线
- `stage_record.go` - 阶段记录

**Agent 相关**:
- `task.go` - Agent 任务
- `agent_execution.go` - 执行记录

**AI Runtime 相关**:
- `provider.go` - AI Provider
- `model.go` - 模型配置
- `model_chain.go` - 模型链
- `call_log.go` - 调用日志

**知识库相关**:
- `document.go` - 文档
- `document_chunk.go` - 文档分片
- `document_index.go` - 文档索引

**Git 相关**:
- `repository.go` - 代码仓库
- `pull_request.go` - PR/MR
- `commit.go` - 提交记录

**通知相关**:
- `notification.go` - 通知
- `subscription.go` - 订阅
- `template.go` - 模板

**克隆池相关**:
- `clone_instance.go` - 克隆实例

#### 2.2 模型设计原则
- 使用 `gorm.Model` 作为基础模型
- 所有表包含 `created_at`, `updated_at` 字段
- 使用 `gorm:"softDelete"` 支持软删除
- 外键关系使用 `gorm:"foreignKey"`
- 添加必要的索引以提升查询性能
- 使用 JSON 类型存储复杂数据结构

---

### 4.3 Phase 3: 数据库迁移系统 (2-3 天)

#### 3.1 迁移文件结构
**目录**: `backend/infrastructure/database/migrations/`

```
migrations/
├── 000001_create_initial_schema.up.sql
├── 000001_create_initial_schema.down.sql
├── 000002_add_indexes.up.sql
├── 000002_add_indexes.down.sql
├── 000003_add_ai_runtime_tables.up.sql
├── 000003_add_ai_runtime_tables.down.sql
└── ...
```

#### 3.2 迁移执行器
**文件**: `backend/infrastructure/database/migrator.go`

```go
type Migrator struct {
    db           *gorm.DB
    migrations   []Migration
    migrationDir string
}

type Migration struct {
    Version     string
    Name        string
    UpScript    string
    DownScript  string
}

func NewMigrator(db *gorm.DB, migrationDir string) *Migrator
func (m *Migrator) Up() error
func (m *Migrator) Down(version string) error
func (m *Migrator) Status() ([]MigrationStatus, error)
func (m *Migrator) CreateMigration(name, up, down string) error
```

#### 3.3 创建迁移脚本

**000001_create_initial_schema**:
- 创建所有基础表
- 定义外键关系
- 设置基础索引

**000002_add_indexes**:
- 为常用查询字段添加索引
- 优化查询性能

**000003_add_ai_runtime_tables**:
- AI Runtime 相关表
- Provider、Model、ModelChain、CallLog

---

### 4.4 Phase 4: SQLite Repository 实现 (4-5 天)

#### 4.1 Repository 接口实现
**目录**: `backend/infrastructure/persistence/sqlite/`

为所有 22 个 Repository 接口创建 SQLite 实现：

**基础实现模板**:
```go
type SQLiteCardRepo struct {
    db *gorm.DB
}

func NewSQLiteCardRepo(db *gorm.DB) *SQLiteCardRepo
func (r *SQLiteCardRepo) Save(card *requirement.Card) error
func (r *SQLiteCardRepo) FindByID(id string) (*requirement.Card, error)
func (r *SQLiteCardRepo) FindByProjectID(projectID string) ([]*requirement.Card, error)
func (r *SQLiteCardRepo) FindAll() ([]*requirement.Card, error)
func (r *SQLiteCardRepo) Update(card *requirement.Card) error
func (r *SQLiteCardRepo) Delete(id string) error
```

#### 4.2 需要实现的 Repository 列表

**Project Domain (5 个)**:
- `project_repo.go`
- `project_config_repo.go`
- `standards_repo.go`
- `dev_standard_repo.go`
- `branch_policy_config_repo.go`
- `project_lifecycle_repo.go`
- `config_change_log_repo.go`

**Requirement Domain (3 个)**:
- `card_repo.go`
- `card_dependency_repo.go`

**Pipeline Domain (2 个)**:
- `pipeline_repo.go`
- `stage_record_repo.go`

**Agent Domain (3 个)**:
- `task_repo.go`
- `agent_execution_repo.go`
- `skill_repo.go`

**AI Runtime Domain (3 个)**:
- `provider_repo.go`
- `model_repo.go`
- `model_chain_repo.go`
- `call_log_repo.go`

**Knowledge Domain (4 个)**:
- `document_repo.go`
- `document_chunk_repo.go`
- `document_index_repo.go`
- `directory_repo.go`

**Clone Pool Domain (2 个)**:
- `clone_instance_repo.go`
- `repository_repo.go`

**Notify Domain (5 个)**:
- `notification_repo.go`
- `subscription_repo.go`
- `subscription_preference_repo.go`
- `notification_template_repo.go`
- `notification_history_repo.go`

**Git Domain (2 个)**:
- `pull_request_repo.go`
- `commit_repo.go`

---

### 4.5 Phase 5: 集成与测试 (2-3 天)

#### 5.1 主程序集成
**文件**: `backend/cmd/server/main.go`

```go
func main() {
    // 初始化数据库
    dbConfig := database.DatabaseConfig{
        DSN:          "./data/nbcoder.db",
        MaxOpenConns: 25,
        MaxIdleConns: 5,
        MaxLifetime:  time.Hour,
    }

    dbManager, err := database.NewDatabaseManager(dbConfig)
    if err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    defer dbManager.Close()

    // 运行数据库迁移
    migrator := database.NewMigrator(dbManager.GetDB(), "./migrations")
    if err := migrator.Up(); err != nil {
        log.Fatalf("Failed to run migrations: %v", err)
    }

    // 使用 SQLite Repositories 替换 InMemory Repositories
    projectRepo := persistence.NewSQLiteProjectRepo(dbManager.GetDB())
    // ... 其他 repositories
}
```

#### 5.2 集成测试
**目录**: `backend/integration/persistence/`

创建集成测试验证：
- 数据库连接和事务
- CRUD 操作
- 复杂查询
- 外键关系
- 并发访问
- 数据迁移

---

## 5. 详细实施计划

### 5.1 Day 1-2: 数据库基础设施
- [ ] 创建 `infrastructure/database` 目录
- [ ] 实现 DatabaseManager
- [ ] 实现 TransactionManager
- [ ] 实现基础数据库初始化
- [ ] 单元测试

### 5.2 Day 3-6: 数据模型定义
- [ ] 创建 20 个 GORM 模型
- [ ] 定义外键关系
- [ ] 添加索引
- [ ] 模型验证测试

### 5.3 Day 7-9: 数据库迁移系统
- [ ] 设计迁移目录结构
- [ ] 实现 Migrator
- [ ] 创建初始迁移脚本
- [ ] 实现迁移工具 CLI

### 5.4 Day 10-14: SQLite Repository 实现
- [ ] 实现 22 个 SQLite Repository
- [ ] Repository 接口测试
- [ ] 性能优化

### 5.5 Day 15-16: 集成和测试
- [ ] 主程序集成
- [ ] 数据迁移验证
- [ ] 集成测试
- [ ] 性能测试
- [ ] 文档编写

---

## 6. 数据库表设计

### 6.1 核心表结构

#### projects 表
```sql
CREATE TABLE projects (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);

CREATE INDEX idx_projects_deleted_at ON projects(deleted_at);
CREATE INDEX idx_projects_name ON projects(name);
```

#### cards 表
```sql
CREATE TABLE cards (
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    original TEXT,
    status VARCHAR(50) NOT NULL,
    priority VARCHAR(20) NOT NULL,
    structured_output TEXT,
    pipeline_id VARCHAR(36),
    project_id VARCHAR(36) NOT NULL,
    superseded_by VARCHAR(36),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    FOREIGN KEY (project_id) REFERENCES projects(id),
    FOREIGN KEY (pipeline_id) REFERENCES pipelines(id),
    FOREIGN KEY (superseded_by) REFERENCES cards(id)
);

CREATE INDEX idx_cards_project_id ON cards(project_id);
CREATE INDEX idx_cards_status ON cards(status);
CREATE INDEX idx_cards_pipeline_id ON cards(pipeline_id);
CREATE INDEX idx_cards_deleted_at ON cards(deleted_at);
```

#### tasks 表
```sql
CREATE TABLE tasks (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(500) NOT NULL,
    description TEXT,
    task_type VARCHAR(100) NOT NULL,
    agent_type VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    priority INT DEFAULT 5,
    assigned_to VARCHAR(36),
    pipeline_id VARCHAR(36),
    project_id VARCHAR(36) NOT NULL,
    started_at DATETIME,
    completed_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    FOREIGN KEY (project_id) REFERENCES projects(id),
    FOREIGN KEY (pipeline_id) REFERENCES pipelines(id)
);

CREATE INDEX idx_tasks_project_id ON tasks(project_id);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_agent_type ON tasks(agent_type);
CREATE INDEX idx_tasks_deleted_at ON tasks(deleted_at);
```

### 6.2 完整表列表 (20 个表)

1. projects
2. project_configs
3. standards
4. dev_standards
5. branch_policy_configs
6. project_lifecycles
7. config_change_logs
8. cards
9. card_dependencies
10. pipelines
11. stage_records
12. tasks
13. agent_executions
14. providers
15. models
16. model_chains
17. call_logs
18. repositories
19. clone_instances
20. documents
21. document_chunks
22. document_indexes
23. directories
24. notifications
25. subscriptions
26. subscription_preferences
27. notification_templates
28. notification_histories
29. pull_requests
30. commits

---

## 7. 测试策略

### 7.1 单元测试
- DatabaseManager 测试
- TransactionManager 测试
- Migrator 测试
- 每个 Repository CRUD 测试

### 7.2 集成测试
- 数据库连接和迁移
- 多表关联查询
- 事务完整性
- 并发访问

### 7.3 性能测试
- 查询性能测试
- 并发访问测试
- 大数据量测试

---

## 8. 风险与缓解

### 8.1 技术风险
**风险**: SQLite 并发写入性能限制
**缓解**: 
- 使用连接池控制并发
- 优化查询和索引
- 考虑读写分离

### 8.2 数据迁移风险
**风险**: 现有数据迁移失败
**缓解**:
- 提供数据导出/导入工具
- 提供迁移回滚机制
- 充分测试迁移脚本

### 8.3 性能风险
**风险**: 大数据量查询性能
**缓解**:
- 合理设计索引
- 使用查询优化
- 考虑分页和缓存

---

## 9. 交付物

### 9.1 代码文件
- `infrastructure/database/manager.go`
- `infrastructure/database/transaction.go`
- `infrastructure/database/schema.go`
- `infrastructure/database/migrator.go`
- `infrastructure/database/models/*.go`
- `infrastructure/persistence/sqlite/*.go`

### 9.2 脚本文件
- `migrations/*.sql`
- 数据库备份/恢复脚本
- 数据迁移工具

### 9.3 文档
- 数据库架构文档
- API 使用文档
- 迁移指南
- 性能优化指南

---

## 10. 验收标准

### 10.1 功能验收
- [ ] 所有 22 个 Repository 正常工作
- [ ] 数据迁移正常执行
- [ ] 数据持久化正常
- [ ] 服务重启数据不丢失
- [ ] 事务正常工作

### 10.2 性能验收
- [ ] 常用查询响应时间 < 100ms
- [ ] 支持并发连接数 ≥ 10
- [ ] 数据库文件大小合理
- [ ] 查询优化正常

### 10.3 质量验收
- [ ] 单元测试覆盖率 ≥ 80%
- [ ] 集成测试覆盖主要场景
- [ ] 代码审查通过
- [ ] 文档完整

---

## 11. 后续工作

### 11.1 Phase 2 优化
- 数据库连接池优化
- 查询性能调优
- 数据备份自动化
- 数据库监控

### 11.2 Phase 3 扩展
- 考虑支持 PostgreSQL
- 数据库读写分离
- 数据库分片

---

## 12. 注意事项

1. **兼容性**: 确保 GORM 版本兼容
2. **备份**: 数据库变更前必须备份
3. **测试**: 充分测试后再部署
4. **文档**: 及时更新文档
5. **监控**: 增加数据库监控

---

## 13. 时间安排

| 阶段 | 工作内容 | 预计时间 |
|------|----------|----------|
| Phase 1 | 基础设施搭建 | 2-3 天 |
| Phase 2 | 数据模型定义 | 3-4 天 |
| Phase 3 | 数据库迁移系统 | 2-3 天 |
| Phase 4 | SQLite Repository 实现 | 4-5 天 |
| Phase 5 | 集成与测试 | 2-3 天 |

**总计**: 1-2 周
