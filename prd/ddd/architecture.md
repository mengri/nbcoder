# NBCoder DDD领域架构设计

本文件对各领域进行架构设计，包含领域模型、核心接口、聚合、限界上下文、交互关系与技术选型建议。

---

## 1. 需求领域（Requirement Domain）
- **聚合根**：CardAggregate
- **核心实体**：Card, CardDependency, CardStatus
- **限界上下文**：仅管理需求卡片及依赖，事件驱动与Pipeline、Agent解耦
- **核心接口**：
  - 创建/编辑/废弃卡片
  - 卡片状态流转（状态机）
  - 依赖关系管理与变更分析
  - 领域事件发布（CardCreated等）
- **存储建议**：SQLite表（card, card_dependency），支持事件溯源
- **交互关系**：
  - Pipeline消费CardAggregate
  - 变更事件通知Agent、Notify

---

## 2. 流水线领域（Pipeline Domain）
- **聚合根**：PipelineAggregate
- **核心实体**：Pipeline, Stage, StageRecord, StageStatus
- **限界上下文**：编排卡片开发全流程，阶段与审核配置可扩展
- **核心接口**：
  - 流水线/阶段定义与配置
  - 阶段状态流转与执行控制
  - 阶段审核（AI/人工/跳过）
  - 阶段执行记录与追溯
  - 领域事件发布（StageStarted等）
- **存储建议**：SQLite表（pipeline, stage, stage_record）
- **交互关系**：
  - 监听Card事件启动流水线
  - 阶段事件驱动Agent任务

---

## 3. 智能体领域（Agent Domain）
- **聚合根**：TaskAggregate
- **核心实体**：Task, Skill, AgentExecution, TaskStatus, AgentType
- **限界上下文**：任务调度与Skill封装，独立上下文与模型配置
- **核心接口**：
  - 任务分派与状态流转
  - 多类型Agent调度与Skill调用
  - 任务执行日志与追溯
  - 领域事件发布（TaskAssigned等）
- **存储建议**：SQLite表（task, agent_execution, skill）+ 文件日志
- **交互关系**：
  - 消费Pipeline阶段事件
  - 通过端口访问ClonePool、AIRuntime

---

## 4. 知识库领域（Knowledge Domain）
- **聚合根**：Document
- **核心实体**：Document, DocumentIndex, Chunk
- **限界上下文**：文档管理与检索，RAG上下文注入
- **核心接口**：
  - 文档上传/分片/索引/检索
  - 片段注入Agent上下文
  - 变更追溯与血缘管理
- **存储建议**：本地文件+SQLite索引表
- **交互关系**：
  - Agent按需检索知识片段

---

## 5. 克隆池领域（ClonePool Domain）
- **聚合根**：ClonePool
- **核心实体**：CloneInstance, Repository
- **限界上下文**：仓库克隆实例管理，状态流转与异常恢复
- **核心接口**：
  - 克隆实例分配/回收/校验
  - 状态流转（idle/busy/dirty）
  - 事件发布（CloneAcquired等）
- **存储建议**：本地目录结构+SQLite表（clone_instance, repository）
- **交互关系**：
  - Agent通过端口获取/归还克隆实例

---

## 6. AI运行时领域（AI Runtime Domain）
- **聚合根**：Provider
- **核心实体**：Provider, Model, Chain
- **限界上下文**：AI Provider与模型链管理，Token计费与降级
- **核心接口**：
  - Provider/模型注册与配置
  - 模型链路由与降级
  - Token计费与调用日志
  - 事件发布（ModelCalled等）
- **存储建议**：SQLite表（provider, model, chain, call_log）
- **交互关系**：
  - Agent通过端口调用模型

---

## 7. 项目领域（Project Domain）
- **聚合根**：Project
- **核心实体**：Project, ProjectConfig, Standards
- **限界上下文**：项目元数据与配置管理，多项目隔离
- **核心接口**：
  - 项目创建/配置/生命周期管理
  - 开发规范与分支策略配置
- **存储建议**：本地目录+SQLite表（project, project_config, standards）
- **交互关系**：
  - 领域配置注入Requirement、Git等

---

## 8. Git领域（Git Domain）
- **聚合根**：Repository
- **核心实体**：PullRequest, Branch, CommitHistory
- **限界上下文**：仓库抽象与PR/MR管理，分支策略
- **核心接口**：
  - 分支策略与命名规范
  - PR/MR创建与描述生成
  - 提交历史追溯与squash合并
- **存储建议**：Git仓库+SQLite表（pr, branch, commit）
- **交互关系**：
  - 与Project、Agent领域通过端口集成

---

## 9. 通知领域（Notify Domain）
- **聚合根**：Notification
- **核心实体**：Notification, Channel
- **限界上下文**：多渠道通知推送，事件订阅与模板
- **核心接口**：
  - 通知推送（WebSocket/邮件/系统）
  - 事件订阅与屏蔽
  - 通知模板与历史追溯
- **存储建议**：SQLite表（notification, channel, template）
- **交互关系**：
  - 监听各领域事件，推送通知

---

> 各领域均建议采用领域服务+事件驱动+接口适配器分层，聚合根严格控制不变式，限界上下文通过端口/事件解耦。