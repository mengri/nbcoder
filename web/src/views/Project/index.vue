<template>
  <div class="project-list">
    <el-page-header @back="goBack">
      <template #content>
        <div class="page-header-content">
          <span class="title">项目管理</span>
          <el-button
            type="primary"
            @click="handleCreate"
          >
            <el-icon><Plus /></el-icon>
            新建项目
          </el-button>
        </div>
      </template>
    </el-page-header>

    <div class="search-bar">
      <el-input
        v-model="searchKeyword"
        placeholder="搜索项目名称或描述"
        clearable
        @input="handleSearch"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
    </div>

    <el-row
      v-loading="loading"
      :gutter="20"
    >
      <el-col
        v-for="project in filteredProjects"
        :key="project.id"
        :xs="24"
        :sm="12"
        :md="8"
        :lg="6"
        class="project-col"
      >
        <el-card
          class="project-card"
          :body-style="{ padding: '20px' }"
        >
          <template #header>
            <div class="card-header">
              <span class="project-name">{{ project.name }}</span>
              <el-dropdown @command="(cmd) => handleCommand(cmd, project)">
                <el-icon class="more-icon">
                  <MoreFilled />
                </el-icon>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="edit">
                      编辑
                    </el-dropdown-item>
                    <el-dropdown-item command="settings">
                      设置
                    </el-dropdown-item>
                    <el-dropdown-item
                      command="delete"
                      divided
                    >
                      删除
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </template>

          <div class="project-info">
            <p class="description">
              {{ truncateText(project.description, 80) }}
            </p>
            <div class="meta">
              <el-tag
                :type="getStatusType(project.status)"
                size="small"
              >
                {{ getStatusText(project.status) }}
              </el-tag>
              <span class="time">{{ formatDate(project.createdAt) }}</span>
            </div>
          </div>

          <template #footer>
            <el-button
              type="primary"
              size="small"
              @click="handleEnter(project)"
            >
              进入项目
            </el-button>
          </template>
        </el-card>
      </el-col>

      <el-col
        v-if="filteredProjects.length === 0"
        :span="24"
      >
        <el-empty description="暂无项目" />
      </el-col>
    </el-row>

    <ProjectDialog
      v-model="dialogVisible"
      :project="editingProject"
      @success="handleDialogSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search, MoreFilled } from '@element-plus/icons-vue'
import { useProjectStore } from '@/stores/project'
import ProjectDialog from './components/ProjectDialog.vue'
import { formatDate, truncateText } from '@/utils/format'
import type { Project } from '@/types/project'

const router = useRouter()
const projectStore = useProjectStore()

const searchKeyword = ref('')
const loading = ref(false)
const dialogVisible = ref(false)
const editingProject = ref<Project | null>(null)

const filteredProjects = computed(() => {
  const projects = projectStore.projects || []
  console.log('Computing filteredProjects, projects:', projects)
  console.log('Projects type:', typeof projects, 'Is array:', Array.isArray(projects))
  if (!searchKeyword.value) {
    return projects
  }
  const keyword = searchKeyword.value.toLowerCase()
  return projects.filter(
    p =>
      p.name.toLowerCase().includes(keyword) ||
      p.description?.toLowerCase().includes(keyword)
  )
})

const loadProjects = async () => {
  loading.value = true
  try {
    await projectStore.loadProjects()
    console.log('✅ Load completed')
    console.log('Store projects:', projectStore.projects)
    console.log('Projects length:', projectStore.projects?.length)
    console.log('Filtered projects:', filteredProjects.value)
  } catch (error) {
    console.error('❌ Failed to load projects:', error)
    ElMessage.error('加载项目列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
}

const handleCreate = () => {
  editingProject.value = null
  dialogVisible.value = true
}

const handleEdit = (project: Project) => {
  editingProject.value = project
  dialogVisible.value = true
}

const handleSettings = (project: Project) => {
  router.push(`/projects/${project.id}/settings`)
}

const handleDelete = async (project: Project) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除项目 "${project.name}" 吗？此操作不可恢复。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await projectStore.deleteProject(project.id)
    ElMessage.success('删除成功')
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const handleCommand = (command: string, project: Project) => {
  switch (command) {
    case 'edit':
      handleEdit(project)
      break
    case 'settings':
      handleSettings(project)
      break
    case 'delete':
      handleDelete(project)
      break
  }
}

const handleEnter = (project: Project) => {
  projectStore.setCurrentProject(project)
  router.push(`/projects/${project.id}/cards`)
}

const handleDialogSuccess = () => {
  dialogVisible.value = false
  loadProjects()
}

const getStatusType = (status: string) => {
  const types: Record<string, any> = {
    ACTIVE: 'success',
    ARCHIVED: 'info',
    DELETED: 'danger'
  }
  return types[status] || 'info'
}

const getStatusText = (status: string) => {
  const texts: Record<string, string> = {
    ACTIVE: '活跃',
    ARCHIVED: '已归档',
    DELETED: '已删除'
  }
  return texts[status] || status
}

const goBack = () => {
  router.push('/')
}

onMounted(() => {
  loadProjects()
})
</script>

<style scoped lang="scss">
.project-list {
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

  .search-bar {
    margin: 20px 0;

    .el-input {
      max-width: 400px;
    }
  }

  .project-col {
    margin-bottom: 20px;

    .project-card {
      height: 100%;
      cursor: pointer;
      transition: all 0.3s;

      &:hover {
        box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
        transform: translateY(-2px);
      }

      .card-header {
        display: flex;
        justify-content: space-between;
        align-items: center;

        .project-name {
          font-weight: 500;
          font-size: 16px;
        }

        .more-icon {
          cursor: pointer;
          color: #999;

          &:hover {
            color: #409eff;
          }
        }
      }

      .project-info {
        .description {
          color: #666;
          font-size: 14px;
          line-height: 1.5;
          margin-bottom: 12px;
          min-height: 42px;
        }

        .meta {
          display: flex;
          justify-content: space-between;
          align-items: center;

          .time {
            font-size: 12px;
            color: #999;
          }
        }
      }
    }
  }
}
</style>
