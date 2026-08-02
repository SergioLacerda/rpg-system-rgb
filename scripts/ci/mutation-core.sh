#!/usr/bin/env bash
set -euo pipefail

go_bin="${GO:-go}"
export GOCACHE="${GOCACHE:-/tmp/go-cache}"

repo_root="$(pwd)"
tmp_root="$(mktemp -d "${TMPDIR:-/tmp}/rgb-core-mutation-XXXXXX")"
trap 'rm -rf "${tmp_root}"' EXIT

copy_worktree() {
  local target="$1"
  mkdir -p "${target}"
  find . \
    -path './.git' -prune -o \
    -path './web/landing/node_modules' -prune -o \
    -path './web/landing/dist' -prune -o \
    -path './web/landing/coverage' -prune -o \
    -type f -print | tar -cf - -T - | tar -xf - -C "${target}"
}

run_mutation() {
  local id="$1"
  local file="$2"
  local from="$3"
  local to="$4"
  local worktree="${tmp_root}/${id}"

  copy_worktree "${worktree}"
  perl -0pi -e "s/${from}/${to}/" "${worktree}/${file}"

  if (cd "${worktree}" && "${go_bin}" test ./internal/components/core ./tests/properties >/tmp/rgb-mutation-"${id}".log 2>&1); then
    printf 'mutation survived: %s (%s)\n' "${id}" "${file}" >&2
    sed -n '1,120p' /tmp/rgb-mutation-"${id}".log >&2
    return 1
  fi

  printf 'mutation killed: %s\n' "${id}"
}

cd "${repo_root}"
run_mutation "damage-penetration-direction" "internal/components/core/damage.go" \
  "target\\.Resources\\.Armor-input\\.Penetration" \
  "target.Resources.Armor+input.Penetration"
run_mutation "shield-derivation" "internal/components/core/resources.go" \
  "vectors\\.B \\* 3" \
  "vectors.B * 2"
run_mutation "strong-success-boundary" "internal/components/core/resolution.go" \
  "margin >= 3" \
  "margin > 3"
