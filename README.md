# Go Google Auth

A full-stack web application for Google OAuth authentication with comprehensive Cookie/Session management testing capabilities.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21%2B-blue)](https://go.dev/)
[![Vue Version](https://img.shields.io/badge/Vue-3.5-green)](https://vuejs.org/)
[![Node Version](https://img.shields.io/badge/Node-22.15.0%2B-brightgreen)](https://nodejs.org/)

## 📋 Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Test Coverage](#test-coverage)
- [Technology Stack](#technology-stack)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [Usage](#usage)
- [API Documentation](#api-documentation)
- [Development](#development)
- [AWS Deployment](#aws-deployment)
- [Docker Commands](#docker-commands)
- [Environment Variables](#environment-variables)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [License](#license)
- [Author](#author)
- [Resources](#resources)
- [Project Status](#project-status)

## 🎯 Overview

This project demonstrates a modern full-stack application architecture with:
- **Backend**: Go-based REST API using Gin framework
- **Frontend**: Vue.js 3 with TypeScript and Vite
- **Database**: DynamoDB (AWS DynamoDB Local for development)
- **Authentication**: Google Identity Services (GIS) with JWT tokens
- **Session Management**: Secure HttpOnly cookie-based JWT sessions

## ✨ Features

### Current Features
- ✅ **Google Identity Services (GIS) authentication**
- ✅ **JWT-based session management** (access + refresh tokens)
- ✅ **Protected routes and authorization**
- ✅ **Secure HttpOnly cookies**
- ✅ Cookie/Session testing interface
- ✅ Set-Cookie header validation
- ✅ Cookie transmission verification
- ✅ Real-time cookie display and management
- ✅ CORS configuration for cross-origin requests
- ✅ Docker-based development environment
- ✅ Hot reload for both frontend and backend

### Planned Features
- ✅ AWS deployment with CDK (Network, Backend, Frontend, Secrets stacks)

## 🧪 Test Coverage

This project maintains comprehensive unit test coverage using Go's testing framework with `testify` and `gomock` for clean, type-safe mocking.

![Coverage Status](https://img.shields.io/badge/coverage-95.6%25-brightgreen)
![Tests](https://img.shields.io/badge/tests-133%20passing-brightgreen)
![Quality](https://img.shields.io/badge/quality-production%20ready-blue)

### Current Coverage by Layer

| Layer | Coverage | Description |
|-------|----------|-------------|
| **Domain** | 100.0% | User entities and value objects |
| **Application** | 93.6% | Use cases and business logic |
| **Infrastructure** | 94.1% | Config, Auth (JWT: 90.7%), Persistence (97.6%) |
| **Presentation** | 100.0% | HTTP handlers |
| **Overall** | **95.6%** ✅ | **133 tests** across all layers |

**Achievement**: Exceeded the 80% coverage requirement with **95.6%** overall coverage through comprehensive testing of all layers using gomock for clean, maintainable tests.

### Running Tests

```bash
cd backend

# Run all tests
make test

# Run tests with coverage
make test-coverage

# Generate HTML coverage report
make test-coverage-html

# Check if coverage meets threshold (80%)
make coverage-check

# View detailed coverage summary
make coverage-summary

# Clean coverage files
make clean-coverage
```

**Manual commands** (if not using Make):
```bash
# Run all tests
go test -v ./internal/...

# Run tests with coverage
go test -coverprofile=coverage.out -covermode=atomic ./internal/...

# View coverage report
go tool cover -html=coverage.out

# View total coverage
go tool cover -func=coverage.out | grep total:
```

### Test Organization

Tests are organized following the DDD (Domain-Driven Design) architecture:

- **Domain Layer**: `services/jwt_test.go` - Pure business logic tests
- **Infrastructure Layer**: `config/config_test.go` - Configuration and environment tests
- **Application Layer**: `middleware/auth_test.go` - Middleware and request flow tests
- **Presentation Layer**: `handlers/*_test.go` - HTTP handler tests

### CI/CD Integration

The GitHub Actions workflow automatically:
- Runs all tests on every push and pull request
- Generates coverage reports
- Enforces minimum 70% coverage threshold (currently at 95.6%)
- Comments coverage results on pull requests
- Uploads coverage to Codecov

**Coverage Achievement** 🎉:
- **Current**: 95.6% overall coverage
- **Target**: 80%+ ✅ **ACHIEVED**
- **Quality**: Production-ready with comprehensive test coverage across all layers

### Key Testing Features

- **Interface-based Design**: Dependency injection for full testability
- **Gomock Integration**: Type-safe, generated mocks using go.uber.org/mock
- **Comprehensive Scenarios**: Success paths, error cases, edge conditions, and security validations
- **Layer Testing**: Complete coverage of Domain (100%), Application (93.6%), Infrastructure (94.1%), and Presentation (100%) layers

## 🛠 Technology Stack

### Backend
| Technology | Version | Purpose |
|------------|---------|---------|
| **Go** | 1.21+ | Primary backend language |
| **Gin** | Latest | Web framework |
| **DynamoDB** | Local | NoSQL database |
| **Air** | Latest | Hot reload for development |
| **Docker** | Latest | Containerization |

### Frontend
| Technology | Version | Purpose |
|------------|---------|---------|
| **Vue.js** | 3.5+ | Frontend framework |
| **TypeScript** | 5.8+ | Type-safe JavaScript |
| **Vite** | 6.0+ | Build tool and dev server |
| **Vue Router** | 4.5+ | Client-side routing |
| **Vue DevTools** | 8.0+ | Development tools |

### Infrastructure
| Technology | Purpose |
|------------|---------|
| **Docker Compose** | Multi-container orchestration |
| **AWS CDK** | Infrastructure as Code (TypeScript) |
| **GitHub Actions** | CI/CD pipeline |

## 📁 Project Structure

```
go-google-auth-public/
│
├── backend/                    # Go backend application
│   ├── cmd/
│   │   └── api/
│   │       └── main.go        # Application entry point
│   ├── internal/
│   │   ├── config/            # Configuration management
│   │   ├── handlers/          # HTTP request handlers
│   │   │   ├── auth.go        # Google OAuth & JWT handlers
│   │   │   ├── cookie.go      # Cookie management handlers
│   │   │   ├── health.go      # Health check endpoints
│   │   │   └── hello.go       # Example endpoints
│   │   ├── middleware/        # HTTP middleware
│   │   │   └── auth.go        # JWT authentication middleware
│   │   └── services/          # Business logic services
│   │       └── jwt.go         # JWT token generation/validation
│   ├── dockers/
│   │   ├── Dockerfile.local   # Development Docker image
│   │   └── Dockerfile.prod    # Production Docker image
│   ├── compose.yml            # Backend Docker Compose config
│   ├── .air.toml             # Air hot reload configuration
│   ├── .env.example          # Environment variables template
│   ├── go.mod                # Go module definition
│   └── Makefile              # Backend build commands
│
├── frontend/                  # Vue.js frontend application
│   ├── vue-app/
│   │   ├── src/
│   │   │   ├── components/   # Reusable Vue components
│   │   │   │   └── AppHeader.vue
│   │   │   ├── composables/  # Vue composition functions
│   │   │   │   └── useAuth.ts    # Authentication state management
│   │   │   ├── views/        # Page-level components
│   │   │   │   ├── HomeView.vue      # Cookie testing page
│   │   │   │   ├── LoginView.vue     # Google Sign-In page
│   │   │   │   ├── DashboardView.vue # Protected user dashboard
│   │   │   │   └── AboutView.vue     # About page
│   │   │   ├── router/       # Vue Router configuration
│   │   │   │   └── index.ts  # Routes with auth guards
│   │   │   ├── assets/       # Static assets
│   │   │   │   ├── main.css
│   │   │   │   └── styles/   # Organized CSS
│   │   │   │       ├── base/
│   │   │   │       ├── components/
│   │   │   │       ├── layouts/
│   │   │   │       ├── pages/
│   │   │   │       └── utilities/
│   │   │   └── main.ts       # Application entry point
│   │   ├── index.html        # HTML template (includes GIS script)
│   │   ├── .env.example      # Environment variables template
│   │   ├── vite.config.ts
│   │   ├── tsconfig.json
│   │   └── package.json
│   ├── Dockerfile            # Frontend Docker image
│   ├── compose.yml           # Frontend Docker Compose config
│   └── docker-entrypoint.sh  # Container startup script
│
├── iac/                       # Infrastructure as Code
│   ├── bin/
│   │   ├── backend.ts        # Backend infrastructure stack
│   │   └── network.ts        # Network infrastructure stack
│   ├── cdk.json
│   └── package.json
│
├── .github/
│   └── workflows/            # GitHub Actions CI/CD
│       ├── backend.yml
│       ├── backend-setup.yml
│       └── network.yml
│
├── Makefile                  # Root-level commands
├── README.md                 # This file
├── .cursorrules             # Development guidelines
├── .gitignore
└── LICENSE
```

## 🚀 Getting Started

### Prerequisites

Ensure you have the following installed:

- **Docker Desktop** (latest version)
- **Docker Compose** (v2.0+)
- **Git**

Optional (for local development without Docker):
- **Go** (1.21+)
- **Node.js** (22.15.0+)
- **npm** (10.9+)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/yuki5155/go-google-auth-public.git
   cd go-google-auth-public
   ```

2. **Set up Google OAuth credentials** (see [Google OAuth Setup](#google-oauth-setup) below)

3. **Configure environment variables**
   ```bash
   # Backend
   cd backend
   make env  # Creates .env from .env.example (requires Docker)
   # Edit .env with your Google credentials
   
   # Frontend
   cd ../frontend/vue-app
   cp .env.example .env.development
   # Edit .env.development with your Google Client ID
   ```

4. **Start all services using Docker Compose**
   ```bash
   # Start backend
   cd backend && docker compose up -d
   
   # Start frontend (in a new terminal)
   cd frontend && docker compose up -d
   ```

5. **Verify services are running**
   ```bash
   docker ps | grep go-google-auth
   ```

   You should see three containers running:
   - `go-google-auth-frontend` (Port 5173)
   - `go-google-auth-app` (Port 8080)
   - `go-google-auth-dynamodb` (Port 8000)

6. **Access the application**
   - Frontend: http://localhost:5173
   - Backend API: http://localhost:8080
   - DynamoDB Local: http://localhost:8000

### Google OAuth Setup

To enable Google Sign-In, you need to create OAuth 2.0 credentials in Google Cloud Console:

1. **Go to Google Cloud Console**
   - Visit https://console.cloud.google.com/apis/credentials

2. **Create a new project** (or select an existing one)

3. **Configure OAuth consent screen**
   - Go to "OAuth consent screen"
   - Select "External" user type
   - Fill in the required fields (App name, User support email, Developer contact)
   - Add scopes: `email`, `profile`, `openid`
   - Save and continue

4. **Create OAuth 2.0 Client ID**
   - Go to "Credentials" → "Create Credentials" → "OAuth client ID"
   - Application type: **Web application**
   - Name: `Web client 1` (or any name)
   - **Authorized JavaScript origins:**
     ```
     http://localhost:5173
     http://localhost:8080
     ```
   - **Authorized redirect URIs:**
     ```
     http://localhost:8080/auth/google/callback
     ```
   - Click "Create"

5. **Copy your credentials**
   - Copy the **Client ID** (ends with `.apps.googleusercontent.com`)
   - Copy the **Client Secret**

6. **Update environment files**

   **Backend `.env`:**
   ```env
   GOOGLE_CLIENT_ID=your-client-id.apps.googleusercontent.com
   GOOGLE_CLIENT_SECRET=your-client-secret
   JWT_SECRET=your-random-secret-key  # Generate with: openssl rand -base64 32
   ```

   **Frontend `.env.development`:**
   ```env
   VITE_BACKEND_URL=http://localhost:8080
   VITE_GOOGLE_CLIENT_ID=your-client-id.apps.googleusercontent.com
   ```

7. **Restart containers**
   ```bash
   cd backend && docker compose down && docker compose up -d
   cd ../frontend && docker compose down && docker compose up -d
   ```

## 💡 Usage

### Cookie/Session Testing

The application provides an interactive interface for testing cookie behavior:

1. **Open the application**
   - Navigate to http://localhost:5173 in your browser

2. **Test Set-Cookie**
   - Click the **"Test Set-Cookie"** button
   - The backend will send a `Set-Cookie` header
   - Check the response in the UI

3. **View Current Cookies**
   - The **"Current Cookies"** section displays all cookies stored in your browser
   - Click **"Refresh"** to update the display

4. **Test Cookie Sending**
   - Click **"Check Cookie Sending"** to verify cookies are sent to the backend
   - The backend will confirm receipt and display cookie contents

5. **Clear Cookies**
   - Use the **"Clear All Cookies"** button to remove test cookies

### Backend URL Configuration

You can change the backend API URL in the frontend interface if your backend is running on a different port or host.

## 📚 API Documentation

### Health Endpoints

#### `GET /health`
Health check endpoint for monitoring.

**Response:**
```json
{
  "status": "ok"
}
```

#### `GET /health/ready`
Readiness probe for Kubernetes/ECS deployments.

### Cookie Testing Endpoints

#### `GET /api/set-cookie`
Sets a test cookie in the response.

**Headers:**
- `Access-Control-Allow-Origin`: Frontend URL
- `Access-Control-Allow-Credentials`: true

**Response:**
```json
{
  "message": "Cookie has been set",
  "cookieName": "test_session",
  "cookieValue": "session_value_12345",
  "timestamp": "2025-12-14T10:09:50Z"
}
```

**Set-Cookie Header:**
```
test_session=session_value_12345; Path=/; Max-Age=3600; HttpOnly
```

#### `GET /api/check-cookie`
Verifies if cookies are received from the client.

**Required Headers:**
- `Cookie`: Must include cookies from browser

**Response:**
```json
{
  "cookieReceived": true,
  "cookies": "test_session=session_value_12345",
  "testSession": "session_value_12345",
  "timestamp": "2025-12-14T10:10:00Z"
}
```

### Authentication Endpoints

#### `POST /auth/google`
Authenticates user with Google ID token and creates JWT session.

**Request Body:**
```json
{
  "credential": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response:**
```json
{
  "message": "Login successful",
  "user": {
    "id": "123456789",
    "email": "user@example.com",
    "name": "John Doe",
    "picture": "https://lh3.googleusercontent.com/..."
  }
}
```

**Cookies Set:**
- `access_token` - JWT access token (15 min expiry, HttpOnly)
- `refresh_token` - JWT refresh token (7 days expiry, HttpOnly)

#### `POST /auth/refresh`
Refreshes the access token using the refresh token cookie.

**Response:**
```json
{
  "message": "Token refreshed successfully"
}
```

#### `POST /auth/logout`
Logs out the user by clearing authentication cookies.

**Response:**
```json
{
  "message": "Logged out successfully"
}
```

#### `GET /api/me` (Protected)
Returns the current authenticated user's information.

**Required:** Valid `access_token` cookie

**Response:**
```json
{
  "user": {
    "id": "123456789",
    "email": "user@example.com",
    "name": "John Doe",
    "picture": "https://lh3.googleusercontent.com/..."
  }
}
```

**Error Response (401):**
```json
{
  "error": "unauthorized",
  "message": "Access token not found"
}
```

## 🔧 Development

### Backend Development

#### Running Locally (with Air hot reload)
```bash
cd backend

# Install dependencies
go mod download

# Run with hot reload
make dev

# Run tests
make test

# Build binary
make build
```

#### Project Structure
```go
// Handler pattern
type Handler struct {
    Path string
}

func NewHandler() *Handler {
    return &Handler{
        Path: "/api/endpoint",
    }
}

func (h *Handler) Handle(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "success",
    })
}
```

### Frontend Development

#### Running Locally
```bash
cd frontend/vue-app

# Install dependencies
npm install

# Start dev server
npm run dev

# Build for production
npm run build

# Run linter
npm run lint

# Type check
npm run type-check
```

#### Component Structure
```vue
<script setup lang="ts">
import { ref } from 'vue'

interface Props {
  title: string
}

const props = defineProps<Props>()
const count = ref(0)

const increment = () => {
  count.value++
}
</script>

<template>
  <div>
    <h1>{{ props.title }}</h1>
    <button @click="increment">Count: {{ count }}</button>
  </div>
</template>

<style scoped>
/* Component-specific styles */
</style>
```

## ☁️ AWS Deployment

### Prerequisites

1. AWS CLI configured with appropriate credentials
2. AWS CDK bootstrapped in your account
3. GitHub Actions OIDC provider configured

### Initial Setup

#### 1. Configure GitHub Actions IAM Role

```bash
cd iac

# Update CloudFormation template with your GitHub org/repo from git remote
make update-cf-defaults

# Or run init-cf to also see deployment instructions
make init-cf
```

#### 2. Deploy GitHub Actions IAM Role

```bash
aws cloudformation create-stack \
  --stack-name go-google-auth-github-actions \
  --template-body file://iac/cloudformations/github-actions-role.yaml \
  --capabilities CAPABILITY_NAMED_IAM \
  --parameters \
    ParameterKey=RoleName,ParameterValue=GoGoogleAuthGitHubActionsRole
```

#### 3. Configure GitHub Repository Secrets

Add these secrets in GitHub repository settings (Settings → Secrets and variables → Actions):

| Secret | Description |
|--------|-------------|
| `AWS_ROLE_TO_ASSUME` | IAM Role ARN from CloudFormation output |
| `AWS_REGION` | AWS Region (e.g., `ap-northeast-1`) |
| `PROJECT_NAME` | Project name (e.g., `go-google-auth`) |
| `ROOT_DOMAIN` | Your domain (e.g., `example.com`) |

### Deployment Order

Deploy stacks in this order:

```bash
# 1. Network Stack
# GitHub Actions: Run "Network Stack" workflow

# 2. Secrets Stack (creates placeholder secrets)
# GitHub Actions: Run "Secrets Stack" workflow

# 3. Update secret values in AWS Secrets Manager
aws secretsmanager put-secret-value \
  --secret-id "go-google-auth/dev/google-auth" \
  --secret-string '{
    "GOOGLE_CLIENT_ID": "your-client-id.apps.googleusercontent.com",
    "GOOGLE_CLIENT_SECRET": "your-client-secret",
    "JWT_SECRET": "your-jwt-secret"
  }'

# 4. Backend Stack
# GitHub Actions: Run "Backend Stack" workflow

# 5. Frontend Stack
# GitHub Actions: Run "Deploy Frontend" workflow
```

### IAC Makefile Commands

```bash
cd iac

make update-cf-defaults  # Update CloudFormation defaults from git remote
make init-cf             # Update defaults + show deployment instructions
make help                # Show available commands
```

## 🐳 Docker Commands

### Basic Operations

```bash
# Start all services
cd backend && docker compose up -d
cd frontend && docker compose up -d

# Stop all services
cd backend && docker compose down
cd frontend && docker compose down

# Stop with volume cleanup
docker compose down -v

# View logs
cd backend && docker compose logs -f
cd frontend && docker compose logs -f

# Rebuild containers
cd frontend && docker compose build --no-cache
```

### Container Management

```bash
# List running containers
docker ps

# List all containers
docker ps -a

# Inspect a container
docker inspect go-google-auth-frontend

# Execute command in container
docker exec -it go-google-auth-frontend bash

# Remove stopped containers
docker container prune

# View container logs
docker logs go-google-auth-app
```

### Volume Management

```bash
# List volumes
docker volume ls

# Remove unused volumes
docker volume prune

# Inspect volume
docker volume inspect frontend_node_modules
```

## 🔐 Environment Variables

### Backend Configuration

The backend uses environment variables for configuration:

- **Local Development**: `.env` file in `backend/` directory
- **Staging/Production**: AWS ECS Task Definition, Kubernetes Secrets, or CI/CD pipeline

> ⚠️ **Note:** `.env` files are for **local development only**. Never commit `.env` files to version control.

#### Required Environment Variables

```env
# Application Configuration
GO_ENV=development                # Options: development, staging, production
PORT=8080                        # Server port (default: 8080)

# CORS Configuration
ALLOWED_ORIGINS=http://localhost:5173,https://yourdomain.com  # Comma-separated list
FRONTEND_URL=http://localhost:5173                            # Main frontend URL

# AWS Configuration
AWS_REGION=ap-northeast-1
DYNAMODB_ENDPOINT=http://dynamodb:8000  # Local only
AWS_ACCESS_KEY_ID=dummy                 # Local only
AWS_SECRET_ACCESS_KEY=dummy             # Local only

# Google OAuth (required for authentication features)
GOOGLE_CLIENT_ID=your_google_client_id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your_google_client_secret

# JWT Configuration
JWT_SECRET=your-secret-key        # Generate with: openssl rand -base64 32
```

#### Environment-Specific Configuration

**Local Development (`.env` file):**
```env
ALLOWED_ORIGINS=http://localhost:5173
FRONTEND_URL=http://localhost:5173
```

**Staging/Production (AWS Secrets Manager):**

Secrets are stored in AWS Secrets Manager and automatically injected into ECS containers.

#### Option 1: Deploy Secrets Stack (Recommended)

The project includes a CDK secrets stack that creates the secret with default placeholder values:

```bash
cd iac

# Deploy secrets stack (creates secret with placeholder values)
npx cdk deploy --app "npx ts-node --prefer-ts-exts bin/secrets.ts" \
  --context projectName=go-google-auth \
  --context environment=dev

# After deployment, update with your actual credentials
aws secretsmanager put-secret-value \
  --secret-id "go-google-auth/dev/google-auth" \
  --secret-string '{
    "GOOGLE_CLIENT_ID": "your-client-id.apps.googleusercontent.com",
    "GOOGLE_CLIENT_SECRET": "your-client-secret",
    "JWT_SECRET": "your-jwt-secret"
  }'
```

#### Option 2: Create Secret Manually

```bash
# Create secret manually
aws secretsmanager create-secret \
  --name "go-google-auth/dev/google-auth" \
  --secret-string '{
    "GOOGLE_CLIENT_ID": "your-client-id.apps.googleusercontent.com",
    "GOOGLE_CLIENT_SECRET": "your-client-secret",
    "JWT_SECRET": "your-jwt-secret-generated-with-openssl"
  }'
```

#### Secret Structure

| Key | Description |
|-----|-------------|
| `GOOGLE_CLIENT_ID` | Google OAuth Client ID |
| `GOOGLE_CLIENT_SECRET` | Google OAuth Client Secret |
| `JWT_SECRET` | Secret key for JWT token signing |

> 💡 **Note:** The backend CDK stack automatically references secrets from `{projectName}/{environment}/google-auth`

### Frontend Configuration

The frontend uses Vite environment variables:

- **Local Development**: `.env.development` file in `frontend/vue-app/`
- **Staging/Production**: Build-time environment variables via CI/CD pipeline

> ⚠️ **Note:** `.env.*` files are for **local development only**. For production builds, set environment variables in your CI/CD pipeline.

#### `.env.development` (Local Development Only)
```env
VITE_BACKEND_URL=http://localhost:8080
VITE_PORT=5173
VITE_APP_ENV=development
VITE_GOOGLE_CLIENT_ID=your-client-id.apps.googleusercontent.com
```

#### Staging/Production (CI/CD Environment Variables)

Set these in your CI/CD pipeline (GitHub Actions, etc.):

```bash
# Staging
VITE_BACKEND_URL=https://api.staging.yourdomain.com
VITE_APP_ENV=staging
VITE_GOOGLE_CLIENT_ID=your-staging-client-id.apps.googleusercontent.com

# Production
VITE_BACKEND_URL=https://api.yourdomain.com
VITE_APP_ENV=production
VITE_GOOGLE_CLIENT_ID=your-production-client-id.apps.googleusercontent.com
```

#### Available Frontend Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `VITE_BACKEND_URL` | Backend API URL | `http://localhost:8080` |
| `VITE_PORT` | Development server port | `5173` |
| `VITE_APP_ENV` | Application environment | `development` |
| `VITE_GOOGLE_CLIENT_ID` | Google OAuth Client ID | (required) |

#### Using Environment Variables in Vue Components

```typescript
// Access environment variables
const backendUrl = import.meta.env.VITE_BACKEND_URL
const appEnv = import.meta.env.VITE_APP_ENV

console.log('Backend URL:', backendUrl)
```

### Docker Compose Environment Variables

Both `backend/compose.yml` and `frontend/compose.yml` support environment variable substitution:

```bash
# Set environment variables before starting
export ALLOWED_ORIGINS="http://localhost:5173,http://localhost:3000"
export VITE_BACKEND_URL="http://localhost:8080"

# Start services
cd backend && docker compose up -d
cd frontend && docker compose up -d
```

### CI/CD Environment Variables (GitHub Actions)

Configure these secrets in your GitHub repository:

- `AWS_ROLE_ARN` - AWS IAM role for OIDC authentication
- `PROJECT_NAME` - Your project name
- `ROOT_DOMAIN` - Your root domain (e.g., `example.com`)

> 💡 **Note:** Google OAuth credentials and JWT secret are stored in **AWS Secrets Manager**, not in GitHub secrets. The ECS task automatically retrieves them at runtime.

## 🐛 Troubleshooting

### Frontend Not Starting

**Symptom:** Container exits immediately or logs show npm errors

**Solution:**
```bash
cd frontend
docker compose down -v           # Remove volumes
docker compose build --no-cache  # Rebuild image
docker compose up -d             # Start container
```

### Permission Errors in Frontend

**Symptom:** `EACCES: permission denied` errors

**Solution:** Already handled in the Dockerfile with proper user permissions. If issues persist:
```bash
cd frontend
docker compose down -v
docker compose up -d
```

### CORS Errors

**Symptom:** Browser console shows CORS policy errors

**Solution:**
1. Verify backend CORS configuration includes `http://localhost:5173`
2. Ensure frontend uses `credentials: 'include'` in fetch requests
3. Check `Access-Control-Allow-Credentials` header is `true`

### Cookie Not Being Set

**Symptom:** Set-Cookie header sent but cookie not stored in browser

**Checklist:**
- ✅ Backend returns `Access-Control-Allow-Credentials: true`
- ✅ Frontend uses `credentials: 'include'` in fetch
- ✅ Frontend and backend are on allowed domains
- ✅ Cookie flags are appropriate (HttpOnly, Secure in production)

### Port Already in Use

**Symptom:** Cannot bind to port (8080, 5173, or 8000)

**Solution:**
```bash
# Find process using the port
lsof -i :8080

# Kill the process
kill -9 <PID>

# Or stop Docker containers
docker ps | grep go-google-auth
docker stop <CONTAINER_ID>
```

### Docker Build Fails

**Symptom:** Build context transfer errors or build failures

**Solution:**
```bash
# Clean Docker system
docker system prune -a

# Rebuild from scratch
cd frontend
docker compose build --no-cache --pull
```

## 🤝 Contributing

Contributions are welcome! Please follow these guidelines:

1. **Fork the repository**
2. **Create a feature branch**
   ```bash
   git checkout -b feature/amazing-feature
   ```
3. **Commit your changes**
   ```bash
   git commit -m 'Add amazing feature'
   ```
4. **Push to the branch**
   ```bash
   git push origin feature/amazing-feature
   ```
5. **Open a Pull Request**

### Code Style

- **Backend (Go)**: Use `gofmt` and follow [Effective Go](https://go.dev/doc/effective_go)
- **Frontend (Vue/TS)**: Use Prettier and ESLint configurations provided

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👤 Author

**yuki5155**

## 🔗 Resources

- [Go Documentation](https://go.dev/doc/)
- [Gin Framework](https://gin-gonic.com/)
- [Vue.js 3 Documentation](https://vuejs.org/)
- [Vite Documentation](https://vitejs.dev/)
- [Docker Documentation](https://docs.docker.com/)
- [AWS DynamoDB Local](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/DynamoDBLocal.html)

## 📊 Project Status

This project is fully functional with Google Identity Services (GIS) authentication implemented. The authentication flow includes:

- ✅ Google Sign-In with GIS library
- ✅ JWT-based session management (access + refresh tokens)
- ✅ Secure HttpOnly cookie storage
- ✅ Automatic token refresh
- ✅ Protected routes with authentication guards
- ✅ User dashboard with profile information

---

**Note:** This is a development project. For production deployment, ensure proper security configurations, use HTTPS, and enable appropriate cookie security flags.
