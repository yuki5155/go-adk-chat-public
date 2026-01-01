import { createRouter, createWebHistory } from 'vue-router'
import type { RouteLocationNormalized, NavigationGuardNext } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import AboutView from '../views/AboutView.vue'
import LoginView from '../views/LoginView.vue'
import DashboardView from '../views/DashboardView.vue'

// Backend URL for auth check
const backendUrl = import.meta.env.VITE_BACKEND_URL || 'http://localhost:8080'

// Check if user is authenticated by calling the backend
async function isAuthenticated(): Promise<boolean> {
  try {
    const response = await fetch(`${backendUrl}/api/me`, {
      method: 'GET',
      credentials: 'include',
    })
    return response.ok
  } catch {
    return false
  }
}

// Navigation guard for protected routes
async function requireAuth(
  _to: RouteLocationNormalized,
  _from: RouteLocationNormalized,
  next: NavigationGuardNext
) {
  const authenticated = await isAuthenticated()
  if (authenticated) {
    next()
  } else {
    next('/login')
  }
}

// Navigation guard for guest routes (redirect if already logged in)
async function requireGuest(
  _to: RouteLocationNormalized,
  _from: RouteLocationNormalized,
  next: NavigationGuardNext
) {
  const authenticated = await isAuthenticated()
  if (authenticated) {
    next('/dashboard')
  } else {
    next()
  }
}

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/about',
      name: 'about',
      component: AboutView,
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView,
      beforeEnter: requireGuest,
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: DashboardView,
      beforeEnter: requireAuth,
      meta: {
        requiresAuth: true,
      },
    },
  ],
})

export default router
