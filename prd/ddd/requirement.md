# 需求领域（Requirement Domain）需求说明

## 领域职责
- 负责需求卡片的全生命周期管理，包括创建、状态流转、依赖关系、变更与废弃。
- 支持碎片化输入、依赖分析、批量生成、变更预览、卡片池管理。

## 主要实体
- Card（卡片）
- CardDependency（卡片依赖）
- CardStatus（卡片状态，值对象）
- CardAggregate（聚合根）

## 主要功能需求
1. 支持自然语言输入生成结构化需求卡片，保留原始表达。
2. 卡片状态机：草稿、已确认、进行中、已完成、已变更、已废弃。
3. 支持卡片间依赖关系，自动阻塞/解锁。
4. 支持卡片变更、废弃、被替代，保留原文与血缘关系。
5. 支持批量生成卡片与依赖分析。
6. 变更预览：净变更、新增/变更/废弃/影响警告。
7. 卡片池视图：多状态筛选、排序、批量操作。
8. 领域事件：CardCreated、CardConfirmed、CardSuperseded、CardAbandoned。

## 领域事件
- CardCreated：新建卡片
- CardConfirmed：卡片确认
- CardSuperseded：卡片被替代
- CardAbandoned：卡片废弃

## 限界上下文
- 仅管理需求卡片及其依赖，不涉及具体实现、任务拆解、代码开发。
- 与 Pipeline、Agent 领域通过聚合和事件解耦。
