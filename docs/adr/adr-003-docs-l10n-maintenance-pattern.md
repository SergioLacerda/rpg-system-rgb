# ADR-003: Documentation L10n Maintenance Pattern

## Status

Accepted.

## Context

RGB System keeps English documentation as the canonical technical source and
Portuguese documentation as a localized reader-facing surface. Path parity is
useful, but it does not show whether a localized document is current, stale, or
missing after an English source changes.

Core V2 documentation updates exposed this maintenance cost: English rule
documents changed first, and PT-br alignment required a separate review pass.

## Decision

Adopt a lightweight advisory L10n manifest for documentation maintenance.

The project path model is:

```text
Canonical source      docs/core/en/**
Localized surface     docs/core/PT-br/**
Traceability layer    docs/core/semantic/l10n-manifest.v0.1.json
```

The manifest records:

- canonical English source path;
- localized path;
- target locale;
- authority type;
- translation status;
- source revision marker;
- optional notes.

English remains the canonical technical source. PT-br documents are localized
reader-facing projections and must not introduce mechanics that are absent from
the English source.

The first manifest is advisory. It supports review and tooling, but it does not
generate localized documentation and does not block publication by itself.

## Authority Types

- `concept`
- `rule`
- `procedure`
- `example`
- `design_note`
- `translation`
- `generated_artifact`

## Translation Statuses

- `current`
- `needs_review`
- `stale`
- `missing`
- `intentionally_divergent`
- `not_applicable`

## Review Workflow

1. Update the canonical English document.
2. Identify affected localized documents through the manifest.
3. Mark affected PT-br entries as `needs_review` or `stale`.
4. Update PT-br content in a separate documentation task.
5. Return the manifest entry to `current` only after review.

## Consequences

- Maintainers can see source-to-localized traceability without reading every
  pair manually.
- Future Library, PDF, bundle, search, and Specialist consumers can disclose
  localization status before relying on localized content.
- Full semantic-source migration, per-document frontmatter, and automatic
  translation remain separate future decisions.
