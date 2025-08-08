#!/bin/bash
set -e

if ! command -v oapi-codegen &> /dev/null; then
  echo "oapi-codegen not found, installing..."
  go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
fi