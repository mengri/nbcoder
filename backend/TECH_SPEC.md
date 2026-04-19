# 后端 Golang 技术规范

## 目录结构
- 遵循 DDD（领域驱动设计）分层：domain、application、infrastructure、interfaces
- 每层单独目录，领域模型、服务、接口分离

## 代码风格
- 统一使用 go mod 管理依赖
- 包名、文件名、变量名、函数名统一小写+下划线（snake_case）或小驼峰（camelCase）
- 结构体、接口、方法命名采用大驼峰（PascalCase）
- 每个包应有独立的 doc.go 或 README 注释说明

## 领域建模
- 领域实体、值对象、聚合根、领域服务分离
- 领域事件、仓储接口定义在 domain 层
- 禁止 domain 直接依赖 infrastructure

## 错误处理
- 错误优先返回，使用 errors.Wrap/WithStack 保留上下文
- 日志统一用 zap/logrus 等主流库

## 单元测试
- 每个包应有 _test.go 文件，覆盖主要业务逻辑
- 推荐使用 testify/assert 进行断言

## 版本管理
- 统一使用 git，分支命名 feature/xxx、fix/xxx、refactor/xxx
- 每次开发按任务拆分，完成后 commit+push

## 依赖管理
- 禁止直接 vendor，统一 go mod tidy
- 外部依赖需在 PR/commit 说明用途

## 其他
- 代码需包含必要注释，重要业务逻辑需配合注释说明
- agents.md、TECH_SPEC.md 持续维护，团队成员有新约定需及时补充
