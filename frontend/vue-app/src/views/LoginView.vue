<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '@/composables/useAuth'

// Google Identity Services types
declare global {
  interface Window {
    google: {
      accounts: {
        id: {
          initialize: (config: GoogleIdConfig) => void
          renderButton: (element: HTMLElement, config: GoogleButtonConfig) => void
          prompt: () => void
        }
      }
    }
  }
}

interface GoogleIdConfig {
  client_id: string
  callback: (response: GoogleCredentialResponse) => void
  auto_select?: boolean
  cancel_on_tap_outside?: boolean
}

interface GoogleButtonConfig {
  theme?: 'outline' | 'filled_blue' | 'filled_black'
  size?: 'large' | 'medium' | 'small'
  type?: 'standard' | 'icon'
  text?: 'signin_with' | 'signup_with' | 'continue_with' | 'signin'
  shape?: 'rectangular' | 'pill' | 'circle' | 'square'
  logo_alignment?: 'left' | 'center'
  width?: number
  locale?: string
}

interface GoogleCredentialResponse {
  credential: string
  select_by: string
  clientId?: string
}

const router = useRouter()
const { loginWithGoogle, isLoading, error, isAuthenticated } = useAuth()

const googleButtonRef = ref<HTMLDivElement | null>(null)
const googleClientId = import.meta.env.VITE_GOOGLE_CLIENT_ID || ''
const isGoogleLoaded = ref(false)
const loadError = ref<string | null>(null)

// Handle Google credential response
async function handleCredentialResponse(response: GoogleCredentialResponse) {
  const success = await loginWithGoogle(response.credential)
  if (success) {
    router.push('/dashboard')
  }
}

// Initialize Google Sign-In
function initializeGoogleSignIn() {
  if (!googleClientId) {
    loadError.value = 'Google Client ID is not configured. Please set VITE_GOOGLE_CLIENT_ID environment variable.'
    return
  }

  if (!window.google?.accounts?.id) {
    // Retry after a short delay if GIS is not loaded yet
    setTimeout(initializeGoogleSignIn, 100)
    return
  }

  try {
    window.google.accounts.id.initialize({
      client_id: googleClientId,
      callback: handleCredentialResponse,
      auto_select: false,
      cancel_on_tap_outside: true,
    })

    if (googleButtonRef.value) {
      window.google.accounts.id.renderButton(googleButtonRef.value, {
        theme: 'outline',
        size: 'large',
        type: 'standard',
        text: 'signin_with',
        shape: 'rectangular',
        logo_alignment: 'left',
        width: 300,
      })
    }

    isGoogleLoaded.value = true
  } catch (err) {
    console.error('Failed to initialize Google Sign-In:', err)
    loadError.value = 'Failed to initialize Google Sign-In'
  }
}

onMounted(() => {
  // Redirect if already authenticated
  if (isAuthenticated.value) {
    router.push('/dashboard')
    return
  }

  // Wait for GIS script to load
  initializeGoogleSignIn()
})
</script>

<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-header">
        <h1>Welcome</h1>
        <p class="text-muted">Sign in to access your dashboard</p>
      </div>

      <div class="login-content">
        <!-- Loading state -->
        <div v-if="isLoading" class="login-loading">
          <div class="spinner"></div>
          <p>Signing you in...</p>
        </div>

        <!-- Error message -->
        <div v-else-if="loadError || error" class="login-error">
          <p class="text-danger">{{ loadError || error }}</p>
          <button v-if="error" class="btn btn-secondary mt-2" @click="() => {}">
            Try Again
          </button>
        </div>

        <!-- Google Sign-In button -->
        <div v-else class="google-signin-container">
          <div ref="googleButtonRef" class="google-button"></div>
          
          <div v-if="!isGoogleLoaded && !loadError" class="loading-google">
            <p class="text-muted">Loading Google Sign-In...</p>
          </div>
        </div>

        <!-- Additional info -->
        <div class="login-info">
          <p class="text-muted text-small">
            By signing in, you agree to our terms of service and privacy policy.
          </p>
        </div>
      </div>
    </div>

    <!-- Debug info (development only) -->
    <div v-if="!googleClientId" class="debug-info">
      <h3>Configuration Missing</h3>
      <p>Add <code>VITE_GOOGLE_CLIENT_ID</code> to your environment variables.</p>
      <p>Example in <code>.env.development</code>:</p>
      <pre>VITE_GOOGLE_CLIENT_ID=your-client-id.apps.googleusercontent.com</pre>
    </div>
  </div>
</template>

<style scoped>
.login-page {
  min-height: calc(100vh - 120px);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 2rem;
}

.login-card {
  background: var(--color-surface, #ffffff);
  border-radius: 12px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.1);
  padding: 2.5rem;
  width: 100%;
  max-width: 400px;
}

.login-header {
  text-align: center;
  margin-bottom: 2rem;
}

.login-header h1 {
  font-size: 1.75rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
  color: var(--color-text, #1a1a1a);
}

.login-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1.5rem;
}

.google-signin-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100%;
}

.google-button {
  display: flex;
  justify-content: center;
}

.loading-google {
  margin-top: 1rem;
}

.login-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  padding: 2rem;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid var(--color-border, #e0e0e0);
  border-top-color: var(--color-primary, #4285f4);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.login-error {
  text-align: center;
  padding: 1rem;
}

.login-info {
  text-align: center;
  padding-top: 1rem;
  border-top: 1px solid var(--color-border, #e0e0e0);
  width: 100%;
}

.text-small {
  font-size: 0.875rem;
}

.debug-info {
  margin-top: 2rem;
  padding: 1.5rem;
  background: #fff3cd;
  border: 1px solid #ffc107;
  border-radius: 8px;
  max-width: 500px;
}

.debug-info h3 {
  color: #856404;
  margin-bottom: 0.5rem;
}

.debug-info p {
  color: #856404;
  margin-bottom: 0.5rem;
}

.debug-info pre {
  background: #f8f9fa;
  padding: 0.5rem;
  border-radius: 4px;
  overflow-x: auto;
  font-size: 0.875rem;
}

.debug-info code {
  background: rgba(0, 0, 0, 0.1);
  padding: 0.125rem 0.25rem;
  border-radius: 3px;
}
</style>
