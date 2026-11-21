package fixtures

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/SergioLacerda/rpg-system-rgb/internal/components"
	"github.com/SergioLacerda/rpg-system-rgb/internal/components/maker"
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

			ep := fixtureToEntityPackage(pkg)
			if errs := maker.ValidateEntityPackage(ep); len(errs) != 0 {
				t.Errorf("%s: maker.ValidateEntityPackage rejected valid fixture: %v", pkg.id, errs)
			}

			dup := fixtureToEntityPackage(pkg)
			dup.InferredSuggestions = append(dup.InferredSuggestions, maker.InferredSuggestion{
				Field:      pkg.explicitFactFields[0],
				Value:      "duplicate",
				Confidence: "low",
				Rationale:  "parity check",
			})
			if errs := maker.ValidateEntityPackage(dup); len(errs) == 0 {
				t.Errorf("%s: maker.ValidateEntityPackage unexpectedly accepted duplicate field %q", pkg.id, pkg.explicitFactFields[0])
			}

			idx, err := components.LoadSemanticIndex("../../docs/core/semantic/core-v2.index.json")
			if err != nil {
				t.Fatalf("loading semantic index: %v", err)
			}
			withDeps := fixtureToEntityPackage(pkg)
			withDeps.CanonicalDeps = []string{"core.ability.contract"}
			if errs := maker.ValidateCanonicalDependencies(withDeps, idx); len(errs) != 0 {
				t.Errorf("%s: maker.ValidateCanonicalDependencies rejected valid canonical unit: %v", pkg.id, errs)
			}
			withDeps.CanonicalDeps = []string{"core.nonexistent.unit"}
			if errs := maker.ValidateCanonicalDependencies(withDeps, idx); len(errs) == 0 {
				t.Errorf("%s: maker.ValidateCanonicalDependencies accepted nonexistent unit", pkg.id)
			}
		})
	}
}

func fixtureToEntityPackage(pkg makerEntityPackageFixture) maker.EntityPackage {
	entityType := strings.SplitN(pkg.id, ".", 2)[0]
	ep := maker.EntityPackage{
		ID:         pkg.id,
		EntityType: entityType,
		Status:     "needs_review",
		SourceRefs: []string{"examples/maker-entity-package.yaml"},
	}
	for _, f := range pkg.explicitFactFields {
		ep.ExplicitFacts = append(ep.ExplicitFacts, maker.ExplicitFact{
			Field:     f,
			Value:     "fixture",
			SourceRef: "examples/maker-entity-package.yaml",
		})
	}
	for _, f := range pkg.inferredSuggestionFields {
		ep.InferredSuggestions = append(ep.InferredSuggestions, maker.InferredSuggestion{
			Field:      f,
			Value:      "fixture",
			Confidence: "medium",
			Rationale:  "fixture",
		})
	}
	for _, f := range pkg.unresolvedFields {
		ep.UnresolvedFields = append(ep.UnresolvedFields, maker.UnresolvedField{
			Field:  f,
			Reason: "fixture",
		})
	}
	return ep
}
