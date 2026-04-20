<template>
  <div class="pipeline-detail" v-loading="loading">
    <el-page-header @back="goBack">
      <template #content>
        <div class="page-header-content">
          <span class="title">{{ pipeline?.name }}</span>
          <el-tag :type="getStatusType(pipeline?.status)" size="small">
            {{ getStatusText(pipeline?.status) }}
          </el-tag>
        </div>
      </template>
    </el-page-header>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="16">
        <el-card class="stages-card">
          <template #header>
            <span>执行阶段</span>
          </template>
          <el-steps :active="activeStage" finish-status="success" process-status="process">
            <el-step
              v-for="(stage, index) in pipeline?.stages"
              :key="stage.id"
              :title="stage.name"
              :description="getStageDescription(stage)"
              :status="getStageStatus(stage.status)"
              @click="handleStageClick(stage)"
            />
          </el-steps>
        </el-card>

        <el-card class="timeline-card" style="margin-top: 20px">
          <template #header>
            <span>执行时间线</span>
          </template>
          <el-timeline>
            <el-timeline-item
              v-for="(stage, index) in pipeline?.stages"
              :key="stage.id"
              :timestamp="formatDateTime(stage.startTime, stage.endTime)"
              :type="getTimelineType(stage.status)"
            >
              <div class="timeline-content">
                <div class="stage-name">{{ stage.name }}</div>
                <div class="stage-status">{{ getStageStatusText(stage.status) }}</div>
                <div v-if="stage.steps.length > 0" class="steps-list">
                  <div
                    v-for="step in stage.steps"
                    :key="step.id"
                    class="step-item"
                  >
                    <el-icon class="step-icon">
                      <component :is="getStepIcon(step.status)" />
                    </el-icon>
                    <span class="step-name">{{ step.name }}</span>
                    <span class="step-time">
                      {{ formatDuration(step.startTime, step.endTime) }}
                    </span>
                  </div>
                </div>
              </div>
            </el-timeline-item>
          </el-timeline>
        </el-card>
      </el-col>

      <el-col :span="8">
        <el-card class="info-card">
          <template #header>
            <span>基本信息</span>
          </template>
          <el-descriptions :column="1" border>
            <el-descriptions-item label="Pipeline ID">
              {{ pipeline?.id }}
            </el-descriptions-item>
            <el-descriptions-item label="关联卡片">
              {{ pipeline?.cardId }}
            </el-descriptions-item>
            <el-descriptions-item label="创建时间">
              {{ formatDate(pipeline?.createdAt) }}
            </el-descriptions-item>
            <el-descriptions-item label="总阶段数">
              {{ pipeline?.stages.length || 0 }}
            </el-descriptions-item>
            <el-descriptions-item label="已完成阶段">
              {{ completedStages }} / {{ pipeline?.stages.length || 0 }}
            </el-descriptions-item>
          </el-descriptions>
        </el-card>

        <el-card class="stats-card" style="margin-top: 20px">
          <template #header>
            <span>执行统计</span>
          </template>
          <el-statistic title="总耗时" :value="totalDuration" suffix="s" />
          <el-statistic title="总步骤数" :value="totalSteps" style="margin-top: 20px" />
          <el-statistic
            title="失败步骤"
            :value="failedSteps"
            style="margin-top: 20px"
          />
        </el-card>

        <el-card class="actions-card" style="margin-top: 20px">
          <template #header>
            <span>操作</span>
          </template>
          <el-space wrap>
            <el-button
              type="warning"
              size="small"
              :disabled="!canPause"
              @click="handlePause"
            >
              暂停
            </el-button>
            <el-button
              type="success"
              size="small"
              :disabled="!canResume"
              @click="handleResume"
            >
              继续
            </el-button>
            <el-button
              type="danger"
              size="small"
              :disabled="!canCancel"
              @click="handleCancel"
            >
              取消
            </el-button>
            <el-button
              type="primary"
              size="small"
              @click="viewLogs"
            >
              查看日志
            </el-button>
          </el-space>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { CircleCheck, CircleClose, Clock, Loading } from '@element-plus/icons-vue'
import { usePipelineStore } from '@/stores/pipeline'
import { formatDate, formatDuration } from '@/utils/format'
import type { Pipeline, PipelineStage, StageStatus } from '@/types/pipeline'

const route = useRoute()
const router = useRouter()
const pipelineStore = usePipelineStore()

const loading = ref(false)
const pipeline = ref<Pipeline | null>(null)

const pipelineId = route.params.id as string

const activeStage = computed(() => {
  if (!pipeline.value) return -1
  return pipeline.value.stages.findIndex(s => s.status === 'RUNNING')
})

const completedStages = computed(() => {
  if (!pipeline.value) return 0
  return pipeline.value.stages.filter(s => s.status === 'COMPLETED').length
})

const totalDuration = computed(() => {
  if (!pipeline.value) return 0
  let total = 0
  pipeline.value.stages.forEach(stage => {
    if (stage.startTime && stage.endTime) {
      total += new Date(stage.endTime).getTime() - new Date(stage.startTime).getTime()
    }
  })
  return Math.round(total / 1000)
})

const totalSteps = computed(() => {
  if (!pipeline.value) return 0
  return pipeline.value.stages.reduce((sum, stage) => sum + stage.steps.length, 0)
})

const failedSteps = computed(() => {
  if (!pipeline.value) return 0
  return pipeline.value.stages.reduce((sum, stage) => {
    return sum + stage.steps.filter(s => s.status === 'FAILED').length
  }, 0)
})

const canPause = computed(() => pipeline.value?.status === 'RUNNING')
const canResume = computed(() => pipeline.value?.status === 'PAUSED')
const canCancel = computed(() =>
  pipeline.value && ['RUNNING', 'PAUSED', 'PENDING'].includes(pipeline.value.status)
)

const loadPipeline = async () => {
  loading.value = true
  try {
    pipeline.value = await pipelineStore.getPipeline(pipelineId)
  } catch (error) {
    ElMessage.error('加载 Pipeline 详情失败')
  } finally {
    loading.value = false
  }
}

const handleStageClick = (stage: PipelineStage) => {
}

const handlePause = async () => {
  try {
    await pipelineStore.pausePipeline(pipelineId)
    ElMessage.success('已暂停')
    loadPipeline()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleResume = async () => {
  try {
    await pipelineStore.resumePipeline(pipelineId)
    ElMessage.success('已继续')
    loadPipeline()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleCancel = async () => {
  try {
    await ElMessageBox.confirm('确定要取消此 Pipeline 吗？', '确认', {
      type: 'warning'
    })
    await pipelineStore.cancelPipeline(pipelineId)
    ElMessage.success('已取消')
    loadPipeline()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('操作失败')
    }
  }
}

const viewLogs = () => {
  ElMessage.info('日志查看功能开发中')
}

const getStatusType = (status?: string) => {
  if (!status) return ''
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

const getStatusText = (status?: string) => {
  if (!status) return ''
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

const getStageStatus = (status: StageStatus) => {
  const statusMap: Record<string, any> = {
    COMPLETED: 'success',
    FAILED: 'error',
    RUNNING: 'process',
    PENDING: 'wait',
    REVIEW: 'warning',
    SKIPPED: 'info'
  }
  return statusMap[status] || 'wait'
}

const getStageDescription = (stage: PipelineStage) => {
  if (stage.status === 'COMPLETED') {
    return `已完成 - ${formatDuration(stage.startTime, stage.endTime)}`
  }
  if (stage.status === 'FAILED') {
    return '失败'
  }
  if (stage.status === 'RUNNING') {
    return '执行中...'
  }
  if (stage.status === 'REVIEW') {
    return '待审核'
  }
  return '等待中'
}

const getStageStatusText = (status: StageStatus) => {
  const texts: Record<string, string> = {
    PENDING: '等待中',
    RUNNING: '执行中',
    REVIEW: '待审核',
    COMPLETED: '已完成',
    FAILED: '失败',
    SKIPPED: '已跳过'
  }
  return texts[status] || status
}

const getTimelineType = (status: StageStatus) => {
  const types: Record<string, any> = {
    COMPLETED: 'primary',
    FAILED: 'danger',
    RUNNING: 'warning',
    PENDING: 'info'
  }
  return types[status] || 'info'
}

const getStepIcon = (status: StageStatus) => {
  const icons: Record<string, any> = {
    COMPLETED: CircleCheck,
    FAILED: CircleClose,
    RUNNING: Loading,
    PENDING: Clock
  }
  return icons[status] || Clock
}

const formatDateTime = (start?: string, end?: string) => {
  if (!start) return ''
  if (end) {
    return `${formatDate(start)} - ${formatDate(end)}`
  }
  return `${formatDate(start)} - 进行中`
}

const goBack = () => {
  router.push('/projects')
}

onMounted(() => {
  loadPipeline()
})
</script>

<style scoped lang="scss">
.pipeline-detail {
  .page-header-content {
    display: flex;
    justify-content: space-between;
    align-items: center;
    width: 100%;

    .title {
      font-size: 20px;
      font-weight: 500;
    }
  }

  .timeline-content {
    .stage-name {
      font-weight: 500;
      margin-bottom: 4px;
    }

    .stage-status {
      font-size: 12px;
      color: #666;
      margin-bottom: 8px;
    }

    .steps-list {
      margin-top: 12px;
      padding-left: 12px;

      .step-item {
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: 13px;
        margin-bottom: 6px;

        .step-icon {
          font-size: 14px;
        }

        .step-name {
          flex: 1;
        }

        .step-time {
          color: #999;
          font-size: 12px;
        }
      }
    }
  }

  .el-steps {
    cursor: pointer;

    :deep(.el-step__title) {
      cursor: pointer;
    }
  }
}
</style>
