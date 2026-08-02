# Procedure: Explain

Answer a question about a rule, resource, or concept in the RGB System.

## Input

A natural-language question, or a semantic ID directly (e.g.
`core.combat.attack-margin`).

## Steps

1. Resolve the question to one or more semantic unit IDs.
   - Prefer `generated/ai-context/core-specialist-pack.json` (fast path,
     already scoped to Specialist).
   - Fall back to `docs/core/semantic/core-v2.index.json` if the ID isn't in
     the pack.
2. Prefer `kind: rule` or `kind: procedure` units over `kind: example`
   units — an example illustrates a rule, it does not replace it.
3. If no unit matches, stop. Do not guess. Report:
   `"no canonical rule found for <question>"`.
4. Read the unit's `retrieval_summary` for the short answer, and its
   `projection_paths.markdown_en` (and `markdown_pt_br` if the user asked in
   Portuguese) for the full canonical text.
5. Compose the answer from the retrieval summary and Markdown source —
   never from memory of the rule, since the Markdown file is authoritative
   and may have changed.
6. Attach a citation: the semantic ID(s) used and their `source_path`.

## Output Shape

```text
Answer: <explanation, grounded in the cited unit(s)>
Cited: <semantic id> (<source_path>)
```

## Worked Example

See [`../examples/explain-attack-margin.md`](../examples/explain-attack-margin.md).

## Forbidden

- Answering from a `core.example.*` unit when a `rule`/`procedure` unit
  exists for the same question (explain the rule, then optionally
  illustrate with the example).
- Answering without a citation.
