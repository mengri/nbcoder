<template>
  <el-menu
    :default-active="activeMenu"
    :router="true"
    background-color="#001529"
    text-color="#fff"
    active-text-color="#1890ff"
    class="sidebar-menu"
  >
    <el-menu-item index="/projects">
      <el-icon><FolderOpened /></el-icon>
      <span>项目管理</span>
    </el-menu-item>

    <el-menu-item
      v-if="currentProjectId"
      index="/projects/:id/cards"
    >
      <el-icon><Document /></el-icon>
      <span>需求卡片</span>
    </el-menu-item>

    <el-menu-item
      v-if="currentProjectId"
      index="/projects/:id/pipelines"
    >
      <el-icon><Operation /></el-icon>
      <span>Pipeline</span>
    </el-menu-item>

    <el-menu-item
      v-if="currentProjectId"
      index="/projects/:id/agents"
    >
      <el-icon><View /></el-icon>
      <span>Agent 任务</span>
    </el-menu-item>

    <el-menu-item
      v-if="currentProjectId"
      index="/projects/:id/monitor"
    >
      <el-icon><Monitor /></el-icon>
      <span>Agent 监控</span>
    </el-menu-item>

    <el-menu-item
      v-if="currentProjectId"
      index="/projects/:id/knowledge"
    >
      <el-icon><Reading /></el-icon>
      <span>知识库</span>
    </el-menu-item>

    <el-divider style="margin: 8px 0" />

    <el-menu-item index="/ai">
      <el-icon><MagicStick /></el-icon>
      <span>AI Runtime</span>
    </el-menu-item>

    <el-menu-item index="/ai/statistics">
      <el-icon><TrendCharts /></el-icon>
      <span>调用统计</span>
    </el-menu-item>
  </el-menu>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useProjectStore } from '@/stores/project'
import {
  FolderOpened,
  Document,
  Operation,
  View,
  Monitor,
  Reading,
  MagicStick,
  TrendCharts
} from '@element-plus/icons-vue'

const route = useRoute()
const projectStore = useProjectStore()

const activeMenu = computed(() => route.path)
const currentProjectId = computed(() => projectStore.currentProject?.id)
</script>

<style scoped lang="scss">
.sidebar-menu {
  height: 100%;
  border-right: none;

  .el-divider {
    border-color: #1f2937;
  }
}
</style>
