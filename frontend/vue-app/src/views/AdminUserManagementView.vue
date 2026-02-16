<template>
  <div class="admin-user-management">
    <div class="header">
      <h1>User Management</h1>
      <div class="header-actions">
        <select v-model="selectedRole" class="role-filter" @change="fetchUsers">
          <option value="user">Users</option>
          <option value="subscriber">Subscribers</option>
        </select>
        <button @click="refreshUsers" class="refresh-btn" :disabled="isLoading">
          {{ isLoading ? 'Loading...' : 'Refresh' }}
        </button>
      </div>
    </div>

    <!-- Info message -->
    <div class="info-message">
      <p>
        <strong>Note:</strong> Admin and Root users are system-level accounts and do not appear in this list.
        They cannot request roles through the standard flow.
      </p>
    </div>

    <!-- Error Message -->
    <div v-if="roleError" class="alert alert-error">
      {{ roleError }}
    </div>

    <!-- Loading State -->
    <div v-if="isLoading && !users.length" class="loading">
      <p>Loading users...</p>
    </div>

    <!-- Empty State -->
    <div v-else-if="!users.length" class="empty-state">
      <p>No {{ selectedRole }}s found.</p>
    </div>

    <!-- Users Table -->
    <div v-else class="users-table-container">
      <div class="table-info">
        <p>Total {{ selectedRole }}s: <strong>{{ users.length }}</strong></p>
      </div>
      <table class="users-table">
        <thead>
          <tr>
            <th>Email</th>
            <th>Role</th>
            <th>Status</th>
            <th>Created At</th>
            <th>Last Updated</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="user in users" :key="user.user_id">
            <td>{{ user.email }}</td>
            <td>
              <span class="role-badge" :class="`role-${user.role}`">
                {{ user.role }}
              </span>
            </td>
            <td>
              <span class="status-badge" :class="`status-${user.status}`">
                {{ user.status }}
              </span>
            </td>
            <td>{{ formatDate(user.created_at) }}</td>
            <td>{{ formatDate(user.updated_at) }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoleManagement } from '../composables/useRoleManagement'
import type { UserRoleData, UserRole } from '../types/role'

const { listUsersByRole, isLoading, error: roleError, clearError } = useRoleManagement()

const selectedRole = ref<UserRole>('user')
const users = ref<UserRoleData[]>([])

onMounted(() => {
  fetchUsers()
})

async function fetchUsers() {
  clearError()
  users.value = await listUsersByRole(selectedRole.value)
}

async function refreshUsers() {
  await fetchUsers()
}

function formatDate(dateString: string): string {
  return new Date(dateString).toLocaleString()
}
</script>

<style scoped>
.admin-user-management {
  padding: 2rem;
  max-width: 1200px;
  margin: 0 auto;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

h1 {
  font-size: 2rem;
  color: #333;
}

.header-actions {
  display: flex;
  gap: 1rem;
  align-items: center;
}

.role-filter {
  padding: 0.75rem 1rem;
  border: 2px solid #e0e0e0;
  border-radius: 8px;
  font-size: 1rem;
  background: white;
  cursor: pointer;
  transition: border-color 0.3s;
}

.role-filter:focus {
  outline: none;
  border-color: #4CAF50;
}

.refresh-btn {
  background: #4CAF50;
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 600;
  transition: background-color 0.3s;
}

.refresh-btn:hover:not(:disabled) {
  background: #45a049;
}

.refresh-btn:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.info-message {
  background-color: #d1ecf1;
  color: #0c5460;
  border: 1px solid #bee5eb;
  padding: 1rem;
  border-radius: 8px;
  margin-bottom: 1.5rem;
}

.info-message p {
  margin: 0;
  font-size: 0.95rem;
}

.alert {
  padding: 1rem;
  border-radius: 8px;
  margin-bottom: 1.5rem;
}

.alert-error {
  background-color: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
}

.loading,
.empty-state {
  text-align: center;
  padding: 3rem;
  color: #666;
  font-size: 1.1rem;
}

.users-table-container {
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.table-info {
  padding: 1rem 1.5rem;
  background: #f8f9fa;
  border-bottom: 1px solid #e0e0e0;
}

.table-info p {
  margin: 0;
  color: #666;
  font-size: 0.95rem;
}

.users-table {
  width: 100%;
  border-collapse: collapse;
}

.users-table thead {
  background: #f8f9fa;
}

.users-table th {
  padding: 1rem;
  text-align: left;
  font-weight: 600;
  color: #333;
  border-bottom: 2px solid #e0e0e0;
}

.users-table td {
  padding: 1rem;
  border-bottom: 1px solid #e0e0e0;
}

.users-table tbody tr:hover {
  background: #f8f9fa;
}

.role-badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.875rem;
  font-weight: 600;
}

.role-subscriber {
  background: #e3f2fd;
  color: #1976d2;
}

.role-user {
  background: #f3e5f5;
  color: #7b1fa2;
}

.role-admin {
  background: #fff3e0;
  color: #f57c00;
}

.status-badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.875rem;
  font-weight: 600;
}

.status-active {
  background: #e8f5e9;
  color: #2e7d32;
}

.status-suspended {
  background: #ffebee;
  color: #c62828;
}

/* Mobile responsive */
@media (max-width: 768px) {
  .admin-user-management {
    padding: 1rem;
  }

  .header {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }

  h1 {
    font-size: 1.5rem;
  }

  .header-actions {
    width: 100%;
  }

  .role-filter {
    flex: 1;
    min-width: 0;
  }

  .refresh-btn {
    white-space: nowrap;
  }

  .users-table-container {
    overflow-x: auto;
  }

  .users-table th,
  .users-table td {
    padding: 0.75rem 0.5rem;
    font-size: 0.875rem;
  }
}
</style>
