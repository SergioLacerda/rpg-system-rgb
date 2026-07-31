# RGB System

![status](https://img.shields.io/badge/status-active-brightgreen)
![license](https://img.shields.io/badge/license-open-blue)
![rpg](https://img.shields.io/badge/system-RGB-orange)

A lightweight **tabletop role‑playing system** designed to keep rules simple while enabling **tactical combat**, **modular design**, and **adaptable worlds**.

## Languages

This project documentation is available in multiple languages.

- 🇬🇧 English (official) → [English Documentation](docs/core/en/)
- 🇧🇷 Português → [Documentação em Português](docs/core/PT-br/)

## Quick Navigation

- Quick Start → [Quick Start](docs/core/en/introduction/quick_start.md)
- System Overview → [System Overview](docs/core/en/introduction/system_overview.md)
- One page Rules → [One page Rules](docs/core/en/introduction/rgb_one_page_rules.md)

- Core Rules → [Core Rules](docs/core/en/core/)
- Combat System → [Combat System](docs/core/en/combat/)
- Damage Model → [Damage Model](docs/core/en/combat/damage_model.md)
- Combat Decision Model → [Combat Decision Model](docs/core/en/combat/combat_decision_model.md)

- Equipment → [Equipment](docs/core/en/equipment/)
- Weapons → [Weapons](docs/core/en/weapons/)

- Examples & Reference → [Examples & Reference](docs/core/en/reference/)
- Damage Model and Interaction Design Notes → [Damage Model and Interaction Design Notes](docs/core/en/reference/rgb_damage_interaction_model.md)
- Gameplay Loop → [Gameplay Loop](docs/core/en/reference/gameplay_loop.md)
- RGB System Architecture → [RGB System Architecture](docs/core/en/reference/rgb_system_architecture_notes.md)
- RGB Interaction Model → [RGB Interaction Model](docs/core/en/reference/rgb_damage_interaction_model.md)
- RGB System Engine → [RGB System Engine](docs/core/en/reference/rgb_system_engine.md)

## Core Concept

The **RGB System** is based on three fundamental vectors.

```text
Vector  Name     Function
------  -------- -------------------------------------------------------
R       Red      attacks, power, strength tests, hit chance, damage dealt
G       Green    agility, mobility, dodge, reaction speed
B       Blue     magical shields, special protection, energy manipulation,
                 intellect, wisdom, resistance
```

These vectors define how characters interact with the world and combat.

---

## Design Principles

The RGB system was designed with four guiding principles:

- **Structural simplicity**
- **Tactical depth**
- **Modularity**
- **Scalability**

The goal is to allow game masters and players to adapt the system to different genres without rewriting core mechanics.

---

## Supported Campaign Types

Because the RGB System separates **core mechanics** from **setting rules**, it can easily support multiple genres:

- Modern campaigns
- Fantasy settings
- Science fiction worlds
- Super‑powered universes

The same rule structure can be reused across different campaign types with minimal adjustments.

## Philosophy

RGB intentionally avoids excessive rule complexity.  
Instead it focuses on:

- clear numeric relationships
- tactical positioning
- flexible character design
- modular mechanics

This allows the system to work both for **story‑driven campaigns** and **tactical combat playstyles**.

## Project Structure

```text
rpg-system-rgb/
cmd/
   rgb/              core rule engine entrypoint
   rgb-compiler/     renders docs/core/**/*.md into HTML + print pages
   rgb-tooling/      validates docs/core/semantic/**, generates projections

internal/
   app/              thin orchestration layer (only import path for cmd/*)
   components/
      core/          RGB rule engine: characters, actions, combat, states
      tooling/       semantic-docs validators + projection generator
      compiler/      Markdown parser + HTML/print renderer
      library/       site assembly + PDF export instructions
      maker/         campaign/content structuring (stub, not yet implemented)
      specialist/    rule consultation and validation (stub, not yet implemented)
      bundles/       consolidated bundle format (stub, not yet implemented)

docs/
   core/
      en/            canonical bilingual rule source (bridge)
         combat/
         core/
         equipment/
         reference/
         weapons/

      PT-br/         localized projection
         combat/
         core/
         equipment/
         reference/
         weapons/

      semantic/      stable IDs, l10n manifest, projection manifest, contracts

   adr/               architecture decision records
   architecture/       C4 diagrams (context, containers, components)
   engineering/         review workflow and guardrail documentation

skills/
   maker/            RGB Maker skill (placeholder, not yet implemented)
   specialist/       RGB Specialist skill (placeholder, not yet implemented)

web/
   landing/          unified landing page (placeholder, not yet implemented)

generated/            derived artifacts only, never edited by hand
   ai-context/       AI context packs
   bundles/          consolidated projection bundle
   landing/          landing page summary data
   library/          compiled HTML site + print-ready pages
   pdf/              PDF export manifest
   search/           search index

tests/                repo-wide guardrails (architecture, semantic docs, features)

Makefile
go.mod
.golangci.yml
LICENSE
README.md
```

The **English documentation under `docs/core/en` is the official reference
version** of the system.

Translations aim to remain synchronized with the English documentation.

## Development

Common commands (see `Makefile` for the full list):

### Tests

```bash
make test        # run the full test suite
make test-arch   # run the Clean Architecture boundary guardrail tests
make cover       # run tests with a coverage report
```

### Documentation guardrails

```bash
make validate    # validate docs/core/semantic/** (l10n, index, contracts, projections)
make generate    # regenerate the semantic projection JSON artifacts
```

### HTML site and PDF

```bash
make compile     # render docs/core/**/*.md into a static HTML site and print-ready pages
make docs-pdf    # generate latest PDF downloads from MkDocs (requires docs-pdf-install)
```

`make compile` writes:

- `generated/library/html/{en,PT-br}/**` — static HTML site (1:1 mirror of `docs/core/**/*.md`), with per-locale navigation index pages
- `generated/library/print/core-v2-rules-{en,PT-br}.html` — print-ready pages, ordered per `docs/core/semantic/projection-manifest.v0.1.json`
- `generated/library/PDF_EXPORT.md` — instructions for exporting the print-ready pages to PDF (print-to-PDF from a browser; PDF export is a documented manual step by design, not an automated build step)

After exporting a PDF, publish it to the landing static assets with:

```bash
make pdf-publish PDF_SRC=RGB-Core-V2-Rules-PT-br.pdf PDF_LOCALE=pt-br PDF_VERSION=v0.2
make pdf-publish PDF_SRC=RGB-Core-V2-Rules-en.pdf PDF_LOCALE=en PDF_VERSION=v0.2
```

The command creates a versioned PDF and updates the locale-specific `latest`
alias consumed by the landing page.

Automated PDF generation is available through a separate dependency profile:

```bash
make docs-pdf-install
make docs-pdf
```

`make docs-pdf` uses MkDocs PDF configuration files and writes latest aliases
under `web/landing/public/downloads/`. This path is intentionally separate from
the base `make docs-build` flow because the PDF plugin depends on additional
rendering libraries documented by the plugin/WeasyPrint stack.

### Full gate

```bash
make review-structure   # vet + lint + test + test-arch + validate
```

## Contributing

Contributions are welcome.

Examples of useful contributions:

- rule clarifications
- additional combat examples
- translation improvements
- optional modular systems (magic, technology, powers)

Please open an issue or submit a pull request.

## License

See the **LICENSE** file for details.
