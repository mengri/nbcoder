# 前端界面开发实施计划

> 任务编号: 2. 前端界面开发 [P0]
> 预估工作量: 2-4 周
> 阻塞影响: 用户无法使用系统

---

## 1. 概述

开发 NBCoder 系统的 Web 前端界面，为用户提供友好的交互界面，支持项目管理、需求卡片管理、Pipeline 执行监控、Agent 执行监控、知识库管理和 AI Runtime 管理等核心功能。

---

## 2. 技术选型

### 2.1 核心技术栈
- **框架**: Vue 3.4+ (Composition API)
- **语言**: TypeScript 5.0+
- **构建工具**: Vite 5.0+
- **UI 组件库**: Element Plus 2.4+
- **状态管理**: Pinia 2.1+
- **路由**: Vue Router 4.2+
- **HTTP 客户端**: Axios 1.6+
- **WebSocket 客户端**: 原生 WebSocket API
- **图表库**: ECharts 5.4+ (可选)

### 2.2 开发工具
- **代码规范**: ESLint + Prettier
- **Git 规范**: Husky + lint-staged
- **类型检查**: TypeScript strict mode
- **组件文档**: VitePress (可选)

### 2.3 选型理由
- Vue 3 生态成熟，性能优秀
- TypeScript 提供类型安全
- Vite 构建速度快，开发体验好
- Element Plus 组件丰富，企业级 UI
- Pinia 轻量级，API 简洁
- 全部技术栈国产化，易于获取支持

---

## 3. 项目结构

### 3.1 目录结构
```
web/
├── public/                   # 静态资源
│   ├── favicon.ico
│   └── logo.png
├── src/
│   ├── api/                 # API 接口
│   │   ├── index.ts
│   │   ├── project.ts
│   │   ├── card.ts
│   │   ├── pipeline.ts
│   │   ├── agent.ts
│   │   ├── knowledge.ts
│   │   └── ai-runtime.ts
│   ├── assets/              # 资源文件
│   │   ├── images/
│   │   ├── styles/
│   │   └── fonts/
│   ├── components/          # 通用组件
│   │   ├── common/          # 基础组件
│   │   │   ├── Button/
│   │   │   ├── Input/
│   │   │   └── Modal/
│   │   ├── layout/          # 布局组件
│   │   │   ├── Header/
│   │   │   ├── Sidebar/
│   │   │   └── Footer/
│   │   └── business/        # 业务组件
│   │       ├── CardList/
│   │       ├── PipelineView/
│   │       └── AgentMonitor/
│   ├── composables/         # 组合式函数
│   │   ├── useProject.ts
│   │   ├── useCard.ts
│   │   └── useWebSocket.ts
│   ├── layouts/             # 页面布局
│   │   ├── DefaultLayout.vue
│   │   └── EmptyLayout.vue
│   ├── router/              # 路由配置
│   │   └── index.ts
│   ├── stores/              # 状态管理
│   │   ├── project.ts
│   │   ├── card.ts
│   │   ├── pipeline.ts
│   │   └── user.ts
│   ├── types/               # TypeScript 类型
│   │   ├── project.ts
│   │   ├── card.ts
│   │   └── api.ts
│   ├── utils/               # 工具函数
│   │   ├── request.ts       # HTTP 请求封装
│   │   ├── validate.ts      # 表单验证
│   │   └── format.ts        # 格式化工具
│   ├── views/               # 页面组件
│   │   ├── Login/
│   │   ├── Project/
│   │   ├── Card/
│   │   ├── Pipeline/
│   │   ├── Agent/
│   │   ├── Knowledge/
│   │   ├── AI/
│   │   └── Settings/
│   ├── App.vue
│   └── main.ts
├── index.html
├── package.json
├── tsconfig.json
├── vite.config.ts
├── .eslintrc.cjs
├── .prettierrc
└── README.md
```

---

## 4. 实施步骤

### 4.1 Phase 1: 项目初始化 (2-3 天)

#### 1.1 项目脚手架创建
```bash
# 创建 Vite + Vue3 + TypeScript 项目
npm create vite@latest nbcoder-web -- --template vue-ts

# 进入项目目录
cd nbcoder-web

# 安装依赖
npm install

# 安装 UI 框架和状态管理
npm install element-plus @element-plus/icons-vue pinia vue-router axios

# 安装开发依赖
npm install -D @types/node sass
```

#### 1.2 配置文件设置

**vite.config.ts**:
```typescript
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'src')
    }
  },
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  }
})
```

**tsconfig.json**:
```json
{
  "compilerOptions": {
    "target": "ES2020",
    "useDefineForClassFields": true,
    "module": "ESNext",
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "skipLibCheck": true,
    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": true,
    "jsx": "preserve",
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noFallthroughCasesInSwitch": true,
    "baseUrl": ".",
    "paths": {
      "@/*": ["src/*"]
    }
  },
  "include": ["src/**/*.ts", "src/**/*.d.ts", "src/**/*.tsx", "src/**/*.vue"],
  "references": [{ "path": "./tsconfig.node.json" }]
}
```

#### 1.3 基础布局和路由设置

**router/index.ts**:
```typescript
import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import DefaultLayout from '@/layouts/DefaultLayout.vue'
import EmptyLayout from '@/layouts/EmptyLayout.vue'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login/index.vue'),
    meta: { layout: EmptyLayout }
  },
  {
    path: '/',
    component: DefaultLayout,
    children: [
      {
        path: '',
        redirect: '/projects'
      },
      {
        path: 'projects',
        name: 'Projects',
        component: () => import('@/views/Project/index.vue')
      },
      // ... 其他路由
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
```

---

### 4.2 Phase 2: 基础设施开发 (3-4 天)

#### 2.1 API 客户端封装
**utils/request.ts**:
```typescript
import axios from 'axios'
import type { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'

const instance: AxiosInstance = axios.create({
  baseURL: '/api/v1',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

instance.interceptors.request.use(
  (config) => {
    // 添加 token
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

instance.interceptors.response.use(
  (response: AxiosResponse) => {
    return response.data
  },
  (error) => {
    const message = error.response?.data?.error || error.message || '请求失败'
    ElMessage.error(message)
    return Promise.reject(error)
  }
)

export default instance
```

#### 2.2 WebSocket 客户端封装
**composables/useWebSocket.ts**:
```typescript
import { ref, onUnmounted } from 'vue'

export function useWebSocket(url: string) {
  const ws = ref<WebSocket | null>(null)
  const connected = ref(false)
  const message = ref<any>(null)

  const connect = () => {
    ws.value = new WebSocket(url)

    ws.value.onopen = () => {
      connected.value = true
    }

    ws.value.onmessage = (event) => {
      message.value = JSON.parse(event.data)
    }

    ws.value.onclose = () => {
      connected.value = false
      // 自动重连
      setTimeout(connect, 3000)
    }

    ws.value.onerror = (error) => {
      console.error('WebSocket error:', error)
    }
  }

  const disconnect = () => {
    if (ws.value) {
      ws.value.close()
      ws.value = null
    }
  }

  const send = (data: any) => {
    if (ws.value && connected.value) {
      ws.value.send(JSON.stringify(data))
    }
  }

  onUnmounted(() => {
    disconnect()
  })

  return {
    connected,
    message,
    connect,
    disconnect,
    send
  }
}
```

#### 2.3 状态管理设置
**stores/project.ts**:
```typescript
import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Project } from '@/types/project'

export const useProjectStore = defineStore('project', () => {
  const projects = ref<Project[]>([])
  const currentProject = ref<Project | null>(null)

  const setProjects = (data: Project[]) => {
    projects.value = data
  }

  const setCurrentProject = (project: Project) => {
    currentProject.value = project
  }

  return {
    projects,
    currentProject,
    setProjects,
    setCurrentProject
  }
})
```

---

### 4.3 Phase 3: 项目管理界面 (3-4 天)

#### 3.1 项目列表页面
**views/Project/index.vue**:
```vue
<template>
  <div class="project-list">
    <el-row :gutter="20">
      <el-col :span="4">
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          新建项目
        </el-button>
      </el-col>
      <el-col :span="20">
        <el-input
          v-model="searchText"
          placeholder="搜索项目..."
          clearable
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="project-grid">
      <el-col
        v-for="project in filteredProjects"
        :key="project.id"
        :span="8"
        class="project-item"
      >
        <el-card shadow="hover" @click="handleOpen(project)">
          <template #header>
            <div class="card-header">
              <span>{{ project.name }}</span>
              <el-tag>{{ project.status }}</el-tag>
            </div>
          </template>
          <div class="project-info">
            <p>{{ project.description }}</p>
            <el-descriptions :column="1" size="small" border>
              <el-descriptions-item label="创建时间">
                {{ formatDate(project.created_at) }}
              </el-descriptions-item>
              <el-descriptions-item label="卡片数量">
                {{ project.card_count }}
              </el-descriptions-item>
            </el-descriptions>
          </div>
          <div class="project-actions">
            <el-button link type="primary" @click.stop="handleEdit(project)">
              编辑
            </el-button>
            <el-button link type="danger" @click.stop="handleDelete(project)">
              删除
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <ProjectDialog
      v-model:visible="dialogVisible"
      :project="currentProject"
      @confirm="handleConfirm"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useProjectStore } from '@/stores/project'
import { useRouter } from 'vue-router'
import { formatDate } from '@/utils/format'
import { ElMessageBox } from 'element-plus'
import { Plus, Search } from '@element-plus/icons-vue'
import ProjectDialog from './components/ProjectDialog.vue'

const router = useRouter()
const projectStore = useProjectStore()

const searchText = ref('')
const dialogVisible = ref(false)
const currentProject = ref<any>(null)

const filteredProjects = computed(() => {
  return projectStore.projects.filter(p =>
    p.name.toLowerCase().includes(searchText.value.toLowerCase())
  )
})

const handleCreate = () => {
  currentProject.value = null
  dialogVisible.value = true
}

const handleOpen = (project: any) => {
  projectStore.setCurrentProject(project)
  router.push(`/projects/${project.id}/cards`)
}

const handleEdit = (project: any) => {
  currentProject.value = project
  dialogVisible.value = true
}

const handleDelete = async (project: any) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除项目 "${project.name}" 吗？`,
      '删除确认',
      { type: 'warning' }
    )
    // 调用删除 API
    await deleteProject(project.id)
    ElMessage.success('删除成功')
  } catch {
    // 用户取消
  }
}

const handleConfirm = async (data: any) => {
  if (currentProject.value) {
    await updateProject(currentProject.value.id, data)
    ElMessage.success('更新成功')
  } else {
    await createProject(data)
    ElMessage.success('创建成功')
  }
  dialogVisible.value = false
}
</script>
```

#### 3.2 项目设置页面
**views/Project/Settings.vue**:
```vue
<template>
  <div class="project-settings">
    <el-tabs v-model="activeTab">
      <el-tab-pane label="基本信息" name="basic">
        <ProjectBasicForm :project="project" @save="handleSaveBasic" />
      </el-tab-pane>
      <el-tab-pane label="开发规范" name="standards">
        <StandardsConfig :project="project" @save="handleSaveStandards" />
      </el-tab-pane>
      <el-tab-pane label="分支策略" name="branch">
        <BranchPolicyConfig :project="project" @save="handleSaveBranch" />
      </el-tab-pane>
      <el-tab-pane label="Git 仓库" name="git">
        <GitRepositoryConfig :project="project" @save="handleSaveGit" />
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import { useProjectStore } from '@/stores/project'

const route = useRoute()
const projectStore = useProjectStore()

const activeTab = ref('basic')
const project = ref(projectStore.currentProject)

const handleSaveBasic = async (data: any) => {
  // 保存基本信息
}

const handleSaveStandards = async (data: any) => {
  // 保存开发规范
}

const handleSaveBranch = async (data: any) => {
  // 保存分支策略
}

const handleSaveGit = async (data: any) => {
  // 保存 Git 仓库配置
}
</script>
```

---

### 4.4 Phase 4: 需求卡片管理界面 (4-5 天)

#### 4.1 卡片池视图
**views/Card/Pool.vue**:
```vue
<template>
  <div class="card-pool">
    <el-row :gutter="20" class="filter-bar">
      <el-col :span="6">
        <el-select v-model="filters.status" placeholder="状态筛选" clearable>
          <el-option label="全部" value="" />
          <el-option label="草稿" value="DRAFT" />
          <el-option label="已确认" value="CONFIRMED" />
          <el-option label="进行中" value="IN_PROGRESS" />
          <el-option label="已完成" value="COMPLETED" />
        </el-select>
      </el-col>
      <el-col :span="6">
        <el-select v-model="filters.priority" placeholder="优先级筛选" clearable>
          <el-option label="全部" value="" />
          <el-option label="P0" value="CRITICAL" />
          <el-option label="P1" value="HIGH" />
          <el-option label="P2" value="MEDIUM" />
          <el-option label="P3" value="LOW" />
        </el-select>
      </el-col>
      <el-col :span="8">
        <el-input
          v-model="filters.search"
          placeholder="搜索卡片..."
          clearable
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
      </el-col>
      <el-col :span="4">
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          新建卡片
        </el-button>
      </el-col>
    </el-row>

    <el-table :data="filteredCards" stripe>
      <el-table-column prop="id" label="卡片ID" width="120" />
      <el-table-column prop="title" label="标题" min-width="200" />
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)">
            {{ getStatusLabel(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="priority" label="优先级" width="80" />
      <el-table-column prop="created_at" label="创建时间" width="180" />
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="handleView(row)">
            查看
          </el-button>
          <el-button link type="primary" @click="handleEdit(row)">
            编辑
          </el-button>
          <el-dropdown @command="handleAction($event, row)">
            <el-button link>
              更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="confirm">确认</el-dropdown-item>
                <el-dropdown-item command="start">启动</el-dropdown-item>
                <el-dropdown-item command="complete">完成</el-dropdown-item>
                <el-dropdown-item command="abandon" divided>废弃</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-model:current-page="pagination.page"
      v-model:page-size="pagination.size"
      :total="pagination.total"
      :page-sizes="[10, 20, 50, 100]"
      layout="total, sizes, prev, pager, next, jumper"
      @size-change="handleSizeChange"
      @current-change="handlePageChange"
    />
  </div>
</template>
```

#### 4.2 卡片详情页面
**views/Card/Detail.vue**:
```vue
<template>
  <div class="card-detail">
    <el-page-header @back="handleBack">
      <template #content>
        <div class="header-content">
          <h2>{{ card?.title }}</h2>
          <el-tag :type="getStatusType(card?.status)">
            {{ getStatusLabel(card?.status) }}
          </el-tag>
          <el-tag>优先级: {{ card?.priority }}</el-tag>
        </div>
      </template>
      <template #extra>
        <el-button-group>
          <el-button @click="handleEdit">编辑</el-button>
          <el-button @click="handleStart" :disabled="!canStart">启动</el-button>
          <el-button @click="handleComplete" :disabled="!canComplete">完成</el-button>
        </el-button-group>
      </template>
    </el-page-header>

    <el-row :gutter="20">
      <el-col :span="16">
        <el-card shadow="never">
          <template #header>
            <span>需求描述</span>
          </template>
          <div class="card-content">
            <div class="original-input">
              <h4>原始输入</h4>
              <p>{{ card?.original }}</p>
            </div>
            <div class="structured-output">
              <h4>结构化产出</h4>
              <pre>{{ card?.structured_output }}</pre>
            </div>
          </div>
        </el-card>

        <el-card shadow="never" style="margin-top: 20px">
          <template #header>
            <span>依赖关系</span>
          </template>
          <DependencyGraph :card-id="cardId" />
        </el-card>

        <el-card shadow="never" style="margin-top: 20px">
          <template #header>
            <span>Pipeline 执行记录</span>
          </template>
          <PipelineTimeline :card-id="cardId" />
        </el-card>
      </el-col>

      <el-col :span="8">
        <el-card shadow="never">
          <template #header>
            <span>卡片信息</span>
          </template>
          <el-descriptions :column="1" border>
            <el-descriptions-item label="卡片ID">
              {{ card?.id }}
            </el-descriptions-item>
            <el-descriptions-item label="创建时间">
              {{ formatDate(card?.created_at) }}
            </el-descriptions-item>
            <el-descriptions-item label="更新时间">
              {{ formatDate(card?.updated_at) }}
            </el-descriptions-item>
            <el-descriptions-item label="Pipeline">
              {{ card?.pipeline_id || '未分配' }}
            </el-descriptions-item>
            <el-descriptions-item label="关联任务">
              {{ card?.task_count || 0 }} 个
            </el-descriptions-item>
          </el-descriptions>
        </el-card>

        <el-card shadow="never" style="margin-top: 20px">
          <template #header>
            <span>快速操作</span>
          </template>
          <el-button-group vertical>
            <el-button @click="handleAddDependency">添加依赖</el-button>
            <el-button @click="handleCreateTask">创建任务</el-button>
            <el-button @click="handleViewHistory">查看历史</el-button>
            <el-button @click="handleDuplicate">复制卡片</el-button>
          </el-button-group>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>
```

---

### 4.5 Phase 5: Pipeline 执行监控界面 (3-4 天)

#### 5.1 Pipeline 列表
**views/Pipeline/List.vue**:
```vue
<template>
  <div class="pipeline-list">
    <el-table :data="pipelines" stripe>
      <el-table-column prop="id" label="Pipeline ID" width="200" />
      <el-table-column prop="name" label="名称" min-width="200" />
      <el-table-column prop="current_stage" label="当前阶段" width="120" />
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getPipelineStatusType(row.status)">
            {{ row.status }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="进度" width="200">
        <template #default="{ row }">
          <el-progress :percentage="getProgress(row)" />
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200">
        <template #default="{ row }">
          <el-button link @click="handleView(row)">查看</el-button>
          <el-button link @click="handlePause" :disabled="!canPause(row)">
            暂停
          </el-button>
          <el-button link @click="handleResume" :disabled="!canResume(row)">
            继续
          </el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>
```

#### 5.2 Pipeline 详情和阶段可视化
**views/Pipeline/Detail.vue**:
```vue
<template>
  <div class="pipeline-detail">
    <el-page-header @back="handleBack">
      <template #content>
        <h2>{{ pipeline?.name }}</h2>
        <el-tag>{{ pipeline?.status }}</el-tag>
      </template>
    </el-page-header>

    <el-steps :active="currentStageIndex" finish-status="success" process-status="wait">
      <el-step
        v-for="stage in stages"
        :key="stage.id"
        :title="stage.name"
        :description="stage.description"
      />
    </el-steps>

    <el-row :gutter="20" style="margin-top: 40px">
      <el-col :span="16">
        <el-card shadow="never">
          <template #header>
            <span>阶段执行详情</span>
          </template>
          <StageTimeline :stages="stages" />
        </el-card>
      </el-col>

      <el-col :span="8">
        <el-card shadow="never">
          <template #header>
            <span>执行统计</span>
          </template>
          <el-statistic title="总耗时" :value="totalDuration" suffix="分钟" />
          <el-statistic title="成功阶段" :value="successStages" />
          <el-statistic title="失败阶段" :value="failedStages" />
          <el-statistic title="待审核" :value="pendingReview" />
        </el-card>

        <el-card shadow="never" style="margin-top: 20px">
          <template #header>
            <span>操作日志</span>
          </template>
          <el-timeline>
            <el-timeline-item
              v-for="log in logs"
              :key="log.id"
              :timestamp="log.timestamp"
            >
              {{ log.message }}
            </el-timeline-item>
          </el-timeline>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>
```

---

### 4.6 Phase 6: Agent 执行界面 (3-4 天)

#### 6.1 Agent 任务列表
**views/Agent/Tasks.vue**:
```vue
<template>
  <div class="agent-tasks">
    <el-table :data="tasks" stripe>
      <el-table-column prop="id" label="任务ID" width="200" />
      <el-table-column prop="name" label="任务名称" min-width="200" />
      <el-table-column prop="agent_type" label="Agent 类型" width="120" />
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getTaskStatusType(row.status)">
            {{ row.status }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="priority" label="优先级" width="80" />
      <el-table-column prop="duration" label="耗时" width="100" />
      <el-table-column label="操作" width="200">
        <template #default="{ row }">
          <el-button link @click="handleViewLog(row)">日志</el-button>
          <el-button link @click="handleRetry(row)" :disabled="!canRetry(row)">
            重试
          </el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>
```

#### 6.2 Skill 调用监控
**views/Agent/Monitor.vue**:
```vue
<template>
  <div class="agent-monitor">
    <el-row :gutter="20">
      <el-col :span="12">
        <el-card shadow="never">
          <template #header>
            <span>实时任务</span>
          </template>
          <RealTimeTaskMonitor />
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="never">
          <template #header>
            <span>Skill 调用统计</span>
          </template>
          <SkillStatistics />
        </el-card>
      </el-col>
    </el-row>

    <el-card shadow="never" style="margin-top: 20px">
      <template #header>
        <span>Agent 执行历史</span>
      </template>
      <ExecutionTimeline />
    </el-card>
  </div>
</template>
```

---

### 4.7 Phase 7: 知识库管理界面 (3-4 天)

#### 7.1 文档管理
**views/Knowledge/Documents.vue**:
```vue
<template>
  <div class="knowledge-documents">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card shadow="never">
          <template #header>
            <span>文档目录</span>
          </template>
          <DocumentTree @select="handleSelectDocument" />
        </el-card>
      </el-col>
      <el-col :span="18">
        <el-card shadow="never">
          <template #header>
            <span>文档列表</span>
            <el-button style="float: right" type="primary" @click="handleUpload">
              上传文档
            </el-button>
          </template>
          <DocumentTable :documents="documents" @edit="handleEditDocument" />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>
```

#### 7.2 索引状态监控
**views/Knowledge/IndexStatus.vue**:
```vue
<template>
  <div class="index-status">
    <el-row :gutter="20">
      <el-col :span="8">
        <el-card shadow="never">
          <template #header>
            <span>索引概览</span>
          </template>
          <el-statistic title="总文档数" :value="totalDocs" />
          <el-statistic title="已索引" :value="indexedDocs" />
          <el-statistic title="待索引" :value="pendingDocs" />
          <el-statistic title="索引错误" :value="errorDocs" />
        </el-card>
      </el-col>
      <el-col :span="16">
        <el-card shadow="never">
          <template #header>
            <span>索引状态</span>
            <el-button style="float: right" @click="handleReindexAll">
              全部重建索引
            </el-button>
          </template>
          <IndexStatusTable :documents="documents" @reindex="handleReindex" />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>
```

---

### 4.8 Phase 8: AI Runtime 管理界面 (3-4 天)

#### 8.1 Provider 管理
**views/AI/Providers.vue**:
```vue
<template>
  <div class="ai-providers">
    <el-button type="primary" @click="handleAddProvider">
      添加 Provider
    </el-button>

    <el-table :data="providers" stripe style="margin-top: 20px">
      <el-table-column prop="name" label="名称" width="200" />
      <el-table-column prop="type" label="类型" width="120" />
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 'available' ? 'success' : 'danger'">
            {{ row.status }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="model_count" label="模型数量" width="100" />
      <el-table-column prop="last_used" label="最后使用" width="180" />
      <el-table-column label="操作" width="200">
        <template #default="{ row }">
          <el-button link @click="handleEdit(row)">编辑</el-button>
          <el-button link @click="handleTest(row)">测试</el-button>
          <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <ProviderDialog
      v-model:visible="dialogVisible"
      :provider="currentProvider"
      @confirm="handleConfirm"
    />
  </div>
</template>
```

#### 8.2 调用统计仪表板
**views/AI/Statistics.vue**:
```vue
<template>
  <div class="ai-statistics">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card shadow="never">
          <el-statistic title="总调用次数" :value="totalCalls" />
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="never">
          <el-statistic title="总 Token 使用" :value="totalTokens" />
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="never">
          <el-statistic title="总费用" :value="totalCost" :precision="2" suffix="元" />
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="never">
          <el-statistic title="今日调用" :value="todayCalls" />
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="16">
        <el-card shadow="never">
          <template #header>
            <span>调用趋势</span>
          </template>
          <div ref="chartRef" style="width: 100%; height: 400px;"></div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="never">
          <template #header>
            <span>模型使用分布</span>
          </template>
          <div ref="pieChartRef" style="width: 100%; height: 400px;"></div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>
```

---

## 5. 时间安排

| 阶段 | 工作内容 | 预计时间 |
|------|----------|----------|
| Phase 1 | 项目初始化 | 2-3 天 |
| Phase 2 | 基础设施开发 | 3-4 天 |
| Phase 3 | 项目管理界面 | 3-4 天 |
| Phase 4 | 需求卡片管理界面 | 4-5 天 |
| Phase 5 | Pipeline 执行监控界面 | 3-4 天 |
| Phase 6 | Agent 执行界面 | 3-4 天 |
| Phase 7 | 知识库管理界面 | 3-4 天 |
| Phase 8 | AI Runtime 管理界面 | 3-4 天 |

**总计**: 2-4 周

---

## 6. 验收标准

### 6.1 功能验收
- [ ] 所有 8 个主要模块界面正常工作
- [ ] 用户可以创建、编辑、删除项目
- [ ] 用户可以管理需求卡片
- [ ] 用户可以监控 Pipeline 执行
- [ ] 用户可以查看 Agent 执行状态
- [ ] 用户可以管理知识库
- [ ] 用户可以配置 AI Runtime

### 6.2 性能验收
- [ ] 页面加载时间 < 2 秒
- [ ] 列表查询响应时间 < 500ms
- [ ] 表单提交响应时间 < 1 秒
- [ ] 实时更新延迟 < 1 秒

### 6.3 质量验收
- [ ] 代码符合 ESLint 规范
- [ ] TypeScript 类型检查通过
- [ ] 组件可复用性良好
- [ ] 响应式设计适配主流屏幕

---

## 7. 注意事项

1. **用户体验**: 注重用户体验，提供友好的错误提示
2. **响应式设计**: 适配不同屏幕尺寸
3. **可访问性**: 支持键盘导航和屏幕阅读器
4. **国际化**: 预留国际化支持
5. **安全性**: 防止 XSS 攻击，验证用户输入
