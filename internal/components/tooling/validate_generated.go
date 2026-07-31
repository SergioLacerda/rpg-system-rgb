package tooling

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type generatedValidationManifest struct {
	Projections []generatedValidationManifestProjection `json:"projections"`
}

type generatedValidationManifestProjection struct {
	ID          string   `json:"id"`
	Surface     string   `json:"surface"`
	Owner       string   `json:"owner"`
	SourceUnits []string `json:"source_units"`
	OutputPath  string   `json:"output_path"`
}

type generatedValidationOutput struct {
	Schema        string   `json:"schema"`
	ProjectionID  string   `json:"projection_id"`
	Surface       string   `json:"surface"`
	Owner         string   `json:"owner"`
	AuthorityType string   `json:"authority_type"`
	GeneratedFrom []string `json:"generated_from"`
	Units         []struct {
		ID string `json:"id"`
	} `json:"units"`
}

// ValidateGeneratedProjections validates that every generated/*/*.json
// output matches its declaration in the projection manifest, migrated from
// scripts/validate_generated_projections.go. Unlike the original script,
// this resolves output paths against the repository root explicitly
// (the original assumed the process working directory was the repo root).
func ValidateGeneratedProjections(manifestFile string) error {
	repoRoot, err := repoRootFromFile(manifestFile)
	if err != nil {
		return err
	}

	var manifest generatedValidationManifest
	if err := readGeneratedJSON(manifestFile, &manifest); err != nil {
		return err
	}
	if len(manifest.Projections) == 0 {
		return errors.New("manifest projections must be non-empty")
	}
	for _, proj := range manifest.Projections {
		outputPath := filepath.Join(repoRoot, filepath.FromSlash(proj.OutputPath))
		var output generatedValidationOutput
		if err := readGeneratedJSON(outputPath, &output); err != nil {
			return err
		}
		if err := validateGeneratedProjectionOutput(proj, output); err != nil {
			return err
		}
	}
	fmt.Printf("generated-projection validation passed: %s (%d outputs)\n", manifestFile, len(manifest.Projections))
	return nil
}

// validateGeneratedProjectionOutput checks a single generated projection
// output file against its declaration in the projection manifest.
func validateGeneratedProjectionOutput(proj generatedValidationManifestProjection, output generatedValidationOutput) error {
	if err := validateGeneratedProjectionIdentity(proj, output); err != nil {
		return err
	}
	return validateGeneratedProjectionShape(proj, output)
}

// validateGeneratedProjectionIdentity validates that the output's schema
// and projection identity match its manifest declaration.
func validateGeneratedProjectionIdentity(proj generatedValidationManifestProjection, output generatedValidationOutput) error {
	if output.Schema != "rgb-docs-generated-projection/0.1" {
		return fmt.Errorf("%s: unexpected generated schema `%s`", proj.OutputPath, output.Schema)
	}
	if output.ProjectionID != proj.ID {
		return fmt.Errorf("%s: projection_id mismatch", proj.OutputPath)
	}
	if output.Surface != proj.Surface {
		return fmt.Errorf("%s: surface mismatch", proj.OutputPath)
	}
	if output.Owner != proj.Owner {
		return fmt.Errorf("%s: owner mismatch", proj.OutputPath)
	}
	return nil
}

// validateGeneratedProjectionShape validates the output's generated-artifact
// authority type and that its unit counts match the manifest declaration.
func validateGeneratedProjectionShape(proj generatedValidationManifestProjection, output generatedValidationOutput) error {
	if output.AuthorityType != "generated_artifact" {
		return fmt.Errorf("%s: authority_type must be generated_artifact", proj.OutputPath)
	}
	if len(output.GeneratedFrom) != len(proj.SourceUnits) {
		return fmt.Errorf("%s: generated_from count mismatch", proj.OutputPath)
	}
	if len(output.Units) != len(proj.SourceUnits) {
		return fmt.Errorf("%s: units count mismatch", proj.OutputPath)
	}
	return nil
}

func readGeneratedJSON(path string, target any) error {
	bytes, err := os.ReadFile(path) //nolint:gosec // G304: path is a caller-supplied validation target, by design
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, target); err != nil {
		return fmt.Errorf("%s: invalid JSON: %w", path, err)
	}
	return nil
}
