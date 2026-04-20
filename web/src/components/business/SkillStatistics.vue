<template>
  <div class="skill-statistics">
    <el-empty v-if="skills.length === 0" description="暂无统计数据" />
    <div v-else class="stats-list">
      <div
        v-for="(skill, index) in sortedSkills"
        :key="skill.name"
        class="stat-item"
      >
        <div class="stat-header">
          <div class="stat-name">
            <el-icon><Star /></el-icon>
            <span>{{ skill.name }}</span>
          </div>
          <div class="stat-count">
            <span class="count-value">{{ skill.count }}</span>
            <span class="count-label">次调用</span>
          </div>
        </div>

        <div class="stat-bar">
          <div
            class="bar-fill"
            :style="{ width: getPercentage(skill.count) + '%', backgroundColor: getBarColor(index) }"
          />
        </div>
      </div>
    </div>

    <div v-if="skills.length > 0" class="stats-summary">
      <el-statistic title="总调用次数" :value="totalCount" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Star } from '@element-plus/icons-vue'

interface SkillStats {
  name: string
  count: number
}

interface Props {
  skills: SkillStats[]
}

const props = defineProps<Props>()

const sortedSkills = computed(() => {
  return [...props.skills].sort((a, b) => b.count - a.count)
})

const totalCount = computed(() => {
  return props.skills.reduce((sum, skill) => sum + skill.count, 0)
})

const getPercentage = (count: number) => {
  if (totalCount.value === 0) return 0
  return (count / totalCount.value) * 100
}

const getBarColor = (index: number) => {
  const colors = [
    '#409eff',
    '#67c23a',
    '#e6a23c',
    '#f56c6c',
    '#909399',
    '#c71585',
    '#ff6347',
    '#00ced1'
  ]
  return colors[index % colors.length]
}
</script>

<style scoped lang="scss">
.skill-statistics {
  .stats-list {
    display: flex;
    flex-direction: column;
    gap: 16px;
    max-height: 400px;
    overflow-y: auto;
    margin-bottom: 20px;
  }

  .stat-item {
    .stat-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 8px;

      .stat-name {
        display: flex;
        align-items: center;
        gap: 6px;
        font-size: 14px;
        font-weight: 500;

        .el-icon {
          color: #e6a23c;
        }
      }

      .stat-count {
        display: flex;
        align-items: baseline;
        gap: 4px;

        .count-value {
          font-size: 20px;
          font-weight: bold;
          color: #409eff;
        }

        .count-label {
          font-size: 12px;
          color: #999;
        }
      }
    }

    .stat-bar {
      width: 100%;
      height: 8px;
      background-color: #f0f0f0;
      border-radius: 4px;
      overflow: hidden;

      .bar-fill {
        height: 100%;
        transition: width 0.3s ease;
        border-radius: 4px;
      }
    }
  }

  .stats-summary {
    padding-top: 16px;
    border-top: 1px solid #e8e8e8;
  }
}
</style>
