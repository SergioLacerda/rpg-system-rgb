# ADR-004: Modular Monolith In A Single Monorepo

## Status

Proposed.

## Context

RGB System V2 spans several planned products: Core rules, Maker, Specialist,
Tooling, Library (HTML/PDF), and a landing page. The V2 initial analysis
(base_project intake, section 20) lists two viable repository strategies: a
monorepo or separate repositories per product.

The current state favors consolidation:

- the products share one canonical source (`docs/` Markdown + YAML) and one
  contract vocabulary that is still in formation;
- the Go scaffold already models products as component boundaries under
  `internal/components/`;
- there is a single maintainer and no divergent release cycles;
- validation must run across products (docs ↔ tooling ↔ skills), which is
  cheap in one repository and expensive across several.

## Decision

Adopt a **modular monolith inside a single monorepo**.

- One Go module (`github.com/SergioLacerda/rpg-system-rgb`) with one binary
  entrypoint per product surface under `cmd/` as they materialize
  (`rgb`, later `rgb-compile`, `rgb-validate`, `rgb-query`).
- Each product is a component boundary under `internal/components/<name>`,
  isolated by enforced import rules (see below), so a future repository split
  is a directory move, not a rewrite.
- Non-Go surfaces stay as monorepo siblings: `docs/` (canonical content),
  `skills/` (lightweight Markdown/YAML skills), `web/` (Astro landing),
  `generated/` (derived artifacts only), `scripts/` (standalone validation
  module).

Boundaries are enforced mechanically, not by convention:

- `tests/architecture/` fails the build when a component imports `app`,
  `cmd`, or a sibling component;
- `.golangci.yml` depguard rules repeat the same constraints in the linter;
- `make review-structure` is the aggregate local gate
  (`docs/engineering/base-structure-review-workflow.md`).

## Split Criteria

Extract a product into its own repository only when at least one of these
holds (mirrors the V2 analysis, section 20.3):

- release cycles genuinely diverge;
- distinct teams or communities own different products;
- independent distribution produces real benefit (e.g. a public SDK);
- repository size becomes an operational obstacle.

Because component boundaries are import-clean by construction, a split
requires moving `internal/components/<name>` plus its `cmd/` surface and
publishing the shared contract package — nothing else should break.

## Consequences

Positive:

- contracts between products stay synchronized in one place;
- cross-product validation (docs, schemas, bundles) runs in one gate;
- refactoring across boundaries is atomic;
- the architecture tests keep the monolith modular instead of letting it
  degrade into an entangled monolith.

Negative / accepted costs:

- one repository carries heterogeneous toolchains (Go, Node/Astro, Markdown
  pipelines);
- releases of independent products must be tagged within one history;
- contributors see the whole ecosystem even when touching one product.
