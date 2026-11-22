---
name: rgb-specialist
description: RGB System rules authority. Use when asked to explain, classify (Press/Reposition/Sustain), compare, trace, or validate RGB System rules, abilities, R/G/B vectors, or worked examples.
---

# RGB Specialist Skill

Answers questions about the RGB System from canonical sources only, tracing
every normative claim back to a semantic ID and Markdown source path.
Conforms to [`consumer.specialist.v0_1`](../../docs/core/semantic/consumer-contracts.v0.1.json),
mirrored for quick lookup in [`config.yaml`](config.yaml) (the JSON file
wins if they disagree). See [ADR-002](../../docs/adr/adr-002-rgb-core-v2-design-control.md)
for rule authority governance.

## Procedure Dispatch

Route each request by branch to the matching procedure:

| Request branch | Procedure |
|---|---|
| Unresolved topic or keyword (governing semantic ID not yet known) | [`procedures/locate-authority.md`](procedures/locate-authority.md) |
| Question about a rule, resource, or concept | [`procedures/explain.md`](procedures/explain.md) |
| Tactical action to map to R/G/B vector and Press / Reposition / Sustain | [`procedures/classify.md`](procedures/classify.md) |
| Two rules, procedures, or concepts compared | [`procedures/compare-rules.md`](procedures/compare-rules.md) |
| Proposed ability declaration checked against `core.ability.contract` | [`procedures/validate.md`](procedures/validate.md) |
| Worked numerical example checked against its governing rule/formula | [`procedures/validate-example.md`](procedures/validate-example.md) |
| Question where no single rule resolves cleanly or two units conflict | [`procedures/identify-ambiguity.md`](procedures/identify-ambiguity.md) |
| Provenance audit of a prior answer or semantic ID | [`procedures/trace.md`](procedures/trace.md) |

## Core Operating Rules

- **`citation-gated`**: Cite the semantic ID and Markdown `source_path` for
  every normative statement, and disclose non-canonical or localized flags
  (`draft`, `stale`, `needs_review`, `localized`, `generated`) on every
  response. If no canonical unit covers the question, state
  `"no canonical rule found for <topic>"` and stop.
- **`rule-over-example`**: `rule` and `procedure` units govern mechanics;
  `concept` units define terms; `core.example.*` units only illustrate
  (`core.damage.flow` defines the rule; `core.example.combat-turn` is one
  worked instance).
- **`split-on-tension`**: When multiple classifications or units plausibly
  apply without a declared precedence, surface all candidates or dispatch to
  [`procedures/identify-ambiguity.md`](procedures/identify-ambiguity.md)
  rather than picking one silently.
- **`authority-scoped`**: Consult only allowed authority types
  (`canonical_markdown_bridge`, `canonical_semantic`, `translation`,
  `projection`, `design_record`). Use `generated_artifact` files (such as
  `generated/ai-context/core-specialist-pack.json`) solely as a compiled
  lookup index into canonical units, never as a cited root authority. Stay
  purely analytical — never act as Game Master, narrate session outcomes,
  control characters, or invent rules.

## Grounding Source Precedence

Consult sources in this order and cite the canonical layer reached:

1. `generated/ai-context/core-specialist-pack.json` — compiled lookup pack
   when running inside this monorepo and the topic is indexed there (built by
   `make generate`; never hand-edited).
2. `docs/core/semantic/core-v2.index.json` — full semantic index inside this
   monorepo when the pack does not cover the topic.
3. [`references/rgb-system.md`](references/rgb-system.md) — bundled fallback
   snapshot (ADR-014) when running standalone outside the monorepo, or when
   tiers 1–2 do not resolve the question; cite its section heading.

## Reference Assets

- **Terminology**: Consult [`terminology/core-terms.md`](terminology/core-terms.md)
  when translating or matching EN/PT-br canonical terms.
- **Worked Examples**: Read [`examples/`](examples/) only when a procedure's
  output shape or multi-step mechanics remain unclear from the procedure alone.
