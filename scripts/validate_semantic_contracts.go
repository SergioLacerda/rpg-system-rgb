package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const expectedContractsSchema = "rgb-docs-consumer-contracts/0.1"

type consumerContractsFile struct {
	Schema      string             `json:"schema"`
	SourceIndex string             `json:"source_index"`
	Description string             `json:"description"`
	Contracts   []consumerContract `json:"contracts"`
}

type consumerContract struct {
	ID                        string   `json:"id"`
	Component                 string   `json:"component"`
	Status                    string   `json:"status"`
	Description               string   `json:"description"`
	AllowedInputs             []string `json:"allowed_inputs"`
	RequiredUnitFields        []string `json:"required_unit_fields"`
	AllowedProjectionSurfaces []string `json:"allowed_projection_surfaces"`
	ForbiddenAuthorityTypes   []string `json:"forbidden_authority_types"`
	RequiredDisclosures       []string `json:"required_disclosures"`
	Outputs                   []string `json:"outputs"`
	ValidationChecks          []string `json:"validation_checks"`
	SourceRefs                []string `json:"source_refs"`
}

type contractSourceIndex struct {
	AuthorityTypes     []string `json:"authority_types"`
	SourceStatuses     []string `json:"source_statuses"`
	ProjectionSurfaces []string `json:"projection_surfaces"`
	ComponentConsumers []string `json:"component_consumers"`
}

func main() {
	if len(os.Args) != 3 {
		fatal(errors.New("usage: go run scripts/validate_semantic_contracts.go <contracts.json> <index.json>"))
	}
	if err := validateContracts(os.Args[1], os.Args[2]); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "semantic-contract validation failed: %v\n", err)
	os.Exit(1)
}

func validateContracts(contractsFile, indexFile string) error {
	contractsBytes, err := os.ReadFile(contractsFile)
	if err != nil {
		return err
	}
	indexBytes, err := os.ReadFile(indexFile)
	if err != nil {
		return err
	}

	var contracts consumerContractsFile
	if err := json.Unmarshal(contractsBytes, &contracts); err != nil {
		return fmt.Errorf("invalid contracts JSON: %w", err)
	}
	var index contractSourceIndex
	if err := json.Unmarshal(indexBytes, &index); err != nil {
		return fmt.Errorf("invalid source index JSON: %w", err)
	}

	if contracts.Schema != expectedContractsSchema {
		return fmt.Errorf("schema must be `%s`", expectedContractsSchema)
	}
	if contracts.SourceIndex != indexFile {
		return fmt.Errorf("source_index must match validation index path `%s`", indexFile)
	}
	if contracts.Description == "" {
		return errors.New("description must be non-empty")
	}
	if len(contracts.Contracts) == 0 {
		return errors.New("contracts must be non-empty")
	}

	repoRoot, err := repoRootFromContracts(contractsFile)
	if err != nil {
		return err
	}

	authorityTypes := stringSet(index.AuthorityTypes)
	sourceStatuses := stringSet(index.SourceStatuses)
	projectionSurfaces := stringSet(index.ProjectionSurfaces)
	componentConsumers := stringSet(index.ComponentConsumers)
	knownUnitFields := stringSet([]string{
		"id",
		"kind",
		"locale",
		"authority_type",
		"source_status",
		"title",
		"source_path",
		"projection_paths",
		"relationships",
		"index",
		"provenance",
		"component_consumers",
		"source_unit",
		"translation_status",
	})

	seenIDs := map[string]bool{}
	for _, contract := range contracts.Contracts {
		if err := validateContract(repoRoot, contract, seenIDs, authorityTypes, sourceStatuses, projectionSurfaces, componentConsumers, knownUnitFields); err != nil {
			return err
		}
		seenIDs[contract.ID] = true
	}

	fmt.Printf("semantic-contract validation passed: %s (%d contracts)\n", contractsFile, len(contracts.Contracts))
	return nil
}

func repoRootFromContracts(contractsFile string) (string, error) {
	abs, err := filepath.Abs(contractsFile)
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

func validateContract(
	repoRoot string,
	contract consumerContract,
	seenIDs map[string]bool,
	authorityTypes map[string]bool,
	sourceStatuses map[string]bool,
	projectionSurfaces map[string]bool,
	componentConsumers map[string]bool,
	knownUnitFields map[string]bool,
) error {
	if contract.ID == "" {
		return errors.New("contract id must be non-empty")
	}
	if seenIDs[contract.ID] {
		return fmt.Errorf("duplicate contract id: %s", contract.ID)
	}
	if contract.Component == "" || !componentConsumers[contract.Component] {
		return fmt.Errorf("%s: unknown component `%s`", contract.ID, contract.Component)
	}
	if contract.Status == "" || !sourceStatuses[contract.Status] {
		return fmt.Errorf("%s: unknown status `%s`", contract.ID, contract.Status)
	}
	if contract.Description == "" {
		return fmt.Errorf("%s: description must be non-empty", contract.ID)
	}
	if err := validateKnownValues(contract.ID, "allowed_inputs", contract.AllowedInputs, authorityTypes); err != nil {
		return err
	}
	if err := validateKnownValues(contract.ID, "required_unit_fields", contract.RequiredUnitFields, knownUnitFields); err != nil {
		return err
	}
	if err := validateKnownValues(contract.ID, "allowed_projection_surfaces", contract.AllowedProjectionSurfaces, projectionSurfaces); err != nil {
		return err
	}
	if len(contract.ForbiddenAuthorityTypes) > 0 {
		if err := validateKnownValues(contract.ID, "forbidden_authority_types", contract.ForbiddenAuthorityTypes, authorityTypes); err != nil {
			return err
		}
	}
	if err := validateNonEmptyStrings(contract.ID, "required_disclosures", contract.RequiredDisclosures); err != nil {
		return err
	}
	if err := validateNonEmptyStrings(contract.ID, "outputs", contract.Outputs); err != nil {
		return err
	}
	if err := validateNonEmptyStrings(contract.ID, "validation_checks", contract.ValidationChecks); err != nil {
		return err
	}
	if err := validateNonEmptyStrings(contract.ID, "source_refs", contract.SourceRefs); err != nil {
		return err
	}
	for index, ref := range contract.SourceRefs {
		if err := validateRepoPath(repoRoot, ref, contract.ID, fmt.Sprintf("source_refs[%d]", index)); err != nil {
			return err
		}
	}
	return nil
}

func validateKnownValues(contractID, field string, values []string, allowed map[string]bool) error {
	if err := validateNonEmptyStrings(contractID, field, values); err != nil {
		return err
	}
	for _, value := range values {
		if !allowed[value] {
			return fmt.Errorf("%s: %s contains unknown value `%s`", contractID, field, value)
		}
	}
	return nil
}

func validateNonEmptyStrings(contractID, field string, values []string) error {
	if len(values) == 0 {
		return fmt.Errorf("%s: %s must be non-empty", contractID, field)
	}
	for index, value := range values {
		if value == "" {
			return fmt.Errorf("%s: %s[%d] must be non-empty", contractID, field, index)
		}
	}
	return nil
}

func validateRepoPath(repoRoot, pathValue, contractID, field string) error {
	if filepath.IsAbs(pathValue) {
		return fmt.Errorf("%s: %s must be repository-relative", contractID, field)
	}
	fullPath := filepath.Join(repoRoot, filepath.FromSlash(pathValue))
	if _, err := os.Stat(fullPath); err != nil {
		return fmt.Errorf("%s: %s does not exist: %s", contractID, field, pathValue)
	}
	return nil
}

func stringSet(values []string) map[string]bool {
	result := make(map[string]bool, len(values))
	for _, value := range values {
		result[value] = true
	}
	return result
}
