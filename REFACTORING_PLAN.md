# 仓储重构完成计划

## ✅ 已完成的任务

### 1. 核心架构重构
- [x] 创建DBProvider接口
- [x] 实现DatabaseManager的动态数据库管理
- [x] 实现项目独立SQLite数据库机制
- [x] 修改项目实体结构（去掉ID字段，使用Name作为标识）

### 2. 已完成重构的仓储文件（17个）
- [x] project_repo.go - 项目信息（文件系统存储）
- [x] project_config_repo.go - 项目配置
- [x] standards_repo.go - 开发规范
- [x] dev_standard_repo.go - 开发标准
- [x] branch_policy_config_repo.go - 分支策略配置
- [x] project_lifecycle_repo.go - 项目生命周期
- [x] config_change_log_repo.go - 配置变更日志
- [x] card_repo.go - 需求卡片
- [x] card_dependency_repo.go - 卡片依赖
- [x] task_repo.go - 任务
- [x] pipeline_repo.go - 流水线
- [x] stage_record_repo.go - 阶段记录
- [x] agent_execution_repo.go - 代理执行
- [x] knowledge_repo.go - 文档和目录
- [x] knowledge_chunk_repo.go - 文档分块和索引
- [x] git_clone_repo.go - 拉拉实例和PR
- [x] skill_repo.go - 代理技能

## 🔧 待完成的任务

### 任务1：修改Domain层接口签名（高优先级）

**问题：** 仓储接口参数已变更，但domain层接口未同步更新

**需要修改的文件：**
1. `domain/project/repository.go`
   - DevStandardRepo.Delete参数：Delete(id string) → Delete(id string, projectName string)
   - BranchPolicyConfigRepo.Delete参数：Delete(id string) → Delete(id string, projectName string)
   - ConfigChangeLogRepo接口方法参数调整

2. `domain/requirement/card.go`
   - Card实体字段：ProjectID → ProjectName
   - CardRepo接口方法参数调整：
     - FindByID(id string) → FindByID(id string, projectName string)
     - Delete(id string) → Delete(id string, projectName string)
     - FindByProjectID(projectID string) → FindByProjectName(projectName string)

3. `domain/git/repository.go`
   - PullRequest实体字段：ProjectID → ProjectName
   - PullRequestRepo接口方法参数调整：
     - FindByID(id string) → FindByID(id string, projectName string)
     - FindByProjectID(projectID string) → FindByProjectName(projectName string)

**验证方式：** 运行 `go build ./cmd/server` 检查编译错误

### 任务2：完成剩余4个全局仓储重构（中优先级）

**这些仓储使用全局数据库，无需projectName参数：**
1. `notify_repo.go` - 通知和订阅
2. `notify_pref_repo.go` - 通知偏好、模板和历史
3. `airuntime_repo.go` - 提供者和模型
4. `airuntime_chain_repo.go` - 模型链和调用日志

**重构模式：**
```go
type XXXRepo struct {
    dbProvider DBProvider
}

func NewXXXRepo(dbProvider DBProvider) XXXRepo {
    return &XXXRepo{dbProvider: dbProvider}
}

func (r *XXXRepo) getDB() (*gorm.DB, error) {
    return r.dbProvider.GetGlobalDB(), nil
}

// 所有方法：
func (r *XXXRepo) SomeMethod(...) error {
    db, err := r.getDB()
    if err != nil {
        return err
    }
    // 使用db而不是r.db
}
```

### 任务3：修复仓储实现中的接口不匹配问题（高优先级）

**当前编译错误需要修复：**

1. `dev_standard_repo.go` - Delete方法签名不匹配
   - 当前：Delete(id string, projectName string) error
   - 需要：检查domain接口期望的签名

2. `branch_policy_config_repo.go` - Delete方法签名不匹配
   - 当前：Delete(id string, projectName string) error
   - 需要：检查domain接口期望的签名

3. `card_repo.go` - 实体字段不匹配
   - 需要：修改requirement.Card实体的ProjectID字段为ProjectName

4. `git_clone_repo.go` - FindByID方法签名不匹配
   - 当前：FindByID(id string, projectName string)
   - 需要：检查domain接口期望的签名

### 任务4：更新Repositories构造函数（已完成✅）

**已经完成：**
- [x] 修改NewRepositories接受DBProvider参数
- [x] 更新所有仓储实例化调用

### 任务5：更新主程序启动代码（已完成✅）

**已经完成：**
- [x] 修改cmd/server/main.go使用DatabaseManager
- [x] 传递dbManager给ProjectService

### 任务6：编译和测试验证（高优先级）

**步骤：**
1. 运行 `go build ./cmd/server` 修复所有编译错误
2. 运行 `go test ./...` 确保测试通过
3. 创建测试项目，验证数据库隔离功能：
   - 创建项目test_project
   - 检查是否生成 `projects/test_project/nbcoder.db`
   - 验证项目数据存储在独立数据库中

### 任务7：前端API适配（中优先级）

**需要修改的文件：**
1. `web/src/api/project.ts` - 调用项目API
2. `web/src/types/project.ts` - 项目类型定义
3. 其他需要projectId的地方改为使用projectName

**API路由变更：**
- `/api/v1/projects/:id` → `/api/v1/projects/:name`
- 所有相关接口路径参数从id改为name

## 📊 完成进度

- **总任务数：** 7个主要任务
- **已完成：** 3个任务（任务4、任务5、部分任务1）
- **进行中：** 1个任务（任务1）
- **待开始：** 3个任务（任务2、任务3、任务6、任务7）

## 🎯 优先级建议

**P0（必须立即完成）：**
1. 任务1：修改Domain层接口签名（解决编译错误）
2. 任务3：修复仓储实现中的接口不匹配问题
3. 任务6：编译和测试验证

**P1（重要但不紧急）：**
4. 任务2：完成剩余4个全局仓储重构
5. 任务7：前端API适配

## 🚀 执行顺序建议

1. **第一步：** 修改Domain层接口签名（任务1）
2. **第二步：** 修复仓储实现中的接口不匹配（任务3）
3. **第三步：** 编译验证，解决所有编译错误
4. **第四步：** 完成剩余4个全局仓储重构（任务2）
5. **第五步：** 最终编译和功能测试（任务6）
6. **第六步：** 前端API适配（任务7）

## 📝 注意事项

1. **接口一致性：** 确保domain层接口和仓储实现完全匹配
2. **字段映射：** 所有ProjectID字段都已改为ProjectName
3. **数据库隔离：** 确保每个项目使用独立的SQLite数据库
4. **向后兼容：** 考虑是否需要保留API兼容性
5. **测试覆盖：** 完成后需要添加数据库隔离的测试用例

## 🔗 相关文件清单

### 需要修改的Domain层文件
- `domain/project/repository.go`
- `domain/requirement/card.go`
- `domain/git/repository.go`
- `domain/clonepool/repository.go`（如有ProjectID）
- `domain/knowledge/repository.go`（如有ProjectID）

### 需要重构的仓储文件
- `infrastructure/persistence/sqlite/notify_repo.go`
- `infrastructure/persistence/sqlite/notify_pref_repo.go`
- `infrastructure/persistence/sqlite/airuntime_repo.go`
- `infrastructure/persistence/sqlite/airuntime_chain_repo.go`

### 需要适配的前端文件
- `web/src/api/project.ts`
- `web/src/types/project.ts`
- 其他引用projectId的组件

## ✨ 完成后的架构特性

1. **项目完全隔离：** 每个项目独立的数据库文件
2. **数据存储分离：** 项目信息存储在文件系统，项目数据存储在数据库
3. **动态数据库管理：** 运行时动态创建和获取项目数据库
4. **统一接口：** 通过DBProvider接口统一管理数据库连接
5. **易于扩展：** 新增项目时自动创建数据库和表结构

---

**创建时间：** 2024-04-21
**当前状态：** 核心架构完成，接口适配和剩余重构待完成
**预计完成时间：** 1-2小时
