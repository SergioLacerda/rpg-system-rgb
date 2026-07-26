package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type generationManifest struct {
	Schema      string                 `json:"schema"`
	Projections []generationProjection `json:"projections"`
}

type generationProjection struct {
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

type generationIndex struct {
	Units []generationIndexUnit `json:"units"`
}

type generationIndexUnit struct {
	ID              string              `json:"id"`
	Kind            string              `json:"kind"`
	Locale          string              `json:"locale"`
	AuthorityType   string              `json:"authority_type"`
	SourceStatus    string              `json:"source_status"`
	Title           string              `json:"title"`
	SourcePath      string              `json:"source_path"`
	ProjectionPaths map[string]string   `json:"projection_paths"`
	Relationships   map[string][]string `json:"relationships"`
	Index           struct {
		Tags             []string `json:"tags"`
		RetrievalSummary string   `json:"retrieval_summary"`
	} `json:"index"`
	Provenance map[string]any `json:"provenance"`
}

type generationSource struct {
	Units []generationSourceUnit `json:"units"`
}

type generationSourceUnit struct {
	ID              string              `json:"id"`
	Kind            string              `json:"kind"`
	Locale          string              `json:"locale"`
	AuthorityType   string              `json:"authority_type"`
	SourceStatus    string              `json:"source_status"`
	Title           string              `json:"title"`
	Statement       string              `json:"statement"`
	ProjectionPaths map[string]string   `json:"projection_paths"`
	Relationships   map[string][]string `json:"relationships"`
	Provenance      map[string]any      `json:"provenance"`
}

type generatedProjection struct {
	Schema              string                    `json:"schema"`
	ProjectionID        string                    `json:"projection_id"`
	Surface             string                    `json:"surface"`
	Owner               string                    `json:"owner"`
	Status              string                    `json:"status"`
	AuthorityType       string                    `json:"authority_type"`
	Description         string                    `json:"description"`
	GeneratedFrom       []string                  `json:"generated_from"`
	SourceIndex         string                    `json:"source_index"`
	SemanticSource      string                    `json:"semantic_source"`
	Provenance          map[string]any            `json:"provenance"`
	RequiredDisclosures []string                  `json:"required_disclosures"`
	GenerationGate      []string                  `json:"generation_gate"`
	Units               []generatedProjectionUnit `json:"units"`
}

type generatedProjectionUnit struct {
	ID               string              `json:"id"`
	Kind             string              `json:"kind"`
	Locale           string              `json:"locale"`
	AuthorityType    string              `json:"authority_type"`
	SourceStatus     string              `json:"source_status"`
	Title            string              `json:"title"`
	Statement        string              `json:"statement,omitempty"`
	RetrievalSummary string              `json:"retrieval_summary,omitempty"`
	ProjectionPaths  map[string]string   `json:"projection_paths,omitempty"`
	Relationships    map[string][]string `json:"relationships,omitempty"`
	Tags             []string            `json:"tags,omitempty"`
	Provenance       map[string]any      `json:"provenance,omitempty"`
}

func main() {
	if len(os.Args) != 4 {
		fatal(errors.New("usage: go run scripts/generate_semantic_projections.go <manifest.json> <index.json> <source.json>"))
	}
	if err := generateProjections(os.Args[1], os.Args[2], os.Args[3]); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "semantic-projection generation failed: %v\n", err)
	os.Exit(1)
}

func generateProjections(manifestFile, indexFile, sourceFile string) error {
	var manifest generationManifest
	if err := readJSON(manifestFile, &manifest); err != nil {
		return err
	}
	var index generationIndex
	if err := readJSON(indexFile, &index); err != nil {
		return err
	}
	var source generationSource
	if err := readJSON(sourceFile, &source); err != nil {
		return err
	}

	indexUnits := map[string]generationIndexUnit{}
	for _, unit := range index.Units {
		indexUnits[unit.ID] = unit
	}
	sourceUnits := map[string]generationSourceUnit{}
	for _, unit := range source.Units {
		sourceUnits[unit.ID] = unit
	}

	for _, projection := range manifest.Projections {
		output := generatedProjection{
			Schema:              "rgb-docs-generated-projection/0.1",
			ProjectionID:        projection.ID,
			Surface:             projection.Surface,
			Owner:               projection.Owner,
			Status:              projection.Status,
			AuthorityType:       "generated_artifact",
			Description:         projection.Description,
			GeneratedFrom:       projection.SourceUnits,
			SourceIndex:         indexFile,
			SemanticSource:      sourceFile,
			Provenance:          projection.Provenance,
			RequiredDisclosures: projection.RequiredDisclosures,
			GenerationGate:      projection.GenerationGate,
			Units:               make([]generatedProjectionUnit, 0, len(projection.SourceUnits)),
		}
		for _, unitID := range projection.SourceUnits {
			indexUnit, ok := indexUnits[unitID]
			if !ok {
				return fmt.Errorf("%s: unknown source unit %s", projection.ID, unitID)
			}
			output.Units = append(output.Units, projectUnit(indexUnit, sourceUnits[unitID]))
		}
		if err := writeGeneratedProjection(projection.OutputPath, output); err != nil {
			return err
		}
		fmt.Printf("generated %s\n", projection.OutputPath)
	}

	return nil
}

func projectUnit(indexUnit generationIndexUnit, sourceUnit generationSourceUnit) generatedProjectionUnit {
	unit := generatedProjectionUnit{
		ID:               indexUnit.ID,
		Kind:             indexUnit.Kind,
		Locale:           indexUnit.Locale,
		AuthorityType:    indexUnit.AuthorityType,
		SourceStatus:     indexUnit.SourceStatus,
		Title:            indexUnit.Title,
		RetrievalSummary: indexUnit.Index.RetrievalSummary,
		ProjectionPaths:  indexUnit.ProjectionPaths,
		Relationships:    indexUnit.Relationships,
		Tags:             indexUnit.Index.Tags,
		Provenance:       indexUnit.Provenance,
	}
	if sourceUnit.ID != "" {
		unit.Statement = sourceUnit.Statement
		unit.AuthorityType = sourceUnit.AuthorityType
		unit.SourceStatus = sourceUnit.SourceStatus
		unit.Provenance = sourceUnit.Provenance
	}
	return unit
}

func readJSON(path string, target any) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, target); err != nil {
		return fmt.Errorf("%s: invalid JSON: %w", path, err)
	}
	return nil
}

func writeGeneratedProjection(path string, projection generatedProjection) error {
	bytes, err := json.MarshalIndent(projection, "", "  ")
	if err != nil {
		return err
	}
	bytes = append(bytes, '\n')
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, bytes, 0o644)
}
