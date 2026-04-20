<template>
  <div class="pipeline-view">
    <div class="pipeline-container" v-if="stages.length > 0">
      <div
        v-for="(stage, index) in stages"
        :key="stage.id"
        class="stage-node"
        :class="getStageClass(stage)"
        @click="$emit('stage-click', stage)"
      >
        <div class="stage-header">
          <div class="stage-icon">
            <el-icon>
              <component :is="getStageIcon(stage.status)" />
            </el-icon>
          </div>
          <div class="stage-info">
            <div class="stage-name">{{ stage.name }}</div>
            <div class="stage-status">{{ getStatusText(stage.status) }}</div>
          </div>
        </div>

        <div class="stage-body">
          <div class="stage-meta">
            <span v-if="stage.startTime">
              开始: {{ formatDate(stage.startTime) }}
            </span>
            <span v-if="stage.endTime">
              结束: {{ formatDate(stage.endTime) }}
            </span>
            <span v-if="stage.startTime && stage.endTime">
              耗时: {{ formatDuration(stage.startTime, stage.endTime) }}
            </span>
          </div>

          <div v-if="stage.steps.length > 0" class="steps-list">
            <div
              v-for="(step, stepIndex) in stage.steps"
              :key="step.id"
              class="step-item"
              :class="getStepClass(step.status)"
            >
              <div class="step-icon">
                <el-icon>
                  <component :is="getStepIcon(step.status)" />
                </el-icon>
              </div>
              <div class="step-name">{{ step.name }}</div>
            </div>
          </div>
        </div>

        <div v-if="index < stages.length - 1" class="stage-connector">
          <div class="connector-line"></div>
          <div class="connector-arrow">
            <el-icon><ArrowDown /></el-icon>
          </div>
        </div>
      </div>
    </div>

    <el-empty v-else description="暂无 Pipeline 数据" />
  </div>
</template>

<script setup lang="ts">
import {
  CircleCheck,
  CircleClose,
  Clock,
  Loading,
  Warning,
  ArrowDown
} from '@element-plus/icons-vue'
import { formatDate, formatDuration } from '@/utils/format'
import type { PipelineStage } from '@/types/pipeline'

interface Props {
  stages: PipelineStage[]
}

defineProps<Props>()

defineEmits<{
  (e: 'stage-click', stage: PipelineStage): void
}>()

const getStageClass = (stage: PipelineStage) => {
  return {
    'is-completed': stage.status === 'COMPLETED',
    'is-running': stage.status === 'RUNNING',
    'is-failed': stage.status === 'FAILED',
    'is-review': stage.status === 'REVIEW',
    'is-pending': stage.status === 'PENDING',
    'is-skipped': stage.status === 'SKIPPED'
  }
}

const getStageIcon = (status: string) => {
  const icons: Record<string, any> = {
    COMPLETED: CircleCheck,
    FAILED: CircleClose,
    RUNNING: Loading,
    REVIEW: Warning,
    PENDING: Clock,
    SKIPPED: Clock
  }
  return icons[status] || Clock
}

const getStepClass = (status: string) => {
  return {
    'step-completed': status === 'COMPLETED',
    'step-running': status === 'RUNNING',
    'step-failed': status === 'FAILED',
    'step-pending': status === 'PENDING'
  }
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
</script>

<style scoped lang="scss">
.pipeline-view {
  width: 100%;
  min-height: 400px;
  padding: 20px;

  .pipeline-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0;
  }

  .stage-node {
    width: 100%;
    max-width: 600px;
    background: #fff;
    border: 2px solid #e8e8e8;
    border-radius: 8px;
    padding: 16px;
    cursor: pointer;
    transition: all 0.3s;
    position: relative;

    &:hover {
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    }

    &.is-completed {
      border-color: #67c23a;
      background-color: #f0f9ff;
    }

    &.is-running {
      border-color: #409eff;
      background-color: #ecf5ff;
      animation: pulse 2s infinite;
    }

    &.is-failed {
      border-color: #f56c6c;
      background-color: #fef0f0;
    }

    &.is-review {
      border-color: #e6a23c;
      background-color: #fdf6ec;
    }

    &.is-pending {
      border-color: #dcdfe6;
      background-color: #f5f7fa;
    }

    &.is-skipped {
      border-color: #909399;
      background-color: #f4f4f5;
      opacity: 0.6;
    }

    .stage-header {
      display: flex;
      align-items: flex-start;
      gap: 12px;
      margin-bottom: 12px;

      .stage-icon {
        width: 40px;
        height: 40px;
        display: flex;
        align-items: center;
        justify-content: center;
        border-radius: 50%;
        background-color: #fff;
        flex-shrink: 0;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);

        .el-icon {
          font-size: 20px;

          &.is-completed {
            color: #67c23a;
          }

          &.is-failed {
            color: #f56c6c;
          }

          &.is-running {
            color: #409eff;
            animation: rotate 1s linear infinite;
          }

          &.is-review {
            color: #e6a23c;
          }

          &.is-pending {
            color: #909399;
          }
        }
      }

      .stage-info {
        flex: 1;

        .stage-name {
          font-size: 16px;
          font-weight: 500;
          margin-bottom: 4px;
        }

        .stage-status {
          font-size: 14px;
          color: #666;
        }
      }
    }

    .stage-body {
      .stage-meta {
        display: flex;
        flex-wrap: wrap;
        gap: 16px;
        font-size: 13px;
        color: #666;
        margin-bottom: 12px;
        padding: 8px 12px;
        background-color: rgba(255, 255, 255, 0.6);
        border-radius: 4px;
      }

      .steps-list {
        display: flex;
        flex-wrap: wrap;
        gap: 8px;

        .step-item {
          display: flex;
          align-items: center;
          gap: 6px;
          padding: 6px 12px;
          background-color: rgba(255, 255, 255, 0.6);
          border-radius: 16px;
          font-size: 13px;

          &.step-completed {
            color: #67c23a;
            background-color: rgba(103, 194, 58, 0.1);
          }

          &.step-running {
            color: #409eff;
            background-color: rgba(64, 158, 255, 0.1);
          }

          &.step-failed {
            color: #f56c6c;
            background-color: rgba(245, 108, 108, 0.1);
          }

          &.step-pending {
            color: #909399;
          }

          .step-icon {
            font-size: 14px;
          }

          .step-name {
            font-size: 13px;
          }
        }
      }
    }

    .stage-connector {
      position: absolute;
      bottom: -40px;
      left: 50%;
      transform: translateX(-50%);
      width: 2px;
      height: 40px;

      .connector-line {
        width: 100%;
        height: 100%;
        background-color: #dcdfe6;
      }

      .connector-arrow {
        position: absolute;
        bottom: -12px;
        left: 50%;
        transform: translateX(-50%);
        width: 24px;
        height: 24px;
        display: flex;
        align-items: center;
        justify-content: center;
        background-color: #fff;
        border-radius: 50%;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);

        .el-icon {
          color: #dcdfe6;
        }
      }
    }

    &:last-child {
      .stage-connector {
        display: none;
      }
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

@keyframes rotate {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>
