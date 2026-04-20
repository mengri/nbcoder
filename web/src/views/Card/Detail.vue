<template>
  <div
    v-loading="loading"
    class="card-detail"
  >
    <el-page-header @back="goBack">
      <template #content>
        <div class="page-header-content">
          <span class="title">{{ card?.title }}</span>
          <div class="actions">
            <el-tag
              :type="getStatusType(card?.status)"
              size="small"
            >
              {{ getStatusText(card?.status) }}
            </el-tag>
            <el-tag
              :type="getPriorityType(card?.priority)"
              size="small"
            >
              {{ getPriorityText(card?.priority) }}
            </el-tag>
          </div>
        </div>
      </template>
    </el-page-header>

    <el-row
      :gutter="20"
      style="margin-top: 20px"
    >
      <el-col :span="16">
        <el-card class="info-card">
          <template #header>
            <span>基本信息</span>
          </template>
          <el-descriptions
            :column="2"
            border
          >
            <el-descriptions-item label="卡片ID">
              {{ card?.id }}
            </el-descriptions-item>
            <el-descriptions-item label="创建时间">
              {{ formatDate(card?.createdAt) }}
            </el-descriptions-item>
            <el-descriptions-item label="更新时间">
              {{ formatDate(card?.updatedAt) }}
            </el-descriptions-item>
            <el-descriptions-item label="Pipeline">
              {{ card?.pipelineId || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="Agent 任务">
              {{ card?.agentTaskId || '-' }}
            </el-descriptions-item>
          </el-descriptions>
        </el-card>

        <el-card
          class="description-card"
          style="margin-top: 20px"
        >
          <template #header>
            <span>描述</span>
          </template>
          <p>{{ card?.description }}</p>
        </el-card>

        <el-card
          class="raw-input-card"
          style="margin-top: 20px"
        >
          <template #header>
            <span>原始输入</span>
          </template>
          <pre class="code-block">{{ card?.rawInput }}</pre>
        </el-card>

        <el-card
          v-if="card?.structuredOutput"
          class="output-card"
          style="margin-top: 20px"
        >
          <template #header>
            <span>结构化产出</span>
          </template>
          <pre class="code-block">{{ formatJson(card.structuredOutput) }}</pre>
        </el-card>
      </el-col>

      <el-col :span="8">
        <el-card class="dependencies-card">
          <template #header>
            <div class="card-header">
              <span>依赖关系</span>
              <el-button
                size="small"
                @click="handleAddDependency"
              >
                添加依赖
              </el-button>
            </div>
          </template>
          <div
            v-if="card?.dependencies.length === 0"
            class="empty"
          >
            暂无依赖
          </div>
          <el-tag
            v-for="depId in card?.dependencies"
            :key="depId"
            class="dependency-tag"
            closable
            @close="handleRemoveDependency(depId)"
          >
            {{ depId }}
          </el-tag>
        </el-card>

        <el-card
          class="operations-card"
          style="margin-top: 20px"
        >
          <template #header>
            <span>快速操作</span>
          </template>
          <el-space wrap>
            <el-button
              type="primary"
              size="small"
              :disabled="!canConfirm"
              @click="handleConfirm"
            >
              确认卡片
            </el-button>
            <el-button
              type="success"
              size="small"
              :disabled="!canStart"
              @click="handleStart"
            >
              启动开发
            </el-button>
            <el-button
              type="success"
              size="small"
              :disabled="!canComplete"
              @click="handleComplete"
            >
              完成卡片
            </el-button>
            <el-button
              type="danger"
              size="small"
              :disabled="!canAbandon"
              @click="handleAbandon"
            >
              废弃卡片
            </el-button>
          </el-space>
        </el-card>

        <el-card
          class="pipeline-card"
          style="margin-top: 20px"
        >
          <template #header>
            <span>Pipeline 执行</span>
          </template>
          <div v-if="card?.pipelineId">
            <el-button
              type="primary"
              size="small"
              link
              @click="viewPipeline"
            >
              查看 Pipeline
            </el-button>
          </div>
          <div
            v-else
            class="empty"
          >
            尚未启动 Pipeline
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useCardStore } from '@/stores/card'
import { formatDate, formatJson } from '@/utils/format'
import type { Card } from '@/types/card'

const route = useRoute()
const router = useRouter()
const cardStore = useCardStore()

const loading = ref(false)
const card = ref<Card | null>(null)

const projectId = route.params.id as string
const cardId = route.params.cardId as string

const canConfirm = computed(() => {
  return card.value?.status === 'PENDING' || card.value?.status === 'DRAFT'
})

const canStart = computed(() => {
  return card.value?.status === 'CONFIRMED'
})

const canComplete = computed(() => {
  return card.value?.status === 'IN_PROGRESS'
})

const canAbandon = computed(() => {
  return card.value && !['COMPLETED', 'ABANDONED'].includes(card.value.status)
})

const loadCard = async () => {
  loading.value = true
  try {
    card.value = await cardStore.getCard(cardId)
  } catch (error) {
    ElMessage.error('加载卡片详情失败')
  } finally {
    loading.value = false
  }
}

const handleAddDependency = () => {
  ElMessage.info('添加依赖功能开发中')
}

const handleRemoveDependency = async (depId: string) => {
  if (!card.value) return

  try {
    await ElMessageBox.confirm('确定要移除此依赖吗？', '确认', {
      type: 'warning'
    })

    const dependencies = card.value.dependencies.filter(id => id !== depId)
    await cardStore.updateCard(cardId, { dependencies })
    ElMessage.success('已移除依赖')
    loadCard()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('操作失败')
    }
  }
}

const handleConfirm = async () => {
  try {
    await cardStore.confirmCard(cardId)
    ElMessage.success('已确认')
    loadCard()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleStart = async () => {
  try {
    await cardStore.startCard(cardId)
    ElMessage.success('已启动')
    loadCard()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleComplete = async () => {
  try {
    await cardStore.completeCard(cardId)
    ElMessage.success('已完成')
    loadCard()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleAbandon = async () => {
  try {
    await ElMessageBox.confirm('确定要废弃此卡片吗？', '确认', {
      type: 'warning'
    })
    await cardStore.abandonCard(cardId)
    ElMessage.success('已废弃')
    loadCard()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('操作失败')
    }
  }
}

const viewPipeline = () => {
  if (card.value?.pipelineId) {
    router.push(`/pipelines/${card.value.pipelineId}`)
  }
}

const getStatusType = (status?: string) => {
  if (!status) return ''
  const types: Record<string, any> = {
    DRAFT: 'info',
    PENDING: 'warning',
    CONFIRMED: 'success',
    IN_PROGRESS: 'primary',
    COMPLETED: 'success',
    SUPPLANTED: 'info',
    ABANDONED: 'danger'
  }
  return types[status] || 'info'
}

const getStatusText = (status?: string) => {
  if (!status) return ''
  const texts: Record<string, string> = {
    DRAFT: '草稿',
    PENDING: '待确认',
    CONFIRMED: '已确认',
    IN_PROGRESS: '进行中',
    COMPLETED: '已完成',
    SUPPLANTED: '被取代',
    ABANDONED: '已废弃'
  }
  return texts[status] || status
}

const getPriorityType = (priority?: string) => {
  if (!priority) return ''
  const types: Record<string, any> = {
    LOW: 'info',
    MEDIUM: '',
    HIGH: 'warning',
    CRITICAL: 'danger'
  }
  return types[priority] || ''
}

const getPriorityText = (priority?: string) => {
  if (!priority) return ''
  const texts: Record<string, string> = {
    LOW: '低',
    MEDIUM: '中',
    HIGH: '高',
    CRITICAL: '紧急'
  }
  return texts[priority] || priority
}

const goBack = () => {
  router.push(`/projects/${projectId}/cards`)
}

onMounted(() => {
  loadCard()
})
</script>

<style scoped lang="scss">
.card-detail {
  .page-header-content {
    display: flex;
    justify-content: space-between;
    align-items: center;
    width: 100%;

    .title {
      font-size: 20px;
      font-weight: 500;
    }

    .actions {
      display: flex;
      gap: 8px;
    }
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .code-block {
    background-color: #f5f7fa;
    padding: 12px;
    border-radius: 4px;
    font-size: 13px;
    line-height: 1.6;
    max-height: 400px;
    overflow-y: auto;
  }

  .dependency-tag {
    margin: 4px;
  }

  .empty {
    color: #999;
    text-align: center;
    padding: 20px 0;
  }
}
</style>
