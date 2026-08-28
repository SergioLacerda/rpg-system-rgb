#!/usr/bin/env bash
# Read-only advanced coverage inventory report: prints mutation case,
# golden entry, feature/scenario, and mirror-anchor counts. Informational
# only — never modifies tracked files and never fails the build.
set -euo pipefail

repo_root="$(pwd)"
cd "${repo_root}"

echo "== Mutation coverage (scripts/ci/mutation-core.sh) =="
grep -c '^run_mutation "' scripts/ci/mutation-core.sh

echo
echo "== Golden coverage (skills/specialist/benchmark/golden-qa.yaml) =="
grep -c '^  - id:' skills/specialist/benchmark/golden-qa.yaml

echo
echo "== Gherkin/BDD coverage (tests/features, by domain) =="
total_files=0
total_scenarios=0
for domain_dir in tests/features/*/; do
  domain="$(basename "${domain_dir}")"
  files=$(find "${domain_dir}" -maxdepth 1 -name '*.feature' | wc -l | tr -d ' ')
  scenarios=$(grep -rh '^\s*Scenario\(\s*Outline\)\?:' "${domain_dir}" 2>/dev/null | wc -l | tr -d ' ')
  printf '%-12s %3s feature files, %3s scenarios\n' "${domain}" "${files}" "${scenarios}"
  total_files=$((total_files + files))
  total_scenarios=$((total_scenarios + scenarios))
done
printf '%-12s %3s feature files, %3s scenarios\n' "total" "${total_files}" "${total_scenarios}"

echo
echo "== Mirror anchors (tests/, excluding tests/features) =="
find tests -type f -name '*.go' -not -path 'tests/features/*' -print0 \
  | xargs -0 grep -h '// mirrors:' 2>/dev/null \
  | wc -l | tr -d ' '
