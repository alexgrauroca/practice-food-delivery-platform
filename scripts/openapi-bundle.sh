#!/bin/bash

# Generate the OpenAPI distribution file from modular files and HTML documentation
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

# Ensure dist directory exists
mkdir -p "$SERVICE_PATH/docs/dist"

echo "Generating OpenAPI distribution file for $SERVICE_PATH..."
docker run --rm \
  -v "$(pwd)":/spec \
  redocly/cli bundle \
  "/spec/$SERVICE_PATH/docs/openapi.yaml" \
  -o "/spec/$SERVICE_PATH/docs/dist/openapi.yaml"

echo "Generating HTML documentation..."
docker run --rm \
  -v "$(pwd)":/spec \
  redocly/cli build-docs \
  "/spec/$SERVICE_PATH/docs/dist/openapi.yaml" \
  -o "/spec/$SERVICE_PATH/docs/dist/index.html"

echo "Files successfully generated at $SERVICE_PATH/docs/dist/:"
echo "- openapi.yaml (bundled OpenAPI specification)"
echo "- index.html (API documentation)"