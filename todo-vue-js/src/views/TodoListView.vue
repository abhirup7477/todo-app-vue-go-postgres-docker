<template>
  <div>
    <div style="display:flex; justify-content:space-between; margin-bottom:15px">
      <h2>All Todos</h2>
      <el-button type="primary" @click="goAdd">Add Todo</el-button>
    </div>

    <TodoItem v-for="t in store.todos" :key="t.id" :todo="t" />
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useTodoStore } from '../store/todoStore'
import TodoItem from '../components/TodoItem.vue'

const store = useTodoStore()
const router = useRouter()
const goAdd = () => router.push('/add')
</script>