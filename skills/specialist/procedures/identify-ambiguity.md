# Procedure: Identify Ambiguity

Name a question the canonical material genuinely does not resolve, rather
than picking a plausible-sounding answer.

## Input

A question that [`locate-authority.md`](locate-authority.md) or
[`compare-rules.md`](compare-rules.md) surfaced as unresolved across
candidate units (distinct from a clean zero-match topic, which stops in
[`locate-authority.md`](locate-authority.md)).

## When This Applies

- Zero units address an edge case directly, while adjacent rules invite an
  extrapolation.
- Two or more units plausibly apply and disagree, or interact in an order
  `relationships` does not declare.
- A rule's `statement`/`formula`/`ordered_steps` covers the general case
  but leaves the asked-about boundary unspecified.

## Steps

1. State the exact question and tension.
2. List **every** unit consulted and explain why none resolves it (cite IDs
   even when noting they did not resolve the case).
3. If a reasonable ruling can be inferred by analogy, present it strictly
   under `Possible inference (non-authoritative)` as a GM call, never as a
   cited rule (`split-on-tension`).

## Output Shape

```text
Question: <question>
Status: AMBIGUOUS
Consulted: <ids> — <why each didn't resolve it>
Possible inference (non-authoritative): <text, or "none offered">
```

## Worked Example

Read [`../examples/identify-ambiguity-sprint-reaction.md`](../examples/identify-ambiguity-sprint-reaction.md)
only when documenting edge cases where two action timings or rules overlap.
