package tooling

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// sourceFixture bundles a valid semanticSourceFile and its supporting
// on-disk files (repo root markers, accepted ADR, EN/PT-br projection
// pages) so tests can mutate one field at a time and re-validate.
type sourceFixture struct {
	scratch    string
	sourceFile string
	indexFile  string
	source     semanticSourceFile
}

// newValidSourceFixture builds a minimal, fully valid semantic source
// scenario in a scratch directory: repo-root markers (README.md + docs/),
// an accepted ADR, an EN projection page, and a PT-br projection page
// containing the required translation term.
func newValidSourceFixture(t *testing.T) sourceFixture {
	t.Helper()
	scratch := t.TempDir()

	mustWriteFile(t, filepath.Join(scratch, "README.md"), "placeholder")
	mustMkdirAll(t, filepath.Join(scratch, "docs"))
	mustWriteFile(t, filepath.Join(scratch, "docs", "adr", "decision.md"), "## Status\n\nAccepted.\n")
	mustWriteFile(t, filepath.Join(scratch, "docs", "core", "en", "core", "topic.md"), "# Topic\n\nEnglish content.\n")
	mustWriteFile(t, filepath.Join(scratch, "docs", "core", "PT-br", "core", "topic.md"), "# Tópico\n\nConteúdo em português.\n")

	indexFile := filepath.Join(scratch, "docs", "core", "semantic", "index.json")
	mustWriteJSON(t, indexFile, sourceIndexForSourceValidation{
		Units: []struct {
			ID string `json:"id"`
		}{{ID: "U-001"}},
	})

	sourceFile := filepath.Join(scratch, "docs", "core", "semantic", "source", "core-v2-rules.v0.1.json")
	source := semanticSourceFile{
		Schema:                 expectedSemanticSourceSchema,
		SourceIndex:            "docs/core/semantic/index.json",
		AuthorityDecision:      "docs/adr/decision.md",
		SourceLocale:           "en",
		DefaultLocalizedLocale: "PT-br",
		Description:            "test description",
		Units: []semanticSourceUnit{
			{
				ID:            "U-001",
				Kind:          "rule",
				Locale:        "en",
				AuthorityType: "canonical_semantic",
				SourceStatus:  "canonical",
				Title:         "Test Unit",
				Statement:     "Test statement.",
				ProjectionPaths: map[string]string{
					"en": "docs/core/en/core/topic.md",
				},
				Relationships: map[string][]string{},
				Provenance: map[string]any{
					"source_revision": "abc123",
					"decision_refs":   []any{"docs/adr/decision.md"},
				},
				TranslationExpectations: map[string]translationExpectation{
					"pt_br": {
						ProjectionPath: "docs/core/PT-br/core/topic.md",
						RequiredTerms: []localizedTerm{
							{EN: "topic", PTBR: "Tópico"},
						},
					},
				},
				ComponentConsumers: []string{"specialist"},
			},
		},
	}

	fixture := sourceFixture{scratch: scratch, sourceFile: sourceFile, indexFile: indexFile, source: source}
	fixture.write(t)
	return fixture
}

// write (re-)serializes the fixture's source struct to sourceFile, so
// mutation tests can change a field and call write again before
// validating.
func (fixture sourceFixture) write(t *testing.T) {
	t.Helper()
	mustWriteJSON(t, fixture.sourceFile, fixture.source)
}

func (fixture sourceFixture) validate() error {
	return ValidateSemanticSource(fixture.sourceFile, fixture.indexFile)
}

func mustWriteFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func mustMkdirAll(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatal(err)
	}
}

func mustWriteJSON(t *testing.T, path string, value any) {
	t.Helper()
	bytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	mustWriteFile(t, path, string(bytes))
}

func TestValidateSemanticSourceAcceptsValidFixture(t *testing.T) {
	fixture := newValidSourceFixture(t)
	if err := fixture.validate(); err != nil {
		t.Fatalf("expected valid fixture to pass, got: %v", err)
	}
}

func TestValidateSemanticSourceLoadErrors(t *testing.T) {
	fixture := newValidSourceFixture(t)

	t.Run("missing source file", func(t *testing.T) {
		if err := ValidateSemanticSource(filepath.Join(fixture.scratch, "missing.json"), fixture.indexFile); err == nil {
			t.Fatal("expected error for missing source file")
		}
	})

	t.Run("missing index file", func(t *testing.T) {
		if err := ValidateSemanticSource(fixture.sourceFile, filepath.Join(fixture.scratch, "missing.json")); err == nil {
			t.Fatal("expected error for missing index file")
		}
	})

	t.Run("invalid source JSON", func(t *testing.T) {
		mustWriteFile(t, fixture.sourceFile, "not json")
		defer fixture.write(t)
		if err := fixture.validate(); err == nil {
			t.Fatal("expected error for invalid source JSON")
		}
	})

	t.Run("invalid index JSON", func(t *testing.T) {
		original, err := os.ReadFile(fixture.indexFile)
		if err != nil {
			t.Fatal(err)
		}
		mustWriteFile(t, fixture.indexFile, "not json")
		defer mustWriteFile(t, fixture.indexFile, string(original))
		if err := fixture.validate(); err == nil {
			t.Fatal("expected error for invalid index JSON")
		}
	})
}

func TestValidateSemanticSourceRejectsFileOutsideRepoRoot(t *testing.T) {
	// No README.md/docs/ marker anywhere under this directory, so
	// repoRootFromFile must fail before any field validation runs.
	orphan := t.TempDir()
	sourceFile := filepath.Join(orphan, "source.json")
	indexFile := filepath.Join(orphan, "index.json")
	mustWriteJSON(t, sourceFile, semanticSourceFile{})
	mustWriteJSON(t, indexFile, sourceIndexForSourceValidation{})

	if err := ValidateSemanticSource(sourceFile, indexFile); err == nil {
		t.Fatal("expected error when the source file has no resolvable repository root")
	}
}

func TestValidateSemanticSourceTopLevelFields(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*semanticSourceFile)
	}{
		{"wrong schema", func(s *semanticSourceFile) { s.Schema = "wrong/0.0" }},
		{"wrong source_index", func(s *semanticSourceFile) { s.SourceIndex = "docs/wrong/index.json" }},
		{"wrong source_locale", func(s *semanticSourceFile) { s.SourceLocale = "pt-br" }},
		{"wrong default_localized_locale", func(s *semanticSourceFile) { s.DefaultLocalizedLocale = "en" }},
		{"empty description", func(s *semanticSourceFile) { s.Description = "" }},
		{"empty units", func(s *semanticSourceFile) { s.Units = nil }},
		{"empty authority_decision", func(s *semanticSourceFile) { s.AuthorityDecision = "" }},
		{"absolute authority_decision", func(s *semanticSourceFile) { s.AuthorityDecision = "/docs/adr/decision.md" }},
		{"missing authority_decision file", func(s *semanticSourceFile) { s.AuthorityDecision = "docs/adr/missing.md" }},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			fixture := newValidSourceFixture(t)
			testCase.mutate(&fixture.source)
			fixture.write(t)
			if err := fixture.validate(); err == nil {
				t.Fatalf("expected error for %s", testCase.name)
			}
		})
	}
}

func TestValidateSemanticSourceAuthorityDecisionNotAccepted(t *testing.T) {
	fixture := newValidSourceFixture(t)
	mustWriteFile(t, filepath.Join(fixture.scratch, "docs", "adr", "draft.md"), "## Status\n\nDraft.\n")
	fixture.source.AuthorityDecision = "docs/adr/draft.md"
	fixture.write(t)
	if err := fixture.validate(); err == nil {
		t.Fatal("expected error when authority_decision ADR is not Accepted")
	}
}

func TestValidateSemanticSourceUnitIdentityFields(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*semanticSourceUnit)
	}{
		{"empty id", func(u *semanticSourceUnit) { u.ID = "" }},
		{"unit not in index", func(u *semanticSourceUnit) { u.ID = "U-999" }},
		{"empty kind", func(u *semanticSourceUnit) { u.Kind = "" }},
		{"wrong locale", func(u *semanticSourceUnit) { u.Locale = "pt-br" }},
		{"wrong authority_type", func(u *semanticSourceUnit) { u.AuthorityType = "derived" }},
		{"wrong source_status", func(u *semanticSourceUnit) { u.SourceStatus = "draft" }},
		{"empty title", func(u *semanticSourceUnit) { u.Title = "" }},
		{"empty statement", func(u *semanticSourceUnit) { u.Statement = "" }},
		{"nil relationships", func(u *semanticSourceUnit) { u.Relationships = nil }},
		{"empty provenance", func(u *semanticSourceUnit) { u.Provenance = nil }},
		{"empty component_consumers", func(u *semanticSourceUnit) { u.ComponentConsumers = nil }},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			fixture := newValidSourceFixture(t)
			testCase.mutate(&fixture.source.Units[0])
			fixture.write(t)
			if err := fixture.validate(); err == nil {
				t.Fatalf("expected error for %s", testCase.name)
			}
		})
	}
}

func TestValidateSemanticSourceDuplicateUnitID(t *testing.T) {
	fixture := newValidSourceFixture(t)
	mustWriteJSON(t, fixture.indexFile, sourceIndexForSourceValidation{
		Units: []struct {
			ID string `json:"id"`
		}{{ID: "U-001"}},
	})
	fixture.source.Units = append(fixture.source.Units, fixture.source.Units[0])
	fixture.write(t)
	if err := fixture.validate(); err == nil {
		t.Fatal("expected error for duplicate unit id")
	}
}

func TestValidateSemanticSourceProjectionPathFields(t *testing.T) {
	t.Run("empty projection_paths", func(t *testing.T) {
		fixture := newValidSourceFixture(t)
		fixture.source.Units[0].ProjectionPaths = map[string]string{}
		fixture.write(t)
		if err := fixture.validate(); err == nil {
			t.Fatal("expected error for empty projection_paths")
		}
	})

	t.Run("empty projection path value", func(t *testing.T) {
		fixture := newValidSourceFixture(t)
		fixture.source.Units[0].ProjectionPaths = map[string]string{"en": ""}
		fixture.write(t)
		if err := fixture.validate(); err == nil {
			t.Fatal("expected error for empty projection path value")
		}
	})

	t.Run("absolute projection path", func(t *testing.T) {
		fixture := newValidSourceFixture(t)
		fixture.source.Units[0].ProjectionPaths = map[string]string{"en": "/docs/core/en/core/topic.md"}
		fixture.write(t)
		if err := fixture.validate(); err == nil {
			t.Fatal("expected error for absolute projection path")
		}
	})

	t.Run("missing projection path file", func(t *testing.T) {
		fixture := newValidSourceFixture(t)
		fixture.source.Units[0].ProjectionPaths = map[string]string{"en": "docs/core/en/core/missing.md"}
		fixture.write(t)
		if err := fixture.validate(); err == nil {
			t.Fatal("expected error for missing projection path file")
		}
	})
}

func TestValidateSemanticSourceProvenanceFields(t *testing.T) {
	t.Run("missing source_revision", func(t *testing.T) {
		fixture := newValidSourceFixture(t)
		fixture.source.Units[0].Provenance = map[string]any{
			"decision_refs": []any{"docs/adr/decision.md"},
		}
		fixture.write(t)
		if err := fixture.validate(); err == nil {
			t.Fatal("expected error for missing provenance.source_revision")
		}
	})

	t.Run("empty source_revision", func(t *testing.T) {
		fixture := newValidSourceFixture(t)
		fixture.source.Units[0].Provenance["source_revision"] = ""
		fixture.write(t)
		if err := fixture.validate(); err == nil {
			t.Fatal("expected error for empty provenance.source_revision")
		}
	})

	t.Run("source_revision wrong type", func(t *testing.T) {
		fixture := newValidSourceFixture(t)
		fixture.source.Units[0].Provenance["source_revision"] = 123
		fixture.write(t)
		if err := fixture.validate(); err == nil {
			t.Fatal("expected error when provenance.source_revision is not a string")
		}
	})

	t.Run("missing decision_refs", func(t *testing.T) {
		fixture := newValidSourceFixture(t)
		fixture.source.Units[0].Provenance = map[string]any{
			"source_revision": "abc123",
		}
		fixture.write(t)
		if err := fixture.validate(); err == nil {
			t.Fatal("expected error for missing provenance.decision_refs")
		}
	})

	t.Run("decision_refs wrong type", func(t *testing.T) {
		fixture := newValidSourceFixture(t)
		fixture.source.Units[0].Provenance["decision_refs"] = "not-a-list"
		fixture.write(t)
		if err := fixture.validate(); err == nil {
			t.Fatal("expected error when provenance.decision_refs is not a list")
		}
	})

	t.Run("decision_refs empty list", func(t *testing.T) {
		fixture := newValidSourceFixture(t)
		fixture.source.Units[0].Provenance["decision_refs"] = []any{}
		fixture.write(t)
		if err := fixture.validate(); err == nil {
			t.Fatal("expected error when provenance.decision_refs is empty")
		}
	})

	t.Run("decision_refs entry wrong type", func(t *testing.T) {
		fixture := newValidSourceFixture(t)
		fixture.source.Units[0].Provenance["decision_refs"] = []any{123}
		fixture.write(t)
		if err := fixture.validate(); err == nil {
			t.Fatal("expected error when a decision_refs entry is not a string")
		}
	})

	t.Run("decision_refs entry missing file", func(t *testing.T) {
		fixture := newValidSourceFixture(t)
		fixture.source.Units[0].Provenance["decision_refs"] = []any{"docs/adr/missing.md"}
		fixture.write(t)
		if err := fixture.validate(); err == nil {
			t.Fatal("expected error when a decision_refs entry does not exist on disk")
		}
	})
}

func TestValidateSemanticSourceTranslationExpectations(t *testing.T) {
	t.Run("missing pt_br entry", func(t *testing.T) {
		fixture := newValidSourceFixture(t)
		fixture.source.Units[0].TranslationExpectations = map[string]translationExpectation{}
		fixture.write(t)
		if err := fixture.validate(); err == nil {
			t.Fatal("expected error for missing translation_expectations.pt_br")
		}
	})

	t.Run("empty projection_path", func(t *testing.T) {
		fixture := newValidSourceFixture(t)
		expectation := fixture.source.Units[0].TranslationExpectations["pt_br"]
		expectation.ProjectionPath = ""
		fixture.source.Units[0].TranslationExpectations["pt_br"] = expectation
		fixture.write(t)
		if err := fixture.validate(); err == nil {
			t.Fatal("expected error for empty translation_expectations.pt_br.projection_path")
		}
	})

	t.Run("missing PT-br projection file", func(t *testing.T) {
		fixture := newValidSourceFixture(t)
		expectation := fixture.source.Units[0].TranslationExpectations["pt_br"]
		expectation.ProjectionPath = "docs/core/PT-br/core/missing.md"
		fixture.source.Units[0].TranslationExpectations["pt_br"] = expectation
		fixture.write(t)
		if err := fixture.validate(); err == nil {
			t.Fatal("expected error when the PT-br projection file does not exist")
		}
	})

	t.Run("empty required_terms", func(t *testing.T) {
		fixture := newValidSourceFixture(t)
		expectation := fixture.source.Units[0].TranslationExpectations["pt_br"]
		expectation.RequiredTerms = nil
		fixture.source.Units[0].TranslationExpectations["pt_br"] = expectation
		fixture.write(t)
		if err := fixture.validate(); err == nil {
			t.Fatal("expected error for empty translation_expectations.pt_br.required_terms")
		}
	})

	t.Run("required term missing en", func(t *testing.T) {
		fixture := newValidSourceFixture(t)
		expectation := fixture.source.Units[0].TranslationExpectations["pt_br"]
		expectation.RequiredTerms = []localizedTerm{{EN: "", PTBR: "Tópico"}}
		fixture.source.Units[0].TranslationExpectations["pt_br"] = expectation
		fixture.write(t)
		if err := fixture.validate(); err == nil {
			t.Fatal("expected error when a required term is missing its en field")
		}
	})

	t.Run("required term missing pt_br", func(t *testing.T) {
		fixture := newValidSourceFixture(t)
		expectation := fixture.source.Units[0].TranslationExpectations["pt_br"]
		expectation.RequiredTerms = []localizedTerm{{EN: "topic", PTBR: ""}}
		fixture.source.Units[0].TranslationExpectations["pt_br"] = expectation
		fixture.write(t)
		if err := fixture.validate(); err == nil {
			t.Fatal("expected error when a required term is missing its pt_br field")
		}
	})

	t.Run("required term not present in content", func(t *testing.T) {
		fixture := newValidSourceFixture(t)
		expectation := fixture.source.Units[0].TranslationExpectations["pt_br"]
		expectation.RequiredTerms = []localizedTerm{{EN: "topic", PTBR: "termo-inexistente"}}
		fixture.source.Units[0].TranslationExpectations["pt_br"] = expectation
		fixture.write(t)
		if err := fixture.validate(); err == nil {
			t.Fatal("expected error when a required pt_br term is not present in the projection content")
		}
	})
}
