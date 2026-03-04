// import { defineStore } from 'pinia'

// export const useTodoStore = defineStore('todo', {
// 	state: () => {
//         return {
//             todos: JSON.parse(localStorage.getItem("todos")) || []
//         }
//     },
// 	getters:{
//         getCount(){
//             return this.todos.length
//         }
//     },
// 	actions: {
//         getById(id) {
//             return this.todos.find(t => String(t.id) === id)
//         },
// 		saveToStorage(){
//             localStorage.setItem("todos", JSON.stringify(this.todos))
//         },
// 		toggleTodo(id) {
// 			const t = this.todos.find(t => t.id === id)
// 			if (t) {
// 				t.completed = !t.completed
// 			}
// 			this.saveToStorage()
// 		},
// 		addTask(title, description){
//             const exists = this.todos.find(t => t.title === title)
//             if(exists){
//                 alert("Task already exists in the list!")
//                 return false
//             }
//             const id = Date.now()
//             this.todos.push({
//                 id: id,
//                 title: title,
//                 description: description,
//                 completed: false,
//             })
//             this.saveToStorage()
//             return true
//         },
// 		removeById(id){
//             const index = this.todos.findIndex(t => t.id === id)
//             if(index != -1){
//                 this.todos.splice(index, 1)
//                 this.saveToStorage()
//             }
//         },
//         update(todo){
//             const index = this.todos.findIndex(t => String(t.id) === String(todo.id))
//             if (index != -1){
//                 this.todos[index] = {...todo}
//             }
//             this.saveToStorage()
//         }
//   	}
// })


import { defineStore } from 'pinia'

const API_BASE = '/todos'

export const useTodoStore = defineStore('todos', {
  state: () => ({
    todos: [],
    loading: false,
    error: null,
  }),

  actions: {
    async fetchTodos() {
      const res = await fetch(API_BASE)
      this.todos = await res.json()
    },

    async addTask(title, description) {
      const res = await fetch(`${API_BASE}/add`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ title, description })
      })

      if (!res.ok) {
        const err = await res.json()
        alert(err.error || 'Failed to add task')
        return false
      }

      const todo = await res.json()
      this.todos.push(todo)
      return true
    },

    async update(todo) {
      const res = await fetch(API_BASE, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(todo)
      })

      if (!res.ok) {
        alert('Failed to update')
        return
      }

      const index = this.todos.findIndex(t => t.id === todo.id)
      if (index !== -1) this.todos[index] = todo
    },

    async getById(id) {
        this.loading = true
        try {
            const res = await fetch(`${API_BASE}/${id}`)
            if (!res.ok) throw new Error('Todo not found')

            const todo = await res.json()

            // ✅ keep store in sync
            const index = this.todos.findIndex(t => t.id === todo.id)
            if (index === -1) {
            this.todos.push(todo)
            } else {
            this.todos[index] = todo
            }

            return todo
        } finally {
            this.loading = false
        }
    },

    async removeById(id) {
      const res = await fetch(`${API_BASE}/${id}`, {
        method: 'DELETE'
      })

      if (!res.ok) {
        alert('Failed to delete')
        return
      }

      this.todos = this.todos.filter(t => t.id !== id)
    }
  }
})