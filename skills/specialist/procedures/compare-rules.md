# Procedure: Compare Rules

Explain how two rules, procedures, or concepts differ or interact.

## Input

Two semantic IDs (or two questions that [`locate-authority.md`](locate-authority.md)
resolves to two IDs).

## Steps

1. Locate both units in `docs/core/semantic/core-v2.index.json` (and, for
   `rule`/`procedure`/`ability` kinds, `docs/core/semantic/source/core-v2-rules.v0.1.json`
   for the full `statement`/`formula`/`ordered_steps`).
2. Report each unit's `kind` and one-line summary side by side — do not
   blend them into a single paragraph that loses which fact came from
   which unit.
3. Check `relationships`: if one `depends_on` or `clarifies` the other,
   say so explicitly — that's usually the actual answer to "how do these
   interact" (e.g. `core.damage.armor-reduction` and
   `core.damage.shield-absorption` are both steps `core.damage.flow`
   `depends_on`, in a fixed order — the comparison is "sequence," not
   "alternative").
4. If the two units are siblings with no declared relationship, say that
   too — do not invent a relationship the index doesn't declare.
5. If the comparison reveals the two rules could plausibly both apply to
   the same situation in conflicting ways, hand off to
   [`identify-ambiguity.md`](identify-ambiguity.md).

## Output Shape

```text
A: <id> (<kind>) — <one-line summary>
B: <id> (<kind>) — <one-line summary>
Relationship: <depends_on | clarifies | sequence | "no declared relationship">
Cited: <both ids and their source_paths>
```

## Worked Example

See [`../examples/compare-armor-vs-shield.md`](../examples/compare-armor-vs-shield.md).

## Forbidden

- Declaring a relationship between two units that their `relationships`
  fields do not declare.
