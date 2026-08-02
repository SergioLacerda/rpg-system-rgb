# Components

```mermaid
C4Component
    title RGB System V2 — internal/components Boundary

    Container_Boundary(app, "internal/app") {
        Component(app_c, "internal/app", "Go package", "Thin orchestration layer; the only import path for cmd/*")
    }

    Container_Boundary(components, "internal/components") {
        Component(contract, "internal/components (root)", "Go package", "Shared contract types; stdlib-only leaf")
        Component(core, "core", "Go package", "RGB rule engine: characters, actions, combat, states")
        Component(tooling_c, "tooling", "Go package", "Semantic-docs validators + projection generator")
        Component(publication_c, "publication", "Go package", "HTML Library + PDF publication")
        Component(maker_c, "maker", "Go package", "Campaign/content structuring (stub: boundary only, no behavior)")
        Component(specialist_c, "specialist", "Go package", "Rule consultation and validation (stub: boundary only, no behavior)")
        Component(bundles_c, "bundles", "Go package", "Semantic index projection to generated/bundle/rgb.bundle.json (implemented, minimal — see ADR-007)")
    }

    Rel(app_c, core, "imports")
    Rel(app_c, tooling_c, "imports")
    Rel(app_c, publication_c, "imports")
    Rel(core, contract, "imports")
    Rel(tooling_c, contract, "imports")
    Rel(publication_c, contract, "imports")
```

The retired `compiler` and `library` packages are no longer part of this
boundary. Active public publication now runs through the `publication`
component and Astro packaging.

This diagram matches, not just illustrates, what
`tests/architecture/architecture_test.go` mechanically enforces
(`TestComponentsDoNotImportOuterLayers`, `TestComponentsDoNotImportSiblings`,
`TestSharedContractStaysLeaf`, `TestEntrypointsGoThroughApp`) and what
`.golangci.yml`'s `depguard` rules repeat:

- Every arrow shown above is allowed. Any arrow **not** shown — most
  importantly, any edge directly between two sibling components under
  `internal/components/` (e.g. `compiler` importing `tooling`) — is
  forbidden and fails both the linter and `make test-arch`.
- `internal/components` (the root/contract package) must stay a leaf:
  it may depend on the standard library only, never on `app`, `cmd`, or
  any sibling component.
- `cmd/*` binaries may only import `internal/app` — never
  `internal/components/*` directly.

These import boundaries are non-negotiable: they predate and are
independent of any CLI topology or bundle-scope decision (ADR-006,
ADR-007). No future integration work may shortcut them with a
sibling-to-sibling import, even temporarily.

← [Back to Architecture Index](README.md)
