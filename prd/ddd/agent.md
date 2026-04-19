# 智能体领域（Agent Domain）需求说明

## 领域职责
- 负责 Agent 调度、Skill 管理、任务分派与执行。
- 管理任务生命周期、状态、执行日志。

## 主要实体
- Task（任务）
- Skill（技能）
- AgentExecution（Agent 执行记录）
- TaskStatus（任务状态，值对象）
- AgentType（Agent 类型，值对象）
- TaskAggregate（聚合根）

## 主要功能需求
1. 支持多类型 Agent（产品/架构/管理/技术栈），按任务类型自动调度。
2. 任务状态机：待分配、进行中、已完成、失败、中断。
3. Skill 封装与 tool 调用，支持独立上下文与模型配置。
4. 任务执行日志完整记录，支持追溯与审计。
5. 领域事件：TaskAssigned、TaskStarted、TaskCompleted、TaskFailed、TaskInterrupted。

## 领域事件
- TaskAssigned：任务分配
- TaskStarted：任务开始
- TaskCompleted：任务完成
- TaskFailed：任务失败
- TaskInterrupted：任务中断

## 限界上下文
- 仅负责任务调度与执行，不直接管理卡片内容、流水线编排、代码实现。
- 与 Pipeline、ClonePool、AIRuntime 领域通过端口和事件解耦。
