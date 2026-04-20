<template>
  <div class="index-status-table">
    <el-table :data="documents" style="width: 100%" v-loading="loading">
      <el-table-column prop="name" label="文件名" min-width="200" show-overflow-tooltip />
      <el-table-column prop="path" label="路径" min-width="250" show-overflow-tooltip />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)" size="small">
            {{ getStatusText(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="分片数" width="100">
        <template #default="{ row }">
          {{ row.chunks || 0 }}
        </template>
      </el-table-column>
      <el-table-column label="索引进度" width="180">
        <template #default="{ row }">
          <el-progress
            v-if="row.status === 'INDEXING'"
            :percentage="50"
            :indeterminate="true"
            :stroke-width="8"
          />
          <span v-else-if="row.status === 'INDEXED'">100%</span>
          <span v-else>0%</span>
        </template>
      </el-table-column>
      <el-table-column label="索引时间" width="180">
        <template #default="{ row }">
          {{ row.indexedAt ? formatDate(row.indexedAt) : '-' }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="150" fixed="right">
        <template #default="{ row }">
          <el-button
            link
            type="primary"
            size="small"
            @click="handleViewChunks(row)"
          >
            查看分片
          </el-button>
          <el-button
            link
            type="warning"
            size="small"
            :disabled="row.status !== 'INDEXED'"
            @click="$emit('reindex', row)"
          >
            重建索引
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-empty v-if="!loading && documents.length === 0" description="暂无文档" />
  </div>
</template>

<script setup lang="ts">
import { ElMessage } from 'element-plus'
import { formatDate } from '@/utils/format'
import type { Document } from '@/types/knowledge'

interface Props {
  documents: Document[]
  loading?: boolean
}

withDefaults(defineProps<Props>(), {
  loading: false
})

const emit = defineEmits<{
  (e: 'reindex', document: Document): void
}>()

const handleViewChunks = (document: Document) => {
  ElMessage.info(`查看 "${document.name}" 的分片信息功能开发中`)
}

const getStatusType = (status: string) => {
  const types: Record<string, any> = {
    PENDING: 'info',
    INDEXING: 'warning',
    INDEXED: 'success',
    FAILED: 'danger'
  }
  return types[status] || 'info'
}

const getStatusText = (status: string) => {
  const texts: Record<string, string> = {
    PENDING: '待索引',
    INDEXING: '索引中',
    INDEXED: '已索引',
    FAILED: '失败'
  }
  return texts[status] || status
}
</script>

<style scoped lang="scss">
.index-status-table {
}
</style>
