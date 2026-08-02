#!/usr/bin/env bash
set -euo pipefail

make generate

if ! git diff --exit-code -- generated docs/core/semantic; then
  echo "::error::generated artifacts drifted; run 'make generate' and review the diff"
  exit 1
fi
