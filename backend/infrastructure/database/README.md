# NBCoder 数据库持久化系统

## 概述

NBCoder 系统已从 InMemory 内存存储迁移到 SQLite 持久化存储，确保数据在服务重启后不会丢失，支持生产环境部署。

## 架构

### 目录结构

```
backend/infrastructure/
├── database/
│   ├── manager.go          # 数据库连接管理器
│   ├── transaction.go      # 事务管理器
│   ├── schema.go           # 数据库初始化
│   ├── models/             # GORM 数据模型
│   │   ├── project_models.go
│   │   ├── requirement_models.go
│   │   ├── pipeline_models.go
│   │   ├── agent_models.go
│   │   ├── airuntime_models.go
│   │   ├── knowledge_models.go
│   │   ├── git_models.go
│   │   ├── notify_models.go
│   │   └── clonepool_models.go
│   └── migrations/         # 数据库迁移脚本
│       ├── migrator.go
│       ├── 000001_create_initial_schema.up.sql
│       ├── 000001_create_initial_schema.down.sql
│       ├── 000002_add_indexes.up.sql
│       ├── 000002_add_indexes.down.sql
│       ├── 000003_add_ai_runtime_tables.up.sql
│       ├── 000003_add_ai_runtime_tables.down.sql
│       ├── 000004_add_knowledge_git_notify_tables.up.sql
│       └── 000004_add_knowledge_git_notify_tables.down.sql
└── persistence/
    └── sqlite/             # SQLite Repository 实现
        ├── project_repo.go
        ├── project_config_repo.go
        ├── standards_repo.go
        ├── dev_standard_repo.go
        ├── branch_policy_config_repo.go
        ├── project_lifecycle_repo.go
        ├── config_change_log_repo.go
        ├── card_repo.go
        ├── card_dependency_repo.go
        ├── pipeline_repo.go
        ├── stage_record_repo.go
        ├── task_repo.go
        ├── agent_execution_repo.go
        ├── skill_repo.go
        ├── airuntime_repo.go
        ├── airuntime_chain_repo.go
        ├── knowledge_repo.go
        ├── knowledge_chunk_repo.go
        ├── git_clone_repo.go
        ├── notify_repo.go
        ├── notify_pref_repo.go
        ├── repositories.go
        └── repo_test.go
```

### 数据库表结构

系统包含 30 个数据表，分为以下领域：

1. **项目相关** (7个表)
   - projects: 项目主表
   - project_configs: 项目配置
   - standards: 开发标准
   - dev_standards: 开发标准详情
   - branch_policy_configs: 分支策略配置
   - project_lifecycles: 项目生命周期
   - config_change_logs: 配置变更日志

2. **需求相关** (2个表)
   - cards: 需求卡片
   - card_dependencies: 卡片依赖关系

3. **流水线相关** (2个表)
   - pipelines: 流水线
   - stage_records: 阶段记录

4. **Agent 相关** (3个表)
   - tasks: 任务
   - agent_executions: Agent 执行记录
   - skills: 技能

5. **AI Runtime 相关** (4个表)
   - providers: AI 提供商
   - models: AI 模型
   - model_chains: 模型链
   - call_logs: 调用日志

6. **知识库相关** (4个表)
   - documents: 文档
   - document_chunks: 文档分块
   - document_indices: 文档索引
   - directories: 目录

7. **Git 相关** (3个表)
   - repositories: 代码仓库
   - pull_requests: 拉取请求
   - commits: 提交记录

8. **通知相关** (5个表)
   - notifications: 通知
   - subscriptions: 订阅
   - subscription_preferences: 订阅偏好
   - notification_templates: 通知模板
   - notification_histories: 通知历史

## 使用方法

### 1. 初始化数据库连接

```go
import "github.com/mengri/nbcoder/infrastructure/database"

// 使用默认配置
dbManager, err := database.NewDatabaseManager(database.DefaultDatabaseConfig())
if err != nil {
    log.Fatal(err)
}
defer dbManager.Close()

// 初始化数据库表结构
err = database.InitSchema(dbManager.GetDB())
if err != nil {
    log.Fatal(err)
}
```

### 2. 使用 Repository

```go
import "github.com/mengri/nbcoder/infrastructure/persistence/sqlite"

// 创建所有 Repository
repos := sqlite.NewRepositories(dbManager.GetDB())

// 使用 Project Repository
project := domain.NewProject("id", "name", "desc", "repo_url")
err := repos.Project.Save(project)

// 查询项目
found, err := repos.Project.FindByID("id")
```

### 3. 使用事务

```go
import "github.com/mengri/nbcoder/infrastructure/database"

tm := database.NewTransactionManager(dbManager.GetDB())

err := tm.Execute(nil, func(ctx context.Context, txDB *gorm.DB) error {
    // 在事务中执行操作
    project := domain.NewProject("id", "name", "desc", "repo_url")
    err := repos.Project.Save(project)
    return err
})
```

### 4. 数据库迁移

```go
import "github.com/mengri/nbcoder/infrastructure/database/migrations"

// 创建迁移器
migrator := migrations.NewMigrator(dbManager.GetDB(), nil)

// 初始化迁移系统
err := migrator.Init()
if err != nil {
    log.Fatal(err)
}

// 执行迁移
err = migrator.Up()
if err != nil {
    log.Fatal(err)
}

// 查看迁移状态
status, err := migrator.Status()
if err != nil {
    log.Fatal(err)
}

// 回滚最后一个迁移
err = migrator.Down()
if err != nil {
    log.Fatal(err)
}
```

## 配置

### 数据库配置

```go
config := &database.DatabaseConfig{
    DSN:                ".NBCoder-global/nbcoder.db",  // 数据库文件路径
    MaxOpenConnections: 25,                             // 最大打开连接数
    MaxIdleConnections: 5,                              // 最大空闲连接数
    ConnectionMaxLife:  time.Hour,                      // 连接最大生命周期
    MaxIdleTime:        10 * time.Minute,               // 最大空闲时间
    DisableForeignKey:  false,                          // 是否禁用外键约束
    DisableAutoMigrate: false,                          // 是否禁用自动迁移
    LogLevel:           gorm.Info,                      // 日志级别
}
```

### 默认配置

```go
config := database.DefaultDatabaseConfig()
```

## 数据模型

所有数据模型都遵循以下原则：

1. 使用 `gorm.Model` 作为基础模型
2. 所有表包含 `created_at`、`updated_at` 字段
3. 使用 `gorm:"softDelete"` 支持软删除
4. 外键关系使用 `gorm:"foreignKey"`
5. 添加必要的索引以提升查询性能
6. 使用 JSON 类型存储复杂数据结构

## Repository 实现

系统为每个领域提供了完整的 SQLite Repository 实现：

- Project Domain: 7 个 Repository
- Requirement Domain: 2 个 Repository
- Pipeline Domain: 2 个 Repository
- Agent Domain: 3 个 Repository
- AI Runtime Domain: 4 个 Repository
- Knowledge Domain: 4 个 Repository
- Clone Pool Domain: 2 个 Repository
- Notify Domain: 5 个 Repository
- Git Domain: 3 个 Repository

## 测试

### 运行所有测试

```bash
cd backend
go test ./infrastructure/database/...
go test ./infrastructure/persistence/sqlite/...
```

### 运行特定测试

```bash
go test -run TestProjectRepo_Save ./infrastructure/persistence/sqlite/
```

## 性能优化

1. **索引优化**: 所有常用查询字段都已添加索引
2. **连接池**: 使用 GORM 内置连接池管理
3. **批量操作**: 支持批量插入和更新
4. **软删除**: 使用软删除避免数据丢失
5. **查询优化**: 使用预加载减少 N+1 查询

## 注意事项

1. **并发访问**: SQLite 默认支持并发读，写操作会获取锁
2. **事务处理**: 复杂操作建议使用事务保证数据一致性
3. **备份**: 定期备份数据库文件
4. **迁移**: 使用迁移脚本管理数据库结构变更
5. **测试**: 生产环境部署前务必进行完整测试

## 迁移指南

从 InMemory 迁移到 SQLite：

1. 安装依赖: `go get gorm.io/gorm gorm.io/driver/sqlite`
2. 初始化数据库连接
3. 初始化数据库表结构
4. 替换 InMemory Repository 为 SQLite Repository
5. 运行测试验证功能
6. 部署到生产环境

## 故障排除

### 数据库锁定

如果遇到数据库锁定错误：
- 检查是否有长时间运行的写操作
- 增加连接超时时间
- 使用事务确保操作的原子性

### 性能问题

如果遇到性能问题：
- 检查慢查询日志
- 优化索引
- 考虑使用批量操作
- 检查连接池配置

## 未来改进

1. 支持其他数据库 (PostgreSQL, MySQL)
2. 添加数据库连接池监控
3. 实现数据库备份和恢复
4. 添加查询性能分析
5. 支持数据库读写分离
