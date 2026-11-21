# RGB Maker Skill

Status: Contract defined in `skills/maker/**`; deterministic Go input
routing, schema validation, canonical dependency validation, structured-fact
diffing, report aggregation, Markdown rendering (with `conflicts/resolved/`
archival and stale-file pruning), workspace statistics, compact chat-pack
export, and `rgb maker process|stats|export` CLI orchestration live in
`internal/components/maker/`, `internal/app/`, and `cmd/rgb/` per
[ADR-015](../../docs/adr/adr-015-skill-runtime-split.md). End-to-end LLM
runtime wiring remains deferred.

See [`SKILL.md`](SKILL.md) for the agent contract, [`config.yaml`](config.yaml)
for the consumer contract restatement, [`prompts/`](prompts/) for extraction
and conflict prompts, [`schemas/`](schemas/) and [`templates/`](templates/)
for the four package kinds, and [`examples/`](examples/) for a worked
example.

The LLM-side skill workflow does not yet run end-to-end — it should not be
presented as an available end-user product until an invocable LLM runtime
is wired.

## Sequencing & Architecture

**Runtime architecture:** the LLM-instruction vs. Go-component split is
decided in [ADR-015](../../docs/adr/adr-015-skill-runtime-split.md). Read
that decision before extending the runtime rather than re-deriving the
split.

**Generated projection boundary:** Maker reads
`docs/core/semantic/core-v2.index.json` directly rather than adding a
Maker-specific generated projection, avoiding unused derived artifacts per
[ADR-007](../../docs/adr/adr-007-bundle-consumer-contract-policy.md).

**Sequencing:** per the internal `specialist-first-skill-roadmap` planning
analysis, Maker's Go deterministic layer is built while full LLM workflow
wiring remains gated behind [`../specialist/`](../specialist/)'s runtime
path (source-trace benchmark in
[`../specialist/benchmark/golden-qa.yaml`](../specialist/benchmark/golden-qa.yaml)
and stable bundle/search context shape).
