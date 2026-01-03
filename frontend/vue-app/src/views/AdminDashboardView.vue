<template>
  <div class="admin-dashboard">
    <h1>Admin Dashboard</h1>
    <div v-if="loading">Loading...</div>
    <div v-else-if="error" class="error">
      {{ error }}
    </div>
    <div v-else-if="dashboardData">
      <p>{{ dashboardData.message }}</p>
      <div class="user-info">
        <h2>Your Info</h2>
        <p><strong>ID:</strong> {{ dashboardData.user.id }}</p>
        <p><strong>Email:</strong> {{ dashboardData.user.email }}</p>
        <p><strong>Name:</strong> {{ dashboardData.user.name }}</p>
        <p><strong>Root User:</strong> {{ dashboardData.user.isRoot ? 'Yes' : 'No' }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const loading = ref(true)
const error = ref('')
const dashboardData = ref<any>(null)

const backendUrl = import.meta.env.VITE_BACKEND_URL || 'http://localhost:8080'

onMounted(async () => {
  try {
    const response = await fetch(`${backendUrl}/admin/dashboard`, {
      method: 'GET',
      credentials: 'include',
    })

    if (response.status === 403) {
      error.value = 'Access denied: Root privileges required'
      setTimeout(() => {
        router.push('/dashboard')
      }, 2000)
      return
    }

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }

    dashboardData.value = await response.json()
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to load admin dashboard'
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.admin-dashboard {
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
}

.error {
  color: red;
  padding: 10px;
  background-color: #fee;
  border-radius: 4px;
  margin: 10px 0;
}

.user-info {
  background-color: #f5f5f5;
  padding: 20px;
  border-radius: 8px;
  margin-top: 20px;
}

.user-info h2 {
  margin-top: 0;
}

.user-info p {
  margin: 10px 0;
}
</style>
