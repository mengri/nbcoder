# AI 运行时领域（AI Runtime Domain）需求说明

## 领域职责
- 负责 AI Provider 管理、模型路由、Token 计费与降级。
- 支持多 Provider、多模型、模型链配置。

## 主要实体
- Provider（AI 服务商）
- Model（模型）
- Chain（模型链）

## 主要功能需求
1. 支持多 Provider 管理，API Key 加密存储。
2. 支持每类 Agent 配置模型链，按类型路由。
3. 支持 Provider/模型可用性检测与自动降级。
4. 支持 Token 计费与调用日志。
5. 领域事件：ModelCalled、ModelFailed。

## 领域事件
- ModelCalled：模型被调用
- ModelFailed：模型调用失败

## 限界上下文
- 仅管理 AI Provider 与模型，不直接参与任务调度、代码开发。
- 与 Agent 领域通过端口和事件解耦。
