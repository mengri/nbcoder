# NBCoder Web Frontend

NBCoder 系统的 Web 前端界面，基于 Vue 3 + TypeScript 开发。

## 技术栈

- Vue 3.4+ (Composition API)
- TypeScript 5.0+
- Vite 5.0+
- Element Plus 2.4+
- Pinia 2.1+
- Vue Router 4.2+
- Axios 1.6+
- ECharts 5.4+

## 开发

```bash
# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 构建生产版本
npm run build

# 预览生产构建
npm run preview

# 代码检查
npm run lint

# 代码格式化
npm run format
```

## 项目结构

```
web/
├── src/
│   ├── api/           # API 接口
│   ├── assets/        # 静态资源
│   ├── components/    # 组件
│   │   ├── business/  # 业务组件
│   │   ├── common/    # 通用组件
│   │   └── layout/    # 布局组件
│   ├── composables/   # 组合式函数
│   ├── layouts/       # 页面布局
│   ├── router/        # 路由配置
│   ├── stores/        # 状态管理
│   ├── styles/        # 样式文件
│   ├── types/         # TypeScript 类型定义
│   ├── utils/         # 工具函数
│   └── views/         # 页面视图
│       ├── Agent/     # Agent 监控
│       ├── AI/        # AI Runtime 管理
│       ├── Card/      # 需求卡片
│       ├── Knowledge/ # 知识库
│       ├── Pipeline/  # Pipeline 监控
│       ├── Project/   # 项目管理
│       └── Login/     # 登录页
├── index.html
├── package.json
├── tsconfig.json
└── vite.config.ts
```

## 功能模块

- 项目管理：创建、编辑、删除项目，配置项目信息
- 需求卡片：管理需求卡片，支持卡片池视图和详情
- Pipeline 监控：查看 Pipeline 执行状态和详细日志
- Agent 执行：实时监控 Agent 任务执行情况
- 知识库：管理文档和索引状态
- AI Runtime：配置 AI Provider 和查看调用统计

## 开发规范

- 使用 Composition API (`<script setup>`)
- 所有组件都要有 TypeScript 类型定义
- 使用 Pinia 进行状态管理
- 代码风格符合 ESLint 和 Prettier 规范
- 组件命名使用 PascalCase
- 文件命名使用 kebab-case
