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
