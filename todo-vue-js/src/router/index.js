import { createRouter, createWebHistory } from 'vue-router'
import TodoListView from '../views/TodoListView.vue'
import TodoDetailsView from '../views/TodoDetailsView.vue'
import AddTodoView from '../views/AddTodoView.vue'
import EditTodoView from '../views/EditTodoView.vue'

const todoRoutes =  [
  { path: '/', component: TodoListView },
  { path: '/todo/:id', component: TodoDetailsView, props: true },
  { path: '/add', component: AddTodoView },
  { path: '/edit/:id', component: EditTodoView, props: true }
]

export default createRouter({
  history: createWebHistory('/frontend/'),
  routes: todoRoutes
})