<template>
  <div class="header">
    <div class="header-left">
      <h1 class="logo">NBCoder</h1>
      <el-tag v-if="currentProject" type="success" size="small">
        {{ currentProject.name }}
      </el-tag>
    </div>
    <div class="header-right">
      <el-dropdown @command="handleCommand">
        <span class="user-info">
          <el-icon><User /></el-icon>
          {{ userName }}
          <el-icon class="el-icon--right"><arrow-down /></el-icon>
        </span>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="settings">设置</el-dropdown-item>
            <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useProjectStore } from '@/stores/project'
import { User, ArrowDown } from '@element-plus/icons-vue'

const router = useRouter()
const projectStore = useProjectStore()
const userName = ref('Admin')
const currentProject = computed(() => projectStore.currentProject)

const handleCommand = (command: string) => {
  if (command === 'logout') {
    localStorage.removeItem('token')
    router.push('/login')
  }
}

onMounted(() => {
  projectStore.loadProjects()
})
</script>

<style scoped lang="scss">
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 100%;

  .header-left {
    display: flex;
    align-items: center;
    gap: 12px;

    .logo {
      font-size: 24px;
      font-weight: bold;
      color: #1890ff;
      margin: 0;
    }
  }

  .header-right {
    .user-info {
      cursor: pointer;
      display: flex;
      align-items: center;
      gap: 4px;

      &:hover {
        color: #1890ff;
      }
    }
  }
}
</style>
