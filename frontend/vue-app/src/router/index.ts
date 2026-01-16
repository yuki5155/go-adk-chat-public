import { createRouter, createWebHistory } from 'vue-router'
import type { RouteLocationNormalized, NavigationGuardNext } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import AboutView from '../views/AboutView.vue'
import LoginView from '../views/LoginView.vue'
import DashboardView from '../views/DashboardView.vue'
import AdminDashboardView from '../views/AdminDashboardView.vue'
import RoleRequestView from '../views/RoleRequestView.vue'
import AdminPendingRequestsView from '../views/AdminPendingRequestsView.vue'
import AdminUserManagementView from '../views/AdminUserManagementView.vue'
import ChatbotView from '../views/ChatbotView.vue'

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

// Check if user has root privileges
async function isRootUser(): Promise<boolean> {
  try {
    const response = await fetch(`${backendUrl}/api/me`, {
      method: 'GET',
      credentials: 'include',
    })
    if (!response.ok) return false
    const data = await response.json()
    return data.user?.role === 'root'
  } catch {
    return false
  }
}

// Check if user has admin or root privileges
async function isAdminUser(): Promise<boolean> {
  try {
    const response = await fetch(`${backendUrl}/api/me`, {
      method: 'GET',
      credentials: 'include',
    })
    if (!response.ok) return false
    const data = await response.json()
    const role = data.user?.role
    return role === 'admin' || role === 'root'
  } catch {
    return false
  }
}

// Check if user has subscriber, admin, or root privileges
async function isSubscriberOrAdmin(): Promise<boolean> {
  try {
    const response = await fetch(`${backendUrl}/api/me`, {
      method: 'GET',
      credentials: 'include',
    })
    if (!response.ok) return false
    const data = await response.json()
    const role = data.user?.role
    return role === 'subscriber' || role === 'admin' || role === 'root'
  } catch {
    return false
  }
}

// Navigation guard for subscriber routes (requires subscriber, admin, or root)
async function requireSubscriber(
  _to: RouteLocationNormalized,
  _from: RouteLocationNormalized,
  next: NavigationGuardNext
) {
  const authenticated = await isAuthenticated()
  if (!authenticated) {
    next('/login')
    return
  }

  const hasAccess = await isSubscriberOrAdmin()
  if (hasAccess) {
    next()
  } else {
    next('/dashboard')
  }
}

// Navigation guard for admin routes (requires admin or root privileges)
async function requireAdmin(
  _to: RouteLocationNormalized,
  _from: RouteLocationNormalized,
  next: NavigationGuardNext
) {
  const authenticated = await isAuthenticated()
  if (!authenticated) {
    next('/login')
    return
  }

  const isAdmin = await isAdminUser()
  if (isAdmin) {
    next()
  } else {
    next('/dashboard')
  }
}

// Navigation guard for admin routes (requires root privileges)
async function requireRoot(
  _to: RouteLocationNormalized,
  _from: RouteLocationNormalized,
  next: NavigationGuardNext
) {
  const authenticated = await isAuthenticated()
  if (!authenticated) {
    next('/login')
    return
  }

  const isRoot = await isRootUser()
  if (isRoot) {
    next()
  } else {
    next('/dashboard')
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
    {
      path: '/role-request',
      name: 'role-request',
      component: RoleRequestView,
      beforeEnter: requireAuth,
      meta: {
        requiresAuth: true,
      },
    },
    {
      path: '/chatbot',
      name: 'chatbot',
      component: ChatbotView,
      beforeEnter: requireSubscriber,
      meta: {
        requiresAuth: true,
        requiresSubscriber: true,
      },
    },
    {
      path: '/admin/dashboard',
      name: 'admin-dashboard',
      component: AdminDashboardView,
      beforeEnter: requireRoot,
      meta: {
        requiresAuth: true,
        requiresRoot: true,
      },
    },
    {
      path: '/admin/role-requests',
      name: 'admin-role-requests',
      component: AdminPendingRequestsView,
      beforeEnter: requireAdmin,
      meta: {
        requiresAuth: true,
        requiresAdmin: true,
      },
    },
    {
      path: '/admin/users',
      name: 'admin-users',
      component: AdminUserManagementView,
      beforeEnter: requireAdmin,
      meta: {
        requiresAuth: true,
        requiresAdmin: true,
      },
    },
  ],
})

export default router
