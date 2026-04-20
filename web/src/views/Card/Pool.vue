<template>
  <div class="card-pool">
    <el-page-header @back="goBack">
      <template #content>
        <span class="title">需求卡片池</span>
      </template>
    </el-page-header>

    <el-card class="filter-card">
      <el-form :inline="true" :model="filters">
        <el-form-item label="状态">
          <el-select
            v-model="filters.status"
            placeholder="全部状态"
            clearable
            @change="loadCards"
          >
            <el-option
              v-for="status in cardStatuses"
              :key="status.value"
              :label="status.label"
              :value="status.value"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="优先级">
          <el-select
            v-model="filters.priority"
            placeholder="全部优先级"
            clearable
            @change="loadCards"
          >
            <el-option
              v-for="priority in priorities"
              :key="priority.value"
              :label="priority.label"
              :value="priority.value"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="搜索">
          <el-input
            v-model="filters.keyword"
            placeholder="搜索卡片标题或描述"
            clearable
            @input="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            新建卡片
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="table-card" v-loading="loading">
      <el-table :data="cards" style="width: 100%">
        <el-table-column prop="title" label="标题" min-width="200" />
        <el-table-column prop="description" label="描述" min-width="250" show-overflow-tooltip />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="优先级" width="100">
          <template #default="{ row }">
            <el-tag :type="getPriorityType(row.priority)" size="small">
              {{ getPriorityText(row.priority) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="依赖" width="80">
          <template #default="{ row }">
            <span>{{ row.dependencies.length }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280" fixed="right">
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
                更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="confirm" :disabled="!canConfirm(row)">
                    确认
                  </el-dropdown-item>
                  <el-dropdown-item command="start" :disabled="!canStart(row)">
                    启动
                  </el-dropdown-item>
                  <el-dropdown-item command="complete" :disabled="!canComplete(row)">
                    完成
                  </el-dropdown-item>
                  <el-dropdown-item command="abandon" :disabled="!canAbandon(row)">
                    废弃
                  </el-dropdown-item>
                  <el-dropdown-item command="delete" divided :disabled="!canDelete(row)">
                    删除
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
        @size-change="loadCards"
        @current-change="loadCards"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Plus, ArrowDown } from '@element-plus/icons-vue'
import { useCardStore } from '@/stores/card'
import { useProjectStore } from '@/stores/project'
import { formatDate } from '@/utils/format'
import type { Card, CardStatus, CardPriority } from '@/types/card'

const route = useRoute()
const router = useRouter()
const cardStore = useCardStore()
const projectStore = useProjectStore()

const loading = ref(false)
const cards = ref<Card[]>([])

const filters = ref({
  status: '',
  priority: '',
  keyword: ''
})

const pagination = ref({
  page: 1,
  size: 20,
  total: 0
})

const cardStatuses = [
  { label: '草稿', value: 'DRAFT' },
  { label: '待确认', value: 'PENDING' },
  { label: '已确认', value: 'CONFIRMED' },
  { label: '进行中', value: 'IN_PROGRESS' },
  { label: '已完成', value: 'COMPLETED' },
  { label: '被取代', value: 'SUPPLANTED' },
  { label: '已废弃', value: 'ABANDONED' }
]

const priorities = [
  { label: '低', value: 'LOW' },
  { label: '中', value: 'MEDIUM' },
  { label: '高', value: 'HIGH' },
  { label: '紧急', value: 'CRITICAL' }
]

const projectId = route.params.id as string

const loadCards = async () => {
  loading.value = true
  try {
    const result = await cardStore.loadCards(projectId, {
      page: pagination.value.page,
      size: pagination.value.size,
      keyword: filters.value.keyword
    })
    cards.value = result.items
    pagination.value.total = result.total

    if (filters.value.status) {
      cards.value = cards.value.filter(c => c.status === filters.value.status)
    }
    if (filters.value.priority) {
      cards.value = cards.value.filter(c => c.priority === filters.value.priority)
    }
  } catch (error) {
    ElMessage.error('加载卡片列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.value.page = 1
  loadCards()
}

const handleCreate = () => {
  router.push(`/projects/${projectId}/cards/new`)
}

const handleView = (card: Card) => {
  cardStore.currentCard = card
  router.push(`/projects/${projectId}/cards/${card.id}`)
}

const handleCommand = async (command: string, card: Card) => {
  switch (command) {
    case 'confirm':
      await handleConfirm(card)
      break
    case 'start':
      await handleStart(card)
      break
    case 'complete':
      await handleComplete(card)
      break
    case 'abandon':
      await handleAbandon(card)
      break
    case 'delete':
      await handleDelete(card)
      break
  }
}

const handleConfirm = async (card: Card) => {
  try {
    await cardStore.confirmCard(card.id)
    ElMessage.success('已确认')
    loadCards()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleStart = async (card: Card) => {
  try {
    await cardStore.startCard(card.id)
    ElMessage.success('已启动')
    loadCards()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleComplete = async (card: Card) => {
  try {
    await cardStore.completeCard(card.id)
    ElMessage.success('已完成')
    loadCards()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleAbandon = async (card: Card) => {
  try {
    await ElMessageBox.confirm('确定要废弃此卡片吗？', '确认', {
      type: 'warning'
    })
    await cardStore.abandonCard(card.id)
    ElMessage.success('已废弃')
    loadCards()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('操作失败')
    }
  }
}

const handleDelete = async (card: Card) => {
  try {
    await ElMessageBox.confirm('确定要删除此卡片吗？此操作不可恢复。', '确认删除', {
      type: 'warning'
    })
    await cardStore.deleteCard(card.id)
    ElMessage.success('已删除')
    loadCards()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('操作失败')
    }
  }
}

const canConfirm = (card: Card) => card.status === 'PENDING' || card.status === 'DRAFT'
const canStart = (card: Card) => card.status === 'CONFIRMED'
const canComplete = (card: Card) => card.status === 'IN_PROGRESS'
const canAbandon = (card: Card) => !['COMPLETED', 'ABANDONED'].includes(card.status)
const canDelete = (card: Card) => ['DRAFT', 'PENDING', 'ABANDONED'].includes(card.status)

const getStatusType = (status: string) => {
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

const getStatusText = (status: string) => {
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

const getPriorityType = (priority: string) => {
  const types: Record<string, any> = {
    LOW: 'info',
    MEDIUM: '',
    HIGH: 'warning',
    CRITICAL: 'danger'
  }
  return types[priority] || ''
}

const getPriorityText = (priority: string) => {
  const texts: Record<string, string> = {
    LOW: '低',
    MEDIUM: '中',
    HIGH: '高',
    CRITICAL: '紧急'
  }
  return texts[priority] || priority
}

const goBack = () => {
  router.push('/projects')
}

onMounted(() => {
  const project = projectStore.projects.find(p => p.id === projectId)
  if (project) {
    projectStore.setCurrentProject(project)
  }
  loadCards()
})
</script>

<style scoped lang="scss">
.card-pool {
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
