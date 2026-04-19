# 后端技术架构说明

## 总体架构
- 采用 DDD（领域驱动设计）分层，分为 domain、application、infrastructure、interfaces 四层
- 领域层（domain）聚焦核心业务建模，定义实体、聚合、值对象、领域服务、仓储接口
- 应用层（application）负责用例编排、DTO、应用服务，协调领域对象完成业务流程
- 基础设施层（infrastructure）实现外部资源适配，如数据库、消息队列、第三方服务
- 接口层（interfaces）对外暴露 API（REST/gRPC）、消息、UI 适配等

## 技术选型
- 语言：Golang 1.20+
- Web 框架：Gin（github.com/gin-gonic/gin）
- ORM：GORM、sqlx 或原生 database/sql
- 配置管理：Viper、env
- 日志：zap、logrus
- 测试：Go test、testify
- 依赖注入：github.com/mengri/autowire
- API 文档：Swagger/OpenAPI
- 任务调度/异步：可选集成 cron、RabbitMQ、Kafka
- 任务调度/定时：Golang 内置 time 包，无需第三方依赖

## 关键设计原则
- 领域层无基础设施依赖，依赖倒置通过接口实现
- 领域事件驱动，支持扩展与解耦
- 统一错误处理与日志
- 单元测试与集成测试并重
- 支持本地开发与容器化部署

## 典型调用链
1. 外部请求（API/消息）进入 interfaces 层
2. interfaces 调用 application 层用例服务
3. application 协调 domain 层实体/服务完成业务
4. 需要持久化/外部资源时，调用 infrastructure 层实现
5. 结果返回 interfaces 层，响应外部

## 目录结构示例
```
backend/
├── domain/          # 领域模型
├── application/     # 应用服务
├── infrastructure/  # 基础设施实现
├── interfaces/      # API/消息/适配层
├── TECH_SPEC.md     # 技术规范
├── ARCHITECTURE.md  # 架构说明
```

## 扩展建议
- 可根据业务复杂度增加 docs、scripts、tests、deploy 等目录
- agents.md、TECH_SPEC.md、ARCHITECTURE.md 持续维护
