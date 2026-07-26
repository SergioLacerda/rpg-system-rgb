package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type generatedValidationManifest struct {
	Projections []struct {
		ID          string   `json:"id"`
		Surface     string   `json:"surface"`
		Owner       string   `json:"owner"`
		SourceUnits []string `json:"source_units"`
		OutputPath  string   `json:"output_path"`
	} `json:"projections"`
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

func main() {
	if len(os.Args) != 2 {
		fatal(errors.New("usage: go run scripts/validate_generated_projections.go <manifest.json>"))
	}
	if err := validateGeneratedProjections(os.Args[1]); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "generated-projection validation failed: %v\n", err)
	os.Exit(1)
}

func validateGeneratedProjections(manifestFile string) error {
	var manifest generatedValidationManifest
	if err := readGeneratedJSON(manifestFile, &manifest); err != nil {
		return err
	}
	if len(manifest.Projections) == 0 {
		return errors.New("manifest projections must be non-empty")
	}
	for _, projection := range manifest.Projections {
		var output generatedValidationOutput
		if err := readGeneratedJSON(projection.OutputPath, &output); err != nil {
			return err
		}
		if output.Schema != "rgb-docs-generated-projection/0.1" {
			return fmt.Errorf("%s: unexpected generated schema `%s`", projection.OutputPath, output.Schema)
		}
		if output.ProjectionID != projection.ID {
			return fmt.Errorf("%s: projection_id mismatch", projection.OutputPath)
		}
		if output.Surface != projection.Surface {
			return fmt.Errorf("%s: surface mismatch", projection.OutputPath)
		}
		if output.Owner != projection.Owner {
			return fmt.Errorf("%s: owner mismatch", projection.OutputPath)
		}
		if output.AuthorityType != "generated_artifact" {
			return fmt.Errorf("%s: authority_type must be generated_artifact", projection.OutputPath)
		}
		if len(output.GeneratedFrom) != len(projection.SourceUnits) {
			return fmt.Errorf("%s: generated_from count mismatch", projection.OutputPath)
		}
		if len(output.Units) != len(projection.SourceUnits) {
			return fmt.Errorf("%s: units count mismatch", projection.OutputPath)
		}
	}
	fmt.Printf("generated-projection validation passed: %s (%d outputs)\n", manifestFile, len(manifest.Projections))
	return nil
}

func readGeneratedJSON(path string, target any) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, target); err != nil {
		return fmt.Errorf("%s: invalid JSON: %w", path, err)
	}
	return nil
}
