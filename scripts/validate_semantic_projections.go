package main

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

func main() {
	if len(os.Args) != 4 {
		fatal(errors.New("usage: go run scripts/validate_semantic_projections.go <manifest.json> <index.json> <contracts.json>"))
	}
	if err := validateProjections(os.Args[1], os.Args[2], os.Args[3]); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "semantic-projection validation failed: %v\n", err)
	os.Exit(1)
}

func validateProjections(manifestFile, indexFile, contractsFile string) error {
	manifestBytes, err := os.ReadFile(manifestFile)
	if err != nil {
		return err
	}
	indexBytes, err := os.ReadFile(indexFile)
	if err != nil {
		return err
	}
	contractsBytes, err := os.ReadFile(contractsFile)
	if err != nil {
		return err
	}

	var manifest projectionManifest
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		return fmt.Errorf("invalid manifest JSON: %w", err)
	}
	var index projectionSourceIndex
	if err := json.Unmarshal(indexBytes, &index); err != nil {
		return fmt.Errorf("invalid source index JSON: %w", err)
	}
	var contracts projectionContractsFile
	if err := json.Unmarshal(contractsBytes, &contracts); err != nil {
		return fmt.Errorf("invalid consumer contracts JSON: %w", err)
	}

	if manifest.Schema != expectedProjectionSchema {
		return fmt.Errorf("schema must be `%s`", expectedProjectionSchema)
	}
	if manifest.SourceIndex != indexFile {
		return fmt.Errorf("source_index must match validation index path `%s`", indexFile)
	}
	if manifest.ConsumerContracts != contractsFile {
		return fmt.Errorf("consumer_contracts must match validation contracts path `%s`", contractsFile)
	}
	if manifest.Description == "" {
		return errors.New("description must be non-empty")
	}
	if len(manifest.Projections) == 0 {
		return errors.New("projections must be non-empty")
	}

	repoRoot, err := repoRootFromManifest(manifestFile)
	if err != nil {
		return err
	}

	statuses := projectionStringSet(index.SourceStatuses)
	surfaces := projectionStringSet(index.ProjectionSurfaces)
	components := projectionStringSet(index.ComponentConsumers)
	unitIDs := map[string]bool{}
	for _, unit := range index.Units {
		unitIDs[unit.ID] = true
	}
	allowedSurfacesByOwner := map[string]map[string]bool{}
	for _, contract := range contracts.Contracts {
		allowedSurfacesByOwner[contract.Component] = projectionStringSet(contract.AllowedProjectionSurfaces)
	}

	seenIDs := map[string]bool{}
	for _, projection := range manifest.Projections {
		if err := validateProjection(repoRoot, projection, seenIDs, statuses, surfaces, components, unitIDs, allowedSurfacesByOwner); err != nil {
			return err
		}
		seenIDs[projection.ID] = true
	}

	fmt.Printf("semantic-projection validation passed: %s (%d projections)\n", manifestFile, len(manifest.Projections))
	return nil
}

func repoRootFromManifest(manifestFile string) (string, error) {
	abs, err := filepath.Abs(manifestFile)
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

func validateProjection(
	repoRoot string,
	projection projection,
	seenIDs map[string]bool,
	statuses map[string]bool,
	surfaces map[string]bool,
	components map[string]bool,
	unitIDs map[string]bool,
	allowedSurfacesByOwner map[string]map[string]bool,
) error {
	if projection.ID == "" {
		return errors.New("projection id must be non-empty")
	}
	if seenIDs[projection.ID] {
		return fmt.Errorf("duplicate projection id: %s", projection.ID)
	}
	if projection.Surface == "" || !surfaces[projection.Surface] {
		return fmt.Errorf("%s: unknown surface `%s`", projection.ID, projection.Surface)
	}
	if projection.Owner == "" || !components[projection.Owner] {
		return fmt.Errorf("%s: unknown owner `%s`", projection.ID, projection.Owner)
	}
	ownerSurfaces, ok := allowedSurfacesByOwner[projection.Owner]
	if !ok {
		return fmt.Errorf("%s: owner `%s` has no consumer contract", projection.ID, projection.Owner)
	}
	if !ownerSurfaces[projection.Surface] {
		return fmt.Errorf("%s: owner `%s` contract does not allow surface `%s`", projection.ID, projection.Owner, projection.Surface)
	}
	if projection.Status == "" || !statuses[projection.Status] {
		return fmt.Errorf("%s: unknown status `%s`", projection.ID, projection.Status)
	}
	if projection.Description == "" {
		return fmt.Errorf("%s: description must be non-empty", projection.ID)
	}
	if err := validateProjectionSourceUnits(projection, unitIDs); err != nil {
		return err
	}
	if projection.OutputPath == "" {
		return fmt.Errorf("%s: output_path must be non-empty", projection.ID)
	}
	if filepath.IsAbs(projection.OutputPath) {
		return fmt.Errorf("%s: output_path must be repository-relative", projection.ID)
	}
	if len(projection.Provenance) == 0 {
		return fmt.Errorf("%s: provenance must be non-empty", projection.ID)
	}
	if err := validateProjectionProvenance(repoRoot, projection); err != nil {
		return err
	}
	if err := validateProjectionStrings(projection.ID, "required_disclosures", projection.RequiredDisclosures); err != nil {
		return err
	}
	if err := validateProjectionStrings(projection.ID, "generation_gate", projection.GenerationGate); err != nil {
		return err
	}
	return nil
}

func validateProjectionSourceUnits(projection projection, unitIDs map[string]bool) error {
	if len(projection.SourceUnits) == 0 {
		return fmt.Errorf("%s: source_units must be non-empty", projection.ID)
	}
	for index, unitID := range projection.SourceUnits {
		if unitID == "" {
			return fmt.Errorf("%s: source_units[%d] must be non-empty", projection.ID, index)
		}
		if !unitIDs[unitID] {
			return fmt.Errorf("%s: source_units[%d] points to unknown unit `%s`", projection.ID, index, unitID)
		}
	}
	return nil
}

func validateProjectionProvenance(repoRoot string, projection projection) error {
	sourceRevision, ok := projection.Provenance["source_revision"].(string)
	if !ok || sourceRevision == "" {
		return fmt.Errorf("%s: provenance.source_revision must be a non-empty string", projection.ID)
	}

	refsValue, exists := projection.Provenance["decision_refs"]
	if !exists {
		return nil
	}
	refs, ok := refsValue.([]any)
	if !ok {
		return fmt.Errorf("%s: provenance.decision_refs must be a string list", projection.ID)
	}
	for index, refValue := range refs {
		ref, ok := refValue.(string)
		if !ok {
			return fmt.Errorf("%s: provenance.decision_refs[%d] must be a string", projection.ID, index)
		}
		if filepath.IsAbs(ref) {
			return fmt.Errorf("%s: provenance.decision_refs[%d] must be repository-relative", projection.ID, index)
		}
		fullPath := filepath.Join(repoRoot, filepath.FromSlash(ref))
		if _, err := os.Stat(fullPath); err != nil {
			return fmt.Errorf("%s: provenance.decision_refs[%d] does not exist: %s", projection.ID, index, ref)
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

func projectionStringSet(values []string) map[string]bool {
	result := make(map[string]bool, len(values))
	for _, value := range values {
		result[value] = true
	}
	return result
}
