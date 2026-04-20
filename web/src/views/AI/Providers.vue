<template>
  <div class="ai-providers">
    <el-page-header @back="goBack">
      <template #content>
        <span class="title">AI Provider 管理</span>
      </template>
    </el-page-header>

    <el-card
      class="table-card"
      style="margin-top: 20px"
    >
      <template #header>
        <div class="card-header">
          <span>Provider 列表</span>
          <el-button
            type="primary"
            size="small"
            @click="handleCreate"
          >
            <el-icon><Plus /></el-icon>
            添加 Provider
          </el-button>
        </div>
      </template>
      <el-table
        v-loading="loading"
        :data="providers"
        style="width: 100%"
      >
        <el-table-column
          prop="name"
          label="名称"
          min-width="150"
        />
        <el-table-column
          prop="type"
          label="类型"
          width="120"
        />
        <el-table-column
          label="状态"
          width="100"
        >
          <template #default="{ row }">
            <el-tag
              :type="getStatusType(row.status)"
              size="small"
            >
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          label="模型数量"
          width="100"
        >
          <template #default="{ row }">
            {{ row.models?.length || 0 }}
          </template>
        </el-table-column>
        <el-table-column
          label="可用模型"
          width="150"
        >
          <template #default="{ row }">
            <div class="model-list">
              <el-tag
                v-for="model in getAvailableModels(row)"
                :key="model.id"
                size="small"
                style="margin: 2px"
              >
                {{ model.name }}
              </el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column
          prop="createdAt"
          label="创建时间"
          width="180"
        >
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column
          label="操作"
          width="250"
          fixed="right"
        >
          <template #default="{ row }">
            <el-button
              link
              type="primary"
              size="small"
              @click="handleEdit(row)"
            >
              编辑
            </el-button>
            <el-button
              link
              type="success"
              size="small"
              @click="handleTest(row)"
            >
              测试
            </el-button>
            <el-button
              link
              type="danger"
              size="small"
              @click="handleDelete(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-empty
        v-if="!loading && providers.length === 0"
        description="暂无 Provider"
      />
    </el-card>

    <ProviderDialog
      v-model="dialogVisible"
      :provider="editingProvider"
      @success="handleDialogSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { formatDate } from '@/utils/format'
import ProviderDialog from './components/ProviderDialog.vue'
import type { AIProvider } from '@/types/ai'

const router = useRouter()

const loading = ref(false)
const providers = ref<AIProvider[]>([])
const dialogVisible = ref(false)
const editingProvider = ref<AIProvider | null>(null)

const loadProviders = async () => {
  loading.value = true
  try {
    const data = await fetch('/api/v1/ai/providers')
    const result = await data.json()
    providers.value = result.items || []
  } catch (error) {
    ElMessage.error('加载 Provider 列表失败')
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  editingProvider.value = null
  dialogVisible.value = true
}

const handleEdit = (provider: AIProvider) => {
  editingProvider.value = provider
  dialogVisible.value = true
}

const handleTest = async (provider: AIProvider) => {
  try {
    ElMessage.info('测试 Provider 连接...')
    await fetch(`/api/v1/ai/providers/${provider.id}/test`, { method: 'POST' })
    ElMessage.success('连接测试成功')
  } catch (error) {
    ElMessage.error('连接测试失败')
  }
}

const handleDelete = async (provider: AIProvider) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除 Provider "${provider.name}" 吗？此操作不可恢复。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await fetch(`/api/v1/ai/providers/${provider.id}`, { method: 'DELETE' })
    ElMessage.success('删除成功')
    loadProviders()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const handleDialogSuccess = () => {
  dialogVisible.value = false
  loadProviders()
}

const getAvailableModels = (provider: AIProvider) => {
  return provider.models?.filter(m => m.status === 'ACTIVE') || []
}

const getStatusType = (status: string) => {
  const types: Record<string, any> = {
    ACTIVE: 'success',
    INACTIVE: 'info',
    UNAVAILABLE: 'danger'
  }
  return types[status] || 'info'
}

const getStatusText = (status: string) => {
  const texts: Record<string, string> = {
    ACTIVE: '活跃',
    INACTIVE: '未激活',
    UNAVAILABLE: '不可用'
  }
  return texts[status] || status
}

const goBack = () => {
  router.push('/projects')
}

onMounted(() => {
  loadProviders()
})
</script>

<style scoped lang="scss">
.ai-providers {
  .title {
    font-size: 20px;
    font-weight: 500;
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .model-list {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
  }
}
</style>
