<template>
  <el-dialog
    :model-value="modelValue"
    :title="isEdit ? '编辑项目' : '新建项目'"
    width="600px"
    @update:model-value="handleClose"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="100px"
    >
      <el-form-item label="项目名称" prop="name">
        <el-input v-model="form.name" placeholder="请输入项目名称" />
      </el-form-item>

      <el-form-item label="项目描述" prop="description">
        <el-input
          v-model="form.description"
          type="textarea"
          :rows="3"
          placeholder="请输入项目描述"
        />
      </el-form-item>

      <el-form-item label="Git 仓库" prop="gitRepoUrl">
        <el-input
          v-model="form.gitRepoUrl"
          placeholder="请输入 Git 仓库 URL（可选）"
        />
      </el-form-item>

      <el-form-item label="分支" prop="gitBranch">
        <el-input
          v-model="form.gitBranch"
          placeholder="请输入分支名称（可选）"
        />
      </el-form-item>

      <el-form-item label="开发规范" prop="standards">
        <el-input
          v-model="form.standards"
          type="textarea"
          :rows="2"
          placeholder="请输入开发规范（可选）"
        />
      </el-form-item>

      <el-form-item label="分支策略" prop="branchPolicy">
        <el-input
          v-model="form.branchPolicy"
          type="textarea"
          :rows="2"
          placeholder="请输入分支策略（可选）"
        />
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" :loading="loading" @click="handleSubmit">
        {{ isEdit ? '保存' : '创建' }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { useProjectStore } from '@/stores/project'
import type { Project, CreateProjectDto } from '@/types/project'

interface Props {
  modelValue: boolean
  project?: Project | null
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void
  (e: 'success'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const projectStore = useProjectStore()
const formRef = ref<FormInstance>()
const loading = ref(false)

const isEdit = computed(() => !!props.project)

const form = ref<CreateProjectDto>({
  name: '',
  description: '',
  gitRepoUrl: '',
  gitBranch: '',
  standards: '',
  branchPolicy: ''
})

const rules: FormRules = {
  name: [
    { required: true, message: '请输入项目名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  description: [
    { required: true, message: '请输入项目描述', trigger: 'blur' },
    { min: 5, max: 500, message: '长度在 5 到 500 个字符', trigger: 'blur' }
  ]
}

watch(
  () => props.project,
  (newProject) => {
    if (newProject) {
      form.value = {
        name: newProject.name,
        description: newProject.description,
        gitRepoUrl: newProject.gitRepoUrl || '',
        gitBranch: newProject.gitBranch || '',
        standards: newProject.standards || '',
        branchPolicy: newProject.branchPolicy || ''
      }
    } else {
      form.value = {
        name: '',
        description: '',
        gitRepoUrl: '',
        gitBranch: '',
        standards: '',
        branchPolicy: ''
      }
    }
  },
  { immediate: true }
)

const handleClose = () => {
  emit('update:modelValue', false)
  formRef.value?.resetFields()
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        if (isEdit.value && props.project) {
          await projectStore.updateProject(props.project.id, form.value)
          ElMessage.success('更新成功')
        } else {
          await projectStore.createProject(form.value)
          ElMessage.success('创建成功')
        }
        emit('success')
      } catch (error) {
        ElMessage.error(isEdit.value ? '更新失败' : '创建失败')
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped lang="scss">
</style>
