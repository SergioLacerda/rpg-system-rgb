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
        Component(compiler_c, "compiler", "Go package", "Markdown parser + HTML/print renderer")
        Component(library_c, "library", "Go package", "Site assembly + PDF export instructions")
        Component(maker_c, "maker", "Go package", "Campaign/content structuring (stub)")
        Component(specialist_c, "specialist", "Go package", "Rule consultation and validation (stub)")
        Component(bundles_c, "bundles", "Go package", "Consolidated bundle format (stub)")
    }

    Rel(app_c, core, "imports")
    Rel(app_c, tooling_c, "imports")
    Rel(app_c, compiler_c, "imports")
    Rel(app_c, library_c, "imports")
    Rel(core, contract, "imports")
    Rel(tooling_c, contract, "imports")
    Rel(compiler_c, contract, "imports")
    Rel(library_c, contract, "imports")
```

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

← [Back to Architecture Index](README.md)
