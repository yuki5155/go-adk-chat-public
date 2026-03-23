# Frontend

Node.js/Vue.js development environment

## Setup

```bash
# Start the development server
docker compose up -d

# Install dependencies
docker compose exec frontend npm install

# Stop the server
docker compose down
```

## Accessing the Application

- Frontend: http://localhost:5173

## Environment Variables

The following variables are read from `frontend/vue-app/.env.development` (created automatically by `make setup-env`):

| Variable | Description | Default |
|----------|-------------|---------|
| `VITE_BACKEND_URL` | URL of the Go backend API | `http://localhost:8080` |
| `VITE_GOOGLE_CLIENT_ID` | Google OAuth 2.0 Client ID | — |

See `frontend/vue-app/.env.example` for the full list.

## Development

The `vue-app` directory is mounted as a volume, so changes to the code will be reflected immediately.

For more details on the Vue app structure and state management, see [`vue-app/README.md`](./vue-app/README.md).
