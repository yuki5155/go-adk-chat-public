<script setup lang="ts">
import { onMounted } from 'vue'
import AppHeader from './components/AppHeader.vue'
import { useAuth } from '@/composables/useAuth'

const { initAuth } = useAuth()

// Initialize authentication state on app mount
onMounted(async () => {
  await initAuth()
})
</script>

<template>
  <div id="app">
    <AppHeader />
    <main class="main-content">
      <RouterView />
    </main>
  </div>
</template>

<style>
.main-content {
  flex: 1;
  width: 100%;
  padding: 2rem;
}
</style>
