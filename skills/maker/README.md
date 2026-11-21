# RGB Maker Skill

Status: Contract defined in `skills/maker/**`; deterministic Go validation,
diffing, aggregation, and input routing are implemented in
`internal/components/maker/` per
[ADR-015](../../docs/adr/adr-015-skill-runtime-split.md), while LLM-side
runtime wiring remains deferred. See
[`SKILL.md`](SKILL.md) for the full contract (truth layers, the
image/text authority rule, and what Maker must never do),
[`config.yaml`](config.yaml) for the consumer contract restatement,
[`prompts/`](prompts/) for its two core prompts (extract entities, detect
conflicts), [`schemas/`](schemas/) and [`templates/`](templates/) for its
four package kinds (entity package, maker report, visual observation,
conflict report), and [`examples/`](examples/) for a worked example.

The LLM-side skill workflow does not yet run end-to-end — it should not be
presented as an available end-user product until an invocable LLM runtime
is wired.

**Runtime architecture:** the LLM-instruction vs. Go-component split for
Maker's eventual runtime is decided in
[ADR-015](../../docs/adr/adr-015-skill-runtime-split.md). Read that decision
before starting the runtime implementation, rather than re-deriving the split.

**Sequencing:** per the internal `specialist-first-skill-roadmap` planning
analysis, Maker's *runtime* work is deferred until
[`../specialist/`](../specialist/) has a
source-trace benchmark (see
[`../specialist/benchmark/golden-qa.yaml`](../specialist/benchmark/golden-qa.yaml)),
the bundle/search context shape is stable, and Maker's provenance/canon
schemas (already drafted in [`schemas/`](schemas/)) are accepted. The
contract in this package — including the authority rule and truth layers —
is written and stable; what's deferred is building the engine that
executes it.
