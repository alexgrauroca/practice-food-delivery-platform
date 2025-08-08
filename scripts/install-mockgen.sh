#!/bin/bash
set -e

if ! command -v mockgen &> /dev/null; then
  echo "mockgen not found, installing..."
  go install go.uber.org/mock/mockgen@latest
fi