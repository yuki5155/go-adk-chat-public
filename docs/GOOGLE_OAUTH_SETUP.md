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

```bash
# Create a new project
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

You need to enable the Google Identity services API for OAuth authentication.

### Using gcloud CLI

```bash
# Enable the required APIs
gcloud services enable identitytoolkit.googleapis.com
gcloud services enable oauth2.googleapis.com

# Verify APIs are enabled
gcloud services list --enabled | grep -E "identitytoolkit|oauth2"
```

### Using Google Cloud Console

1. Navigate to [APIs & Services > Library](https://console.cloud.google.com/apis/library)
2. Search for "Google Identity Toolkit API"
3. Click on it and click "Enable"
4. Search for "Google OAuth2 API"
5. Click on it and click "Enable"

---

## 3. Configure OAuth Consent Screen

The OAuth consent screen is what users see when they're asked to authorize your application.

### Using Google Cloud Console

1. Go to [APIs & Services > OAuth consent screen](https://console.cloud.google.com/apis/credentials/consent)

2. **Choose User Type:**
   - **Internal**: Only for Google Workspace users in your organization
   - **External**: For anyone with a Google account (recommended for most apps)
   - Select "External" and click "Create"

3. **App Information:**
   - **App name**: Your application name (e.g., "Go ADK Chat")
   - **User support email**: Select your email
   - **App logo**: (Optional) Upload your app logo
   - **App domain**: (Optional for development)
   - **Developer contact information**: Enter your email address
   - Click "Save and Continue"

4. **Scopes:**
   - Click "Add or Remove Scopes"
   - Add the following scopes (or keep defaults for basic profile):
     - `openid`
     - `email`
     - `profile`
   - These are already included in the default scopes
   - Click "Update" and then "Save and Continue"

5. **Test Users (for External apps in testing):**
   - Click "Add Users"
   - Add your Google account email addresses for testing
   - Click "Save and Continue"

6. **Summary:**
   - Review your settings
   - Click "Back to Dashboard"

**Note**: Your app will be in "Testing" mode by default, which limits users to only those you've added as test users. To make it public, you'll need to publish the app (click "Publish App" on the OAuth consent screen page).

---

## 4. Create OAuth 2.0 Credentials

Now you'll create the Client ID and Client Secret that your application will use.

### Using Google Cloud Console

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

### Using gcloud CLI (Alternative - creates desktop app credentials)

```bash
# Note: This creates a desktop app credential. For web apps, use the Console method above.
gcloud alpha iap oauth-clients create \
  --project=YOUR-PROJECT-ID \
  --display-name="Go ADK Chat"
```

**For web applications, it's recommended to use the Console method above to properly configure redirect URIs.**

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

# List OAuth client IDs
gcloud alpha iap oauth-clients list --project=$(gcloud config get-value project)

# Check enabled APIs
gcloud services list --enabled | grep -E "identity|oauth"
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
# Enable required APIs for OAuth
gcloud services enable identitytoolkit.googleapis.com oauth2.googleapis.com

# List all enabled services in current project
gcloud services list --enabled

# Check if specific API is enabled
gcloud services list --enabled --filter="name:identitytoolkit OR name:oauth2"

# List available APIs
gcloud services list --available
```

### OAuth Credentials

```bash
# View OAuth clients (if created via CLI)
gcloud alpha iap oauth-clients list

# List OAuth brands (consent screen)
gcloud alpha iap oauth-brands list

# Describe a specific OAuth client
gcloud alpha iap oauth-clients describe CLIENT_ID --brand=BRAND_ID
```

### Quick Setup for Existing Project

```bash
# Complete setup in one go
gcloud config set project YOUR-EXISTING-PROJECT-ID && \
gcloud services enable identitytoolkit.googleapis.com oauth2.googleapis.com && \
echo "Project: $(gcloud config get-value project)" && \
echo "Account: $(gcloud config get-value account)" && \
echo "APIs enabled successfully!"
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
