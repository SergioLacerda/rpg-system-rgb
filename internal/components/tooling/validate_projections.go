package tooling

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const expectedProjectionSchema = "rgb-docs-projection-manifest/0.1"

type projectionManifest struct {
	Schema            string       `json:"schema"`
	SourceIndex       string       `json:"source_index"`
	ConsumerContracts string       `json:"consumer_contracts"`
	Description       string       `json:"description"`
	Projections       []projection `json:"projections"`
}

type projection struct {
	ID                  string         `json:"id"`
	Surface             string         `json:"surface"`
	Owner               string         `json:"owner"`
	Status              string         `json:"status"`
	Description         string         `json:"description"`
	SourceUnits         []string       `json:"source_units"`
	OutputPath          string         `json:"output_path"`
	Provenance          map[string]any `json:"provenance"`
	RequiredDisclosures []string       `json:"required_disclosures"`
	GenerationGate      []string       `json:"generation_gate"`
}

type projectionSourceIndex struct {
	SourceStatuses     []string         `json:"source_statuses"`
	ProjectionSurfaces []string         `json:"projection_surfaces"`
	ComponentConsumers []string         `json:"component_consumers"`
	Units              []projectionUnit `json:"units"`
}

type projectionUnit struct {
	ID string `json:"id"`
}

type projectionContractsFile struct {
	Contracts []projectionConsumerContract `json:"contracts"`
}

type projectionConsumerContract struct {
	Component                 string   `json:"component"`
	AllowedProjectionSurfaces []string `json:"allowed_projection_surfaces"`
}

// ValidateProjectionManifest validates docs/core/semantic/projection-manifest.v0.1.json
// against its index and consumer contracts, migrated from
// scripts/validate_semantic_projections.go.
func ValidateProjectionManifest(manifestFile, indexFile, contractsFile string) error {
	manifest, index, contracts, err := loadProjectionInputs(manifestFile, indexFile, contractsFile)
	if err != nil {
		return err
	}

	repoRoot, err := repoRootFromFile(manifestFile)
	if err != nil {
		return err
	}
	if err := validateProjectionManifestTopLevel(manifest, repoRoot, indexFile, contractsFile); err != nil {
		return err
	}

	statuses, surfaces, components, unitIDs, allowedSurfacesByOwner := buildProjectionLookups(index, contracts)

	seenIDs := map[string]bool{}
	for _, proj := range manifest.Projections {
		if err := validateProjectionEntry(repoRoot, proj, seenIDs, statuses, surfaces, components, unitIDs, allowedSurfacesByOwner); err != nil {
			return err
		}
		seenIDs[proj.ID] = true
	}

	fmt.Printf("semantic-projection validation passed: %s (%d projections)\n", manifestFile, len(manifest.Projections))
	return nil
}

// loadProjectionInputs reads and parses the projection manifest, its parent
// source index, and the consumer contracts file.
func loadProjectionInputs(manifestFile, indexFile, contractsFile string) (projectionManifest, projectionSourceIndex, projectionContractsFile, error) {
	manifestBytes, err := os.ReadFile(manifestFile) //nolint:gosec // G304: path is a caller-supplied validation target, by design
	if err != nil {
		return projectionManifest{}, projectionSourceIndex{}, projectionContractsFile{}, err
	}
	indexBytes, err := os.ReadFile(indexFile) //nolint:gosec // G304: path is a caller-supplied validation target, by design
	if err != nil {
		return projectionManifest{}, projectionSourceIndex{}, projectionContractsFile{}, err
	}
	contractsBytes, err := os.ReadFile(contractsFile) //nolint:gosec // G304: path is a caller-supplied validation target, by design
	if err != nil {
		return projectionManifest{}, projectionSourceIndex{}, projectionContractsFile{}, err
	}

	var manifest projectionManifest
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		return projectionManifest{}, projectionSourceIndex{}, projectionContractsFile{}, fmt.Errorf("invalid manifest JSON: %w", err)
	}
	var index projectionSourceIndex
	if err := json.Unmarshal(indexBytes, &index); err != nil {
		return projectionManifest{}, projectionSourceIndex{}, projectionContractsFile{}, fmt.Errorf("invalid source index JSON: %w", err)
	}
	var contracts projectionContractsFile
	if err := json.Unmarshal(contractsBytes, &contracts); err != nil {
		return projectionManifest{}, projectionSourceIndex{}, projectionContractsFile{}, fmt.Errorf("invalid consumer contracts JSON: %w", err)
	}
	return manifest, index, contracts, nil
}

// buildProjectionLookups builds the allowed-value sets and ID-keyed lookup
// maps used to validate every projection entry.
func buildProjectionLookups(index projectionSourceIndex, contracts projectionContractsFile) (
	statuses, surfaces, components map[string]bool,
	unitIDs map[string]bool,
	allowedSurfacesByOwner map[string]map[string]bool,
) {
	statuses = stringSet(index.SourceStatuses)
	surfaces = stringSet(index.ProjectionSurfaces)
	components = stringSet(index.ComponentConsumers)
	unitIDs = map[string]bool{}
	for _, unit := range index.Units {
		unitIDs[unit.ID] = true
	}
	allowedSurfacesByOwner = map[string]map[string]bool{}
	for _, contract := range contracts.Contracts {
		allowedSurfacesByOwner[contract.Component] = stringSet(contract.AllowedProjectionSurfaces)
	}
	return statuses, surfaces, components, unitIDs, allowedSurfacesByOwner
}

func validateProjectionManifestTopLevel(manifest projectionManifest, repoRoot, indexFile, contractsFile string) error {
	if manifest.Schema != expectedProjectionSchema {
		return fmt.Errorf("schema must be `%s`", expectedProjectionSchema)
	}
	if manifest.SourceIndex != repoRelative(repoRoot, indexFile) {
		return fmt.Errorf("source_index must match validation index path `%s`", indexFile)
	}
	if manifest.ConsumerContracts != repoRelative(repoRoot, contractsFile) {
		return fmt.Errorf("consumer_contracts must match validation contracts path `%s`", contractsFile)
	}
	if manifest.Description == "" {
		return errors.New("description must be non-empty")
	}
	if len(manifest.Projections) == 0 {
		return errors.New("projections must be non-empty")
	}
	return nil
}

func validateProjectionEntry(
	repoRoot string,
	proj projection,
	seenIDs map[string]bool,
	statuses map[string]bool,
	surfaces map[string]bool,
	components map[string]bool,
	unitIDs map[string]bool,
	allowedSurfacesByOwner map[string]map[string]bool,
) error {
	if err := validateProjectionOwnership(proj, seenIDs, statuses, surfaces, components, allowedSurfacesByOwner); err != nil {
		return err
	}
	if err := validateProjectionSourceUnits(proj, unitIDs); err != nil {
		return err
	}
	return validateProjectionOutput(repoRoot, proj)
}

func validateProjectionOwnership(
	proj projection,
	seenIDs map[string]bool,
	statuses map[string]bool,
	surfaces map[string]bool,
	components map[string]bool,
	allowedSurfacesByOwner map[string]map[string]bool,
) error {
	if err := validateProjectionIdentity(proj, seenIDs, surfaces, components); err != nil {
		return err
	}
	return validateProjectionContractCompliance(proj, statuses, allowedSurfacesByOwner)
}

// validateProjectionIdentity validates the projection's own ID, uniqueness,
// surface, and owner fields.
func validateProjectionIdentity(proj projection, seenIDs map[string]bool, surfaces, components map[string]bool) error {
	if proj.ID == "" {
		return errors.New("projection id must be non-empty")
	}
	if seenIDs[proj.ID] {
		return fmt.Errorf("duplicate projection id: %s", proj.ID)
	}
	if proj.Surface == "" || !surfaces[proj.Surface] {
		return fmt.Errorf("%s: unknown surface `%s`", proj.ID, proj.Surface)
	}
	if proj.Owner == "" || !components[proj.Owner] {
		return fmt.Errorf("%s: unknown owner `%s`", proj.ID, proj.Owner)
	}
	return nil
}

// validateProjectionContractCompliance validates that the projection's
// surface is allowed by its owner's consumer contract, and its status and
// description are well-formed.
func validateProjectionContractCompliance(proj projection, statuses map[string]bool, allowedSurfacesByOwner map[string]map[string]bool) error {
	ownerSurfaces, ok := allowedSurfacesByOwner[proj.Owner]
	if !ok {
		return fmt.Errorf("%s: owner `%s` has no consumer contract", proj.ID, proj.Owner)
	}
	if !ownerSurfaces[proj.Surface] {
		return fmt.Errorf("%s: owner `%s` contract does not allow surface `%s`", proj.ID, proj.Owner, proj.Surface)
	}
	if proj.Status == "" || !statuses[proj.Status] {
		return fmt.Errorf("%s: unknown status `%s`", proj.ID, proj.Status)
	}
	if proj.Description == "" {
		return fmt.Errorf("%s: description must be non-empty", proj.ID)
	}
	return nil
}

func validateProjectionOutput(repoRoot string, proj projection) error {
	if proj.OutputPath == "" {
		return fmt.Errorf("%s: output_path must be non-empty", proj.ID)
	}
	if filepath.IsAbs(proj.OutputPath) {
		return fmt.Errorf("%s: output_path must be repository-relative", proj.ID)
	}
	if len(proj.Provenance) == 0 {
		return fmt.Errorf("%s: provenance must be non-empty", proj.ID)
	}
	if err := validateProjectionProvenance(repoRoot, proj); err != nil {
		return err
	}
	if err := validateProjectionStrings(proj.ID, "required_disclosures", proj.RequiredDisclosures); err != nil {
		return err
	}
	return validateProjectionStrings(proj.ID, "generation_gate", proj.GenerationGate)
}

func validateProjectionSourceUnits(proj projection, unitIDs map[string]bool) error {
	if len(proj.SourceUnits) == 0 {
		return fmt.Errorf("%s: source_units must be non-empty", proj.ID)
	}
	for index, unitID := range proj.SourceUnits {
		if unitID == "" {
			return fmt.Errorf("%s: source_units[%d] must be non-empty", proj.ID, index)
		}
		if !unitIDs[unitID] {
			return fmt.Errorf("%s: source_units[%d] points to unknown unit `%s`", proj.ID, index, unitID)
		}
	}
	return nil
}

func validateProjectionProvenance(repoRoot string, proj projection) error {
	sourceRevision, ok := proj.Provenance["source_revision"].(string)
	if !ok || sourceRevision == "" {
		return fmt.Errorf("%s: provenance.source_revision must be a non-empty string", proj.ID)
	}

	refsValue, exists := proj.Provenance["decision_refs"]
	if !exists {
		return nil
	}
	return validateProjectionDecisionRefs(repoRoot, proj.ID, refsValue)
}

// validateProjectionDecisionRefs validates the optional
// provenance.decision_refs list: each entry must be a string and a valid
// repository-relative path.
func validateProjectionDecisionRefs(repoRoot, projectionID string, refsValue any) error {
	refs, ok := refsValue.([]any)
	if !ok {
		return fmt.Errorf("%s: provenance.decision_refs must be a string list", projectionID)
	}
	for index, refValue := range refs {
		ref, ok := refValue.(string)
		if !ok {
			return fmt.Errorf("%s: provenance.decision_refs[%d] must be a string", projectionID, index)
		}
		if filepath.IsAbs(ref) {
			return fmt.Errorf("%s: provenance.decision_refs[%d] must be repository-relative", projectionID, index)
		}
		fullPath := filepath.Join(repoRoot, filepath.FromSlash(ref))
		if _, err := os.Stat(fullPath); err != nil {
			return fmt.Errorf("%s: provenance.decision_refs[%d] does not exist: %s", projectionID, index, ref)
		}
	}
	return nil
}

func validateProjectionStrings(projectionID, field string, values []string) error {
	if len(values) == 0 {
		return fmt.Errorf("%s: %s must be non-empty", projectionID, field)
	}
	for index, value := range values {
		if value == "" {
			return fmt.Errorf("%s: %s[%d] must be non-empty", projectionID, field, index)
		}
	}
	return nil
}
