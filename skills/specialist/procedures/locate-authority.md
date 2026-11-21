# Procedure: Locate Authority

Find which semantic unit(s) govern a topic before answering.

## Input

A topic, keyword, or question not yet mapped to a semantic ID.

## Steps

1. Search `generated/ai-context/core-specialist-pack.json` first (fast,
   pre-scoped to Specialist), falling back to
   `docs/core/semantic/core-v2.index.json` by `title`, `index.tags`, and
   `index.retrieval_summary`.
2. Rank matching candidates per `rule-over-example`: `kind: rule` and
   `kind: procedure` outrank `kind: concept`, which outrank `kind: example`.
   Never settle on the first search hit or a `core.example.*` unit when a
   `rule` or `procedure` unit also matches.
3. If exactly one unit clearly holds authority, hand off to
   [`explain.md`](explain.md) with that semantic ID.
4. If zero units match, report `"no canonical unit found for <topic>"` and
   stop — a clean zero-match is an uncovered topic, not an ambiguity or an
   invitation to substitute a nearby unit.
5. If two or more units plausibly hold authority and neither is clearly more
   specific, hand off to [`identify-ambiguity.md`](identify-ambiguity.md).

## Output Shape

```text
Topic: <topic>
Authority: <semantic id> (<kind>) | "none found" | "ambiguous — see identify-ambiguity"
Source: <source_path>
```

## Worked Example

Read [`../examples/locate-authority-shield-vs-armor.md`](../examples/locate-authority-shield-vs-armor.md)
only when disambiguating two adjacent defensive or mechanical units.
