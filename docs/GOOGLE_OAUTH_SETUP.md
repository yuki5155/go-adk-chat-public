# Google OAuth Setup Guide for GCP

This guide will walk you through setting up Google OAuth 2.0 authentication for your application using Google Cloud Platform (GCP).

## Prerequisites

- Google Cloud SDK (gcloud CLI) installed and authenticated
- A Google account with GCP access
- Access to the [Google Cloud Console](https://console.cloud.google.com/)

## Table of Contents

1. [Create or Select a GCP Project](#1-create-or-select-a-gcp-project)
2. [Enable Required APIs](#2-enable-required-apis)
3. [Configure OAuth Consent Screen](#3-configure-oauth-consent-screen)
4. [Create OAuth 2.0 Credentials](#4-create-oauth-20-credentials)
5. [Configure Application Environment Variables](#5-configure-application-environment-variables)
6. [Verify Setup](#6-verify-setup)

---

## 1. Create or Select a GCP Project

You can either create a new project or use an existing one.

### Option A: Use an Existing Project

If you already have a GCP project you want to use:

#### Using gcloud CLI

```bash
# List all your available projects
gcloud projects list

# Set an existing project as active
gcloud config set project YOUR-EXISTING-PROJECT-ID

# Verify the project is set
gcloud config get-value project

# Check if you have necessary permissions
gcloud projects get-iam-policy YOUR-EXISTING-PROJECT-ID --flatten="bindings[].members" --filter="bindings.members:user:$(gcloud config get-value account)"
```

Replace `YOUR-EXISTING-PROJECT-ID` with one of your project IDs from the list (e.g., `glassy-rush-297801`).

#### Using Google Cloud Console

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Click on the project dropdown at the top
3. Select your existing project from the list

### Option B: Create a New Project

#### Using gcloud CLI

**Option 1: Auto-generate from repository name**

```bash
# Auto-generate project ID from repository name with timestamp
REPO_NAME=$(basename $(git rev-parse --show-toplevel 2>/dev/null || pwd))
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
PROJECT_ID="${REPO_NAME}-${TIMESTAMP}"
PROJECT_NAME=$(echo "$REPO_NAME" | sed 's/-/ /g' | awk '{for(i=1;i<=NF;i++) $i=toupper(substr($i,1,1)) tolower(substr($i,2));}1')

# Show what will be created
echo "Repository: $REPO_NAME"
echo "Project ID: $PROJECT_ID"
echo "Project Name: $PROJECT_NAME"
echo ""

# Create the project
gcloud projects create "$PROJECT_ID" --name="$PROJECT_NAME"

# Set as active project
gcloud config set project "$PROJECT_ID"

# Verify
gcloud config get-value project
```

**Example output:**
```
Repository: go-adk-chat
Project ID: go-adk-chat-20260101-143022
Project Name: Go Adk Chat
```

**Option 2: Simpler version (without timestamp)**

```bash
# Generate project ID from repository name
REPO_NAME=$(basename $(git rev-parse --show-toplevel 2>/dev/null || pwd))
PROJECT_ID="${REPO_NAME}-oauth"
PROJECT_NAME=$(echo "$REPO_NAME" | sed 's/-/ /g' | awk '{for(i=1;i<=NF;i++) $i=toupper(substr($i,1,1)) tolower(substr($i,2));}1')

# Create and set project
gcloud projects create "$PROJECT_ID" --name="$PROJECT_NAME OAuth" && \
gcloud config set project "$PROJECT_ID"
```

**Option 3: Manual specification**

```bash
# Create a new project with custom ID
gcloud projects create YOUR-PROJECT-ID --name="Your Project Name"

# Set the new project as active
gcloud config set project YOUR-PROJECT-ID

# Verify the project is set
gcloud config get-value project
```

Replace `YOUR-PROJECT-ID` with a unique project ID (e.g., `my-app-oauth-2025`).

#### Using Google Cloud Console

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Click on the project dropdown at the top
3. Click "New Project"
4. Enter project name and click "Create"
5. Select the newly created project from the dropdown

---

## 2. Enable Required APIs

For basic Google OAuth authentication (Sign in with Google), you typically **don't need to enable any APIs**. OAuth 2.0 works by default once you configure the consent screen and credentials.

### Optional: Identity Toolkit API (Advanced Features Only)

Only enable this if you're using advanced Firebase Authentication or Google Identity Platform features:

### Using gcloud CLI

```bash
# Optional: Only if using Identity Toolkit features
gcloud services enable identitytoolkit.googleapis.com

# Verify API is enabled
gcloud services list --enabled | grep identitytoolkit
```

### Using Google Cloud Console

1. Navigate to [APIs & Services > Library](https://console.cloud.google.com/apis/library)
2. Search for "Identity Toolkit API"
3. Click on it and click "Enable"

**Note**: For most web applications using standard Google OAuth (Sign in with Google), you can skip this step entirely and proceed directly to configuring the OAuth consent screen.

---

## 3. Configure OAuth Consent Screen

The OAuth consent screen is what users see when they're asked to authorize your application.

**Note**: This step must be done through the Google Cloud Console web interface. There is no CLI alternative for configuring the OAuth consent screen.

1. Go to [APIs & Services > OAuth consent screen](https://console.cloud.google.com/apis/credentials/consent)

2. **Click "Get started"** button to begin configuration

   You'll see a 4-step wizard: App Information → Audience → Contact Information → Finish

3. **Step 1: App Information**
   - **App name**: Your application name (e.g., "Go ADK Chat")
   - **User support email**: Select your email from the dropdown
   - Click "Save and Continue"

4. **Step 2: Audience** - Choose user type:
   - **Internal**: Only for Google Workspace users in your organization
   - **External**: For anyone with a Google account (recommended for most apps)
   - Select **"External"**
   - Click "Save and Continue"

5. **Step 3: Contact Information**
   - **Developer contact information**: Enter your email address
   - Click "Save and Continue"

6. **Step 4: Finish**
   - Review your settings
   - Click "Finish" or "Back to Dashboard"

### Configure Test Users and Scopes

After completing the initial setup, your app will be in **"Testing" mode** by default. You need to add test users:

1. **Go back to the OAuth consent screen page** if you're not already there
   - Navigate to [APIs & Services > OAuth consent screen](https://console.cloud.google.com/apis/credentials/consent)

2. **Add Test Users:**
   - Scroll down to the "Test users" section
   - Click "Add Users"
   - Enter the email addresses of Google accounts you want to allow for testing
   - Click "Save"
   - **Important**: Only these test users will be able to sign in while your app is in Testing mode

3. **Configure Scopes (Optional):**
   - Click "Edit App" or find the "Scopes" section
   - The default scopes (`openid`, `email`, `profile`) are sufficient for basic authentication
   - You can add additional scopes if your app needs more Google API access
   - Click "Save and Continue"

4. **Publishing Your App (When Ready for Production):**
   - While in Testing mode, only test users can sign in
   - To make your app available to all users:
     - Go to the OAuth consent screen page
     - Click "Publish App"
     - **Note**: Depending on the scopes you request, you may need to go through Google's verification process

---

## 4. Create OAuth 2.0 Credentials

Now you'll create the Client ID and Client Secret that your application will use.

**⚠️ IMPORTANT - No CLI Support**:
- This step **MUST** be done through the Google Cloud Console web interface
- There is **NO** gcloud CLI command to create OAuth credentials for web applications
- The gcloud CLI only supports IAP (Identity-Aware Proxy) OAuth clients, which are different from standard web app OAuth
- Web applications require custom redirect URIs, which can only be configured via the Console

1. Go to [APIs & Services > Credentials](https://console.cloud.google.com/apis/credentials)

2. Click "Create Credentials" > "OAuth client ID"

3. **Application type:**
   - Select "Web application"

4. **Name:**
   - Enter a name (e.g., "Go ADK Chat Web Client")

5. **Authorized JavaScript origins** (for frontend):
   - Click "Add URI"
   - Add your frontend URLs:
     ```
     http://localhost:5173
     https://yourdomain.com
     ```

6. **Authorized redirect URIs** (where Google sends users after login):
   - Click "Add URI"
   - Add your callback URLs:
     ```
     http://localhost:5173
     http://localhost:5173/auth/callback
     https://yourdomain.com
     https://yourdomain.com/auth/callback
     ```

7. Click "Create"

8. **Save Your Credentials:**
   - A modal will appear with your Client ID and Client Secret
   - **IMPORTANT**: Copy both values immediately
   - Client ID: `XXXXXXXXX.apps.googleusercontent.com`
   - Client Secret: `GOCSPX-XXXXXXXXX`

---

## 5. Configure Application Environment Variables

### Backend Configuration

#### Option A: Using CLI Commands (Automated)

After obtaining your Client ID and Client Secret from Step 4, run these commands:

```bash
# Set your credentials as variables (replace with your actual values)
export CLIENT_ID="YOUR-CLIENT-ID.apps.googleusercontent.com"
export CLIENT_SECRET="YOUR-CLIENT-SECRET"
export FRONTEND_URL="http://localhost:5173"
export ALLOWED_ORIGINS="http://localhost:5173,https://yourdomain.com"

# Copy the example file
cp backend/.env.example backend/.env

# Update the .env file with your credentials using sed
sed -i '' "s|GOOGLE_CLIENT_ID=.*|GOOGLE_CLIENT_ID=${CLIENT_ID}|g" backend/.env
sed -i '' "s|GOOGLE_CLIENT_SECRET=.*|GOOGLE_CLIENT_SECRET=${CLIENT_SECRET}|g" backend/.env
sed -i '' "s|FRONTEND_URL=.*|FRONTEND_URL=${FRONTEND_URL}|g" backend/.env
sed -i '' "s|ALLOWED_ORIGINS=.*|ALLOWED_ORIGINS=${ALLOWED_ORIGINS}|g" backend/.env

# Verify the changes
echo "Backend .env configured:"
grep -E "GOOGLE_CLIENT_ID|GOOGLE_CLIENT_SECRET|FRONTEND_URL|ALLOWED_ORIGINS" backend/.env
```

**Note for Linux users:** Remove the empty string `''` after `-i` flag:
```bash
sed -i "s|GOOGLE_CLIENT_ID=.*|GOOGLE_CLIENT_ID=${CLIENT_ID}|g" backend/.env
```

#### Option B: Manual Configuration

1. Copy the backend environment example file:
   ```bash
   cp backend/.env.example backend/.env
   ```

2. Edit `backend/.env` and update the following:
   ```env
   # Google OAuth - Replace with your actual credentials
   GOOGLE_CLIENT_ID=YOUR-CLIENT-ID.apps.googleusercontent.com
   GOOGLE_CLIENT_SECRET=YOUR-CLIENT-SECRET

   # Update other settings as needed
   ALLOWED_ORIGINS=http://localhost:5173,https://yourdomain.com
   FRONTEND_URL=http://localhost:5173
   ```

### Frontend Configuration

#### Option A: Using CLI Commands (Automated)

```bash
# Set your Client ID (use the same one from backend)
export CLIENT_ID="YOUR-CLIENT-ID.apps.googleusercontent.com"
export BACKEND_URL="http://localhost:8080"

# Copy the example file
cp frontend/vue-app/.env.example frontend/vue-app/.env.development

# Update the .env.development file
sed -i '' "s|VITE_GOOGLE_CLIENT_ID=.*|VITE_GOOGLE_CLIENT_ID=${CLIENT_ID}|g" frontend/vue-app/.env.development
sed -i '' "s|VITE_BACKEND_URL=.*|VITE_BACKEND_URL=${BACKEND_URL}|g" frontend/vue-app/.env.development

# Verify the changes
echo "Frontend .env configured:"
grep -E "VITE_GOOGLE_CLIENT_ID|VITE_BACKEND_URL" frontend/vue-app/.env.development
```

**Note for Linux users:** Remove the empty string `''` after `-i` flag:
```bash
sed -i "s|VITE_GOOGLE_CLIENT_ID=.*|VITE_GOOGLE_CLIENT_ID=${CLIENT_ID}|g" frontend/vue-app/.env.development
```

#### Option B: Manual Configuration

1. Copy the frontend environment example file:
   ```bash
   cp frontend/vue-app/.env.example frontend/vue-app/.env.development
   ```

2. Edit `frontend/vue-app/.env.development` and update:
   ```env
   # Google OAuth - Replace with your actual Client ID
   VITE_GOOGLE_CLIENT_ID=YOUR-CLIENT-ID.apps.googleusercontent.com

   # Backend API URL
   VITE_BACKEND_URL=http://localhost:8080
   ```

### Complete Setup Script (All-in-One)

For convenience, here's a complete script to set up both frontend and backend:

```bash
#!/bin/bash

# Set your OAuth credentials here
export CLIENT_ID="YOUR-CLIENT-ID.apps.googleusercontent.com"
export CLIENT_SECRET="YOUR-CLIENT-SECRET"

# Backend setup
echo "Setting up backend environment..."
cp backend/.env.example backend/.env
if [[ "$OSTYPE" == "darwin"* ]]; then
  # macOS
  sed -i '' "s|GOOGLE_CLIENT_ID=.*|GOOGLE_CLIENT_ID=${CLIENT_ID}|g" backend/.env
  sed -i '' "s|GOOGLE_CLIENT_SECRET=.*|GOOGLE_CLIENT_SECRET=${CLIENT_SECRET}|g" backend/.env
  sed -i '' "s|FRONTEND_URL=.*|FRONTEND_URL=http://localhost:5173|g" backend/.env
  sed -i '' "s|ALLOWED_ORIGINS=.*|ALLOWED_ORIGINS=http://localhost:5173|g" backend/.env
else
  # Linux
  sed -i "s|GOOGLE_CLIENT_ID=.*|GOOGLE_CLIENT_ID=${CLIENT_ID}|g" backend/.env
  sed -i "s|GOOGLE_CLIENT_SECRET=.*|GOOGLE_CLIENT_SECRET=${CLIENT_SECRET}|g" backend/.env
  sed -i "s|FRONTEND_URL=.*|FRONTEND_URL=http://localhost:5173|g" backend/.env
  sed -i "s|ALLOWED_ORIGINS=.*|ALLOWED_ORIGINS=http://localhost:5173|g" backend/.env
fi

# Frontend setup
echo "Setting up frontend environment..."
cp frontend/vue-app/.env.example frontend/vue-app/.env.development
if [[ "$OSTYPE" == "darwin"* ]]; then
  # macOS
  sed -i '' "s|VITE_GOOGLE_CLIENT_ID=.*|VITE_GOOGLE_CLIENT_ID=${CLIENT_ID}|g" frontend/vue-app/.env.development
  sed -i '' "s|VITE_BACKEND_URL=.*|VITE_BACKEND_URL=http://localhost:8080|g" frontend/vue-app/.env.development
else
  # Linux
  sed -i "s|VITE_GOOGLE_CLIENT_ID=.*|VITE_GOOGLE_CLIENT_ID=${CLIENT_ID}|g" frontend/vue-app/.env.development
  sed -i "s|VITE_BACKEND_URL=.*|VITE_BACKEND_URL=http://localhost:8080|g" frontend/vue-app/.env.development
fi

echo "✓ Environment configuration complete!"
echo ""
echo "Backend configuration:"
grep -E "GOOGLE_CLIENT_ID|FRONTEND_URL" backend/.env
echo ""
echo "Frontend configuration:"
grep -E "VITE_GOOGLE_CLIENT_ID|VITE_BACKEND_URL" frontend/vue-app/.env.development
```

Save this as `setup-env.sh`, update the CLIENT_ID and CLIENT_SECRET values, and run:
```bash
chmod +x setup-env.sh
./setup-env.sh
```

---

## 6. Verify Setup

### Check GCP Configuration

```bash
# Verify you're using the correct project
gcloud config get-value project

# Note: OAuth credentials for web apps must be viewed in the Cloud Console
# https://console.cloud.google.com/apis/credentials

# Optional: Check if Identity Toolkit API is enabled (only needed if using advanced features)
gcloud services list --enabled | grep identitytoolkit
```

### Test Locally

1. **Start the backend:**
   ```bash
   cd backend
   make run
   # Or: go run cmd/api/main.go
   ```

2. **Start the frontend:**
   ```bash
   cd frontend/vue-app
   npm run dev
   ```

3. **Test Google Login:**
   - Open your browser to `http://localhost:5173`
   - Click on the Google login button
   - You should see the Google OAuth consent screen
   - After authorizing, you should be redirected back to your app

### Common Issues

#### Issue: "Access blocked: This app's request is invalid"
- **Solution**: Make sure you've added your test user email in the OAuth consent screen's test users section (if app is in testing mode)

#### Issue: "Redirect URI mismatch"
- **Solution**: Verify that the redirect URI in your OAuth client credentials exactly matches the one your application is using
- Check both the protocol (http vs https) and the path

#### Issue: "API has not been used in project before"
- **Solution**: Make sure you've enabled the Google Identity Toolkit API and OAuth2 API

#### Issue: "The OAuth client was not found"
- **Solution**: Verify you're using the correct project ID and that credentials were created successfully

---

## Security Best Practices

1. **Never commit credentials to version control:**
   - Add `.env` files to `.gitignore`
   - Use `.env.example` files with placeholder values

2. **Use different credentials for development and production:**
   - Create separate OAuth clients for dev and prod
   - Use different redirect URIs

3. **Restrict your OAuth client:**
   - Only add necessary redirect URIs
   - Keep the list of authorized JavaScript origins minimal

4. **Keep your Client Secret secure:**
   - Never expose it in frontend code
   - Only use it in backend services
   - Rotate it periodically

5. **Review OAuth scopes:**
   - Only request the minimum scopes needed
   - Users are more likely to trust apps that request fewer permissions

---

## Quick Reference Commands

### Project Management

```bash
# List all your GCP projects
gcloud projects list

# Show detailed info about current project
gcloud projects describe $(gcloud config get-value project)

# Switch to an existing project
gcloud config set project YOUR-EXISTING-PROJECT-ID

# Check current active project
gcloud config get-value project

# Check your account
gcloud config get-value account

# List all gcloud configurations
gcloud config configurations list

# Check project permissions
gcloud projects get-iam-policy YOUR-PROJECT-ID --flatten="bindings[].members" --filter="bindings.members:user:$(gcloud config get-value account)"
```

### API Management

```bash
# Optional: Enable Identity Toolkit API (only if needed for advanced features)
gcloud services enable identitytoolkit.googleapis.com

# List all enabled services in current project
gcloud services list --enabled

# Check if Identity Toolkit API is enabled
gcloud services list --enabled --filter="name:identitytoolkit"

# List available APIs
gcloud services list --available
```

### OAuth Credentials

**Note**: OAuth credentials for web applications must be managed through the [Google Cloud Console](https://console.cloud.google.com/apis/credentials).

The gcloud CLI commands below only work for Identity-Aware Proxy (IAP) OAuth clients, not for general web application OAuth credentials:

```bash
# List IAP OAuth clients (not applicable for standard web app OAuth)
gcloud alpha iap oauth-clients list

# List OAuth brands (consent screen) - read-only
gcloud alpha iap oauth-brands list
```

### Quick Setup for Existing Project

```bash
# Set your existing project as active
gcloud config set project YOUR-EXISTING-PROJECT-ID && \
echo "Project: $(gcloud config get-value project)" && \
echo "Account: $(gcloud config get-value account)" && \
echo "Project configured successfully!"

# Optional: Enable Identity Toolkit API (only if using advanced features)
# gcloud services enable identitytoolkit.googleapis.com
```

---

## Resources

- [Google Cloud Console](https://console.cloud.google.com/)
- [Google OAuth 2.0 Documentation](https://developers.google.com/identity/protocols/oauth2)
- [Google Identity Platform](https://cloud.google.com/identity-platform)
- [OAuth Playground (for testing)](https://developers.google.com/oauthplayground/)

---

## Next Steps

After completing this setup:

1. Test the authentication flow thoroughly
2. Configure production OAuth credentials when ready to deploy
3. Set up proper error handling for OAuth failures
4. Consider implementing refresh token rotation
5. Set up monitoring and logging for authentication events

For production deployment, remember to:
- Publish your OAuth consent screen (if not published yet)
- Use HTTPS for all redirect URIs
- Set up proper environment variable management (e.g., AWS Secrets Manager, GCP Secret Manager)
- Enable Google Cloud audit logging
