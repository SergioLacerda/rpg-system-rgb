package tooling

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const expectedSemanticSourceSchema = "rgb-docs-semantic-source/0.1"

type semanticSourceFile struct {
	Schema                 string               `json:"schema"`
	SourceIndex            string               `json:"source_index"`
	AuthorityDecision      string               `json:"authority_decision"`
	SourceLocale           string               `json:"source_locale"`
	DefaultLocalizedLocale string               `json:"default_localized_locale"`
	Description            string               `json:"description"`
	Units                  []semanticSourceUnit `json:"units"`
}

type semanticSourceUnit struct {
	ID                      string                            `json:"id"`
	Kind                    string                            `json:"kind"`
	Locale                  string                            `json:"locale"`
	AuthorityType           string                            `json:"authority_type"`
	SourceStatus            string                            `json:"source_status"`
	Title                   string                            `json:"title"`
	Statement               string                            `json:"statement"`
	ProjectionPaths         map[string]string                 `json:"projection_paths"`
	Relationships           map[string][]string               `json:"relationships"`
	Provenance              map[string]any                    `json:"provenance"`
	TranslationExpectations map[string]translationExpectation `json:"translation_expectations"`
	ComponentConsumers      []string                          `json:"component_consumers"`
}

type translationExpectation struct {
	ProjectionPath string          `json:"projection_path"`
	RequiredTerms  []localizedTerm `json:"required_terms"`
}

type localizedTerm struct {
	EN   string `json:"en"`
	PTBR string `json:"pt_br"`
}

type sourceIndexForSourceValidation struct {
	Units []struct {
		ID string `json:"id"`
	} `json:"units"`
}

// ValidateSemanticSource validates a promoted semantic source file (e.g.
// docs/core/semantic/source/core-v2-rules.v0.1.json) against its parent
// index, migrated from scripts/validate_semantic_source.go.
func ValidateSemanticSource(sourceFile, indexFile string) error {
	sourceBytes, err := os.ReadFile(sourceFile) //nolint:gosec // G304: path is a caller-supplied validation target, by design
	if err != nil {
		return err
	}
	indexBytes, err := os.ReadFile(indexFile) //nolint:gosec // G304: path is a caller-supplied validation target, by design
	if err != nil {
		return err
	}

	var source semanticSourceFile
	if err := json.Unmarshal(sourceBytes, &source); err != nil {
		return fmt.Errorf("invalid source JSON: %w", err)
	}
	var index sourceIndexForSourceValidation
	if err := json.Unmarshal(indexBytes, &index); err != nil {
		return fmt.Errorf("invalid index JSON: %w", err)
	}

	repoRoot, err := repoRootFromFile(sourceFile)
	if err != nil {
		return err
	}
	if err := validateSemanticSourceTopLevel(source, repoRoot, indexFile); err != nil {
		return err
	}

	indexIDs := map[string]bool{}
	for _, unit := range index.Units {
		indexIDs[unit.ID] = true
	}

	seenIDs := map[string]bool{}
	for _, unit := range source.Units {
		if err := validateSemanticSourceUnit(repoRoot, unit, indexIDs, seenIDs); err != nil {
			return err
		}
		seenIDs[unit.ID] = true
	}

	fmt.Printf("semantic-source validation passed: %s (%d units)\n", sourceFile, len(source.Units))
	return nil
}

func validateSemanticSourceTopLevel(source semanticSourceFile, repoRoot, indexFile string) error {
	if source.Schema != expectedSemanticSourceSchema {
		return fmt.Errorf("schema must be `%s`", expectedSemanticSourceSchema)
	}
	if source.SourceIndex != repoRelative(repoRoot, indexFile) {
		return fmt.Errorf("source_index must match validation index path `%s`", indexFile)
	}
	if source.SourceLocale != "en" {
		return errors.New("source_locale must be `en`")
	}
	if source.DefaultLocalizedLocale != "PT-br" {
		return errors.New("default_localized_locale must be `PT-br`")
	}
	if source.Description == "" {
		return errors.New("description must be non-empty")
	}
	if len(source.Units) == 0 {
		return errors.New("units must be non-empty")
	}
	return validateAcceptedDecision(repoRoot, source.AuthorityDecision)
}

// validateAcceptedDecision confirms a referenced ADR exists and has been
// accepted. Shared by the semantic-source and l10n-manifest validators
// (previously two near-identical copies in scripts/validate_semantic_source.go
// and scripts/validate_docs_l10n_manifest.go).
func validateAcceptedDecision(repoRoot, pathValue string) error {
	if pathValue == "" {
		return errors.New("authority_decision must be non-empty")
	}
	if filepath.IsAbs(pathValue) {
		return errors.New("authority_decision must be repository-relative")
	}
	decisionPath := filepath.Join(repoRoot, filepath.FromSlash(pathValue))
	bytes, err := os.ReadFile(decisionPath) //nolint:gosec // G304: path is repo-root-joined and validated repository-relative above
	if err != nil {
		return fmt.Errorf("authority_decision does not exist: %s", pathValue)
	}
	if !strings.Contains(string(bytes), "## Status\n\nAccepted.") {
		return fmt.Errorf("authority_decision must have status Accepted: %s", pathValue)
	}
	return nil
}

func validateSemanticSourceUnit(repoRoot string, unit semanticSourceUnit, indexIDs, seenIDs map[string]bool) error {
	if err := validateSemanticSourceUnitIdentity(unit, indexIDs, seenIDs); err != nil {
		return err
	}
	if len(unit.ProjectionPaths) == 0 {
		return fmt.Errorf("%s: projection_paths must be non-empty", unit.ID)
	}
	for surface, pathValue := range unit.ProjectionPaths {
		if err := validateExistingSemanticSourcePath(repoRoot, unit.ID, "projection_paths."+surface, pathValue); err != nil {
			return err
		}
	}
	if unit.Relationships == nil {
		return fmt.Errorf("%s: relationships must be an object", unit.ID)
	}
	if len(unit.Provenance) == 0 {
		return fmt.Errorf("%s: provenance must be non-empty", unit.ID)
	}
	if err := validateSemanticSourceProvenance(repoRoot, unit); err != nil {
		return err
	}
	if len(unit.ComponentConsumers) == 0 {
		return fmt.Errorf("%s: component_consumers must be non-empty", unit.ID)
	}
	return validateTranslationExpectations(repoRoot, unit)
}

func validateSemanticSourceUnitIdentity(unit semanticSourceUnit, indexIDs, seenIDs map[string]bool) error {
	if unit.ID == "" {
		return errors.New("unit id must be non-empty")
	}
	if seenIDs[unit.ID] {
		return fmt.Errorf("duplicate unit id: %s", unit.ID)
	}
	if !indexIDs[unit.ID] {
		return fmt.Errorf("%s: unit must exist in source index", unit.ID)
	}
	if unit.Kind == "" {
		return fmt.Errorf("%s: kind must be non-empty", unit.ID)
	}
	if unit.Locale != "en" {
		return fmt.Errorf("%s: locale must be `en`", unit.ID)
	}
	if unit.AuthorityType != "canonical_semantic" {
		return fmt.Errorf("%s: authority_type must be canonical_semantic", unit.ID)
	}
	if unit.SourceStatus != "canonical" {
		return fmt.Errorf("%s: source_status must be canonical", unit.ID)
	}
	if unit.Title == "" {
		return fmt.Errorf("%s: title must be non-empty", unit.ID)
	}
	if unit.Statement == "" {
		return fmt.Errorf("%s: statement must be non-empty", unit.ID)
	}
	return nil
}

func validateSemanticSourceProvenance(repoRoot string, unit semanticSourceUnit) error {
	sourceRevision, ok := unit.Provenance["source_revision"].(string)
	if !ok || sourceRevision == "" {
		return fmt.Errorf("%s: provenance.source_revision must be a non-empty string", unit.ID)
	}
	refsValue, exists := unit.Provenance["decision_refs"]
	if !exists {
		return fmt.Errorf("%s: provenance.decision_refs must be present", unit.ID)
	}
	refs, ok := refsValue.([]any)
	if !ok || len(refs) == 0 {
		return fmt.Errorf("%s: provenance.decision_refs must be a non-empty string list", unit.ID)
	}
	for index, refValue := range refs {
		ref, ok := refValue.(string)
		if !ok {
			return fmt.Errorf("%s: provenance.decision_refs[%d] must be a string", unit.ID, index)
		}
		if err := validateExistingSemanticSourcePath(repoRoot, unit.ID, fmt.Sprintf("provenance.decision_refs[%d]", index), ref); err != nil {
			return err
		}
	}
	return nil
}

func validateTranslationExpectations(repoRoot string, unit semanticSourceUnit) error {
	expectation, ok := unit.TranslationExpectations["pt_br"]
	if !ok {
		return fmt.Errorf("%s: translation_expectations.pt_br must be present", unit.ID)
	}
	if expectation.ProjectionPath == "" {
		return fmt.Errorf("%s: translation_expectations.pt_br.projection_path must be non-empty", unit.ID)
	}
	fullPath := filepath.Join(repoRoot, filepath.FromSlash(expectation.ProjectionPath))
	bytes, err := os.ReadFile(fullPath) //nolint:gosec // G304: path is repo-root-joined from a declared projection path, by design
	if err != nil {
		return fmt.Errorf("%s: PT-br projection does not exist: %s", unit.ID, expectation.ProjectionPath)
	}
	content := string(bytes)
	if len(expectation.RequiredTerms) == 0 {
		return fmt.Errorf("%s: translation_expectations.pt_br.required_terms must be non-empty", unit.ID)
	}
	for index, term := range expectation.RequiredTerms {
		if term.EN == "" || term.PTBR == "" {
			return fmt.Errorf("%s: translation term %d must include en and pt_br", unit.ID, index)
		}
		if !strings.Contains(content, term.PTBR) {
			return fmt.Errorf("%s: PT-br projection missing expected term `%s`", unit.ID, term.PTBR)
		}
	}
	return nil
}

func validateExistingSemanticSourcePath(repoRoot, unitID, field, pathValue string) error {
	if pathValue == "" {
		return fmt.Errorf("%s: %s must be non-empty", unitID, field)
	}
	if filepath.IsAbs(pathValue) {
		return fmt.Errorf("%s: %s must be repository-relative", unitID, field)
	}
	fullPath := filepath.Join(repoRoot, filepath.FromSlash(pathValue))
	if _, err := os.Stat(fullPath); err != nil {
		return fmt.Errorf("%s: %s does not exist: %s", unitID, field, pathValue)
	}
	return nil
}
