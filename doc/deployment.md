# Deployment Guide

This guide covers deploying go-adk-chat to AWS using GitHub Actions workflows.

The project supports two backend deployment models:

- **Lambda** (recommended) -- Serverless Go functions behind API Gateway with SSE streaming
- **ECS Fargate** -- Containerized backend with Application Load Balancer

## Prerequisites

- AWS account with a Route53 hosted zone for your domain
- Google OAuth 2.0 credentials ([console](https://console.cloud.google.com/apis/credentials))
- Google AI API key for Gemini ([ai.google.dev](https://ai.google.dev/))

## Deployment Order

Stacks must be deployed in this order due to dependencies:

```
1. GitHub Actions IAM Role  (one-time setup)
2. Secrets Stack
3. DynamoDB Stack
4. Network Stack
5. Backend: Lambda Stack OR ECS Stack
6. Frontend Stack
```

---

## 1. GitHub Actions IAM Role (One-time)

Create an OIDC IAM role so GitHub Actions can deploy to your AWS account.

```bash
cd iac
make update-cf-defaults
```

Then deploy the role via the AWS Console or CLI:

```bash
aws cloudformation create-stack \
  --stack-name go-adk-chat-github-actions \
  --template-body file://iac/cloudformations/github-actions-role.yaml \
  --capabilities CAPABILITY_NAMED_IAM \
  --parameters \
    ParameterKey=GitHubOrganization,ParameterValue=YOUR_GITHUB_ORG \
    ParameterKey=GitHubRepository,ParameterValue=go-adk-chat
```

After the stack completes, configure these GitHub repository secrets (Settings > Secrets and variables > Actions):

| Secret | Description | Example |
|--------|-------------|---------|
| `AWS_ROLE_TO_ASSUME` | IAM role ARN from the CloudFormation output | `arn:aws:iam::123456789012:role/...` |
| `AWS_REGION` | AWS region | `ap-northeast-1` |
| `PROJECT_NAME` | Project identifier | `go-adk-chat` |
| `ROOT_DOMAIN` | Your domain | `yourdomain.com` |

---

## 2. Secrets Stack

Creates an AWS Secrets Manager secret with placeholder values.

**Workflow:** Run the **Secrets Stack** workflow from the Actions tab.

- Input: `environment` (dev / stg / prod)

After the workflow completes, populate the secret with your actual credentials:

```bash
aws secretsmanager put-secret-value \
  --secret-id "go-adk-chat/dev/google-auth" \
  --secret-string '{
    "GOOGLE_CLIENT_ID": "your-id.apps.googleusercontent.com",
    "GOOGLE_CLIENT_SECRET": "your-secret",
    "JWT_SECRET": "your-jwt-secret-min-32-chars",
    "ROOT_USER_EMAIL": "admin@yourdomain.com",
    "GOOGLE_AI_API_KEY": "your-gemini-api-key"
  }'
```

---

## 3. DynamoDB Stack

Creates 6 DynamoDB tables for users, roles, and chat data.

**Workflow:** Run the **DynamoDB Stack** workflow from the Actions tab.

- Input: `environment` (dev / staging / prod)

### Tables created

| Table | Partition Key | Sort Key |
|-------|--------------|----------|
| `{project}-{env}-user-roles` | user_id | - |
| `{project}-{env}-role-requests` | request_id | - |
| `{project}-{env}-chat-threads` | user_id | thread_id |
| `{project}-{env}-chat-sessions` | thread_id | session_id |
| `{project}-{env}-chat-events` | session_id | event_id |
| `{project}-{env}-chat-memories` | thread_id | memory_id |

---

## 4. Network Stack

Creates VPC, subnets, NAT gateways, and security groups.

**Workflow:** Run the **Network Stack** workflow from the Actions tab.

- Inputs: `environment`, `cost-level`, `action`

### Cost levels

| Level | Description |
|-------|-------------|
| `minimal` | Single AZ, minimum resources |
| `standard` | Multi-AZ (2-3 AZs) |
| `high-availability` | Full redundancy across all AZs |

---

## 5a. Lambda Stack (Recommended)

Deploys 18 Lambda functions behind API Gateway with SSE streaming support for chat.

**Workflow:** Run the **Lambda Stack** workflow from the Actions tab.

- Inputs: `environment`, `subdomain` (default: `lambda`)

The workflow automatically builds all Go Lambda functions and deploys them.

### Endpoints deployed

| Group | Endpoints |
|-------|-----------|
| Auth | `POST /auth/google`, `POST /auth/refresh`, `POST /auth/logout` |
| User | `GET /api/me`, `POST /api/role/request` |
| Admin | `GET /admin/dashboard`, `GET /api/admin/role-requests`, `POST /api/admin/role/approve`, `POST /api/admin/role/reject`, `GET /api/admin/users` |
| Chat | `GET /api/chat/models`, `POST /api/chat/threads`, `GET /api/chat/threads`, `GET /api/chat/threads/{id}`, `DELETE /api/chat/threads/{id}`, `POST /api/chat/threads/{id}/message`, `POST /api/chat/threads/{id}/stream` |
| Health | `GET /health`, `GET /hello` |

### Verify

```bash
curl https://lambda.dev.yourdomain.com/health
# {"status":"ok"}
```

---

## 5b. ECS Fargate Stack (Alternative)

Deploys the backend as a containerized service with an Application Load Balancer.

**Workflow:** Run the **Backend Stack** workflow from the Actions tab.

- Input: `environment`

The workflow automatically builds the Docker image, pushes it to ECR, and deploys the ECS service.

> If the ECR repository doesn't exist yet, run the **Backend Setup** workflow first to create it.

---

## 6. Frontend Stack

Deploys the Vue.js SPA to S3 with CloudFront distribution.

**Workflow:** Run the **Frontend** workflow from the Actions tab.

- Inputs: `environment`, `backendType` (`ecs` or `lambda`)

The workflow automatically:
1. Fetches `GOOGLE_CLIENT_ID` from Secrets Manager
2. Builds the frontend with correct environment variables
3. Deploys the S3 + CloudFront infrastructure
4. Uploads files with proper cache headers
5. Invalidates the CloudFront cache

The workflow also triggers automatically on push:
- Push to `main` → production deployment
- Push to `develop` → development deployment

---

## Domain Structure

| Environment | Frontend | Backend (Lambda) | Backend (ECS) |
|-------------|----------|-------------------|---------------|
| dev | `dev.yourdomain.com` | `lambda.dev.yourdomain.com` | `api.dev.yourdomain.com` |
| stg | `stg.yourdomain.com` | `lambda.stg.yourdomain.com` | `api.stg.yourdomain.com` |
| prod | `yourdomain.com` | `lambda.yourdomain.com` | `api.yourdomain.com` |

---

## Rollback

### Lambda / ECS

Re-run the corresponding workflow from a previous commit on the Actions tab.

### Frontend

Re-run the **Frontend** workflow from a previous commit, or push a revert commit to trigger automatic deployment.

---

## Monitoring

### CloudWatch Logs

| Service | Log Group |
|---------|-----------|
| Lambda | `/aws/lambda/{project}-{env}-lambda-{function}` |
| ECS | `/aws/ecs/{project}-{env}-backend` |
| API Gateway | Configured in stage settings |

### Key Metrics

- **Lambda:** Invocations, Duration, Errors, Throttles
- **ECS:** CPU/Memory utilization, Task count
- **CloudFront:** Requests, Cache hit ratio, Error rate
