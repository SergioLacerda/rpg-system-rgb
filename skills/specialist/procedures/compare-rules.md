# Procedure: Compare Rules

Explain how two rules, procedures, or concepts differ or interact.

## Input

Two semantic IDs (or two questions resolved first via
[`locate-authority.md`](locate-authority.md)).

## Steps

1. Locate both units in `docs/core/semantic/core-v2.index.json` (and, for
   `rule`/`procedure`/`ability` kinds, `docs/core/semantic/source/core-v2-rules.v0.1.json`
   for full `statement`/`formula`/`ordered_steps`).
2. Present each unit's `kind` and one-line summary side by side so facts from
   unit A and unit B stay clearly attributed.
3. Inspect `relationships` verbatim (`depends_on`, `clarifies`, etc.) to
   explain declared interactions (e.g. `core.damage.armor-reduction` and
   `core.damage.shield-absorption` are sequenced steps under
   `core.damage.flow`).
4. If `relationships` declares no link between the two units, report
   `"no declared relationship"` explicitly rather than inferring one.
5. If both rules plausibly apply to the same situation in conflicting ways,
   hand off to [`identify-ambiguity.md`](identify-ambiguity.md) per
   `split-on-tension`.
6. Attach both semantic IDs, their `source_path`s, and all applicable
   status/locale disclosures (`citation-gated`).

## Output Shape

```text
A: <id> (<kind>) — <one-line summary>
B: <id> (<kind>) — <one-line summary>
Relationship: <depends_on | clarifies | sequence | "no declared relationship">
Cited: <both ids and their source_paths>
Disclosures: <draft | stale | needs_review | localized | generated, or "none (canonical en)">
```

## Worked Example

Read [`../examples/compare-armor-vs-shield.md`](../examples/compare-armor-vs-shield.md)
only when comparing two sequenced combat or damage-resolution steps.
