#!/bin/bash
set -e

if ! command -v mockgen &> /dev/null; then
  echo "mockgen not found, installing..."
  go install github.com/golang/mock/mockgen@latest
fi