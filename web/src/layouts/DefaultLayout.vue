<template>
  <el-container class="default-layout">
    <el-aside width="240px">
      <Sidebar />
    </el-aside>
    <el-container>
      <el-header>
        <Header />
      </el-header>
      <el-main>
        <router-view />
      </el-main>
      <el-footer>
        <Footer />
      </el-footer>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import Header from '@/components/layout/Header.vue'
import Sidebar from '@/components/layout/Sidebar.vue'
import Footer from '@/components/layout/Footer.vue'
</script>

<style scoped lang="scss">
.default-layout {
  height: 100vh;

  .el-aside {
    background-color: #001529;
    color: #fff;
  }

  .el-header {
    background-color: #fff;
    border-bottom: 1px solid #e8e8e8;
    padding: 0 20px;
    height: 60px;
    line-height: 60px;
  }

  .el-main {
    background-color: #f5f7fa;
    padding: 20px;
  }

  .el-footer {
    background-color: #fff;
    border-top: 1px solid #e8e8e8;
    text-align: center;
    line-height: 40px;
    color: #999;
    font-size: 12px;
  }
}
</style>
