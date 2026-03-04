<template>
  <el-container class="app-layout">
    
    <el-header class="app-header">
      <div class="header-content">
        Todo Application
      </div>
    </el-header>

    <el-main class="app-main">
      <div class="content-wrapper">
        <router-view />
      </div>
    </el-main>

  </el-container>
</template>

<script setup>
import { onMounted } from 'vue'
import { useTodoStore } from '@/store/todoStore';

const store = useTodoStore()

onMounted(() => {
  store.fetchTodos()
})
</script>

<style scoped>
.app-layout {
  min-height: 100vh;
}

.app-header {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  font-weight: 600;
  border-bottom: 1px solid #eaeaea;
  background: #fff;
}

.app-main {
  display: flex;
  justify-content: center;   
  align-items: flex-start;  
  padding: 24px 12px;
  background: #f5f7fa;
}

.content-wrapper {
  width: 100%;
  max-width: 700px;        
}
</style>