#!/usr/bin/env bash
set -euo pipefail

required_files=(
  "CHANGELOG.md"
  "CODEOWNERS"
  "CODE_OF_CONDUCT.md"
  "CONTRIBUTING.md"
  "LICENSE"
  "LICENSES.md"
  "SECURITY.md"
  ".github/dependabot.yml"
  ".github/pull_request_template.md"
  ".github/ISSUE_TEMPLATE/bug_report.yml"
  ".github/ISSUE_TEMPLATE/docs_update.yml"
)

missing=()
for path in "${required_files[@]}"; do
  if [ ! -s "${path}" ]; then
    missing+=("${path}")
  fi
done

if [ "${#missing[@]}" -gt 0 ]; then
  printf '::error::missing required governance file: %s\n' "${missing[@]}"
  exit 1
fi
