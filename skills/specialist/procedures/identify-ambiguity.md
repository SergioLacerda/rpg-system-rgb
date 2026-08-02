# Procedure: Identify Ambiguity

Name a question the canonical material genuinely does not resolve, instead
of picking a plausible-sounding answer.

## Input

A question that [`locate-authority.md`](locate-authority.md) or
[`compare-rules.md`](compare-rules.md) could not cleanly resolve to one
authority.

## When This Applies

- Zero units address the question directly, but a plausible answer could
  be inferred from adjacent rules (inference is not authority — say so).
- Two or more units plausibly apply and disagree, or apply in an order the
  index does not declare.
- A rule's `statement`/`formula`/`ordered_steps` covers the general case
  but the question asks about an edge case the text does not mention.

## Steps

1. State the question precisely.
2. State every unit consulted and why none of them resolves it (cite IDs
   even when citing "did not help" — this keeps the ambiguity report
   itself traceable).
3. If an inference is possible, offer it labeled clearly as an inference,
   never as a rule ("the text doesn't say, but by analogy with X, a
   reasonable ruling would be Y — this is a GM call, not a rule").
4. Do not silently pick the inference as "the answer."

## Output Shape

```text
Question: <question>
Status: AMBIGUOUS
Consulted: <ids> — <why each didn't resolve it>
Possible inference (non-authoritative): <text, or "none offered">
```

## Worked Example

See [`../examples/identify-ambiguity-sprint-reaction.md`](../examples/identify-ambiguity-sprint-reaction.md).

## Forbidden

- Reporting an inference as if it were a cited rule.
- Treating "no unit found" the same as "ambiguous" — see
  [`locate-authority.md`](locate-authority.md) step 5: a clean zero-match
  is "not covered," not automatically an ambiguity. This procedure is for
  cases with a real, describable tension — not simply an unaddressed gap.
