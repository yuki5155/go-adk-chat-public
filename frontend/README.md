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

## Development

The `vue-app` directory is mounted as a volume, so changes to the code will be reflected immediately.
