# Licensing

RGB System uses a split licensing model recorded in
[ADR-011](docs/adr/adr-011-split-license-model.md).

| Surface | License |
| --- | --- |
| Go source under `cmd/`, `internal/`, and `tests/` | MIT or Apache-2.0 |
| Shell scripts, CI workflows, configuration, and build tooling | MIT or Apache-2.0 |
| Web application code under `web/landing/src/` | MIT or Apache-2.0 |
| Rule text, examples, documentation, diagrams, web copy, generated Library/PDF content, and art | CC BY 4.0 |
| Generated semantic artifacts and bundles derived from rule/docs source | CC BY 4.0 |
| Dependency lockfiles and third-party package metadata | Their upstream licenses |

The root [LICENSE](LICENSE) file is an index for this split model. Full terms
are kept in:

- [LICENSE-MIT](LICENSE-MIT)
- [LICENSE-APACHE-2.0](LICENSE-APACHE-2.0)
- [LICENSE-CC-BY-4.0](LICENSE-CC-BY-4.0)

Unless a file states a different SPDX expression, contributions follow the
license of the surface they modify.
