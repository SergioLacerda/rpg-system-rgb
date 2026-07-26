package tooling

import (
	"encoding/json"
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

// Generate produces the derived generated/*/*.json projection artifacts
// from the projection manifest, semantic index, and semantic source,
// migrated from scripts/generate_semantic_projections.go. Unlike the
// original script, output paths are resolved against repoRoot explicitly
// instead of the process working directory.
func Generate(repoRoot, manifestFile, indexFile, sourceFile string) error {
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

	for _, proj := range manifest.Projections {
		output := generatedProjection{
			Schema:              "rgb-docs-generated-projection/0.1",
			ProjectionID:        proj.ID,
			Surface:             proj.Surface,
			Owner:               proj.Owner,
			Status:              proj.Status,
			AuthorityType:       "generated_artifact",
			Description:         proj.Description,
			GeneratedFrom:       proj.SourceUnits,
			SourceIndex:         repoRelative(repoRoot, indexFile),
			SemanticSource:      repoRelative(repoRoot, sourceFile),
			Provenance:          proj.Provenance,
			RequiredDisclosures: proj.RequiredDisclosures,
			GenerationGate:      proj.GenerationGate,
			Units:               make([]generatedProjectionUnit, 0, len(proj.SourceUnits)),
		}
		for _, unitID := range proj.SourceUnits {
			indexUnit, ok := indexUnits[unitID]
			if !ok {
				return fmt.Errorf("%s: unknown source unit %s", proj.ID, unitID)
			}
			output.Units = append(output.Units, projectUnit(indexUnit, sourceUnits[unitID]))
		}
		outputPath := filepath.Join(repoRoot, filepath.FromSlash(proj.OutputPath))
		if err := writeGeneratedProjection(outputPath, output); err != nil {
			return err
		}
		fmt.Printf("generated %s\n", proj.OutputPath)
	}

	return nil
}

// GenerateDefault runs Generate against the standard
// docs/core/semantic/** manifest, index, and source paths under repoRoot.
func GenerateDefault(repoRoot string) error {
	semantic := filepath.Join(repoRoot, "docs", "core", "semantic")
	return Generate(
		repoRoot,
		filepath.Join(semantic, "projection-manifest.v0.1.json"),
		filepath.Join(semantic, "core-v2.index.json"),
		filepath.Join(semantic, "source", "core-v2-rules.v0.1.json"),
	)
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
	bytes, err := os.ReadFile(path) //nolint:gosec // G304: path is a caller-supplied source file, by design for this generator
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
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return err
	}
	return os.WriteFile(path, bytes, 0o644) //nolint:gosec // G306: generated JSON artifacts are intended to be readable, matching the existing generated/*.json files
}
