package tooling

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const expectedSemanticIndexSchema = "rgb-docs-semantic-index/0.1"

// SemanticIndex is the decoded shape of docs/core/semantic/core-v2.index.json.
type SemanticIndex struct {
	Schema                 string              `json:"schema"`
	SourceLocale           string              `json:"source_locale"`
	DefaultLocalizedLocale string              `json:"default_localized_locale"`
	AuthorityTypes         []string            `json:"authority_types"`
	SourceStatuses         []string            `json:"source_statuses"`
	Kinds                  []string            `json:"kinds"`
	ProjectionSurfaces     []string            `json:"projection_surfaces"`
	ComponentConsumers     []string            `json:"component_consumers"`
	Units                  []SemanticIndexUnit `json:"units"`
}

// SemanticIndexUnit is one semantic unit entry in the index.
type SemanticIndexUnit struct {
	ID                 string                 `json:"id"`
	Kind               string                 `json:"kind"`
	Locale             string                 `json:"locale"`
	AuthorityType      string                 `json:"authority_type"`
	SourceStatus       string                 `json:"source_status"`
	Title              string                 `json:"title"`
	SourcePath         string                 `json:"source_path"`
	ProjectionPaths    map[string]string      `json:"projection_paths"`
	Relationships      map[string][]string    `json:"relationships"`
	Index              SemanticIndexUnitIndex `json:"index"`
	Provenance         map[string]any         `json:"provenance"`
	ComponentConsumers []string               `json:"component_consumers"`
	SourceUnit         string                 `json:"source_unit,omitempty"`
	TranslationStatus  string                 `json:"translation_status,omitempty"`
}

// SemanticIndexUnitIndex is the retrieval metadata block of a unit.
type SemanticIndexUnitIndex struct {
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

// ValidateSemanticIndex validates docs/core/semantic/core-v2.index.json,
// migrated from scripts/validate_semantic_index.go.
func ValidateSemanticIndex(indexFile string) error {
	index, err := loadSemanticIndex(indexFile)
	if err != nil {
		return err
	}
	if err := validateIndexTopLevelFields(index); err != nil {
		return err
	}

	authorityTypes, sourceStatuses, kinds, projectionSurfaces, componentConsumers, err := buildIndexLookupSets(index)
	if err != nil {
		return err
	}

	repoRoot, err := repoRootFromFile(indexFile)
	if err != nil {
		return err
	}

	if err := validateIndexUnits(repoRoot, index.Units, authorityTypes, sourceStatuses, kinds, projectionSurfaces, componentConsumers); err != nil {
		return err
	}

	fmt.Printf("semantic-index validation passed: %s (%d units)\n", indexFile, len(index.Units))
	return nil
}

// loadSemanticIndex reads and parses the semantic index file.
func loadSemanticIndex(indexFile string) (SemanticIndex, error) {
	bytes, err := os.ReadFile(indexFile) //nolint:gosec // G304: path is a caller-supplied validation target, by design
	if err != nil {
		return SemanticIndex{}, err
	}
	var index SemanticIndex
	if err := json.Unmarshal(bytes, &index); err != nil {
		return SemanticIndex{}, fmt.Errorf("invalid JSON: %w", err)
	}
	return index, nil
}

// buildIndexLookupSets builds the five top-level allowed-value sets used to
// validate every unit, returning the first error encountered.
func buildIndexLookupSets(index SemanticIndex) (authorityTypes, sourceStatuses, kinds, projectionSurfaces, componentConsumers map[string]bool, err error) {
	authorityTypes, err = requiredSet("authority_types", index.AuthorityTypes)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	sourceStatuses, err = requiredSet("source_statuses", index.SourceStatuses)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	kinds, err = requiredSet("kinds", index.Kinds)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	projectionSurfaces, err = requiredSet("projection_surfaces", index.ProjectionSurfaces)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	componentConsumers, err = requiredSet("component_consumers", index.ComponentConsumers)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	return authorityTypes, sourceStatuses, kinds, projectionSurfaces, componentConsumers, nil
}

// validateIndexUnits validates every unit's own contract, then (in a second
// pass, once every unit ID is known) every unit's cross-unit relationships.
func validateIndexUnits(
	repoRoot string,
	units []SemanticIndexUnit,
	authorityTypes, sourceStatuses, kinds, projectionSurfaces, componentConsumers map[string]bool,
) error {
	unitIDs := make(map[string]bool)
	for _, unit := range units {
		if err := validateIndexUnit(repoRoot, unit, authorityTypes, sourceStatuses, kinds, projectionSurfaces, componentConsumers, unitIDs); err != nil {
			return err
		}
		unitIDs[unit.ID] = true
	}
	for _, unit := range units {
		if err := validateIndexRelationships(unit, unitIDs); err != nil {
			return err
		}
	}
	return nil
}

func validateIndexTopLevelFields(index SemanticIndex) error {
	if index.Schema != expectedSemanticIndexSchema {
		return fmt.Errorf("schema must be `%s`", expectedSemanticIndexSchema)
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
	return nil
}

func validateIndexUnit(
	repoRoot string,
	unit SemanticIndexUnit,
	authorityTypes map[string]bool,
	sourceStatuses map[string]bool,
	kinds map[string]bool,
	projectionSurfaces map[string]bool,
	componentConsumers map[string]bool,
	seenIDs map[string]bool,
) error {
	if err := validateIndexUnitIdentity(repoRoot, unit, authorityTypes, sourceStatuses, kinds, seenIDs); err != nil {
		return err
	}
	if err := validateIndexUnitProjections(repoRoot, unit, projectionSurfaces); err != nil {
		return err
	}
	if err := validateIndexUnitRetrieval(repoRoot, unit, componentConsumers); err != nil {
		return err
	}
	if err := validateIndexUnitConsumersAndContract(unit, sourceStatuses, componentConsumers); err != nil {
		return err
	}
	if unit.AuthorityType == "generated_artifact" && unit.SourceStatus == "canonical" {
		return fmt.Errorf("%s: generated artifacts must not use canonical source_status", unit.ID)
	}
	return nil
}

func validateIndexUnitIdentity(
	repoRoot string,
	unit SemanticIndexUnit,
	authorityTypes map[string]bool,
	sourceStatuses map[string]bool,
	kinds map[string]bool,
	seenIDs map[string]bool,
) error {
	if err := validateIndexUnitIdentityCore(unit, kinds, seenIDs); err != nil {
		return err
	}
	return validateIndexUnitClassification(repoRoot, unit, authorityTypes, sourceStatuses)
}

// validateIndexUnitIdentityCore validates the unit's own identity fields:
// ID, uniqueness, kind, and locale.
func validateIndexUnitIdentityCore(unit SemanticIndexUnit, kinds map[string]bool, seenIDs map[string]bool) error {
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
	return nil
}

// validateIndexUnitClassification validates the unit's authority/status/
// title classification and its source_path.
func validateIndexUnitClassification(repoRoot string, unit SemanticIndexUnit, authorityTypes, sourceStatuses map[string]bool) error {
	if unit.AuthorityType == "" || !authorityTypes[unit.AuthorityType] {
		return fmt.Errorf("%s: unknown authority_type `%s`", unit.ID, unit.AuthorityType)
	}
	if unit.SourceStatus == "" || !sourceStatuses[unit.SourceStatus] {
		return fmt.Errorf("%s: unknown source_status `%s`", unit.ID, unit.SourceStatus)
	}
	if unit.Title == "" {
		return fmt.Errorf("%s: title must be non-empty", unit.ID)
	}
	return validateIndexPath(repoRoot, unit.SourcePath, unit.ID, "source_path")
}

func validateIndexUnitProjections(repoRoot string, unit SemanticIndexUnit, projectionSurfaces map[string]bool) error {
	if len(unit.ProjectionPaths) == 0 {
		return fmt.Errorf("%s: projection_paths must be non-empty", unit.ID)
	}
	for surface, path := range unit.ProjectionPaths {
		if !projectionSurfaces[surface] {
			return fmt.Errorf("%s: unknown projection surface `%s`", unit.ID, surface)
		}
		if err := validateIndexPath(repoRoot, path, unit.ID, "projection_paths."+surface); err != nil {
			return err
		}
	}
	if unit.AuthorityType == "canonical_markdown_bridge" && unit.ProjectionPaths["markdown_en"] == "" {
		return fmt.Errorf("%s: canonical Markdown bridge units must expose projection_paths.markdown_en", unit.ID)
	}
	return nil
}

func validateIndexUnitRetrieval(repoRoot string, unit SemanticIndexUnit, componentConsumers map[string]bool) error {
	if unit.Relationships == nil {
		return fmt.Errorf("%s: relationships must be an object", unit.ID)
	}
	if unit.Index.RetrievalSummary == "" {
		return fmt.Errorf("%s: index.retrieval_summary must be non-empty", unit.ID)
	}
	if len(unit.Index.Track) == 0 {
		return fmt.Errorf("%s: index.track must be non-empty", unit.ID)
	}
	if err := validateIndexTrackConsumers(unit, componentConsumers); err != nil {
		return err
	}
	if len(unit.Index.Tags) == 0 {
		return fmt.Errorf("%s: index.tags must be non-empty", unit.ID)
	}
	if len(unit.Provenance) == 0 {
		return fmt.Errorf("%s: provenance must be non-empty", unit.ID)
	}
	return validateIndexProvenance(repoRoot, unit)
}

// validateIndexTrackConsumers validates that every index.track consumer is
// a known component consumer.
func validateIndexTrackConsumers(unit SemanticIndexUnit, componentConsumers map[string]bool) error {
	for _, consumer := range unit.Index.Track {
		if !componentConsumers[consumer] {
			return fmt.Errorf("%s: unknown index track consumer `%s`", unit.ID, consumer)
		}
	}
	return nil
}

func validateIndexUnitConsumersAndContract(unit SemanticIndexUnit, sourceStatuses map[string]bool, componentConsumers map[string]bool) error {
	if len(unit.ComponentConsumers) == 0 {
		return fmt.Errorf("%s: component_consumers must be non-empty", unit.ID)
	}
	for _, consumer := range unit.ComponentConsumers {
		if !componentConsumers[consumer] {
			return fmt.Errorf("%s: unknown component consumer `%s`", unit.ID, consumer)
		}
	}
	return validateIndexTranslationContract(unit, sourceStatuses)
}

func validateIndexProvenance(repoRoot string, unit SemanticIndexUnit) error {
	sourceRevision, ok := unit.Provenance["source_revision"].(string)
	if !ok || sourceRevision == "" {
		return fmt.Errorf("%s: provenance.source_revision must be a non-empty string", unit.ID)
	}

	refsValue, exists := unit.Provenance["decision_refs"]
	if !exists {
		return nil
	}
	return validateIndexDecisionRefs(repoRoot, unit.ID, refsValue)
}

// validateIndexDecisionRefs validates the optional provenance.decision_refs
// list: each entry must be a string and a valid repository-relative path.
func validateIndexDecisionRefs(repoRoot, unitID string, refsValue any) error {
	refs, ok := refsValue.([]any)
	if !ok {
		return fmt.Errorf("%s: provenance.decision_refs must be a string list", unitID)
	}
	for index, refValue := range refs {
		ref, ok := refValue.(string)
		if !ok {
			return fmt.Errorf("%s: provenance.decision_refs[%d] must be a string", unitID, index)
		}
		if err := validateIndexPath(repoRoot, ref, unitID, fmt.Sprintf("provenance.decision_refs[%d]", index)); err != nil {
			return err
		}
	}
	return nil
}

func validateIndexTranslationContract(unit SemanticIndexUnit, sourceStatuses map[string]bool) error {
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

func validateIndexRelationships(unit SemanticIndexUnit, unitIDs map[string]bool) error {
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

func validateIndexPath(repoRoot, pathValue, unitID, field string) error {
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
