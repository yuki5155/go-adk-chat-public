<template>
  <div class="role-request-container">
    <div class="role-request-card">
      <h1>Request Role Access</h1>
      <p class="description">
        Request additional access to features by selecting a role below.
      </p>

      <!-- Success Message -->
      <div v-if="successMessage" class="alert alert-success">
        {{ successMessage }}
      </div>

      <!-- Error Message -->
      <div v-if="roleError" class="alert alert-error">
        {{ roleError }}
      </div>

      <!-- Role Selection Form -->
      <form v-if="!successMessage" @submit.prevent="handleSubmit" class="role-form">
        <div class="form-group">
          <label for="role-select">Select Role:</label>
          <select
            id="role-select"
            v-model="selectedRole"
            class="role-select"
            :disabled="isLoading"
            required
          >
            <option value="">-- Choose a role --</option>
            <option value="subscriber">Subscriber (Access to chatbot features)</option>
            <option value="user">User (Basic access)</option>
          </select>
        </div>

        <div class="role-info" v-if="selectedRole">
          <h3>{{ selectedRole === 'subscriber' ? 'Subscriber' : 'User' }} Role</h3>
          <p v-if="selectedRole === 'subscriber'">
            Subscribers get access to AI chatbot features and premium content.
          </p>
          <p v-else>
            Users have basic access to the platform features.
          </p>
        </div>

        <button
          type="submit"
          class="submit-btn"
          :disabled="!selectedRole || isLoading"
        >
          {{ isLoading ? 'Submitting...' : 'Submit Request' }}
        </button>
      </form>

      <div v-else class="success-actions">
        <button @click="goToDashboard" class="btn-secondary">
          Go to Dashboard
        </button>
        <button @click="resetForm" class="btn-link">
          Submit Another Request
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '../composables/useAuth'
import { useRoleManagement } from '../composables/useRoleManagement'

const router = useRouter()
const { user } = useAuth()
const { requestRole, isLoading, error: roleError, clearError } = useRoleManagement()

const selectedRole = ref<'user' | 'subscriber' | ''>('')
const successMessage = ref<string | null>(null)

// Block admin/root users from accessing this page
onMounted(() => {
  if (user.value && ['admin', 'root'].includes(user.value.role)) {
    router.push('/dashboard')
  }
})

async function handleSubmit() {
  if (!selectedRole.value) return

  clearError()
  successMessage.value = null

  const result = await requestRole(selectedRole.value)

  if (result) {
    successMessage.value = `Your request for ${selectedRole.value} role has been submitted successfully! An admin will review it soon.`
  }
}

function resetForm() {
  selectedRole.value = ''
  successMessage.value = null
  clearError()
}

function goToDashboard() {
  router.push('/dashboard')
}
</script>

<style scoped>
.role-request-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 80vh;
  padding: 2rem;
}

.role-request-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  padding: 3rem;
  max-width: 600px;
  width: 100%;
}

h1 {
  font-size: 2rem;
  color: #333;
  margin-bottom: 0.5rem;
}

.description {
  color: #666;
  margin-bottom: 2rem;
}

.alert {
  padding: 1rem;
  border-radius: 8px;
  margin-bottom: 1.5rem;
}

.alert-success {
  background-color: #d4edda;
  color: #155724;
  border: 1px solid #c3e6cb;
}

.alert-error {
  background-color: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
}

.role-form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.form-group {
  display: flex;
  flex-direction: column;
}

label {
  font-weight: 600;
  margin-bottom: 0.5rem;
  color: #333;
}

.role-select {
  padding: 0.75rem;
  border: 2px solid #e0e0e0;
  border-radius: 8px;
  font-size: 1rem;
  transition: border-color 0.3s;
}

.role-select:focus {
  outline: none;
  border-color: #4CAF50;
}

.role-select:disabled {
  background-color: #f5f5f5;
  cursor: not-allowed;
}

.role-info {
  background: #f8f9fa;
  padding: 1.5rem;
  border-radius: 8px;
  border-left: 4px solid #4CAF50;
}

.role-info h3 {
  margin-top: 0;
  color: #333;
  font-size: 1.25rem;
}

.role-info p {
  margin-bottom: 0;
  color: #666;
}

.submit-btn {
  background: #4CAF50;
  color: white;
  border: none;
  padding: 1rem 2rem;
  font-size: 1rem;
  font-weight: 600;
  border-radius: 8px;
  cursor: pointer;
  transition: background-color 0.3s, transform 0.2s;
}

.submit-btn:hover:not(:disabled) {
  background: #45a049;
  transform: translateY(-2px);
}

.submit-btn:disabled {
  background: #ccc;
  cursor: not-allowed;
  transform: none;
}

.success-actions {
  display: flex;
  gap: 1rem;
  margin-top: 2rem;
}

.btn-secondary {
  flex: 1;
  background: #4CAF50;
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  font-size: 1rem;
  border-radius: 8px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.btn-secondary:hover {
  background: #45a049;
}

.btn-link {
  flex: 1;
  background: transparent;
  color: #4CAF50;
  border: 2px solid #4CAF50;
  padding: 0.75rem 1.5rem;
  font-size: 1rem;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
}

.btn-link:hover {
  background: #4CAF50;
  color: white;
}
</style>
