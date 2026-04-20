<template>
  <div class="pipeline-timeline">
    <el-timeline>
      <el-timeline-item
        v-for="stage in stages"
        :key="stage.id"
        :timestamp="formatTimestamp(stage)"
        :type="getTimelineType(stage.status)"
        :hollow="stage.status === 'PENDING'"
        placement="top"
      >
        <el-card class="stage-card">
          <div class="stage-header">
            <div class="stage-info">
              <span class="stage-name">{{ stage.name }}</span>
              <el-tag
                :type="getStatusType(stage.status)"
                size="small"
              >
                {{ getStatusText(stage.status) }}
              </el-tag>
            </div>
            <div class="stage-duration">
              {{ formatDuration(stage.startTime, stage.endTime) }}
            </div>
          </div>

          <div
            v-if="stage.steps.length > 0"
            class="steps-container"
          >
            <div
              v-for="step in stage.steps"
              :key="step.id"
              class="step-item"
            >
              <div class="step-icon">
                <el-icon :class="getStepIconClass(step.status)">
                  <component :is="getStepIcon(step.status)" />
                </el-icon>
              </div>
              <div class="step-content">
                <div class="step-header">
                  <span class="step-name">{{ step.name }}</span>
                  <el-tag
                    v-if="step.status"
                    :type="getStatusType(step.status)"
                    size="small"
                  >
                    {{ getStatusText(step.status) }}
                  </el-tag>
                </div>
                <div class="step-time">
                  {{ formatStepTime(step) }}
                </div>
                <div
                  v-if="step.error"
                  class="step-error"
                >
                  <el-alert
                    type="error"
                    :closable="false"
                  >
                    {{ step.error }}
                  </el-alert>
                </div>
                <div
                  v-if="step.output"
                  class="step-output"
                >
                  <el-collapse>
                    <el-collapse-item
                      title="输出"
                      name="output"
                    >
                      <pre class="code-block">{{ formatJson(step.output) }}</pre>
                    </el-collapse-item>
                  </el-collapse>
                </div>
              </div>
            </div>
          </div>
        </el-card>
      </el-timeline-item>
    </el-timeline>
  </div>
</template>

<script setup lang="ts">
import { CircleCheck, CircleClose, Clock, Loading } from '@element-plus/icons-vue'
import { formatDuration, formatJson, formatDate } from '@/utils/format'
import type { PipelineStage } from '@/types/pipeline'

interface Props {
  stages: PipelineStage[]
}

defineProps<Props>()

const formatTimestamp = (stage: PipelineStage) => {
  if (stage.startTime && stage.endTime) {
    return `${formatDate(stage.startTime)} - ${formatDate(stage.endTime)}`
  }
  if (stage.startTime) {
    return `${formatDate(stage.startTime)} - 进行中`
  }
  return ''
}

const getTimelineType = (status: string) => {
  const types: Record<string, any> = {
    COMPLETED: 'primary',
    FAILED: 'danger',
    RUNNING: 'warning',
    PENDING: 'info',
    REVIEW: 'warning',
    SKIPPED: 'info'
  }
  return types[status] || 'info'
}

const getStatusType = (status: string) => {
  const types: Record<string, any> = {
    COMPLETED: 'success',
    FAILED: 'danger',
    RUNNING: 'primary',
    PENDING: 'info',
    REVIEW: 'warning',
    SKIPPED: 'info'
  }
  return types[status] || 'info'
}

const getStatusText = (status: string) => {
  const texts: Record<string, string> = {
    PENDING: '等待中',
    RUNNING: '执行中',
    REVIEW: '待审核',
    COMPLETED: '已完成',
    FAILED: '失败',
    SKIPPED: '已跳过'
  }
  return texts[status] || status
}

const getStepIcon = (status: string) => {
  const icons: Record<string, any> = {
    COMPLETED: CircleCheck,
    FAILED: CircleClose,
    RUNNING: Loading,
    PENDING: Clock
  }
  return icons[status] || Clock
}

const getStepIconClass = (status: string) => {
  const classes: Record<string, string> = {
    COMPLETED: 'step-success',
    FAILED: 'step-error',
    RUNNING: 'step-running',
    PENDING: 'step-pending'
  }
  return classes[status] || 'step-pending'
}

const formatStepTime = (step: any) => {
  if (step.startTime && step.endTime) {
    return `${formatDate(step.startTime)} - ${formatDate(step.endTime)} (${formatDuration(step.startTime, step.endTime)})`
  }
  if (step.startTime) {
    return `${formatDate(step.startTime)} - 进行中`
  }
  return '等待中'
}
</script>

<style scoped lang="scss">
.pipeline-timeline {
  .stage-card {
    margin-bottom: 0;

    .stage-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 16px;

      .stage-info {
        display: flex;
        align-items: center;
        gap: 12px;

        .stage-name {
          font-size: 16px;
          font-weight: 500;
        }
      }

      .stage-duration {
        font-size: 14px;
        color: #666;
      }
    }

    .steps-container {
      .step-item {
        display: flex;
        gap: 12px;
        padding: 12px 0;
        border-bottom: 1px solid #f0f0f0;

        &:last-child {
          border-bottom: none;
        }

        .step-icon {
          width: 32px;
          height: 32px;
          display: flex;
          align-items: center;
          justify-content: center;
          border-radius: 50%;
          background-color: #f5f7fa;
          flex-shrink: 0;

          .el-icon {
            font-size: 16px;

            &.step-success {
              color: #67c23a;
            }

            &.step-error {
              color: #f56c6c;
            }

            &.step-running {
              color: #409eff;
            }

            &.step-pending {
              color: #909399;
            }
          }
        }

        .step-content {
          flex: 1;

          .step-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 8px;

            .step-name {
              font-weight: 500;
            }
          }

          .step-time {
            font-size: 13px;
            color: #666;
            margin-bottom: 8px;
          }

          .step-error {
            margin-top: 8px;
          }

          .step-output {
            margin-top: 8px;

            .code-block {
              background-color: #f5f7fa;
              padding: 12px;
              border-radius: 4px;
              font-size: 13px;
              line-height: 1.6;
              max-height: 300px;
              overflow-y: auto;
            }
          }
        }
      }
    }
  }
}
</style>
