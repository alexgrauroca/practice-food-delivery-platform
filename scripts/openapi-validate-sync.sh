#!/bin/bash

# Validate that the distribution file is in sync with source files by comparing hashes
# Usage: ./openapi-validate-sync.sh <service-base-path>
# Example: ./openapi-validate-sync.sh services/authentication-service

set -e

if [ -z "$1" ]; then
    echo "Error: Service base path is required"
    echo "Usage: $0 <service-base-path>"
    echo "Example: $0 services/authentication-service"
    exit 1
fi

# Verify sha256sum is available
if ! command -v sha256sum &> /dev/null; then
    echo "Error: sha256sum command not found"
    echo "Please install coreutils or equivalent package"
    exit 1
fi

SERVICE_PATH="$1"
DOCS_PATH="$SERVICE_PATH/docs"
DIST_PATH="$SERVICE_PATH/docs/dist"
DIST_FILE="$DIST_PATH/openapi.yaml"

if [ ! -f "$DIST_FILE" ]; then
    echo "Error: Distribution file not found at $DIST_FILE"
    echo "Please run 'scripts/openapi-bundle.sh $SERVICE_PATH' first"
    exit 1
fi

# Get hash of current distribution file
DIST_HASH=$(sha256sum "$DIST_FILE" | cut -d' ' -f1)

# Calculate expected hash by bundling source files and calculating hash directly from stdout
EXPECTED_HASH=$(docker run --rm -v "$(pwd)":/spec redocly/cli bundle \
    "/spec/$DOCS_PATH/openapi.yaml" --format=yaml 2>/dev/null | sha256sum | cut -d' ' -f1)

if [ "$DIST_HASH" != "$EXPECTED_HASH" ]; then
    echo "Error: OpenAPI distribution file is out of sync with source files"
    echo "Current distribution file hash: $DIST_HASH"
    echo "Expected hash from source files: $EXPECTED_HASH"
    echo "Please run 'scripts/openapi-bundle.sh $SERVICE_PATH' to update the distribution file"
    exit 1
fi

echo "Distribution file is properly synchronized with source files"