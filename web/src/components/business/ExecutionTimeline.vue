<template>
  <div class="execution-timeline">
    <el-empty
      v-if="events.length === 0"
      description="暂无执行事件"
    />
    <el-timeline v-else>
      <el-timeline-item
        v-for="(event, index) in events"
        :key="index"
        :timestamp="formatTime(event.timestamp)"
        :type="getEventType(event.type)"
        :hollow="event.type === 'start'"
        placement="top"
      >
        <el-card class="event-card">
          <div class="event-header">
            <div class="event-info">
              <el-tag
                :type="getEventTagType(event.type)"
                size="small"
              >
                {{ getEventText(event.type) }}
              </el-tag>
              <span class="task-type">{{ getTypeText(event.task?.type) }}</span>
            </div>
            <div class="event-time">
              {{ formatRelativeTime(event.timestamp) }}
            </div>
          </div>

          <div class="event-body">
            <div
              v-if="event.task?.skill"
              class="event-skill"
            >
              <el-icon><MagicStick /></el-icon>
              <span>Skill: {{ event.task.skill }}</span>
            </div>

            <div
              v-if="event.task?.cardId"
              class="event-card"
            >
              <el-icon><Document /></el-icon>
              <span>卡片: {{ event.task.cardId }}</span>
            </div>

            <div
              v-if="event.task?.error"
              class="event-error"
            >
              <el-alert
                type="error"
                :closable="false"
              >
                {{ event.task.error }}
              </el-alert>
            </div>

            <div
              v-if="event.task?.duration"
              class="event-duration"
            >
              <el-icon><Clock /></el-icon>
              <span>耗时: {{ formatDuration(event.task.duration) }}</span>
            </div>
          </div>
        </el-card>
      </el-timeline-item>
    </el-timeline>
  </div>
</template>

<script setup lang="ts">
import { MagicStick, Document, Clock } from '@element-plus/icons-vue'
import { formatTime, formatRelativeTime, formatDuration } from '@/utils/format'

interface Props {
  events: any[]
}

defineProps<Props>()

const getEventType = (type: string) => {
  const types: Record<string, any> = {
    start: 'primary',
    complete: 'success',
    fail: 'danger',
    info: 'info'
  }
  return types[type] || 'info'
}

const getEventTagType = (type: string) => {
  const types: Record<string, any> = {
    start: 'primary',
    complete: 'success',
    fail: 'danger',
    info: 'info'
  }
  return types[type] || 'info'
}

const getEventText = (type: string) => {
  const texts: Record<string, string> = {
    start: '任务开始',
    complete: '任务完成',
    fail: '任务失败',
    info: '信息'
  }
  return texts[type] || type
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
</script>

<style scoped lang="scss">
.execution-timeline {
  .event-card {
    margin-bottom: 0;

    .event-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 12px;

      .event-info {
        display: flex;
        align-items: center;
        gap: 8px;

        .task-type {
          font-size: 14px;
          font-weight: 500;
        }
      }

      .event-time {
        font-size: 12px;
        color: #999;
      }
    }

    .event-body {
      display: flex;
      flex-direction: column;
      gap: 8px;

      .event-skill,
      .event-card,
      .event-duration {
        display: flex;
        align-items: center;
        gap: 6px;
        font-size: 13px;
        color: #666;

        .el-icon {
          font-size: 14px;
        }
      }

      .event-error {
        margin-top: 8px;
      }
    }
  }
}
</style>
