package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type SemanticIndex struct {
	Schema                 string   `json:"schema"`
	SourceLocale           string   `json:"source_locale"`
	DefaultLocalizedLocale string   `json:"default_localized_locale"`
	AuthorityTypes         []string `json:"authority_types"`
	SourceStatuses         []string `json:"source_statuses"`
	Kinds                  []string `json:"kinds"`
	ProjectionSurfaces     []string `json:"projection_surfaces"`
	ComponentConsumers     []string `json:"component_consumers"`
	Units                  []Unit   `json:"units"`
}

type Unit struct {
	ID                 string              `json:"id"`
	Kind               string              `json:"kind"`
	Locale             string              `json:"locale"`
	AuthorityType      string              `json:"authority_type"`
	SourceStatus       string              `json:"source_status"`
	Title              string              `json:"title"`
	SourcePath         string              `json:"source_path"`
	ProjectionPaths    map[string]string   `json:"projection_paths"`
	Relationships      map[string][]string `json:"relationships"`
	Index              UnitIndex           `json:"index"`
	Provenance         map[string]any      `json:"provenance"`
	ComponentConsumers []string            `json:"component_consumers"`
	SourceUnit         string              `json:"source_unit,omitempty"`
	TranslationStatus  string              `json:"translation_status,omitempty"`
}

const expectedSchema = "rgb-docs-semantic-index/0.1"

type UnitIndex struct {
	Track            []string `json:"track"`
	Tags             []string `json:"tags"`
	RetrievalSummary string   `json:"retrieval_summary"`
}

var relationshipFieldsWithUnitIDs = []string{
	"depends_on",
	"clarifies",
	"implements",
	"illustrated_by",
	"translated_by",
	"supersedes",
}

func main() {
	if len(os.Args) != 2 {
		fatal(errors.New("usage: go run scripts/validate_semantic_index.go <index.json>"))
	}
	if err := validate(os.Args[1]); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "semantic-index validation failed: %v\n", err)
	os.Exit(1)
}

func validate(indexFile string) error {
	bytes, err := os.ReadFile(indexFile)
	if err != nil {
		return err
	}

	var index SemanticIndex
	if err := json.Unmarshal(bytes, &index); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	if index.Schema != expectedSchema {
		return fmt.Errorf("schema must be `%s`", expectedSchema)
	}
	if index.SourceLocale == "" {
		return errors.New("missing top-level field `source_locale`")
	}
	if index.DefaultLocalizedLocale == "" {
		return errors.New("missing top-level field `default_localized_locale`")
	}
	if len(index.Units) == 0 {
		return errors.New("missing or empty top-level field `units`")
	}

	authorityTypes, err := requiredSet("authority_types", index.AuthorityTypes)
	if err != nil {
		return err
	}
	sourceStatuses, err := requiredSet("source_statuses", index.SourceStatuses)
	if err != nil {
		return err
	}
	kinds, err := requiredSet("kinds", index.Kinds)
	if err != nil {
		return err
	}
	projectionSurfaces, err := requiredSet("projection_surfaces", index.ProjectionSurfaces)
	if err != nil {
		return err
	}
	componentConsumers, err := requiredSet("component_consumers", index.ComponentConsumers)
	if err != nil {
		return err
	}

	repoRoot, err := repoRootFromIndex(indexFile)
	if err != nil {
		return err
	}

	unitIDs := make(map[string]bool)
	for _, unit := range index.Units {
		if err := validateUnit(repoRoot, unit, authorityTypes, sourceStatuses, kinds, projectionSurfaces, componentConsumers, unitIDs); err != nil {
			return err
		}
		unitIDs[unit.ID] = true
	}

	for _, unit := range index.Units {
		if err := validateRelationships(unit, unitIDs); err != nil {
			return err
		}
	}

	fmt.Printf("semantic-index validation passed: %s (%d units)\n", indexFile, len(index.Units))
	return nil
}

func repoRootFromIndex(indexFile string) (string, error) {
	abs, err := filepath.Abs(indexFile)
	if err != nil {
		return "", err
	}
	return findRepoRoot(filepath.Dir(abs))
}

func findRepoRoot(start string) (string, error) {
	for dir := start; ; dir = filepath.Dir(dir) {
		if fileExists(filepath.Join(dir, "README.md")) && dirExists(filepath.Join(dir, "docs")) {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find repository root from %s", start)
		}
	}
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func validateUnit(
	repoRoot string,
	unit Unit,
	authorityTypes map[string]bool,
	sourceStatuses map[string]bool,
	kinds map[string]bool,
	projectionSurfaces map[string]bool,
	componentConsumers map[string]bool,
	seenIDs map[string]bool,
) error {
	if unit.ID == "" {
		return errors.New("unit id must be a non-empty string")
	}
	if seenIDs[unit.ID] {
		return fmt.Errorf("duplicate unit id: %s", unit.ID)
	}
	if unit.Kind == "" || !kinds[unit.Kind] {
		return fmt.Errorf("%s: unknown kind `%s`", unit.ID, unit.Kind)
	}
	if unit.Locale == "" {
		return fmt.Errorf("%s: locale must be non-empty", unit.ID)
	}
	if unit.AuthorityType == "" || !authorityTypes[unit.AuthorityType] {
		return fmt.Errorf("%s: unknown authority_type `%s`", unit.ID, unit.AuthorityType)
	}
	if unit.SourceStatus == "" || !sourceStatuses[unit.SourceStatus] {
		return fmt.Errorf("%s: unknown source_status `%s`", unit.ID, unit.SourceStatus)
	}
	if unit.Title == "" {
		return fmt.Errorf("%s: title must be non-empty", unit.ID)
	}
	if err := validatePath(repoRoot, unit.SourcePath, unit.ID, "source_path"); err != nil {
		return err
	}
	if len(unit.ProjectionPaths) == 0 {
		return fmt.Errorf("%s: projection_paths must be non-empty", unit.ID)
	}
	for surface, path := range unit.ProjectionPaths {
		if !projectionSurfaces[surface] {
			return fmt.Errorf("%s: unknown projection surface `%s`", unit.ID, surface)
		}
		if err := validatePath(repoRoot, path, unit.ID, "projection_paths."+surface); err != nil {
			return err
		}
	}
	if unit.AuthorityType == "canonical_markdown_bridge" && unit.ProjectionPaths["markdown_en"] == "" {
		return fmt.Errorf("%s: canonical Markdown bridge units must expose projection_paths.markdown_en", unit.ID)
	}
	if unit.Relationships == nil {
		return fmt.Errorf("%s: relationships must be an object", unit.ID)
	}
	if unit.Index.RetrievalSummary == "" {
		return fmt.Errorf("%s: index.retrieval_summary must be non-empty", unit.ID)
	}
	if len(unit.Index.Track) == 0 {
		return fmt.Errorf("%s: index.track must be non-empty", unit.ID)
	}
	for _, consumer := range unit.Index.Track {
		if !componentConsumers[consumer] {
			return fmt.Errorf("%s: unknown index track consumer `%s`", unit.ID, consumer)
		}
	}
	if len(unit.Index.Tags) == 0 {
		return fmt.Errorf("%s: index.tags must be non-empty", unit.ID)
	}
	if len(unit.Provenance) == 0 {
		return fmt.Errorf("%s: provenance must be non-empty", unit.ID)
	}
	if err := validateProvenance(repoRoot, unit); err != nil {
		return err
	}
	if len(unit.ComponentConsumers) == 0 {
		return fmt.Errorf("%s: component_consumers must be non-empty", unit.ID)
	}
	for _, consumer := range unit.ComponentConsumers {
		if !componentConsumers[consumer] {
			return fmt.Errorf("%s: unknown component consumer `%s`", unit.ID, consumer)
		}
	}
	if err := validateTranslationContract(unit, sourceStatuses); err != nil {
		return err
	}
	if unit.AuthorityType == "generated_artifact" && unit.SourceStatus == "canonical" {
		return fmt.Errorf("%s: generated artifacts must not use canonical source_status", unit.ID)
	}
	return nil
}

func validateProvenance(repoRoot string, unit Unit) error {
	sourceRevision, ok := unit.Provenance["source_revision"].(string)
	if !ok || sourceRevision == "" {
		return fmt.Errorf("%s: provenance.source_revision must be a non-empty string", unit.ID)
	}

	refsValue, exists := unit.Provenance["decision_refs"]
	if !exists {
		return nil
	}
	refs, ok := refsValue.([]any)
	if !ok {
		return fmt.Errorf("%s: provenance.decision_refs must be a string list", unit.ID)
	}
	for index, refValue := range refs {
		ref, ok := refValue.(string)
		if !ok {
			return fmt.Errorf("%s: provenance.decision_refs[%d] must be a string", unit.ID, index)
		}
		if err := validatePath(repoRoot, ref, unit.ID, fmt.Sprintf("provenance.decision_refs[%d]", index)); err != nil {
			return err
		}
	}
	return nil
}

func validateTranslationContract(unit Unit, sourceStatuses map[string]bool) error {
	if unit.Kind != "translation" && unit.AuthorityType != "translation" {
		return nil
	}
	if unit.SourceUnit == "" {
		return fmt.Errorf("%s: translation units must declare source_unit", unit.ID)
	}
	if unit.TranslationStatus == "" || !sourceStatuses[unit.TranslationStatus] {
		return fmt.Errorf("%s: unknown translation_status `%s`", unit.ID, unit.TranslationStatus)
	}
	if unit.ProjectionPaths["markdown_pt_br"] == "" {
		return fmt.Errorf("%s: translation units must expose projection_paths.markdown_pt_br", unit.ID)
	}
	return nil
}

func validateRelationships(unit Unit, unitIDs map[string]bool) error {
	for _, field := range relationshipFieldsWithUnitIDs {
		for _, targetID := range unit.Relationships[field] {
			if !unitIDs[targetID] {
				return fmt.Errorf("%s: relationship `%s` points to unknown unit `%s`", unit.ID, field, targetID)
			}
		}
	}
	if unit.SourceUnit != "" && !unitIDs[unit.SourceUnit] {
		return fmt.Errorf("%s: source_unit points to unknown unit `%s`", unit.ID, unit.SourceUnit)
	}
	return nil
}

func validatePath(repoRoot, pathValue, unitID, field string) error {
	if pathValue == "" {
		return fmt.Errorf("%s: %s must be non-empty", unitID, field)
	}
	if filepath.IsAbs(pathValue) {
		return fmt.Errorf("%s: %s must be repository-relative", unitID, field)
	}
	if pathValue == ".git" || pathValue == ".git/" {
		return fmt.Errorf("%s: %s must not point at git internals", unitID, field)
	}
	fullPath := filepath.Join(repoRoot, filepath.FromSlash(pathValue))
	if _, err := os.Stat(fullPath); err != nil {
		return fmt.Errorf("%s: %s does not exist: %s", unitID, field, pathValue)
	}
	return nil
}

func requiredSet(field string, values []string) (map[string]bool, error) {
	if len(values) == 0 {
		return nil, fmt.Errorf("missing or empty top-level field `%s`", field)
	}
	result := make(map[string]bool, len(values))
	for _, value := range values {
		if value == "" {
			return nil, fmt.Errorf("top-level field `%s` must not contain empty values", field)
		}
		result[value] = true
	}
	return result, nil
}
