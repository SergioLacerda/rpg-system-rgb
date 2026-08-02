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

## Public Status

| Surface | Status | Notes |
| --- | --- | --- |
| RGB Core V2 | Implemented | Playable rules, examples, Go rule-engine tests, and semantic docs validation exist. |
| Semantic tooling | Implemented | `make validate` and `make generate` validate and produce indexed projections from `docs/core/semantic/**`. |
| Bundle output | Experimental | `make bundle` produces a minimal semantic bundle; ADR-007 records the current consumer-contract limits. |
| Library HTML | Implemented | MkDocs builds the public Library into `web/landing/public/library/`; Astro packages it into the landing build. |
| PDF downloads | Implemented | `make docs-pdf` publishes latest and versioned PDFs with manifest, checksums, and editorial smoke checks. |
| Landing | Implemented | Astro routes, bilingual UI, Library links, PDF links, tests, and static builds exist. |
| Maker | Contract defined | Contract, schemas, templates, prompts, and examples exist; runtime behavior is not implemented. |
| Specialist | Contract defined | Contract, procedures, terminology, examples, and benchmark scaffolding exist; runtime behavior is not implemented. |
| PDF/UA tagging | Deferred | Full accessibility tagging is outside the first editorial PDF quality gate. |

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
   rgb/              unified CLI: status, validate, generate, bundle (ADR-006)
   rgb-tooling/      deprecated compatibility path for validate/generate/bundle

internal/
   app/              thin orchestration layer (only import path for cmd/*)
   components/
      core/          RGB rule engine: characters, actions, combat, states
      tooling/       semantic-docs validators + projection generator
      maker/         planned campaign/content structuring boundary
      specialist/    planned rule consultation and validation boundary
      bundles/       experimental semantic bundle output (ADR-007)

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
   maker/            RGB Maker skill contract; runtime not implemented
   specialist/       RGB Specialist skill contract; runtime not implemented

web/
   landing/          implemented Astro landing page, Library shell, and PDF download surface

generated/            derived artifacts only, never edited by hand
   ai-context/       AI context packs
   bundles/          consolidated projection bundle
   landing/          landing page summary data
   library/          core-v2-rules.json semantic projection (make generate)
   pdf/              core-v2-rules.manifest.json PDF export manifest (make generate)
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

MkDocs + Astro is the current publication path for the HTML Library and PDF
downloads; there is no separate Go-rendered HTML/print output to generate first.

```bash
make docs-build  # build the MkDocs Library into web/landing/public/library
make landing-build # embed the Library into the Astro static site
make docs-pdf    # generate latest PDF downloads from MkDocs (requires docs-pdf-install)
```

Automated PDF generation is available through a separate dependency profile:

```bash
make docs-pdf-install
make docs-pdf
```

`make docs-pdf` uses MkDocs PDF configuration files and writes latest aliases
and versioned PDFs under `web/landing/public/downloads/`. It also refreshes
the release manifest and checksums used by `make pdf-editorial-check`.

`make pdf-publish` remains available for publishing a reviewed external PDF
source when needed, but the normal release path is automated.

### Full gate

```bash
make check-fast       # fast Go, architecture, semantic, and coverage gate
make check            # check-fast + web, docs, generated drift, workflow, shell, governance files
make release-check    # check + PDF build and editorial artifact gate
```

## Contributing

Contributions are welcome.

Examples of useful contributions:

- rule clarifications
- additional combat examples
- translation improvements
- optional modular systems (magic, technology, powers)

Please read [CONTRIBUTING.md](CONTRIBUTING.md) before opening an issue or pull
request.

## License

See [LICENSE](LICENSE) and [LICENSES.md](LICENSES.md) for the current licensing
model and the deferred split-license decision.
