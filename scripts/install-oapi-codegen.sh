#!/bin/bash
set -e

if ! command -v oapi-codegen &> /dev/null; then
  echo "oapi-codegen not found, installing..."
  go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
fi