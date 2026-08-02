# Containers

```mermaid
C4Container
    title RGB System V2 — Containers (What's Actually Built)

    Person(gm, "Game Master / Player")

    System_Boundary(rgb, "RGB System V2") {
        Container(docs, "docs/core/**", "Markdown + YAML", "Canonical bilingual rule source (EN bridge, PT-br projection)")
        Container(semantic, "docs/core/semantic/**", "JSON", "Stable IDs, l10n manifest, projection manifest, consumer contracts")
        Container(rgb_cli, "cmd/rgb", "Go binary", "Unified CLI: status, validate, generate, bundle (see ADR-006)")
        Container(tooling, "cmd/rgb-tooling", "Go binary", "Deprecated compatibility path for validate/generate/bundle (ADR-006)")
        Container(mkdocs, "MkDocs + Astro", "Python + Node", "HTML Library + PDF publication path")
        Container(generated, "generated/**", "JSON", "Derived: bundles, ai-context, landing summary, semantic projections, search")
    }

    Rel(gm, docs, "Reads")
    Rel(docs, semantic, "Indexed by")
    Rel(rgb_cli, docs, "Validates")
    Rel(rgb_cli, semantic, "Reads and writes")
    Rel(rgb_cli, generated, "Generates projections and bundle into")
    Rel(tooling, docs, "Validates (deprecated path)")
    Rel(tooling, semantic, "Reads and writes (deprecated path)")
    Rel(mkdocs, docs, "Renders")
    Rel(gm, mkdocs, "Reads published Library / downloads PDF")
```

`cmd/rgb` does not yet consume `docs/core/**` at runtime for rule
resolution — RGB Core V2's rule logic lives in `internal/components/core`
and is validated against canonical docs via Gherkin/BDD tests
(`tests/features/`, `tests/core_behavior/`), not by `cmd/rgb` reading
Markdown directly. Its `validate`/`generate`/`bundle` subcommands do read
and write `docs/**` and `generated/**`, as shown above; only rule-engine
consumption of docs at runtime is still absent.

`cmd/rgb-tooling` is deprecated, not removed — see
[ADR-006](../adr/adr-006-unified-cli-topology.md) for the compatibility
decision and M006 removal requirement. The former Go-rendered
HTML/print/PDF pipeline (`cmd/rgb-compiler`, `internal/components/compiler`,
`internal/components/library`) is no longer the public publication path.
MkDocs + Astro owns HTML Library and PDF publication.

← [Back to Architecture Index](README.md)
