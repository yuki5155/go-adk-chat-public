# Go Google Auth Frontend

Frontend application built with Vue.js 3 + TypeScript

## Features

- AI chat interface with real-time SSE streaming and typewriter effect
- Tool/function-calling activity display (`tool_start` / `tool_end` events)
- Chat thread management (create, select, delete)
- Model selection per thread (lists available Gemini models from backend)
- Google Sign-In via GIS library
- Role request workflow (request subscriber/admin access)
- Admin dashboard (approve/reject role requests, manage users)
- Mobile-responsive layout

## Setup

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Build
npm run build

# Preview
npm run preview
```

## Development Environment

- Node.js >= 22.15.0
- npm

## Tech Stack

- Vue 3
- TypeScript
- Vite
- Vue Router

## State Management

This project does not use a dedicated state management library (e.g. Pinia or Vuex). Instead, state is managed via **composables with module-level reactive refs**.

Each composable (`useAuth`, `useChat`, `useRoleManagement`) declares its state outside the exported function, making it a shared singleton across all components:

```ts
// State is declared at module scope — shared across all consumers
const user = ref<User | null>(null)

export function useAuth() {
  return { user, ... }
}
```

This pattern provides the same shared-state behavior as a store, without the added dependency. It is a lightweight, idiomatic Vue 3 approach suitable for an app of this size.

| Composable | Responsibility |
|---|---|
| `useAuth` | Authenticated user, login/logout, token refresh |
| `useChat` | Chat threads, messages, streaming |
| `useRoleManagement` | Role requests and admin approval flow |
