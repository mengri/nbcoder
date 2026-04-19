# 项目目录结构与规范

## 目录结构

```
nbcoder/
├── backend/         # 后端代码（DDD 规范）
│   ├── domain/          # 领域层
│   ├── application/     # 应用层
│   ├── infrastructure/  # 基础设施层
│   └── interfaces/      # 接口层
├── web/             # 前端代码根目录
└── agents.md        # 目录结构与开发规范说明
```

## 规范说明

- 后端 backend 采用 DDD（领域驱动设计）分层：
  - domain：核心业务、实体、聚合、值对象、领域服务、仓储接口
  - application：用例、应用服务、DTO、协调领域对象完成业务流程
  - infrastructure：与外部资源交互，如数据库、消息中间件、第三方服务
  - interfaces：对外暴露 API、UI、消息等，接收外部请求并调用应用层
- 前端 web 目录为前端项目根目录，结构和技术栈可根据实际需求扩展
- agents.md 用于记录项目结构、命名规范、开发约定等说明

## 命名规范

- 目录、文件、类、方法、变量命名统一采用小写+下划线（snake_case）或小驼峰（camelCase），保持风格一致。
- 领域层（domain）实体、值对象、服务等建议用英文单数名词。
- 应用层（application）用例、服务以“*Service”或“*UseCase”结尾。
- 基础设施层（infrastructure）实现类以“*Repository”、“*Adapter”等结尾。
- 接口层（interfaces）API 文件以“api_*.py”或“*_controller.js”等命名。

## 开发约定

- 每次开发按任务拆分，完成一个任务后需 commit 并 push。
- 代码需包含必要注释，重要业务逻辑需配合注释说明。
- 领域层禁止直接依赖基础设施层，依赖倒置通过接口实现。
- 前后端分离，接口通过 API 文档（如 OpenAPI/Swagger）约定。
- 统一使用 git 进行版本管理，分支命名建议 feature/xxx、fix/xxx、refactor/xxx。
- 推荐使用单元测试，测试代码与业务代码分离。

## 其他说明

- 可根据实际业务扩展目录结构，如增加 docs、scripts、tests 等。
- agents.md 持续维护，团队成员有新约定需及时补充。
