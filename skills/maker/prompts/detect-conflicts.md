# Prompt: Detect Conflicts

Record and classify single-field disagreements between two or more sources
(notes or images) for the same entity.

## Steps

1. Isolate each disputed field into its own
   [`conflict-report`](../templates/conflict-report.md); route whole-entity
   reboots or multi-field drafts to `variants` per `variant-preserved` in
   [`../SKILL.md`](../SKILL.md).
2. Record every competing claim with `value`, `source_ref`, and `source_kind`
   (`text` or `image`).
3. Set `resolution` per `text-over-image` and `conflict-on-dispute`:
   - Mixed `text` + `image` on a **non-appearance** fact → set
     `resolution: resolved_by_authority_rule` and document why text wins in
     `resolution_note`.
   - Text vs. text, image vs. image, or appearance disputes → keep
     `resolution: unresolved` for the user regardless of recency or detail level.
4. Link the conflict report ID from every affected `entity-package`.
5. After resolution (`resolved_by_user` or `resolved_by_authority_rule`),
   archive the report under `packages/conflicts/resolved/` and keep only
   `unresolved` reports in `packages/conflicts/`.

## Completion Check

- Every conflict report has at least two claims.
- `resolution_note` is present whenever `resolution` is
  `resolved_by_authority_rule` or `resolved_by_user`.
- Zero resolved reports remain in `packages/conflicts/`; all resolved reports
  live in `packages/conflicts/resolved/`.

## Worked Example

Read [`../examples/harbor-quartermaster.md`](../examples/harbor-quartermaster.md)
only when verifying how `conflict-report` interacts with an entity's
`unresolved_fields`.
