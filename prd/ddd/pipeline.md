# 流水线领域（Pipeline Domain）需求说明

## 领域职责
- 负责需求卡片的 7 阶段开发流水线编排与执行控制。
- 管理每个阶段的状态、审核模式、执行记录。

## 主要实体
- Pipeline（流水线）
- Stage（阶段）
- StageRecord（阶段执行记录）
- StageStatus（阶段状态，值对象）
- PipelineAggregate（聚合根）

## 主要功能需求
1. 支持每张卡片独立走 7 阶段流水线，阶段包括需求分析、方案设计、任务拆解、测试用例、代码开发、测试验证、评审合并。
2. 阶段可配置启用/禁用、审核模式（AI/人工/跳过）、最大重试次数。
3. 阶段状态机：未开始、进行中、已完成、失败、等待审核。
4. 支持依赖拓扑排序、批量启动、排队与解锁。
5. 阶段执行记录可追溯，支持人工/AI 审核。
6. 领域事件：StageStarted、StageCompleted、StageFailed、StageReviewRequired。

## 领域事件
- StageStarted：阶段开始
- StageCompleted：阶段完成
- StageFailed：阶段失败
- StageReviewRequired：需要人工审核

## 限界上下文
- 仅负责流水线编排与阶段控制，不直接管理卡片内容、任务执行、代码开发。
- 与 Requirement、Agent 领域通过事件和聚合解耦。
