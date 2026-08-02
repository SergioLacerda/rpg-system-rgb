#!/usr/bin/env bash
set -euo pipefail

generate_output="$(make generate)"
echo "$generate_output"

if ! git diff --exit-code -- generated docs/core/semantic web/landing/generated; then
  echo "::error::generated artifacts drifted; run 'make generate' and review the diff"
  exit 1
fi

# git diff is blind to gitignored-but-untracked files, so a declared
# projection output that was never `git add -f`'d (like
# web/landing/generated/core-v2-summary.json once was) passes silently
# here even though a fresh checkout is missing it. Cross-check every
# path `make generate` reports against the git index.
missing=0
while IFS= read -r path; do
  if ! git ls-files --error-unmatch -- "$path" >/dev/null 2>&1; then
    echo "::error::declared projection output '$path' is not tracked in git (run: git add -f $path)"
    missing=1
  fi
done < <(sed -n 's/^generated //p' <<< "$generate_output")

if [ "$missing" -ne 0 ]; then
  exit 1
fi
