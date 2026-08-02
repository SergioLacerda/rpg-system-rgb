# Prompt: Detect Conflicts

Use whenever two or more sources (two notes, a note and an image, two
images) describe the same field of the same entity differently.

## Instructions

1. Identify the disputed field precisely — a conflict report covers
   exactly one field. "The quartermaster's name" and "the quartermaster's
   loyalty" are two separate conflict reports even if they come from the
   same pair of sources.
2. Record every competing claim with its source and `source_kind`
   (`text` or `image`) in a [`conflict-report`](../templates/conflict-report.md).
3. Apply the authority rule from `../SKILL.md` (`image = appearance, text
   = facts`) to decide resolution:
   - If one claim's `source_kind` is `image` and the disputed field is
     **not** an appearance property (name, loyalty, backstory, canon
     status, relationships, stats) — the image was never authoritative
     over that field. Set `resolution: resolved_by_authority_rule` and
     explain in `resolution_note` why the image claim doesn't apply.
   - If both claims are `text`, or the field genuinely is an appearance
     property both sources address, the rule does not resolve it — leave
     `resolution: unresolved` for the user.
4. Never default to "most recent source wins" or "most detailed source
   wins" — those are not the authority rule, and using them silently
   would misattribute a Maker guess as a resolved fact.
5. Link the conflict report from every entity-package it affects.

## Do Not

- Resolve a conflict between two text sources on Maker's own judgment —
  that decision belongs to the user.
- Apply `resolved_by_authority_rule` to a field the image genuinely could
  speak to (e.g. "what does the merchant look like" — if the dispute is
  image vs. image, the authority rule doesn't apply either, since neither
  side is text).
