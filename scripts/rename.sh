#!/bin/bash

# Script to replace "go-foundation-v2" with your service name across 
# the Go codebase only
# Will not touch YAML files, Dockerfiles, or GitHub workflow files

set -e  # Exit on error

# Check if an argument is provided
if [ -z "$1" ]; then
  echo -e "${RED}Error: No service name provided.${NC}"
  echo "Usage: $0 <new_service_name>"
  exit 1
fi

NEW_SERVICE_NAME=$1
OLD_MODULE_PATH="github.com/razorpay/go-foundation-v2"
NEW_MODULE_PATH="github.com/razorpay/${NEW_SERVICE_NAME}"
OLD_SERVICE_NAME="go-foundation-v2"

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Starting replacement of '${OLD_SERVICE_NAME}' with '${NEW_SERVICE_NAME}' in Go files only...${NC}"

# 1. First replace in go.mod file
echo -e "Updating ${GREEN}go.mod${NC}..."
sed -i '' "s|${OLD_MODULE_PATH}|${NEW_MODULE_PATH}|g" go.mod

# 2. Replace in all .go files
echo -e "Updating ${GREEN}Go source files${NC}..."
find . -name "*.go" -type f -print0 | xargs -0 sed -i '' "s|${OLD_MODULE_PATH}|${NEW_MODULE_PATH}|g"

# 3. Update git remote origin URL
NEW_REMOTE_URL="git@github.com:razorpay/${NEW_SERVICE_NAME}.git"
echo -e "Updating ${GREEN}git remote origin${NC} URL to ${NEW_REMOTE_URL}..."
git remote set-url origin "${NEW_REMOTE_URL}"

echo -e "${YELLOW}Replacement complete in Go files, go.mod, Makefile, and git remote origin URL.${NC}"
echo -e "${YELLOW}Note: You will need to manually update Dockerfiles and GitHub workflow files.${NC}"
echo -e "${YELLOW}Note: You may need to run 'go mod tidy' after this script to update dependencies.${NC}" 