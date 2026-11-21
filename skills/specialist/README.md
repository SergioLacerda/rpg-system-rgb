# RGB Specialist Skill

Status: Installable package, procedure runtime not implemented.
`SKILL.md` carries front-matter (`name`/`description`), and the package is
self-contained — [`references/rgb-system.md`](references/rgb-system.md)
bundles a snapshot of the canonical rule documentation from `docs/core/**`
(see ADR-014), so the folder keeps working when copied outside this
monorepo. See [`SKILL.md`](SKILL.md) for the agent contract and dispatch
table, [`procedures/`](procedures/) for its eight core procedures
(`explain`, `classify`, `validate`, `trace`, `locate-authority`,
`compare-rules`, `validate-example`, `identify-ambiguity`),
[`terminology/`](terminology/) for canonical EN/PT-br terms,
[`examples/`](examples/) for worked examples, and [`config.yaml`](config.yaml)
for the runtime contract summary. Required status disclosures (`draft`,
`stale`, `needs_review`, `localized`, `generated`) are enforced across all
procedure outputs rather than as a separate procedure.

This package does not yet run — it should not be presented as an available
product until a runtime implementation exists.

Once a runtime exists, distribution follows
[ADR-013: Skill Distribution Via The Existing Publication Pipeline](../../docs/adr/adr-013-skill-distribution-via-existing-publication-pipeline.md) —
a downloadable `.zip` published alongside the PDF/Library downloads, not a
separate GitHub Release or installer CLI.

## Roadmap Sequencing

Per the internal `specialist-first-skill-roadmap` planning analysis,
Specialist is the first skill built out — it consumes existing canonical
material (semantic IDs, validation, the bundle) and carries lower risk than
Maker, which creates/restructures content and can accidentally invent or
overwrite canon. Maker's contract (`../maker/SKILL.md`) is written alongside
its Go validation/aggregation layer (`internal/components/maker/`), while
end-to-end LLM runtime wiring stays gated on Specialist's benchmark and
runtime path — see [`../maker/README.md`](../maker/README.md).

## Benchmark

[`benchmark/golden-qa.yaml`](benchmark/golden-qa.yaml) is an offline golden
Q/A dataset covering rule lookup, action classification, ambiguous
questions, example validation, and EN/PT-br parity. Every entry outside
the `ambiguous` category carries at least one `semantic_ids` entry present
in `docs/core/semantic/core-v2.index.json` and at least one
`expected_answer_traits` entry. `tests/fixtures/specialist_golden_qa_test.go`
validates both structurally ahead of a runtime Specialist implementation.

## Installation

This package (`SKILL.md` plus `procedures/`, `references/`, `terminology/`,
`examples/`, `benchmark/`, `config.yaml`) is plain Markdown/YAML and does not
depend on any single vendor runtime. To install it with any LLM or agent
tool that loads local skill folders:

1. Copy or link the `skills/specialist/` folder into your tool's local
   skills/context directory.
2. Point the assistant at [`SKILL.md`](SKILL.md) as the entry point.
3. Outside this monorepo, the agent automatically falls back to
   [`references/rgb-system.md`](references/rgb-system.md) (see Grounding
   Source Precedence in `SKILL.md`).

If your host tool uses custom slash commands, generate a thin host-specific
wrapper that forwards the user's question to `SKILL.md`.
