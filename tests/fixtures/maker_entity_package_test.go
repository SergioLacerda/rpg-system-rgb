package fixtures

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

type makerEntityPackageFixture struct {
	id                       string
	explicitFactFields       []string
	inferredSuggestionFields []string
	unresolvedFields         []string
}

var makerFieldListLine = regexp.MustCompile(`(?m)^    (\w+):\s*\[(.*)\]\s*$`)

// parseMakerEntityPackageFixture reads examples/maker-entity-package.yaml
// and extracts each package's three field-name lists. Line-oriented
// reader, matching this repo's other hand-authored fixture files.
func parseMakerEntityPackageFixture(t *testing.T) []makerEntityPackageFixture {
	t.Helper()
	body, err := os.ReadFile("examples/maker-entity-package.yaml")
	if err != nil {
		t.Fatal(err)
	}

	blocks := strings.Split(string(body), "\n  - id:")
	var packages []makerEntityPackageFixture
	for i, block := range blocks {
		if i == 0 {
			continue // header/comment content before the first package entry
		}
		full := "  id:" + block
		id := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(strings.SplitN(full, "\n", 2)[0]), "id:"))

		fields := map[string][]string{}
		for _, match := range makerFieldListLine.FindAllStringSubmatch(full, -1) {
			key, inner := match[1], strings.TrimSpace(match[2])
			var values []string
			if inner != "" {
				for _, part := range strings.Split(inner, ",") {
					values = append(values, strings.TrimSpace(part))
				}
			}
			fields[key] = values
		}

		packages = append(packages, makerEntityPackageFixture{
			id:                       id,
			explicitFactFields:       fields["explicit_fact_fields"],
			inferredSuggestionFields: fields["inferred_suggestion_fields"],
			unresolvedFields:         fields["unresolved_fields"],
		})
	}
	return packages
}

// mirrors: .analysis/pending/20260802-skill-contract-deterministic-validation-analysis.md
// TestMakerEntityPackageFieldsStaySeparate proves, for the fixture backing
// skills/maker/examples/harbor-quartermaster.md, that no field name
// appears in more than one of explicit_fact_fields,
// inferred_suggestion_fields, unresolved_fields — the invariant declared
// in skills/maker/schemas/entity-package.schema.yaml but never previously
// checked against any fixture data.
func TestMakerEntityPackageFieldsStaySeparate(t *testing.T) {
	packages := parseMakerEntityPackageFixture(t)
	if len(packages) == 0 {
		t.Fatal("no packages parsed from maker-entity-package.yaml")
	}

	for _, pkg := range packages {
		t.Run(pkg.id, func(t *testing.T) {
			if len(pkg.explicitFactFields) == 0 {
				t.Fatalf("%s: explicit_fact_fields must be non-empty (schema invariant: a package with zero explicit facts is pure invention)", pkg.id)
			}

			seen := map[string]string{}
			categories := []struct {
				name   string
				fields []string
			}{
				{"explicit_fact_fields", pkg.explicitFactFields},
				{"inferred_suggestion_fields", pkg.inferredSuggestionFields},
				{"unresolved_fields", pkg.unresolvedFields},
			}
			for _, category := range categories {
				for _, field := range category.fields {
					if priorCategory, ok := seen[field]; ok {
						t.Errorf("%s: field %q appears in both %s and %s — entity-package.schema.yaml requires each field in exactly one category", pkg.id, field, priorCategory, category.name)
						continue
					}
					seen[field] = category.name
				}
			}
		})
	}
}
