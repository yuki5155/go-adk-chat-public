<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '@/composables/useAuth'

const router = useRouter()
const { user, isLoading, isAuthenticated, logout, initAuth } = useAuth()

// Handle logout
async function handleLogout() {
  await logout()
  router.push('/login')
}

onMounted(async () => {
  // Refresh user data on mount
  if (!isAuthenticated.value) {
    await initAuth()
    // If still not authenticated after init, redirect to login
    if (!isAuthenticated.value) {
      router.push('/login')
    }
  }
})
</script>

<template>
  <div class="dashboard-page">
    <!-- Loading state -->
    <div v-if="isLoading" class="dashboard-loading">
      <div class="spinner"></div>
      <p>Loading your dashboard...</p>
    </div>

    <!-- Dashboard content -->
    <div v-else-if="user" class="dashboard-content">
      <!-- User profile card -->
      <section class="profile-card">
        <div class="profile-header">
          <img
            v-if="user.picture"
            :src="user.picture"
            :alt="user.name"
            class="profile-avatar"
            referrerpolicy="no-referrer"
          />
          <div v-else class="profile-avatar-placeholder">
            {{ user.name?.charAt(0)?.toUpperCase() || '?' }}
          </div>
          <div class="profile-info">
            <h1>{{ user.name }}</h1>
            <p class="text-muted">{{ user.email }}</p>
          </div>
        </div>
        <button class="btn btn-danger" @click="handleLogout">
          Sign Out
        </button>
      </section>

      <!-- Welcome section -->
      <section class="welcome-section">
        <h2>Welcome to your Dashboard!</h2>
        <p class="text-muted">
          You've successfully authenticated with Google. This is a protected page
          that only authenticated users can access.
        </p>
      </section>

      <!-- User details card -->
      <section class="details-card">
        <h3>Your Account Information</h3>
        <div class="details-grid">
          <div class="detail-item">
            <label>User ID</label>
            <span class="detail-value">{{ user.id }}</span>
          </div>
          <div class="detail-item">
            <label>Email</label>
            <span class="detail-value">{{ user.email }}</span>
          </div>
          <div class="detail-item">
            <label>Name</label>
            <span class="detail-value">{{ user.name }}</span>
          </div>
        </div>
      </section>

      <!-- Authentication info -->
      <section class="auth-info-card">
        <h3>Authentication Details</h3>
        <ul class="auth-features">
          <li>
            <span class="feature-icon">✓</span>
            Google Identity Services (GIS) authentication
          </li>
          <li>
            <span class="feature-icon">✓</span>
            JWT-based session management
          </li>
          <li>
            <span class="feature-icon">✓</span>
            Secure HttpOnly cookies
          </li>
          <li>
            <span class="feature-icon">✓</span>
            Automatic token refresh
          </li>
        </ul>
      </section>
    </div>

    <!-- Not authenticated state -->
    <div v-else class="not-authenticated">
      <p>You are not authenticated. Redirecting to login...</p>
    </div>
  </div>
</template>

<style scoped>
.dashboard-page {
  padding: 2rem;
  max-width: 900px;
  margin: 0 auto;
}

.dashboard-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 50vh;
  gap: 1rem;
}

.spinner {
  width: 48px;
  height: 48px;
  border: 4px solid var(--color-border, #e0e0e0);
  border-top-color: var(--color-primary, #4285f4);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.dashboard-content {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.profile-card {
  background: var(--color-surface, #ffffff);
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  padding: 1.5rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 1rem;
}

.profile-header {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.profile-avatar {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  object-fit: cover;
  border: 3px solid var(--color-border, #e0e0e0);
}

.profile-avatar-placeholder {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: var(--color-primary, #4285f4);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  font-weight: 600;
}

.profile-info h1 {
  font-size: 1.5rem;
  font-weight: 600;
  margin: 0 0 0.25rem 0;
  color: var(--color-text, #1a1a1a);
}

.profile-info p {
  margin: 0;
}

.welcome-section {
  background: linear-gradient(135deg, #4285f4 0%, #34a853 100%);
  border-radius: 12px;
  padding: 2rem;
  color: white;
}

.welcome-section h2 {
  margin: 0 0 0.5rem 0;
  font-size: 1.5rem;
}

.welcome-section p {
  margin: 0;
  opacity: 0.9;
}

.details-card,
.auth-info-card {
  background: var(--color-surface, #ffffff);
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  padding: 1.5rem;
}

.details-card h3,
.auth-info-card h3 {
  margin: 0 0 1rem 0;
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--color-text, #1a1a1a);
}

.details-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
}

.detail-item {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.detail-item label {
  font-size: 0.875rem;
  color: var(--color-text-muted, #666);
  font-weight: 500;
}

.detail-value {
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 0.875rem;
  background: var(--color-background, #f5f5f5);
  padding: 0.5rem;
  border-radius: 6px;
  word-break: break-all;
}

.auth-features {
  list-style: none;
  padding: 0;
  margin: 0;
  display: grid;
  gap: 0.75rem;
}

.auth-features li {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.5rem;
  background: var(--color-background, #f5f5f5);
  border-radius: 6px;
}

.feature-icon {
  color: var(--color-success, #34a853);
  font-weight: bold;
}

.not-authenticated {
  text-align: center;
  padding: 4rem;
  color: var(--color-text-muted, #666);
}

@media (max-width: 600px) {
  .profile-card {
    flex-direction: column;
    text-align: center;
  }

  .profile-header {
    flex-direction: column;
  }
}
</style>
