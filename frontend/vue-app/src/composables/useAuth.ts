import { ref, computed, readonly } from 'vue'

// User interface
export interface User {
  id: string
  email: string
  name: string
  picture: string
}

// Authentication state
const user = ref<User | null>(null)
const isLoading = ref(false)
const error = ref<string | null>(null)

// Backend URL from environment
const backendUrl = import.meta.env.VITE_BACKEND_URL
if (!backendUrl) {
  console.warn('VITE_BACKEND_URL is not set! Falling back to localhost:8080 for development.')
}
const finalBackendUrl = backendUrl || 'http://localhost:8080'

// Computed properties
const isAuthenticated = computed(() => user.value !== null)

// Initialize authentication state by checking with backend
async function initAuth(): Promise<void> {
  isLoading.value = true
  error.value = null

  console.log('Initializing auth with backend:', finalBackendUrl)

  try {
    const response = await fetch(`${finalBackendUrl}/api/me`, {
      method: 'GET',
      credentials: 'include',
    })

    if (response.ok) {
      const data = await response.json()
      user.value = data.user
    } else if (response.status === 401) {
      // Try to refresh the token
      const refreshed = await refreshToken()
      if (refreshed) {
        // Retry fetching user info
        const retryResponse = await fetch(`${finalBackendUrl}/api/me`, {
          method: 'GET',
          credentials: 'include',
        })
        if (retryResponse.ok) {
          const data = await retryResponse.json()
          user.value = data.user
        }
      }
    } else {
      const errorData = await response.json().catch(() => ({}))
      console.error('Auth initialization failed:', response.status, errorData)
    }
  } catch (err) {
    console.error('Failed to initialize auth:', err)
    error.value = err instanceof Error ? `Connection Error: ${err.message}` : 'Failed to connect to backend'
    user.value = null
  } finally {
    isLoading.value = false
  }
}

// Login with Google credential (ID token from GIS)
async function loginWithGoogle(credential: string): Promise<boolean> {
  isLoading.value = true
  error.value = null

  try {
    const response = await fetch(`${finalBackendUrl}/auth/google`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify({ credential }),
    })

    const data = await response.json()

    if (response.ok) {
      user.value = data.user
      return true
    } else {
      error.value = data.message || 'Login failed'
      return false
    }
  } catch (err) {
    console.error('Login error:', err)
    error.value = err instanceof Error ? err.message : 'An error occurred during login'
    return false
  } finally {
    isLoading.value = false
  }
}

// Refresh the access token
async function refreshToken(): Promise<boolean> {
  try {
    const response = await fetch(`${finalBackendUrl}/auth/refresh`, {
      method: 'POST',
      credentials: 'include',
    })

    return response.ok
  } catch (err) {
    console.error('Token refresh error:', err)
    return false
  }
}

// Logout
async function logout(): Promise<void> {
  isLoading.value = true

  try {
    await fetch(`${finalBackendUrl}/auth/logout`, {
      method: 'POST',
      credentials: 'include',
    })
  } catch (err) {
    console.error('Logout error:', err)
  } finally {
    user.value = null
    isLoading.value = false
  }
}

// Clear error
function clearError(): void {
  error.value = null
}

// Export composable
export function useAuth() {
  return {
    // State (readonly to prevent external mutations)
    user: readonly(user),
    isLoading: readonly(isLoading),
    error: readonly(error),
    isAuthenticated,

    // Actions
    initAuth,
    loginWithGoogle,
    refreshToken,
    logout,
    clearError,
  }
}
