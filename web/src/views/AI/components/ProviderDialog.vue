<template>
  <el-dialog
    :model-value="modelValue"
    :title="isEdit ? '编辑 Provider' : '添加 Provider'"
    width="600px"
    @update:model-value="handleClose"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="100px"
    >
      <el-form-item label="名称" prop="name">
        <el-input v-model="form.name" placeholder="请输入 Provider 名称" />
      </el-form-item>

      <el-form-item label="类型" prop="type">
        <el-select v-model="form.type" placeholder="请选择 Provider 类型" style="width: 100%">
          <el-option label="OpenAI" value="openai" />
          <el-option label="Azure OpenAI" value="azure-openai" />
          <el-option label="Anthropic" value="anthropic" />
          <el-option label="Google" value="google" />
          <el-option label="Hugging Face" value="huggingface" />
          <el-option label="自定义" value="custom" />
        </el-select>
      </el-form-item>

      <el-form-item label="API Key" prop="apiKey">
        <el-input
          v-model="form.apiKey"
          type="password"
          placeholder="请输入 API Key"
          show-password
        />
      </el-form-item>

      <el-form-item label="Base URL">
        <el-input v-model="form.baseUrl" placeholder="请输入 Base URL（可选）" />
      </el-form-item>

      <el-form-item label="测试连接">
        <el-button :loading="testing" @click="testConnection">
          <el-icon><Connection /></el-icon>
          测试连接
        </el-button>
        <span v-if="testResult" :class="['test-result', testResult.success ? 'success' : 'error']">
          {{ testResult.message }}
        </span>
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
import { Connection } from '@element-plus/icons-vue'
import type { AIProvider, CreateProviderDto } from '@/types/ai'

interface Props {
  modelValue: boolean
  provider?: AIProvider | null
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void
  (e: 'success'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const formRef = ref<FormInstance>()
const loading = ref(false)
const testing = ref(false)
const testResult = ref<{ success: boolean; message: string } | null>(null)

const isEdit = computed(() => !!props.provider)

const form = ref<CreateProviderDto>({
  name: '',
  type: 'openai',
  apiKey: '',
  baseUrl: ''
})

const rules: FormRules = {
  name: [
    { required: true, message: '请输入 Provider 名称', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择 Provider 类型', trigger: 'change' }
  ],
  apiKey: [
    { required: true, message: '请输入 API Key', trigger: 'blur' }
  ]
}

watch(
  () => props.provider,
  (newProvider) => {
    if (newProvider) {
      form.value = {
        name: newProvider.name,
        type: newProvider.type,
        apiKey: newProvider.apiKey,
        baseUrl: newProvider.baseUrl || ''
      }
    } else {
      form.value = {
        name: '',
        type: 'openai',
        apiKey: '',
        baseUrl: ''
      }
    }
    testResult.value = null
  },
  { immediate: true }
)

const handleClose = () => {
  emit('update:modelValue', false)
  formRef.value?.resetFields()
  testResult.value = null
}

const testConnection = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      testing.value = true
      testResult.value = null
      try {
        const response = await fetch('/api/v1/ai/providers/test', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify(form.value)
        })
        const result = await response.json()

        if (result.success) {
          testResult.value = { success: true, message: '连接成功' }
          ElMessage.success('连接测试成功')
        } else {
          testResult.value = { success: false, message: result.message || '连接失败' }
          ElMessage.error(result.message || '连接失败')
        }
      } catch (error) {
        testResult.value = { success: false, message: '连接失败' }
        ElMessage.error('连接失败')
      } finally {
        testing.value = false
      }
    }
  })
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const url = isEdit.value
          ? `/api/v1/ai/providers/${props.provider?.id}`
          : '/api/v1/ai/providers'
        const method = isEdit.value ? 'PUT' : 'POST'

        await fetch(url, {
          method,
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify(form.value)
        })

        ElMessage.success(isEdit.value ? '更新成功' : '创建成功')
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
.test-result {
  margin-left: 12px;
  font-size: 13px;

  &.success {
    color: #67c23a;
  }

  &.error {
    color: #f56c6c;
  }
}
</style>
