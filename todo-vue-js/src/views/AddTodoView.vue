<template>
  <el-card>
    <h2>Add Todo</h2>

    <el-form label-position="top">
      <el-form-item label="Title">
        <el-input v-model.trim="title" />
      </el-form-item>

      <el-form-item label="Description">
        <el-input v-model.trim="description" />
      </el-form-item>

      <el-button type="primary" @click="add">Add</el-button>
      <el-button @click="back">Cancel</el-button>
    </el-form>
  </el-card>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useTodoStore } from '../store/todoStore'

const title = ref('')
const description = ref('')

const store = useTodoStore()
const router = useRouter()

const add = () => {
  if (!title.value) return
  if(store.addTask(title.value, description.value)){
    router.push('/')
  }
}

const back = () => router.push('/')
</script>