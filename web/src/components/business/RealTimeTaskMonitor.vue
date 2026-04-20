<template>
  <div class="real-time-task-monitor">
    <el-empty
      v-if="tasks.length === 0"
      description="暂无任务"
    />
    <div
      v-else
      class="task-list"
    >
      <div
        v-for="task in tasks"
        :key="task.id"
        class="task-item"
        :class="getTaskClass(task.status)"
      >
        <div class="task-header">
          <div class="task-info">
            <el-tag
              :type="getStatusType(task.status)"
              size="small"
            >
              {{ getStatusText(task.status) }}
            </el-tag>
            <span class="task-type">{{ getTypeText(task.type) }}</span>
          </div>
          <div class="task-time">
            {{ formatRelativeTime(task.createdAt) }}
          </div>
        </div>

        <div class="task-body">
          <div
            v-if="task.skill"
            class="task-skill"
          >
            <el-icon><MagicStick /></el-icon>
            <span>Skill: {{ task.skill }}</span>
          </div>

          <div
            v-if="task.cardId"
            class="task-card"
          >
            <el-icon><Document /></el-icon>
            <span>卡片: {{ task.cardId }}</span>
          </div>

          <div
            v-if="task.error"
            class="task-error"
          >
            <el-alert
              type="error"
              :closable="false"
            >
              {{ task.error }}
            </el-alert>
          </div>
        </div>

        <div
          v-if="task.status === 'RUNNING'"
          class="task-progress"
        >
          <el-progress
            :percentage="getProgress(task)"
            :indeterminate="true"
            :stroke-width="6"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { MagicStick, Document } from '@element-plus/icons-vue'
import { formatRelativeTime } from '@/utils/format'
import type { AgentTask } from '@/types/agent'

interface Props {
  tasks: AgentTask[]
}

defineProps<Props>()

const getTaskClass = (status: string) => {
  return {
    'task-running': status === 'RUNNING',
    'task-completed': status === 'COMPLETED',
    'task-failed': status === 'FAILED',
    'task-pending': status === 'PENDING' || status === 'ASSIGNED',
    'task-cancelled': status === 'CANCELLED'
  }
}

const getStatusType = (status: string) => {
  const types: Record<string, any> = {
    PENDING: 'info',
    ASSIGNED: 'warning',
    RUNNING: 'primary',
    COMPLETED: 'success',
    FAILED: 'danger',
    CANCELLED: 'info'
  }
  return types[status] || 'info'
}

const getStatusText = (status: string) => {
  const texts: Record<string, string> = {
    PENDING: '等待中',
    ASSIGNED: '已分配',
    RUNNING: '运行中',
    COMPLETED: '已完成',
    FAILED: '失败',
    CANCELLED: '已取消'
  }
  return texts[status] || status
}

const getTypeText = (type: string) => {
  const texts: Record<string, string> = {
    CODE_GENERATION: '代码生成',
    TEST: '测试',
    DOCUMENTATION: '文档',
    REVIEW: '审查',
    REFACTORING: '重构',
    CUSTOM: '自定义'
  }
  return texts[type] || type
}

const getProgress = (_task: AgentTask) => {
  return 50
}
</script>

<style scoped lang="scss">
.real-time-task-monitor {
  .task-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
    max-height: 600px;
    overflow-y: auto;
  }

  .task-item {
    background: #fff;
    border: 1px solid #e8e8e8;
    border-radius: 8px;
    padding: 16px;
    transition: all 0.3s;

    &.task-running {
      border-color: #409eff;
      background-color: #ecf5ff;
      animation: pulse 2s infinite;
    }

    &.task-completed {
      border-color: #67c23a;
      background-color: #f0f9ff;
    }

    &.task-failed {
      border-color: #f56c6c;
      background-color: #fef0f0;
    }

    &.task-pending {
      border-color: #dcdfe6;
      background-color: #f5f7fa;
    }

    &.task-cancelled {
      border-color: #909399;
      background-color: #f4f4f5;
      opacity: 0.6;
    }

    .task-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 12px;

      .task-info {
        display: flex;
        align-items: center;
        gap: 8px;

        .task-type {
          font-weight: 500;
          font-size: 14px;
        }
      }

      .task-time {
        font-size: 12px;
        color: #999;
      }
    }

    .task-body {
      .task-skill,
      .task-card {
        display: flex;
        align-items: center;
        gap: 6px;
        font-size: 13px;
        color: #666;
        margin-bottom: 8px;

        .el-icon {
          font-size: 14px;
        }
      }

      .task-error {
        margin-top: 8px;
      }
    }

    .task-progress {
      margin-top: 12px;
    }
  }
}

@keyframes pulse {
  0% {
    box-shadow: 0 0 0 0 rgba(64, 158, 255, 0.4);
  }
  70% {
    box-shadow: 0 0 0 10px rgba(64, 158, 255, 0);
  }
  100% {
    box-shadow: 0 0 0 0 rgba(64, 158, 255, 0);
  }
}
</style>
