# Lambda Deployment Checklist

Use this checklist to verify everything is ready before deploying the Lambda stack.

## ✅ Prerequisites

### 1. Route53 Hosted Zone (REQUIRED)
Custom domain is mandatory for this project.

```bash
# Check if hosted zone exists
aws route53 list-hosted-zones --query "HostedZones[?Name=='yourdomain.com.']"
```

**Expected output:**
```json
[
    {
        "Id": "/hostedzone/Z1234567890ABC",
        "Name": "yourdomain.com.",
        ...
    }
]
```

❌ **If empty:** Create hosted zone first:
```bash
aws route53 create-hosted-zone \
  --name yourdomain.com \
  --caller-reference $(date +%s)
```

---

### 2. Google OAuth Credentials (REQUIRED)
You need Google OAuth 2.0 credentials from Google Cloud Console.

**Steps:**
1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create or select a project
3. Enable Google+ API
4. Go to "Credentials" → "Create Credentials" → "OAuth 2.0 Client ID"
5. Application type: "Web application"
6. Authorized JavaScript origins: `https://dev.yourdomain.com`
7. Authorized redirect URIs: `https://dev.yourdomain.com`
8. Save the **Client ID** and **Client Secret**

---

### 3. Deploy GoogleAuthSecretsStack (REQUIRED)

```bash
cd iac

# Deploy secrets stack
npx cdk deploy --app "npx ts-node --prefer-ts-exts bin/secrets.ts" \
  --context projectName=go-google-auth \
  --context environment=dev
```

**Then update with actual values:**
```bash
aws secretsmanager put-secret-value \
  --secret-id "go-google-auth/dev/google-auth" \
  --secret-string '{
    "GOOGLE_CLIENT_ID":"1234567890-abcdefg.apps.googleusercontent.com",
    "GOOGLE_CLIENT_SECRET":"GOCSPX-your-actual-secret",
    "JWT_SECRET":"your-super-secret-random-string-min-32-chars"
  }'
```

**Verify secret exists:**
```bash
aws secretsmanager describe-secret \
  --secret-id "go-google-auth/dev/google-auth"
```

---

### 4. Configure GitHub Secrets (for GitHub Actions)

If deploying via GitHub Actions, configure these secrets in your repository:

| Secret Name | Value | Example |
|-------------|-------|---------|
| `AWS_ROLE_TO_ASSUME` | IAM role ARN for OIDC | `arn:aws:iam::123456789012:role/GitHubActionsRole` |
| `PROJECT_NAME` | Project identifier | `go-google-auth` |
| `ROOT_DOMAIN` | Your domain | `yourdomain.com` |
| `AWS_REGION` | AWS region (optional) | `ap-northeast-1` |

**How to add secrets:**
1. Go to repository Settings → Secrets and variables → Actions
2. Click "New repository secret"
3. Add each secret

---

## 🚀 Deployment Steps

### Option 1: Local Deployment (Manual)

#### Step 1: Build Lambda Functions
```bash
cd backend
./scripts/build-lambda.sh all
```

**Expected output:**
```
=== Building Lambda Functions (ZIP) ===
✓ Built auth-google.zip (8.2M)
✓ Built auth-refresh.zip (8.1M)
✓ Built auth-logout.zip (8.0M)
✓ Built get-user.zip (8.3M)
✓ Built health.zip (7.9M)
✓ Built hello.zip (7.8M)

=== Build Summary ===
✓ Success: 6/6
```

#### Step 2: Deploy Lambda Stack
```bash
cd ../iac
npx cdk deploy --app "npx ts-node --prefer-ts-exts bin/lambda.ts" \
  --context projectName=go-google-auth \
  --context environment=dev \
  --context rootDomain=yourdomain.com \
  --context subdomain=lambda
```

**Expected output:**
```
=== Lambda Backend Configuration ===
API Domain: lambda.dev.yourdomain.com
✓ Custom domain configured: https://lambda.dev.yourdomain.com
✅ Successfully created go-google-auth-dev-lambda
✅ Deployed 6 Lambda functions
```

#### Step 3: Verify Deployment
```bash
# Check stack exists
aws cloudformation describe-stacks \
  --stack-name go-google-auth-dev-lambda

# Test health endpoint
curl https://lambda.dev.yourdomain.com/health
```

**Expected response:**
```json
{"status":"ok"}
```

---

### Option 2: GitHub Actions (Automated)

#### Step 1: Push to GitHub
```bash
git add .
git commit -m "Add Lambda deployment"
git push origin feature/12-add-lambda-handlers
```

#### Step 2: Run Lambda Stack Workflow
1. Go to **Actions** → **Lambda Stack**
2. Click **Run workflow**
3. Select:
   - **environment**: `dev`
   - **subdomain**: `lambda` (default)
4. Click **Run workflow**

#### Step 3: Wait for Deployment
Monitor the workflow progress. It will:
- ✓ Build all 6 Lambda functions (~30 seconds)
- ✓ Deploy Lambda stack (~3-5 minutes)
- ✓ Display deployment summary

**Expected output:**
```
=== Deployment Summary ===
Environment: dev
Subdomain: lambda
Region: ap-northeast-1
Lambda Functions: 6
API URL: https://lambda.dev.yourdomain.com

✅ Lambda deployment complete!
```

---

## 🧪 Testing

### 1. Test Health Endpoint
```bash
curl https://lambda.dev.yourdomain.com/health
```
**Expected:** `{"status":"ok"}`

### 2. Test Hello Endpoint
```bash
curl https://lambda.dev.yourdomain.com/hello
```
**Expected:** `{"message":"Hello World"}`

### 3. Test CORS
```bash
curl -H "Origin: https://dev.yourdomain.com" \
     -H "Access-Control-Request-Method: POST" \
     -H "Access-Control-Request-Headers: Content-Type" \
     -X OPTIONS \
     https://lambda.dev.yourdomain.com/auth/google
```
**Expected:** CORS headers in response

### 4. Test Authentication (requires frontend)
1. Deploy frontend with backend type: `lambda`
2. Open `https://dev.yourdomain.com`
3. Click "Sign in with Google"
4. Verify authentication works

---

## 🎯 Deploy Frontend to Use Lambda

Once Lambda is deployed and tested:

### Via GitHub Actions
1. Go to **Actions** → **Deploy Frontend**
2. Click **Run workflow**
3. Select:
   - **environment**: `dev`
   - **backendType**: `lambda` ⬅️ Important!
4. Click **Run workflow**

**Frontend will now connect to:** `https://lambda.dev.yourdomain.com`

---

## ✅ Final Verification Checklist

After deployment, verify:

- [ ] Route53 has A record for `lambda.dev.yourdomain.com`
- [ ] ACM certificate is validated
- [ ] API Gateway is accessible at custom domain
- [ ] All 6 Lambda functions are deployed
- [ ] Health endpoint responds: `GET /health`
- [ ] Hello endpoint responds: `GET /hello`
- [ ] CORS headers are present in responses
- [ ] Frontend can authenticate users
- [ ] JWT tokens are set in cookies
- [ ] Protected endpoint works: `GET /api/me`

---

## 🔍 Troubleshooting

### Issue: "rootDomain is required" error
**Cause:** Missing rootDomain parameter
**Fix:** Add `--context rootDomain=yourdomain.com`

### Issue: "Secret not found" error
**Cause:** GoogleAuthSecretsStack not deployed
**Fix:** Deploy secrets stack first (see step 3)

### Issue: "Hosted zone not found" error
**Cause:** Route53 hosted zone doesn't exist
**Fix:** Create hosted zone (see step 1)

### Issue: Certificate validation pending
**Cause:** DNS records not created yet
**Fix:** Wait 5-10 minutes for DNS propagation

### Issue: "Access Denied" on secret
**Cause:** Lambda role doesn't have permission
**Fix:** Check Lambda execution role has `secretsmanager:GetSecretValue` permission

---

## 📊 Cost Estimate

Lambda deployment cost (dev environment):
- **Lambda**: ~$0-5/month (within free tier for low traffic)
- **API Gateway**: ~$3.50/month per million requests
- **Route53**: ~$0.50/month per hosted zone
- **ACM Certificate**: FREE

**Total estimated cost: $0-10/month** (significantly cheaper than ECS)

---

## 🎉 You're Ready!

If all items in the prerequisites section are ✅, you're ready to deploy!

**Quick deploy command:**
```bash
cd backend && make -f Makefile.lambda deploy ENV=dev
```

This will build and deploy everything in one command.
