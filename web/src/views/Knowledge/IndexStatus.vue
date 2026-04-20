<template>
  <div class="index-status">
    <el-page-header @back="goBack">
      <template #content>
        <span class="title">索引状态监控</span>
      </template>
    </el-page-header>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="6">
        <el-card>
          <el-statistic title="总文档数" :value="indexStatus.totalDocuments">
            <template #prefix>
              <el-icon style="vertical-align: -0.125em">
                <Document />
              </el-icon>
            </template>
          </el-statistic>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <el-statistic title="已索引文档" :value="indexStatus.indexedDocuments">
            <template #prefix>
              <el-icon style="vertical-align: -0.125em">
                <CircleCheck />
              </el-icon>
            </template>
          </el-statistic>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <el-statistic title="总分片数" :value="indexStatus.totalChunks">
            <template #prefix>
              <el-icon style="vertical-align: -0.125em">
                <Grid />
              </el-icon>
            </template>
          </el-statistic>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <el-statistic title="已索引分片" :value="indexStatus.indexedChunks">
            <template #prefix>
              <el-icon style="vertical-align: -0.125em">
                <Finished />
              </el-icon>
            </template>
          </el-statistic>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>索引进度</span>
              <el-tag :type="getStatusType(indexStatus.status)" size="small">
                {{ getStatusText(indexStatus.status) }}
              </el-tag>
            </div>
          </template>
          <div class="progress-section">
            <div class="progress-item">
              <div class="progress-label">文档索引进度</div>
              <el-progress
                :percentage="documentProgress"
                :status="getProgressStatus(indexStatus.status)"
              />
            </div>
            <div class="progress-item" style="margin-top: 24px">
              <div class="progress-label">分片索引进度</div>
              <el-progress
                :percentage="chunkProgress"
                :status="getProgressStatus(indexStatus.status)"
              />
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>索引操作</span>
            </div>
          </template>
          <el-space direction="vertical" style="width: 100%">
            <el-button
              type="primary"
              :loading="rebuilding"
              @click="handleRebuildAll"
            >
              <el-icon><Refresh /></el-icon>
              重建所有索引
            </el-button>
            <el-button
              type="success"
              :loading="optimizing"
              @click="handleOptimize"
            >
              <el-icon><Operation /></el-icon>
              优化索引
            </el-button>
            <el-button
              type="warning"
              @click="handleExportStats"
            >
              <el-icon><Download /></el-icon>
              导出统计信息
            </el-button>
          </el-space>
        </el-card>
      </el-col>
    </el-row>

    <el-card style="margin-top: 20px">
      <template #header>
        <span>文档索引状态</span>
      </template>
      <IndexStatusTable :documents="documents" :loading="loading" />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Document,
  CircleCheck,
  Grid,
  Finished,
  Refresh,
  Operation,
  Download
} from '@element-plus/icons-vue'
import IndexStatusTable from '@/components/business/IndexStatusTable.vue'
import type { Document, IndexStatus as IndexStatusType } from '@/types/knowledge'

const route = useRoute()
const router = useRouter()

const projectId = route.params.id as string

const loading = ref(false)
const rebuilding = ref(false)
const optimizing = ref(false)

const indexStatus = ref<IndexStatusType>({
  id: '',
  projectId,
  totalDocuments: 0,
  indexedDocuments: 0,
  totalChunks: 0,
  indexedChunks: 0,
  status: 'IDLE'
})

const documents = ref<Document[]>([])

const documentProgress = computed(() => {
  if (indexStatus.value.totalDocuments === 0) return 0
  return Math.round((indexStatus.value.indexedDocuments / indexStatus.value.totalDocuments) * 100)
})

const chunkProgress = computed(() => {
  if (indexStatus.value.totalChunks === 0) return 0
  return Math.round((indexStatus.value.indexedChunks / indexStatus.value.totalChunks) * 100)
})

const loadIndexStatus = async () => {
  loading.value = true
  try {
    const [statusData, docsData] = await Promise.all([
      fetch(`/api/v1/projects/${projectId}/knowledge/index-status`),
      fetch(`/api/v1/projects/${projectId}/knowledge/documents`)
    ])
    indexStatus.value = await statusData.json()
    const docsResult = await docsData.json()
    documents.value = docsResult.items || []
  } catch (error) {
    ElMessage.error('加载索引状态失败')
  } finally {
    loading.value = false
  }
}

const handleRebuildAll = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要重建所有索引吗？此操作可能需要较长时间。',
      '确认重建',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    rebuilding.value = true
    await fetch(`/api/v1/projects/${projectId}/knowledge/rebuild-index`, {
      method: 'POST'
    })
    ElMessage.success('已提交重建索引请求')
    loadIndexStatus()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('操作失败')
    }
  } finally {
    rebuilding.value = false
  }
}

const handleOptimize = async () => {
  optimizing.value = true
  try {
    await fetch(`/api/v1/projects/${projectId}/knowledge/optimize-index`, {
      method: 'POST'
    })
    ElMessage.success('索引优化完成')
    loadIndexStatus()
  } catch (error) {
    ElMessage.error('优化失败')
  } finally {
    optimizing.value = false
  }
}

const handleExportStats = () => {
  ElMessage.info('导出统计信息功能开发中')
}

const getStatusType = (status: string) => {
  const types: Record<string, any> = {
    IDLE: 'info',
    INDEXING: 'warning',
    FAILED: 'danger'
  }
  return types[status] || 'info'
}

const getStatusText = (status: string) => {
  const texts: Record<string, string> = {
    IDLE: '空闲',
    INDEXING: '索引中',
    FAILED: '失败'
  }
  return texts[status] || status
}

const getProgressStatus = (status: string) => {
  if (status === 'FAILED') return 'exception'
  if (status === 'INDEXING') return undefined
  return 'success'
}

const goBack = () => {
  router.push('/projects')
}

onMounted(() => {
  loadIndexStatus()
})
</script>

<style scoped lang="scss">
.index-status {
  .title {
    font-size: 20px;
    font-weight: 500;
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .progress-section {
    .progress-item {
      .progress-label {
        margin-bottom: 12px;
        font-size: 14px;
        color: #666;
      }
    }
  }
}
</style>
