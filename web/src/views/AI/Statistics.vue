<template>
  <div class="ai-statistics">
    <el-page-header @back="goBack">
      <template #content>
        <span class="title">AI 调用统计</span>
      </template>
    </el-page-header>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="6">
        <el-card>
          <el-statistic title="总调用次数" :value="statistics.totalCalls">
            <template #prefix>
              <el-icon style="vertical-align: -0.125em">
                <DataLine />
              </el-icon>
            </template>
          </el-statistic>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <el-statistic title="总 Token 使用" :value="statistics.totalTokens">
            <template #prefix>
              <el-icon style="vertical-align: -0.125em">
                <Coin />
              </el-icon>
            </template>
          </el-statistic>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <el-statistic title="总费用" :value="statistics.totalCost" :precision="4" prefix="$">
            <template #prefix>
              <el-icon style="vertical-align: -0.125em">
                <Wallet />
              </el-icon>
            </template>
          </el-statistic>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <el-statistic title="今日调用" :value="statistics.todayCalls">
            <template #prefix>
              <el-icon style="vertical-align: -0.125em">
                <Calendar />
              </el-icon>
            </template>
          </el-statistic>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="16">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>调用趋势</span>
              <el-radio-group v-model="trendPeriod" size="small">
                <el-radio-button label="7">7天</el-radio-button>
                <el-radio-button label="30">30天</el-radio-button>
                <el-radio-button label="90">90天</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <div ref="trendChartRef" style="width: 100%; height: 400px"></div>
        </el-card>
      </el-col>

      <el-col :span="8">
        <el-card>
          <template #header>
            <span>模型使用分布</span>
          </template>
          <div ref="modelChartRef" style="width: 100%; height: 400px"></div>
        </el-card>
      </el-col>
    </el-row>

    <el-card style="margin-top: 20px">
      <template #header>
        <span>模型使用详情</span>
      </template>
      <el-table :data="statistics.modelDistribution" style="width: 100%">
        <el-table-column prop="modelName" label="模型名称" min-width="200" />
        <el-table-column prop="calls" label="调用次数" width="120" />
        <el-table-column prop="tokens" label="Token 使用" width="150" />
        <el-table-column prop="cost" label="费用" width="120">
          <template #default="{ row }">
            {{ formatCost(row.cost) }}
          </template>
        </el-table-column>
        <el-table-column label="费用占比" width="120">
          <template #default="{ row }">
            {{ formatPercentage(row.cost, statistics.totalCost) }}
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { DataLine, Coin, Wallet, Calendar } from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import { formatCost, formatPercentage } from '@/utils/format'
import type { Statistics } from '@/types/ai'

const router = useRouter()

const trendPeriod = ref('7')
const trendChartRef = ref<HTMLElement | null>(null)
const modelChartRef = ref<HTMLElement | null>(null)

let trendChart: echarts.ECharts | null = null
let modelChart: echarts.ECharts | null = null

const statistics = ref<Statistics>({
  totalCalls: 0,
  totalTokens: 0,
  totalCost: 0,
  todayCalls: 0,
  modelDistribution: [],
  dailyTrend: []
})

const loadStatistics = async () => {
  try {
    const data = await fetch(`/api/v1/ai/statistics?period=${trendPeriod.value}`)
    const result = await data.json()
    statistics.value = result
    updateCharts()
  } catch (error) {
    ElMessage.error('加载统计数据失败')
  }
}

const updateCharts = () => {
  updateTrendChart()
  updateModelChart()
}

const updateTrendChart = () => {
  if (!trendChartRef.value) return

  if (!trendChart) {
    trendChart = echarts.init(trendChartRef.value)
  }

  const dates = statistics.value.dailyTrend.map(item => item.date)
  const calls = statistics.value.dailyTrend.map(item => item.calls)
  const tokens = statistics.value.dailyTrend.map(item => item.tokens)
  const costs = statistics.value.dailyTrend.map(item => item.cost)

  const option = {
    tooltip: {
      trigger: 'axis'
    },
    legend: {
      data: ['调用次数', 'Token 使用', '费用']
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
      },
      {
        name: '费用',
        type: 'line',
        yAxisIndex: 1,
        data: costs,
        smooth: true
      }
    ]
  }

  trendChart.setOption(option)
}

const updateModelChart = () => {
  if (!modelChartRef.value) return

  if (!modelChart) {
    modelChart = echarts.init(modelChartRef.value)
  }

  const data = statistics.value.modelDistribution.map(item => ({
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
        data: data,
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)'
          }
        }
      }
    ]
  }

  modelChart.setOption(option)
}

const handleResize = () => {
  trendChart?.resize()
  modelChart?.resize()
}

watch(trendPeriod, () => {
  loadStatistics()
})

const goBack = () => {
  router.push('/ai')
}

onMounted(() => {
  loadStatistics()
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
