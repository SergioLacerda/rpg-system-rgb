# Advanced Test Coverage Strategy

## Purpose

Numeric statement coverage floors (`make cover-check`) answer "how much code
ran," not "would a plausible mistake be caught." This document describes the
three advanced layers built on top of that baseline, per the
`20260827-advanced-test-coverage-strategy` planning mission: mutation,
golden/oracle, and Gherkin/BDD traceability.

## Coverage Model

```text
┌───────────────────────────┐
│ Numeric coverage gates     │  package and web thresholds
├───────────────────────────┤
│ Mutation gates             │  named risky implementation changes must be killed
├───────────────────────────┤
│ Golden oracle gates        │  canonical inputs produce accepted outputs
├───────────────────────────┤
│ Gherkin/BDD traceability   │  scenarios stay mirrored by executable tests
└───────────────────────────┘
```

## Mutation Layer

`scripts/ci/mutation-core.sh` copies the worktree to `/tmp`, applies one
hand-written mutation at a time, and requires a named kill-suite of Go
packages to fail. `run_mutation` takes an optional 5th argument overriding
the default kill-suite (`./internal/components/core ./tests/properties`) so
each mutant can target the packages most likely to catch it.

Current risk-tier matrix:

| Group | Mutants | Kill suite |
| --- | --- | --- |
| Core resolution | `strong-success-boundary`, `success-with-cost-dropped`, `resolve-modifier-sign` | core, properties, core_behavior |
| Damage | `damage-penetration-direction`, `damage-shield-pre-armor`, `damage-injure-on-zero` | core, core_behavior, properties |
| Resources | `shield-derivation`, `health-derivation`, `resources-isdown-boundary` | core, fixtures, core_behavior, properties |
| Encounter flow | `objective-failure-round`, `action-declared-order`, `surprise-priority-inverted` | core, core_behavior, simulation, properties |

`action-declared-order` reverses `RunEncounter`'s action-processing loop;
it only became killable after adding
`TestRunEncounterResolvesActionsInDeclaredOrder`
(`internal/components/core/encounter_order_test.go`), which locks the
`RunEncounter` doc comment's own "in list order" contract — no prior test
depended on `ActionResults` lining up with `Actions` by position.
`surprise-priority-inverted` targets `InitiativeOrder`'s surprise-goes-first
tie-break in `internal/components/core/initiative.go`.

The tooling/publication group from the original risk matrix is deferred
per design.md's "only include if deterministic" guidance.

Add a mutant only when it represents a plausible business-rule defect, and
confirm it is killed (`make mutation-core`) before keeping it. A mutation
framework is deliberately not adopted while the hand-written matrix stays
small and fast (per `analysis.md` uncertainty U-01).

## Golden Layer

Golden tests assert stable structured or normalized values, not full-prose
snapshots:

- **Specialist Q/A** (`skills/specialist/benchmark/golden-qa.yaml`): every
  normative entry (outside `category: ambiguous`) cites real
  `semantic_ids` and declares `expected_answer_traits` — short
  keyword/phrase fragments a correct answer must contain, ahead of any
  runtime Specialist implementation to check them against
  (`tests/fixtures/specialist_golden_qa_test.go`).
- **Semantic JSON fixtures**
  (`tests/fixtures/semantic_projection_golden_test.go`): a curated set of
  canonical semantic IDs keep stable `authority_type`/`source_status` and
  stay covered by at least one entry in
  `docs/core/semantic/projection-manifest.v0.1.json`.
- **CLI output** (`cmd/rgb/golden_output_test.go`): exact-text assertions
  for the unified CLI's intentionally-public usage/error strings.

Reserve exact snapshots for output that is intentionally canonical; prefer
structured assertions everywhere else (per `analysis.md` uncertainty U-02).

## Gherkin/BDD Layer

`.feature` files under `tests/features/<domain>/` stay the readable
specification; `// mirrors: <feature path>#<scenario name>` anchor comments
in Go tests under `tests/` are the executable enforcement:

- `TestEveryScenarioHasAMirror` (`tests/features/sync_test.go`) fails if a
  declared scenario has no matching anchor.
- `TestNoMirrorAnchorReferencesAMissingScenario` is the reverse-trace
  guard: it fails if an anchor still points at a scenario that was renamed
  or deleted.
- Scenario taxonomy is the existing domain directory layout (`combat`,
  `core`, `damage`, `encounters`, `equipment`, `fixtures`) — no separate
  taxonomy file is maintained.

Adopting a real Gherkin runner (e.g. `godog`) stays declined; the mirror
model is kept until step-reuse or drift becomes a real maintenance cost
(per `tasks.md` side quest SQ-003).

## Validation Commands

```bash
GOCACHE=/tmp/go-cache go test ./...
GOCACHE=/tmp/go-cache go test ./internal/components/core ./tests/core_behavior ./tests/properties ./tests/simulation ./tests/features ./tests/fixtures ./cmd/rgb
make mutation-core
make cover-check
make coverage-inventory   # read-only inventory: mutation/golden/BDD/mirror counts
make check-fast
make check
```

`make coverage-inventory` is informational only (like
`make go-file-size-report`) and is not wired into `check-fast` or `check`;
it exists so coverage movement across the three advanced layers is
reviewable without re-deriving the counts by hand.
