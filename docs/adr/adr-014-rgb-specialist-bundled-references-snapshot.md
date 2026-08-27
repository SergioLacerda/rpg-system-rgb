# ADR-014: RGB Specialist Bundled References Snapshot

## Status

Accepted.

## Context

`skills/specialist/` grounds its answers in `generated/ai-context/core-specialist-pack.json`
and `docs/core/semantic/**` (see ADR-002). Both are repo-relative: if the
`skills/specialist/` folder is copied outside this monorepo — which
ADR-013 now makes possible by publishing it as a downloadable `.zip` — it
loses access to the canonical rule documentation it depends on to answer
questions and cite a source.

The user asked for the package to be self-contained, carrying canonical
documentation material with it, rather than depending on sibling paths in
this monorepo.

## Decision

`skills/specialist/references/rgb-system.md` is added as a bundled copy of
the four canonical rule documents (`docs/core/en/**` and
`docs/core/PT-br/**`), consolidated verbatim with each section labeled by
its original source path.

This is a **hand-maintained snapshot**, not a generated artifact:

- `docs/core/**` remains the single source of truth (ADR-002) — the
  bundled copy never overrides it;
- `SKILL.md`'s grounding-source precedence uses the bundled copy only as a
  fallback (inside the monorepo) or as the sole available source (running
  standalone, outside the monorepo) — citation is still required in every
  case;
- no automated regeneration step exists yet for this file (unlike
  `generated/ai-context/*`, which is regenerated via `make generate`).
  Adding one is a config/tooling change, deferred to a separate task
  (Task 6 of the parallel skill-materialization mission) — it touches
  `config.yaml`, which is out of this mission's documentation-only scope.

## Consequences

Positive:

- the skill package answers questions and cites a source even when copied
  outside this monorepo, satisfying the self-contained requirement without
  changing the existing consumer contract or procedures;
- the precedence order keeps the higher-trust, ID-traceable sources
  (compiled pack, semantic index) as first choice whenever they are
  available, so nothing regresses for in-monorepo use.

Accepted costs:

- `references/rgb-system.md` and `docs/core/**` can drift until the
  automation from Task 6 exists — anyone updating `docs/core/**` should
  know the bundled copy needs a manual refresh until then;
- this is intentionally scoped to content only. It does not decide how the
  folder gets packaged or published; that is ADR-013's decision.

## Relationship To ADR-013

ADR-013 decided how `skills/specialist/` gets zipped and published
(reusing the existing Pages-based publication pipeline). This ADR decides
what the folder needs to contain to still be useful once distributed that
way — a self-contained copy of the canonical rules it grounds answers in.
The two are independent and do not amend each other: ADR-013 governs
packaging/publication, this ADR governs package content.
