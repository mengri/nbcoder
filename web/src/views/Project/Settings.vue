<template>
  <div class="project-settings">
    <el-page-header @back="goBack">
      <template #content>
        <span class="title">项目设置</span>
      </template>
    </el-page-header>

    <el-card class="settings-card" v-loading="loading">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="基本信息" name="basic">
          <el-form
            ref="basicFormRef"
            :model="basicForm"
            :rules="basicRules"
            label-width="120px"
          >
            <el-form-item label="项目名称" prop="name">
              <el-input v-model="basicForm.name" />
            </el-form-item>

            <el-form-item label="项目描述" prop="description">
              <el-input
                v-model="basicForm.description"
                type="textarea"
                :rows="4"
              />
            </el-form-item>

            <el-form-item label="项目状态">
              <el-select v-model="basicForm.status">
                <el-option label="活跃" value="ACTIVE" />
                <el-option label="归档" value="ARCHIVED" />
              </el-select>
            </el-form-item>

            <el-form-item>
              <el-button type="primary" :loading="saving" @click="saveBasic">
                保存
              </el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="开发规范" name="standards">
          <el-form
            ref="standardsFormRef"
            :model="standardsForm"
            label-width="120px"
          >
            <el-form-item label="开发规范">
              <el-input
                v-model="standardsForm.standards"
                type="textarea"
                :rows="10"
                placeholder="请输入开发规范说明..."
              />
            </el-form-item>

            <el-form-item>
              <el-button type="primary" :loading="saving" @click="saveStandards">
                保存
              </el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="分支策略" name="branch">
          <el-form
            ref="branchFormRef"
            :model="branchForm"
            label-width="120px"
          >
            <el-form-item label="分支策略">
              <el-input
                v-model="branchForm.branchPolicy"
                type="textarea"
                :rows="10"
                placeholder="请输入分支策略说明..."
              />
            </el-form-item>

            <el-form-item>
              <el-button type="primary" :loading="saving" @click="saveBranch">
                保存
              </el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="Git 仓库" name="git">
          <el-form
            ref="gitFormRef"
            :model="gitForm"
            label-width="120px"
          >
            <el-form-item label="仓库 URL">
              <el-input v-model="gitForm.gitRepoUrl" />
            </el-form-item>

            <el-form-item label="默认分支">
              <el-input v-model="gitForm.gitBranch" />
            </el-form-item>

            <el-form-item>
              <el-button type="primary" :loading="saving" @click="saveGit">
                保存
              </el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { useProjectStore } from '@/stores/project'
import type { Project } from '@/types/project'

const route = useRoute()
const router = useRouter()
const projectStore = useProjectStore()

const activeTab = ref('basic')
const loading = ref(false)
const saving = ref(false)
const basicFormRef = ref<FormInstance>()
const standardsFormRef = ref<FormInstance>()
const branchFormRef = ref<FormInstance>()
const gitFormRef = ref<FormInstance>()

const projectId = route.params.id as string
let project: Project | null = null

const basicForm = ref({
  name: '',
  description: '',
  status: 'ACTIVE'
})

const basicRules: FormRules = {
  name: [
    { required: true, message: '请输入项目名称', trigger: 'blur' }
  ],
  description: [
    { required: true, message: '请输入项目描述', trigger: 'blur' }
  ]
}

const standardsForm = ref({
  standards: ''
})

const branchForm = ref({
  branchPolicy: ''
})

const gitForm = ref({
  gitRepoUrl: '',
  gitBranch: ''
})

const loadProject = async () => {
  loading.value = true
  try {
    project = await projectStore.getProject(projectId)

    if (project) {
      basicForm.value = {
        name: project.name,
        description: project.description,
        status: project.status
      }

      standardsForm.value = {
        standards: project.standards || ''
      }

      branchForm.value = {
        branchPolicy: project.branchPolicy || ''
      }

      gitForm.value = {
        gitRepoUrl: project.gitRepoUrl || '',
        gitBranch: project.gitBranch || ''
      }
    }
  } catch (error) {
    ElMessage.error('加载项目信息失败')
  } finally {
    loading.value = false
  }
}

const saveBasic = async () => {
  if (!basicFormRef.value) return
  await basicFormRef.value.validate(async (valid) => {
    if (valid) {
      saving.value = true
      try {
        await projectStore.updateProject(projectId, {
          name: basicForm.value.name,
          description: basicForm.value.description,
          status: basicForm.value.status as any
        })
        ElMessage.success('保存成功')
      } catch (error) {
        ElMessage.error('保存失败')
      } finally {
        saving.value = false
      }
    }
  })
}

const saveStandards = async () => {
  saving.value = true
  try {
    await projectStore.updateProject(projectId, {
      standards: standardsForm.value.standards
    })
    ElMessage.success('保存成功')
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

const saveBranch = async () => {
  saving.value = true
  try {
    await projectStore.updateProject(projectId, {
      branchPolicy: branchForm.value.branchPolicy
    })
    ElMessage.success('保存成功')
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

const saveGit = async () => {
  saving.value = true
  try {
    await projectStore.updateProject(projectId, {
      gitRepoUrl: gitForm.value.gitRepoUrl,
      gitBranch: gitForm.value.gitBranch
    })
    ElMessage.success('保存成功')
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

const goBack = () => {
  router.push('/projects')
}

onMounted(() => {
  loadProject()
})
</script>

<style scoped lang="scss">
.project-settings {
  .title {
    font-size: 20px;
    font-weight: 500;
  }

  .settings-card {
    margin-top: 20px;
    min-height: 500px;
  }
}
</style>
