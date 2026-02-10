import { ref } from 'vue'
import type {
  RoleRequest,
  UserRoleData,
  ApiResponse,
  RoleRequestListResponse,
  UserRoleListResponse,
  RequestRolePayload,
  ApproveRequestPayload,
  RejectRequestPayload,
  UserRole,
} from '../types/role'

// Backend URL from environment
const backendUrl = import.meta.env.VITE_BACKEND_URL || 'http://localhost:8080'

// Shared state
const isLoading = ref(false)
const error = ref<string | null>(null)

export function useRoleManagement() {
  // Request a role (user endpoint)
  async function requestRole(requestedRole: 'user' | 'subscriber'): Promise<RoleRequest | null> {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch(`${backendUrl}/api/role/request`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({ requested_role: requestedRole } as RequestRolePayload),
      })

      const data: ApiResponse<RoleRequest> = await response.json()

      if (response.ok && data.success) {
        return data.data || null
      } else {
        error.value = data.error?.message || 'Failed to request role'
        return null
      }
    } catch (err) {
      console.error('Role request error:', err)
      error.value = err instanceof Error ? err.message : 'An error occurred'
      return null
    } finally {
      isLoading.value = false
    }
  }

  // List pending role requests (admin endpoint)
  async function listPendingRequests(): Promise<RoleRequest[]> {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch(`${backendUrl}/api/admin/role/requests`, {
        method: 'GET',
        credentials: 'include',
      })

      const data: ApiResponse<RoleRequestListResponse> = await response.json()

      if (response.ok && data.success) {
        return data.data?.requests || []
      } else {
        error.value = data.error?.message || 'Failed to fetch pending requests'
        return []
      }
    } catch (err) {
      console.error('Fetch pending requests error:', err)
      error.value = err instanceof Error ? err.message : 'An error occurred'
      return []
    } finally {
      isLoading.value = false
    }
  }

  // Approve a role request (admin endpoint)
  async function approveRequest(requestId: string, notes?: string): Promise<boolean> {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch(`${backendUrl}/api/admin/role/approve`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({ request_id: requestId, notes } as ApproveRequestPayload),
      })

      const data: ApiResponse<{ message: string }> = await response.json()

      if (response.ok && data.success) {
        return true
      } else {
        error.value = data.error?.message || 'Failed to approve request'
        return false
      }
    } catch (err) {
      console.error('Approve request error:', err)
      error.value = err instanceof Error ? err.message : 'An error occurred'
      return false
    } finally {
      isLoading.value = false
    }
  }

  // Reject a role request (admin endpoint)
  async function rejectRequest(requestId: string, notes: string): Promise<boolean> {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch(`${backendUrl}/api/admin/role/reject`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({ request_id: requestId, notes } as RejectRequestPayload),
      })

      const data: ApiResponse<{ message: string }> = await response.json()

      if (response.ok && data.success) {
        return true
      } else {
        error.value = data.error?.message || 'Failed to reject request'
        return false
      }
    } catch (err) {
      console.error('Reject request error:', err)
      error.value = err instanceof Error ? err.message : 'An error occurred'
      return false
    } finally {
      isLoading.value = false
    }
  }

  // List users by role (admin endpoint)
  async function listUsersByRole(role: UserRole): Promise<UserRoleData[]> {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch(`${backendUrl}/api/admin/role/users?role=${role}`, {
        method: 'GET',
        credentials: 'include',
      })

      const data: ApiResponse<UserRoleListResponse> = await response.json()

      if (response.ok && data.success) {
        return data.data?.users || []
      } else {
        error.value = data.error?.message || 'Failed to fetch users'
        return []
      }
    } catch (err) {
      console.error('Fetch users error:', err)
      error.value = err instanceof Error ? err.message : 'An error occurred'
      return []
    } finally {
      isLoading.value = false
    }
  }

  // Clear error
  function clearError(): void {
    error.value = null
  }

  return {
    // State
    isLoading,
    error,

    // Actions
    requestRole,
    listPendingRequests,
    approveRequest,
    rejectRequest,
    listUsersByRole,
    clearError,
  }
}
