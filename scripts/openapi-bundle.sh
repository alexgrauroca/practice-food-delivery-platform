#!/bin/bash

# Generate the OpenAPI distribution file from modular files
# Usage: ./openapi-bundle.sh <service-base-path>
# Example: ./openapi-bundle.sh services/authentication-service

set -e

if [ -z "$1" ]; then
    echo "Error: Service base path is required"
    echo "Usage: $0 <service-base-path>"
    echo "Example: $0 services/authentication-service"
    exit 1
fi

SERVICE_PATH="$1"
DOCS_PATH="$SERVICE_PATH/docs"
DIST_PATH="$SERVICE_PATH/docs/dist"

# Ensure dist directory exists
mkdir -p "$DIST_PATH"

echo "Generating OpenAPI distribution file for $SERVICE_PATH..."
docker run --rm -v "$(pwd)":/spec redocly/cli bundle \
  "/spec/$DOCS_PATH/openapi.yaml" \
  -o "/spec/$DIST_PATH/openapi.yaml"

echo "Distribution file successfully generated at $DIST_PATH/openapi.yaml"
