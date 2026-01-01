# Deployment Guide

This guide explains how to deploy the Go Google Auth application to AWS using CDK.

## Prerequisites

- AWS Account with appropriate permissions
- AWS CLI configured
- Node.js 22.15.0+
- Docker Desktop (for local testing)
- GitHub repository with OIDC authentication configured

## Architecture Overview

### Frontend (Static Site)
- **S3 Bucket**: Hosts the built Vue.js application static files
- **CloudFront**: CDN for global distribution with caching and compression
- **Route53** (Optional): Custom domain configuration
- **Certificate Manager**: SSL/TLS certificates in us-east-1 (auto-detected or specified)

**Deployment Flow:**
1. Build Vue.js app: `npm run build` → generates `dist/` directory
2. Upload to S3: `aws s3 sync dist/ s3://bucket-name`
3. Invalidate CloudFront cache
4. CloudFront serves files globally with optimizations

### Backend (Container Service)
- **ECS Fargate**: Runs Go application containers
- **Application Load Balancer**: Routes traffic to ECS tasks
- **ECR**: Stores Docker images
- **Route53** (Optional): Custom domain for API
- **Certificate Manager**: SSL/TLS for HTTPS

### Shared Infrastructure
- **VPC**: Network isolation
- **DynamoDB**: Session storage
- **Secrets Manager**: Stores sensitive configuration

## Environment Setup

### 1. Configure GitHub Secrets

Add these secrets to your GitHub repository (`Settings` > `Secrets and variables` > `Actions`):

```
AWS_ROLE_ARN=arn:aws:iam::ACCOUNT_ID:role/GithubActionsRole
PROJECT_NAME=go-google-auth
ROOT_DOMAIN=yourdomain.com
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
```

### 2. AWS IAM Role for GitHub Actions

Create an IAM role with OIDC trust policy:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Federated": "arn:aws:iam::ACCOUNT_ID:oidc-provider/token.actions.githubusercontent.com"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringEquals": {
          "token.actions.githubusercontent.com:aud": "sts.amazonaws.com"
        },
        "StringLike": {
          "token.actions.githubusercontent.com:sub": "repo:your-username/go-google-auth:*"
        }
      }
    }
  ]
}
```

Attach policies:
- `AdministratorAccess` (or more restrictive policies based on your needs)

### 3. Domain Configuration (Optional)

If using a custom domain:

1. **Register domain** in Route53 or transfer existing domain
2. **Create hosted zone** (auto-created during deployment if using rootDomain)
3. **Certificate**: Certificates are auto-detected if they exist in your account

## Deployment Steps

### Step 1: Deploy Network Infrastructure

```bash
cd iac

# Development
npx cdk deploy \
  --app "npx ts-node --prefer-ts-exts bin/network.ts" \
  --context projectName=go-google-auth \
  --context environment=dev \
  --context costLevel=minimal

# Production
npx cdk deploy \
  --app "npx ts-node --prefer-ts-exts bin/network.ts" \
  --context projectName=go-google-auth \
  --context environment=prod \
  --context costLevel=standard
```

### Step 2: Deploy Backend

```bash
cd iac

# Build and push Docker image first
cd ../backend
docker build -t go-google-auth:latest -f dockers/Dockerfile.prod .

# Tag and push to ECR (replace ACCOUNT_ID and REGION)
aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin ACCOUNT_ID.dkr.ecr.ap-northeast-1.amazonaws.com
docker tag go-google-auth:latest ACCOUNT_ID.dkr.ecr.ap-northeast-1.amazonaws.com/go-google-auth:latest
docker push ACCOUNT_ID.dkr.ecr.ap-northeast-1.amazonaws.com/go-google-auth:latest

# Deploy backend stack
cd ../iac
npx cdk deploy \
  --app "npx ts-node --prefer-ts-exts bin/backend.ts" \
  --context projectName=go-google-auth \
  --context environment=dev \
  --context rootDomain=yourdomain.com
```

### Step 3: Deploy Frontend

```bash
cd iac

# Development
npx cdk deploy \
  --app "npx ts-node --prefer-ts-exts bin/frontend.ts" \
  --context projectName=go-google-auth \
  --context environment=dev \
  --context costLevel=minimal \
  --context rootDomain=yourdomain.com

# Production
npx cdk deploy \
  --app "npx ts-node --prefer-ts-exts bin/frontend.ts" \
  --context projectName=go-google-auth \
  --context environment=prod \
  --context costLevel=standard \
  --context rootDomain=yourdomain.com
```

### Step 4: Upload Frontend Files

```bash
# Build frontend with production settings
cd frontend/vue-app
VITE_BACKEND_URL=https://api.yourdomain.com npm run build

# Get S3 bucket name from CloudFormation
BUCKET_NAME=$(aws cloudformation describe-stacks \
  --stack-name go-google-auth-prod-frontend \
  --query 'Stacks[0].Outputs[?OutputKey==`BucketName`].OutputValue' \
  --output text)

# Upload to S3 with optimized cache settings
aws s3 sync dist s3://$BUCKET_NAME \
  --delete \
  --cache-control "public, max-age=31536000, immutable" \
  --exclude "index.html" \
  --exclude "*.json"

# Upload index.html with no-cache (for SPA routing)
aws s3 cp dist/index.html s3://$BUCKET_NAME/index.html \
  --cache-control "no-cache, no-store, must-revalidate" \
  --content-type "text/html"

# Invalidate CloudFront cache to serve new files immediately
DISTRIBUTION_ID=$(aws cloudformation describe-stacks \
  --stack-name go-google-auth-prod-frontend \
  --query 'Stacks[0].Outputs[?OutputKey==`CloudFrontDistributionId`].OutputValue' \
  --output text)

aws cloudfront create-invalidation \
  --distribution-id $DISTRIBUTION_ID \
  --paths "/*"
```

**Note:** The FrontendStack (from `automation-deploy-template-iac`) automatically configures:
- CloudFront error pages for SPA routing (404 → index.html)
- Appropriate cache behaviors for static assets
- Gzip/Brotli compression
- Security headers via CloudFront functions

## CI/CD with GitHub Actions

### Automatic Deployments

The project includes GitHub Actions workflows for automated deployments:

- **`.github/workflows/network.yml`**: Deploys network infrastructure
- **`.github/workflows/backend.yml`**: Builds and deploys backend
- **`.github/workflows/frontend.yml`**: Builds and deploys frontend

### Triggering Deployments

**Automatic:**
- Push to `main` branch → Production deployment
- Push to `develop` branch → Development deployment

**Manual:**
- Go to `Actions` tab in GitHub
- Select the workflow
- Click `Run workflow`
- Choose environment (dev/staging/prod)

## Environment-Specific Configuration

### Development
```bash
# Context
projectName=go-google-auth
environment=dev
costLevel=minimal
rootDomain=yourdomain.com

# URLs
Frontend: https://dev.yourdomain.com
Backend: https://api.dev.yourdomain.com
```

### Staging
```bash
# Context
projectName=go-google-auth
environment=staging
costLevel=standard
rootDomain=yourdomain.com

# URLs
Frontend: https://staging.yourdomain.com
Backend: https://api.staging.yourdomain.com
```

### Production
```bash
# Context
projectName=go-google-auth
environment=prod
costLevel=standard
rootDomain=yourdomain.com

# URLs
Frontend: https://yourdomain.com
Backend: https://api.yourdomain.com
```

## Post-Deployment Tasks

### 1. Configure DNS

If using Route53 and a custom domain:
- Name servers are automatically configured
- Wait for DNS propagation (up to 48 hours, usually much faster)

### 2. Verify SSL Certificates

Certificates are auto-detected or created:
```bash
aws acm list-certificates --region us-east-1  # CloudFront certificates
aws acm list-certificates --region ap-northeast-1  # ALB certificates
```

### 3. Test Endpoints

```bash
# Frontend
curl https://dev.yourdomain.com

# Backend health check
curl https://api.dev.yourdomain.com/health

# Backend API
curl https://api.dev.yourdomain.com/api/set-cookie
```

### 4. Monitor Resources

- **CloudWatch Logs**: Check application logs
- **CloudWatch Metrics**: Monitor ECS tasks, ALB, CloudFront
- **X-Ray** (if enabled): Trace requests

## Updating Deployments

### Backend Updates

1. **Make code changes**
2. **Push to GitHub** → Automatically builds and deploys
3. **Or manually deploy:**
   ```bash
   cd backend
   docker build -t go-google-auth:latest -f dockers/Dockerfile.prod .
   # Tag and push to ECR
   # ECS will automatically deploy new image
   ```

### Frontend Updates

1. **Make code changes**
2. **Push to GitHub** → Automatically builds and deploys
3. **Or manually deploy:**
   ```bash
   cd frontend/vue-app
   VITE_BACKEND_URL=https://api.yourdomain.com npm run build
   aws s3 sync dist s3://$BUCKET_NAME --delete
   aws cloudfront create-invalidation --distribution-id $DISTRIBUTION_ID --paths "/*"
   ```

### Infrastructure Updates

```bash
cd iac
npx cdk deploy --all \
  --context projectName=go-google-auth \
  --context environment=prod \
  --context rootDomain=yourdomain.com
```

## Rollback Procedures

### Backend Rollback

1. **Find previous task definition:**
   ```bash
   aws ecs list-task-definitions --family-prefix go-google-auth-prod-backend
   ```

2. **Update service:**
   ```bash
   aws ecs update-service \
     --cluster go-google-auth-prod-backend-cluster \
     --service go-google-auth-prod-backend-service \
     --task-definition go-google-auth-prod-backend:REVISION
   ```

### Frontend Rollback

1. **Get previous S3 version** (if versioning enabled)
2. **Or restore from git:**
   ```bash
   git checkout PREVIOUS_COMMIT
   cd frontend/vue-app
   npm run build
   # Upload to S3
   ```

## Troubleshooting

### Deployment Fails

**CDK Stack Creation Failed:**
```bash
# Check CloudFormation events
aws cloudformation describe-stack-events --stack-name go-google-auth-prod-frontend

# View CDK diff
npx cdk diff --context projectName=go-google-auth --context environment=prod
```

### CORS Errors in Production

1. Check backend `ALLOWED_ORIGINS` environment variable
2. Verify CloudFront domain is included
3. Check ALB health checks

### SSL Certificate Issues

```bash
# Check certificate status
aws acm describe-certificate --certificate-arn CERT_ARN

# Certificate must be in us-east-1 for CloudFront
aws acm list-certificates --region us-east-1
```

### ECS Tasks Not Starting

```bash
# Check task logs
aws ecs describe-tasks --cluster CLUSTER_NAME --tasks TASK_ARN

# Check CloudWatch logs
aws logs tail /aws/ecs/go-google-auth-prod-backend --follow
```

## Cost Optimization

### Development Environment
- Use `costLevel=minimal`
- Single AZ deployment
- Smaller ECS task sizes
- Consider stopping non-essential resources overnight

### Production Environment
- Use `costLevel=standard` or `high-availability`
- Multi-AZ deployment
- Enable auto-scaling
- Use Reserved Instances for predictable workloads

## Security Best Practices

1. **Use Secrets Manager** for sensitive data
2. **Enable CloudTrail** for audit logs
3. **Use WAF** for CloudFront (optional)
4. **Enable VPC Flow Logs**
5. **Regular security scans** with AWS Inspector
6. **Keep dependencies updated**
7. **Use least-privilege IAM policies**

## Cleanup

To destroy all resources:

```bash
cd iac

# Destroy in reverse order
npx cdk destroy go-google-auth-prod-frontend
npx cdk destroy go-google-auth-prod-backend
npx cdk destroy go-google-auth-prod-network

# Empty S3 buckets manually if needed
aws s3 rm s3://BUCKET_NAME --recursive
```

## Support

For issues or questions:
- Check [GitHub Issues](https://github.com/yuki5155/go-google-auth/issues)
- Review [AWS CDK Documentation](https://docs.aws.amazon.com/cdk/)
- Consult [Troubleshooting](#troubleshooting) section

---

**Last Updated:** December 2025
