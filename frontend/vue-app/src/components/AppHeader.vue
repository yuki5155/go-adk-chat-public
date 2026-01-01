<script setup lang="ts">
import { RouterLink } from 'vue-router'
import { useAuth } from '@/composables/useAuth'

const { user, isAuthenticated } = useAuth()
</script>

<template>
  <header class="header">
    <div class="header-content">
      <h1 class="header-title">Go Google Auth</h1>
      <nav class="header-nav">
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
