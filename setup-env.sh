#!/bin/bash

# Google OAuth Environment Setup Script
# This script configures the .env files for both backend and frontend

set -e  # Exit on error

echo "==================================="
echo "Google OAuth Environment Setup"
echo "==================================="
echo ""

# Check if CLIENT_ID and CLIENT_SECRET are provided
if [ -z "$CLIENT_ID" ] || [ -z "$CLIENT_SECRET" ]; then
  echo "Please set your OAuth credentials first:"
  echo ""
  echo "  export CLIENT_ID=\"YOUR-CLIENT-ID.apps.googleusercontent.com\""
  echo "  export CLIENT_SECRET=\"YOUR-CLIENT-SECRET\""
  echo ""
  echo "Then run this script again: ./setup-env.sh"
  echo ""
  echo "Alternatively, you can run:"
  echo "  CLIENT_ID=\"your-id\" CLIENT_SECRET=\"your-secret\" ./setup-env.sh"
  exit 1
fi

# Backend setup
echo "📦 Setting up backend environment..."
if [ ! -f "backend/.env.example" ]; then
  echo "❌ Error: backend/.env.example not found"
  exit 1
fi

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

echo "✓ Backend .env created and configured"

# Frontend setup
echo "📦 Setting up frontend environment..."
if [ ! -f "frontend/vue-app/.env.example" ]; then
  echo "❌ Error: frontend/vue-app/.env.example not found"
  exit 1
fi

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

echo "✓ Frontend .env.development created and configured"
echo ""
echo "==================================="
echo "✓ Environment configuration complete!"
echo "==================================="
echo ""
echo "📄 Backend configuration (backend/.env):"
echo "-----------------------------------"
grep -E "GOOGLE_CLIENT_ID|FRONTEND_URL|ALLOWED_ORIGINS" backend/.env | sed 's/GOOGLE_CLIENT_SECRET=.*/GOOGLE_CLIENT_SECRET=***hidden***/g'
echo ""
echo "📄 Frontend configuration (frontend/vue-app/.env.development):"
echo "-----------------------------------"
grep -E "VITE_GOOGLE_CLIENT_ID|VITE_BACKEND_URL" frontend/vue-app/.env.development
echo ""
echo "🚀 Next steps:"
echo "  1. Review the .env files to ensure all values are correct"
echo "  2. Start the backend: cd backend && make run"
echo "  3. Start the frontend: cd frontend/vue-app && npm run dev"
echo "  4. Test Google login at http://localhost:5173"
echo ""
