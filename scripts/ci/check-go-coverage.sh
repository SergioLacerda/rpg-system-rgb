#!/usr/bin/env bash
set -euo pipefail

go_bin="${GO:-go}"
export GOCACHE="${GOCACHE:-/tmp/go-cache}"

check_package() {
  local label="$1"
  local package="$2"
  local floor="$3"
  local output
  local pct

  output="$("${go_bin}" test -cover "${package}")"
  printf '%s\n' "${output}"
  pct="$(printf '%s\n' "${output}" | sed -nE 's/.*coverage: ([0-9]+\.[0-9]+)% of statements.*/\1/p' | tail -1)"
  if [[ -z "${pct}" ]]; then
    printf 'coverage unavailable for %s (%s)\n' "${label}" "${package}" >&2
    return 1
  fi
  printf '%s coverage: %s%% (floor: %s%%)\n' "${label}" "${pct}" "${floor}"
  awk -v pct="${pct}" -v floor="${floor}" 'BEGIN { exit !(pct + 0 >= floor + 0) }' \
    || {
      printf '%s coverage %s%% is below the %s%% floor\n' "${label}" "${pct}" "${floor}" >&2
      return 1
    }
}

check_package "component registry" "./internal/components" "${COMPONENTS_COVER_THRESHOLD:-90}"
check_package "bundle component" "./internal/components/bundles" "${BUNDLES_COVER_THRESHOLD:-90}"
check_package "core formula engine" "./internal/components/core" "${CORE_COVER_THRESHOLD:-90}"
check_package "core fixtures" "./internal/components/core/fixtures" "${FIXTURES_COVER_THRESHOLD:-90}"
check_package "maker component descriptor" "./internal/components/maker" "${MAKER_COVER_THRESHOLD:-90}"
check_package "specialist component descriptor" "./internal/components/specialist" "${SPECIALIST_COVER_THRESHOLD:-90}"
check_package "semantic tooling" "./internal/components/tooling" "${TOOLING_COVER_THRESHOLD:-90}"
check_package "publication component" "./internal/components/publication" "${PUBLICATION_COVER_THRESHOLD:-90}"
check_package "application orchestration" "./internal/app" "${APP_COVER_THRESHOLD:-90}"
check_package "unified CLI" "./cmd/rgb" "${CLI_COVER_THRESHOLD:-90}"
check_package "legacy tooling CLI" "./cmd/rgb-tooling" "${TOOLING_CLI_COVER_THRESHOLD:-90}"
