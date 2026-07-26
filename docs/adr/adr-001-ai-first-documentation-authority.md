---
mission_id: ai-first-docs-source-model-20260726
date: 2026-07-26
phase: architecture_decision
status: accepted
---

# ADR-001: AI-First Documentation Authority

## Status

Accepted.

## Context

RGB Core V2 needs stable documentation identity before Maker, Specialist,
tooling, Library, PDF, landing, bundles, AI context packs, or search indexes can
consume rules safely.

The project uses Markdown as a human review surface, but AI-first retrieval,
provenance, localized drift checks, projection manifests, and generation gates
need stable source IDs and machine-readable contracts.

## Decision

Adopt a bridge-first documentation authority model:

- English Markdown under `docs/core/en/**` remains the canonical bridge for Core V2
  rule text until a specific semantic source unit is promoted.
- PT-br Markdown under `docs/core/PT-br/**` remains a localized projection and cannot
  introduce mechanics not present in canonical English or canonical semantic
  source units.
- `docs/core/semantic/core-v2.index.json` provides stable IDs, authority type, source
  status, projection paths, relationships, retrieval metadata, and provenance.
- Accepted semantic source files under `docs/core/semantic/source/**` may become
  canonical for specific promoted units when they preserve Markdown EN/PT-br
  projection paths, source IDs, provenance, and validation coverage.
- Generated artifacts, indexes, context packs, Library pages, PDF output,
  landing summaries, and bundles are derived projections only.

## Consequences

- Consumers get stable IDs and provenance without waiting for a repository-wide
  rewrite.
- Tooling can validate source paths, projection paths, relationships,
  translations, consumer contracts, projection ownership, and generated output
  authority.
- Promoted semantic source files can replace Markdown as canonical rule
  authority for their specific source IDs.
- Generated artifacts remain non-canonical and must trace back to source IDs.

## Rejected Alternatives

- Big-bang semantic rewrite.
- Markdown-only long-term source model.
- Generated artifacts as rule authority.
