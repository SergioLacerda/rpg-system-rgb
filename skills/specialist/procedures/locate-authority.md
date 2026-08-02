# Procedure: Locate Authority

Find which semantic unit(s) actually govern a topic, before answering.
This is the step that precedes [`explain.md`](explain.md) whenever the
right ID isn't already known — `explain.md` assumes resolution already
happened; this procedure is how it happens.

## Input

A topic, keyword, or question that hasn't yet been mapped to a semantic ID.

## Steps

1. Search `docs/core/semantic/core-v2.index.json` units by `title`,
   `index.tags`, and `index.retrieval_summary` for matches to the topic.
2. Prefer `generated/ai-context/core-specialist-pack.json` first (smaller,
   pre-scoped to Specialist) before falling back to the full index.
3. Rank candidates: `kind: rule` or `kind: procedure` outrank `kind:
   concept`, which outrank `kind: example` — the most authoritative unit
   for a mechanical question is a rule or procedure, not a concept
   definition or a worked instance.
4. If exactly one unit clearly has authority, that is the answer to
   "locate authority" — hand off to `explain.md` with that ID.
5. If zero units match, this is not automatically an ambiguity — it may
   simply mean the topic isn't covered yet. Report "no canonical unit
   found for `<topic>`" and stop; do not guess a nearby unit as a
   substitute.
6. If two or more units plausibly have authority and neither is clearly
   more specific, hand off to
   [`identify-ambiguity.md`](identify-ambiguity.md) instead of picking one.

## Output Shape

```text
Topic: <topic>
Authority: <semantic id> (<kind>) | "none found" | "ambiguous — see identify-ambiguity"
Source: <source_path>
```

## Worked Example

See [`../examples/locate-authority-shield-vs-armor.md`](../examples/locate-authority-shield-vs-armor.md).

## Forbidden

- Picking the first search hit without checking `kind` — a `core.example.*`
  hit is never the located authority when a `rule`/`procedure` unit also
  matches.
