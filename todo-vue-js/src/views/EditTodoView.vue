<template>
  <!-- show form ONLY when todo exists -->
  <el-card v-if="todo">
    <h2>Edit Todo</h2>

    <el-form label-position="top">
      <el-form-item label="Title">
        <el-input v-model.trim="title" />
      </el-form-item>

      <el-form-item label="Description">
        <el-input v-model.trim="description" />
      </el-form-item>

      <el-button type="primary" @click="update">Update</el-button>
      <el-button @click="back">Cancel</el-button>
    </el-form>
  </el-card>

  <!-- loading fallback -->
  <el-skeleton v-else animated />
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useTodoStore } from '@/store/todoStore'

const route = useRoute()
const router = useRouter()
const store = useTodoStore()

const todo = ref(null)          // 🔑 THIS WAS MISSING
const title = ref('')
const description = ref('')

onMounted(async () => {
  try {
    const data = await store.getById(route.params.id)
    todo.value = data            // 🔑 assign fetched data
    title.value = data.title
    description.value = data.description
  } catch (err) {
    router.push('/')             // safety fallback
  }
})

const update = async () => {
  if (!todo.value || !title.value.trim()) return

  await store.update({
    id: todo.value.id,
    title: title.value,
    description: description.value,
    completed: todo.value.completed
  })

  router.push(`/todo/${todo.value.id}`)
}

const back = () => router.push('/')
</script>