# 任务：AR1.1 Provider与模型管理

## 用户故事
作为一名AI平台管理员，我希望能够管理不同Provider和模型，便于灵活切换和扩展AI能力。

## 领域上下文
- 所属领域：AI运行时领域（AI Runtime Domain）
- 领域职责：Provider管理、模型注册与切换。
- 相关聚合根/实体：Provider, Model
- 依赖任务：无

## 任务目标
实现Provider与模型的注册、管理与切换机制。

## 输入
- 领域建模说明（见 prd/ddd/airuntime.md）
- Provider与模型元数据

## 输出
- Provider与模型管理实现代码

## 实现要点
- 支持多Provider与多模型注册
- 切换与扩展灵活
- 与Agent/流水线集成

## 交付标准
- 代码可独立运行，单元测试覆盖主要管理流程
- 结构与领域建模一致
