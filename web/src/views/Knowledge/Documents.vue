<template>
  <div class="knowledge-documents">
    <el-page-header @back="goBack">
      <template #content>
        <span class="title">知识库文档</span>
      </template>
    </el-page-header>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="6">
        <el-card class="tree-card">
          <template #header>
            <div class="card-header">
              <span>文档目录</span>
              <el-button size="small" @click="refreshTree">
                <el-icon><Refresh /></el-icon>
              </el-button>
            </div>
          </template>
          <DocumentTree
            :documents="documents"
            :loading="loading"
            @select="handleTreeSelect"
          />
        </el-card>
      </el-col>

      <el-col :span="18">
        <el-card class="table-card">
          <template #header>
            <div class="card-header">
              <span>文档列表</span>
              <el-button type="primary" size="small" @click="handleUpload">
                <el-icon><Upload /></el-icon>
                上传文档
              </el-button>
            </div>
          </template>
          <DocumentTable
            :documents="filteredDocuments"
            :loading="loading"
            @view="handleView"
            @delete="handleDelete"
            @reindex="handleReindex"
          />
        </el-card>
      </el-col>
    </el-row>

    <el-dialog
      v-model="uploadDialogVisible"
      title="上传文档"
      width="500px"
    >
      <el-upload
        ref="uploadRef"
        :auto-upload="false"
        :on-change="handleFileChange"
        :on-remove="handleFileRemove"
        :limit="10"
        drag
        multiple
      >
        <el-icon class="el-icon--upload"><upload-filled /></el-icon>
        <div class="el-upload__text">
          将文件拖到此处，或<em>点击上传</em>
        </div>
        <template #tip>
          <div class="el-upload__tip">
            支持 .md, .txt, .pdf, .doc, .docx 格式，单个文件不超过 10MB
          </div>
        </template>
      </el-upload>

      <template #footer>
        <el-button @click="uploadDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="uploading" @click="handleUploadSubmit">
          开始上传
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox, type UploadInstance, type UploadUserFile } from 'element-plus'
import { Refresh, Upload, UploadFilled } from '@element-plus/icons-vue'
import DocumentTree from '@/components/business/DocumentTree.vue'
import DocumentTable from '@/components/business/DocumentTable.vue'
import type { Document } from '@/types/knowledge'

const route = useRoute()
const router = useRouter()

const projectId = route.params.id as string

const loading = ref(false)
const uploading = ref(false)
const documents = ref<Document[]>([])
const selectedDocumentId = ref<string | null>(null)

const uploadDialogVisible = ref(false)
const uploadRef = ref<UploadInstance>()
const uploadFiles = ref<UploadUserFile[]>([])

const filteredDocuments = computed(() => {
  if (!selectedDocumentId.value) {
    return documents.value
  }
  return documents.value.filter(d => d.id === selectedDocumentId.value)
})

const loadDocuments = async () => {
  loading.value = true
  try {
    const data = await fetch(`/api/v1/projects/${projectId}/knowledge/documents`)
    const result = await data.json()
    documents.value = result.items || []
  } catch (error) {
    ElMessage.error('加载文档列表失败')
  } finally {
    loading.value = false
  }
}

const refreshTree = () => {
  loadDocuments()
}

const handleTreeSelect = (document: Document) => {
  selectedDocumentId.value = document.id
}

const handleView = (document: Document) => {
  ElMessage.info('查看文档功能开发中')
}

const handleDelete = async (document: Document) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除文档 "${document.name}" 吗？此操作不可恢复。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await fetch(`/api/v1/knowledge/documents/${document.id}`, { method: 'DELETE' })
    ElMessage.success('删除成功')
    loadDocuments()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const handleReindex = async (document: Document) => {
  try {
    await fetch(`/api/v1/knowledge/documents/${document.id}/reindex`, { method: 'POST' })
    ElMessage.success('已提交重建索引请求')
    loadDocuments()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleUpload = () => {
  uploadDialogVisible.value = true
}

const handleFileChange = (file: UploadUserFile, fileList: UploadUserFile[]) => {
  uploadFiles.value = fileList
}

const handleFileRemove = (file: UploadUserFile, fileList: UploadUserFile[]) => {
  uploadFiles.value = fileList
}

const handleUploadSubmit = async () => {
  if (uploadFiles.value.length === 0) {
    ElMessage.warning('请选择要上传的文件')
    return
  }

  uploading.value = true
  try {
    const formData = new FormData()
    uploadFiles.value.forEach(file => {
      if (file.raw) {
        formData.append('files', file.raw)
      }
    })

    await fetch(`/api/v1/projects/${projectId}/knowledge/upload`, {
      method: 'POST',
      body: formData
    })

    ElMessage.success('上传成功')
    uploadDialogVisible.value = false
    uploadFiles.value = []
    loadDocuments()
  } catch (error) {
    ElMessage.error('上传失败')
  } finally {
    uploading.value = false
  }
}

const goBack = () => {
  router.push('/projects')
}

onMounted(() => {
  loadDocuments()
})
</script>

<style scoped lang="scss">
.knowledge-documents {
  .title {
    font-size: 20px;
    font-weight: 500;
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .tree-card {
    min-height: 600px;
  }

  .table-card {
    min-height: 600px;
  }
}
</style>
