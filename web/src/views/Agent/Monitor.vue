<template>
  <div class="agent-monitor">
    <el-page-header @back="goBack">
      <template #content>
        <span class="title">Agent 实时监控</span>
      </template>
    </el-page-header>

    <el-row
      :gutter="20"
      style="margin-top: 20px"
    >
      <el-col :span="8">
        <el-card>
          <el-statistic
            title="运行中任务"
            :value="stats.runningTasks"
          >
            <template #suffix>
              <el-icon style="vertical-align: -0.125em">
                <Loading />
              </el-icon>
            </template>
          </el-statistic>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card>
          <el-statistic
            title="今日完成"
            :value="stats.todayCompleted"
          />
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card>
          <el-statistic
            title="平均耗时"
            :value="stats.avgDuration"
            suffix="s"
          />
        </el-card>
      </el-col>
    </el-row>

    <el-row
      :gutter="20"
      style="margin-top: 20px"
    >
      <el-col :span="16">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>执行时间线</span>
              <el-tag
                :type="wsConnected ? 'success' : 'danger'"
                size="small"
              >
                {{ wsConnected ? '已连接' : '未连接' }}
              </el-tag>
            </div>
          </template>
          <ExecutionTimeline :events="executionEvents" />
        </el-card>
      </el-col>

      <el-col :span="8">
        <el-card>
          <template #header>
            <span>Skill 调用统计</span>
          </template>
          <SkillStatistics :skills="skillStats" />
        </el-card>
      </el-col>
    </el-row>

    <el-card style="margin-top: 20px">
      <template #header>
        <div class="card-header">
          <span>实时任务监控</span>
          <el-button
            size="small"
            @click="clearHistory"
          >
            清空历史
          </el-button>
        </div>
      </template>
      <RealTimeTaskMonitor :tasks="realTimeTasks" />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Loading } from '@element-plus/icons-vue'
import { useWebSocket } from '@/composables/useWebSocket'
import ExecutionTimeline from '@/components/business/ExecutionTimeline.vue'
import SkillStatistics from '@/components/business/SkillStatistics.vue'
import RealTimeTaskMonitor from '@/components/business/RealTimeTaskMonitor.vue'

const route = useRoute()
const router = useRouter()

const projectId = route.params.id as string
const wsUrl = `ws://localhost:8080/ws/projects/${projectId}/agent-monitor`

const { connected: wsConnected, connect, disconnect, onMessage } = useWebSocket(wsUrl)

const stats = ref({
  runningTasks: 0,
  todayCompleted: 0,
  avgDuration: 0
})

const executionEvents = ref<any[]>([])
const skillStats = ref<any[]>([])
const realTimeTasks = ref<any[]>([])

const handleWebSocketMessage = (data: any) => {
  switch (data.type) {
    case 'task_started':
      handleTaskStarted(data.payload)
      break
    case 'task_completed':
      handleTaskCompleted(data.payload)
      break
    case 'task_failed':
      handleTaskFailed(data.payload)
      break
    case 'skill_called':
      handleSkillCalled(data.payload)
      break
    case 'stats_update':
      handleStatsUpdate(data.payload)
      break
  }
}

const handleTaskStarted = (task: any) => {
  stats.value.runningTasks++
  executionEvents.value.unshift({
    type: 'start',
    task,
    timestamp: new Date()
  })
  realTimeTasks.value.unshift({ ...task, status: 'RUNNING' })

  if (executionEvents.value.length > 100) {
    executionEvents.value.pop()
  }
  if (realTimeTasks.value.length > 50) {
    realTimeTasks.value.pop()
  }
}

const handleTaskCompleted = (task: any) => {
  stats.value.runningTasks--
  stats.value.todayCompleted++
  executionEvents.value.unshift({
    type: 'complete',
    task,
    timestamp: new Date()
  })

  const index = realTimeTasks.value.findIndex(t => t.id === task.id)
  if (index !== -1) {
    realTimeTasks.value[index] = { ...task, status: 'COMPLETED' }
  }
}

const handleTaskFailed = (task: any) => {
  stats.value.runningTasks--
  executionEvents.value.unshift({
    type: 'fail',
    task,
    timestamp: new Date()
  })

  const index = realTimeTasks.value.findIndex(t => t.id === task.id)
  if (index !== -1) {
    realTimeTasks.value[index] = { ...task, status: 'FAILED' }
  }
}

const handleSkillCalled = (data: any) => {
  const { skill, count } = data
  const existing = skillStats.value.find(s => s.name === skill)
  if (existing) {
    existing.count += count
  } else {
    skillStats.value.push({ name: skill, count })
  }
}

const handleStatsUpdate = (data: any) => {
  stats.value = { ...stats.value, ...data }
}

const clearHistory = () => {
  executionEvents.value = []
  realTimeTasks.value = []
  skillStats.value = []
}

const goBack = () => {
  router.push('/projects')
}

onMounted(() => {
  onMessage(handleWebSocketMessage)
  connect()
})

onUnmounted(() => {
  disconnect()
})
</script>

<style scoped lang="scss">
.agent-monitor {
  .title {
    font-size: 20px;
    font-weight: 500;
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
}
</style>
