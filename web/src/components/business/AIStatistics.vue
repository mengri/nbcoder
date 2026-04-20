<template>
  <div class="ai-statistics">
    <el-row :gutter="16">
      <el-col :span="8">
        <el-card class="stat-card">
          <el-statistic title="总调用次数" :value="totalCalls" />
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card class="stat-card">
          <el-statistic title="总 Token 使用" :value="totalTokens" />
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card class="stat-card">
          <el-statistic
            title="总费用"
            :value="totalCost"
            :precision="4"
            prefix="$"
          />
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :span="12">
        <el-card class="chart-card">
          <template #header>
            <span>调用趋势</span>
          </template>
          <div ref="trendChartRef" style="width: 100%; height: 300px"></div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card class="chart-card">
          <template #header>
            <span>模型分布</span>
          </template>
          <div ref="modelChartRef" style="width: 100%; height: 300px"></div>
        </el-card>
      </el-col>
    </el-row>

    <el-card class="table-card" style="margin-top: 16px">
      <template #header>
        <span>模型使用详情</span>
      </template>
      <el-table :data="modelDistribution" style="width: 100%" max-height="300">
        <el-table-column prop="modelName" label="模型名称" />
        <el-table-column prop="calls" label="调用次数" width="100" />
        <el-table-column prop="tokens" label="Token 使用" width="120" />
        <el-table-column prop="cost" label="费用" width="100">
          <template #default="{ row }">
            {{ formatCost(row.cost) }}
          </template>
        </el-table-column>
        <el-table-column label="占比" width="80">
          <template #default="{ row }">
            {{ formatPercentage(row.cost, totalCost) }}
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import * as echarts from 'echarts'
import { formatCost, formatPercentage } from '@/utils/format'
import type { ModelStats, DailyStats } from '@/types/ai'

interface Props {
  totalCalls: number
  totalTokens: number
  totalCost: number
  modelDistribution: ModelStats[]
  dailyTrend: DailyStats[]
}

const props = defineProps<Props>()

const trendChartRef = ref<HTMLElement | null>(null)
const modelChartRef = ref<HTMLElement | null>(null)

let trendChart: echarts.ECharts | null = null
let modelChart: echarts.ECharts | null = null

const totalCost = computed(() => props.totalCost)

const initCharts = () => {
  initTrendChart()
  initModelChart()
}

const initTrendChart = () => {
  if (!trendChartRef.value) return

  trendChart = echarts.init(trendChartRef.value)

  const dates = props.dailyTrend.map(item => item.date)
  const calls = props.dailyTrend.map(item => item.calls)
  const tokens = props.dailyTrend.map(item => item.tokens)

  const option = {
    tooltip: {
      trigger: 'axis'
    },
    legend: {
      data: ['调用次数', 'Token 使用']
    },
    xAxis: {
      type: 'category',
      data: dates
    },
    yAxis: [
      {
        type: 'value',
        name: '次数'
      },
      {
        type: 'value',
        name: 'Token'
      }
    ],
    series: [
      {
        name: '调用次数',
        type: 'line',
        data: calls,
        smooth: true
      },
      {
        name: 'Token 使用',
        type: 'bar',
        yAxisIndex: 1,
        data: tokens
      }
    ]
  }

  trendChart.setOption(option)
}

const initModelChart = () => {
  if (!modelChartRef.value) return

  modelChart = echarts.init(modelChartRef.value)

  const data = props.modelDistribution.map(item => ({
    name: item.modelName,
    value: item.cost
  }))

  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{b}: ${c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      left: 'left'
    },
    series: [
      {
        name: '费用分布',
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 10,
          borderColor: '#fff',
          borderWidth: 2
        },
        label: {
          show: false,
          position: 'center'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: 20,
            fontWeight: 'bold'
          }
        },
        labelLine: {
          show: false
        },
        data: data
      }
    ]
  }

  modelChart.setOption(option)
}

const handleResize = () => {
  trendChart?.resize()
  modelChart?.resize()
}

onMounted(() => {
  initCharts()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  trendChart?.dispose()
  modelChart?.dispose()
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped lang="scss">
.ai-statistics {
  .stat-card {
    height: 100px;
  }

  .chart-card {
    min-height: 380px;
  }

  .table-card {
    min-height: 400px;
  }
}
</style>
