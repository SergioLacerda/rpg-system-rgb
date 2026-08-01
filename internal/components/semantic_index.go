package components

import (
	"encoding/json"
	"fmt"
	"os"
)

// SemanticIndex is the decoded shape of docs/core/semantic/core-v2.index.json.
// It lives in the shared contract package (not in tooling) so any component
// boundary — tooling's own validator, bundles, or future consumers — can
// read the semantic index without importing a sibling component, which
// tests/architecture/architecture_test.go's TestComponentsDoNotImportSiblings
// forbids.
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

// LoadSemanticIndex reads and parses the semantic index file.
func LoadSemanticIndex(indexFile string) (SemanticIndex, error) {
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
