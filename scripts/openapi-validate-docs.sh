#!/bin/bash

# Validate OpenAPI documentation syntax and structure
# Usage: ./openapi-validate-docs.sh <service-base-path>
# Example: ./openapi-validate-docs.sh services/authentication-service

set -e

if [ -z "$1" ]; then
    echo "Error: Service base path is required"
    echo "Usage: $0 <service-base-path>"
    echo "Example: $0 services/authentication-service"
    exit 1
fi

SERVICE_PATH="$1"
DIST_PATH="$SERVICE_PATH/dist"

echo "Validating OpenAPI documentation for $SERVICE_PATH..."
docker run --rm -v "$(pwd)":/spec redocly/cli lint "/spec/$DIST_PATH/openapi.yaml"

echo "OpenAPI documentation is valid"