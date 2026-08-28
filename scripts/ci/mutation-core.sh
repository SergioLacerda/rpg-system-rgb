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
  local kill_suite_text="${5:-./internal/components/core ./tests/properties}"
  local -a kill_suite
  read -r -a kill_suite <<< "${kill_suite_text}"
  local worktree="${tmp_root}/${id}"

  copy_worktree "${worktree}"
  perl -0pi -e "s/${from}/${to}/" "${worktree}/${file}"

  if (cd "${worktree}" && "${go_bin}" test "${kill_suite[@]}" >/tmp/rgb-mutation-"${id}".log 2>&1); then
    printf 'mutation survived: %s (%s)\n' "${id}" "${file}" >&2
    sed -n '1,120p' /tmp/rgb-mutation-"${id}".log >&2
    return 1
  fi

  printf 'mutation killed: %s\n' "${id}"
}

cd "${repo_root}"

# --- Core resolution ---
run_mutation "damage-penetration-direction" "internal/components/core/damage.go" \
  "target\\.Resources\\.Armor-input\\.Penetration" \
  "target.Resources.Armor+input.Penetration"
run_mutation "shield-derivation" "internal/components/core/resources.go" \
  "vectors\\.B \\* 3" \
  "vectors.B * 2"
run_mutation "strong-success-boundary" "internal/components/core/resolution.go" \
  "margin >= 3" \
  "margin > 3"
run_mutation "success-with-cost-dropped" "internal/components/core/resolution.go" \
  "outcome == OutcomeStrongSuccess \\|\\| outcome == OutcomeSuccess \\|\\| outcome == OutcomeSuccessWithCost" \
  "outcome == OutcomeStrongSuccess || outcome == OutcomeSuccess" \
  "./internal/components/core ./tests/properties ./tests/core_behavior"
run_mutation "resolve-modifier-sign" "internal/components/core/resolution.go" \
  "actingValue \\+ modifier - opposingValue" \
  "actingValue - modifier - opposingValue" \
  "./internal/components/core ./tests/properties ./tests/core_behavior"

# --- Damage ---
run_mutation "damage-shield-pre-armor" "internal/components/core/damage.go" \
  "shieldAbsorbed := min\\(target\\.Resources\\.CurrentShield, afterArmor\\)" \
  "shieldAbsorbed := min(target.Resources.CurrentShield, input.Impact)" \
  "./internal/components/core ./tests/core_behavior ./tests/properties"
run_mutation "damage-injure-on-zero" "internal/components/core/damage.go" \
  "else if healthDamage > 0" \
  "else if healthDamage >= 0" \
  "./internal/components/core ./tests/core_behavior ./tests/properties"

# --- Resources ---
run_mutation "health-derivation" "internal/components/core/resources.go" \
  "health := 4 \\+ vectors\\.R \\+ vectors\\.B" \
  "health := 4 + vectors.R - vectors.B" \
  "./internal/components/core ./tests/fixtures ./tests/core_behavior ./tests/properties"
run_mutation "resources-isdown-boundary" "internal/components/core/resources.go" \
  "resources\\.CurrentHealth <= 0" \
  "resources.CurrentHealth < 0" \
  "./internal/components/core ./tests/fixtures"

# --- Encounter flow ---
run_mutation "objective-failure-round" "internal/components/core/encounter.go" \
  "outcome\\.ResolvedRound = max\\(currentRound, objective\\.DeadlineRounds\\)" \
  "outcome.ResolvedRound = min(currentRound, objective.DeadlineRounds)" \
  "./internal/components/core ./tests/core_behavior ./tests/simulation"
run_mutation "action-declared-order" "internal/components/core/encounter.go" \
  "for _, action := range encounter\\.Actions \\{" \
  "for i := len(encounter.Actions) - 1; i >= 0; i-- {\n\t\taction := encounter.Actions[i]" \
  "./internal/components/core ./tests/core_behavior ./tests/simulation ./tests/properties"
run_mutation "surprise-priority-inverted" "internal/components/core/initiative.go" \
  "return sorted\\[i\\]\\.Surprise" \
  "return sorted[j].Surprise" \
  "./internal/components/core ./tests/core_behavior"
