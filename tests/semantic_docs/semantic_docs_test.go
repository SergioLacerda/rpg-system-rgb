package semanticdocs

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type sourceFile struct {
	Schema string       `json:"schema"`
	Units  []sourceUnit `json:"units"`
}

type generatedFile struct {
	Schema        string `json:"schema"`
	ProjectionID  string `json:"projection_id"`
	AuthorityType string `json:"authority_type"`
	Units         []struct {
		ID string `json:"id"`
	} `json:"units"`
}

type l10nManifestFile struct {
	Schema            string         `json:"schema"`
	SourceLocale      string         `json:"source_locale"`
	LocalizedLocales  []string       `json:"localized_locales"`
	AuthorityDecision string         `json:"authority_decision"`
	Documents         []l10nDocument `json:"documents"`
}

type l10nDocument struct {
	Source            string `json:"source"`
	Locale            string `json:"locale"`
	Localized         string `json:"localized"`
	AuthorityType     string `json:"authority_type"`
	TranslationStatus string `json:"translation_status"`
}

type sourceUnit struct {
	ID                      string                            `json:"id"`
	AuthorityType           string                            `json:"authority_type"`
	SourceStatus            string                            `json:"source_status"`
	TranslationExpectations map[string]translationExpectation `json:"translation_expectations"`
}

type translationExpectation struct {
	ProjectionPath string          `json:"projection_path"`
	RequiredTerms  []localizedTerm `json:"required_terms"`
}

type localizedTerm struct {
	EN   string `json:"en"`
	PTBR string `json:"pt_br"`
}

func TestSemanticDocsValidationRunner(t *testing.T) {
	root := repoRoot(t)
	cmd := exec.Command("go", "run", "./cmd/rgb-tooling", "validate")
	cmd.Dir = root
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("semantic docs validation failed: %v\n%s", err, output)
	}

	text := string(output)
	for _, expected := range []string{
		"project-path validation passed",
		"semantic-index validation passed",
		"docs-l10n validation passed",
		"semantic-source validation passed",
		"semantic-contract validation passed",
		"semantic-projection validation passed",
		"generated-projection validation passed",
		"semantic docs validation passed",
	} {
		if !strings.Contains(text, expected) {
			t.Fatalf("validation output missing %q:\n%s", expected, text)
		}
	}
}

func TestAcceptedADRPromotesSemanticSource(t *testing.T) {
	root := repoRoot(t)
	adr := readText(t, filepath.Join(root, "docs/adr/adr-001-ai-first-documentation-authority.md"))
	if !strings.Contains(adr, "## Status\n\nAccepted.") {
		t.Fatal("ADR-001 must be Accepted before canonical semantic source promotion")
	}

	source := readSource(t, filepath.Join(root, "docs/core/semantic/source/core-v2-rules.v0.1.json"))
	if source.Schema != "rgb-docs-semantic-source/0.1" {
		t.Fatalf("unexpected semantic source schema: %s", source.Schema)
	}

	requiredIDs := map[string]bool{
		"core.resource.health":          false,
		"core.combat.attack-margin":     false,
		"core.damage.flow":              false,
		"core.damage.impact-source":     false,
		"core.damage.penetration":       false,
		"core.damage.armor-reduction":   false,
		"core.damage.shield-absorption": false,
		"core.ability.contract":         false,
	}

	for _, unit := range source.Units {
		if unit.AuthorityType != "canonical_semantic" {
			t.Fatalf("%s must be canonical_semantic, got %s", unit.ID, unit.AuthorityType)
		}
		if unit.SourceStatus != "canonical" {
			t.Fatalf("%s must be canonical, got %s", unit.ID, unit.SourceStatus)
		}
		if _, ok := requiredIDs[unit.ID]; ok {
			requiredIDs[unit.ID] = true
		}
	}

	for id, found := range requiredIDs {
		if !found {
			t.Fatalf("semantic source missing promoted unit %s", id)
		}
	}
}

func TestPTBRProjectionTermsMatchSemanticExpectations(t *testing.T) {
	root := repoRoot(t)
	source := readSource(t, filepath.Join(root, "docs/core/semantic/source/core-v2-rules.v0.1.json"))

	for _, unit := range source.Units {
		expectation, ok := unit.TranslationExpectations["pt_br"]
		if !ok {
			t.Fatalf("%s missing pt_br translation expectations", unit.ID)
		}
		content := readText(t, filepath.Join(root, filepath.FromSlash(expectation.ProjectionPath)))
		for _, term := range expectation.RequiredTerms {
			if term.EN == "" || term.PTBR == "" {
				t.Fatalf("%s has incomplete translation expectation: %#v", unit.ID, term)
			}
			if !strings.Contains(content, term.PTBR) {
				t.Fatalf("%s PT-br projection missing %q for %q", unit.ID, term.PTBR, term.EN)
			}
		}
	}
}

func TestDocsL10nManifestTracksCurrentPTBRDocuments(t *testing.T) {
	root := repoRoot(t)
	manifest := readL10nManifest(t, filepath.Join(root, "docs/core/semantic/l10n-manifest.v0.1.json"))
	if manifest.Schema != "rgb-docs-l10n/0.1" {
		t.Fatalf("unexpected L10n manifest schema: %s", manifest.Schema)
	}
	if manifest.SourceLocale != "en" {
		t.Fatalf("unexpected source locale: %s", manifest.SourceLocale)
	}
	if !containsString(manifest.LocalizedLocales, "PT-br") {
		t.Fatal("L10n manifest must include PT-br as a localized locale")
	}

	adr := readText(t, filepath.Join(root, filepath.FromSlash(manifest.AuthorityDecision)))
	if !strings.Contains(adr, "## Status\n\nAccepted.") {
		t.Fatal("L10n manifest authority decision must be Accepted")
	}

	for _, document := range manifest.Documents {
		if !strings.HasPrefix(document.Source, "docs/core/en/") {
			t.Fatalf("%s must be under docs/core/en", document.Source)
		}
		if document.Locale != "PT-br" {
			t.Fatalf("%s has unexpected locale %s", document.Source, document.Locale)
		}
		if document.TranslationStatus != "current" {
			t.Fatalf("%s must start as current, got %s", document.Source, document.TranslationStatus)
		}
		if document.AuthorityType == "" {
			t.Fatalf("%s missing authority_type", document.Source)
		}
		if !strings.HasPrefix(document.Localized, "docs/core/PT-br/") {
			t.Fatalf("%s localized path must be under docs/core/PT-br", document.Source)
		}
		if !exists(filepath.Join(root, filepath.FromSlash(document.Source))) {
			t.Fatalf("%s source path does not exist", document.Source)
		}
		if !exists(filepath.Join(root, filepath.FromSlash(document.Localized))) {
			t.Fatalf("%s localized path does not exist", document.Localized)
		}
	}
}

func TestGeneratedProjectionOutputsRemainDerivedArtifacts(t *testing.T) {
	root := repoRoot(t)
	outputs := []string{
		"generated/library/core-v2-rules.json",
		"generated/pdf/core-v2-rules.manifest.json",
		"generated/landing/core-v2-summary.json",
		"generated/bundles/core-v2-rules.bundle.json",
		"generated/ai-context/core-specialist-pack.json",
		"generated/search/core-v2.index.json",
	}

	for _, output := range outputs {
		bytes, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(output)))
		if err != nil {
			t.Fatal(err)
		}
		var generated generatedFile
		if err := json.Unmarshal(bytes, &generated); err != nil {
			t.Fatal(err)
		}
		if generated.Schema != "rgb-docs-generated-projection/0.1" {
			t.Fatalf("%s has unexpected schema %s", output, generated.Schema)
		}
		if generated.AuthorityType != "generated_artifact" {
			t.Fatalf("%s must remain generated_artifact, got %s", output, generated.AuthorityType)
		}
		if generated.ProjectionID == "" {
			t.Fatalf("%s missing projection_id", output)
		}
		if len(generated.Units) == 0 {
			t.Fatalf("%s must include projected units", output)
		}
	}
}

func readSource(t *testing.T, path string) sourceFile {
	t.Helper()
	bytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var source sourceFile
	if err := json.Unmarshal(bytes, &source); err != nil {
		t.Fatal(err)
	}
	return source
}

func readL10nManifest(t *testing.T, path string) l10nManifestFile {
	t.Helper()
	bytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var manifest l10nManifestFile
	if err := json.Unmarshal(bytes, &manifest); err != nil {
		t.Fatal(err)
	}
	return manifest
}

func readText(t *testing.T, path string) string {
	t.Helper()
	bytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(bytes)
}

func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if exists(filepath.Join(dir, "docs/core/semantic/core-v2.index.json")) {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not find repository root")
		}
		dir = parent
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
