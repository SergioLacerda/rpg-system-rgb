#!/usr/bin/env bash
set -euo pipefail

limit=${1:-200}
if [[ ! $limit =~ ^[0-9]+$ ]]; then
  echo "go-file-size-report: limit must be a non-negative integer" >&2
  exit 1
fi

results=$(find cmd internal -type f -name '*.go' ! -name '*_test.go' \
  -exec wc -l {} + |
  awk -v limit="$limit" '
    $2 != "total" && $1 > limit {
      count = $1
      sub(/^[[:space:]]*[0-9]+[[:space:]]+/, "")
      printf "%s\t%s\n", count, $0
    }
  ' | LC_ALL=C sort -t $'\t' -k1,1nr -k2,2 | awk -F '\t' '{ print $2 " " $1 }')

if [[ -n $results ]]; then
  printf '%s\n' "$results"
else
  echo none
fi
