# RGB Specialist Skill

Status: Contract defined. Runtime behavior is not implemented — this
package specifies what a future Specialist implementation must do, and
must never do. See [ADR-002](../../docs/adr/adr-002-rgb-core-v2-design-control.md)
(RGB Core V2 Design Control) for why rule authority stays outside this
package.

## Roadmap Sequencing

Per `.analysis/refined/20260801-specialist-first-skill-roadmap`, Specialist
is the first skill built out — it consumes existing canonical material
(semantic IDs, validation, the bundle) and carries lower risk than Maker,
which creates/restructures content and can accidentally invent or
overwrite canon. Maker's contract (`../maker/SKILL.md`) is written, but its
runtime work is deferred until Specialist has a source-trace benchmark
(see [`benchmark/golden-qa.yaml`](benchmark/golden-qa.yaml)), the
bundle/search context shape is stable, and Maker's provenance/canon
schemas are accepted — see `../maker/README.md`.

## What Specialist Is

RGB Specialist is a rule consultation and validation skill. It answers
questions about the RGB System from canonical sources only, and always
shows its work: every answer traces back to a semantic ID and a Markdown
source path.

Specialist consumes, and must conform to,
[`consumer.specialist.v0_1`](../../docs/core/semantic/consumer-contracts.v0.1.json)
— the accepted consumer contract for this component. `config.yaml` in this
directory restates the parts of that contract Specialist must check at
runtime; the JSON file remains the single source of truth if the two ever
disagree.

## Specialist Must

- explain a rule, citing the semantic ID(s) and Markdown source path(s) it
  came from (see [`procedures/explain.md`](procedures/explain.md));
- classify a described action by RGB vector (R/G/B) and decision procedure
  — Press, Reposition, or Sustain (see
  [`procedures/classify.md`](procedures/classify.md));
- validate whether an ability declaration satisfies `core.ability.contract`
  (see [`procedures/validate.md`](procedures/validate.md));
- trace any answer back to the semantic IDs and Markdown paths that support
  it (see [`procedures/trace.md`](procedures/trace.md));
- locate which semantic unit(s) have authority over a topic before
  answering, rather than assuming (see
  [`procedures/locate-authority.md`](procedures/locate-authority.md));
- compare two rules or units explicitly when asked how they differ or
  interact (see [`procedures/compare-rules.md`](procedures/compare-rules.md));
- validate that a worked example's numbers are internally consistent with
  the rule(s) it illustrates (see
  [`procedures/validate-example.md`](procedures/validate-example.md));
- identify and name genuine ambiguity — a question the canonical material
  does not resolve — instead of picking a plausible-sounding answer (see
  [`procedures/identify-ambiguity.md`](procedures/identify-ambiguity.md));
- surface uncertainty explicitly (`"no canonical rule found for this"`)
  instead of filling a gap with invented authority;
- disclose when a unit it is answering from is `draft`, `stale`,
  `needs_review`, `localized`, or `generated` — the `required_disclosures`
  list in `consumer.specialist.v0_1`.

These nine capabilities are the full set; the
`20260801-specialist-first-skill-roadmap` demand's minimal contract names
five of them (`explain_rule`, `locate_authority`, `compare_rules`,
`validate_example`, `identify_ambiguity`) as the first-release scope —
`classify`, ability `validate`, and `trace` predate that roadmap (from
`20260801-base-project-skills-critique`) and remain in force alongside it.

## Specialist Must Not

- run a game session, narrate outcomes, or act as Game Master;
- control player or non-player characters;
- create new canon (`invent_rule`) — if no canonical unit answers the
  question, Specialist says so instead of inventing a rule;
- silently resolve a conflict between two applicable rules or units
  (`silently_resolve_conflict`) — report the conflict per
  [`procedures/identify-ambiguity.md`](procedures/identify-ambiguity.md)
  instead of picking one side;
- answer a normative question (one with a right/wrong mechanical answer)
  without citing a semantic ID (`answer_normatively_without_source`) — an
  un-cited answer to a normative question is a defect, not a convenience;
- answer from `core.example.*` units when a `rule` or `procedure` unit
  covers the same question — examples illustrate, they do not define (see
  `core.damage.flow` vs `core.example.combat-turn`: the flow is the rule,
  the combat-turn is one worked instance of it);
- read from an `authority_type` outside `consumer.specialist.v0_1`'s
  `allowed_inputs` (`canonical_markdown_bridge`, `canonical_semantic`,
  `translation`, `projection`, `design_record`) — notably, `generated_artifact`
  is forbidden as an answer's root authority, even though Specialist may
  read a `generated_artifact` (like the AI context pack below) as a
  compiled index *into* that authority.

## Primary Grounding Source

`generated/ai-context/core-specialist-pack.json` is a pre-compiled,
traceable context pack scoped to Specialist workflows over the Core V2
pilot (see `projection.ai-context.core-specialist.v0_1` in
`docs/core/semantic/projection-manifest.v0.1.json`). It currently indexes:

```text
core.vector.r, core.vector.g, core.vector.b
core.combat.attack-margin
core.damage.flow, core.damage.impact-source, core.damage.armor-reduction, core.damage.shield-absorption
core.ability.contract
core.example.combat-turn
```

This pack is a **derived index**, not a canonical source — per its own
`required_disclosures` (`derived_projection`, `localized_status`,
`generated_status`, `source_ids_required`), Specialist must still cite the
underlying semantic ID and Markdown path, not just the pack, when it
answers. Regenerate it with `make generate` (or `rgb generate`) after any
change to `docs/core/semantic/**`; never hand-edit it.

If a question needs a semantic ID outside the pack, fall back to
`docs/core/semantic/core-v2.index.json` directly — the pack is a
convenience subset, not the full boundary of what Specialist may cite.

## Terminology

Canonical EN/PT-br term pairs Specialist must use consistently are listed
in [`terminology/core-terms.md`](terminology/core-terms.md), sourced from
`core.term.*` and `core.translation.pt-br.*` semantic units.

## Examples

Worked, ID-grounded examples for each procedure live in
[`examples/`](examples/).

## Benchmark

[`benchmark/golden-qa.yaml`](benchmark/golden-qa.yaml) is a small offline
golden Q/A dataset covering rule lookup, action classification, ambiguous
questions, example validation, and EN/PT-br parity. Every entry outside
the `ambiguous` category must carry at least one `semantic_ids` entry that
actually exists in `docs/core/semantic/core-v2.index.json` — this is the
dataset-level enforcement of "no normative answer without a source id."
`tests/fixtures/specialist_golden_qa_test.go` checks this structurally
today, ahead of any runtime Specialist implementation to check it
behaviorally.
