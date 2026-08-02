# ADR-007: Bundle Expansion Requires A Consumer Contract First

## Status

Accepted.

## Context

`internal/components/bundles` currently projects a minimal, real subset of
the semantic index (`id`, `kind`, `title`, `source_status`,
`relationships`) into `generated/bundle/rgb.bundle.json`. A base project
architecture review confirms this is a useful minimum but flags that the
plan's full bundle shape (`rules`, `procedures`, `graph`, `search`, locale
packs, validation reports) has no defined consumer yet:

- `internal/components/maker` and `internal/components/specialist` are
  boundary stubs only (`Descriptor()` and nothing else) — there is no
  running code to derive a real field-level contract from;
- MkDocs + Astro (the actual Library/PDF publication path) consumes
  canonical `docs/` Markdown and the semantic index directly today, not the
  bundle — `internal/components/library` was retired by
  [ADR-008](adr-008-retire-go-publication-pipeline.md) for being an
  uncalled parallel path, not a bundle consumer to design a contract for;
- no external/third-party consumer exists yet.

Both the review's proposal and design notes are explicit: expanding the
bundle before consumers are defined risks producing "attractive but unused
generated artifacts," and non-goals explicitly forbid treating the
Maker/Specialist stubs as if they were implemented products.

## Decision

Bundle expansion is **contract-first**: no field is added to `Bundle` /
`BundleUnit` (`internal/components/bundles/bundles.go`) for a given
consumer until that consumer has a written contract stating:

1. which bundle fields it reads and why (not "might use");
2. the read pattern (one-shot load, incremental query, streaming);
3. a compatibility expectation (can it tolerate additive-only changes, or
   does it need a schema version bump on any change);
4. a test that exercises the consumer against a bundle fixture, per
   `tasks.md`'s "add tests for any new CLI or bundle contract before
   implementation."

This repo does not yet have a contract meeting that bar for any consumer.
Concretely, that means:

- Maker and Specialist stay boundary stubs until their first real use case
  exists; writing a bundle contract for a stub product would be
  speculation, not a contract.
- Library's future bundle consumption (if any) needs its own contract
  before `bundles` grows fields for it — Library reading `docs/` directly
  today is not itself a reason to add unrelated fields.
- External/third-party consumers need the same bar; `bundleSchema =
  "rgb-system-bundle/0.1"` stays additive-compatible (new optional fields
  only) until a real external consumer forces a versioning decision.

## Consequences

Positive:

- prevents the bundle from accumulating unused speculative structure;
- keeps `bundles` importing only `internal/components` (the shared leaf
  contract), per the existing `TestComponentsDoNotImportSiblings`
  boundary — a contract-first bundle has no reason to reach into a sibling
  component to guess at its needs.

Negative / accepted costs:

- Maker/Specialist integration stays blocked until product work defines
  what they actually need, which this ADR does not do — it is a process
  gate, not a schedule.

## Non-negotiable: component import boundaries

Independent of bundle scope, `docs/architecture/components.md`'s enforced
rule set (`tests/architecture/architecture_test.go`:
`TestComponentsDoNotImportOuterLayers`, `TestComponentsDoNotImportSiblings`,
`TestSharedContractStaysLeaf`, `TestEntrypointsGoThroughApp`) stays exactly
as-is. No bundle contract, CLI change, or integration work may introduce a
sibling-to-sibling import (e.g. `maker` importing `bundles` directly) to
shortcut this policy — it must go through `internal/app`, unchanged from
ADR-004's original boundary decision.
