<!--
Template: Maker Report
Schema: ../schemas/maker-report.schema.yaml
One report per input batch (a GM's notes/images from one session or drop).
-->

# Maker Report: {{id}}

**Status:** {{status}}
**Input batch:** {{input_batch}}

## Entity Packages Produced

- {{entity-package id}} — {{one-line description}}

(If none: "No extractable entities found in this batch." — do not omit
this section.)

## Conflict Reports

- {{conflict-report id}} — {{one-line description}}

## Visual Observations

- {{visual-observation id}} — {{one-line description}}

## Open Questions For The GM

Cross-cutting questions that don't belong to a single entity package.

- {{question}}
