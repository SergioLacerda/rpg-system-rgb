# Containers

```mermaid
C4Container
    title RGB System V2 — Containers (What's Actually Built)

    Person(gm, "Game Master / Player")

    System_Boundary(rgb, "RGB System V2") {
        Container(docs, "docs/core/**", "Markdown + YAML", "Canonical bilingual rule source (EN bridge, PT-br projection)")
        Container(semantic, "docs/core/semantic/**", "JSON", "Stable IDs, l10n manifest, projection manifest, consumer contracts")
        Container(tooling, "cmd/rgb-tooling", "Go binary", "Validates docs/semantic contracts; generates projection JSON")
        Container(compiler, "cmd/rgb-compiler", "Go binary", "Renders docs/core/**/*.md into static HTML + print-ready pages")
        Container(rgb_cli, "cmd/rgb", "Go binary", "Core rule engine entrypoint")
        Container(generated, "generated/**", "JSON + HTML", "Derived: bundles, ai-context, landing, library HTML/print, pdf manifest, search")
    }

    Rel(gm, docs, "Reads")
    Rel(docs, semantic, "Indexed by")
    Rel(tooling, docs, "Validates")
    Rel(tooling, semantic, "Reads and writes")
    Rel(tooling, generated, "Generates projections into")
    Rel(compiler, docs, "Renders")
    Rel(compiler, generated, "Writes HTML/print pages into")
    Rel(gm, generated, "Reads compiled site or PDF export instructions")
```

`cmd/rgb` (the Core rule engine entrypoint) does not yet consume
`docs/core/**` at runtime — RGB Core V2's rule logic lives in
`internal/components/core` and is validated against canonical docs via
Gherkin/BDD tests (`tests/features/`, `tests/core_behavior/`), not by
`cmd/rgb` reading Markdown directly. This is shown as a separate,
currently-unconnected container rather than omitted, since it is a real,
built binary.

← [Back to Architecture Index](README.md)
