<template>
  <div class="dependency-graph">
    <el-empty v-if="!cards || cards.length === 0" description="暂无数据" />
    <div v-else class="graph-container">
      <div
        v-for="card in cards"
        :key="card.id"
        class="card-node"
        :class="{ 'is-selected': selectedCardId === card.id }"
        @click="handleCardClick(card)"
      >
        <div class="card-header">
          <el-tag :type="getStatusType(card.status)" size="small">
            {{ getStatusText(card.status) }}
          </el-tag>
          <el-tag :type="getPriorityType(card.priority)" size="small">
            {{ getPriorityText(card.priority) }}
          </el-tag>
        </div>
        <div class="card-title">{{ card.title }}</div>
        <div class="card-id">{{ card.id }}</div>

        <div v-if="card.dependencies.length > 0" class="dependencies">
          <div class="dep-label">依赖:</div>
          <el-tag
            v-for="depId in card.dependencies"
            :key="depId"
            size="small"
            class="dep-tag"
          >
            {{ getCardTitle(depId) }}
          </el-tag>
        </div>
      </div>

      <svg class="connections" v-if="cards.length > 0">
        <line
          v-for="(connection, index) in connections"
          :key="index"
          :x1="connection.x1"
          :y1="connection.y1"
          :x2="connection.x2"
          :y2="connection.y2"
          stroke="#409eff"
          stroke-width="2"
          marker-end="url(#arrowhead)"
        />
        <defs>
          <marker
            id="arrowhead"
            markerWidth="10"
            markerHeight="7"
            refX="9"
            refY="3.5"
            orient="auto"
          >
            <polygon points="0 0, 10 3.5, 0 7" fill="#409eff" />
          </marker>
        </defs>
      </svg>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import type { Card } from '@/types/card'

interface Props {
  cards: Card[]
  selectedCardId?: string
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'card-click', card: Card): void
}>()

const graphContainer = ref<HTMLElement | null>(null)
const cardPositions = ref<Map<string, { x: number; y: number }>>(new Map())

const connections = computed(() => {
  const result: any[] = []
  props.cards.forEach(card => {
    card.dependencies.forEach(depId => {
      const fromPos = cardPositions.value.get(depId)
      const toPos = cardPositions.value.get(card.id)
      if (fromPos && toPos) {
        result.push({
          x1: fromPos.x,
          y1: fromPos.y,
          x2: toPos.x,
          y2: toPos.y
        })
      }
    })
  })
  return result
})

const handleCardClick = (card: Card) => {
  emit('card-click', card)
}

const getCardTitle = (cardId: string) => {
  const card = props.cards.find(c => c.id === cardId)
  return card ? card.title : cardId
}

const getStatusType = (status: string) => {
  const types: Record<string, any> = {
    DRAFT: 'info',
    PENDING: 'warning',
    CONFIRMED: 'success',
    IN_PROGRESS: 'primary',
    COMPLETED: 'success',
    SUPPLANTED: 'info',
    ABANDONED: 'danger'
  }
  return types[status] || 'info'
}

const getStatusText = (status: string) => {
  const texts: Record<string, string> = {
    DRAFT: '草稿',
    PENDING: '待确认',
    CONFIRMED: '已确认',
    IN_PROGRESS: '进行中',
    COMPLETED: '已完成',
    SUPPLANTED: '被取代',
    ABANDONED: '已废弃'
  }
  return texts[status] || status
}

const getPriorityType = (priority: string) => {
  const types: Record<string, any> = {
    LOW: 'info',
    MEDIUM: '',
    HIGH: 'warning',
    CRITICAL: 'danger'
  }
  return types[priority] || ''
}

const getPriorityText = (priority: string) => {
  const texts: Record<string, string> = {
    LOW: '低',
    MEDIUM: '中',
    HIGH: '高',
    CRITICAL: '紧急'
  }
  return texts[priority] || priority
}

const updatePositions = () => {
  if (!graphContainer.value) return

  const containerRect = graphContainer.value.getBoundingClientRect()
  cardPositions.value.clear()

  const cardElements = graphContainer.value.querySelectorAll('.card-node')
  cardElements.forEach((el, index) => {
    const rect = el.getBoundingClientRect()
    const cardId = props.cards[index]?.id
    if (cardId) {
      cardPositions.value.set(cardId, {
        x: rect.left - containerRect.left + rect.width / 2,
        y: rect.top - containerRect.top
      })
    }
  })
}

const handleResize = () => {
  nextTick(updatePositions)
}

onMounted(() => {
  window.addEventListener('resize', handleResize)
  nextTick(updatePositions)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped lang="scss">
.dependency-graph {
  width: 100%;
  min-height: 400px;
  position: relative;

  .graph-container {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 20px;
    position: relative;
  }

  .card-node {
    background: #fff;
    border: 1px solid #e8e8e8;
    border-radius: 8px;
    padding: 16px;
    cursor: pointer;
    transition: all 0.3s;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
    position: relative;
    z-index: 2;

    &:hover {
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
      transform: translateY(-2px);
    }

    &.is-selected {
      border-color: #409eff;
      box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.2);
    }

    .card-header {
      display: flex;
      gap: 8px;
      margin-bottom: 12px;
    }

    .card-title {
      font-size: 16px;
      font-weight: 500;
      margin-bottom: 8px;
      color: #333;
    }

    .card-id {
      font-size: 12px;
      color: #999;
      margin-bottom: 12px;
    }

    .dependencies {
      .dep-label {
        font-size: 12px;
        color: #666;
        margin-bottom: 8px;
      }

      .dep-tag {
        margin: 4px;
      }
    }
  }

  .connections {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    pointer-events: none;
    z-index: 1;
  }
}
</style>
