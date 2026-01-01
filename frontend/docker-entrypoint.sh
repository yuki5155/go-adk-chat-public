#!/bin/bash
set -e

echo "Checking for node_modules..."

# Check if node_modules exists and if package.json has been modified
if [ ! -d "node_modules" ]; then
  echo "node_modules not found. Installing dependencies..."
  npm install
elif [ "package.json" -nt "node_modules" ]; then
  echo "package.json is newer than node_modules. Reinstalling dependencies..."
  npm install
else
  echo "Dependencies already installed."
fi

echo "Starting application..."

# Execute the main command
exec "$@"
