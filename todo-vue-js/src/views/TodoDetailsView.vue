<template>
  <div class="page-wrapper" v-if="todo">
    <el-card class="todo-card">

      <!-- TITLE -->
      <h2 class="todo-title">
        {{ todo.title }}
      </h2>

      <!-- DESCRIPTION -->
      <div class="todo-description">
        {{ todo.description }}
      </div>

      <!-- STATUS -->
      <el-checkbox
        :model-value="todo.completed"
        @update:model-value="toggle"
      >
        Completed
      </el-checkbox>

      <!-- ACTIONS -->
      <div class="actions">
        <el-button type="primary" @click="back">
          Close
        </el-button>

        <el-button type="danger" @click="remove">
          Delete
        </el-button>
      </div>

    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useTodoStore } from '@/store/todoStore'
import { ElMessageBox, ElMessage } from 'element-plus'

const API_BASE = 'http://localhost:8080/todos'  // ✅ define API_BASE

const route = useRoute()
const router = useRouter()
const store = useTodoStore()

const todo = ref(null)

// Load todo by ID
onMounted(async () => {
  try {
    todo.value = await store.getById(route.params.id)
  } catch {
    router.push('/') // invalid ID / deleted todo
  }
})

const back = () => router.push('/')

// Delete todo
const remove = async () => {
  if (!todo.value) return

  try {
    await ElMessageBox.confirm(
      'This will permanently delete the todo. Continue?',
      'Warning',
      {
        confirmButtonText: 'Delete',
        cancelButtonText: 'Cancel',
        type: 'warning'
      }
    )

    await store.removeById(todo.value.id)
    ElMessage.success('Todo deleted')
    router.push('/')

  } catch {
    ElMessage.info('Deletion cancelled')
  }
}

// Toggle completed status
const toggle = async () => {
  if (!todo.value) return

  try {
    const res = await fetch(`${API_BASE}/toggle/${todo.value.id}`, {
      method: 'PATCH',
    })

    if (!res.ok) {
      const err = await res.json()
      ElMessage.error(err.error || 'Failed to toggle status')
      return
    }

    const updated = await res.json()
    todo.value.completed = updated.completed

    // sync store
    const index = store.todos.findIndex(t => t.id === todo.value.id)
    if (index !== -1) store.todos[index].completed = updated.completed

  } catch (err) {
    ElMessage.error('Failed to toggle status')
  }
}
</script>

<style scoped>
.page-wrapper {
  display: flex;
  justify-content: center;
  padding: 24px 12px;
}
.todo-card {
  width: 100%;
  max-width: 700px;
}
:deep(.el-card__body) {
  overflow: hidden;
}
.todo-title {
  margin: 0 0 10px 0;
  line-height: 1.4;
  white-space: normal;
  word-break: break-word;
  overflow-wrap: anywhere;
}
.todo-description {
  margin: 12px 0 16px 0;
  white-space: normal;
  word-break: break-word;
  overflow-wrap: anywhere;
}
.actions {
  margin-top: 16px;
  display: flex;
  gap: 12px;
}
</style>