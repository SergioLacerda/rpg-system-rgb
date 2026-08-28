# ADR-015: Skill Runtime Split — LLM Instruction Layer vs. Go Component Layer

## Status

Accepted.

## Context

`skills/specialist/**` and `skills/maker/**` both ship a complete skill
contract (`SKILL.md`, procedures/schemas, templates, config) but neither has
a working runtime. `internal/components/specialist/specialist.go` and
`internal/components/maker/maker.go` are both 13-line stubs exposing only a
`Descriptor()` — there is no behavior on either side of the boundary today.

The accepted Specialist-first skill roadmap gated Maker's runtime work behind
three conditions: (a) Specialist has a source-trace benchmark, (b) the
bundle/search context shape is stable, (c) Maker's provenance/canon schemas
are accepted. A 2026-08-27 runtime-base review re-checked these against the
current repository:

- (b) and (c) are met — `generated/ai-context/core-specialist-pack.json` /
  `docs/core/semantic/core-v2.index.json` are in stable use, and Maker's four
  schemas (`entity-package`, `maker-report`, `visual-observation`,
  `conflict-report`) exist and are internally consistent.
- (a) is met only **structurally** — `skills/specialist/benchmark/
  golden-qa.yaml` plus `tests/fixtures/specialist_golden_qa_test.go` validate
  that the benchmark dataset cites real semantic IDs, but no Specialist
  runtime exists to exercise the benchmark behaviorally, because Specialist
  has no runtime either.

That review also found that "use Specialist as the base" only holds for the
contract layer (`SKILL.md`/schemas/procedures shape) — there is no proven
runtime pattern in this repo to copy, for either skill.

## Decision

Split any skill's runtime, once it needs one, into two layers:

```text
skills/<name>/**              stable contract — SKILL.md, schemas, templates
                               (unchanged by runtime work)
        │
        ▼  LLM reads instructions, drafts output
LLM layer                     semantic extraction, classification, image
                               reading — no Go equivalent is possible here
        │
        ▼  structured draft, already tagged by layer
internal/components/<name>/** Go runtime — schema validation, structured-fact
                               diffing, deterministic aggregation. Expected
                               to iterate as heuristics are tuned.
```

Applied to Maker specifically, per package kind:

| Kind | LLM layer | Go layer (`internal/components/maker/**`) |
|---|---|---|
| Entity Package | extract facts/inferences from notes | validate against `entity-package.schema.yaml` |
| Maker Report | thin prose rollup | aggregate counts over already-classified packages |
| Visual Observation | read appearance from an image | enforce appearance-only field scope (image never claims a fact field) |
| Conflict Report | flag disagreement during extraction | diff structured fields across packages sharing an entity ID |

The contract layer (`skills/maker/**`) is the stable surface and is not
expected to change as part of this work. The Go component layer
(`internal/components/maker/**`) is where iteration is expected — differ
heuristics, aggregation buckets, and validation strictness will be tuned as
real usage surfaces edge cases.

Gate status: conditions (b) and (c) of the Maker Deferral Criteria are
treated as met; condition (a) is treated as met for starting Maker's Go-side
runtime work (validators, differ, aggregator — none of which require a
working Specialist runtime to build or test), but not as evidence that
Maker's own LLM-side behavior has been benchmarked.

## Consequences

Positive:

- Isolates expected instability (differ/aggregator/validator heuristics) to
  one package boundary, so iterating on it does not put the already-accepted
  `skills/maker/**` contract at risk of drifting out of sync with runtime
  behavior.
- Establishes a pattern reusable by Specialist (and any future skill) once it
  grows a real runtime, rather than a Maker-specific one-off.
- Continues the split already used for Specialist's existing tooling (bundle
  generation and the golden-qa structural check are Go; explaining a rule is
  LLM) instead of introducing an unrelated third pattern.

Negative / accepted costs:

- No implementation of this split exists yet for either skill — this ADR
  sets policy ahead of a first concrete build.
- Building `internal/components/maker/**`'s validators, differ, and
  aggregator, and wiring Maker's LLM-side procedures into an invocable
  runtime, remain open implementation work — out of scope for this ADR and
  for the Strategist mission that produced it (implementation_handoff, not
  documentation).

## Non-Goals

This ADR does not authorize or perform any Go source change under
`internal/components/maker/**`, nor any change to `skills/maker/**`'s
existing contract. It records the architectural decision only.
