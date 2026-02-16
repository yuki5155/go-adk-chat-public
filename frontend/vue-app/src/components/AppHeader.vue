<script setup lang="ts">
import { ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { useAuth } from '@/composables/useAuth'

const { user, isAuthenticated } = useAuth()
const menuOpen = ref(false)
const router = useRouter()

router.afterEach(() => {
  menuOpen.value = false
})
</script>

<template>
  <header class="header">
    <div class="header-content">
      <h1 class="header-title">Go Google Auth</h1>
      <button class="menu-toggle" @click="menuOpen = !menuOpen" aria-label="Toggle menu">
        <span class="menu-bar"></span>
        <span class="menu-bar"></span>
        <span class="menu-bar"></span>
      </button>
      <nav class="header-nav" :class="{ open: menuOpen }">
        <RouterLink to="/">Session Test</RouterLink>
        <RouterLink to="/about">About</RouterLink>
        <template v-if="isAuthenticated">
          <RouterLink to="/dashboard" class="nav-dashboard">
            <img
              v-if="user?.picture"
              :src="user.picture"
              :alt="user.name"
              class="nav-avatar"
              referrerpolicy="no-referrer"
            />
            <span>Dashboard</span>
          </RouterLink>
        </template>
        <template v-else>
          <RouterLink to="/login" class="nav-login">Sign In</RouterLink>
        </template>
      </nav>
    </div>
  </header>
</template>

<style scoped>
.menu-toggle {
  display: none;
  flex-direction: column;
  justify-content: center;
  gap: 5px;
  background: none;
  border: none;
  cursor: pointer;
  padding: 8px;
}

.menu-bar {
  display: block;
  width: 22px;
  height: 2px;
  background: var(--color-text, #333);
  border-radius: 2px;
  transition: transform 0.2s, opacity 0.2s;
}

@media (max-width: 768px) {
  .menu-toggle {
    display: flex;
  }
}

.nav-dashboard {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.nav-avatar {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  object-fit: cover;
}

.nav-login {
  background: var(--color-primary, #4285f4);
  color: white !important;
  padding: 0.5rem 1rem;
  border-radius: 6px;
  font-weight: 500;
}

.nav-login:hover {
  background: #3367d6;
}
</style>
