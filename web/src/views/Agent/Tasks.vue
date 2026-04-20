<template>
  <div class="agent-tasks">
    <el-page-header @back="goBack">
      <template #content>
        <span class="title">Agent 任务列表</span>
      </template>
    </el-page-header>

    <el-card class="filter-card">
      <el-form :inline="true" :model="filters">
        <el-form-item label="任务类型">
          <el-select
            v-model="filters.type"
            placeholder="全部类型"
            clearable
            @change="loadTasks"
          >
            <el-option
              v-for="type in agentTypes"
              :key="type.value"
              :label="type.label"
              :value="type.value"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="状态">
          <el-select
            v-model="filters.status"
            placeholder="全部状态"
            clearable
            @change="loadTasks"
          >
            <el-option
              v-for="status in taskStatuses"
              :key="status.value"
              :label="status.label"
              :value="status.value"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="搜索">
          <el-input
            v-model="filters.keyword"
            placeholder="搜索任务"
            clearable
            @input="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="table-card" v-loading="loading">
      <el-table :data="tasks" style="width: 100%">
        <el-table-column prop="type" label="类型" width="150">
          <template #default="{ row }">
            {{ getTypeText(row.type) }}
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="Skill" width="120">
          <template #default="{ row }">
            {{ row.skill || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="耗时" width="120">
          <template #default="{ row }">
            {{ row.duration ? formatDuration(row.duration) : '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button
              link
              type="primary"
              size="small"
              @click="handleViewLogs(row)"
            >
              查看日志
            </el-button>
            <el-button
              link
              type="primary"
              size="small"
              :disabled="!canRetry(row)"
              @click="handleRetry(row)"
            >
              重试
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.size"
        :total="pagination.total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        style="margin-top: 20px; justify-content: flex-end"
        @size-change="loadTasks"
        @current-change="loadTasks"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import { formatDate, formatDuration } from '@/utils/format'
import type { AgentTask, AgentType, AgentTaskStatus } from '@/types/agent'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const tasks = ref<AgentTask[]>([])

const filters = ref({
  type: '',
  status: '',
  keyword: ''
})

const pagination = ref({
  page: 1,
  size: 20,
  total: 0
})

const agentTypes = [
  { label: '代码生成', value: 'CODE_GENERATION' },
  { label: '测试', value: 'TEST' },
  { label: '文档', value: 'DOCUMENTATION' },
  { label: '审查', value: 'REVIEW' },
  { label: '重构', value: 'REFACTORING' },
  { label: '自定义', value: 'CUSTOM' }
]

const taskStatuses = [
  { label: '等待中', value: 'PENDING' },
  { label: '已分配', value: 'ASSIGNED' },
  { label: '运行中', value: 'RUNNING' },
  { label: '已完成', value: 'COMPLETED' },
  { label: '失败', value: 'FAILED' },
  { label: '已取消', value: 'CANCELLED' }
]

const projectId = route.params.id as string

const loadTasks = async () => {
  loading.value = true
  try {
    const data = await fetch(`/api/v1/projects/${projectId}/agent-tasks?page=${pagination.value.page}&size=${pagination.value.size}`)
    const result = await data.json()
    tasks.value = result.items || []
    pagination.value.total = result.total || 0

    if (filters.value.type) {
      tasks.value = tasks.value.filter(t => t.type === filters.value.type)
    }
    if (filters.value.status) {
      tasks.value = tasks.value.filter(t => t.status === filters.value.status)
    }
  } catch (error) {
    ElMessage.error('加载任务列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.value.page = 1
  loadTasks()
}

const handleViewLogs = (task: AgentTask) => {
  router.push(`/projects/${projectId}/agent-tasks/${task.id}/logs`)
}

const handleRetry = (task: AgentTask) => {
  ElMessage.info('重试功能开发中')
}

const canRetry = (task: AgentTask) => {
  return task.status === 'FAILED' || task.status === 'CANCELLED'
}

const getTypeText = (type: string) => {
  const texts: Record<string, string> = {
    CODE_GENERATION: '代码生成',
    TEST: '测试',
    DOCUMENTATION: '文档',
    REVIEW: '审查',
    REFACTORING: '重构',
    CUSTOM: '自定义'
  }
  return texts[type] || type
}

const getStatusType = (status: string) => {
  const types: Record<string, any> = {
    PENDING: 'info',
    ASSIGNED: 'warning',
    RUNNING: 'primary',
    COMPLETED: 'success',
    FAILED: 'danger',
    CANCELLED: 'info'
  }
  return types[status] || 'info'
}

const getStatusText = (status: string) => {
  const texts: Record<string, string> = {
    PENDING: '等待中',
    ASSIGNED: '已分配',
    RUNNING: '运行中',
    COMPLETED: '已完成',
    FAILED: '失败',
    CANCELLED: '已取消'
  }
  return texts[status] || status
}

const goBack = () => {
  router.push('/projects')
}

onMounted(() => {
  loadTasks()
})
</script>

<style scoped lang="scss">
.agent-tasks {
  .title {
    font-size: 20px;
    font-weight: 500;
  }

  .filter-card {
    margin: 20px 0;
  }

  .table-card {
    min-height: 400px;
  }
}
</style>
