<template>
  <div class="pipeline-list">
    <el-page-header @back="goBack">
      <template #content>
        <span class="title">Pipeline 列表</span>
      </template>
    </el-page-header>

    <el-card class="table-card" v-loading="loading">
      <el-table :data="pipelines" style="width: 100%">
        <el-table-column prop="name" label="名称" min-width="200" />
        <el-table-column prop="cardId" label="关联卡片" width="150" />
        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="当前阶段" width="150">
          <template #default="{ row }">
            {{ row.currentStage || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="进度" width="200">
          <template #default="{ row }">
            <el-progress
              :percentage="getProgress(row)"
              :status="getProgressStatus(row.status)"
              :stroke-width="8"
            />
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
              @click="handleView(row)"
            >
              查看
            </el-button>
            <el-dropdown @command="(cmd) => handleCommand(cmd, row)">
              <el-button link type="primary" size="small">
                操作<el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item
                    command="pause"
                    :disabled="!canPause(row)"
                  >
                    暂停
                  </el-dropdown-item>
                  <el-dropdown-item
                    command="resume"
                    :disabled="!canResume(row)"
                  >
                    继续
                  </el-dropdown-item>
                  <el-dropdown-item
                    command="cancel"
                    :disabled="!canCancel(row)"
                    divided
                  >
                    取消
                  </el-dropdown-item>
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
        style="margin-top: 20px; justify-content: flex-end"
        @size-change="loadPipelines"
        @current-change="loadPipelines"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowDown } from '@element-plus/icons-vue'
import { usePipelineStore } from '@/stores/pipeline'
import { useProjectStore } from '@/stores/project'
import { formatDate } from '@/utils/format'
import type { Pipeline, PipelineStatus } from '@/types/pipeline'

const route = useRoute()
const router = useRouter()
const pipelineStore = usePipelineStore()
const projectStore = useProjectStore()

const loading = ref(false)
const pipelines = ref<Pipeline[]>([])

const pagination = ref({
  page: 1,
  size: 20,
  total: 0
})

const projectId = route.params.id as string

const loadPipelines = async () => {
  loading.value = true
  try {
    const result = await pipelineStore.loadPipelines(projectId, {
      page: pagination.value.page,
      size: pagination.value.size
    })
    pipelines.value = result.items
    pagination.value.total = result.total
  } catch (error) {
    ElMessage.error('加载 Pipeline 列表失败')
  } finally {
    loading.value = false
  }
}

const handleView = (pipeline: Pipeline) => {
  router.push(`/pipelines/${pipeline.id}`)
}

const handleCommand = async (command: string, pipeline: Pipeline) => {
  switch (command) {
    case 'pause':
      await handlePause(pipeline)
      break
    case 'resume':
      await handleResume(pipeline)
      break
    case 'cancel':
      await handleCancel(pipeline)
      break
  }
}

const handlePause = async (pipeline: Pipeline) => {
  try {
    await pipelineStore.pausePipeline(pipeline.id)
    ElMessage.success('已暂停')
    loadPipelines()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleResume = async (pipeline: Pipeline) => {
  try {
    await pipelineStore.resumePipeline(pipeline.id)
    ElMessage.success('已继续')
    loadPipelines()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleCancel = async (pipeline: Pipeline) => {
  try {
    await ElMessageBox.confirm('确定要取消此 Pipeline 吗？', '确认', {
      type: 'warning'
    })
    await pipelineStore.cancelPipeline(pipeline.id)
    ElMessage.success('已取消')
    loadPipelines()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('操作失败')
    }
  }
}

const canPause = (pipeline: Pipeline) => pipeline.status === 'RUNNING'
const canResume = (pipeline: Pipeline) => pipeline.status === 'PAUSED'
const canCancel = (pipeline: Pipeline) =>
  ['RUNNING', 'PAUSED', 'PENDING'].includes(pipeline.status)

const getProgress = (pipeline: Pipeline) => {
  const total = pipeline.stages.length
  if (total === 0) return 0

  const completed = pipeline.stages.filter(
    s => s.status === 'COMPLETED'
  ).length
  return Math.round((completed / total) * 100)
}

const getProgressStatus = (status: PipelineStatus) => {
  const statusMap: Record<string, any> = {
    COMPLETED: 'success',
    FAILED: 'exception',
    CANCELLED: 'warning'
  }
  return statusMap[status] || ''
}

const getStatusType = (status: string) => {
  const types: Record<string, any> = {
    PENDING: 'info',
    RUNNING: 'primary',
    PAUSED: 'warning',
    COMPLETED: 'success',
    FAILED: 'danger',
    CANCELLED: 'info'
  }
  return types[status] || 'info'
}

const getStatusText = (status: string) => {
  const texts: Record<string, string> = {
    PENDING: '等待中',
    RUNNING: '运行中',
    PAUSED: '已暂停',
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
  const project = projectStore.projects.find(p => p.id === projectId)
  if (project) {
    projectStore.setCurrentProject(project)
  }
  loadPipelines()
})
</script>

<style scoped lang="scss">
.pipeline-list {
  .title {
    font-size: 20px;
    font-weight: 500;
  }

  .table-card {
    margin-top: 20px;
    min-height: 400px;
  }
}
</style>
