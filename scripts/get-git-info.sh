#!/bin/bash

# Exit on any error
set -euo pipefail

# Color variables
YELLOW='\033[1;33m'
RED='\033[1;31m'
GREEN='\033[1;32m'
NC='\033[0m' # No Color

# Check if ENV_FILE parameter is provided
if [ $# -eq 0 ]; then
    printf "${RED}Error: ENV_FILE parameter is required${NC}\n"
    echo "Usage: $0 <ENV_FILE>"
    echo "Example: $0 .env.dev"
    exit 1
fi

ENV_FILE="$1"

# Check if we're in a git repository
if ! git rev-parse --git-dir >/dev/null 2>&1; then
    printf "${RED}Error: Not in a git repository${NC}\n"
    exit 1
fi

# Get git commit information
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(git rev-parse --short HEAD)
GIT_COMMIT_HASH=$(git rev-parse --short HEAD)
GIT_COMMIT_FULL=$(git rev-parse HEAD)
GIT_COMMIT_DATE=$(git log -1 --format=%cd --date=format:%Y%m%d%H%M%S)
GIT_COMMIT_DATE_SHORT=$(git log -1 --format=%cd --date=format:"%Y-%m-%d")
GIT_COMMIT_MESSAGE=$(git log -1 --format=%s)
# Escape for .env: backslashes and double quotes
ESCAPED_GIT_COMMIT_MESSAGE="$(printf '%s' "$GIT_COMMIT_MESSAGE" | sed -e 's/\\/\\\\/g' -e 's/"/\\"/g')"
GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
GIT_TAG=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_DIRTY=$(git status --porcelain | wc -l | tr -d ' ')
BUILD_TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Output environment variables
echo
echo "Git commit information:"
printf "${YELLOW}GIT_COMMIT_HASH=%s${NC}\n" "$GIT_COMMIT_HASH"
printf "${YELLOW}GIT_COMMIT_FULL=%s${NC}\n" "$GIT_COMMIT_FULL"
printf "${YELLOW}GIT_COMMIT_DATE=%s${NC}\n" "$GIT_COMMIT_DATE"
printf "${YELLOW}GIT_COMMIT_DATE_SHORT=%s${NC}\n" "$GIT_COMMIT_DATE_SHORT"
printf "${YELLOW}GIT_COMMIT_MESSAGE=\"%s\"${NC}\n" "$ESCAPED_GIT_COMMIT_MESSAGE"
printf "${YELLOW}GIT_BRANCH=%s${NC}\n" "$GIT_BRANCH"
printf "${YELLOW}GIT_TAG=%s${NC}\n" "$GIT_TAG"
printf "${YELLOW}GIT_DIRTY=%s${NC}\n" "$GIT_DIRTY"
printf "${YELLOW}BUILD_TIMESTAMP=%s${NC}\n" "$BUILD_TIMESTAMP"

# Update or insert git commit information in the specified env file
if [ -f "$ENV_FILE" ]; then
    # Create a temporary file
    tmp_env=$(mktemp)
    
    # Remove any existing git commit info lines and write the rest to the temp file
    grep -vE '^(GIT_COMMIT_HASH|GIT_COMMIT_FULL|GIT_COMMIT_DATE|GIT_COMMIT_DATE_SHORT|GIT_COMMIT_MESSAGE|GIT_BRANCH|GIT_TAG|GIT_DIRTY|BUILD_TIMESTAMP)=' "$ENV_FILE" > "$tmp_env" 2>/dev/null || true
    
    # Append the updated git commit info
    {
        echo "GIT_COMMIT_HASH=$GIT_COMMIT_HASH"
        echo "GIT_COMMIT_FULL=$GIT_COMMIT_FULL"
        echo "GIT_COMMIT_DATE=$GIT_COMMIT_DATE"
        echo "GIT_COMMIT_DATE_SHORT=$GIT_COMMIT_DATE_SHORT"
        printf 'GIT_COMMIT_MESSAGE="%s"\n' "$ESCAPED_GIT_COMMIT_MESSAGE"
        echo "GIT_BRANCH=$GIT_BRANCH"
        echo "GIT_TAG=$GIT_TAG"
        echo "GIT_DIRTY=$GIT_DIRTY"
        echo "BUILD_TIMESTAMP=$BUILD_TIMESTAMP"
    } >> "$tmp_env"
    
    # Move the temp file back to the target file
    mv "$tmp_env" "$ENV_FILE"
    
    printf "${GREEN}Git commit information successfully updated in $ENV_FILE${NC}\n"
else
    echo
    printf "${RED}No $ENV_FILE file found. You can copy the above variables to your $ENV_FILE file.${NC}\n"
fi
