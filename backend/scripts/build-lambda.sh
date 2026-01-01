#!/bin/bash
# Build script for Lambda functions (ZIP deployment)
# Usage: ./scripts/build-lambda.sh [function-name]
# Example: ./scripts/build-lambda.sh auth-google
# Or build all: ./scripts/build-lambda.sh all

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Lambda functions to build
LAMBDA_FUNCTIONS=(
  "auth-google"
  "auth-refresh"
  "auth-logout"
  "get-user"
  "health"
  "hello"
)

# Build directory
BUILD_DIR="build/lambda"
CMD_DIR="cmd/lambda"

# Function to build a single Lambda function
build_lambda() {
  local func_name=$1
  echo -e "${YELLOW}Building Lambda function: ${func_name}${NC}"

  # Create build directory
  mkdir -p "${BUILD_DIR}/${func_name}"

  # Build the Go binary
  echo "  → Compiling Go binary..."
  GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
    -ldflags="-s -w" \
    -o "${BUILD_DIR}/${func_name}/bootstrap" \
    "./${CMD_DIR}/${func_name}"

  if [ $? -ne 0 ]; then
    echo -e "${RED}✗ Failed to build ${func_name}${NC}"
    return 1
  fi

  # Create ZIP file
  echo "  → Creating ZIP file..."
  cd "${BUILD_DIR}/${func_name}"
  zip -q ../../../"${BUILD_DIR}/${func_name}.zip" bootstrap
  cd ../../..

  # Clean up temporary directory
  rm -rf "${BUILD_DIR}/${func_name}"

  # Get file size
  local size=$(du -h "${BUILD_DIR}/${func_name}.zip" | cut -f1)
  echo -e "${GREEN}✓ Built ${func_name}.zip (${size})${NC}"

  return 0
}

# Main script
echo "=== Lambda Build Script (ZIP Deployment) ==="
echo ""

# Clean build directory
if [ -d "$BUILD_DIR" ]; then
  echo "Cleaning build directory..."
  rm -rf "$BUILD_DIR"
fi
mkdir -p "$BUILD_DIR"

# Check if building specific function or all
if [ "$1" = "all" ] || [ -z "$1" ]; then
  echo "Building all Lambda functions..."
  echo ""

  failed_builds=()
  for func in "${LAMBDA_FUNCTIONS[@]}"; do
    if ! build_lambda "$func"; then
      failed_builds+=("$func")
    fi
    echo ""
  done

  # Summary
  echo "=== Build Summary ==="
  total=${#LAMBDA_FUNCTIONS[@]}
  failed=${#failed_builds[@]}
  success=$((total - failed))

  echo -e "${GREEN}✓ Success: ${success}/${total}${NC}"

  if [ $failed -gt 0 ]; then
    echo -e "${RED}✗ Failed: ${failed}/${total}${NC}"
    echo "Failed functions: ${failed_builds[*]}"
    exit 1
  fi

  echo ""
  echo "All Lambda functions built successfully!"
  echo "ZIP files location: ${BUILD_DIR}/"
  ls -lh "${BUILD_DIR}"/*.zip

else
  # Build specific function
  func_name=$1

  # Check if function exists
  if [[ ! " ${LAMBDA_FUNCTIONS[@]} " =~ " ${func_name} " ]]; then
    echo -e "${RED}Error: Unknown Lambda function: ${func_name}${NC}"
    echo ""
    echo "Available functions:"
    for func in "${LAMBDA_FUNCTIONS[@]}"; do
      echo "  - $func"
    done
    exit 1
  fi

  build_lambda "$func_name"

  echo ""
  echo "Lambda function built successfully!"
  echo "ZIP file: ${BUILD_DIR}/${func_name}.zip"
  ls -lh "${BUILD_DIR}/${func_name}.zip"
fi

echo ""
echo "Next steps:"
echo "  cd ../iac"
echo "  cdk deploy --app 'npx ts-node --prefer-ts-exts bin/lambda.ts' \\"
echo "    --context environment=dev \\"
echo "    --context projectName=go-google-auth"
