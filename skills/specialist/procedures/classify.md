# Procedure: Classify

Classify a described action by RGB vector and tactical decision procedure.

## Input

A short description of what a character does (e.g. "the fighter charges
in and swings a greatsword" or "the scout ducks behind cover").

## Steps

1. Match the described action against the three decision procedures:

   | Procedure | Vector | Semantic ID | Meaning |
   |---|---|---|---|
   | Press | R | `core.action.press` | changes the source of pressure — attack, force, disruption |
   | Reposition | G | `core.action.reposition` | changes relation to pressure — movement, evasion, positioning |
   | Sustain | B | `core.action.sustain` | preserves continuity under pressure — endurance, shields, stabilization |

2. If the action clearly maps to one procedure, classify it and cite the
   matching `core.action.*` unit and its `depends_on` vector
   (`core.vector.r`, `core.vector.g`, or `core.vector.b`) per
   `citation-gated`.
3. Per `split-on-tension`, if the action plausibly spans multiple procedures
   (e.g. "shove an enemy off a ledge" touches Press and Reposition), list all
   plausible classifications and cite each corresponding `core.action.*` unit
   rather than picking one silently.
4. If the action is pure roleplay with no mechanical resolution, state that
   tactical classification does not apply.
5. Record any non-canonical or localized status flags on the cited units in
   `Disclosures`.

## Output Shape

```text
Action: <restated action>
Classification: <Press|Reposition|Sustain> (<vector>) | <multiple candidates> | "none (non-mechanical)"
Cited: <core.action.* id> (<source_path>)
Disclosures: <draft | stale | needs_review | localized | generated, or "none (canonical en)">
```

## Worked Example

Read [`../examples/classify-charge-attack.md`](../examples/classify-charge-attack.md)
only when handling a multi-clause action description or checking citation
format.
