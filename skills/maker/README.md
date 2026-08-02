# RGB Maker Skill

Status: Contract defined, runtime not implemented. See
[`SKILL.md`](SKILL.md) for the full contract (truth layers, the
image/text authority rule, and what Maker must never do),
[`prompts/`](prompts/) for its two core prompts (extract entities, detect
conflicts), [`schemas/`](schemas/) and [`templates/`](templates/) for its
four package kinds (entity package, maker report, visual observation,
conflict report), and [`examples/`](examples/) for a worked example.

This package does not yet run — it should not be presented as an available
product until a runtime implementation exists.

**Sequencing:** per
`.analysis/refined/20260801-specialist-first-skill-roadmap`, Maker's
*runtime* work is deferred until [`../specialist/`](../specialist/) has a
source-trace benchmark (see
[`../specialist/benchmark/golden-qa.yaml`](../specialist/benchmark/golden-qa.yaml)),
the bundle/search context shape is stable, and Maker's provenance/canon
schemas (already drafted in [`schemas/`](schemas/)) are accepted. The
contract in this package — including the authority rule and truth layers —
is written and stable; what's deferred is building the engine that
executes it.
