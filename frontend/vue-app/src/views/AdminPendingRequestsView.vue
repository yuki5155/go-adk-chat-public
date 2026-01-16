<template>
  <div class="admin-pending-requests">
    <div class="header">
      <h1>Pending Role Requests</h1>
      <button @click="refreshRequests" class="refresh-btn" :disabled="isLoading">
        {{ isLoading ? 'Loading...' : 'Refresh' }}
      </button>
    </div>

    <!-- Error Message -->
    <div v-if="roleError" class="alert alert-error">
      {{ roleError }}
    </div>

    <!-- Success Message -->
    <div v-if="successMessage" class="alert alert-success">
      {{ successMessage }}
    </div>

    <!-- Loading State -->
    <div v-if="isLoading && !requests.length" class="loading">
      <p>Loading pending requests...</p>
    </div>

    <!-- Empty State -->
    <div v-else-if="!requests.length" class="empty-state">
      <p>No pending role requests at the moment.</p>
    </div>

    <!-- Requests Table -->
    <div v-else class="requests-table-container">
      <table class="requests-table">
        <thead>
          <tr>
            <th>User Email</th>
            <th>Requested Role</th>
            <th>Requested At</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="request in requests" :key="request.request_id">
            <td>{{ request.user_email }}</td>
            <td>
              <span class="role-badge" :class="`role-${request.requested_role}`">
                {{ request.requested_role }}
              </span>
            </td>
            <td>{{ formatDate(request.requested_at) }}</td>
            <td class="actions">
              <button
                @click="openApproveModal(request)"
                class="btn-approve"
                :disabled="isLoading"
              >
                Approve
              </button>
              <button
                @click="openRejectModal(request)"
                class="btn-reject"
                :disabled="isLoading"
              >
                Reject
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Approve Modal -->
    <div v-if="showApproveModal" class="modal-overlay" @click="closeModals">
      <div class="modal" @click.stop>
        <h2>Approve Role Request</h2>
        <p>
          Approve <strong>{{ selectedRequest?.user_email }}</strong> for
          <strong>{{ selectedRequest?.requested_role }}</strong> role?
        </p>
        <div class="form-group">
          <label for="approve-notes">Notes (optional):</label>
          <textarea
            id="approve-notes"
            v-model="notes"
            rows="3"
            placeholder="Add any notes about this approval..."
          ></textarea>
        </div>
        <div class="modal-actions">
          <button @click="confirmApprove" class="btn-approve" :disabled="isLoading">
            {{ isLoading ? 'Approving...' : 'Confirm Approve' }}
          </button>
          <button @click="closeModals" class="btn-cancel" :disabled="isLoading">
            Cancel
          </button>
        </div>
      </div>
    </div>

    <!-- Reject Modal -->
    <div v-if="showRejectModal" class="modal-overlay" @click="closeModals">
      <div class="modal" @click.stop>
        <h2>Reject Role Request</h2>
        <p>
          Reject <strong>{{ selectedRequest?.user_email }}</strong> for
          <strong>{{ selectedRequest?.requested_role }}</strong> role?
        </p>
        <div class="form-group">
          <label for="reject-notes">Rejection reason (required):</label>
          <textarea
            id="reject-notes"
            v-model="notes"
            rows="3"
            placeholder="Explain why this request is being rejected..."
            required
          ></textarea>
        </div>
        <div class="modal-actions">
          <button
            @click="confirmReject"
            class="btn-reject"
            :disabled="!notes.trim() || isLoading"
          >
            {{ isLoading ? 'Rejecting...' : 'Confirm Reject' }}
          </button>
          <button @click="closeModals" class="btn-cancel" :disabled="isLoading">
            Cancel
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoleManagement } from '../composables/useRoleManagement'
import type { RoleRequest } from '../types/role'

const { listPendingRequests, approveRequest, rejectRequest, isLoading, error: roleError, clearError } = useRoleManagement()

const requests = ref<RoleRequest[]>([])
const showApproveModal = ref(false)
const showRejectModal = ref(false)
const selectedRequest = ref<RoleRequest | null>(null)
const notes = ref('')
const successMessage = ref<string | null>(null)

onMounted(() => {
  fetchRequests()
})

async function fetchRequests() {
  clearError()
  requests.value = await listPendingRequests()
}

async function refreshRequests() {
  successMessage.value = null
  await fetchRequests()
}

function openApproveModal(request: RoleRequest) {
  selectedRequest.value = request
  notes.value = ''
  showApproveModal.value = true
  successMessage.value = null
  clearError()
}

function openRejectModal(request: RoleRequest) {
  selectedRequest.value = request
  notes.value = ''
  showRejectModal.value = true
  successMessage.value = null
  clearError()
}

function closeModals() {
  showApproveModal.value = false
  showRejectModal.value = false
  selectedRequest.value = null
  notes.value = ''
}

async function confirmApprove() {
  if (!selectedRequest.value) return

  const success = await approveRequest(selectedRequest.value.request_id, notes.value)
  if (success) {
    successMessage.value = `Successfully approved ${selectedRequest.value.user_email}'s request`
    closeModals()
    await fetchRequests()
  }
}

async function confirmReject() {
  if (!selectedRequest.value || !notes.value.trim()) return

  const success = await rejectRequest(selectedRequest.value.request_id, notes.value)
  if (success) {
    successMessage.value = `Rejected ${selectedRequest.value.user_email}'s request`
    closeModals()
    await fetchRequests()
  }
}

function formatDate(dateString: string): string {
  return new Date(dateString).toLocaleString()
}
</script>

<style scoped>
.admin-pending-requests {
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

.loading,
.empty-state {
  text-align: center;
  padding: 3rem;
  color: #666;
  font-size: 1.1rem;
}

.requests-table-container {
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.requests-table {
  width: 100%;
  border-collapse: collapse;
}

.requests-table thead {
  background: #f8f9fa;
}

.requests-table th {
  padding: 1rem;
  text-align: left;
  font-weight: 600;
  color: #333;
  border-bottom: 2px solid #e0e0e0;
}

.requests-table td {
  padding: 1rem;
  border-bottom: 1px solid #e0e0e0;
}

.requests-table tbody tr:hover {
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

.actions {
  display: flex;
  gap: 0.5rem;
}

.btn-approve,
.btn-reject {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 600;
  transition: all 0.3s;
}

.btn-approve {
  background: #4CAF50;
  color: white;
}

.btn-approve:hover:not(:disabled) {
  background: #45a049;
}

.btn-reject {
  background: #f44336;
  color: white;
}

.btn-reject:hover:not(:disabled) {
  background: #da190b;
}

.btn-approve:disabled,
.btn-reject:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal {
  background: white;
  border-radius: 12px;
  padding: 2rem;
  max-width: 500px;
  width: 90%;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.2);
}

.modal h2 {
  margin-top: 0;
  color: #333;
}

.form-group {
  margin: 1.5rem 0;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 600;
  color: #333;
}

.form-group textarea {
  width: 100%;
  padding: 0.75rem;
  border: 2px solid #e0e0e0;
  border-radius: 8px;
  font-family: inherit;
  font-size: 1rem;
  resize: vertical;
}

.form-group textarea:focus {
  outline: none;
  border-color: #4CAF50;
}

.modal-actions {
  display: flex;
  gap: 1rem;
  margin-top: 1.5rem;
}

.btn-cancel {
  flex: 1;
  background: #f5f5f5;
  color: #333;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 600;
  transition: background-color 0.3s;
}

.btn-cancel:hover:not(:disabled) {
  background: #e0e0e0;
}

.modal-actions .btn-approve,
.modal-actions .btn-reject {
  flex: 1;
  padding: 0.75rem 1.5rem;
}
</style>
