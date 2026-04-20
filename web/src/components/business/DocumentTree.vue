<template>
  <div class="document-tree">
    <el-tree
      :data="treeData"
      :props="defaultProps"
      :loading="loading"
      node-key="id"
      default-expand-all
      :expand-on-click-node="false"
      @node-click="handleNodeClick"
    >
      <template #default="{ node, data }">
        <div class="tree-node">
          <el-icon>
            <component :is="getFileIcon(data)" />
          </el-icon>
          <span class="node-label">{{ node.label }}</span>
          <el-tag v-if="data.status" :type="getStatusType(data.status)" size="small">
            {{ getStatusText(data.status) }}
          </el-tag>
        </div>
      </template>
    </el-tree>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Folder, FolderOpened, Document } from '@element-plus/icons-vue'
import type { Document } from '@/types/knowledge'

interface Props {
  documents: Document[]
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  loading: false
})

const emit = defineEmits<{
  (e: 'select', document: Document): void
}>()

const defaultProps = {
  children: 'children',
  label: 'name'
}

const treeData = computed(() => {
  const tree: any[] = []
  const pathMap = new Map<string, any>()

  props.documents.forEach(doc => {
    const parts = doc.path.split('/').filter(p => p)
    let currentPath = ''
    let currentLevel = tree

    parts.forEach((part, index) => {
      currentPath += '/' + part
      const isFile = index === parts.length - 1

      if (!pathMap.has(currentPath)) {
        const node: any = isFile
          ? {
              id: doc.id,
              name: part,
              status: doc.status,
              document: doc
            }
          : {
              id: currentPath,
              name: part,
              children: []
            }

        currentLevel.push(node)
        pathMap.set(currentPath, node)
      }

      const existingNode = pathMap.get(currentPath)
      if (!isFile && existingNode.children) {
        currentLevel = existingNode.children
      }
    })
  })

  return tree
})

const getFileIcon = (data: any) => {
  if (data.children && data.children.length > 0) {
    return Folder
  }
  return Document
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

const handleNodeClick = (data: any) => {
  if (data.document) {
    emit('select', data.document)
  }
}
</script>

<style scoped lang="scss">
.document-tree {
  .tree-node {
    display: flex;
    align-items: center;
    gap: 8px;
    flex: 1;
    overflow: hidden;

    .node-label {
      flex: 1;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
  }

  :deep(.el-tree-node__content) {
    height: 36px;
  }
}
</style>
