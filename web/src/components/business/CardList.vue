<template>
  <div class="card-list">
    <el-table
      v-loading="loading"
      :data="cards"
      style="width: 100%"
    >
      <el-table-column
        prop="title"
        label="标题"
        min-width="200"
      />
      <el-table-column
        prop="description"
        label="描述"
        min-width="250"
        show-overflow-tooltip
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
        label="优先级"
        width="100"
      >
        <template #default="{ row }">
          <el-tag
            :type="getPriorityType(row.priority)"
            size="small"
          >
            {{ getPriorityText(row.priority) }}
          </el-tag>
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
        v-if="showActions"
        label="操作"
        width="200"
        fixed="right"
      >
        <template #default="{ row }">
          <el-button
            link
            type="primary"
            size="small"
            @click="$emit('view', row)"
          >
            查看
          </el-button>
          <el-dropdown @command="(cmd) => $emit('command', cmd, row)">
            <el-button
              link
              type="primary"
              size="small"
            >
              更多<el-icon class="el-icon--right">
                <ArrowDown />
              </el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <slot
                  name="dropdown-items"
                  :card="row"
                />
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-if="showPagination"
      :current-page="currentPage"
      :page-size="currentPageSize"
      :total="total"
      :page-sizes="[10, 20, 50, 100]"
      layout="total, sizes, prev, pager, next, jumper"
      style="margin-top: 20px; justify-content: flex-end"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { ArrowDown } from '@element-plus/icons-vue'
import { formatDate } from '@/utils/format'
import type { Card } from '@/types/card'

interface Props {
  cards: Card[]
  loading?: boolean
  showActions?: boolean
  showPagination?: boolean
  page?: number
  pageSize?: number
  total?: number
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  showActions: true,
  showPagination: true,
  page: 1,
  pageSize: 20,
  total: 0
})

const emit = defineEmits<{
  (e: 'view', card: Card): void
  (e: 'command', command: string, card: Card): void
  (e: 'page-change', page: number, pageSize: number): void
}>()

const currentPage = ref(props.page)
const currentPageSize = ref(props.pageSize)

watch(() => props.page, (newPage) => {
  currentPage.value = newPage
})

watch(() => props.pageSize, (newSize) => {
  currentPageSize.value = newSize
})

const handleSizeChange = (size: number) => {
  currentPageSize.value = size
  emit('page-change', currentPage.value, currentPageSize.value)
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  emit('page-change', currentPage.value, currentPageSize.value)
}

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
</script>

<style scoped lang="scss">
.card-list {
}
</style>
