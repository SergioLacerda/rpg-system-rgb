# RGB Maker Skill

Status: Contract defined. Runtime behavior is not implemented — this
package specifies what a future Maker implementation must do, and must
never do. See [ADR-001](../../docs/adr/adr-001-ai-first-documentation-authority.md)
(AI-First Documentation Authority) for the authority model this contract
extends into content creation.

**Deferred by design:** per
`.analysis/refined/20260801-specialist-first-skill-roadmap`, Maker's
runtime implementation is deferred until [`../specialist/`](../specialist/)
proves out the rule/source/bundle path (source-trace benchmark, stable
bundle/search context shape, accepted provenance schemas) — see
[`../specialist/SKILL.md`](../specialist/SKILL.md)'s "Roadmap
Sequencing." This contract stays written and ready in the meantime.

## What Maker Is

RGB Maker is a content structuring skill. It turns a GM's notes, text, or
images into campaign/world artifacts — NPCs, locations, items, factions —
while keeping a hard line between what the source material actually said
and what Maker inferred or generated.

Maker consumes, and must conform to,
[`consumer.maker.v0_1`](../../docs/core/semantic/consumer-contracts.v0.1.json)
— the accepted consumer contract for this component.

## Maker Must

- extract entities from notes, text, or images into structured packages
  (see [`schemas/entity-package.schema.yaml`](schemas/entity-package.schema.yaml)
  and [`templates/entity-package.md`](templates/entity-package.md));
- classify every extracted fact into exactly one of four layers:
  **explicit fact**, **inferred suggestion**, **unresolved field**, or
  **conflict** — never blend these into a single undifferentiated
  description;
- generate both a human-readable Markdown artifact and its machine-readable
  metadata for every package (schema + template pair, not one or the
  other);
- preserve the three truth layers from `design.md`'s Maker contract:
  - **canonical truth** — what the RGB System's own rules say (Maker never
    overrides this; see Non-Goals);
  - **public or scoped knowledge** — what characters/factions in the
    fiction are meant to know;
  - **observed representation** — what a specific note or image shows,
    which may be incomplete, biased, or wrong in-fiction;
- report every ambiguity as an unresolved field or a conflict report (see
  [`schemas/conflict-report.schema.yaml`](schemas/conflict-report.schema.yaml)),
  never silently resolved by guessing;
- disclose `draft`, `stale`, `needs_review`, and `generated` status on
  every output, per `consumer.maker.v0_1`'s `required_disclosures`.

## Maker Must Not

- implement runtime behavior in this package (this is a contract, not an
  engine);
- become mandatory for using the RGB System — Maker is an authoring aid,
  not a dependency of play;
- act as rule authority: Maker may *read* `core.ability.contract` and other
  canonical units to keep generated content consistent, but may never
  originate a new rule, and its outputs' `authority_type` is never
  `canonical_semantic` or `canonical_markdown_bridge` (`consumer.maker.v0_1`
  forbids `authority_type: projection` and `generated_artifact` on its
  *inputs* for the same reason — Maker cannot bootstrap authority from
  something already derived);
- silently overwrite an explicit fact with an inference, or an inference
  with a guess — every promotion from a lower-confidence layer to a higher
  one is a decision the user makes, not Maker.

## The Authority Rule

```text
image = appearance
text = facts
```

When a note includes both an image and text describing the same subject,
the **text is the authoritative source for facts** (names, stats,
relationships, canon status); the **image is authoritative only for
appearance** (what something looks like). An image never overrides a
textual fact, and a textual description never overrides what an image
actually shows about appearance. When image and text conflict on a fact
the image cannot actually establish (e.g. text says "the merchant is
lying," an image cannot contradict that), the text wins outright — this is
not a case that needs a conflict report, because the image was never
authoritative over that field.

See [`schemas/visual-observation.schema.yaml`](schemas/visual-observation.schema.yaml)
for how an image's contribution is recorded as appearance-only, and
[`templates/visual-observation.md`](templates/visual-observation.md) for
the corresponding human-readable form.

## Package Kinds

| Kind | Schema | Template | Purpose |
|---|---|---|---|
| Entity Package | [`schemas/entity-package.schema.yaml`](schemas/entity-package.schema.yaml) | [`templates/entity-package.md`](templates/entity-package.md) | one extracted entity (NPC, location, item, faction), with facts/inferences/unresolved fields |
| Maker Report | [`schemas/maker-report.schema.yaml`](schemas/maker-report.schema.yaml) | [`templates/maker-report.md`](templates/maker-report.md) | session-level summary of everything Maker extracted from one input batch |
| Visual Observation | [`schemas/visual-observation.schema.yaml`](schemas/visual-observation.schema.yaml) | [`templates/visual-observation.md`](templates/visual-observation.md) | appearance-only facts read from an image, per the authority rule |
| Conflict Report | [`schemas/conflict-report.schema.yaml`](schemas/conflict-report.schema.yaml) | [`templates/conflict-report.md`](templates/conflict-report.md) | two or more sources disagree on a fact Maker cannot resolve on its own |

## Prompts

Starting prompts for the two core Maker tasks live in
[`prompts/`](prompts/): [`extract-entities.md`](prompts/extract-entities.md)
and [`detect-conflicts.md`](prompts/detect-conflicts.md).

## Examples

A worked example (raw GM notes → entity package → conflict report) lives in
[`examples/`](examples/).

## Go Boundary

Per `design.md`'s Skills Track design, Go is used for Maker only where
deterministic behavior is required — schema validation, semantic ID lookup,
bundle/context generation. There is no Maker-specific generated projection
yet (unlike Specialist's `core-specialist-pack.json`); until Maker has a
real consumer beyond this contract, adding one would be exactly the
"attractive but unused generated artifact" risk
[ADR-007](../../docs/adr/adr-007-bundle-consumer-contract-policy.md) warns
against. Maker reads `docs/core/semantic/core-v2.index.json` directly for
now.
