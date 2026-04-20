<template>
  <div class="document-table">
    <el-table :data="documents" style="width: 100%" v-loading="loading">
      <el-table-column prop="name" label="文件名" min-width="200" show-overflow-tooltip />
      <el-table-column prop="path" label="路径" min-width="250" show-overflow-tooltip />
      <el-table-column label="大小" width="120">
        <template #default="{ row }">
          {{ formatFileSize(row.size) }}
        </template>
      </el-table-column>
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
      <el-table-column label="索引时间" width="180">
        <template #default="{ row }">
          {{ row.indexedAt ? formatDate(row.indexedAt) : '-' }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button
            link
            type="primary"
            size="small"
            @click="$emit('view', row)"
          >
            查看
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
          <el-button
            link
            type="danger"
            size="small"
            @click="$emit('delete', row)"
          >
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-empty v-if="!loading && documents.length === 0" description="暂无文档" />
  </div>
</template>

<script setup lang="ts">
import { formatFileSize, formatDate } from '@/utils/format'
import type { Document } from '@/types/knowledge'

interface Props {
  documents: Document[]
  loading?: boolean
}

withDefaults(defineProps<Props>(), {
  loading: false
})

defineEmits<{
  (e: 'view', document: Document): void
  (e: 'delete', document: Document): void
  (e: 'reindex', document: Document): void
}>()

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
.document-table {
}
</style>
