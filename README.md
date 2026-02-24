# Go ADK Chat

An AI chat application built with Go and Google Gemini, featuring conversation memory, Google OAuth authentication, and serverless AWS Lambda deployment.

## Overview

Go ADK Chat is a full-stack application that provides an AI-powered chat interface with persistent conversation memory. Users authenticate via Google OAuth, create chat threads, and interact with Google Gemini models. Chat responses are streamed in real-time via Server-Sent Events (SSE).

**Key highlights:**

- **AI Chat with Memory** -- Conversations maintain context across messages using a session/memory system backed by DynamoDB
- **Function / Tool Calling** -- Gemini can invoke registered tools (function calling) mid-conversation and stream the results back to the user
- **Real-time Streaming** -- Chat responses are streamed token-by-token via SSE through API Gateway Lambda streaming
- **Role-based Access** -- Admin dashboard for managing user roles (root, admin, subscriber)
- **Google-only Auth** -- We intentionally support only Google / Google Workspace to minimize attack surface and operational complexity
- **Low Cost** -- Runs on AWS serverless with pay-as-you-go pricing. Typical personal usage stays under $5/month
- **Self-hosted** -- No managed SaaS dependencies. Everything runs on your own AWS account
- **Clean Architecture** -- Domain-Driven Design with clear separation of domain, application, infrastructure, and presentation layers

## Tech Stack

| Layer | Technologies |
|-------|-------------|
| **Backend** | Go 1.25, Gin, Google Generative AI SDK, AWS Lambda |
| **Frontend** | Vue 3, TypeScript, Vite |
| **Database** | DynamoDB (local for dev, AWS for prod) |
| **Auth** | Google OAuth (GIS), JWT (HttpOnly cookies) |
| **Infrastructure** | AWS CDK, Lambda, API Gateway, Secrets Manager |
| **CI/CD** | GitHub Actions |

## Project Structure

```
go-adk-chat/
├── backend/
│   ├── cmd/
│   │   ├── api/main.go              # HTTP server entry point
│   │   └── lambda/                   # Lambda function handlers
│   │       ├── chat-stream/          # SSE streaming endpoint
│   │       ├── chat-message/         # Send message
│   │       ├── chat-threads-*/       # Thread CRUD
│   │       ├── chat-models/          # List available models
│   │       ├── auth-google/          # Google OAuth
│   │       └── ...
│   ├── internal/
│   │   ├── domain/                   # Entities & repository interfaces
│   │   │   ├── chat/                 # Thread, Session, Memory, Event
│   │   │   ├── user/
│   │   │   └── role/
│   │   ├── application/              # Use cases
│   │   │   ├── chat/                 # Create/list/get/delete threads, send message
│   │   │   ├── auth/
│   │   │   └── admin/
│   │   ├── infrastructure/           # Implementations
│   │   │   ├── adk/                  # Gemini client & AIRunner adapter
│   │   │   ├── persistence/          # DynamoDB repositories
│   │   │   ├── auth/                 # Google OAuth & JWT
│   │   │   └── container/            # Dependency injection
│   │   └── presentation/http/        # Handlers, middleware, router
│   ├── compose.yml
│   └── Makefile
├── frontend/
│   └── vue-app/                      # Vue 3 + TypeScript SPA
│       ├── src/
│       │   ├── views/ChatbotView.vue # Chat interface
│       │   ├── composables/useChat.ts # Chat state & SSE streaming
│       │   └── ...
│       └── compose.yml
├── iac/                              # AWS CDK stacks
│   └── bin/
│       ├── lambda.ts                 # Lambda + API Gateway
│       ├── dynamodb.ts               # DynamoDB tables
│       ├── network.ts                # VPC
│       └── secrets.ts                # Secrets Manager
├── .github/workflows/                # CI/CD pipelines
└── Makefile                          # Root orchestration
```

## Getting Started

### Prerequisites

- Docker Desktop & Docker Compose v2+
- Git
- Google Cloud OAuth 2.0 credentials ([console](https://console.cloud.google.com/apis/credentials))
- Google AI API key for Gemini ([ai.google.dev](https://ai.google.dev/))

Optional (for local dev without Docker): Go 1.25+, Node.js 22.15+

### Setup

1. **Clone the repository**

   ```bash
   git clone https://github.com/yuki5155/go-adk-chat.git
   cd go-adk-chat
   ```

2. **Configure environment variables**

   ```bash
   # Set your Google credentials and run the setup script
   CLIENT_ID="your-id.apps.googleusercontent.com" \
   CLIENT_SECRET="your-secret" \
   make setup-env
   ```

   Then edit `backend/.env` to add your Gemini API key:

   ```env
   GOOGLE_AI_API_KEY=your-gemini-api-key
   ```

3. **Start all services**

   ```bash
   make up
   ```

4. **Access the application**

   - Frontend: http://localhost:5173
   - Backend API: http://localhost:8080

### Useful Commands

```bash
make help              # Show all available commands
make up                # Start backend + frontend
make down              # Stop all containers
make logs              # Show all logs
make logs-backend      # Backend logs only
make rebuild           # Stop, rebuild, and restart
make status            # Check service status
make dynamodb-tables   # List DynamoDB tables
make clean             # Remove containers and volumes
```

## Architecture

### Chat Flow

1. User creates a chat thread (with optional model selection)
2. User sends a message within the thread
3. Backend retrieves conversation history from the session/memory store
4. Google Gemini API generates a response, optionally invoking registered tools via function calling
5. Tool results (if any) are fed back to Gemini and surfaced to the frontend as `tool_start` / `tool_end` SSE events
6. Final text response is streamed back to the frontend token-by-token via SSE
7. Message and memory artifacts are persisted to DynamoDB

### DynamoDB Tables

| Table | Purpose |
|-------|---------|
| `chat-threads` | Thread metadata (title, model, owner) |
| `chat-sessions` | Session state with message history |
| `chat-events` | Audit log of chat operations |
| `chat-memories` | Conversation memory artifacts |
| `users` | User profiles |
| `roles` | Role assignments and requests |

### Authentication

- Google Sign-In via GIS library on the frontend
- Backend validates the ID token and issues JWT access/refresh tokens as HttpOnly cookies
- Access token: 15-minute expiry / Refresh token: 7-day expiry
- Role hierarchy: root > admin > subscriber > user

## Backend Development

```bash
cd backend

make test                # Run all tests
make test-coverage       # Generate coverage report
make test-coverage-html  # Open HTML coverage report
make dev                 # Run with hot reload (Air)
make build               # Build binary
```

Test coverage: **75%**.

## Frontend Development

```bash
cd frontend/vue-app

npm install
npm run dev              # Start dev server
npm run build            # Production build
npm run lint             # Lint & fix
npm run type-check       # TypeScript check
```

## AWS Deployment

The project deploys to AWS via GitHub Actions workflows. See the [Deployment Guide](doc/deployment.md) for full instructions.

**Deploy order:**

1. GitHub Actions IAM Role (one-time setup)
2. Secrets Stack
3. DynamoDB Stack
4. Network Stack
5. Lambda Stack or ECS Stack
6. Frontend Stack

## Environment Variables

| Variable | Description |
|----------|-------------|
| `GOOGLE_CLIENT_ID` | Google OAuth Client ID |
| `GOOGLE_CLIENT_SECRET` | Google OAuth Client Secret |
| `JWT_SECRET` | Secret for JWT token signing |
| `GOOGLE_AI_API_KEY` | Google Gemini API key |
| `GEMINI_MODEL` | Model name (default: `gemini-2.0-flash`) |
| `ROOT_USER_EMAIL` | Email granted root privileges |
| `VITE_BACKEND_URL` | Backend URL for frontend (default: `http://localhost:8080`) |
| `VITE_GOOGLE_CLIENT_ID` | Google Client ID for frontend |

See `backend/.env.example` for the full list.

## License

MIT
