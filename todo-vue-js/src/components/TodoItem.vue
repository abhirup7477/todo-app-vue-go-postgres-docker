<template>
  <el-card class="todo-card">
    <div class="todo-row">
      
      <!-- LEFT: TITLE -->
      <div class="todo-title" @click="open">
        {{ todo.title }}
      </div>

      <!-- MIDDLE: STATUS -->
      <div class="todo-status">
        <el-text v-if="todo.completed" type="success">Done</el-text>
        <el-text v-else type="warning">Incomplete</el-text>
      </div>

      <!-- RIGHT: ACTIONS -->
      <div class="todo-actions">
        <el-button type="primary" :icon="Edit" circle @click="editTodo"/>
        <el-button type="danger" :icon="Delete" circle @click="remove"/>
      </div>

    </div>
  </el-card>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useTodoStore } from '../store/todoStore'
import { Edit, Delete } from '@element-plus/icons-vue'
import { ElMessageBox, ElMessage } from 'element-plus'

const props = defineProps({ todo: Object })
const router = useRouter()
const store = useTodoStore()

const open = () => router.push(`/todo/${props.todo.id}`)
const remove = async () => {
  try {
    await ElMessageBox.confirm(
      'This will permanently delete the todo. Continue?',
      'Warning',
      {
        confirmButtonText: 'Delete',
        cancelButtonText: 'Cancel',
        type: 'warning',
      }
    )

    store.removeById(props.todo.id)
    ElMessage.success('Todo deleted')

  } catch {
    ElMessage.info('Deletion cancelled')
  }
}
const editTodo = () => router.push(`/edit/${props.todo.id}`)
</script>

<style scoped>
.todo-card {
  margin-bottom: 12px;
}

.todo-row {
  display: flex;
  align-items: center;
  gap: 16px;
}

/* Title takes remaining space */
.todo-title {
  flex: 1;
  cursor: pointer;

  /* wrap long text */
  word-break: break-word;
  overflow-wrap: anywhere;
}

/* Fixed width center alignment */
.todo-status {
  width: 100px;
  text-align: center;
}

/* Buttons always right */
.todo-actions {
  display: flex;
  gap: 8px;
}

/* Mobile responsiveness */
@media (max-width: 600px) {
  .todo-row {
    flex-direction: column;
    align-items: stretch;
  }

  .todo-status {
    width: 100%;
    text-align: left;
  }

  .todo-actions {
    justify-content: flex-end;
  }
}
</style>