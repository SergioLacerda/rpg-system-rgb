# RGB Specialist Skill

Status: Installable package, procedure runtime not implemented.
`SKILL.md` now carries front-matter (`name`/`description`) and an
[Installation](SKILL.md#installation) section, and the package is
self-contained — [`references/rgb-system.md`](references/rgb-system.md)
bundles a snapshot of the canonical rule documentation from `docs/core/**`
(see ADR-014), so the folder keeps working when copied outside this
monorepo. See [`SKILL.md`](SKILL.md) for the full contract (what Specialist
must and must not do), [`procedures/`](procedures/) for its nine core
procedures (explain, classify, validate, trace, locate-authority,
compare-rules, validate-example, identify-ambiguity, disclosure),
[`terminology/`](terminology/) for canonical EN/PT-br terms,
[`examples/`](examples/) for worked examples, and [`config.yaml`](config.yaml)
for the runtime contract summary.

This package does not yet run — it should not be presented as an available
product until a runtime implementation exists.

Once a runtime exists, distribution follows
[ADR-013: Skill Distribution Via The Existing Publication Pipeline](../../docs/adr/adr-013-skill-distribution-via-existing-publication-pipeline.md) —
a downloadable `.zip` published alongside the PDF/Library downloads, not a
separate GitHub Release or installer CLI.
